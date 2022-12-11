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

func TestSnevxvy_skip(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x9AB0)
	pc := vm.registers.PC
	vx := byte(0x00)
	vy := byte(0x01)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.snevxvy(opcode.opcode)

	if vm.registers.PC != (pc + 4) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+4)
	}
}

func TestSnevxvy(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0x9AB0)
	pc := vm.registers.PC
	vx := byte(0x10)
	vy := byte(0x10)

	vm.registers.V[opcode.x] = vx
	vm.registers.V[opcode.y] = vy

	vm.snevxvy(opcode.opcode)

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestLdi(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0xABCD)
	pc := vm.registers.PC

	vm.ldi(opcode.opcode)

	if vm.registers.I != opcode.nnn {
		t.Errorf("got I: 0x%04x, want I: 0x%04x\n", vm.registers.I, opcode.nnn)
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestJpv0(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0xBABC)
	vmz := byte(0x10)

	vm.registers.V[0x00] = vmz
	vm.jpv0(opcode.opcode)

	if vm.registers.PC != (uint16(vmz) + opcode.nnn) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, (uint16(vmz) + opcode.nnn))
	}
}

func TestRnd(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0xCABC)
	pc := vm.registers.PC

	vm.registers.V[opcode.x] = 0xFF // because generate number in the half-open interval [0x00, 0xFF), so it can't be 0xFF
	vm.rnd(opcode.opcode)

	if vm.registers.V[opcode.x] == 0xFF {
		t.Errorf("got V[x]: 0x%02x, want V[x]: [0x00, 0xFF)\n", vm.registers.V[opcode.x])
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestDrw_carry(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0xD125)
	pc := vm.registers.PC
	pixel := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0} // 0

	x := vm.registers.V[opcode.x] // 0x00 by default
	y := vm.registers.V[opcode.y] // 0x00 by default

	vm.memory.Write(vm.registers.I, pixel[0])
	vm.memory.Write(vm.registers.I+1, pixel[1])
	vm.memory.Write(vm.registers.I+2, pixel[2])
	vm.memory.Write(vm.registers.I+3, pixel[3])
	vm.memory.Write(vm.registers.I+4, pixel[4])

	vm.screen.SetPixel(x, y)
	vm.drw(opcode.opcode)
	vm.screen.SetPixel(x, y) // set pixel back, because we toggled it

	if vm.registers.V[0x0F] != 0x01 {
		t.Errorf("got V[0x0F]: 0x%02x, want V[0x0F]: 0x01\n", vm.registers.V[0x0F])
	}

	for i := byte(0); i < opcode.n; i++ {
		for k := 0; k < 8; k++ {
			if pixel[i]&(0x80>>k) != 0 {
				if vm.screen.GetPixel(x+byte(k), y+byte(i)) != 1 {
					t.Errorf("got screen[%d][%d] == 0, want screen[%d][%d] == 1\n", x+byte(k), y+byte(i), x+byte(k), y+byte(i))
				}
			}
		}
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}

func TestDrw(t *testing.T) {
	vm.Reset()

	opcode := newOpcode(0xD125)
	pc := vm.registers.PC
	pixel := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0} // 0

	x := vm.registers.V[opcode.x] // 0x00 by default
	y := vm.registers.V[opcode.y] // 0x00 by default

	vm.memory.Write(vm.registers.I, pixel[0])
	vm.memory.Write(vm.registers.I+1, pixel[1])
	vm.memory.Write(vm.registers.I+2, pixel[2])
	vm.memory.Write(vm.registers.I+3, pixel[3])
	vm.memory.Write(vm.registers.I+4, pixel[4])

	vm.drw(opcode.opcode)

	for i := byte(0); i < opcode.n; i++ {
		for k := 0; k < 8; k++ {
			if pixel[i]&(0x80>>k) != 0 {
				if vm.screen.GetPixel(x+byte(k), y+byte(i)) != 1 {
					t.Errorf("got screen[%d][%d] == 0, want screen[%d][%d] == 1\n", x+byte(k), y+byte(i), x+byte(k), y+byte(i))
				}
			}
		}
	}

	if vm.registers.PC != (pc + 2) {
		t.Errorf("got PC: 0x%04x, want PC: 0x%04x\n", vm.registers.PC, pc+2)
	}
}