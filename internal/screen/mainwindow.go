package screen

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Chip8Screen interface {
	SetPixel(x, y byte)
	GetPixel(x, y byte) byte
	Clear()
}

type MainWindow struct {
	window          *sdl.Window
	renderer        *sdl.Renderer
	backgroundColor sdl.Color
	pixelColor      sdl.Color
	buffer          [64][32]byte
}

func NewMainWindow(title string, width, height int32, backgroundColor, pixelColor uint64) (*MainWindow, error) {
	var mw MainWindow

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, err
	}

	mw.window = window
	mw.renderer = renderer

	mw.backgroundColor = sdl.Color{
		R: uint8((backgroundColor & 0xFF000000) >> 24),
		G: uint8((backgroundColor & 0x00FF0000) >> 16),
		B: uint8((backgroundColor & 0x0000FF00) >> 8),
		A: uint8(backgroundColor & 0x000000FF)}

	mw.pixelColor = sdl.Color{
		R: uint8((pixelColor & 0xFF000000) >> 24),
		G: uint8((pixelColor & 0x00FF0000) >> 16),
		B: uint8((pixelColor & 0x0000FF00) >> 8),
		A: uint8(pixelColor & 0x000000FF)}

	return &mw, nil
}

func (mw *MainWindow) Render() {
	for i := byte(0); i < 32; i++ {
		for k := byte(0); k < 64; k++ {
			if mw.GetPixel(k, i) == 1 {
				mw.renderer.SetDrawColor(mw.pixelColor.R, mw.pixelColor.G, mw.pixelColor.B, mw.pixelColor.A)
			} else {
				mw.renderer.SetDrawColor(mw.backgroundColor.R, mw.backgroundColor.G, mw.backgroundColor.B, mw.backgroundColor.A)
			}

			mw.renderer.FillRect(&sdl.Rect{
				X: int32(k) * 10,
				Y: int32(i) * 10,
				W: 10,
				H: 10,
			})
		}
	}

	mw.renderer.Present()
	mw.renderer.Clear()
}

func (mw *MainWindow) Free() {
	mw.window.Destroy()
	mw.renderer.Destroy()
}

func (mw *MainWindow) SetPixel(x, y byte) {
	if x < 64 && y < 32 {
		mw.buffer[x][y] ^= 1
	}
}

func (mw *MainWindow) GetPixel(x, y byte) byte {
	if x < 64 && y < 32 {
		return mw.buffer[x][y]
	}

	return 0x00
}

func (mw *MainWindow) Clear() {
	for i := 0; i < 32; i++ {
		for k := 0; k < 64; k++ {
			mw.buffer[k][i] = 0x00
		}
	}
}
