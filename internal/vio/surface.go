package vio

import "unicode/utf8"

// Surface is an in-memory grid of Cells. It is the single drawing target the
// view tree composes onto each frame; the active Screen flushes it to the
// terminal at the end of every frame.
//
// Surfaces are not safe for concurrent mutation. The framework is single-
// threaded for UI work; background goroutines must post events instead.
type Surface struct {
	w, h  int
	cells []Cell
}

// NewSurface allocates a surface of the given size, with all cells set to the
// zero Cell (rune 0, Attr{}).
func NewSurface(w, h int) *Surface {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	return &Surface{w: w, h: h, cells: make([]Cell, w*h)}
}

// Size returns the surface dimensions.
func (s *Surface) Size() (w, h int) { return s.w, s.h }

// Resize re-allocates the surface to (w, h). The contents are reset to zero.
// Callers should follow with a Clear or full re-draw.
func (s *Surface) Resize(w, h int) {
	if w < 0 {
		w = 0
	}
	if h < 0 {
		h = 0
	}
	if w == s.w && h == s.h {
		clear(s.cells)
		return
	}
	s.w, s.h = w, h
	s.cells = make([]Cell, w*h)
}

// Clear fills the entire surface with a space character in attr.
func (s *Surface) Clear(attr Attr) {
	c := Cell{Rune: ' ', Attr: attr}
	for i := range s.cells {
		s.cells[i] = c
	}
}

// At returns the cell at (x, y). If (x, y) is out of bounds, the zero Cell is
// returned.
func (s *Surface) At(x, y int) Cell {
	if x < 0 || y < 0 || x >= s.w || y >= s.h {
		return Cell{}
	}
	return s.cells[y*s.w+x]
}

// Set writes a single cell. Out-of-bounds writes are silently dropped, which
// matches Turbo Vision's clipping behavior and removes a class of off-by-one
// crashes.
func (s *Surface) Set(x, y int, ch rune, attr Attr) {
	if x < 0 || y < 0 || x >= s.w || y >= s.h {
		return
	}
	s.cells[y*s.w+x] = Cell{Rune: ch, Attr: attr}
}

// FillRect fills r with ch and attr, clipped to the surface.
func (s *Surface) FillRect(r Rect, ch rune, attr Attr) {
	r = r.Intersect(Rect{X: 0, Y: 0, W: s.w, H: s.h})
	if r.IsEmpty() {
		return
	}
	c := Cell{Rune: ch, Attr: attr}
	for y := r.Y; y < r.Bottom(); y++ {
		row := s.cells[y*s.w+r.X : y*s.w+r.Right()]
		for i := range row {
			row[i] = c
		}
	}
}

// DrawString writes str starting at (x, y) and stops at the right edge of the
// surface. It returns the number of runes drawn.
func (s *Surface) DrawString(x, y int, str string, attr Attr) int {
	if y < 0 || y >= s.h {
		return 0
	}
	count := 0
	for _, r := range str {
		if x >= s.w {
			break
		}
		if x >= 0 {
			s.Set(x, y, r, attr)
		}
		x++
		count += utf8.RuneLen(r)
	}
	return count
}

// Box draws a rectangular frame around r using the supplied border set and
// attribute. The frame is drawn on the boundary cells of r; r's interior is
// untouched.
func (s *Surface) Box(r Rect, b BorderSet, attr Attr) {
	if r.W < 2 || r.H < 2 {
		return
	}
	right, bottom := r.Right()-1, r.Bottom()-1
	for x := r.X + 1; x < right; x++ {
		s.Set(x, r.Y, b.H, attr)
		s.Set(x, bottom, b.H, attr)
	}
	for y := r.Y + 1; y < bottom; y++ {
		s.Set(r.X, y, b.V, attr)
		s.Set(right, y, b.V, attr)
	}
	s.Set(r.X, r.Y, b.TL, attr)
	s.Set(right, r.Y, b.TR, attr)
	s.Set(r.X, bottom, b.BL, attr)
	s.Set(right, bottom, b.BR, attr)
}

// DrawShadow paints the Turbo Vision drop-shadow: two columns to the right
// and one row below r, dimmed to dark-gray-on-black.
func (s *Surface) DrawShadow(r Rect) {
	shadow := Attr{Fg: DarkGray, Bg: Black}
	right := r.Right()
	bottom := r.Bottom()
	// right edge, two columns wide
	for y := r.Y + 1; y < bottom+1 && y < s.h; y++ {
		for dx := 0; dx < 2; dx++ {
			x := right + dx
			if x < 0 || x >= s.w {
				continue
			}
			c := s.At(x, y)
			if c.Rune == 0 {
				c.Rune = ' '
			}
			c.Attr = shadow
			s.Set(x, y, c.Rune, c.Attr)
		}
	}
	// bottom edge
	if bottom >= 0 && bottom < s.h {
		for x := r.X + 1; x < right; x++ {
			if x < 0 || x >= s.w {
				continue
			}
			c := s.At(x, bottom)
			if c.Rune == 0 {
				c.Rune = ' '
			}
			c.Attr = shadow
			s.Set(x, bottom, c.Rune, c.Attr)
		}
	}
}

// Snapshot returns a copy of every cell row, joined into a slice of strings.
// Cells with rune 0 are rendered as spaces. Intended only for tests and
// debugging.
func (s *Surface) Snapshot() []string {
	out := make([]string, s.h)
	buf := make([]rune, 0, s.w)
	for y := 0; y < s.h; y++ {
		buf = buf[:0]
		for x := 0; x < s.w; x++ {
			r := s.cells[y*s.w+x].Rune
			if r == 0 {
				r = ' '
			}
			buf = append(buf, r)
		}
		out[y] = string(buf)
	}
	return out
}
