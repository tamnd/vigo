package view

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
)

// recorder is a Viewer that records every event it sees, used to assert
// dispatch order.
type recorder struct {
	*View
	name string
	log  *[]string
	// consume, when true, causes HandleEvent to clear the event.
	consume bool
}

func newRecorder(name string, log *[]string) *recorder {
	r := &recorder{View: NewView(vio.R(0, 0, 1, 1)), name: name, log: log}
	r.Options |= OptSelectable
	return r
}

func (r *recorder) HandleEvent(e *event.Event) {
	*r.log = append(*r.log, r.name)
	if r.consume {
		e.Clear()
	}
}

func (r *recorder) Draw(_ *vio.Surface) {
	*r.log = append(*r.log, "draw:"+r.name)
}

func TestInsertPromotesFirstSelectable(t *testing.T) {
	g := NewGroup(vio.R(0, 0, 10, 10))
	if g.Current() != nil {
		t.Fatal("new group should have no current")
	}
	a := newRecorder("a", new([]string))
	g.Insert(a)
	if g.Current() != a {
		t.Fatal("first selectable child should be focused")
	}
	if !a.HasState(StateFocused) {
		t.Fatal("first child should have StateFocused")
	}
}

func TestInsertSkipsNonSelectable(t *testing.T) {
	g := NewGroup(vio.R(0, 0, 10, 10))
	v := NewView(vio.R(0, 0, 1, 1))
	g.Insert(v)
	if g.Current() != nil {
		t.Fatal("non-selectable child should not become current")
	}
}

func TestSetCurrentFocusTransfer(t *testing.T) {
	g := NewGroup(vio.R(0, 0, 10, 10))
	a := newRecorder("a", new([]string))
	b := newRecorder("b", new([]string))
	g.Insert(a)
	g.Insert(b)
	g.SetCurrent(1)
	if g.Current() != b {
		t.Fatal("SetCurrent did not switch focus")
	}
	if a.HasState(StateFocused) {
		t.Fatal("old focused should lose focus")
	}
	if !b.HasState(StateFocused) {
		t.Fatal("new focused should gain focus")
	}
	g.SetCurrent(-1)
	if g.Current() != nil {
		t.Fatal("-1 should clear current")
	}
	g.SetCurrent(99) // out of range clears
	if g.Current() != nil {
		t.Fatal("oob should clear current")
	}
}

func TestRemoveShiftsFocus(t *testing.T) {
	g := NewGroup(vio.R(0, 0, 10, 10))
	a := newRecorder("a", new([]string))
	b := newRecorder("b", new([]string))
	c := newRecorder("c", new([]string))
	g.Insert(a)
	g.Insert(b)
	g.Insert(c)
	g.SetCurrent(1) // focus b
	g.Remove(b)
	if g.Current() != a {
		t.Fatalf("removing focused should fall back to first selectable, got %v", g.Current())
	}
	g.Remove(NewView(vio.Rect{})) // not present
	if len(g.Children()) != 2 {
		t.Fatal("remove of missing view should be no-op")
	}
}

func TestRemoveAdjustsCurrentIndex(t *testing.T) {
	g := NewGroup(vio.R(0, 0, 10, 10))
	a := newRecorder("a", new([]string))
	b := newRecorder("b", new([]string))
	g.Insert(a)
	g.Insert(b)
	g.SetCurrent(1)
	g.Remove(a) // a is before current; current should shift to 0
	if g.Current() != b {
		t.Fatal("focus should follow b across removal of earlier sibling")
	}
}

func TestHandleEventDispatchOrder(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 10, 10))

	pre := newRecorder("pre", log)
	pre.Options |= OptPreProcess
	mid := newRecorder("mid", log)
	post := newRecorder("post", log)
	post.Options |= OptPostProcess

	g.Insert(pre)
	g.Insert(mid)
	g.Insert(post)
	g.SetCurrent(1) // focus mid

	e := &event.Event{What: event.ClassKey}
	g.HandleEvent(e)

	want := []string{"pre", "mid", "post"}
	if len(*log) != 3 {
		t.Fatalf("log: %v", *log)
	}
	for i, w := range want {
		if (*log)[i] != w {
			t.Errorf("log[%d]=%s want %s", i, (*log)[i], w)
		}
	}
}

func TestHandleEventStopsOnConsume(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 10, 10))
	pre := newRecorder("pre", log)
	pre.Options |= OptPreProcess
	pre.consume = true
	mid := newRecorder("mid", log)
	g.Insert(pre)
	g.Insert(mid)
	g.SetCurrent(1)

	g.HandleEvent(&event.Event{What: event.ClassKey})
	if len(*log) != 1 || (*log)[0] != "pre" {
		t.Fatalf("post-pre dispatch should stop after consume: %v", *log)
	}
}

func TestHandleEventClassMaskFilters(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 10, 10))
	a := newRecorder("a", log)
	a.EventMask = event.ClassMouseDown // does not include keys
	g.Insert(a)
	g.SetCurrent(0)
	g.HandleEvent(&event.Event{What: event.ClassKey})
	if len(*log) != 0 {
		t.Fatalf("masked child should not see event: %v", *log)
	}
}

func TestHandleEventNothingShortCircuits(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 10, 10))
	a := newRecorder("a", log)
	g.Insert(a)
	g.HandleEvent(&event.Event{What: event.ClassNothing})
	if len(*log) != 0 {
		t.Fatal("ClassNothing should not be dispatched")
	}
}

