package vio

import "github.com/gdamore/tcell/v2"

// Color is the foreground or background color of a cell. It is type-aliased
// to tcell.Color to avoid forcing translations at call sites; vio is the only
// package that should know it is tcell underneath.
type Color = tcell.Color

// Common Borland-style palette colors, named after their CGA/EGA originals.
//
//nolint:revive // BIOS-named constants intentionally use BIOS conventions.
const (
	Black        = tcell.ColorBlack
	Blue         = tcell.ColorNavy // BIOS "blue", the dark TV desktop
	Green        = tcell.ColorGreen
	Cyan         = tcell.ColorTeal
	Red          = tcell.ColorMaroon
	Magenta      = tcell.ColorPurple
	Brown        = tcell.ColorOlive
	LightGray    = tcell.ColorSilver
	DarkGray     = tcell.ColorGray
	LightBlue    = tcell.ColorBlue
	LightGreen   = tcell.ColorLime
	LightCyan    = tcell.ColorAqua
	LightRed     = tcell.ColorRed
	LightMagenta = tcell.ColorFuchsia
	Yellow       = tcell.ColorYellow
	White        = tcell.ColorWhite
)

// Modifier is a bitmask of attribute modifiers (bold, italic, underline...).
type Modifier uint8

// Modifier flags.
const (
	ModNone Modifier = 0
	ModBold Modifier = 1 << iota
	ModItalic
	ModUnderline
	ModReverse
	ModBlink
	ModStrike
)

// Attr describes the visual attributes of a single screen cell.
type Attr struct {
	Fg, Bg Color
	Mod    Modifier
}

// NewAttr is a small convenience constructor.
func NewAttr(fg, bg Color) Attr { return Attr{Fg: fg, Bg: bg} }

// With returns a copy of a with the given modifiers added.
func (a Attr) With(mod Modifier) Attr { a.Mod |= mod; return a }

// Without returns a copy of a with the given modifiers cleared.
func (a Attr) Without(mod Modifier) Attr { a.Mod &^= mod; return a }

// Style converts the attribute to a tcell.Style.
func (a Attr) Style() tcell.Style {
	s := tcell.StyleDefault.Foreground(a.Fg).Background(a.Bg)
	if a.Mod&ModBold != 0 {
		s = s.Bold(true)
	}
	if a.Mod&ModItalic != 0 {
		s = s.Italic(true)
	}
	if a.Mod&ModUnderline != 0 {
		s = s.Underline(true)
	}
	if a.Mod&ModReverse != 0 {
		s = s.Reverse(true)
	}
	if a.Mod&ModBlink != 0 {
		s = s.Blink(true)
	}
	if a.Mod&ModStrike != 0 {
		s = s.StrikeThrough(true)
	}
	return s
}
