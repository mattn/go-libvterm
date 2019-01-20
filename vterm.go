package vterm

/*
#cgo CFLAGS: -I${SRCDIR}/libvterm/include
#cgo LDFLAGS: ${SRCDIR}/libvterm.a
#include <vterm.h>

inline static int _attr_bold(VTermScreenCell *cell) { return cell->attrs.bold; }
inline static int _attr_underline(VTermScreenCell *cell) { return cell->attrs.underline; }
inline static int _attr_italic(VTermScreenCell *cell) { return cell->attrs.italic; }
inline static int _attr_blink(VTermScreenCell *cell) { return cell->attrs.blink; }
inline static int _attr_reverse(VTermScreenCell *cell) { return cell->attrs.reverse; }
inline static int _attr_strike(VTermScreenCell *cell) { return cell->attrs.strike; }
inline static int _attr_font(VTermScreenCell *cell) { return cell->attrs.font; }
inline static int _attr_dwl(VTermScreenCell *cell) { return cell->attrs.dwl; }
inline static int _attr_dhl(VTermScreenCell *cell) { return cell->attrs.dhl; }

int _go_handle_damage(VTermRect, void*);
int _go_handle_bell(void*);
int _go_handle_resize(int, int, void*);
int _go_handle_moverect(VTermRect, VTermRect, void*);
int _go_handle_movecursor(VTermPos, VTermPos, int, void*);

static VTermScreenCallbacks _screen_callbacks = {
  _go_handle_damage,
  _go_handle_moverect,
  _go_handle_movecursor,
  NULL,
  _go_handle_bell,
  _go_handle_resize,
  NULL,
  NULL
};

static void
_vterm_screen_set_callbacks(VTermScreen *screen, void *user) {
  vterm_screen_set_callbacks(screen, &_screen_callbacks, user);
}
*/
import "C"
import (
	"errors"
	"image/color"
	"unsafe"

	"github.com/mattn/go-pointer"
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
	term   *C.VTerm
	screen *Screen
}

type Pos struct {
	pos C.VTermPos
}

func NewPos(row, col int) *Pos {
	var pos Pos
	pos.pos.col = C.int(col)
	pos.pos.row = C.int(row)
	return &pos
}

func (pos *Pos) Col() int {
	return int(pos.pos.col)
}

func (pos *Pos) Row() int {
	return int(pos.pos.row)
}

type Rect struct {
	rect C.VTermRect
}

func (rect *Rect) StartRow() int {
	return int(rect.rect.start_row)

}

func (rect *Rect) EndRow() int {
	return int(rect.rect.end_row)
}

func (rect *Rect) StartCol() int {
	return int(rect.rect.start_col)
}

func (rect *Rect) EndCol() int {
	return int(rect.rect.end_col)
}

func NewRect(start_row, end_row, start_col, end_col int) *Rect {
	var rect Rect
	rect.rect.start_row = C.int(start_row)
	rect.rect.end_row = C.int(end_row)
	rect.rect.start_col = C.int(start_col)
	rect.rect.end_col = C.int(end_col)
	return &rect
}

type ScreenCell struct {
	cell C.VTermScreenCell
}

type ParserCallbacks struct {
	Text func([]byte, interface{}) int
	/*
	  int (*control)(unsigned char control, void *user);
	  int (*control)(unsigned char control, void *user);
	  int (*escape)(const char *bytes, size_t len, void *user);
	  int (*csi)(const char *leader, const long args[], int argcount, const char *intermed, char command, void *user);
	  int (*osc)(const char *command, size_t cmdlen, void *user);
	  int (*dcs)(const char *command, size_t cmdlen, void *user);
	  int (*resize)(int rows, int cols, void *user);
	*/
}

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

func (sc *ScreenCell) Fg() VTermColor {
	return VTermColor{sc.cell.fg}
}

