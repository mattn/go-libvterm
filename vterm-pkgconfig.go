//go:build !libvterm
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
		R: uint8(sc.cell.fg[1]),
		G: uint8(sc.cell.fg[2]),
		B: uint8(sc.cell.fg[3]),
		A: uint8(255),
	}
}

func (sc *ScreenCell) Bg() color.Color {
	return color.RGBA{
		R: uint8(sc.cell.bg[1]),
		G: uint8(sc.cell.bg[2]),
		B: uint8(sc.cell.bg[3]),
		A: uint8(255),
	}
}

func (s *State) SetDefaultColors(fg, bg color.RGBA) {
	C.vterm_state_set_default_colors(s.state, toCVtermColor(fg), toCVtermColor(bg))
}

func toCVtermColor(col color.RGBA) *C.VTermColor {
	return &C.VTermColor{C.VTERM_COLOR_RGB, byte(col.R), byte(col.G), byte(col.B)}
}
