package vterm

// Modifier is a bitmask regarding the states of the modifier keys (Ctrl, Alt, Shift).
type Modifier int

const (
	ModifierNone  Modifier = 1 << iota // No modifier
	ModifierShift                      // Shift key
	ModifierAlt                        // Alt key
	ModifierCtrl                       // Ctrl key

	ModifierAllModsMask = ModifierShift | ModifierAlt | ModifierCtrl // All the modifier keys
)

// Key is a speical key not represented by visible text.
type Key int

// All the keys
const (
	KeyNone Key = iota

	KeyEnter
	KeyTab
	KeyBackspace
	KeyEscape

	KeyUp
	KeyDown
	KeyLeft
	KeyRight

	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown

	KeyFunction0   Key = 256
	KeyFunctionMax Key = KeyFunction0 + 255
)

// Seperated to reset iota
const (
	KeyKeypad0 Key = 512 + iota
	KeyKeypad1
	KeyKeypad2
	KeyKeypad3
	KeyKeypad4
	KeyKeypad5
	KeyKeypad6
	KeyKeypad7
	KeyKeypad8
	KeyKeypad9
	KeyKeypadMultiply
	KeyKeypadPlus
	KeyKeypadComma
	KeyKeypadMinus
	KeyKeypadPeriod
	KeyKeypadDivide
	KeyKeypadEnter
	KeyKeypadEqual

	KeyMax   // Must be last
	KeyCount = KeyMax
)

// KeyFunction returns the n-th function key code.
func KeyFunction(n int) Key {
	return Key(int(KeyFunction0) + n)
}
