// Package vio provides Vigo's terminal screen abstraction: cells, attributes,
// surfaces, palettes, glyph sets, and a Screen interface backed by tcell or a
// fake screen for tests.
package vio

// Point is a row/column coordinate. X grows right, Y grows down. Both are
// 0-indexed.
type Point struct {
	X, Y int
}

// Add returns p shifted by q.
func (p Point) Add(q Point) Point { return Point{X: p.X + q.X, Y: p.Y + q.Y} }

// Sub returns p shifted by -q.
func (p Point) Sub(q Point) Point { return Point{X: p.X - q.X, Y: p.Y - q.Y} }

// Rect is an axis-aligned rectangle in cell coordinates. (X, Y) is the
// top-left corner and W, H are the width and height in cells. A rect with
// W <= 0 or H <= 0 is empty.
type Rect struct {
	X, Y, W, H int
}

// R is a convenience constructor.
func R(x, y, w, h int) Rect { return Rect{X: x, Y: y, W: w, H: h} }

// Right returns the X coordinate of the column immediately past the rect.
func (r Rect) Right() int { return r.X + r.W }

// Bottom returns the Y coordinate of the row immediately past the rect.
func (r Rect) Bottom() int { return r.Y + r.H }

// IsEmpty reports whether the rect has no cells.
func (r Rect) IsEmpty() bool { return r.W <= 0 || r.H <= 0 }

// Contains reports whether p lies inside r.
func (r Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.Right() && p.Y >= r.Y && p.Y < r.Bottom()
}

// Intersect returns the largest rect contained in both r and o; the result is
// empty when there is no overlap.
func (r Rect) Intersect(o Rect) Rect {
	x := max(r.X, o.X)
	y := max(r.Y, o.Y)
	right := min(r.Right(), o.Right())
	bottom := min(r.Bottom(), o.Bottom())
	if right <= x || bottom <= y {
		return Rect{}
	}
	return Rect{X: x, Y: y, W: right - x, H: bottom - y}
}

// Translate returns r shifted by (dx, dy).
func (r Rect) Translate(dx, dy int) Rect {
	return Rect{X: r.X + dx, Y: r.Y + dy, W: r.W, H: r.H}
}

// Inset returns a rect shrunk uniformly by n cells on every side.
func (r Rect) Inset(n int) Rect {
	return Rect{X: r.X + n, Y: r.Y + n, W: r.W - 2*n, H: r.H - 2*n}
}
