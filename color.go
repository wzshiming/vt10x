package vt10x

// ANSI color values
const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	LightGrey
	DarkGrey
	LightRed
	LightGreen
	LightYellow
	LightBlue
	LightMagenta
	LightCyan
	White
)

// Default colors are potentially distinct to allow for special behavior.
// For example, a transparent background. Otherwise, the simple case is to
// map default colors to another color.
const (
	DefaultFG Color = 1<<24 + iota
	DefaultBG
	DefaultCursor
)

// Color maps to the ANSI colors [0, 16) and the xterm colors [16, 256).
type Color uint32

// ANSI returns true if Color is within [0, 16).
func (c Color) ANSI() bool {
	return c < 16
}

func (c Color) RGB() (r, g, b uint8, ok bool) {
	if c < 16 {
		return 0, 0, 0, false
	}
	if c < 232 {
		c -= 16
		r = uint8(c / 36)
		g = uint8((c % 36) / 6)
		b = uint8(c % 6)
		r = r * 51
		g = g * 51
		b = b * 51
		return r, g, b, true
	}
	if c < 256 {
		g = uint8((c-232)*10 + 8)
		return g, g, g, true
	}
	if c < 1<<24 {
		return uint8(c >> 16), uint8(c >> 8), uint8(c), true
	}
	return 0, 0, 0, false
}
