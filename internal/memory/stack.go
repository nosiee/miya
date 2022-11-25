package memory

const CHIP8_STACK_SIZE = 0x10

type Stack struct {
	buffer []byte
	sp     uint8
}

func NewStack(size int) *Stack {
	return &Stack{
		make([]byte, size),
		0x00,
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
