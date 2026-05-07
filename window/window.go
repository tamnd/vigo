package window

import (
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Window is a framed Group hosted on the desktop. It owns a Frame child
// that paints the border and title; the rest of the children are
// content views laid out inside the client rectangle. Phase 1 only
// supports inserting a Window onto the desktop and toggling focus
// (which switches the frame palette to the active slot). Mouse drag,
// resize, zoom, close, and modal dispatch arrive in phase 2.
type Window struct {
	*view.Group

	Title  string
	Number int

	frame *Frame
	flags Flags
}

// New returns a Window with the supplied bounds, title, and flags.
// The frame matches the bounds; the client rectangle is the bounds
// inset by 1 cell on every side. Insert the returned window into a
// Desktop (or any Group) to make it visible.
func New(bounds vio.Rect, title string, flags Flags) *Window {
	g := view.NewGroup(bounds)
	g.GrowMode = view.GrowFixed
	g.Options |= view.OptSelectable
	g.PaletteIndex = 12 // window text

	w := &Window{
		Group: g,
		Title: title,
		flags: flags,
	}
	w.frame = NewFrame(bounds, flags)
	w.frame.Title = title
	w.Insert(w.frame)
	return w
}

// Frame returns the window's frame view, primarily for tests.
func (w *Window) Frame() *Frame { return w.frame }

// Flags returns the window's decoration flags.
func (w *Window) Flags() Flags { return w.flags }

// SetTitle updates the title shown in the frame.
func (w *Window) SetTitle(title string) {
	w.Title = title
	w.frame.Title = title
}

// ClientRect returns the rectangle inside the frame that content views
// should occupy. It is the window bounds inset by 1 cell on each side.
func (w *Window) ClientRect() vio.Rect { return w.Bounds.Inset(1) }

// SetState propagates focus state into the frame so it switches between
// active and inactive palette slots.
func (w *Window) SetState(bits view.State, on bool) {
	w.View.SetState(bits, on)
	if bits&(view.StateFocused|view.StateActive) != 0 {
		w.frame.SetState(view.StateActive, w.HasState(view.StateFocused))
	}
}

// ChangeBounds resizes the window and keeps the frame aligned with the
// outer rectangle.
func (w *Window) ChangeBounds(r vio.Rect) {
	w.Group.ChangeBounds(r)
	w.frame.ChangeBounds(r)
}

// Draw fills the client rectangle with the window's body color, then
// composes the frame and any content children on top.
func (w *Window) Draw(s *vio.Surface) {
	body := w.MapColor(w.PaletteIndex)
	s.FillRect(w.ClientRect(), ' ', body)
	w.Group.Draw(s)
}
