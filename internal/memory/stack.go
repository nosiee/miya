package memory

const CHIP8_STACK_SIZE = 0x10

type Stack struct {
	buffer []uint16
	sp     uint8
}

func NewStack(size int) *Stack {
	return &Stack{
		make([]uint16, size),
		0x00,
	}
}

func (stack *Stack) Reset() {
	for i := 0; i < CHIP8_STACK_SIZE; i++ {
		stack.buffer[i] = 0x00
	}

	stack.sp = 0
}

func (stack *Stack) Push(data uint16) {
	if stack.sp < CHIP8_STACK_SIZE {
		stack.buffer[stack.sp] = data
		stack.sp++
	}
}

func (stack *Stack) Pop() uint16 {
	if stack.sp > 0 {
		data := stack.buffer[stack.sp-1]
		stack.buffer[stack.sp-1] = 0x00
		stack.sp--

		return data
	}

	return stack.buffer[stack.sp]
}
