package screen

import (
	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

type Screen struct {
	videomode window.SfVideoMode
	context   window.SfContextSettings
	window    graphics.Struct_SS_sfRenderWindow
	rects     []graphics.Struct_SS_sfRectangleShape
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

		screen.Clear()
		for _, rect := range screen.rects {
			graphics.SfRenderWindow_drawRectangleShape(screen.window, rect, (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0)))
		}
		graphics.SfRenderWindow_display(screen.window)
	}
}

func (screen *Screen) Clear() {
	graphics.SfRenderWindow_clear(screen.window, graphics.GetSfBlack())
}

func (screen *Screen) Draw(x, y byte, w, h uint16) {
	rect := graphics.SfRectangleShape_create()

	graphics.SfRectangleShape_setSize(rect, makeVector2(float32(w*10), float32(h*10)))
	graphics.SfRectangleShape_setOutlineThickness(rect, 1)
	graphics.SfRectangleShape_setOutlineColor(rect, graphics.GetSfWhite())
	graphics.SfRectangleShape_setFillColor(rect, graphics.SfColor_fromRGB(255, 255, 255))
	graphics.SfRectangleShape_setOrigin(rect, makeVector2(float32(w), float32(h)))
	graphics.SfRectangleShape_setPosition(rect, makeVector2(float32(x*10), float32(y*10)))

	screen.rects = append(screen.rects, rect)
}

func makeVector2(x float32, y float32) graphics.SfVector2f {
	v := graphics.NewSfVector2f()
	v.SetX(x)
	v.SetY(y)
	return v
}
