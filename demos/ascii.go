package demos

import (
	"fmt"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
	"github.com/tamnd/vigo/window"
)

// ASCIITable shows the 256-glyph code page in a 16x16 grid. Arrow
// keys move the cursor; the status line under the grid shows the
// selected codepoint as decimal, hex, and the printable rune (or '.'
// for non-printables).
type ASCIITable struct {
	*window.Window

	grid   *asciiGrid
	status *widget.StaticText
}

// ASCIIDefaultBounds gives the table a comfortable 18-column footprint.
//
//nolint:gochecknoglobals // immutable default
var ASCIIDefaultBounds = vio.R(2, 2, 18+2, 16+4)

// NewASCIITable returns an ASCIITable window sized to bounds.
func NewASCIITable(bounds vio.Rect) *ASCIITable {
	w := window.New(bounds, "ASCII Table", window.FlagMove|window.FlagClose)
	a := &ASCIITable{Window: w}

	client := w.ClientRect()
	a.grid = newASCIIGrid(vio.R(client.X+1, client.Y, 16, 16))
	w.Insert(a.grid)

	a.status = widget.NewStaticText(
		vio.R(client.X+1, client.Y+16+1, client.W-2, 1),
		"",
	)
	w.Insert(a.status)
	a.refresh()
	return a
}

// Cursor returns the current codepoint under the cursor.
func (a *ASCIITable) Cursor() byte { return a.grid.cursor }

// HandleEvent intercepts arrow keys and grid clicks; everything else
// falls through to the Window dispatch.
func (a *ASCIITable) HandleEvent(e *event.Event) {
	if e.What == event.ClassKey && a.handleKey(e) {
		a.refresh()
		e.Clear()
		return
	}
	if e.What == event.ClassMouseDown && a.handleMouse(e) {
		a.refresh()
		e.Clear()
		return
	}
	a.Window.HandleEvent(e)
}

// handleMouse moves the cursor to the codepoint under the click.
// The click is in absolute screen coordinates; the grid is 16x16
// starting at grid.Bounds.{X,Y}.
func (a *ASCIITable) handleMouse(e *event.Event) bool {
	p := vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}
	if !a.grid.Bounds.Contains(p) {
		return false
	}
	col := p.X - a.grid.Bounds.X
	row := p.Y - a.grid.Bounds.Y
	if col < 0 || col >= 16 || row < 0 || row >= 16 {
		return false
	}
	a.grid.cursor = byte(row*16 + col)
	return true
}

func (a *ASCIITable) handleKey(e *event.Event) bool {
	c := int(a.grid.cursor)
	switch e.Key.Key {
	case event.KeyArrowLeft:
		c--
	case event.KeyArrowRight:
		c++
	case event.KeyArrowUp:
		c -= 16
	case event.KeyArrowDown:
		c += 16
	case event.KeyHome:
		c = 0
	case event.KeyEnd:
		c = 255
	default:
		return false
	}
	c = (c + 256) % 256
	a.grid.cursor = byte(c)
	return true
}

func (a *ASCIITable) refresh() {
	c := a.grid.cursor
	display := rune(c)
	if !isPrintable(c) {
		display = '.'
	}
	a.status.Text = fmt.Sprintf("dec %3d  hex %02X  char %c", c, c, display)
}

func isPrintable(c byte) bool {
	if c < 0x20 {
		return false
	}
	if c == 0x7f {
		return false
	}
	if c >= 0x80 && c <= 0x9f {
		return false
	}
	return true
}

// asciiGrid paints the 16x16 grid and remembers the cursor cell.
type asciiGrid struct {
	*view.View

	cursor byte
}

func newASCIIGrid(bounds vio.Rect) *asciiGrid {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.EventMask = 0
	v.PaletteIndex = widget.SlotDialogText
	return &asciiGrid{View: v}
}

// Draw paints the 16x16 codepoint grid. Each cell is one column
// wide; the cursor cell is highlighted.
func (g *asciiGrid) Draw(s *vio.Surface) {
	normal := g.MapColor(g.PaletteIndex)
	selected := g.MapColor(widget.SlotListSelected)

	s.FillRect(g.Bounds, ' ', normal)
	for c := range 256 {
		x := g.Bounds.X + c%16
		y := g.Bounds.Y + c/16
		if x >= g.Bounds.Right() || y >= g.Bounds.Bottom() {
			continue
		}
		ch := rune(c)
		if !isPrintable(byte(c)) {
			ch = '.'
		}
		attr := normal
		if byte(c) == g.cursor {
			attr = selected
		}
		s.Set(x, y, ch, attr)
	}
}
