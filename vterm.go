package vterm

/*
#include <vterm.h>
#cgo pkg-config: vterm

inline static int _attr_bold(VTermScreenCell *cell) { return cell->attrs.bold; }
inline static int _attr_underline(VTermScreenCell *cell) { return cell->attrs.underline; }
inline static int _attr_italic(VTermScreenCell *cell) { return cell->attrs.italic; }
inline static int _attr_blink(VTermScreenCell *cell) { return cell->attrs.blink; }
inline static int _attr_reverse(VTermScreenCell *cell) { return cell->attrs.reverse; }
inline static int _attr_strike(VTermScreenCell *cell) { return cell->attrs.strike; }
inline static int _attr_font(VTermScreenCell *cell) { return cell->attrs.font; }
inline static int _attr_dwl(VTermScreenCell *cell) { return cell->attrs.dwl; }
inline static int _attr_dhl(VTermScreenCell *cell) { return cell->attrs.dhl; }
*/
import "C"
import (
	"errors"
	"image/color"
	"unsafe"
)

type Attr int

const (
	AttrNone       Attr = 0
	AttrBold            = Attr(C.VTERM_ATTR_BOLD)
	AttrUnderline       = Attr(C.VTERM_ATTR_UNDERLINE)
	AttrItalic          = Attr(C.VTERM_ATTR_ITALIC)
	AttrBlink           = Attr(C.VTERM_ATTR_BLINK)
	AttrReverse         = Attr(C.VTERM_ATTR_REVERSE)
	AttrStrike          = Attr(C.VTERM_ATTR_STRIKE)
	AttrFont            = Attr(C.VTERM_ATTR_FONT)
	AttrForeground      = Attr(C.VTERM_ATTR_FOREGROUND)
	AttrBackground      = Attr(C.VTERM_ATTR_BACKGROUND)
	AttrNAttrrs
)

type VTerm struct {
	term *C.VTerm
}

type Pos struct {
	pos C.VTermPos
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

type Attrs struct {
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

func (sc *ScreenCell) Attrs() *Attrs {
	return &Attrs{
		Bold:      int(C._attr_bold(&sc.cell)),
		Underline: int(C._attr_underline(&sc.cell)),
		Italic:    int(C._attr_italic(&sc.cell)),
		Blink:     int(C._attr_blink(&sc.cell)),
		Reverse:   int(C._attr_reverse(&sc.cell)),
		Strike:    int(C._attr_strike(&sc.cell)),
		Font:      int(C._attr_font(&sc.cell)),
		Dwl:       int(C._attr_dwl(&sc.cell)),
		Dhl:       int(C._attr_dhl(&sc.cell)),
	}
}

func New(rows, cols int) *VTerm {
	return &VTerm{
		term: C.vterm_new(C.int(rows), C.int(cols)),
	}
}

func (vt *VTerm) Close() error {
	C.vterm_free(vt.term)
	return nil
}

func (vt *VTerm) Size() (int, int) {
	var rows, cols C.int
	C.vterm_get_size(vt.term, &rows, &cols)
	return int(rows), int(cols)
}

func (vt *VTerm) SetSize(rows, cols int) {
	C.vterm_set_size(vt.term, C.int(rows), C.int(cols))
}

func (vt *VTerm) KeyboardStartPaste() {
	C.vterm_keyboard_start_paste(vt.term)
}

func (vt *VTerm) KeyboardStopPaste() {
	C.vterm_keyboard_end_paste(vt.term)
}

func (vt *VTerm) ObtainState() *State {
	return &State{
		state: C.vterm_obtain_state(vt.term),
	}
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

func (vt *VTerm) UTF8() bool {
	return C.vterm_get_utf8(vt.term) != C.int(0)
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

func (scr *Screen) EnableAltScreen(e bool) {
	var v C.int
	if e {
		v = 1
	}
	C.vterm_screen_enable_altscreen(scr.screen, v)
}

func (scr *Screen) IsEOL(pos *Pos) bool {
	return C.vterm_screen_is_eol(scr.screen, pos.pos) != C.int(0)
}

type State struct {
	state *C.VTermState
}
