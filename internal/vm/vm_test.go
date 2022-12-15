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
			tcase.vm.screen.SetPixel(k, i)
		}
	}

	tcase.vm.clc(opcode)
	tcase.assertEqualScreen([64][32]byte{})
	tcase.assertEqualPC(0x202)
}

func TestClc_EE(t *testing.T) {
	opcode := NewOpcode(0x00EE)
	tcase := newTestCase(t, "CLC 0x00EE")

	tcase.vm.stack.Push(tcase.vm.Registers.PC)
	tcase.vm.Registers.PC = 0x255

	tcase.vm.clc(opcode)
	tcase.assertEqualPC(0x202)
}

func TestJp(t *testing.T) {
	opcode := NewOpcode(0x1ABC)
	tcase := newTestCase(t, "JP")

	tcase.vm.jp(opcode)
	tcase.assertEqualPC(opcode.nnn)
}

func TestCall(t *testing.T) {
	opcode := NewOpcode(0x2ABC)
	tcase := newTestCase(t, "CALL")

	tcase.vm.call(opcode)
	tcase.assertEqualStackHead(0x200)
	tcase.assertEqualPC(opcode.nnn)
}

func TestSevx_skip(t *testing.T) {
	opcode := NewOpcode(0x3ABC)
	tcase := newTestCase(t, "SEVX skip")

	tcase.vm.Registers.V[opcode.x] = opcode.nn

	tcase.vm.sevx(opcode)
	tcase.assertEqualPC(0x204)
}

func TestSevx(t *testing.T) {
	opcode := NewOpcode(0x3ABC)
	tcase := newTestCase(t, "SEVX")

	tcase.vm.sevx(opcode)
	tcase.assertEqualPC(0x202)
}

func TestSne_skip(t *testing.T) {
	opcode := NewOpcode(0x4ABC)
	tcase := newTestCase(t, "SNE skip")

	tcase.vm.sne(opcode)
	tcase.assertEqualPC(0x204)
}

func TestSne(t *testing.T) {
	opcode := NewOpcode(0x4ABC)
	tcase := newTestCase(t, "SNE")

	tcase.vm.Registers.V[opcode.x] = opcode.nn

	tcase.vm.sne(opcode)
	tcase.assertEqualPC(0x202)
}

func TestSevxvy_skip(t *testing.T) {
	opcode := NewOpcode(0x5ABC)
	tcase := newTestCase(t, "SEVXVY skip")

	tcase.vm.Registers.V[opcode.x] = 0x0A
	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.sevxvy(opcode)
	tcase.assertEqualPC(0x204)
}

func TestSevxvy(t *testing.T) {
	opcode := NewOpcode(0x5ABC)
	tcase := newTestCase(t, "SEVXVY")

	tcase.vm.Registers.V[opcode.x] = 0x00
	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.sevxvy(opcode)
	tcase.assertEqualPC(0x202)
}

func TestLdvx(t *testing.T) {
	opcode := NewOpcode(0x6ABC)
	tcase := newTestCase(t, "LDVX")

	tcase.vm.ldvx(opcode)
	tcase.assertEqualVx(opcode.x, opcode.nn)
	tcase.assertEqualPC(0x202)
}

func TestAdd(t *testing.T) {
	opcode := NewOpcode(0x7ABC)
	tcase := newTestCase(t, "ADD")

	tcase.vm.add(opcode)
	tcase.assertEqualVx(opcode.x, opcode.nn)
	tcase.assertEqualPC(0x202)
}

func TestVxvy_0(t *testing.T) {
	opcode := NewOpcode(0x8AB0)
	tcase := newTestCase(t, "VXVY 0x00")

	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(opcode.x, 0x0A)
	tcase.assertEqualPC(0x202)
}

