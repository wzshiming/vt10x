package vt10x

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
)

// Terminal represents the virtual terminal emulator.
type Terminal interface {
	// Write parses input and writes terminal changes to state.
	io.Writer

	// String dumps the virtual terminal contents.
	fmt.Stringer

	// View displays the virtual terminal.
	View

	// Parse blocks on read on pty or io.Reader, then parses sequences until
	// buffer empties. State is locked as soon as first rune is read, and unlocked
	// when buffer is empty.
	Parse(bf *bufio.Reader) error

	// Resize reports new size to pty and updates state.
	Resize(cols, rows int)
}

// View represents the view of the virtual terminal emulator.
type View interface {
	// Size returns the size of the virtual terminal.
	Size() (rows, cols int)

	// Mode returns the current terminal mode.//
	Mode() ModeFlag

	// Title represents the title of the console window.
	Title() string

	// Cell returns the glyph containing the character code, foreground color, and
	// background color at position (x, y) relative to the top left of the terminal.
	Cell(x, y int) Glyph

	// Cursor returns the current position of the cursor.
	Cursor() Cursor

	// CursorVisible returns the visible state of the cursor.
	CursorVisible() bool

	// Lock locks the state object's mutex.
	Lock()

	// Unlock resets change flags and unlocks the state object's mutex.
	Unlock()
}

type TerminalOption func(*TerminalInfo)

type TerminalInfo struct {
	w io.Writer
}

func WithWriter(w io.Writer) TerminalOption {
	return func(info *TerminalInfo) {
		info.w = w
	}
}

// New initializes a virtual terminal emulator with the target state and
// io.ReadWriter input.
func New(opts ...TerminalOption) Terminal {
	info := TerminalInfo{
		w: ioutil.Discard,
	}
	for _, opt := range opts {
		opt(&info)
	}
	return newTerminal(info)
}
