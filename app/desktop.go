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
// rune. Pass 0 for the Turbo Vision default (a light shade block).
func NewBackground(r vio.Rect, pattern rune) *Background {
	if pattern == 0 {
		pattern = vio.ShadeMedium
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
	d := &Desktop{Group: g}
	d.bg = NewBackground(r, 0)
	d.Insert(d.bg)
	return d
}

// HandleEvent forwards the event to the standard Group dispatch. The desktop
// itself does not consume any events in v0.1; it is a pure container.
func (d *Desktop) HandleEvent(e *event.Event) { d.Group.HandleEvent(e) }
