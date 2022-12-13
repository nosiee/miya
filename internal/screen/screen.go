package screen

import (
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Screen struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	delay    uint64
	Keyevt   chan Keyevent
	buffer   [64][32]byte
}

type Keyevent struct {
	Keycode sdl.Keycode
	Etype   uint32
}

func NewScreen(width, height int32, title string, delay uint64) (*Screen, error) {
	var screen Screen
	var err error

	screen.window, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	screen.renderer, err = sdl.CreateRenderer(screen.window, -1, 0)
	if err != nil {
		return nil, err
	}

	screen.Keyevt = make(chan Keyevent)
	screen.delay = delay
	return &screen, nil
}

func (screen *Screen) Show() {
	defer screen.window.Destroy()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			case *sdl.KeyboardEvent:
				screen.Keyevt <- Keyevent{
					Keycode: evt.Keysym.Sym,
					Etype:   evt.Type,
				}
			}
		}

		screen.renderer.Clear()

		for i := byte(0); i < 32; i++ {
			for k := byte(0); k < 64; k++ {
				if screen.GetPixel(k, i) == 1 {
					screen.renderer.SetDrawColor(255, 255, 255, 0)
				} else {
					screen.renderer.SetDrawColor(0, 0, 0, 0)
				}

				screen.renderer.FillRect(&sdl.Rect{
					X: int32(k) * 10,
					Y: int32(i) * 10,
					W: 10,
					H: 10,
				})
			}
		}

		screen.renderer.Present()
		time.Sleep(time.Millisecond * time.Duration(screen.delay))
	}
}

func (screen *Screen) SetPixel(x, y byte) {
	if x < 64 && y < 32 {
		screen.buffer[x][y] ^= 1
	}
}

func (screen *Screen) GetPixel(x, y byte) byte {
	if x < 64 && y < 32 {
		return screen.buffer[x][y]
	}

	return 0x00
}

func (screen *Screen) Clear() {
	for i := 0; i < 32; i++ {
		for k := 0; k < 64; k++ {
			screen.buffer[k][i] = 0x00
		}
	}
}
