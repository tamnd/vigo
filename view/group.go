package view

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
)

// Group is a View that owns a list of child views, dispatches events
// to them in three phases (pre-process, focused, post-process), and
// composes their drawing onto a Surface.
//
// Groups carry an optional palette. Setting Palette on a group makes
// it the resolution point for any descendant View.MapColor lookup.
// In practice only Application sets a palette; intermediate groups
// leave it nil.
type Group struct {
	*View

	children []Viewer
	current  int
	Palette  vio.Palette
}

// NewGroup returns a Group with the given bounds and no children.
func NewGroup(bounds vio.Rect) *Group {
	return &Group{
		View:    NewView(bounds),
		current: -1,
	}
}

// Children returns the current child list. The slice is shared with
// the group; callers must not modify it directly.
func (g *Group) Children() []Viewer { return g.children }

// Current returns the focused child, or nil if no child is focused.
func (g *Group) Current() Viewer {
	if g.current < 0 || g.current >= len(g.children) {
		return nil
	}
	return g.children[g.current]
}

// SetCurrent sets the focused child to the one at index i. Pass -1 to
// clear focus. Out-of-range indices clear focus.
func (g *Group) SetCurrent(i int) {
	if old := g.Current(); old != nil {
		old.Base().SetState(StateFocused|StateSelected, false)
	}
	if i < 0 || i >= len(g.children) {
		g.current = -1
		return
	}
	g.current = i
	g.children[i].Base().SetState(StateFocused|StateSelected, true)
}

// Insert appends a child to the group. If the group has no current
// child yet and the new child is selectable, the new child becomes
// the focused one.
func (g *Group) Insert(v Viewer) {
	v.Base().Owner = g
	g.children = append(g.children, v)
	if g.current == -1 && v.Base().Options&OptSelectable != 0 {
		g.SetCurrent(len(g.children) - 1)
	}
}

// Remove drops v from the group. If v was the current child, focus
// shifts to the next selectable child, or is cleared.
func (g *Group) Remove(v Viewer) {
	idx := -1
	for i, c := range g.children {
		if c == v {
			idx = i
			break
		}
	}
	if idx < 0 {
		return
	}
	v.Base().Owner = nil
	g.children = append(g.children[:idx], g.children[idx+1:]...)
	if g.current == idx {
		g.current = -1
		for i, c := range g.children {
			if c.Base().Options&OptSelectable != 0 {
				g.SetCurrent(i)
				break
			}
		}
	} else if g.current > idx {
		g.current--
	}
}

// Draw composes all visible children onto s, in insertion order. A
// concrete subclass that needs a back-to-front Z order should
// override this.
func (g *Group) Draw(s *vio.Surface) {
	for _, c := range g.children {
		if c.Base().HasState(StateVisible) {
			c.Draw(s)
		}
	}
}

// HandleEvent dispatches e to children in three phases (pre, focused,
// post) and stops at the first child that consumes the event. Tab and
// Shift-Tab cycle focus through the selectable children when the
// focused child did not already consume them.
func (g *Group) HandleEvent(e *event.Event) {
	if e.What == event.ClassNothing {
		return
	}

	// Phase 1: pre-process.
	for _, c := range g.children {
		if c.Base().Options&OptPreProcess == 0 {
			continue
		}
		if !classMatches(c, e) {
			continue
		}
		c.HandleEvent(e)
		if e.What == event.ClassNothing {
			return
		}
	}

	// Phase 2: focused dispatch.
	if cur := g.Current(); cur != nil && classMatches(cur, e) {
		cur.HandleEvent(e)
		if e.What == event.ClassNothing {
			return
		}
	}

	// Phase 2.5: Tab focus cycling. Runs after focused dispatch so a
	// focused widget can opt to consume Tab itself (a future Memo).
	if e.What == event.ClassKey {
		if e.Key.Key == event.KeyTab {
			g.advanceFocus(1)
			e.Clear()
			return
		}
		if e.Key.Key == event.KeyShiftTab {
			g.advanceFocus(-1)
			e.Clear()
			return
		}
	}

	// Phase 3: post-process.
	for _, c := range g.children {
		if c.Base().Options&OptPostProcess == 0 {
			continue
		}
		if !classMatches(c, e) {
			continue
		}
		c.HandleEvent(e)
		if e.What == event.ClassNothing {
			return
		}
	}
}

// advanceFocus moves the current child by step positions through the
// selectable children, wrapping at the ends. step == +1 for Tab, -1
// for Shift-Tab. With fewer than two selectable children the call is
// a no-op.
func (g *Group) advanceFocus(step int) {
	n := len(g.children)
	if n == 0 {
		return
	}
	selectable := 0
	for _, c := range g.children {
		if c.Base().Options&OptSelectable != 0 {
			selectable++
		}
	}
	if selectable < 2 {
		return
	}
	start := g.current
	if start < 0 {
		start = -1
	}
	for i := 1; i <= n; i++ {
		idx := (start + step*i) % n
		if idx < 0 {
			idx += n
		}
		if g.children[idx].Base().Options&OptSelectable != 0 {
			g.SetCurrent(idx)
			return
		}
	}
}

// classMatches reports whether the view subscribes to e's class. A
// view with EventMask == 0 receives everything (a deliberate
// exception, used during bring-up before classes are wired).
func classMatches(v Viewer, e *event.Event) bool {
	mask := v.Base().EventMask
	if mask == 0 {
		return true
	}
	return mask&e.What != 0
}

// ChangeBounds resizes the group and propagates the change to
// children based on each child's GrowMode. v0.1's behavior is the
// minimum needed for the demo: GrowAll children resize to fill the
// new bounds; GrowFixed children stay put. Per-edge growth flags are
// respected by adjusting the relevant edge.
func (g *Group) ChangeBounds(r vio.Rect) {
	old := g.Bounds
	g.View.ChangeBounds(r)
	for _, c := range g.children {
		base := c.Base()
		nb := growChildBounds(base.Bounds, base.GrowMode, old, r)
		c.ChangeBounds(nb)
	}
}

func growChildBounds(b vio.Rect, mode GrowMode, oldOwner, newOwner vio.Rect) vio.Rect {
	dx := newOwner.W - oldOwner.W
	dy := newOwner.H - oldOwner.H
	if mode&GrowLoX != 0 {
		b.X += dx
	}
	if mode&GrowLoY != 0 {
		b.Y += dy
	}
	if mode&GrowHiX != 0 {
		b.W += dx
	}
	if mode&GrowHiY != 0 {
		b.H += dy
	}
	if b.W < 0 {
		b.W = 0
	}
	if b.H < 0 {
		b.H = 0
	}
	return b
}
