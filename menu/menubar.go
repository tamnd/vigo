// Package menu implements vigo's top-of-screen menu bar and bottom-of-screen
// status line. v0.1 paints the bars with a fixed list of items; activation
// and pull-down menus arrive in v0.2.
package menu

import (
	"unicode"

	"github.com/tamnd/vigo/cmd"
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Item is a single label on the menu bar or a row in a pull-down. Hotkey
// is the index of the highlighted character within Title, or -1 for
// none. Cmd is the command fired when the row is activated. Shortcut
// is the right-aligned accelerator label (e.g. "F3"). Children, when
// non-empty, marks this row as a submenu trigger. Sep == true renders
// as a horizontal separator and ignores every other field.
type Item struct {
	Title    string
	Hotkey   int
	Cmd      event.CommandID
	Shortcut string
	Children []Item
	Sep      bool
}

// MenuRunner runs an open MenuBox to completion. The host (typically
// the Application) installs one that wraps ExecView; the Bar invokes
// it whenever the user opens a pull-down. Tests pass a stub that
// drives the box synchronously.
//
//nolint:revive // MenuRunner mirrors the spec name and pairs with MenuBox
type MenuRunner func(*MenuBox) event.CommandID

// Bar is the horizontal strip of menu titles painted at the top of the
// screen. With Runner installed and Children populated on each Item,
// the Bar opens pull-downs on click, F10, and Alt-letter accelerators.
type Bar struct {
	*view.View

	Items   []Item
	Enabler *cmd.Enabler
	Runner  MenuRunner

	// Last command chosen from a pull-down. Hosts read it after a Bar
	// event returns to broadcast a CmdCommand on its behalf.
	LastChosen event.CommandID
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
	v.Options |= view.OptPreProcess
	v.EventMask = event.ClassKey | event.ClassMouseDown
	return &Bar{View: v, Items: items}
}

// itemX returns the column where item i's title starts on the bar.
// Items are laid out left to right with a leading 1-cell margin and a
// 2-cell gutter between titles.
func (b *Bar) itemX(i int) int {
	x := b.Bounds.X + 1
	for j, it := range b.Items {
		if j == i {
			return x
		}
		x += len([]rune(it.Title)) + 2
	}
	return x
}

// hitItem returns the bar item index at column x, or -1 if x falls in
// a gutter.
func (b *Bar) hitItem(x int) int {
	cur := b.Bounds.X + 1
	for i, it := range b.Items {
		w := len([]rune(it.Title))
		if x >= cur && x < cur+w {
			return i
		}
		cur += w + 2
	}
	return -1
}

// findHotkey returns the bar index whose Hotkey rune matches r (case
// insensitive), or -1.
func (b *Bar) findHotkey(r rune) int {
	target := unicode.ToLower(r)
	for i, it := range b.Items {
		if it.Hotkey < 0 {
			continue
		}
		runes := []rune(it.Title)
		if it.Hotkey >= len(runes) {
			continue
		}
		if unicode.ToLower(runes[it.Hotkey]) == target {
			return i
		}
	}
	return -1
}

// HandleEvent intercepts F10 and Alt-letter accelerators, plus mouse
// clicks on the bar itself, and runs the corresponding pull-down via
// Runner. It is a pre-process handler so the accelerators work even
// when focus is in a window.
func (b *Bar) HandleEvent(e *event.Event) {
	switch e.What {
	case event.ClassKey:
		b.handleKey(e)
	case event.ClassMouseDown:
		b.handleMouse(e)
	}
}

func (b *Bar) handleKey(e *event.Event) {
	k := e.Key
	if k.Key == event.KeyF10 {
		b.openAt(b.firstWithChildren())
		e.Clear()
		return
	}
	if k.Mod&event.ModAlt != 0 && k.Key == event.KeyRune {
		if idx := b.findHotkey(k.Rune); idx >= 0 {
			b.openAt(idx)
			e.Clear()
		}
	}
}

func (b *Bar) handleMouse(e *event.Event) {
	if e.Mouse.Y != b.Bounds.Y {
		return
	}
	if idx := b.hitItem(e.Mouse.X); idx >= 0 {
		b.openAt(idx)
		e.Clear()
	}
}

// openAt runs the pull-down for bar item idx and walks left/right
// when the box result is CmdNext/CmdPrev. The chosen command (if any)
// is stored on LastChosen for the host to broadcast.
func (b *Bar) openAt(idx int) {
	if b.Runner == nil || idx < 0 || idx >= len(b.Items) {
		return
	}
	for {
		items := b.Items[idx].Children
		if len(items) == 0 {
			return
		}
		x := b.itemX(idx)
		box := NewMenuBox(x, b.Bounds.Y+1, items, b.Enabler)
		b.Runner(box)
		switch box.Result() {
		case event.CmdNext:
			idx = b.nextWithChildren(idx, 1)
		case event.CmdPrev:
			idx = b.nextWithChildren(idx, -1)
		case event.CmdMenu:
			b.LastChosen = box.Chosen()
			if b.LastChosen != event.CmdNone && b.Owner != nil {
				b.Owner.HandleEvent(&event.Event{
					What: event.ClassCommand,
					Msg:  event.MessageEvent{Command: b.LastChosen},
				})
			}
			return
		default:
			return
		}
	}
}

func (b *Bar) firstWithChildren() int {
	for i, it := range b.Items {
		if len(it.Children) > 0 {
			return i
		}
	}
	return -1
}

func (b *Bar) nextWithChildren(from, step int) int {
	n := len(b.Items)
	if n == 0 {
		return -1
	}
	for k := 1; k <= n; k++ {
		idx := (from + step*k) % n
		if idx < 0 {
			idx += n
		}
		if len(b.Items[idx].Children) > 0 {
			return idx
		}
	}
	return -1
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
