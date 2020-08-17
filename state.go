package vterm

//#include "vterm.h"
//#include "go_colors.h"
import "C"

// DefaultColors returns the default foreground/background colors.
func (vt *VTerm) DefaultColors() (fg, bg *Color) {
	fg = &Color{}
	bg = &Color{}
	C.vterm_state_get_default_colors(vt.state, (*C.VTermColor)(fg), (*C.VTermColor)(bg))
	return
}

// SetDefaultColors sets the default foreground/background colors.
func (vt *VTerm) SetDefaultColors(fg, bg *Color) {
	C.vterm_state_set_default_colors(vt.state, (*C.VTermColor)(fg), (*C.VTermColor)(bg))
}

// PaletteColor gets the palette color of the given index.
func (vt *VTerm) PaletteColor(index int) (col *Color) {
	col = &Color{}
	C.vterm_state_get_palette_color(vt.state, C.int(index), (*C.VTermColor)(col))
	return
}

// SetPaletteColor sets the palette color for the given index.
func (vt *VTerm) SetPaletteColor(index int, col *Color) {
	C.vterm_state_set_palette_color(vt.state, C.int(index), (*C.VTermColor)(col))
}

// ConvertRGB converts the Color instance to RGB format (if it is not) using
// VTermState's internal palette, and then returns the RGB values.
func (vt *VTerm) ConvertRGB(c *Color) (r, g, b uint8) {
	C.vterm_state_convert_color_to_rgb(vt.state, (*C.VTermColor)(c))
	var rc, gc, bc C.uint8_t
	C.goGetColorRGB((*C.VTermColor)(c), &rc, &gc, &bc)
	return uint8(rc), uint8(gc), uint8(bc)
}

// SetBoldHighbright sets whether to use a brighter foreground on text set to bold.
func (vt *VTerm) SetBoldHighbright(enabled bool) {
	if enabled {
		C.vterm_state_set_bold_highbright(vt.state, 1)
	} else {
		C.vterm_state_set_bold_highbright(vt.state, 0)
	}
}
