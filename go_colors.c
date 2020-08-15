
//#include <cassert>
#include "go_colors.h"

void goGetColorRGB(VTermColor* c, uint8_t* r, uint8_t* g, uint8_t* b) {
	//assert(VTERM_COLOR_IS_RGB(c));
	(*r) = c->rgb.red;
	(*g) = c->rgb.green;
	(*b) = c->rgb.blue;
}

bool goIsColorRGB(VTermColor* c) { return VTERM_COLOR_IS_RGB(c); }
bool goIsColorPalette(VTermColor* c) { return VTERM_COLOR_IS_INDEXED(c); }

int goGetColorIndex(VTermColor* c) { return (int)(c->indexed.idx); }

