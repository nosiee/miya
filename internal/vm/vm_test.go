package vm

import (
	"testing"

	"github.com/veandco/go-sdl2/sdl"
)

func TestClc_E0(t *testing.T) {
	opcode := NewOpcode(0x00E0)
	tcase := newTestCase(t, "CLC 0x00E0")

	for i := byte(0); i < 32; i++ {
		for k := byte(0); k < 64; k++ {
			vm.screen.SetPixel(k, i)
		}
	}

	vm.clc(opcode)
	tcase.assertEqualScreen([64][32]byte{})
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestClc_EE(t *testing.T) {
	opcode := NewOpcode(0x00EE)
	tcase := newTestCase(t, "CLC 0x00EE")

	vm.stack.Push(vm.Registers.PC)
	vm.Registers.PC = 0x255

	vm.clc(opcode)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestJp(t *testing.T) {
	opcode := NewOpcode(0x1ABC)
	tcase := newTestCase(t, "JP")

	vm.jp(opcode)
	tcase.assertEqualPC(opcode.nnn)

	vm.Reset()
}

func TestCall(t *testing.T) {
	opcode := NewOpcode(0x2ABC)
	tcase := newTestCase(t, "CALL")

	vm.call(opcode)
	tcase.assertEqualStackHead(0x200)
	tcase.assertEqualPC(opcode.nnn)

	vm.Reset()
}

func TestSevx_skip(t *testing.T) {
	opcode := NewOpcode(0x3ABC)
	tcase := newTestCase(t, "SEVX skip")

	vm.Registers.V[opcode.x] = opcode.nn

	vm.sevx(opcode)
	tcase.assertEqualPC(0x204)

	vm.Reset()
}

func TestSevx(t *testing.T) {
	opcode := NewOpcode(0x3ABC)
	tcase := newTestCase(t, "SEVX")

	vm.sevx(opcode)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestSne_skip(t *testing.T) {
	opcode := NewOpcode(0x4ABC)
	tcase := newTestCase(t, "SNE skip")

	vm.sne(opcode)
	tcase.assertEqualPC(0x204)

	vm.Reset()
}

func TestSne(t *testing.T) {
	opcode := NewOpcode(0x4ABC)
	tcase := newTestCase(t, "SNE")

	vm.Registers.V[opcode.x] = opcode.nn

	vm.sne(opcode)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestSevxvy_skip(t *testing.T) {
	opcode := NewOpcode(0x5ABC)
	tcase := newTestCase(t, "SEVXVY skip")

	vm.Registers.V[opcode.x] = 0x0A
	vm.Registers.V[opcode.y] = 0x0A

	vm.sevxvy(opcode)
	tcase.assertEqualPC(0x204)

	vm.Reset()
}

func TestSevxvy(t *testing.T) {
	opcode := NewOpcode(0x5ABC)
	tcase := newTestCase(t, "SEVXVY")

	vm.Registers.V[opcode.x] = 0x00
	vm.Registers.V[opcode.y] = 0x0A

	vm.sevxvy(opcode)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdvx(t *testing.T) {
	opcode := NewOpcode(0x6ABC)
	tcase := newTestCase(t, "LDVX")

	vm.ldvx(opcode)
	tcase.assertEqualVx(opcode.x, opcode.nn)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestAdd(t *testing.T) {
	opcode := NewOpcode(0x7ABC)
	tcase := newTestCase(t, "ADD")

	vm.add(opcode)
	tcase.assertEqualVx(opcode.x, opcode.nn)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_0(t *testing.T) {
	opcode := NewOpcode(0x8AB0)
	tcase := newTestCase(t, "VXVY 0x00")

	vm.Registers.V[opcode.y] = 0x0A

	vm.vxvy(opcode)
	tcase.assertEqualVx(opcode.x, 0x0A)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_1(t *testing.T) {
	opcode := NewOpcode(0x8AB1)
	tcase := newTestCase(t, "VXVY 0x01")

	vm.Registers.V[opcode.y] = 0x0A

	vm.vxvy(opcode)
	tcase.assertEqualVx(opcode.x, (0x00 | 0x0A))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_2(t *testing.T) {
	opcode := NewOpcode(0x8AB2)
	tcase := newTestCase(t, "VXVY 0x02")

	vm.Registers.V[opcode.y] = 0x0A

	vm.vxvy(opcode)
	tcase.assertEqualVx(opcode.x, (0x00 & 0x0A))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_3(t *testing.T) {
	opcode := NewOpcode(0x8AB3)
	tcase := newTestCase(t, "VXVY 0x03")

	vm.Registers.V[opcode.y] = 0x0A

	vm.vxvy(opcode)
	tcase.assertEqualVx(opcode.x, (0x00 ^ 0x0A))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_4_carry(t *testing.T) {
	opcode := NewOpcode(0x8AB4)
	tcase := newTestCase(t, "VXVY 0x04 carry flag")

	vm.Registers.V[opcode.x] = 0xFF
	vm.Registers.V[opcode.y] = 0x0A

	vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x01)
	tcase.assertEqualVx(opcode.x, ((0xFF + 0x0A) & 0xFF))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_4(t *testing.T) {
	opcode := NewOpcode(0x8AB4)
	tcase := newTestCase(t, "VXVY 0x04")

	vm.Registers.V[opcode.x] = 0x0A
	vm.Registers.V[opcode.y] = 0x0A

	vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x00)
	tcase.assertEqualVx(opcode.x, (0x0A + 0x0A))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_5_carry(t *testing.T) {
	opcode := NewOpcode(0x8AB5)
	tcase := newTestCase(t, "VXVY 0x05 carry flag")

	vm.Registers.V[opcode.x] = 0x10
	vm.Registers.V[opcode.y] = 0x05

	vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x01)
	tcase.assertEqualVx(opcode.x, (0x10 - 0x05))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_5(t *testing.T) {
	opcode := NewOpcode(0x8AB5)
	tcase := newTestCase(t, "VXVY 0x05")

	vm.Registers.V[opcode.x] = 0x05
	vm.Registers.V[opcode.y] = 0x10

	vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x00)
	tcase.assertEqualVx(opcode.x, ((0x05 - 0x10) & 0xFF))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_6(t *testing.T) {
	opcode := NewOpcode(0x8AB6)
	tcase := newTestCase(t, "VXVY 0x06")

	vm.Registers.V[opcode.x] = 0x10

	vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x00)
	tcase.assertEqualVx(opcode.x, (0x10 >> 1))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_7_carry(t *testing.T) {
	opcode := NewOpcode(0x8AB7)
	tcase := newTestCase(t, "VXVY 0x07 carry flag")

	vm.Registers.V[opcode.x] = 0x0A
	vm.Registers.V[opcode.y] = 0xFF

	vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x01)
	tcase.assertEqualVx(opcode.x, (0xFF - 0x0A))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_7(t *testing.T) {
	opcode := NewOpcode(0x8AB7)
	tcase := newTestCase(t, "VXVY 0x07")

	vm.Registers.V[opcode.x] = 0xFF
	vm.Registers.V[opcode.y] = 0x0A

	vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x00)
	tcase.assertEqualVx(opcode.x, ((0x0A - 0xFF) & 0xFF))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestVxvy_e(t *testing.T) {
	opcode := NewOpcode(0x8ABE)
	tcase := newTestCase(t, "VXVY 0x0E")

	vm.Registers.V[opcode.x] = 0x10

	vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, (0x10 & 0x80))
	tcase.assertEqualVx(opcode.x, (0x10 << 1))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestSnevxvy_skip(t *testing.T) {
	opcode := NewOpcode(0x9AB0)
	tcase := newTestCase(t, "SNEVXVY skip")

	vm.Registers.V[opcode.x] = 0x00
	vm.Registers.V[opcode.y] = 0x10

	vm.snevxvy(opcode)
	tcase.assertEqualPC(0x204)

	vm.Reset()
}

func TestSnevxvy(t *testing.T) {
	opcode := NewOpcode(0x9AB0)
	tcase := newTestCase(t, "SNEVXVY")

	vm.Registers.V[opcode.x] = 0x10
	vm.Registers.V[opcode.y] = 0x10

	vm.snevxvy(opcode)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdi(t *testing.T) {
	opcode := NewOpcode(0xABCD)
	tcase := newTestCase(t, "LDI")

	vm.ldi(opcode)
	tcase.assertEqualI(opcode.nnn)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestJpv0(t *testing.T) {
	opcode := NewOpcode(0xBABC)
	tcase := newTestCase(t, "JPV0")

	vm.Registers.V[0x00] = 0x10

	vm.jpv0(opcode)
	tcase.assertEqualPC(0x10 + opcode.nnn)

	vm.Reset()
}

func TestRnd(t *testing.T) {
	opcode := NewOpcode(0xCABC)
	tcase := newTestCase(t, "RND")

	vm.Registers.V[opcode.x] = 0xFF // because we generate number in the half-open interval [0x00, 0xFF), so it can't be 0xFF

	vm.rnd(opcode)
	tcase.assertNotEqualVx(opcode.x, 0xFF)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestDrw_carry(t *testing.T) {
	opcode := NewOpcode(0xD125)
	tcase := newTestCase(t, "DRW carry flag")

	pixel := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0} // 0
	sbuffer := [64][32]byte{
		{1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	vm.memory.WriteArray(vm.Registers.I, pixel)
	vm.screen.SetPixel(0x00, 0x00) // set pixel to trigger carry flag

	vm.drw(opcode)
	vm.screen.SetPixel(0x00, 0x00) // set pixel back, because we toggled it
	tcase.assertEqualScreen(sbuffer)
	tcase.assertEqualVx(0x0F, 0x01)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestDrw(t *testing.T) {
	opcode := NewOpcode(0xD125)
	tcase := newTestCase(t, "DRW")

	pixel := []byte{0xF0, 0x90, 0x90, 0x90, 0xF0} // 0
	sbuffer := [64][32]byte{
		{1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	vm.memory.WriteArray(vm.Registers.I, pixel)

	vm.drw(opcode)
	tcase.assertEqualScreen(sbuffer)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestSkp_9e_skip(t *testing.T) {
	opcode := NewOpcode(0xE29E)
	tcase := newTestCase(t, "SKP 0x09 skip")

	vm.Keys[0x00] = 0x01
	vm.skp(opcode)
	tcase.assertEqualPC(0x204)

	vm.Reset()
}

func TestSkp_9e(t *testing.T) {
	opcode := NewOpcode(0xE29E)
	tcase := newTestCase(t, "SKP 0x09")

	vm.skp(opcode)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestSkp_A1_skip(t *testing.T) {
	opcode := NewOpcode(0xE2A1)
	tcase := newTestCase(t, "SKP 0x01 skip")

	vm.skp(opcode)
	tcase.assertEqualPC(0x204)

	vm.Reset()
}

func TestSkp_A1(t *testing.T) {
	opcode := NewOpcode(0xE2A1)
	tcase := newTestCase(t, "SKP 0x01 skip")

	vm.Keys[0x00] = 0x01

	vm.skp(opcode)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_07(t *testing.T) {
	opcode := NewOpcode(0xF307)
	tcase := newTestCase(t, "LDF 0x07")

	vm.DelayTimer = 0x20

	vm.ldf(opcode)
	tcase.assertEqualVx(opcode.x, 0x20)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_0A(t *testing.T) {
	opcode := NewOpcode(0xF30A)
	tcase := newTestCase(t, "LDF 0x0A")

	go vm.ldf(opcode)
	vm.keypressed <- sdl.K_3

	tcase.assertEqualVx(opcode.x, sdl.K_3)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_15(t *testing.T) {
	opcode := NewOpcode(0xFA15)
	tcase := newTestCase(t, "LDF 0x15")

	vm.Registers.V[opcode.x] = 0x10

	vm.ldf(opcode)
	tcase.assertEqualDelayTimer(0x10)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_18(t *testing.T) {
	opcode := NewOpcode(0xFA18)
	tcase := newTestCase(t, "LDF 0x18")

	vm.Registers.V[opcode.x] = 0x10

	vm.ldf(opcode)
	tcase.assertEqualSoundTimer(0x10)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_1E(t *testing.T) {
	opcode := NewOpcode(0xFC1E)
	tcase := newTestCase(t, "LDF 0x1E")

	vm.Registers.I = 0x10
	vm.Registers.V[opcode.x] = 0x20

	vm.ldf(opcode)
	tcase.assertEqualI(0x30)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_29(t *testing.T) {
	opcode := NewOpcode(0xFC29)
	tcase := newTestCase(t, "LDF 0x29")

	vm.Registers.V[opcode.x] = 0x05

	vm.ldf(opcode)
	tcase.assertEqualI(0x19)
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_33(t *testing.T) {
	opcode := NewOpcode(0xFC33)
	tcase := newTestCase(t, "LDF 0x33")

	vm.Registers.V[opcode.x] = 0xFC

	vm.ldf(opcode)
	tcase.assertEqualMemory(vm.Registers.I, (0xFC / 100))
	tcase.assertEqualMemory(vm.Registers.I+1, ((0xFC / 10) % 10))
	tcase.assertEqualMemory(vm.Registers.I+2, ((0xFC % 100) % 10))
	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_55(t *testing.T) {
	opcode := NewOpcode(0xFA55)
	tcase := newTestCase(t, "LDF 0x55")

	for i := byte(0); i <= opcode.x; i++ {
		vm.Registers.V[i] = i
	}

	vm.ldf(opcode)
	tcase.assertEqualI(uint16(opcode.x) + 1)

	vm.Registers.I -= (uint16(opcode.x) + 1) // because I is incemented opcode.x times
	for i := byte(0); i <= opcode.x; i++ {
		tcase.assertEqualMemory(vm.Registers.I+uint16(i), i)
	}

	tcase.assertEqualPC(0x202)

	vm.Reset()
}

func TestLdf_65(t *testing.T) {
	opcode := NewOpcode(0xFA65)
	tcase := newTestCase(t, "LDF 0x65")

	vm.memory.WriteArray(vm.Registers.I, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	vm.ldf(opcode)

	tcase.assertEqualI(uint16(opcode.x + 1))
	vm.Registers.I -= (uint16(opcode.x) + 1) // because I is incremented opcode.x times

	for i := byte(0); i <= opcode.x; i++ {
		tcase.assertEqualVx(i, i)
	}

	tcase.assertEqualPC(0x202)

	vm.Reset()
}
