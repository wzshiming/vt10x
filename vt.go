package vt10x

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
)

// VT represents the virtual terminal emulator.
type VT interface {
	// Write parses input and writes terminal changes to state.
	io.Writer

	// String dumps the virtual terminal contents.
	fmt.Stringer

	// Parse blocks on read on pty or io.Reader, then parses sequences until
	// buffer empties. State is locked as soon as first rune is read, and unlocked
	// when buffer is empty.
	Parse(bf *bufio.Reader) error

	// Resize reports new size to pty and updates state.
	Resize(cols, rows int)

	// Size returns the size of the virtual terminal.
	Size() (rows, cols int)

	// Cell returns the character code, foreground color, and background
	// color at position (x, y) relative to the top left of the terminal.
	Cell(x, y int) (ch rune, fg Color, bg Color)

	// Cursor returns the current position of the cursor.
	Cursor() (x, y int)

	// CursorVisible returns the visible state of the cursor.
	CursorVisible() bool

	// Lock locks the state object's mutex.
	Lock()

	// Unlock resets change flags and unlocks the state object's mutex.
	Unlock()
}

type VTOption func(*VTInfo)

type VTInfo struct {
	w io.Writer
}

func WithWriter(w io.Writer) VTOption {
	return func(info *VTInfo) {
		info.w = w
	}
}

// New initializes a virtual terminal emulator with the target state and
// io.ReadWriter input.
func New(opts ...VTOption) VT {
	info := VTInfo{
		w: ioutil.Discard,
	}
	for _, opt := range opts {
		opt(&info)
	}
	return newVT(info)
}
