package vterm

// Property describes the global state of a virtual terminal.
type Property int

const (
	_                 Property = iota
	PropCursorVisible           // bool, Is cursor visible?
	PropCursorBlink             // bool, Is cursor blinking?
	PropAltscreen               // bool, Is alt-screen active?
	PropTitle                   // string, Sets the vterm title
	PropIconName                // string, Icon name
	PropReverse                 // bool, Is vterm in reverse-color?
	PropCursorShape             // int(CursorShape), Shape of the cursor
	PropMouse                   // int(MouseProp), State of the mouse cursor/button
)

// CursorShape describes the shape of the cursor.
// It's a vterm Property.
type CursorShape int

const (
	_                    CursorShape = iota
	CursorShapeBlock                 // A solid block
	CursorShapeUnderline             // An underline
	CursorShapeBarLeft               // A vertial bar on the left of the cell
)

// Mouse describes the moving/button states of the mouse.
// A vterm Property.
type Mouse int

const (
	MouseNone  Mouse = iota // No action/Not using mouse
	MouseClick              // A button is pressed
	MouseDrag               // The cursor is moved when a button is pressed
	MouseMove               // The cursor is moved
)

// Underline describes cell underline type
type Underline int

const (
	UnderlineNone   Underline = iota // No underline
	UnderlineSingle                  // Single underline
	UnderlineDouble                  // Double underline
	UnderlineCurly                   // Curly underline
)

// CellAttrs describes the attributes of a screen cell
type CellAttrs struct {
	Underline    Underline // int(Underline), Has underline?/Underline Type
	Bold         bool      // bool, Uses bold(and bright) text?
	Italic       bool      // bool, Is italic?
	Blink        bool      // bool, Is blinking?
	Reverse      bool      // bool, Has background/foreground colors reversed?
	Strike       bool      // bool, Has strikeline?
	Font         int       // int, Uses alternative font? / Font number 0 to 9
	DoubleHeight int       // DECDHL double-height line flag (0=none 1=top 2=bottom)
	DoubleWidth  bool      // Is the cell on a double-width line as in DECDWL or DECDHL
}

// Cell describes a screen cell
type Cell struct {
	Chars                  []rune    // Runes inside the cell (may be more than one)
	Width                  int       // The width of the cell on the vterm (1=one cell, 2=spans across 2 cells, etc.)
	Attrs                  CellAttrs // Attribues of the cell
	Foreground, Background Color     // Colors in the cell
}

// Pos describes a position on the terminal.
type Pos struct {
	Row, Col int
}

// Rect describes a rectangle on the terminal.
type Rect struct {
	StartRow, EndRow int
	StartCol, EndCol int
}

// Value contains one of the values given.
type Value struct {
	Boolean bool
	Number  int
	String  string
	Color   *Color
}
