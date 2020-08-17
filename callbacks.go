package vterm

// Callbacks: We're only implementing OutputCallback and ScreenCallback - others somehow seems to be internal

//#include "vterm.h"
//#include "go_callbacks.h"
//typedef const VTermScreenCell ConstScreenCell; // For const VTermScreenCell* callbacks; see https://stackoverflow.com/questions/32938558/cgo-cant-find-way-to-use-callbacks-with-const-char-argument
import "C"
import (
	"bytes"
	"io"
	"unsafe"
)

// cbData holds callback user data
type cbData struct {
	output struct {
		f    OutputCallback
		data interface{}
	}
	screencb struct {
		ScreenCallback
		data interface{}
	}
}

// OutputCallback is called on new terminal output (into the input of the slave)
type OutputCallback func([]byte, interface{})

// OutputCallback sets the output callback of the vterm object.
// If set to non-nil, VTerm.Read will have no effect.
func (vt *VTerm) OutputCallback(f OutputCallback, data interface{}) {
	vt.output.f = f
	vt.output.data = data
	if f != nil {
		C.goOutputSetCallback(vt.term, C.int(vt.id))
	} else {
		C.goOutputSetCallback(vt.term, 0)
	}
}

// OutputWriteTo is a shortcut for setting a callback that writes to a io.Writer.
func (vt *VTerm) OutputWriteTo(w io.Writer) {
	vt.OutputCallback(func(b []byte, data interface{}) {
		io.Copy(data.(io.Writer), bytes.NewReader(b))
	}, w)
}

//export goOutputCallback
func goOutputCallback(s *C.char, len C.size_t, user unsafe.Pointer) {
	data := C.GoBytes(unsafe.Pointer(s), C.int(len))
	id := int(uintptr(user))
	if vt := idMap[id]; vt != nil && vt.output.f != nil {
		vt.output.f(data, vt.output.data)
	}
}

// ScreenCallback contains multiple callbacks called on screen changes.
// These functions return true if the action is successfully taken.
type ScreenCallback struct {
	Damage             func(r Rect, data interface{}) bool                     // A rectangle area on the screen has changed
	MoveRect           func(dest, src Rect, data interface{}) bool             // A rect on the screen is moved from src to dest
	MoveCursor         func(pos, oldpos Pos, data interface{}) bool            // Moves the cursor from oldpos to pos
	SetTermProp        func(prop Property, value Value, data interface{}) bool // A Property is set
	Bell               func(data interface{}) bool                             // Rings the bell
	Resize             func(rows, cols int, data interface{}) bool             // Terminal is resized
	ScrollbackPushLine func(cells []Cell, data interface{}) bool               // A line in the slice is poped from the top of terminal and pushed the scrollback stack
	ScrollbackPopLine  func(cells []Cell, data interface{}) bool               // Terminal is resized and a line is to be poped from the top of the scrollback stack into the slice; returns false if the stack is empty
}

// ScreenCallbackObj is an interface object that wraps ScreenCallback in Go manner.
type ScreenCallbackObj interface {
	ScreenCbDamage(r Rect) bool
	ScreenCbMoveRect(dest, src Rect) bool
	ScreenCbMoveCursor(pos, oldpos Pos) bool
	ScreenCbSetTermProp(prop Property, value Value) bool
	ScreenCbBell() bool
	ScreenCbResize(rows, cols int) bool
	ScreenCbScrollbackPushLine(cells []Cell) bool
	ScreenCbScrollbackPopLine(cells []Cell) bool
}

// ScreenCallback sets the screen callback of the vterm object.
// TODO Some fields in the ScreenCallback object can be nil while others are not.
func (vt *VTerm) ScreenCallback(f ScreenCallback, data interface{}) {
	vt.screencb.ScreenCallback = f
	vt.screencb.data = data
	vt.callocData = append(vt.callocData, C.goScreenSetCallback(vt.screen, C.int(vt.id)))
}

