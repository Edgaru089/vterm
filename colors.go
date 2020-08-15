package vterm

//#include "vterm.h"
//#include "go_colors.h"
import "C"

// Color is a VTerm color instance; it can be a palette color or a RGB one
type Color C.VTermColor

// NewColorIndexed constructs a new Color using the internal/standard palette.
func NewColorIndexed(pat uint8) *Color {
	var c Color
	C.vterm_color_indexed((*C.VTermColor)(&c), C.uint8_t(pat))
	return &c
}

// NewColorRGB constructs a new Color representing the given RGB values.
func NewColorRGB(r, g, b uint8) *Color {
	var c Color
	C.vterm_color_rgb((*C.VTermColor)(&c), C.uint8_t(r), C.uint8_t(g), C.uint8_t(b))
	return &c
}

// Equal tests whether the two colors are equal.
func (c *Color) Equal(cb *Color) bool {
	return C.vterm_color_is_equal((*C.VTermColor)(c), (*C.VTermColor)(cb)) == 1
}

// ConvertRGB converts the Color instance to RGB format (if it is not) using
// VTermState's internal palette, and then returns the RGB values.
func (vt *VTerm) ConvertRGB(c *Color) (r, g, b uint8) {
	C.vterm_state_convert_color_to_rgb(vt.state, (*C.VTermColor)(c))
	var rc, gc, bc C.uint8_t
	C.goGetColorRGB((*C.VTermColor)(c), &rc, &gc, &bc)
	return uint8(rc), uint8(gc), uint8(bc)
}

// IsRGB tells whether the color is a RGB one.
func (c *Color) IsRGB() bool {
	return bool(C.goIsColorRGB((*C.VTermColor)(c)))
}

// IsIndexed tells whether the color is indexed (from palette).
func (c *Color) IsIndexed() bool {
	return bool(C.goIsColorPalette((*C.VTermColor)(c)))
}

// Index returns the palette index (0~255) if the color is indexed, -1 otherwise.
func (c *Color) Index() int {
	if !c.IsIndexed() {
		return -1
	}
	return int(C.goGetColorIndex((*C.VTermColor)(c)))
}
