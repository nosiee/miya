package vm

import (
	"fmt"
	"math/rand"
	"miya/internal/memory"
	"miya/internal/screen"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func NewVirtualMachine(memory *memory.Memory, stack *memory.Stack, screen screen.Chip8Screen, delay uint64, debugMode bool) *VirtualMachine {
	vm := VirtualMachine{
		registers: registers{
			PC: 0x200,
			V:  make([]byte, 0x10),
		},
		delayTimer:   0,
		soundTimer:   0,
		delay:        delay,
		memory:       memory,
		stack:        stack,
		screen:       screen,
		instructions: make(map[uint16]func(opcode)),
		keys:         make([]byte, 0x10),
		keyPressed:   make(chan byte),
		debugMode:    debugMode,
	}

	vm.memory.WriteArray(0x000, font)

	vm.instructions[CLC] = vm.clc
	vm.instructions[JP] = vm.jp
	vm.instructions[CALL] = vm.call
	vm.instructions[SE_VX] = vm.sevx
	vm.instructions[SNE] = vm.sne
	vm.instructions[SE_VX_VY] = vm.sevxvy
	vm.instructions[LD_VX] = vm.ldvx
	vm.instructions[ADD] = vm.add
	vm.instructions[VX_VY] = vm.vxvy
	vm.instructions[SNE_VX_VY] = vm.snevxvy
	vm.instructions[LD_I] = vm.ldi
	vm.instructions[JP_V0] = vm.jpv0
	vm.instructions[RND] = vm.rnd
	vm.instructions[DRW] = vm.drw
	vm.instructions[SKP] = vm.skp
	vm.instructions[LDF] = vm.ldf

	return &vm
}

func newOpcode(value uint16) opcode {
	return opcode{
		value: value,
		t:     value & 0xF000,
		x:     byte((value & 0x0F00) >> 8),
		y:     byte((value & 0x00F0) >> 4),
		n:     byte(value & 0x000F),
		nn:    byte(value & 0x00FF),
		nnn:   value & 0x0FFF,
	}
}

func (vm *VirtualMachine) Reset() {
	vm.registers.PC = 0x200
	vm.registers.I = 0x000
	vm.registers.V = make([]byte, 0x10)
	vm.keys = make([]byte, 0x10)
	vm.delayTimer = 0
	vm.soundTimer = 0

	vm.memory.Reset()
	vm.stack.Reset()
	vm.screen.Clear()
	vm.memory.WriteArray(0x000, font)
}

func (vm *VirtualMachine) Debug() {
	for {
		opcode := newOpcode(vm.memory.ReadOpcode(vm.registers.PC))
		opcodeInfo := fmt.Sprintf("0x%04x [t: 0x%04x, x: 0x%02x, y: 0x%02x, n: 0x%02x, nn: 0x%02x, nnn: 0x%04x]",
			opcode.value,
			opcode.t,
			opcode.x,
			opcode.y,
			opcode.n,
			opcode.nn,
			opcode.nnn)

		screen.Debug <- fmt.Sprintf("Opcode: %s\nI: 0x%04x\nPC: 0x%04x\nVX: %v\nDelayTimer: %d\nsoundTimer: %d\nKeys: %v\nStack: %v", opcodeInfo, vm.registers.I, vm.registers.PC, vm.registers.V, vm.delayTimer, vm.soundTimer, vm.keys, vm.stack.Dump())
	}
}

func (vm *VirtualMachine) EvalLoop() {
	go vm.keypad()

	for {
		if vm.debugMode {
			<-screen.Next
		}

		opcode := newOpcode(vm.memory.ReadOpcode(vm.registers.PC))
		vm.instructions[opcode.t](opcode)

		if vm.delayTimer > 0 {
			vm.delayTimer--
		}

		if vm.soundTimer > 0 {
			vm.soundTimer--
		}

		time.Sleep(time.Millisecond * time.Duration(vm.delay))
	}
}

func (vm *VirtualMachine) keypad() {
	for {
		keyevent := <-screen.KeyPressed

		if _, ok := keymap[keyevent.Keycode]; ok {
			if keyevent.Etype == sdl.KEYUP {
				vm.keys[keymap[keyevent.Keycode]] = 0
			} else {
				vm.keys[keymap[keyevent.Keycode]] = 1

				if vm.waitForKey {
					vm.keyPressed <- keymap[keyevent.Keycode]
				}
			}
		}
	}
}

func (vm *VirtualMachine) clc(opcode opcode) {
	if opcode.nnn == 0x0E0 {
		vm.screen.Clear()
		vm.registers.PC += 2

		return
	}

	if opcode.nnn == 0x0EE {
		vm.registers.PC = vm.stack.Pop()
		vm.registers.PC += 2
	}
}

func (vm *VirtualMachine) jp(opcode opcode) {
	vm.registers.PC = opcode.nnn
}

func (vm *VirtualMachine) call(opcode opcode) {
	vm.stack.Push(vm.registers.PC)
	vm.registers.PC = opcode.nnn
}

func (vm *VirtualMachine) sevx(opcode opcode) {
	if vm.registers.V[opcode.x] == opcode.nn {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) sne(opcode opcode) {
	if vm.registers.V[opcode.x] != opcode.nn {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) sevxvy(opcode opcode) {
	if vm.registers.V[opcode.x] == vm.registers.V[opcode.y] {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldvx(opcode opcode) {
	vm.registers.V[opcode.x] = opcode.nn
	vm.registers.PC += 2
}

func (vm *VirtualMachine) add(opcode opcode) {
	vm.registers.V[opcode.x] += opcode.nn
	vm.registers.PC += 2
}

func (vm *VirtualMachine) vxvy(opcode opcode) {
	switch opcode.n {
	case 0:
		vm.registers.V[opcode.x] = vm.registers.V[opcode.y]
	case 1:
		vm.registers.V[opcode.x] |= vm.registers.V[opcode.y]
	case 2:
		vm.registers.V[opcode.x] &= vm.registers.V[opcode.y]
	case 3:
		vm.registers.V[opcode.x] ^= vm.registers.V[opcode.y]
	case 4:
		if (uint16(vm.registers.V[opcode.x]) + uint16(vm.registers.V[opcode.y])) > 0xFF {
			vm.registers.V[0x0F] = 1
		} else {
			vm.registers.V[0x0F] = 0
		}

		vm.registers.V[opcode.x] += vm.registers.V[opcode.y]
	case 5:
		if vm.registers.V[opcode.x] > vm.registers.V[opcode.y] {
			vm.registers.V[0x0F] = 1
		} else {
			vm.registers.V[0x0F] = 0
		}

		vm.registers.V[opcode.x] -= vm.registers.V[opcode.y]
	case 6:
		vm.registers.V[0x0F] = (vm.registers.V[opcode.x] & 0x01)
		vm.registers.V[opcode.x] >>= 1
	case 7:
		if vm.registers.V[opcode.y] > vm.registers.V[opcode.x] {
			vm.registers.V[0x0F] = 1
		} else {
			vm.registers.V[0x0F] = 0
		}

		vm.registers.V[opcode.x] = vm.registers.V[opcode.y] - vm.registers.V[opcode.x]
	case 0xe:
		vm.registers.V[0x0F] = (vm.registers.V[opcode.x] & 0x80)
		vm.registers.V[opcode.x] <<= 1
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) snevxvy(opcode opcode) {
	if vm.registers.V[opcode.x] != vm.registers.V[opcode.y] {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldi(opcode opcode) {
	vm.registers.I = opcode.nnn
	vm.registers.PC += 2
}

func (vm *VirtualMachine) jpv0(opcode opcode) {
	vm.registers.PC = uint16(vm.registers.V[0]) + opcode.nnn
}

func (vm *VirtualMachine) rnd(opcode opcode) {
	rand.Seed(time.Now().UnixNano())
	vm.registers.V[opcode.x] = byte(rand.Intn(0xFF)) & opcode.nn
	vm.registers.PC += 2
}

func (vm *VirtualMachine) drw(opcode opcode) {
	x := vm.registers.V[opcode.x]
	y := vm.registers.V[opcode.y]

	vm.registers.V[0x0F] = 0

	for i := uint16(0); i < uint16(opcode.n); i++ {
		pixel := vm.memory.Read(vm.registers.I + i)
		for k := 0; k < 8; k++ {
			if pixel&(0x80>>k) != 0 {
				if vm.screen.GetPixel(x+byte(k), y+byte(i)) == 1 {
					vm.registers.V[0x0F] = 1
				}

				vm.screen.SetPixel(x+byte(k), y+byte(i))
			}
		}
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) skp(opcode opcode) {
	if opcode.nn == 0x9E {
		if vm.keys[vm.registers.V[opcode.x]] == 1 {
			vm.registers.PC += 4
			return
		}
	}

	if opcode.nn == 0xA1 {
		if vm.keys[vm.registers.V[opcode.x]] == 0 {
			vm.registers.PC += 4
			return
		}
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldf(opcode opcode) {
	switch opcode.nn {
	case 0x07:
		vm.registers.V[opcode.x] = vm.delayTimer
	case 0x0A:
		vm.waitForKey = true
		vm.registers.V[opcode.x] = <-vm.keyPressed
		vm.waitForKey = false
	case 0x15:
		vm.delayTimer = vm.registers.V[opcode.x]
	case 0x18:
		vm.soundTimer = vm.registers.V[opcode.x]
	case 0x1E:
		vm.registers.I += uint16(vm.registers.V[opcode.x])
	case 0x29:
		vm.registers.I = uint16(vm.registers.V[opcode.x] * 0x05)
	case 0x33:
		n := vm.registers.V[opcode.x]
		vm.memory.Write(vm.registers.I, n/100)
		vm.memory.Write(vm.registers.I+1, (n/10)%10)
		vm.memory.Write(vm.registers.I+2, (n%100)%10)
	case 0x55:
		for i := byte(0); i <= opcode.x; i++ {
			vm.memory.Write(vm.registers.I, vm.registers.V[i])
			vm.registers.I += 1
		}
	case 0x65:
		for i := byte(0); i <= opcode.x; i++ {
			vm.registers.V[i] = vm.memory.Read(vm.registers.I)
			vm.registers.I += 1
		}
	}

	vm.registers.PC += 2
}
