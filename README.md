## Package vterm

Package vterm wraps the libvterm C library to provide a virtual terminal implementation.

The Go API matches closely to the original C API and constants are interchangeable. Cell, Property and Attributes are in native Go format and converted when calling libvterm, which increases readability and reduces the number of cgo calls (but might be slower).

The Go API supports hooking output callback and VTermScreen callbacks.

### Example

```go
vt := vterm.New(80, 25)
defer vt.Free()

vt.ScreenCallback(vterm.ScreenCallback{
	Bell: func(data interface{}) bool {
		fmt.Println("Bell!!")
		return true
	},
}, nil)

// ESC [1;31m sets the foreground color to bright red
vt.Write([]byte("\x1b[1;31m呼呼呼~"))

c := vt.CellAt(0, 0)
fmt.Printf("Char: %c\n", c.Chars[0])

fr, fg, fb := vt.ConvertRGB(&c.Foreground)
fmt.Printf("Foreground: R%d G%d B%d\n", fr, fg, fb)
br, bg, bb := vt.ConvertRGB(&c.Background)
fmt.Printf("Background: R%d G%d B%d\n", br, bg, bb)

// Rings the bell!
vt.Write([]byte("\x07"))
```

Most of the VTermState object and Parser callbacks tend to be internal and are not exposed by the Go API.
(There are functions *to-be-added* about palettes and colors however)

The VTermScreen and VTermState objects are not separated from the VTerm object; All related
methods has *VTerm as the receiver.

### License

MIT (applies to all Go source, and go\_\*.h, go\_\*.c files by me)

libvterm(modified) - MIT

The original libvterm 0.1.3 included is modified; dynamic arrays are changed to alloca() and
some NULL checks are added. Get the original here: http://www.leonerd.org.uk/code/libvterm/