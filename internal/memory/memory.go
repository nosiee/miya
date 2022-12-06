package memory

const CHIP8_MEMORY_SIZE = 0xFFF

type Memory struct {
	buffer []byte
}

func NewMemory(size int) *Memory {
	return &Memory{
		make([]byte, size),
	}
}

func (memory *Memory) Write(addr uint16, data byte) {
	if addr < CHIP8_MEMORY_SIZE {
		memory.buffer[addr] = data
	}
}

func (memory Memory) Read(addr uint16) byte {
	if addr < CHIP8_MEMORY_SIZE {
		return memory.buffer[addr]
	}

	return 0x00
}

func (memory *Memory) Reset() {
	for i := 0; i < CHIP8_MEMORY_SIZE; i++ {
		memory.buffer[i] = 0x00
	}
}

func (memory Memory) ReadOpcode(addr uint16) uint16 {
	var opcode uint16 = 0x00

	opcode = uint16(memory.Read(addr))
	opcode = (opcode<<8 | uint16(memory.Read(addr+1)))

	return opcode
}

func (memory *Memory) WriteArray(addr uint16, data []byte) {
	if addr < CHIP8_MEMORY_SIZE {
		for i := 0; i < len(data); i++ {
			memory.buffer[addr+uint16(i)] = data[i]
		}
	}
}