// ScreenCallbackObj sets the screen callbacks with a interface object.
func (vt *VTerm) ScreenCallbackObj(i ScreenCallbackObj) {
	vt.ScreenCallback(ScreenCallback{
		Damage: func(r Rect, data interface{}) bool {
			return data.(ScreenCallbackObj).ScreenCbDamage(r)
		},
		MoveRect: func(dest, src Rect, data interface{}) bool {
			return data.(ScreenCallbackObj).ScreenCbMoveRect(dest, src)
		},
		MoveCursor: func(pos, oldpos Pos, data interface{}) bool {
			return data.(ScreenCallbackObj).ScreenCbMoveCursor(pos, oldpos)
		},
		SetTermProp: func(prop Property, value Value, data interface{}) bool {
			return data.(ScreenCallbackObj).ScreenCbSetTermProp(prop, value)
		},
		Bell: func(data interface{}) bool {
			return data.(ScreenCallbackObj).ScreenCbBell()
		},
		Resize: func(rows, cols int, data interface{}) bool {
			return data.(ScreenCallbackObj).ScreenCbResize(rows, cols)
		},
		ScrollbackPushLine: func(cells []Cell, data interface{}) bool {
			return data.(ScreenCallbackObj).ScreenCbScrollbackPushLine(cells)
		},
		ScrollbackPopLine: func(cells []Cell, data interface{}) bool {
			return data.(ScreenCallbackObj).ScreenCbScrollbackPopLine(cells)
		},
	}, i)
}

//export goScreenDamage
func goScreenDamage(r C.VTermRect, user unsafe.Pointer) int {
	vt := idMap[int(uintptr(user))]
	if vt != nil && vt.screencb.Damage != nil {
		result := vt.screencb.Damage(
			Rect{
				StartRow: int(r.start_row),
				EndRow:   int(r.end_row),
				StartCol: int(r.start_col),
				EndCol:   int(r.end_col),
			}, vt.screencb.data)
		if result {
			return 1
		}
	}
	return 0
}

//export goScreenMoveRect
func goScreenMoveRect(dest, src C.VTermRect, user unsafe.Pointer) int {
	vt := idMap[int(uintptr(user))]
	if vt != nil && vt.screencb.MoveRect != nil {
		result := vt.screencb.MoveRect(
			Rect{
				StartRow: int(dest.start_row),
				EndRow:   int(dest.end_row),
				StartCol: int(dest.start_col),
				EndCol:   int(dest.end_col),
			}, Rect{
				StartRow: int(src.start_row),
				EndRow:   int(src.end_row),
				StartCol: int(src.start_col),
				EndCol:   int(src.end_col),
			}, vt.screencb.data)
		if result {
			return 1
		}
	}
	return 0
}

//export goScreenMoveCursor
func goScreenMoveCursor(pos, oldpos C.VTermPos, user unsafe.Pointer) int {
	vt := idMap[int(uintptr(user))]
	if vt != nil && vt.screencb.MoveCursor != nil {
		result := vt.screencb.MoveCursor(
			Pos{
				Row: int(pos.row),
				Col: int(pos.col),
			}, Pos{
				Row: int(oldpos.row),
				Col: int(oldpos.col),
			}, vt.screencb.data)
		if result {
			return 1
		}
	}
	return 0
}

func parsePropValue(prop C.VTermProp, value *C.VTermValue) Value {
	switch Property(prop) {
	case PropCursorVisible, PropCursorBlink, PropAltscreen:
		return Value{Boolean: bool(C.getValueBool(value))}
	case PropCursorShape, PropMouse:
		return Value{Number: int(C.getValueInt(value))}
	case PropTitle, PropIconName:
		return Value{String: C.GoString(C.getValueString(value))}
	}
	return Value{}
}

//export goScreenSetTermProp
func goScreenSetTermProp(prop C.VTermProp, value *C.VTermValue, user unsafe.Pointer) int {
	vt := idMap[int(uintptr(user))]
	if vt != nil && vt.screencb.SetTermProp != nil {
		result := vt.screencb.SetTermProp(Property(prop), parsePropValue(prop, value), vt.screencb.data)
		if result {
			return 1
		}
	}
	return 0
}

//export goScreenBell
func goScreenBell(user unsafe.Pointer) int {
	vt := idMap[int(uintptr(user))]
	if vt != nil && vt.screencb.Bell != nil {
		result := vt.screencb.Bell(vt.screencb.data)
		if result {
			return 1
		}
	}
	return 0
}

//export goScreenResize
func goScreenResize(rols, cols C.int, user unsafe.Pointer) int {
	vt := idMap[int(uintptr(user))]
	if vt != nil && vt.screencb.Resize != nil {
		result := vt.screencb.Resize(int(rols), int(cols), vt.screencb.data)
		if result {
			return 1
		}
	}
	return 0
}

