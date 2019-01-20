// +build !libvterm

package vterm

/*
#cgo pkg-config: vterm
#include <vterm.h>
*/
import "C"
import (
	"image/color"
)

func (sc *ScreenCell) Fg() color.Color {
	return color.RGBA{
		R: uint8(sc.cell.fg.red),
		G: uint8(sc.cell.fg.green),
		B: uint8(sc.cell.fg.blue),
		A: 255,
	}
}

func (sc *ScreenCell) Bg() color.Color {
	return color.RGBA{
		R: uint8(sc.cell.bg.red),
		G: uint8(sc.cell.bg.green),
		B: uint8(sc.cell.bg.blue),
		A: 255,
	}
}

func (s *State) SetDefaultColors(fg, bg color.RGBA) {
	C.vterm_state_set_default_colors(s.state, toCVtermColor(fg), toCVtermColor(bg))
}

func toCVtermColor(col color.RGBA) *C.VTermColor {
	return &C.VTermColor{C.uint8_t(col.R), C.uint8_t(col.G), C.uint8_t(col.B)}
}
