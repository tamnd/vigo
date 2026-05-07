package app

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Background is the dotted-fill view that shows behind every window. It is
// the visual layer of the desktop: the recognizable Borland blue speckle.
type Background struct {
	*view.View

	Pattern rune
}

// NewBackground returns a Background sized to r, drawing pattern as its fill
// rune. Pass 0 for the Turbo Vision default (the desktop dot).
func NewBackground(r vio.Rect, pattern rune) *Background {
	if pattern == 0 {
		pattern = vio.DesktopDot
	}
	v := view.NewView(r)
	v.GrowMode = view.GrowAll
	v.PaletteIndex = 1 // desktop fill
	return &Background{View: v, Pattern: pattern}
}

// Draw fills the bounds with the desktop pattern in the desktop color.
func (b *Background) Draw(s *vio.Surface) {
	s.FillRect(b.Bounds, b.Pattern, b.MapColor(b.PaletteIndex))
}

// Desktop is the Group that hosts windows. It owns a Background child that
// fills its bounds and stays behind every other child.
type Desktop struct {
	*view.Group

	bg *Background
}

// NewDesktop returns an empty desktop sized to r, with the standard dotted
// background already attached.
func NewDesktop(r vio.Rect) *Desktop {
	g := view.NewGroup(r)
	g.GrowMode = view.GrowAll
	g.Options |= view.OptSelectable
	d := &Desktop{Group: g}
	d.bg = NewBackground(r, 0)
	d.Insert(d.bg)
	return d
}

// HandleEvent intercepts F6 (cycle windows), Ctrl-F4 close commands and
// mouse-down focus changes, then defers to the standard Group dispatch.
func (d *Desktop) HandleEvent(e *event.Event) {
	switch e.What {
	case event.ClassKey:
		if e.Key.Key == event.KeyF6 && e.Key.Mod == event.ModNone {
			d.cycle()
			e.Clear()
			return
		}
	case event.ClassMouseDown:
		d.focusAt(vio.Point{X: e.Mouse.X, Y: e.Mouse.Y})
	case event.ClassCommand:
		if e.Msg.Command == event.CmdClose {
			if d.removeFocused() {
				e.Clear()
				return
			}
		}
	}
	d.Group.HandleEvent(e)
}

// cycle advances focus to the next selectable child. Non-selectable
// children (the background) are skipped.
func (d *Desktop) cycle() {
	children := d.Children()
	n := len(children)
	if n == 0 {
		return
	}
	start := -1
	cur := d.Current()
	for i, c := range children {
		if c == cur {
			start = i
			break
		}
	}
	for i := 1; i <= n; i++ {
		idx := (start + i) % n
		if children[idx].Base().Options&view.OptSelectable != 0 {
			d.SetCurrent(idx)
			return
		}
	}
}

// focusAt focuses whichever selectable child contains p. Topmost (last in
// the children slice) wins so a click on an overlapping window picks the
// one in front. Clicks that miss every window leave focus alone.
func (d *Desktop) focusAt(p vio.Point) {
	children := d.Children()
	for i := len(children) - 1; i >= 0; i-- {
		c := children[i]
		if c.Base().Options&view.OptSelectable == 0 {
			continue
		}
		if c.Base().Bounds.Contains(p) {
			if d.Current() != c {
				d.SetCurrent(i)
			}
			return
		}
	}
}

// removeFocused removes the currently focused selectable child. Returns
// true if a child was removed.
func (d *Desktop) removeFocused() bool {
	cur := d.Current()
	if cur == nil {
		return false
	}
	d.Remove(cur)
	return true
}
