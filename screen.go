package vterm

//#include "vterm.h"
import "C"

// EnableAltscreen enables the screen alt-buffer
func (vt *VTerm) EnableAltscreen(enable bool) {
	if enable {
		C.vterm_screen_enable_altscreen((*C.VTermScreen)(vt.screen), 1)
	} else {
		C.vterm_screen_enable_altscreen((*C.VTermScreen)(vt.screen), 0)
	}
}

// TextAt retreives the text on the terminal in the given area.
func (vt *VTerm) TextAt(rect Rect) string {
	// TODO Screen.TextAt()
	return ""
}

// CellAt retreives a cell at the given location.
func (vt *VTerm) CellAt(row, col int) Cell {
	var src C.VTermScreenCell
	C.vterm_screen_get_cell(vt.screen, C.VTermPos{row: C.int(row), col: C.int(col)}, &src)
	return parseScreenCell(&src)
}

// IsEOL tells whether the location is the end of a line.
func (vt *VTerm) IsEOL(row, col int) bool {
	return C.vterm_screen_is_eol(vt.screen, C.VTermPos{row: C.int(row), col: C.int(col)}) == 1
}
