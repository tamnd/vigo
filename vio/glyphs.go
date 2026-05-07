package vio

// BorderSet is the rune set used to draw a rectangular frame.
type BorderSet struct {
	H, V               rune // horizontal, vertical
	TL, TR, BL, BR     rune // four corners
	LJ, RJ, TJ, BJ, IX rune // tee joins (left, right, top, bottom) and intersection
}

// Single-line and double-line border glyphs, matching the Turbo Vision
// classic frames.
//
//nolint:gochecknoglobals // glyph tables are immutable lookups
var (
	Single = BorderSet{
		H: '─', V: '│',
		TL: '┌', TR: '┐', BL: '└', BR: '┘',
		LJ: '├', RJ: '┤', TJ: '┬', BJ: '┴', IX: '┼',
	}
	Double = BorderSet{
		H: '═', V: '║',
		TL: '╔', TR: '╗', BL: '╚', BR: '╝',
		LJ: '╠', RJ: '╣', TJ: '╦', BJ: '╩', IX: '╬',
	}
	// ASCII fallback used when the terminal can't render Unicode box-drawing
	// glyphs (e.g. LANG=C with no UTF-8 support).
	ASCII = BorderSet{
		H: '-', V: '|',
		TL: '+', TR: '+', BL: '+', BR: '+',
		LJ: '+', RJ: '+', TJ: '+', BJ: '+', IX: '+',
	}
)

// Shading runes used for backgrounds, dimming, and gauges.
const (
	ShadeLight  rune = '░'
	ShadeMedium rune = '▒'
	ShadeDark   rune = '▓'
	ShadeFull   rune = '█'
)
