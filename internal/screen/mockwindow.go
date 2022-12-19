package screen

type MockWindow struct {
	buffer [64][32]byte
}

func (mw *MockWindow) SetPixel(x, y byte) {
	if x < 64 && y < 32 {
		mw.buffer[x][y] ^= 1
	}
}

func (mw *MockWindow) GetPixel(x, y byte) byte {
	if x < 64 && y < 32 {
		return mw.buffer[x][y]
	}

	return 0x00
}

func (mw *MockWindow) Clear() {
	for i := 0; i < 32; i++ {
		for k := 0; k < 64; k++ {
			mw.buffer[k][i] = 0x00
		}
	}
}
