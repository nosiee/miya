package memory

import "testing"

var memtest *Memory
var stacktest *Stack

func init() {
	memtest = NewMemory(CHIP8_MEMORY_SIZE)
}

func TestMemoryReset(t *testing.T) {
	for i := 0; i < CHIP8_MEMORY_SIZE; i++ {
		memtest.buffer[i] = 0xFF
	}

	memtest.Reset()

	for i := 0; i < CHIP8_MEMORY_SIZE; i++ {
		if memtest.buffer[i] != 0x00 {
			t.Errorf("got memory[0x%04x]: 0x%04x, want memory[0x%04x]: 0x%04x\n", i, memtest.buffer[i], i, 0x00)
		}
	}
}

func TestMemoryWrite(t *testing.T) {
	addr := uint16(0x200)
	data := byte(0xFF)

	memtest.Write(addr, data)

	if memtest.buffer[addr] != data {
		t.Errorf("got memory[0x%04x]: 0x%04x, want memory[0x%04x]: 0x%04x\n", addr, memtest.buffer[addr], addr, data)
	}

	memtest.Reset()
}

func TestMemoryRead(t *testing.T) {
	addr := uint16(0x200)
	data := byte(0xFF)

	memtest.Write(addr, data)

	if memtest.Read(addr) != data {
		t.Errorf("got memory[0x%04x]: 0x%04x, want memory[0x%04x]: 0x%04x\n", addr, memtest.Read(addr), addr, data)
	}

	memtest.Reset()
}

func TestMemoryRead_oobs(t *testing.T) {
	addr := uint16(CHIP8_MEMORY_SIZE + 0x10)

	if memtest.Read(addr) != 0x00 {
		t.Errorf("got memory[0x%04x]: 0x%04x, want memory[0x%04x]: 0x%04x\n", addr, memtest.Read(addr), addr, 0x00)
	}
}

func TestMemoryReadOpcode(t *testing.T) {
	memtest.Write(0x200, 0xFF)
	memtest.Write(0x201, 0xAB)

	opcode := memtest.ReadOpcode(0x200)

	if opcode != 0xFFAB {
		t.Errorf("got opcode: 0x%04x, want opcode: 0x%04x\n", opcode, 0xFFAB)
	}

	memtest.Reset()
}

func TestMemoryWriteArray(t *testing.T) {
	addr := uint16(0x200)
	data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	memtest.WriteArray(addr, data)

	for i := 0; i < len(data); i++ {
		if memtest.Read(addr+uint16(i)) != data[i] {
			t.Errorf("got memory[0x%04x]: 0x%04x, want memory[0x%04x]: 0x%04x\n", addr, memtest.Read(addr), addr, data[i])
		}
	}
}
