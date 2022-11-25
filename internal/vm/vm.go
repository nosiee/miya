package vm

import (
	"fmt"
	"miya/internal/memory"
	"time"
)

type VirtualMachine struct {
	Registers Registers
	Memory    *memory.Memory
	Stack     *memory.Stack
}

type Registers struct {
	I  uint16
	PC uint16
	V  []byte
}

func NewVirtualMachine(memory *memory.Memory, stack *memory.Stack) *VirtualMachine {
	return &VirtualMachine{
		Registers: Registers{
			PC: 0x200,
			V:  make([]byte, 0x10),
		},
		Memory: memory,
		Stack:  stack,
	}
}

func (vm *VirtualMachine) EvalLoop() {
	for {
		opcode := vm.Memory.ReadOpcode(vm.Registers.PC)

		switch opcode & 0xF000 {
		case 0x0000:
			otype := (opcode & 0x0FFF)

			if otype == 0x0E0 {
				fmt.Printf("TODO: display clear at memory[0x%04x]\n", vm.Registers.PC)
			} else if otype == 0x0EE {
				fmt.Printf("TODO: return from a subroutine at memory[0x%04x]\n", vm.Registers.PC)
			}
		case 0x1000:
			vm.Registers.PC = (opcode & 0x0FFF)
			continue
		case 0x2000:
			fmt.Printf("TODO: call a subroutine at memory[0x%04x]\n", (opcode & 0x0FFF))
		case 0x3000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			// NOTE: nn can't be greater than 0xff, so this type conversion is ok
			if vm.Registers.V[x] == byte(nn) {
				vm.Registers.PC += 2
			}
		case 0x4000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			if uint16(vm.Registers.V[x]) != nn {
				vm.Registers.PC += 2
			}
		case 0x5000:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4

			if vm.Registers.V[x] == vm.Registers.V[y] {
				vm.Registers.PC += 2
			}
		case 0x6000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			vm.Registers.V[x] = byte(nn)
		case 0x7000:
			x := (opcode & 0x0F00) >> 8
			nn := (opcode & 0x00FF)

			vm.Registers.V[x] += byte(nn)
		case 0x8000:
			x := (opcode & 0x0F00) >> 8
			y := (opcode & 0x00F0) >> 4
			n := (opcode & 0x000F)
			fmt.Println(x, y, n)
		default:
			fmt.Printf("unrecognized opcode '0x%04x' at memory[0x%04x]\n", (opcode & 0xF000), vm.Registers.PC)
		}

		time.Sleep(time.Second)
		vm.Registers.PC += 2
	}
}
