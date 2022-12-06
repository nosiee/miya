package vm

import (
	"miya/internal/memory"
	"miya/internal/screen"
	"testing"
)

var vm *VirtualMachine

func init() {
	vm = NewVirtualMachine(memory.NewMemory(memory.CHIP8_MEMORY_SIZE),
		memory.NewStack(memory.CHIP8_STACK_SIZE),
		screen.NewScreen(640, 320, "CHIP8"))
}

func TestClc(t *testing.T) {
	opcode := uint16(0x00E0)
	pc := vm.registers.PC

	vm.clc(opcode)

	for i := byte(0); i < 32; i++ {
		for k := byte(0); k < 64; k++ {
			if vm.screen.GetPixel(k, i) != 0 {
				t.Error("Screen buffer has not been cleared")
			}
		}
	}

	if pc != (pc + 2) {
		t.Error("Programm counter has not been changed")
	}
}
