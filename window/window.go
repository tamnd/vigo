package window

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// minimum dimensions a window may shrink to during a resize drag. Anything
// smaller would clip the frame and leave nothing for content.
const (
	minWindowW = 16
	minWindowH = 6
)

// dragMode is the kind of in-flight pointer drag a window is tracking.
type dragMode uint8

const (
	dragNone dragMode = iota
	dragMove
	dragResize
)

// Window is a framed Group hosted on the desktop. It owns a Frame child
// that paints the border and title; the rest of the children are
// content views laid out inside the client rectangle. v0.2 adds drag,
// resize, F5 zoom, Ctrl-F4 close, and a drop shadow.
type Window struct {
	*view.Group

	Title  string
	Number int

	frame *Frame
	flags Flags

	zoomed  bool
	restore vio.Rect

	drag       dragMode
	dragOrigin vio.Point
	dragStart  vio.Rect
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
	g.SetState(view.StateShadow, true)

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

// Zoomed reports whether the window is currently zoomed to fill its owner.
func (w *Window) Zoomed() bool { return w.zoomed }

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

// Zoom toggles the window between its previous bounds and the bounds of
// its owner. Without an owner Zoom is a no-op.
func (w *Window) Zoom() {
	if w.Owner == nil {
		return
	}
	if w.zoomed {
		w.ChangeBounds(w.restore)
		w.zoomed = false
		return
	}
	w.restore = w.Bounds
	w.ChangeBounds(w.Owner.Bounds)
	w.zoomed = true
}

// Close removes the window from its owning group. Phase 3 will route this
// through a CmdClose command so dialogs can intercept (for example to
// prompt about unsaved changes); for v0.2 we close directly.
func (w *Window) Close() {
	if w.Owner == nil {
		return
	}
	w.Owner.Remove(w)
}

// HandleEvent runs window-level shortcuts (F5 zoom, Ctrl-F4 close) and
// mouse drag/resize before delegating to the embedded Group dispatch.
func (w *Window) HandleEvent(e *event.Event) {
	switch e.What {
	case event.ClassKey:
		w.handleKey(e)
	case event.ClassMouseDown:
		w.handleMouseDown(e)
	case event.ClassMouseUp:
		w.handleMouseUp(e)
	}
	if e.What == event.ClassNothing {
		return
	}
	w.Group.HandleEvent(e)
}

func (w *Window) handleKey(e *event.Event) {
	switch {
	case e.Key.Key == event.KeyF5 && e.Key.Mod == event.ModNone:
		if w.flags&FlagZoom != 0 {
			w.Zoom()
			e.Clear()
		}
	case e.Key.Key == event.KeyF4 && e.Key.Mod&event.ModCtrl != 0:
		if w.flags&FlagClose != 0 {
			w.Close()
			e.Clear()
		}
	}
}

func (w *Window) handleMouseDown(e *event.Event) {
	p := vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}

	if w.drag != dragNone {
		w.continueDrag(p)
		e.Clear()
		return
	}

	if !w.Bounds.Contains(p) {
		return
	}

	if w.flags&FlagGrow != 0 && p.X == w.Bounds.Right()-1 && p.Y == w.Bounds.Bottom()-1 {
		w.startDrag(dragResize, p)
		e.Clear()
		return
	}
	if w.flags&FlagMove != 0 && p.Y == w.Bounds.Y && !w.zoomed {
		w.startDrag(dragMove, p)
		e.Clear()
		return
	}
}

func (w *Window) handleMouseUp(e *event.Event) {
	if w.drag == dragNone {
		return
	}
	w.drag = dragNone
	w.SetState(view.StateDragging, false)
	e.Clear()
}

func (w *Window) startDrag(mode dragMode, at vio.Point) {
	w.drag = mode
	w.dragOrigin = at
	w.dragStart = w.Bounds
	w.SetState(view.StateDragging, true)
}

func (w *Window) continueDrag(at vio.Point) {
	dx := at.X - w.dragOrigin.X
	dy := at.Y - w.dragOrigin.Y
	nb := w.dragStart
	switch w.drag {
	case dragMove:
		nb.X += dx
		nb.Y += dy
	case dragResize:
		nb.W += dx
		nb.H += dy
	case dragNone:
		return
	}
	if nb.W < minWindowW {
		nb.W = minWindowW
	}
	if nb.H < minWindowH {
		nb.H = minWindowH
	}
	if w.Owner != nil {
		nb = clampToOwner(nb, w.Owner.Bounds)
	}
	w.ChangeBounds(nb)
}

// clampToOwner keeps r inside outer. For drag-move (W and H unchanged) the
// window is shifted so it stays fully visible; for drag-resize the size
// is clipped at the owner's edges.
func clampToOwner(r, outer vio.Rect) vio.Rect {
	if r.X < outer.X {
		r.X = outer.X
	}
	if r.Y < outer.Y {
		r.Y = outer.Y
	}
	if r.Right() > outer.Right() {
		if r.W <= outer.W {
			r.X = outer.Right() - r.W
		} else {
			r.W = outer.Right() - r.X
		}
	}
	if r.Bottom() > outer.Bottom() {
		if r.H <= outer.H {
			r.Y = outer.Bottom() - r.H
		} else {
			r.H = outer.Bottom() - r.Y
		}
	}
	return r
}

// Draw paints the drop shadow (when StateShadow is set), fills the client
// rectangle with the body color, and composes the frame and any content
// children on top. It also re-syncs the frame's active state from the
// window's focus, since Group.SetCurrent updates focus through the
// embedded *View and bypasses Window.SetState.
func (w *Window) Draw(s *vio.Surface) {
	w.frame.SetState(view.StateActive, w.HasState(view.StateFocused))

	if w.HasState(view.StateShadow) {
		s.DrawShadow(w.Bounds)
	}
	body := w.MapColor(w.PaletteIndex)
	s.FillRect(w.ClientRect(), ' ', body)
	w.Group.Draw(s)
}
