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
				t.Errorf("got x: %d, y: %d, should be x: 0, y: 0\n", k, i)
			}
		}
	}

	if (pc + 2) != vm.registers.PC {
		t.Errorf("%d want, got %d\n", pc+2, vm.registers.PC)
	}

	opcode = 0x00EE
	pc = vm.registers.PC
	vm.stack.Push(vm.registers.PC)

	vm.registers.PC += 0x10
	vm.clc(opcode)

	if (pc + 2) != vm.registers.PC {
		t.Errorf("%x want, got %x\n", pc+2, vm.registers.PC)
	}
}

func TestJp(t *testing.T) {
	vm.Reset()

	opcode := uint16(0x1222)

	vm.jp(opcode)

	if vm.registers.PC != (opcode & 0x0FFF) {
		t.Errorf("%x want, got %x\n", (opcode & 0x0FFF), vm.registers.PC)
	}
}

func TestCall(t *testing.T) {
	vm.Reset()

	opcode := uint16(0x2222)
	pc := vm.registers.PC

	vm.call(opcode)
	sh := vm.stack.Pop()

	if sh != pc {
		t.Errorf("%x want, got %x\n", pc, sh)
	}

	if vm.registers.PC != (opcode & 0x0FFF) {
		t.Errorf("%x want, got %x\n", (opcode & 0x0FFF), vm.registers.PC)
	}
}