func TestHandleEventZeroMaskMatchesEverything(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 10, 10))
	a := newRecorder("a", log)
	a.EventMask = 0
	g.Insert(a)
	g.SetCurrent(0)
	g.HandleEvent(&event.Event{What: event.ClassKey})
	if len(*log) != 1 {
		t.Fatalf("zero mask should match: %v", *log)
	}
}

func TestDrawSkipsInvisible(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 4, 1))
	a := newRecorder("a", log)
	b := newRecorder("b", log)
	b.SetState(StateVisible, false)
	g.Insert(a)
	g.Insert(b)
	g.Draw(vio.NewSurface(4, 1))
	if len(*log) != 1 || (*log)[0] != "draw:a" {
		t.Fatalf("invisible child drawn: %v", *log)
	}
}

func TestChangeBoundsPropagatesGrow(t *testing.T) {
	g := NewGroup(vio.R(0, 0, 10, 10))
	fixed := NewView(vio.R(1, 1, 3, 3))
	growing := NewView(vio.R(0, 0, 10, 10))
	growing.GrowMode = GrowAll
	g.Insert(fixed)
	g.Insert(growing)
	g.ChangeBounds(vio.R(0, 0, 20, 15))
	if fixed.Bounds != (vio.R(1, 1, 3, 3)) {
		t.Fatalf("fixed moved: %+v", fixed.Bounds)
	}
	if growing.Bounds != (vio.R(10, 5, 20, 15)) {
		t.Fatalf("growAll wrong: %+v", growing.Bounds)
	}
}

// TestMouseDownRoutesByHitTest guards against the regression where a
// mouse-down was dispatched only to the focused child, so a click on
// any other selectable child (e.g. a Button on a Dialog or in the
// Calculator's grid) was silently dropped. The new dispatcher hits
// the topmost selectable child whose bounds contain the click.
func TestMouseDownRoutesByHitTest(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 20, 5))
	a := newRecorder("a", log)
	a.Bounds = vio.R(0, 0, 5, 1)
	b := newRecorder("b", log)
	b.Bounds = vio.R(10, 0, 5, 1)
	g.Insert(a)
	g.Insert(b)
	g.SetCurrent(0) // a is focused

	g.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 12, Y: 0, Buttons: event.MouseLeft},
	})
	if len(*log) != 1 || (*log)[0] != "b" {
		t.Fatalf("hit-test should dispatch to b, got %v", *log)
	}
}

// TestMouseDownFallsBackToFocusWhenNoHit confirms that if no selectable
// child sits under the cursor, the dispatcher still delivers the click
// to the focused child. Without this fall-back, Desktop's focusAt and
// modal sub-loops would lose mouse-down clicks that miss every
// content child.
func TestMouseDownFallsBackToFocusWhenNoHit(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 20, 5))
	a := newRecorder("a", log)
	a.Bounds = vio.R(0, 0, 5, 1)
	g.Insert(a)
	g.SetCurrent(0)

	g.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 15, Y: 4, Buttons: event.MouseLeft},
	})
	if len(*log) != 1 || (*log)[0] != "a" {
		t.Fatalf("fallback should reach focused child, got %v", *log)
	}
}

// TestMouseDownSkipsNonSelectable ensures a frame or static text view
// cannot eat a click meant for a selectable button underneath. The
// hit-test loop only considers OptSelectable children.
func TestMouseDownSkipsNonSelectable(t *testing.T) {
	log := new([]string)
	g := NewGroup(vio.R(0, 0, 20, 5))
	frame := newRecorder("frame", log)
	frame.Options &^= OptSelectable
	frame.Bounds = vio.R(0, 0, 20, 5)
	button := newRecorder("button", log)
	button.Bounds = vio.R(2, 1, 4, 1)
	g.Insert(frame)
	g.Insert(button)

	g.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 3, Y: 1, Buttons: event.MouseLeft},
	})
	if len(*log) != 1 || (*log)[0] != "button" {
		t.Fatalf("hit-test should skip frame and reach button, got %v", *log)
	}
}

// TestChangeBoundsTranslatesChildren guards the regression where dragging
// a window moved the frame but left content children pinned to their
// original screen position. Children store absolute bounds, so a pure
// owner move (size unchanged) must shift every child by the same delta.
func TestChangeBoundsTranslatesChildren(t *testing.T) {
	g := NewGroup(vio.R(10, 5, 20, 10))
	a := NewView(vio.R(11, 6, 5, 3))
	a.GrowMode = GrowFixed
	b := NewView(vio.R(20, 10, 8, 4))
	b.GrowMode = GrowFixed
	g.Insert(a)
	g.Insert(b)
	g.ChangeBounds(vio.R(30, 8, 20, 10))
	if a.Bounds != (vio.R(31, 9, 5, 3)) {
		t.Fatalf("a not translated: %+v", a.Bounds)
	}
	if b.Bounds != (vio.R(40, 13, 8, 4)) {
		t.Fatalf("b not translated: %+v", b.Bounds)
	}
}

func TestChangeBoundsClampsNegative(t *testing.T) {
	g := NewGroup(vio.R(0, 0, 10, 10))
	v := NewView(vio.R(0, 0, 5, 5))
	v.GrowMode = GrowHiX | GrowHiY
	g.Insert(v)
	g.ChangeBounds(vio.R(0, 0, 1, 1))
	if v.Bounds.W < 0 || v.Bounds.H < 0 {
		t.Fatalf("negative bounds: %+v", v.Bounds)
	}
}
