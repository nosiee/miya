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
	return &VirtualMachine{
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
		select {
		case keycode := <-vm.screen.Keyevt:
			switch keycode {
			case window.SfKeyCode(window.SfKeyNum1):
				vm.keys[0x00] = vm.keys[0x00] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x01
				}
			case window.SfKeyCode(window.SfKeyNum2):
				vm.keys[0x01] = vm.keys[0x01] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x02
				}
			case window.SfKeyCode(window.SfKeyNum3):
				vm.keys[0x02] = vm.keys[0x02] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x03
				}
			case window.SfKeyCode(window.SfKeyNum4):
				vm.keys[0x03] = vm.keys[0x03] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x0C
				}
			case window.SfKeyCode(window.SfKeyQ):
				vm.keys[0x04] = vm.keys[0x04] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x4
				}
			case window.SfKeyCode(window.SfKeyW):
				vm.keys[0x05] = vm.keys[0x05] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x05
				}
			case window.SfKeyCode(window.SfKeyE):
				vm.keys[0x06] = vm.keys[0x06] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x06
				}
			case window.SfKeyCode(window.SfKeyR):
				vm.keys[0x07] = vm.keys[0x07] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x0D
				}
			case window.SfKeyCode(window.SfKeyA):
				vm.keys[0x08] = vm.keys[0x08] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x07
				}
			case window.SfKeyCode(window.SfKeyS):
				vm.keys[0x09] = vm.keys[0x09] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x08
				}
			case window.SfKeyCode(window.SfKeyD):
				vm.keys[0xA] = vm.keys[0xA] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x09
				}
			case window.SfKeyCode(window.SfKeyF):
				vm.keys[0xB] = vm.keys[0xB] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x0E
				}
			case window.SfKeyCode(window.SfKeyZ):
				vm.keys[0xC] = vm.keys[0xC] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x0A
				}
			case window.SfKeyCode(window.SfKeyX):
				vm.keys[0xD] = vm.keys[0xD] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x00
				}
			case window.SfKeyCode(window.SfKeyC):
				vm.keys[0xE] = vm.keys[0xE] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x0B
				}
			case window.SfKeyCode(window.SfKeyV):
				vm.keys[0xF] = vm.keys[0xF] ^ 1
				if vm.waitforkey {
					vm.keypressed <- 0x0F
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
		rs := vm.registers.V[x] + vm.registers.V[y]
		if rs > 0xFF {
			vm.registers.V[0x0F] = 1
		} else {
			vm.registers.V[0x0F] = 0
		}

		vm.registers.V[x] = rs
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

		vm.registers.V[x] = vm.registers.V[y] - vm.registers.V[y]
	case 0xe:
		vm.registers.V[0x0F] = (vm.registers.V[x] & 128)
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
