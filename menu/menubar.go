// Package menu implements vigo's top-of-screen menu bar and bottom-of-screen
// status line. v0.1 paints the bars with a fixed list of items; activation
// and pull-down menus arrive in v0.2.
package menu

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Item is a single label on the menu bar. Hotkey is the index of the
// highlighted character within Title, or -1 for none.
type Item struct {
	Title  string
	Hotkey int
}

// Bar is the horizontal strip of menu titles painted at the top of the
// screen. v0.1 only renders it; clicking and pull-downs come in v0.2.
type Bar struct {
	*view.View

	Items []Item
}

// DefaultItems is the Turbo Vision IDE menu list: the system menu marker,
// followed by the nine top-level menus.
//
//nolint:gochecknoglobals // immutable default
var DefaultItems = []Item{
	{Title: "≡", Hotkey: -1}, // system menu
	{Title: "File", Hotkey: 0},
	{Title: "Edit", Hotkey: 0},
	{Title: "Search", Hotkey: 0},
	{Title: "Run", Hotkey: 0},
	{Title: "Compile", Hotkey: 0},
	{Title: "Debug", Hotkey: 0},
	{Title: "Project", Hotkey: 0},
	{Title: "Options", Hotkey: 0},
	{Title: "Window", Hotkey: 0},
	{Title: "Help", Hotkey: 0},
}

// NewBar returns a one-row menu bar pinned to the top of an owner of width
// w, populated with the supplied items. Pass DefaultItems for the IDE list.
func NewBar(w int, items []Item) *Bar {
	v := view.NewView(vio.Rect{X: 0, Y: 0, W: w, H: 1})
	v.GrowMode = view.GrowHiX
	v.EventMask = view.EventMaskDefault
	return &Bar{View: v, Items: items}
}

// Draw paints the bar background and each item separated by spaces. The
// hotkey character is drawn in the menu-shortcut palette color.
func (b *Bar) Draw(s *vio.Surface) {
	normal := b.MapColor(2)   // menu normal
	shortcut := b.MapColor(5) // menu shortcut

	s.FillRect(b.Bounds, ' ', normal)

	x := b.Bounds.X + 1
	y := b.Bounds.Y
	for _, it := range b.Items {
		runes := []rune(it.Title)
		for i, r := range runes {
			if x >= b.Bounds.Right() {
				return
			}
			attr := normal
			if i == it.Hotkey {
				attr = shortcut
			}
			s.Set(x, y, r, attr)
			x++
		}
		x += 2 // gutter between items
	}
}

// HandleEvent is a no-op for v0.1.
func (b *Bar) HandleEvent(e *event.Event) { _ = e }
