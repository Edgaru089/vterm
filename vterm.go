package vterm

//#cgo CFLAGS: -std=c99 -Wno-attributes -Wno-int-to-pointer-cast
//#include "vterm.h"
import "C"

import (
	"unsafe"
)

// VTerm holds the information for a Terminal.
// It must be explcitly freed by calling VTerm.Free().
type VTerm struct {
	term   *C.VTerm
	state  *C.VTermState
	screen *C.VTermScreen
	cbData

	// To keep track of multiple VTerm objects in callbacks without
	// passing pointers across C/Go, the ID was born
	id int

	callocData []unsafe.Pointer // Pointers to memory allocated by C.malloc()
}

// To keep track of multiple VTerm objects in callbacks without
// passing pointers across C/Go, we use an ID to track VTerm objects.
// But the Map contains a copy of the pointer itself so we have to
// free objects manually
var idCnt = 0
var idMap map[int]*VTerm = make(map[int]*VTerm)

// New creates a new VTerm instance
func New(rols, cols int) *VTerm {
	idCnt++
	obj := &VTerm{term: C.vterm_new(C.int(rols), C.int(cols)), id: idCnt}
	idMap[idCnt] = obj
	obj.state = C.vterm_obtain_state(obj.term)
	obj.screen = C.vterm_obtain_screen(obj.term)
	C.vterm_state_reset(obj.state, 1)
	C.vterm_screen_reset(obj.screen, 1)
	obj.UseUTF8(true)
	return obj
}

// Free frees the VTerm instance. It must be called before disposing the object
func (vt *VTerm) Free() {
	delete(idMap, vt.id)
	C.vterm_free(vt.term)
	for _, ptr := range vt.callocData {
		if ptr != nil {
			C.free(ptr)
		}
	}
}

// GetSize gets the apperence size of the terminal.
func (vt *VTerm) GetSize() (rols, cols int) {
	var r, c C.int
	C.vterm_get_size(vt.term, &r, &c)
	return int(r), int(c)
}

// SetSize sets the size of the terminal.
func (vt *VTerm) SetSize(rols, cols int) {
	C.vterm_set_size(vt.term, C.int(rols), C.int(cols))
}

// IsUTF8 tells whether the VTerm is to use UTF-8.
func (vt *VTerm) IsUTF8() bool {
	return C.vterm_get_utf8(vt.term) != 0
}

// UseUTF8 sets whether the VTerm is to use UTF-8.
// The default is true(not the same as libvterm)
func (vt *VTerm) UseUTF8(use bool) {
	if use {
		C.vterm_set_utf8(vt.term, 1)
	} else {
		C.vterm_set_utf8(vt.term, 0)
	}
}

// Write writes UTF-8 encoded data into the input of the VTerm (from the output of the slave)
func (vt *VTerm) Write(p []byte) (int, error) {
	size := C.vterm_input_write(vt.term, (*C.char)(unsafe.Pointer(&p[0])), C.size_t(len(p)))
	return int(size), nil
}

// GetBufferSize returns the current internal buffer size.
// Has no effect when OutputCallback is set.
func (vt *VTerm) GetBufferSize() int {
	return int(C.vterm_output_get_buffer_size(vt.term))
}

// GetBufferCurrent returns the occupied space size of the internal buffer.
// Has no effect when OutputCallback is set.
func (vt *VTerm) GetBufferCurrent() int {
	return int(C.vterm_output_get_buffer_current(vt.term))
}

// GetBufferRemaining returns the remaining space size the internal buffer.
// Has no effect when OutputCallback is set.
func (vt *VTerm) GetBufferRemaining() int {
	return int(C.vterm_output_get_buffer_remaining(vt.term))
}

// Read reads from the internal output buffer VTerm keeps (into the input of the slave).
// Has no effect when OutputCallback is set.
func (vt *VTerm) Read(p []byte) (int, error) {
	size := C.vterm_output_read(vt.term, (*C.char)(unsafe.Pointer(&p[0])), C.size_t(len(p)))
	return int(size), nil
}

// KeyboardUnichar sends a single rune to the vterm.
func (vt *VTerm) KeyboardUnichar(c rune, mod Modifier) {
	C.vterm_keyboard_unichar(vt.term, C.uint32_t(c), C.VTermModifier(mod))
}

// KeyboardKey sends a speical key to the vterm.
func (vt *VTerm) KeyboardKey(key Key, mod Modifier) {
	C.vterm_keyboard_key(vt.term, C.VTermKey(key), C.VTermModifier(mod))
}

// KeyboardStartPaste enables pasting mode.
func (vt *VTerm) KeyboardStartPaste() {
	C.vterm_keyboard_start_paste(vt.term)
}

// KeyboardEndPaste disables pasting mode.
func (vt *VTerm) KeyboardEndPaste() {
	C.vterm_keyboard_end_paste(vt.term)
}

// KeyboardPasteString is a shortcut for pasting a string.
func (vt *VTerm) KeyboardPasteString(str string, mod Modifier) {
	vt.KeyboardStartPaste()
	for _, c := range str {
		vt.KeyboardUnichar(c, mod)
	}
	vt.KeyboardEndPaste()
}

// MouseMove is called when the mouse cursor is moved into another cell.
func (vt *VTerm) MouseMove(row, col int, mod Modifier) {
	C.vterm_mouse_move(vt.term, C.int(row), C.int(row), C.VTermModifier(mod))
}

// MouseButton is used to update mouse button state (pressed/released).
// Left mouse button=1; middle=2; right=3.
func (vt *VTerm) MouseButton(button int, pressed bool, mod Modifier) {
	C.vterm_mouse_button(vt.term, C.int(button), C.bool(pressed), C.VTermModifier(mod))
}
