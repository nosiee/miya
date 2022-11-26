package vm

import (
	"fmt"
	"math/rand"
	"miya/internal/memory"
	"miya/internal/screen"
	"time"
)

type VirtualMachine struct {
	registers Registers
	memory    *memory.Memory
	stack     *memory.Stack
	screen    *screen.Screen
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
		memory: memory,
		stack:  stack,
		screen: screen,
	}
}

func (vm *VirtualMachine) EvalLoop() {
	for {
		opcode := vm.memory.ReadOpcode(vm.registers.PC)

		switch opcode & 0xF000 {
		case 0x0000:
			otype := (opcode & 0x0FFF)

			if otype == 0x0E0 {
				vm.screen.Clear()
			}

			if otype == 0x0EE {
				vm.registers.PC = vm.stack.Pop()
				continue
			}
		case 0x1000:
			vm.registers.PC = (opcode & 0x0FFF)
			continue
		case 0x2000:
			vm.stack.Push(vm.registers.PC)
			vm.registers.PC = (opcode & 0x0FFF)
			continue
		case 0x3000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			if vm.registers.V[x] == byte(nn) {
				vm.registers.PC += 2
			}
		case 0x4000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			if vm.registers.V[x] != byte(nn) {
				vm.registers.PC += 2
			}
		case 0x5000:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4

			if vm.registers.V[x] == vm.registers.V[y] {
				vm.registers.PC += 2
			}
		case 0x6000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			vm.registers.V[x] = byte(nn)
		case 0x7000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			vm.registers.V[x] += byte(nn)
		case 0x8000:
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
				vm.registers.V[x] += vm.registers.V[y]
			case 5:
				vm.registers.V[x] -= vm.registers.V[y]
			case 6:
				vm.registers.V[x] >>= 1
			case 7:
				vm.registers.V[x] = vm.registers.V[y] - vm.registers.V[x]
			case 0xe:
				vm.registers.V[x] <<= 1
			}
		case 0x9000:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4

			if vm.registers.V[x] != vm.registers.V[y] {
				vm.registers.PC += 2
			}
		case 0xA000:
			vm.registers.I = (opcode & 0x0FFF)
		case 0xB000:
			vm.registers.PC = (uint16(vm.registers.V[0]) + (opcode & 0x0FFF))
		case 0xC000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			rand.Seed(time.Now().UnixNano())
			vm.registers.V[x] = byte((rand.Intn(0xff-0x00) + 0x00) & int(nn))
		case 0xD000:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			n := (opcode & 0x000F)

			fmt.Println(x, y, n)
		case 0xE000:
			println("0xe000 TODO")
		case 0xF000:
			println("0xf000 TODO")
		default:
			fmt.Printf("unrecognized opcode '0x%04x' at memory[0x%04x]\n", (opcode & 0xF000), vm.registers.PC)
		}

		vm.registers.PC += 2
		time.Sleep(time.Second / 60)
	}
}
