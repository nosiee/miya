package vm

type opcode struct {
	opcode uint16
	x      byte
	y      byte
	n      byte
	nn     byte
	nnn    uint16
}

func newOpcode(op uint16) opcode {
	return opcode{
		opcode: op,
		x:      byte((op & 0x0F00) >> 8),
		y:      byte((op & 0x00F0) >> 4),
		n:      byte(op & 0x000F),
		nn:     byte(op & 0x00FF),
		nnn:    op & 0x0FFF,
	}
}
