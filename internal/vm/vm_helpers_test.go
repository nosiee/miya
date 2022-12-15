package vm

import (
	"miya/internal/memory"
	"miya/internal/screen"
	"testing"
)

type testCase struct {
	test *testing.T
	name string
	vm   *VirtualMachine
}

func newTestCase(test *testing.T, name string) testCase {
	mem := memory.NewMemory(memory.CHIP8_MEMORY_SIZE)
	stc := memory.NewStack(memory.CHIP8_STACK_SIZE)
	scr, _ := screen.NewScreen(640, 320, "CHIP8-TEST", 10)
	vm := NewVirtualMachine(mem, stc, scr, 10)

	return testCase{
		test: test,
		name: name,
		vm:   vm,
	}
}

func (tcase testCase) assertEqualPC(value uint16) {
	if tcase.vm.Registers.PC != value {
		tcase.test.Errorf("[%s] got PC: 0x%04x, want 0x%04x\n", tcase.name, tcase.vm.Registers.PC, value)
	}
}

func (tcase testCase) assertEqualScreen(value [64][32]byte) {
	for i := byte(0); i < 32; i++ {
		for k := byte(0); k < 64; k++ {
			if tcase.vm.screen.GetPixel(k, i) != value[k][i] {
				tcase.test.Errorf("got screen[x][y] = %d, want screen[x][y] = %d\n", tcase.vm.screen.GetPixel(k, i), value[k][i])
			}
		}
	}
}

func (tcase testCase) assertEqualStackHead(value uint16) {
	head := tcase.vm.stack.Pop()

	if head != value {
		tcase.test.Errorf("[%s] got stack.pop(): 0x%04x, want stack.pop(): 0x%04x\n", tcase.name, head, value)
	}
}

func (tcase testCase) assertEqualVx(x byte, value byte) {
	if tcase.vm.Registers.V[x] != value {
		tcase.test.Errorf("[%s] got V[0x%02x] = 0x%02x, want V[0x%02x] = 0x%02x\n", tcase.name, x, tcase.vm.Registers.V[x], x, value)
	}
}

func (tcase testCase) assertNotEqualVx(x byte, value byte) {
	if tcase.vm.Registers.V[x] == value {
		tcase.test.Errorf("[%s] got V[0x%02x] = 0x%02x, want V[0x%02x] = 0x%02x\n", tcase.name, x, tcase.vm.Registers.V[x], x, value)
	}
}

func (tcase testCase) assertEqualI(value uint16) {
	if tcase.vm.Registers.I != value {
		tcase.test.Errorf("[%s] got I: 0x%04x, want I: 0x%04x\n", tcase.name, tcase.vm.Registers.I, value)
	}
}

func (tcase testCase) assertEqualDelayTimer(value byte) {
	if tcase.vm.DelayTimer != value {
		tcase.test.Errorf("[%s] got DelayTimer: %d, want DelayTimer: %d\n", tcase.name, tcase.vm.DelayTimer, value)
	}
}

func (tcase testCase) assertEqualSoundTimer(value byte) {
	if tcase.vm.SoundTimer != value {
		tcase.test.Errorf("[%s] got SoundTimer: %d, want SoundTimer: %d\n", tcase.name, tcase.vm.SoundTimer, value)
	}
}

func (tcase testCase) assertEqualMemory(addr uint16, value byte) {
	if tcase.vm.memory.Read(addr) != value {
		tcase.test.Errorf("[%s] got memory[0x%04x]: 0x%04x, want memory[0x%04x]: 0x%04x\n", tcase.name, addr, tcase.vm.memory.Read(addr), addr, value)
	}
}
