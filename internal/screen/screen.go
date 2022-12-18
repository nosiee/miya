package screen

import (
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Window interface {
	Render()
	Free()
}

type EmulatorWindow interface {
	Window

	SetPixel(x, y byte)
	GetPixel(x, y byte) byte
	Clear()
}

type Keyevent struct {
	Keycode sdl.Keycode
	Etype   uint32
}

var Keypressed chan Keyevent
var Debug chan string

func init() {
	Keypressed = make(chan Keyevent)
	Debug = make(chan string)
}

func ShowWindows(delay uint64, windows ...Window) {
	defer func() {
		for _, window := range windows {
			window.Free()
		}

		sdl.Quit()
	}()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {
			case *sdl.QuitEvent:
				os.Exit(0)
			case *sdl.KeyboardEvent:
				Keypressed <- Keyevent{
					Keycode: evt.Keysym.Sym,
					Etype:   evt.Type,
				}
			}

		}

		for _, window := range windows {
			window.Render()
		}

		time.Sleep(time.Millisecond * time.Duration(delay))
	}
}
