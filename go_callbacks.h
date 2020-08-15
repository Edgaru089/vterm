
#include <stdint.h>
#include <stdbool.h>
#include "vterm.h"

// if ID==0 then we're setting the output callback to NULL
void goOutputSetCallback(VTerm* vt, int id);

// if ID==0 then we're setting the output callback to NULL
// Returns the memory allocated; it must be freed somehow
// TODO set only the callbacks not nil
void* goScreenSetCallback(VTermScreen* vt, int id);

bool getValueBool(VTermValue* v);
int getValueInt(VTermValue* v);
char* getValueString(VTermValue* v);
void getValueColor(VTermValue* v, VTermColor* col);

// WHY ON EARTH would a sane man use bit fields in the year 2020????
// Well we have to unpack it in C code
typedef struct {
	bool bold, italic, blink, reverse, strike, dwl;
	uint8_t underline, font, dhl;
} UnpackedCellAttrs;

UnpackedCellAttrs unpackCellAttrs(VTermScreenCellAttrs* c);
void packCellAttrs(VTermScreenCellAttrs* c, UnpackedCellAttrs* src);
