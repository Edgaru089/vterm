
#include "go_callbacks.h"

extern void goOutputCallback(const char* s, size_t len, void* user);

extern int goScreenDamage(VTermRect rect, void* user);
extern int goScreenMoveRect(VTermRect dest, VTermRect src, void* user);
extern int goScreenMoveCursor(VTermPos pos, VTermPos oldpos, int visible, void* user);
extern int goScreenSetTermProp(VTermProp prop, VTermValue* val, void* user);
extern int goScreenBell(void* user);
extern int goScreenResize(int rows, int cols, void* user);
extern int goScreenScrollbackPushLine(int cols, const VTermScreenCell* cells, void* user);
extern int goScreenScrollbackPopLine(int cols, VTermScreenCell* cells, void* user);

void goOutputSetCallback(VTerm* vt, int id) {
	if (id != 0)
		vterm_output_set_callback(vt, &goOutputCallback, (void*)(id));
	else
		vterm_output_set_callback(vt, 0, 0);
}

void* goScreenSetCallback(VTermScreen* vt, int id) {
	if (id != 0) {
		VTermScreenCallbacks* cb = malloc(sizeof(VTermScreenCallbacks));
		cb->damage = &goScreenDamage;
		cb->moverect = &goScreenMoveRect;
		cb->movecursor = &goScreenMoveCursor;
		cb->settermprop = &goScreenSetTermProp;
		cb->bell = &goScreenBell;
		cb->resize = &goScreenResize;
		cb->sb_pushline = &goScreenScrollbackPushLine;
		cb->sb_popline = &goScreenScrollbackPopLine;
		vterm_screen_set_callbacks(vt, cb, (void*)(id));
		return (void*)(cb);
	} else {
		vterm_screen_set_callbacks(vt, 0, 0);
		return 0;
	}
}


bool getValueBool(VTermValue* v) { return v->boolean; }
int getValueInt(VTermValue* v) { return v->number; }
char* getValueString(VTermValue* v) { return v->string; }
void getValueColor(VTermValue* v, VTermColor* col) { (*col) = v->color; }

UnpackedCellAttrs unpackCellAttrs(VTermScreenCellAttrs* c) {
	UnpackedCellAttrs attr;
	attr.bold = c->bold;
	attr.italic = c->italic;
	attr.blink = c->blink;
	attr.reverse = c->reverse;
	attr.strike = c->strike;
	attr.dwl = c->dwl;
	attr.underline = c->underline;
	attr.font = c->font;
	attr.dhl = c->dhl;
	return attr;
}
void packCellAttrs(VTermScreenCellAttrs* c, UnpackedCellAttrs* src) {
	c->bold = src->bold;
	c->italic = src->italic;
	c->blink = src->blink;
	c->reverse = src->reverse;
	c->strike = src->strike;
	c->dwl = src->dwl;
	c->underline = src->underline;
	c->font = src->font;
	c->dhl = src->dhl;;
}