func parseCellAttrs(c *C.VTermScreenCellAttrs) (attr CellAttrs) {
	cx := C.unpackCellAttrs(c)
	attr.Bold = bool(cx.bold)
	attr.Italic = bool(cx.italic)
	attr.Blink = bool(cx.blink)
	attr.Reverse = bool(cx.reverse)
	attr.Strike = bool(cx.strike)
	attr.DoubleWidth = bool(cx.dwl)
	attr.Underline = Underline(cx.underline)
	attr.Font = int(cx.font)
	attr.DoubleHeight = int(cx.dhl)
	return
}

func parseRunes(chars *[C.VTERM_MAX_CHARS_PER_CELL]C.uint32_t) (result []rune) {
	var len int
	for _, c := range chars {
		if c == 0 {
			break
		}
		len++
	}
	result = make([]rune, len)
	copy(result, ((*[C.VTERM_MAX_CHARS_PER_CELL]rune)(unsafe.Pointer(chars)))[:len])
	return
}

func parseScreenCell(c *C.VTermScreenCell) (cell Cell) {
	cell.Foreground = Color(c.fg)
	cell.Background = Color(c.bg)
	cell.Attrs = parseCellAttrs(&c.attrs)
	cell.Width = int(c.width)
	cell.Chars = parseRunes(&c.chars)
	return
}

//export goScreenScrollbackPushLine
func goScreenScrollbackPushLine(cols C.int, cells *C. /*ConstScreenCell*/ VTermScreenCell, user unsafe.Pointer) int {
	vt := idMap[int(uintptr(user))]
	if vt != nil && vt.screencb.ScrollbackPushLine != nil {
		ncells := make([]Cell, int(cols))
		for i := range ncells {
			ncells[i] = parseScreenCell((*C.VTermScreenCell)(unsafe.Pointer(uintptr(unsafe.Pointer(cells)) + unsafe.Sizeof(*cells)*uintptr(i))))
		}
		result := vt.screencb.ScrollbackPushLine(ncells, vt.screencb.data)

		if result {
			return 1
		}
	}
	return 0
}

func packCellAttrs(c *C.VTermScreenCellAttrs, src *CellAttrs) {
	var cx C.UnpackedCellAttrs
	cx.bold = C.bool(src.Bold)
	cx.italic = C.bool(src.Italic)
	cx.blink = C.bool(src.Blink)
	cx.reverse = C.bool(src.Reverse)
	cx.strike = C.bool(src.Strike)
	cx.dwl = C.bool(src.DoubleWidth)
	cx.underline = C.uint8_t(src.Underline)
	cx.font = C.uint8_t(src.Font)
	cx.dhl = C.uint8_t(src.DoubleHeight)
	C.packCellAttrs(c, &cx)
}

func packRunes(c *[C.VTERM_MAX_CHARS_PER_CELL]C.uint32_t, src []rune) {
	if len(src) > C.VTERM_MAX_CHARS_PER_CELL {
		src = src[:C.VTERM_MAX_CHARS_PER_CELL]
	}
	copy(((*[C.VTERM_MAX_CHARS_PER_CELL]rune)(unsafe.Pointer(c)))[:len(src)], src)
	if len(src) < C.VTERM_MAX_CHARS_PER_CELL {
		((*[C.VTERM_MAX_CHARS_PER_CELL]rune)(unsafe.Pointer(c)))[len(src)] = 0
	}
}

func packScreenCell(c *C.VTermScreenCell, src *Cell) {
	c.fg = C.VTermColor(src.Foreground)
	c.bg = C.VTermColor(src.Background)
	c.width = C.char(src.Width)
	packRunes(&c.chars, src.Chars)
	packCellAttrs(&c.attrs, &src.Attrs)
}

//export goScreenScrollbackPopLine
func goScreenScrollbackPopLine(cols C.int, cells *C.VTermScreenCell, user unsafe.Pointer) int {
	vt := idMap[int(uintptr(user))]
	if vt != nil && vt.screencb.ScrollbackPopLine != nil {
		ncells := make([]Cell, int(cols))
		result := vt.screencb.ScrollbackPopLine(ncells, vt.screencb.data)
		if result {
			for i := range ncells {
				ic := (*C.VTermScreenCell)(unsafe.Pointer(uintptr(unsafe.Pointer(cells)) + unsafe.Sizeof(*cells)*uintptr(i)))
				packScreenCell(ic, &ncells[i])
			}
			return 1
		}
	}
	return 0
}
