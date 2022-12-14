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

type KeyEvent struct {
	Keycode sdl.Keycode
	Etype   uint32
}

var KeyPressed chan KeyEvent
var Debug chan string
var Next chan struct{}

func init() {
	KeyPressed = make(chan KeyEvent)
	Debug = make(chan string)
	Next = make(chan struct{})
}

func ShowWindows(delay uint64, windows ...Window) {
	var quit bool

	defer func() {
		for _, window := range windows {
			window.Free()
		}

		sdl.Quit()
		os.Exit(0)
	}()

	for !quit {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch evt := event.(type) {
			case *sdl.WindowEvent:
				if evt.Event == sdl.WINDOWEVENT_CLOSE {
					quit = true
				}
			case *sdl.QuitEvent:
				quit = true
			case *sdl.MouseButtonEvent:
				// NOTE: We assume that if WindowID == 2, we are in debug mode
				if evt.WindowID == 2 && (evt.X >= DEBUG_BUTTON_X && evt.X <= (DEBUG_BUTTON_X+DEBUG_BUTTON_W)) && (evt.Y >= DEBUG_BUTTON_Y && evt.Y <= (DEBUG_BUTTON_Y+DEBUG_BUTTON_H)) {
					Next <- struct{}{}
				}
			case *sdl.KeyboardEvent:
				KeyPressed <- KeyEvent{
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
