package vm

import "github.com/veandco/go-sdl2/sdl"

var keymap = map[sdl.Keycode]byte{
	sdl.K_1: 0x01,
	sdl.K_2: 0x02,
	sdl.K_3: 0x03,
	sdl.K_4: 0x0C,
	sdl.K_q: 0x04,
	sdl.K_w: 0x05,
	sdl.K_e: 0x06,
	sdl.K_r: 0x0D,
	sdl.K_a: 0x07,
	sdl.K_s: 0x08,
	sdl.K_d: 0x09,
	sdl.K_f: 0x0E,
	sdl.K_z: 0x0A,
	sdl.K_x: 0x00,
	sdl.K_c: 0x0B,
	sdl.K_v: 0x0F,
}
