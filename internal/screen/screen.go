package screen

import (
	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

type Screen struct {
	videomode window.SfVideoMode
	context   window.SfContextSettings
	window    graphics.Struct_SS_sfRenderWindow
	Keyevt    chan window.SfKeyCode
	buffer    [64][32]byte
}

func NewScreen(width, height uint, title string) *Screen {
	var screen Screen

	screen.videomode = window.NewSfVideoMode()
	screen.videomode.SetWidth(width)
	screen.videomode.SetHeight(height)
	screen.videomode.SetBitsPerPixel(32)

	screen.context = window.NewSfContextSettings()
	screen.window = graphics.SfRenderWindow_create(screen.videomode, title, uint(window.SfResize|window.SfClose), screen.context)
	screen.Keyevt = make(chan window.SfKeyCode)

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

			if ev.GetEvType() == window.SfEventType(window.SfEvtKeyPressed) || ev.GetEvType() == window.SfEventType(window.SfEvtKeyReleased) {
				screen.Keyevt <- ev.GetKey().GetCode()
			}
		}

		graphics.SfRenderWindow_clear(screen.window, graphics.GetSfBlack())
		for i := 0; i < 32; i++ {
			for k := 0; k < 64; k++ {
				rect := graphics.SfRectangleShape_create()

				size := graphics.NewSfVector2f()
				size.SetX(10)
				size.SetY(10)

				pos := graphics.NewSfVector2f()
				pos.SetX(float32(k * 10))
				pos.SetY(float32(i * 10))

				graphics.SfRectangleShape_setSize(rect, size)
				graphics.SfRectangleShape_setOutlineThickness(rect, 1)
				graphics.SfRectangleShape_setPosition(rect, pos)
				if screen.GetPixel(byte(k), byte(i)) == 1 {
					graphics.SfRectangleShape_setOutlineColor(rect, graphics.GetSfWhite())
					graphics.SfRectangleShape_setFillColor(rect, graphics.SfColor_fromRGB(255, 255, 255))
				} else {
					graphics.SfRectangleShape_setOutlineColor(rect, graphics.GetSfBlack())
					graphics.SfRectangleShape_setFillColor(rect, graphics.SfColor_fromRGB(0, 0, 0))
				}

				graphics.SfRenderWindow_drawRectangleShape(screen.window, rect, (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0)))
			}
		}

		graphics.SfRenderWindow_display(screen.window)
	}
}

func (screen *Screen) SetPixel(x, y byte) {
	screen.buffer[x][y] ^= 1
}

func (screen *Screen) GetPixel(x, y byte) byte {
	return screen.buffer[x][y]
}

func (screen *Screen) Clear() {
	for i := 0; i < 32; i++ {
		for k := 0; k < 64; k++ {
			screen.buffer[k][i] = 0x00
		}
	}
}
