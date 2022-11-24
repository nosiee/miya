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
		vm.Registers.PC += 2

		switch opcode {
		default:
			fmt.Printf("unrecognized opcode '0x%04x' at memory[0x%04x]\n", opcode, vm.Registers.PC)
		}

		time.Sleep(time.Second)
	}
}
