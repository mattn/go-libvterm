package vterm

/*
#include <vterm.h>
#cgo pkg-config: vterm
*/
import "C"
import (
	"errors"
	"image/color"
	"unsafe"
)

type VTerm struct {
	term *C.VTerm
}

type Pos struct {
	pos C.VTermPos
}

type ScreenCellAttrs struct {
	Bold      int
	Underline int
	Italic    int
	Blink     int
	Reverse   int
	Strike    int
	Font      int
	Dwl       int
	Dhl       int
}

type ScreenCell struct {
	cell C.VTermScreenCell
}

func (sc *ScreenCell) Chars() []rune {
	chars := make([]rune, int(sc.cell.width))
	for i := 0; i < len(chars); i++ {
		chars[i] = rune(sc.cell.chars[i])
	}
	return chars
}

func (sc *ScreenCell) Width() int {
	return int(sc.cell.width)
}

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

/*
TODO
Attrs ScreenCellAttrs
*/

func New(rows, cols int) *VTerm {
	return &VTerm{
		term: C.vterm_new(C.int(rows), C.int(cols)),
	}
}

func (vt *VTerm) Close() error {
	C.vterm_free(vt.term)
	return nil
}

func (vt *VTerm) GetSize() (int, int) {
	var rows, cols C.int
	C.vterm_get_size(vt.term, &rows, &cols)
	return int(rows), int(cols)
}

func (vt *VTerm) Read(b []byte) (int, error) {
	curlen := C.vterm_output_read(vt.term, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	return int(curlen), nil
}

func (vt *VTerm) Write(b []byte) (int, error) {
	curlen := C.vterm_input_write(vt.term, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))
	return int(curlen), nil
}

func (vt *VTerm) ObtainScreen() *Screen {
	return &Screen{
		screen: C.vterm_obtain_screen(vt.term),
	}
}

func (vt *VTerm) SetUTF8(b bool) {
	var v C.int
	if b {
		v = 1
	}
	C.vterm_set_utf8(vt.term, v)
}

type Screen struct {
	screen *C.VTermScreen
}

func (scr *Screen) Flush() error {
	C.vterm_screen_flush_damage(scr.screen)
	return nil // TODO
}

func (sc *Screen) GetCellAt(row, col int) (*ScreenCell, error) {
	var pos Pos
	pos.pos.col = C.int(col)
	pos.pos.row = C.int(row)
	return sc.GetCell(&pos)
}

func (sc *Screen) GetCell(pos *Pos) (*ScreenCell, error) {
	var cell ScreenCell
	if C.vterm_screen_get_cell(sc.screen, pos.pos, &cell.cell) == 0 {
		return nil, errors.New("GetCell")
	}
	return &cell, nil
}

func (scr *Screen) Reset(hard bool) {
	var v C.int
	if hard {
		v = 1
	}
	C.vterm_screen_reset(scr.screen, v)
}
