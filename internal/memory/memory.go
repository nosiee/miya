package memory

const CHIP8_MEMORY_SIZE = 0xFFF
const CHIP8_STACK_SIZE = 0x10

type Memory struct {
	buffer []byte
}

type Stack struct {
	buffer []byte
	sp     uint8
}

func NewMemory(size int) *Memory {
	return &Memory{
		make([]byte, size),
	}
}

func NewStack(size int) *Stack {
	return &Stack{
		make([]byte, size),
		0x00,
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

func (memory Memory) ReadOpcode(addr uint16) uint16 {
	var opcode uint16 = 0

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

func (stack *Stack) Push(data byte) {
	if stack.sp < CHIP8_STACK_SIZE {
		stack.buffer[stack.sp] = data
		stack.sp++
	}
}

func (stack *Stack) Pop() byte {
	if stack.sp > 0 {
		data := stack.buffer[stack.sp-1]
		stack.buffer[stack.sp-1] = 0x00
		stack.sp--

		return data
	}

	return stack.buffer[stack.sp]
}
