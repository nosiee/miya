package screen

import (
	"strings"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const DEBUG_BUTTON_X = 160
const DEBUG_BUTTON_Y = 90
const DEBUG_BUTTON_W = 60
const DEBUG_BUTTON_H = 20

type DebugWindow struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	font     *ttf.Font
}

func NewDebugWindow(title string, width, height int32) (*DebugWindow, error) {
	var dw DebugWindow

	if err := ttf.Init(); err != nil {
		return nil, err
	}

	font, err := ttf.OpenFont("assets/font.ttf", 10)
	if err != nil {
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

	dw.window = window
	dw.renderer = renderer
	dw.font = font

	return &dw, nil
}

func (dw *DebugWindow) Render() {
	for i, info := range strings.Split(<-Debug, "\n") {
		surface, _ := dw.font.RenderUTF8Solid(info, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		texture, _ := dw.renderer.CreateTextureFromSurface(surface)
		rect := sdl.Rect{
			X: 0,
			Y: int32(i * 10),
			W: surface.W,
			H: surface.H,
		}

		dw.renderer.Copy(texture, nil, &rect)
		surface.Free()
		texture.Destroy()
	}

	dw.drawNextButton()
	dw.renderer.Present()
	dw.renderer.Clear()
}

func (dw *DebugWindow) drawNextButton() {
	dw.renderer.SetDrawColor(255, 255, 255, 255)
	defer dw.renderer.SetDrawColor(0, 0, 0, 0)

	dw.renderer.FillRect(&sdl.Rect{
		X: DEBUG_BUTTON_X,
		Y: DEBUG_BUTTON_Y,
		W: DEBUG_BUTTON_W,
		H: DEBUG_BUTTON_H,
	})

	surface, _ := dw.font.RenderUTF8Solid("Next", sdl.Color{R: 0, G: 0, B: 0, A: 255})
	texture, _ := dw.renderer.CreateTextureFromSurface(surface)
	rect := sdl.Rect{
		X: DEBUG_BUTTON_X + 15,
		Y: DEBUG_BUTTON_Y + 3,
		W: surface.W,
		H: surface.H,
	}

	dw.renderer.Copy(texture, nil, &rect)
	surface.Free()
	texture.Destroy()
}

func (dw *DebugWindow) Free() {
	dw.window.Destroy()
	dw.renderer.Destroy()
}
