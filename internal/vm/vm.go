package vm

import (
	"math/rand"
	"miya/internal/memory"
	"miya/internal/screen"
	"time"

	"github.com/telroshan/go-sfml/v2/window"
)

type VirtualMachine struct {
	registers    Registers
	delayTimer   byte
	soundTimer   byte
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

func NewVirtualMachine(memory *memory.Memory, stack *memory.Stack, screen *screen.Screen) *VirtualMachine {
	vm := VirtualMachine{
		registers: Registers{
			PC: 0x200,
			V:  make([]byte, 0x10),
		},
		delayTimer:   0,
		soundTimer:   0,
		memory:       memory,
		stack:        stack,
		screen:       screen,
		instructions: make(map[uint16]func(uint16)),
		keys:         make([]byte, 0x10),
		keypressed:   make(chan byte),
	}

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
	vm.delayTimer = 0
	vm.soundTimer = 0

	vm.memory.Reset()
	vm.stack.Reset()
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

		time.Sleep(time.Second / 60)
	}
}

func (vm *VirtualMachine) keypad() {
	for {
		key := byte(0x00)

		select {
		case keycode := <-vm.screen.Keyevt:
			switch keycode {
			case window.SfKeyCode(window.SfKeyNum1):
				vm.keys[K01] = vm.keys[K01] ^ 1
				key = K01
			case window.SfKeyCode(window.SfKeyNum2):
				vm.keys[K02] = vm.keys[K02] ^ 1
				key = K02
			case window.SfKeyCode(window.SfKeyNum3):
				vm.keys[K03] = vm.keys[K03] ^ 1
				key = K03
			case window.SfKeyCode(window.SfKeyNum4):
				vm.keys[K0C] = vm.keys[K0C] ^ 1
				key = K0C
			case window.SfKeyCode(window.SfKeyQ):
				vm.keys[K04] = vm.keys[K04] ^ 1
				key = K04
			case window.SfKeyCode(window.SfKeyW):
				vm.keys[K05] = vm.keys[K05] ^ 1
				key = K05
			case window.SfKeyCode(window.SfKeyE):
				vm.keys[K06] = vm.keys[K06] ^ 1
				key = K06
			case window.SfKeyCode(window.SfKeyR):
				vm.keys[K0D] = vm.keys[K0D] ^ 1
				key = K0D
			case window.SfKeyCode(window.SfKeyA):
				vm.keys[K07] = vm.keys[K07] ^ 1
				key = K07
			case window.SfKeyCode(window.SfKeyS):
				vm.keys[K08] = vm.keys[K08] ^ 1
				key = K08
			case window.SfKeyCode(window.SfKeyD):
				vm.keys[K09] = vm.keys[K09] ^ 1
				key = K09
			case window.SfKeyCode(window.SfKeyF):
				vm.keys[K0E] = vm.keys[K0E] ^ 1
				key = K0E
			case window.SfKeyCode(window.SfKeyZ):
				vm.keys[K0A] = vm.keys[K0A] ^ 1
				key = K0A
			case window.SfKeyCode(window.SfKeyX):
				vm.keys[K00] = vm.keys[K00] ^ 1
				key = K00
			case window.SfKeyCode(window.SfKeyC):
				vm.keys[K0B] = vm.keys[K0B] ^ 1
				key = K0B
			case window.SfKeyCode(window.SfKeyV):
				vm.keys[K0F] = vm.keys[K0F] ^ 1
				key = K0F
			}
		}

		if vm.waitforkey {
			vm.keypressed <- key
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
		if vm.registers.V[x] < vm.registers.V[y] {
			vm.registers.V[0x0F] = 1
		} else {
			vm.registers.V[0x0F] = 0
		}

		vm.registers.V[x] -= vm.registers.V[y]
	case 6:
		vm.registers.V[0x0F] = (vm.registers.V[x] & 0x01)
		vm.registers.V[x] >>= 1
	case 7:
		if vm.registers.V[y] < vm.registers.V[x] {
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

				vm.screen.SetPixel(x+byte(k), (y + byte(i)))
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
