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
		R: uint8(sc.cell.fg[0]),
		G: uint8(sc.cell.fg[1]),
		B: uint8(sc.cell.fg[2]),
		A: uint8(sc.cell.fg[3]),
	}
}

func (sc *ScreenCell) Bg() color.Color {
	return color.RGBA{
		R: uint8(sc.cell.bg[0]),
		G: uint8(sc.cell.bg[1]),
		B: uint8(sc.cell.bg[2]),
		A: uint8(sc.cell.bg[3]),
	}
}

func (s *State) SetDefaultColors(fg, bg color.RGBA) {
	C.vterm_state_set_default_colors(s.state, toCVtermColor(fg), toCVtermColor(bg))
}

func toCVtermColor(col color.RGBA) *C.VTermColor {
	return &C.VTermColor{byte(col.R), byte(col.G), byte(col.B)}
}
