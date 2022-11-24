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
		opcode := vm.decodeOpcode(vm.Memory.ReadOpcode(vm.Registers.PC))
		vm.Registers.PC += 2

		switch opcode[0] {
		default:
			fmt.Printf("unrecognized opcode '0x%04x' at memory[0x%04x]\n", opcode[0], vm.Registers.PC)
		}

		time.Sleep(time.Second)
	}
}

func (vm VirtualMachine) decodeOpcode(opcode uint16) []uint16 {
	unpacked := make([]uint16, 4)

	unpacked[0] = opcode & 0xF000
	unpacked[1] = opcode & 0x0F00
	unpacked[2] = opcode & 0x00F0
	unpacked[3] = opcode & 0x000F

	return unpacked
}
