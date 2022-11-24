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
	V0 uint8
	V1 uint8
	V2 uint8
	V3 uint8
	V4 uint8
	V5 uint8
	V6 uint8
	V7 uint8
	V8 uint8
	V9 uint8
	VA uint8
	VB uint8
	VC uint8
	VD uint8
	VF uint8
}

func NewVirtualMachine(memory *memory.Memory, stack *memory.Stack) *VirtualMachine {
	return &VirtualMachine{
		Registers: Registers{},
		Memory:    memory,
		Stack:     stack,
	}
}

func (vm *VirtualMachine) EvalLoop() {
	for {
		opcode := vm.Memory.Read(0x200 + vm.Registers.PC)
		vm.Registers.PC++

		switch opcode {
		default:
			fmt.Printf("unrecognized opcode '0x%02x' at '0x%04x'\n", opcode, vm.Registers.PC)
		}

		time.Sleep(time.Second)
	}
}
