//go:build libvterm
// +build libvterm

package vterm

/*
#cgo CFLAGS: -I${SRCDIR}/libvterm/include
#cgo LDFLAGS: ${SRCDIR}/libvterm/.libs/libvterm.a
#include <vterm.h>
*/
import "C"
import "image/color"

// To get the rgb value from a VTermColor instance, call state.ConvertVTermColorToRGB
type VTermColor struct {
	color C.VTermColor
}

func NewVTermColorRGB(col color.Color) VTermColor {
	var r, g, b uint8
	colRGBA, ok := col.(color.RGBA)
	if ok {
		r, g, b = colRGBA.R, colRGBA.G, colRGBA.B
	} else {
		r16, g16, b16, _ := col.RGBA()
		r = uint8(r16 >> 8)
		g = uint8(g16 >> 8)
		b = uint8(b16 >> 8)
	}
	var t C.VTermColor
	C.vterm_color_rgb(&t, C.uchar(r), C.uchar(g), C.uchar(b))
	return VTermColor{t}
}

func NewVTermColorIndexed(index uint8) VTermColor {
	var t C.VTermColor
	t[0] |= 1
	t[1] = index
	return VTermColor{t}
}

func (c *VTermColor) IsIndex() bool {
	return c.color[0]&1 > 0
}

func (c *VTermColor) IsRGB() bool {
	return c.color[0]&1 == 0
}

func (c *VTermColor) GetRGB() (r, g, b uint8, ok bool) {
	if c.IsRGB() {
		return uint8(c.color[1]), uint8(c.color[2]), uint8(c.color[3]), true
	} else {
		return 0, 0, 0, false
	}
}

func (c *VTermColor) GetIndex() (index uint8, ok bool) {
	if c.IsIndex() {
		return uint8(c.color[1]), true
	} else {
		return 0, false
	}
}
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

func (s *State) ConvertVTermColorToRGB(col VTermColor) color.RGBA {
	if col.IsRGB() {
		arr := col.color
		return color.RGBA{uint8(arr[1]), uint8(arr[2]), uint8(arr[3]), 255}
	}
	cColor := col.color
	C.vterm_state_convert_color_to_rgb(s.state, &cColor)
	return color.RGBA{uint8(cColor[1]), uint8(cColor[2]), uint8(cColor[3]), 255}
}

func (s *State) SetDefaultColors(fg, bg VTermColor) {
	C.vterm_state_set_default_colors(s.state, &fg.color, &bg.color)
}

// index between 0 and 15, 0-7 are normal colors and 8-15 are bright colors.
func (s *State) SetPaletteColor(index int, col VTermColor) {
	if index < 0 || index >= 16 {
		panic("Index out of range")
	}
	C.vterm_state_set_palette_color(s.state, C.int(index), &col.color)
}

func (s *State) GetDefaultColors() (fg, bg VTermColor) {
	c_fg := C.VTermColor{}
	c_bg := C.VTermColor{}
	C.vterm_state_get_default_colors(s.state, &c_fg, &c_bg)
	fg = VTermColor{c_fg}
	bg = VTermColor{c_bg}
	return
}

// index between 0 and 15, 0-7 are normal colors and 8-15 are bright colors.
func (s *State) GetPaletteColor(index int) VTermColor {
	if index < 0 || index >= 16 {
		panic("Index out of range")
	}
	c_color := C.VTermColor{}
	C.vterm_state_get_palette_color(s.state, C.int(index), &c_color)
	return VTermColor{c_color}
}
