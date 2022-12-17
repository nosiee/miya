package vm

import (
	"miya/internal/memory"
	"miya/internal/screen"
	"testing"
)

var vm *VirtualMachine

type testCase struct {
	test *testing.T
	name string
}

func init() {
	mem := memory.NewMemory(memory.CHIP8_MEMORY_SIZE)
	stc := memory.NewStack(memory.CHIP8_STACK_SIZE)
	vm = NewVirtualMachine(mem, stc, &screen.MockWindow{}, 10)
}

func newTestCase(test *testing.T, name string) testCase {
	return testCase{
		test: test,
		name: name,
	}
}

func (tcase testCase) assertEqualPC(value uint16) {
	if vm.Registers.PC != value {
		tcase.test.Errorf("[%s] got PC: 0x%04x, want 0x%04x\n", tcase.name, vm.Registers.PC, value)
	}
}

func (tcase testCase) assertEqualScreen(value [64][32]byte) {
	for i := byte(0); i < 32; i++ {
		for k := byte(0); k < 64; k++ {
			if vm.screen.GetPixel(k, i) != value[k][i] {
				tcase.test.Errorf("got screen[x][y] = %d, want screen[x][y] = %d\n", vm.screen.GetPixel(k, i), value[k][i])
			}
		}
	}
}

func (tcase testCase) assertEqualStackHead(value uint16) {
	head := vm.stack.Pop()

	if head != value {
		tcase.test.Errorf("[%s] got stack.pop(): 0x%04x, want stack.pop(): 0x%04x\n", tcase.name, head, value)
	}
}

func (tcase testCase) assertEqualVx(x byte, value byte) {
	if vm.Registers.V[x] != value {
		tcase.test.Errorf("[%s] got V[0x%02x] = 0x%02x, want V[0x%02x] = 0x%02x\n", tcase.name, x, vm.Registers.V[x], x, value)
	}
}

func (tcase testCase) assertNotEqualVx(x byte, value byte) {
	if vm.Registers.V[x] == value {
		tcase.test.Errorf("[%s] got V[0x%02x] = 0x%02x, want V[0x%02x] = 0x%02x\n", tcase.name, x, vm.Registers.V[x], x, value)
	}
}

func (tcase testCase) assertEqualI(value uint16) {
	if vm.Registers.I != value {
		tcase.test.Errorf("[%s] got I: 0x%04x, want I: 0x%04x\n", tcase.name, vm.Registers.I, value)
	}
}

func (tcase testCase) assertEqualKeys(x byte, value byte) {
	if vm.Keys[x] != value {
		tcase.test.Errorf("[%s] got Keys[0x%02x] = 0x%02x, want Keys[0x%02x] = 0x%02x\n", tcase.name, x, vm.Keys[x], x, value)
	}
}

func (tcase testCase) assertEqualDelayTimer(value byte) {
	if vm.DelayTimer != value {
		tcase.test.Errorf("[%s] got DelayTimer: %d, want DelayTimer: %d\n", tcase.name, vm.DelayTimer, value)
	}
}

func (tcase testCase) assertEqualSoundTimer(value byte) {
	if vm.SoundTimer != value {
		tcase.test.Errorf("[%s] got SoundTimer: %d, want SoundTimer: %d\n", tcase.name, vm.SoundTimer, value)
	}
}

func (tcase testCase) assertEqualMemory(addr uint16, value byte) {
	if vm.memory.Read(addr) != value {
		tcase.test.Errorf("[%s] got memory[0x%04x]: 0x%04x, want memory[0x%04x]: 0x%04x\n", tcase.name, addr, vm.memory.Read(addr), addr, value)
	}
}
