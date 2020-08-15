
#include <stdint.h>
#include <stdbool.h>
#include "vterm.h"

void goGetColorRGB(VTermColor* c, uint8_t* r, uint8_t* g, uint8_t* b);

bool goIsColorRGB(VTermColor* c);
bool goIsColorPalette(VTermColor* c);

int goGetColorIndex(VTermColor* c);

