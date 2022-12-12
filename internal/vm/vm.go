package vm

import (
	"math/rand"
	"miya/internal/memory"
	"miya/internal/screen"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type VirtualMachine struct {
	registers    Registers
	delayTimer   byte
	soundTimer   byte
	delay        uint64
	memory       *memory.Memory
	stack        *memory.Stack
	screen       *screen.Screen
	instructions map[uint16]func(uint16)
	keys         []byte
	keypressed   chan byte
	waitforkey   bool
}

type Registers struct {
	I  uint16
	PC uint16
	V  []byte
}

var font = []byte{0xF0, 0x90, 0x90, 0x90, 0xF0,
	0x20, 0x60, 0x20, 0x20, 0x70,
	0xF0, 0x10, 0xF0, 0x80, 0xF0,
	0xF0, 0x10, 0xF0, 0x10, 0xF0,
	0x90, 0x90, 0xF0, 0x10, 0x10,
	0xF0, 0x80, 0xF0, 0x10, 0xF0,
	0xF0, 0x80, 0xF0, 0x90, 0xF0,
	0xF0, 0x10, 0x20, 0x40, 0x40,
	0xF0, 0x90, 0xF0, 0x90, 0xF0,
	0xF0, 0x90, 0xF0, 0x10, 0xF0,
	0xF0, 0x90, 0xF0, 0x90, 0x90,
	0xE0, 0x90, 0xE0, 0x90, 0xE0,
	0xF0, 0x80, 0x80, 0x80, 0xF0,
	0xE0, 0x90, 0x90, 0x90, 0xE0,
	0xF0, 0x80, 0xF0, 0x80, 0xF0,
	0xF0, 0x80, 0xF0, 0x80, 0x80}

func NewVirtualMachine(memory *memory.Memory, stack *memory.Stack, screen *screen.Screen, delay uint64) *VirtualMachine {
	vm := VirtualMachine{
		registers: Registers{
			PC: 0x200,
			V:  make([]byte, 0x10),
		},
		delayTimer:   0,
		soundTimer:   0,
		delay:        delay,
		memory:       memory,
		stack:        stack,
		screen:       screen,
		instructions: make(map[uint16]func(uint16)),
		keys:         make([]byte, 0x10),
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

func (vm *VirtualMachine) Reset() {
	vm.registers.PC = 0x200
	vm.registers.V = make([]byte, 0x10)
	vm.keys = make([]byte, 0x10)
	vm.delayTimer = 0
	vm.soundTimer = 0

	vm.memory.Reset()
	vm.stack.Reset()
	vm.screen.Clear()
	vm.memory.WriteArray(0x000, font)
}

func (vm *VirtualMachine) EvalLoop() {
	go vm.keypad()

	for {
		opcode := vm.memory.ReadOpcode(vm.registers.PC)
		vm.instructions[(opcode & 0xF000)](opcode)

		if vm.delayTimer > 0 {
			vm.delayTimer--
		}

		if vm.soundTimer > 0 {
			if vm.soundTimer == 1 {
				println("mmm....BEEP!")
			}
			vm.soundTimer--
		}

		time.Sleep(time.Millisecond * time.Duration(vm.delay))
	}
}

func (vm *VirtualMachine) keypad() {
	for {
		keyevent := <-vm.screen.Keyevt
		if _, ok := keymap[keyevent.Keycode]; ok {
			if keyevent.Etype == sdl.KEYUP {
				vm.keys[keymap[keyevent.Keycode]] = 0
			} else {
				vm.keys[keymap[keyevent.Keycode]] = 1

				if vm.waitforkey {
					vm.keypressed <- keymap[keyevent.Keycode]
				}
			}
		}
	}
}

func (vm *VirtualMachine) clc(opcode uint16) {
	otype := (opcode & 0x0FFF)

	if otype == 0x0E0 {
		vm.screen.Clear()
		vm.registers.PC += 2

		return
	}

	if otype == 0x0EE {
		vm.registers.PC = vm.stack.Pop()
		vm.registers.PC += 2
	}
}

func (vm *VirtualMachine) jp(opcode uint16) {
	vm.registers.PC = (opcode & 0x0FFF)
}

func (vm *VirtualMachine) call(opcode uint16) {
	vm.stack.Push(vm.registers.PC)
	vm.registers.PC = (opcode & 0x0FFF)
}

func (vm *VirtualMachine) sevx(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	if vm.registers.V[x] == byte(nn) {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) sne(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	if vm.registers.V[x] != byte(nn) {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) sevxvy(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4

	if vm.registers.V[x] == vm.registers.V[y] {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldvx(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	vm.registers.V[x] = byte(nn)
	vm.registers.PC += 2
}

func (vm *VirtualMachine) add(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	vm.registers.V[x] += byte(nn)
	vm.registers.PC += 2
}

func (vm *VirtualMachine) vxvy(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4

	switch opcode & 0x000F {
	case 0:
		vm.registers.V[x] = vm.registers.V[y]
	case 1:
		vm.registers.V[x] |= vm.registers.V[y]
	case 2:
		vm.registers.V[x] &= vm.registers.V[y]
	case 3:
		vm.registers.V[x] ^= vm.registers.V[y]
	case 4:
		if (uint16(vm.registers.V[x]) + uint16(vm.registers.V[y])) > 0xFF {
			vm.registers.V[0x0F] = 1
		} else {
			vm.registers.V[0x0F] = 0
		}

		vm.registers.V[x] += vm.registers.V[y]
	case 5:
		if vm.registers.V[x] > vm.registers.V[y] {
			vm.registers.V[0x0F] = 1
		} else {
			vm.registers.V[0x0F] = 0
		}

		vm.registers.V[x] -= vm.registers.V[y]
	case 6:
		vm.registers.V[0x0F] = (vm.registers.V[x] & 0x01)
		vm.registers.V[x] >>= 1
	case 7:
		if vm.registers.V[y] > vm.registers.V[x] {
			vm.registers.V[0x0F] = 1
		} else {
			vm.registers.V[0x0F] = 0
		}

		vm.registers.V[x] = vm.registers.V[y] - vm.registers.V[x]
	case 0xe:
		vm.registers.V[0x0F] = (vm.registers.V[x] & 0x80)
		vm.registers.V[x] <<= 1
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) snevxvy(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4

	if vm.registers.V[x] != vm.registers.V[y] {
		vm.registers.PC += 4
		return
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldi(opcode uint16) {
	vm.registers.I = (opcode & 0x0FFF)
	vm.registers.PC += 2
}

func (vm *VirtualMachine) jpv0(opcode uint16) {
	vm.registers.PC = uint16(vm.registers.V[0]) + (opcode & 0x0FFF)
}

func (vm *VirtualMachine) rnd(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	nn := (opcode & 0x00FF)

	rand.Seed(time.Now().UnixNano())
	vm.registers.V[x] = byte(rand.Intn(0xFF)) & byte(nn)
	vm.registers.PC += 2
}

func (vm *VirtualMachine) drw(opcode uint16) {
	x := vm.registers.V[(opcode&0x0F00)>>8]
	y := vm.registers.V[(opcode&0x00F0)>>4]
	h := (opcode & 0x000F)

	vm.registers.V[0x0F] = 0

	for i := uint16(0); i < h; i++ {
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

func (vm *VirtualMachine) skp(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	otype := (opcode & 0x00FF)

	if otype == 0x9E {
		if vm.keys[x] == 1 {
			vm.registers.PC += 4
			return
		}
	}

	if otype == 0xA1 {
		if vm.keys[x] == 0 {
			vm.registers.PC += 4
			return
		}
	}

	vm.registers.PC += 2
}

func (vm *VirtualMachine) ldf(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	otype := (opcode & 0x00FF)

	switch otype {
	case 0x07:
		vm.registers.V[x] = vm.delayTimer
	case 0x0A:
		vm.waitforkey = true
		vm.registers.V[x] = <-vm.keypressed
		vm.waitforkey = false
	case 0x15:
		vm.delayTimer = vm.registers.V[x]
	case 0x18:
		vm.soundTimer = vm.registers.V[x]
	case 0x1E:
		vm.registers.I += uint16(vm.registers.V[x])
	case 0x29:
		vm.registers.I = uint16(vm.registers.V[x] * 0x05)
	case 0x33:
		n := vm.registers.V[x]
		vm.memory.Write(vm.registers.I, n/100)
		vm.memory.Write(vm.registers.I+1, (n/10)%10)
		vm.memory.Write(vm.registers.I+2, (n%100)%10)
	case 0x55:
		for i := uint16(0); i <= x; i++ {
			vm.memory.Write(vm.registers.I, vm.registers.V[i])
			vm.registers.I += 1
		}
	case 0x65:
		for i := uint16(0); i <= x; i++ {
			vm.registers.V[i] = vm.memory.Read(vm.registers.I)
			vm.registers.I += 1
		}
	}

	vm.registers.PC += 2
}
