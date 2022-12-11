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
		screen.NewScreen(640, 320, "CHIP8-TEST"))
}

func TestClc_E0(t *testing.T) {
	opcode := newOpcode(0x00E0)
	pc := vm.registers.PC

	vm.clc(opcode.opcode)

	for i := byte(0); i < 32; i++ {
		for k := byte(0); k < 64; k++ {
			if vm.screen.GetPixel(k, i) != 0 {
				t.Errorf("got x: %d, y: %d, want x: 0, y: 0\n", k, i)
			}
		}
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestClc_EE(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x00EE)
	pc := vm.registers.PC

	vm.stack.Push(vm.registers.PC)
	vm.registers.PC = 0x255

	vm.clc(opcode.opcode)

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, vm.registers.PC+2)
	}
}

func TestJp(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x1ABC)

	vm.jp(opcode.opcode)

	if vm.registers.PC != (opcode.nnn) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, opcode.nnn)
	}
}

func TestCall(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x2ABC)
	pc := vm.registers.PC

	vm.call(opcode.opcode)
	head := vm.stack.Pop()

	if head != pc {
		t.Errorf("got PC from stack: 0x%04x, want PC from stack: 0x%04x\n", head, pc)
	}

	if vm.registers.PC != opcode.nnn {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, opcode.nnn)
	}
}

func TestSevx_skip(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x3ABC)
	pc := vm.registers.PC

	vm.registers.V[opcode.x] = opcode.nn
	vm.sevx(opcode.opcode)

	if vm.registers.PC != (pc + 4) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+4)
	}
}

func TestSevx(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x3ABC)
	pc := vm.registers.PC

	vm.sevx(opcode.opcode)

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestSne_skip(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x4ABC)
	pc := vm.registers.PC

	vm.sne(opcode.opcode)

	if vm.registers.PC != (pc + 4) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+4)
	}
}

func TestSne(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x4ABC)
	pc := vm.registers.PC

	vm.registers.V[opcode.x] = opcode.nn
	vm.sne(opcode.opcode)

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+4)
	}
}

func TestSevxvy_skip(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x5ABC)
	pc := vm.registers.PC
	vx := byte(0x0A)
	vy := byte(0x0A)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.sevxvy(opcode.opcode)

	if vm.registers.PC != (pc + 4) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+4)
	}
}

func TestSevxvy(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x5ABC)
	pc := vm.registers.PC
	vx := byte(0x0A)
	vy := byte(0x0B)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.sevxvy(opcode.opcode)

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestLdvx(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x6ABC)
	pc := vm.registers.PC

	vm.ldvx(opcode.opcode)

	if vm.registers.V[opcode.x] != opcode.nn {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], opcode.nn)
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestAdd(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x7ABC)
	pc := vm.registers.PC
	vx := vm.registers.V[opcode.x]

	vm.add(opcode.opcode)

	if vm.registers.V[opcode.x] != (vx + opcode.nn) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (vx + opcode.nn))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_0(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB0)
	pc := vm.registers.PC
	vx := byte(0x00)
	vy := byte(0x0A)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[opcode.x] != 0x0A {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], 0x0A)
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_1(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB1)
	pc := vm.registers.PC
	vx := byte(0x00)
	vy := byte(0x0A)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[opcode.x] != (vx | vy) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (vx | vy))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_2(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB2)
	pc := vm.registers.PC
	vx := byte(0x00)
	vy := byte(0x0A)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[opcode.x] != (vx & vy) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (vx | vy))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_3(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB3)
	pc := vm.registers.PC
	vx := byte(0x00)
	vy := byte(0x0A)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[opcode.x] != (vx ^ vy) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (vx | vy))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_4_carry(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB4)
	pc := vm.registers.PC
	vx := byte(0xFF)
	vy := byte(0x0A)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[0x0F] != 1 {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x%02x\n", vm.registers.V[0x0F], 0x01)
	}

	if vm.registers.V[opcode.x] != ((0xFF + 0x0A) & 0xFF) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], ((0xFF + 0x0A) & 0xFF))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_4(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB4)
	pc := vm.registers.PC
	vx := byte(0x0A)
	vy := byte(0x0A)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[0x0F] != 0 {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x%02x\n", vm.registers.V[0x0F], 0x00)
	}

	if vm.registers.V[opcode.x] != (0x0A + 0x0A) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (0x0A + 0x0A))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_5_carry(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB5)
	pc := vm.registers.PC
	vx := byte(0x00)
	vy := byte(0x10)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[0x0F] != 1 {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x%02x\n", vm.registers.V[0x0F], 0x00)
	}

	if vm.registers.V[opcode.x] != ((vx - vy) & 0xFF) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], ((vx - vy) & 0xFF))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_5(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB5)
	pc := vm.registers.PC
	vx := byte(0xff)
	vy := byte(0x10)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[0x0F] != 0 {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x%02x\n", vm.registers.V[0x0F], 0x00)
	}

	if vm.registers.V[opcode.x] != (vx - vy) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (vx - vy))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_6(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB6)
	pc := vm.registers.PC
	vx := byte(0x10)

	vm.registers.V[opcode.x] = vx

	vm.vxvy(opcode.opcode)

	if vm.registers.V[0x0F] != (vx & 0x01) {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x%02x\n", vm.registers.V[0x0F], (vx & 0x01))
	}

	if vm.registers.V[opcode.x] != (vx >> 1) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (vx >> 1))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_7_carry(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB7)
	pc := vm.registers.PC
	vx := byte(0xFF)
	vy := byte(0x00)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[0x0F] != 1 {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x%02x\n", vm.registers.V[0x0F], 0x01)
	}

	if vm.registers.V[opcode.x] != ((vy - vx) & 0xFF) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], ((vy - vx) & 0xFF))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_7(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8AB7)
	pc := vm.registers.PC

	vx := byte(0x00)
	vy := byte(0xFF)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.vxvy(opcode.opcode)

	if vm.registers.V[0x0F] != 0 {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x%02x\n", vm.registers.V[0x0F], 0x00)
	}

	if vm.registers.V[opcode.x] != (vy - vx) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (vy - vx))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestVxvy_e(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x8ABE)
	pc := vm.registers.PC
	vx := byte(0x10)

	vm.registers.V[opcode.x] = vx
	vm.vxvy(opcode.opcode)

	if vm.registers.V[0x0F] != (vx & 0x80) {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x%02x\n", vm.registers.V[0x0F], (vx & 0x80))
	}

	if vm.registers.V[opcode.x] != (vx << 1) {
		t.Errorf("got V[x]: 0x%02x, want V[x]: 0x%02x\n", vm.registers.V[opcode.x], (vx << 1))
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}
