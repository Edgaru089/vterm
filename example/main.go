package main

import (
	"fmt"

	"github.com/Edgaru089/vterm"
)

func main() {

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

}
