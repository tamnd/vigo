package view

import (
	"github.com/tamnd/vigo/internal/event"
	"github.com/tamnd/vigo/internal/vio"
)

// Viewer is the interface every view satisfies. Concrete views embed
// *View, which provides default implementations, and override the
// methods they need.
type Viewer interface {
	Base() *View
	Draw(s *vio.Surface)
	HandleEvent(e *event.Event)
	ChangeBounds(r vio.Rect)
}

// View is the embed-friendly base of every visual primitive. It holds
// the state shared by all views (bounds, cursor, owner, options,
// state, palette index) and provides default no-op implementations of
// the Viewer methods.
type View struct {
	Owner *Group

	Bounds vio.Rect
	Cursor vio.Point

	State    State
	Options  Options
	GrowMode GrowMode
	HelpCtx  uint16

	EventMask event.Class

	// PaletteIndex is the 1-based offset into the Application palette
	// that this view uses for its primary fill color. Sub-classes may
	// derive other colors by adding a small offset.
	PaletteIndex byte
}

// NewView returns a *View with the given bounds, the visible state
// set, and the default event mask.
func NewView(bounds vio.Rect) *View {
	return &View{
		Bounds:    bounds,
		State:     StateVisible,
		EventMask: EventMaskDefault,
	}
}

// Base satisfies Viewer; concrete views embed *View, so their Base
// returns the embedded view by default.
func (v *View) Base() *View { return v }

// Draw is the default drawing routine: fill the bounds with the
// view's primary palette color.
func (v *View) Draw(s *vio.Surface) {
	idx := v.PaletteIndex
	if idx == 0 {
		idx = 1
	}
	s.FillRect(v.Bounds, ' ', v.MapColor(idx))
}

// HandleEvent is the default event handler: do nothing.
func (v *View) HandleEvent(e *event.Event) { _ = e }

// ChangeBounds replaces v.Bounds. Sub-classes that need to react to a
// resize (re-laying out children, for example) override this.
func (v *View) ChangeBounds(r vio.Rect) { v.Bounds = r }

// MapColor walks up the owner chain to the nearest Group that has a
// palette, and translates idx through it. If no palette is found
// (orphan view, or no Application above), it falls back to the
// Borland classic palette so something visible still gets drawn.
func (v *View) MapColor(idx byte) vio.Attr {
	for cur := v.Owner; cur != nil; cur = cur.Owner {
		if cur.Palette != nil {
			return cur.Palette.Map(idx)
		}
	}
	return vio.BorlandClassic.Map(idx)
}

// SetState sets or clears state bits.
func (v *View) SetState(bits State, on bool) {
	if on {
		v.State |= bits
	} else {
		v.State &^= bits
	}
}

// HasState reports whether all of the given state bits are set.
func (v *View) HasState(bits State) bool { return v.State&bits == bits }