func (sc *ScreenCell) Bg() VTermColor {
	return VTermColor{sc.cell.bg}
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
	term := C.vterm_new(C.int(rows), C.int(cols))
	vt := &VTerm{
		term: term,
		screen: &Screen{
			screen: C.vterm_obtain_screen(term),
		},
	}
	C._vterm_screen_set_callbacks(C.vterm_obtain_screen(term), pointer.Save(vt))
	return vt
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
	if len(b) == 0 {
		return 0, nil
	}
	return int(C.vterm_input_write(vt.term, (*C.char)(unsafe.Pointer(&b[0])), C.size_t(len(b)))), nil
}

func (vt *VTerm) ObtainScreen() *Screen {
	return vt.screen
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

	UserData     interface{}
	OnDamage     func(*Rect) int
	OnResize     func(int, int) int
	OnMoveRect   func(*Rect, *Rect) int
	OnMoveCursor func(*Pos, *Pos, bool) int
	OnBell       func() int
	/*
	  int (*settermprop)(VTermProp prop, VTermValue *val, void *user);
	  int (*sb_pushline)(int cols, const VTermScreenCell *cells, void *user);
	  int (*sb_popline)(int cols, VTermScreenCell *cells, void *user);
	*/
}

func (scr *Screen) Flush() error {
	C.vterm_screen_flush_damage(scr.screen)
	return nil // TODO
}

func (sc *Screen) GetCellAt(row, col int) (*ScreenCell, error) {
	return sc.GetCell(NewPos(row, col))
}

func (sc *Screen) GetCell(pos *Pos) (*ScreenCell, error) {
	var cell ScreenCell
	if C.vterm_screen_get_cell(sc.screen, pos.pos, &cell.cell) == 0 {
		return nil, errors.New("GetCell")
	}
	return &cell, nil
}

func (scr *Screen) GetChars(r *[]rune, rect *Rect) int {
	l := len(*r)
	buf := make([]C.uint32_t, l)
	ret := int(C.vterm_screen_get_chars(scr.screen, &buf[0], C.size_t(l), rect.rect))
	*r = make([]rune, ret)
	for i := 0; i < ret; i++ {
		(*r)[i] = rune(buf[i])
	}
	return ret
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

//export _go_handle_damage
func _go_handle_damage(rect C.VTermRect, user unsafe.Pointer) C.int {
	onDamage := pointer.Restore(user).(*VTerm).ObtainScreen().OnDamage
	if onDamage != nil {
		return C.int(onDamage(&Rect{rect: rect}))
	}
	return 0
}

//export _go_handle_bell
func _go_handle_bell(user unsafe.Pointer) C.int {
	onBell := pointer.Restore(user).(*VTerm).ObtainScreen().OnBell
	if onBell != nil {
		return C.int(onBell())
	}
	return 0
}

//export _go_handle_resize
func _go_handle_resize(row, col C.int, user unsafe.Pointer) C.int {
	onResize := pointer.Restore(user).(*VTerm).ObtainScreen().OnResize
	if onResize != nil {
		return C.int(onResize(int(row), int(col)))
	}
	return 0
}

//export _go_handle_moverect
func _go_handle_moverect(dest, src C.VTermRect, user unsafe.Pointer) C.int {
	onMoveRect := pointer.Restore(user).(*VTerm).ObtainScreen().OnMoveRect
	if onMoveRect != nil {
		return C.int(onMoveRect(&Rect{rect: dest}, &Rect{rect: src}))
	}
	return 0
}

//export _go_handle_movecursor
func _go_handle_movecursor(pos, oldpos C.VTermPos, visible C.int, user unsafe.Pointer) C.int {
	onMoveCursor := pointer.Restore(user).(*VTerm).ObtainScreen().OnMoveCursor
	if onMoveCursor != nil {
		var b bool
		if visible != C.int(0) {
			b = true
		}
		return C.int(onMoveCursor(&Pos{pos: pos}, &Pos{pos: oldpos}, b))
	}
	return 0
}
