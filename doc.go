// Package vterm wraps the libvterm C library to provide a virtual terminal implementation.
//
// The Go API matches closely to the original C API and constants are interchangeable. However,
// Cell, Properties and Attributes are in native Go format and converted when calling libvterm,
// which increases readability and reduces the number of cgo calls (but might be slower).
//
// The Go API supports hooking output callback and VTermScreen callbacks.
//
// libvterm consists of several components: VTerm, VTermState, VTermScreen and a Parser.
// The VTerm has a State and a Screen object and the State has a built-in Parser. Most of the
// VTermState object and Parser callbacks tend to be internal and are not exposed by the Go API.
// (There are functions about palettes and colors however)
//
// The VTermScreen and VTermState objects are not separated from the VTerm object; All related
// methods has *VTerm as the receiver.
//
// The original libvterm included is modified; dynamic arrays are changed to alloca() and
// some NULL checks are added. Get the original here: http://www.leonerd.org.uk/code/libvterm/
package vterm