func TestVxvy_1(t *testing.T) {
	opcode := NewOpcode(0x8AB1)
	tcase := newTestCase(t, "VXVY 0x01")

	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(opcode.x, (0x00 | 0x0A))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_2(t *testing.T) {
	opcode := NewOpcode(0x8AB2)
	tcase := newTestCase(t, "VXVY 0x02")

	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(opcode.x, (0x00 & 0x0A))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_3(t *testing.T) {
	opcode := NewOpcode(0x8AB3)
	tcase := newTestCase(t, "VXVY 0x03")

	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(opcode.x, (0x00 ^ 0x0A))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_4_carry(t *testing.T) {
	opcode := NewOpcode(0x8AB4)
	tcase := newTestCase(t, "VXVY 0x04 carry flag")

	tcase.vm.Registers.V[opcode.x] = 0xFF
	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x01)
	tcase.assertEqualVx(opcode.x, ((0xFF + 0x0A) & 0xFF))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_4(t *testing.T) {
	opcode := NewOpcode(0x8AB4)
	tcase := newTestCase(t, "VXVY 0x04")

	tcase.vm.Registers.V[opcode.x] = 0x0A
	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x00)
	tcase.assertEqualVx(opcode.x, (0x0A + 0x0A))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_5_carry(t *testing.T) {
	opcode := NewOpcode(0x8AB5)
	tcase := newTestCase(t, "VXVY 0x05 carry flag")

	tcase.vm.Registers.V[opcode.x] = 0x10
	tcase.vm.Registers.V[opcode.y] = 0x05

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x01)
	tcase.assertEqualVx(opcode.x, (0x10 - 0x05))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_5(t *testing.T) {
	opcode := NewOpcode(0x8AB5)
	tcase := newTestCase(t, "VXVY 0x05")

	tcase.vm.Registers.V[opcode.x] = 0x05
	tcase.vm.Registers.V[opcode.y] = 0x10

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x00)
	tcase.assertEqualVx(opcode.x, ((0x05 - 0x10) & 0xFF))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_6(t *testing.T) {
	opcode := NewOpcode(0x8AB6)
	tcase := newTestCase(t, "VXVY 0x06")

	tcase.vm.Registers.V[opcode.x] = 0x10

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x00)
	tcase.assertEqualVx(opcode.x, (0x10 >> 1))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_7_carry(t *testing.T) {
	opcode := NewOpcode(0x8AB7)
	tcase := newTestCase(t, "VXVY 0x07 carry flag")

	tcase.vm.Registers.V[opcode.x] = 0x0A
	tcase.vm.Registers.V[opcode.y] = 0xFF

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x01)
	tcase.assertEqualVx(opcode.x, (0xFF - 0x0A))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_7(t *testing.T) {
	opcode := NewOpcode(0x8AB7)
	tcase := newTestCase(t, "VXVY 0x07")

	tcase.vm.Registers.V[opcode.x] = 0xFF
	tcase.vm.Registers.V[opcode.y] = 0x0A

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, 0x00)
	tcase.assertEqualVx(opcode.x, ((0x0A - 0xFF) & 0xFF))
	tcase.assertEqualPC(0x202)
}

func TestVxvy_e(t *testing.T) {
	opcode := NewOpcode(0x8ABE)
	tcase := newTestCase(t, "VXVY 0x0E")

	tcase.vm.Registers.V[opcode.x] = 0x10

	tcase.vm.vxvy(opcode)
	tcase.assertEqualVx(0x0F, (0x10 & 0x80))
	tcase.assertEqualVx(opcode.x, (0x10 << 1))
	tcase.assertEqualPC(0x202)
}

func TestSnevxvy_skip(t *testing.T) {
	opcode := NewOpcode(0x9AB0)
	tcase := newTestCase(t, "SNEVXVY skip")

	tcase.vm.Registers.V[opcode.x] = 0x00
	tcase.vm.Registers.V[opcode.y] = 0x10

	tcase.vm.snevxvy(opcode)
	tcase.assertEqualPC(0x204)
}

func TestSnevxvy(t *testing.T) {
	opcode := NewOpcode(0x9AB0)
	tcase := newTestCase(t, "SNEVXVY")

	tcase.vm.Registers.V[opcode.x] = 0x10
	tcase.vm.Registers.V[opcode.y] = 0x10

	tcase.vm.snevxvy(opcode)
	tcase.assertEqualPC(0x202)
}

func TestLdi(t *testing.T) {
	opcode := NewOpcode(0xABCD)
	tcase := newTestCase(t, "LDI")

	tcase.vm.ldi(opcode)
	tcase.assertEqualI(opcode.nnn)
	tcase.assertEqualPC(0x202)
}

func TestJpv0(t *testing.T) {
	opcode := NewOpcode(0xBABC)
	tcase := newTestCase(t, "JPV0")

	tcase.vm.Registers.V[0x00] = 0x10

	tcase.vm.jpv0(opcode)
	tcase.assertEqualPC(0x10 + opcode.nnn)
}

func TestRnd(t *testing.T) {
	opcode := NewOpcode(0xCABC)
	tcase := newTestCase(t, "RND")

	tcase.vm.Registers.V[opcode.x] = 0xFF // because generate number in the half-open interval [0x00, 0xFF), so it can't be 0xFF

	tcase.vm.rnd(opcode)
	tcase.assertNotEqualVx(opcode.x, 0xFF)
	tcase.assertEqualPC(0x202)
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

	tcase.vm.memory.WriteArray(tcase.vm.Registers.I, pixel)
	tcase.vm.screen.SetPixel(0x00, 0x00) // set pixel to trigger carry flag

	tcase.vm.drw(opcode)
	tcase.vm.screen.SetPixel(0x00, 0x00) // set pixel back, because we toggled it
	tcase.assertEqualScreen(sbuffer)
	tcase.assertEqualVx(0x0F, 0x01)
	tcase.assertEqualPC(0x202)
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

	tcase.vm.memory.WriteArray(tcase.vm.Registers.I, pixel)

	tcase.vm.drw(opcode)
	tcase.assertEqualScreen(sbuffer)
	tcase.assertEqualPC(0x202)
}

func TestSkp_9e_skip(t *testing.T) {
	opcode := NewOpcode(0xE29E)
	tcase := newTestCase(t, "SKP 0x09 skip")

	tcase.vm.Keys[0x00] = 0x01
	tcase.vm.skp(opcode)
	tcase.assertEqualPC(0x204)
}

func TestSkp_9e(t *testing.T) {
	opcode := NewOpcode(0xE29E)
	tcase := newTestCase(t, "SKP 0x09")

	tcase.vm.skp(opcode)
	tcase.assertEqualPC(0x202)
}

func TestSkp_A1_skip(t *testing.T) {
	opcode := NewOpcode(0xE2A1)
	tcase := newTestCase(t, "SKP 0x01 skip")

	tcase.vm.skp(opcode)
	tcase.assertEqualPC(0x204)
}

func TestSkp_A1(t *testing.T) {
	opcode := NewOpcode(0xE2A1)
	tcase := newTestCase(t, "SKP 0x01 skip")

	tcase.vm.Keys[0x00] = 0x01

	tcase.vm.skp(opcode)
	tcase.assertEqualPC(0x202)
}

func TestLdf_07(t *testing.T) {
	opcode := NewOpcode(0xF307)
	tcase := newTestCase(t, "LDF 0x07")

	tcase.vm.DelayTimer = 0x20

	tcase.vm.ldf(opcode)
	tcase.assertEqualVx(opcode.x, 0x20)
	tcase.assertEqualPC(0x202)
}

func TestLdf_0A(t *testing.T) {
	opcode := NewOpcode(0xF30A)
	tcase := newTestCase(t, "LDF 0x0A")

	go tcase.vm.ldf(opcode)
	tcase.vm.keypressed <- sdl.K_3

	tcase.assertEqualVx(opcode.x, sdl.K_3)
	tcase.assertEqualPC(0x202)
}

func TestLdf_15(t *testing.T) {
	opcode := NewOpcode(0xFA15)
	tcase := newTestCase(t, "LDF 0x15")

	tcase.vm.Registers.V[opcode.x] = 0x10

	tcase.vm.ldf(opcode)
	tcase.assertEqualDelayTimer(0x10)
	tcase.assertEqualPC(0x202)
}

func TestLdf_18(t *testing.T) {
	opcode := NewOpcode(0xFA18)
	tcase := newTestCase(t, "LDF 0x18")

	tcase.vm.Registers.V[opcode.x] = 0x10

	tcase.vm.ldf(opcode)
	tcase.assertEqualSoundTimer(0x10)
	tcase.assertEqualPC(0x202)
}

func TestLdf_1E(t *testing.T) {
	opcode := NewOpcode(0xFC1E)
	tcase := newTestCase(t, "LDF 0x1E")

	tcase.vm.Registers.I = 0x10
	tcase.vm.Registers.V[opcode.x] = 0x20

	tcase.vm.ldf(opcode)
	tcase.assertEqualI(0x30)
	tcase.assertEqualPC(0x202)
}

func TestLdf_29(t *testing.T) {
	opcode := NewOpcode(0xFC29)
	tcase := newTestCase(t, "LDF 0x29")

	tcase.vm.Registers.V[opcode.x] = 0x05

	tcase.vm.ldf(opcode)
	tcase.assertEqualI(0x19)
	tcase.assertEqualPC(0x202)
}

func TestLdf_33(t *testing.T) {
	opcode := NewOpcode(0xFC33)
	tcase := newTestCase(t, "LDF 0x33")

	tcase.vm.Registers.V[opcode.x] = 0xFC

	tcase.vm.ldf(opcode)
	tcase.assertEqualMemory(tcase.vm.Registers.I, (0xFC / 100))
	tcase.assertEqualMemory(tcase.vm.Registers.I+1, ((0xFC / 10) % 10))
	tcase.assertEqualMemory(tcase.vm.Registers.I+2, ((0xFC % 100) % 10))
	tcase.assertEqualPC(0x202)
}

func TestLdf_55(t *testing.T) {
	opcode := NewOpcode(0xFA55)
	tcase := newTestCase(t, "LDF 0x55")

	for i := byte(0); i <= opcode.x; i++ {
		tcase.vm.Registers.V[i] = i
	}

	tcase.vm.ldf(opcode)
	tcase.vm.Registers.I -= (uint16(opcode.x) + 1) // because I is incemented opcode.x times

	for i := byte(0); i <= opcode.x; i++ {
		tcase.assertEqualMemory(tcase.vm.Registers.I+uint16(i), i)
	}

	tcase.assertEqualPC(0x202)
}

func TestLdf_65(t *testing.T) {
	opcode := NewOpcode(0xFA65)
	tcase := newTestCase(t, "LDF 0x65")

	tcase.vm.memory.WriteArray(tcase.vm.Registers.I, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	tcase.vm.ldf(opcode)
	tcase.vm.Registers.I -= (uint16(opcode.x) + 1) // because I is incremented opcode.x times

	for i := byte(0); i <= opcode.x; i++ {
		tcase.assertEqualVx(i, i)
	}

	tcase.assertEqualPC(0x202)
}
