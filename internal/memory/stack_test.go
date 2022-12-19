package memory

import "testing"

func init() {
	stacktest = NewStack(CHIP8_STACK_SIZE)
}

func TestStackReset(t *testing.T) {
	for i := 0; i < CHIP8_STACK_SIZE; i++ {
		stacktest.buffer[i] = 0xFF
	}

	stacktest.Reset()

	for i := 0; i < CHIP8_STACK_SIZE; i++ {
		if stacktest.buffer[i] != 0x00 {
			t.Errorf("got stack[0x%02x]: 0x%02x, want stack[0x%02x]: 0x%02x\n", i, stacktest.buffer[i], i, 0x00)
		}
	}

	if stacktest.sp != 0 {
		t.Errorf("got SP: %d, want SP: %d", stacktest.sp, 0)
	}
}

func TestStackPush(t *testing.T) {
	stacktest.Push(0xFF)
	stacktest.Push(0xAB)

	if stacktest.sp != 2 {
		t.Errorf("got SP: %d, want SP: %d", stacktest.sp, 2)
	}

	if stacktest.buffer[stacktest.sp-1] != 0xAB {
		t.Errorf("got stack[0x%02x]: 0x%02x, want stack[0x%02x]: 0x%02x\n", stacktest.sp-1, stacktest.buffer[stacktest.sp-1], stacktest.sp-1, 0xAB)
	}

	if stacktest.buffer[stacktest.sp-2] != 0xFF {
		t.Errorf("got stack[0x%02x]: 0x%02x, want stack[0x%02x]: 0x%02x\n", stacktest.sp-2, stacktest.buffer[stacktest.sp-2], stacktest.sp-2, 0xFF)
	}

	stacktest.Reset()
}

func TestStackPop(t *testing.T) {
	stacktest.Push(0xFF)
	stacktest.Push(0xAB)

	a := stacktest.Pop()
	b := stacktest.Pop()

	if stacktest.sp != 0 {
		t.Errorf("got SP: %d, want SP: %d", stacktest.sp, 0)
	}

	if a != 0xAB {
		t.Errorf("got stack.pop(): 0x%02x, want stack.pop(): 0x%02x\n", a, 0xAB)
	}

	if b != 0xFF {
		t.Errorf("got stack.pop(): 0x%02x, want stack.pop(): 0x%02x\n", b, 0xFF)
	}
}

func TestStackPop_zsp(t *testing.T) {
	a := stacktest.Pop()
	b := stacktest.Pop()

	if stacktest.sp != 0 {
		t.Errorf("got SP: %d, want SP: %d", stacktest.sp, 0)
	}

	if a != 0x00 {
		t.Errorf("got stack.pop(): 0x%02x, want stack.pop(): 0x%02x\n", a, 0x00)
	}

	if b != 0x00 {
		t.Errorf("got stack.pop(): 0x%02x, want stack.pop(): 0x%02x\n", b, 0x00)
	}
}
