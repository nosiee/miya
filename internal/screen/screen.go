package screen

import (
	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

type Screen struct {
	videomode window.SfVideoMode
	context   window.SfContextSettings
	window    graphics.Struct_SS_sfRenderWindow
}

func NewScreen(width, height uint, title string) *Screen {
	var screen Screen

	screen.videomode = window.NewSfVideoMode()
	screen.videomode.SetWidth(width)
	screen.videomode.SetHeight(height)
	screen.videomode.SetBitsPerPixel(32)

	screen.context = window.NewSfContextSettings()
	screen.window = graphics.SfRenderWindow_create(screen.videomode, title, uint(window.SfResize|window.SfClose), screen.context)

	return &screen
}

func (screen *Screen) Show() {
	defer window.DeleteSfVideoMode(screen.videomode)
	defer window.DeleteSfContextSettings(screen.context)
	defer window.SfWindow_destroy(screen.window)

	ev := window.NewSfEvent()
	defer window.DeleteSfEvent(ev)

	for window.SfWindow_isOpen(screen.window) > 0 {
		for window.SfWindow_pollEvent(screen.window, ev) > 0 {
			if ev.GetEvType() == window.SfEventType(window.SfEvtClosed) {
				return
			}
		}

		graphics.SfRenderWindow_clear(screen.window, graphics.GetSfBlack())
		graphics.SfRenderWindow_display(screen.window)
	}
}

func (screen *Screen) Clear() {
	graphics.SfRenderWindow_clear(screen.window, graphics.GetSfBlack())
}
