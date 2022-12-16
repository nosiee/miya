package vm

import (
	"math/rand"
	"miya/internal/memory"
	"miya/internal/screen"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func NewVirtualMachine(memory *memory.Memory, stack *memory.Stack, screen *screen.Screen, delay uint64) *VirtualMachine {
	vm := VirtualMachine{
		Registers: Registers{
			PC: 0x200,
			V:  make([]byte, 0x10),
		},
		DelayTimer:   0,
		SoundTimer:   0,
		delay:        delay,
		memory:       memory,
		stack:        stack,
		screen:       screen,
		instructions: make(map[uint16]func(Opcode)),
		Keys:         make([]byte, 0x10),
		keypressed:   make(chan byte),
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

func NewOpcode(value uint16) Opcode {
	return Opcode{
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
	vm.Registers.PC = 0x200
	vm.Registers.I = 0x000
	vm.Registers.V = make([]byte, 0x10)
	vm.Keys = make([]byte, 0x10)
	vm.DelayTimer = 0
	vm.SoundTimer = 0

	vm.memory.Reset()
	vm.stack.Reset()
	vm.screen.Clear()
	vm.memory.WriteArray(0x000, font)
}

func (vm *VirtualMachine) EvalLoop() {
	go vm.keypad()

	for {
		opcode := NewOpcode(vm.memory.ReadOpcode(vm.Registers.PC))
		vm.instructions[opcode.t](opcode)

		if vm.DelayTimer > 0 {
			vm.DelayTimer--
		}

		if vm.SoundTimer > 0 {
			if vm.SoundTimer == 1 {
				println("mmm....BEEP!")
			}
			vm.SoundTimer--
		}

		time.Sleep(time.Millisecond * time.Duration(vm.delay))
	}
}

func (vm *VirtualMachine) keypad() {
	for {
		keyevent := <-vm.screen.Keyevt

		if _, ok := keymap[keyevent.Keycode]; ok {
			if keyevent.Etype == sdl.KEYUP {
				vm.Keys[keymap[keyevent.Keycode]] = 0
			} else {
				vm.Keys[keymap[keyevent.Keycode]] = 1

				if vm.waitforkey {
					vm.keypressed <- keymap[keyevent.Keycode]
				}
			}
		}
	}
}

func (vm *VirtualMachine) clc(opcode Opcode) {
	if opcode.nnn == 0x0E0 {
		vm.screen.Clear()
		vm.Registers.PC += 2

		return
	}

	if opcode.nnn == 0x0EE {
		vm.Registers.PC = vm.stack.Pop()
		vm.Registers.PC += 2
	}
}

func (vm *VirtualMachine) jp(opcode Opcode) {
	vm.Registers.PC = opcode.nnn
}

func (vm *VirtualMachine) call(opcode Opcode) {
	vm.stack.Push(vm.Registers.PC)
	vm.Registers.PC = opcode.nnn
}

func (vm *VirtualMachine) sevx(opcode Opcode) {
	if vm.Registers.V[opcode.x] == opcode.nn {
		vm.Registers.PC += 4
		return
	}

	vm.Registers.PC += 2
}

func (vm *VirtualMachine) sne(opcode Opcode) {
	if vm.Registers.V[opcode.x] != opcode.nn {
		vm.Registers.PC += 4
		return
	}

	vm.Registers.PC += 2
}

func (vm *VirtualMachine) sevxvy(opcode Opcode) {
	if vm.Registers.V[opcode.x] == vm.Registers.V[opcode.y] {
		vm.Registers.PC += 4
		return
	}

	vm.Registers.PC += 2
}

func (vm *VirtualMachine) ldvx(opcode Opcode) {
	vm.Registers.V[opcode.x] = opcode.nn
	vm.Registers.PC += 2
}

func (vm *VirtualMachine) add(opcode Opcode) {
	vm.Registers.V[opcode.x] += opcode.nn
	vm.Registers.PC += 2
}

func (vm *VirtualMachine) vxvy(opcode Opcode) {
	switch opcode.n {
	case 0:
		vm.Registers.V[opcode.x] = vm.Registers.V[opcode.y]
	case 1:
		vm.Registers.V[opcode.x] |= vm.Registers.V[opcode.y]
	case 2:
		vm.Registers.V[opcode.x] &= vm.Registers.V[opcode.y]
	case 3:
		vm.Registers.V[opcode.x] ^= vm.Registers.V[opcode.y]
	case 4:
		if (uint16(vm.Registers.V[opcode.x]) + uint16(vm.Registers.V[opcode.y])) > 0xFF {
			vm.Registers.V[0x0F] = 1
		} else {
			vm.Registers.V[0x0F] = 0
		}

		vm.Registers.V[opcode.x] += vm.Registers.V[opcode.y]
	case 5:
		if vm.Registers.V[opcode.x] > vm.Registers.V[opcode.y] {
			vm.Registers.V[0x0F] = 1
		} else {
			vm.Registers.V[0x0F] = 0
		}

		vm.Registers.V[opcode.x] -= vm.Registers.V[opcode.y]
	case 6:
		vm.Registers.V[0x0F] = (vm.Registers.V[opcode.x] & 0x01)
		vm.Registers.V[opcode.x] >>= 1
	case 7:
		if vm.Registers.V[opcode.y] > vm.Registers.V[opcode.x] {
			vm.Registers.V[0x0F] = 1
		} else {
			vm.Registers.V[0x0F] = 0
		}

		vm.Registers.V[opcode.x] = vm.Registers.V[opcode.y] - vm.Registers.V[opcode.x]
	case 0xe:
		vm.Registers.V[0x0F] = (vm.Registers.V[opcode.x] & 0x80)
		vm.Registers.V[opcode.x] <<= 1
	}

	vm.Registers.PC += 2
}

func (vm *VirtualMachine) snevxvy(opcode Opcode) {
	if vm.Registers.V[opcode.x] != vm.Registers.V[opcode.y] {
		vm.Registers.PC += 4
		return
	}

	vm.Registers.PC += 2
}

func (vm *VirtualMachine) ldi(opcode Opcode) {
	vm.Registers.I = opcode.nnn
	vm.Registers.PC += 2
}

func (vm *VirtualMachine) jpv0(opcode Opcode) {
	vm.Registers.PC = uint16(vm.Registers.V[0]) + opcode.nnn
}

func (vm *VirtualMachine) rnd(opcode Opcode) {
	rand.Seed(time.Now().UnixNano())
	vm.Registers.V[opcode.x] = byte(rand.Intn(0xFF)) & opcode.nn
	vm.Registers.PC += 2
}

func (vm *VirtualMachine) drw(opcode Opcode) {
	x := vm.Registers.V[opcode.x]
	y := vm.Registers.V[opcode.y]

	vm.Registers.V[0x0F] = 0

	for i := uint16(0); i < uint16(opcode.n); i++ {
		pixel := vm.memory.Read(vm.Registers.I + i)
		for k := 0; k < 8; k++ {
			if pixel&(0x80>>k) != 0 {
				if vm.screen.GetPixel(x+byte(k), y+byte(i)) == 1 {
					vm.Registers.V[0x0F] = 1
				}

				vm.screen.SetPixel(x+byte(k), y+byte(i))
			}
		}
	}

	vm.Registers.PC += 2
}

func (vm *VirtualMachine) skp(opcode Opcode) {
	if opcode.nn == 0x9E {
		if vm.Keys[vm.Registers.V[opcode.x]] == 1 {
			vm.Registers.PC += 4
			return
		}
	}

	if opcode.nn == 0xA1 {
		if vm.Keys[vm.Registers.V[opcode.x]] == 0 {
			vm.Registers.PC += 4
			return
		}
	}

	vm.Registers.PC += 2
}

func (vm *VirtualMachine) ldf(opcode Opcode) {
	switch opcode.nn {
	case 0x07:
		vm.Registers.V[opcode.x] = vm.DelayTimer
	case 0x0A:
		vm.waitforkey = true
		vm.Registers.V[opcode.x] = <-vm.keypressed
		vm.waitforkey = false
	case 0x15:
		vm.DelayTimer = vm.Registers.V[opcode.x]
	case 0x18:
		vm.SoundTimer = vm.Registers.V[opcode.x]
	case 0x1E:
		vm.Registers.I += uint16(vm.Registers.V[opcode.x])
	case 0x29:
		vm.Registers.I = uint16(vm.Registers.V[opcode.x] * 0x05)
	case 0x33:
		n := vm.Registers.V[opcode.x]
		vm.memory.Write(vm.Registers.I, n/100)
		vm.memory.Write(vm.Registers.I+1, (n/10)%10)
		vm.memory.Write(vm.Registers.I+2, (n%100)%10)
	case 0x55:
		for i := byte(0); i <= opcode.x; i++ {
			vm.memory.Write(vm.Registers.I, vm.Registers.V[i])
			vm.Registers.I += 1
		}
	case 0x65:
		for i := byte(0); i <= opcode.x; i++ {
			vm.Registers.V[i] = vm.memory.Read(vm.Registers.I)
			vm.Registers.I += 1
		}
	}

	vm.Registers.PC += 2
}
