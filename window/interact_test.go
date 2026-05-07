package window

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func mouseDown(x, y int) *event.Event {
	return &event.Event{
		What: event.ClassMouseDown,
		Mouse: event.MouseEvent{
			X: x, Y: y, Buttons: event.MouseLeft,
		},
	}
}

func mouseUp(x, y int) *event.Event {
	return &event.Event{
		What:  event.ClassMouseUp,
		Mouse: event.MouseEvent{X: x, Y: y},
	}
}

func keyEvent(k event.Key, mod event.Modifier) *event.Event {
	return &event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: k, Mod: mod},
	}
}

// hostedWindow inserts w into a sized Group host so it has an Owner with
// known bounds for clamping.
func hostedWindow(host vio.Rect, w *Window) *view.Group {
	g := view.NewGroup(host)
	g.Palette = vio.BorlandClassic
	g.Insert(w)
	return g
}

func TestWindowDragMoveFollowsMouse(t *testing.T) {
	w := New(vio.R(10, 5, 30, 8), "Move", FlagsAll)
	hostedWindow(vio.R(0, 0, 80, 24), w)

	// click on top edge, away from the close/zoom glyphs
	w.HandleEvent(mouseDown(20, 5))
	if w.drag != dragMove {
		t.Fatalf("expected dragMove, got %v", w.drag)
	}
	// drag right and down
	w.HandleEvent(mouseDown(25, 8))
	if w.Bounds != (vio.R(15, 8, 30, 8)) {
		t.Fatalf("after drag: %+v", w.Bounds)
	}
	// release
	w.HandleEvent(mouseUp(25, 8))
	if w.drag != dragNone {
		t.Fatalf("drag should clear on mouse-up: %v", w.drag)
	}
	if w.HasState(view.StateDragging) {
		t.Fatal("StateDragging should clear on mouse-up")
	}
}

func TestWindowDragMoveClampsToOwner(t *testing.T) {
	w := New(vio.R(10, 5, 30, 8), "T", FlagsAll)
	hostedWindow(vio.R(0, 0, 40, 12), w)

	w.HandleEvent(mouseDown(20, 5))
	// try to drag way off-screen to the right and down
	w.HandleEvent(mouseDown(200, 200))
	if w.Bounds.Right() > 40 || w.Bounds.Bottom() > 12 {
		t.Fatalf("clamp failed: %+v", w.Bounds)
	}
	if w.Bounds.X < 0 || w.Bounds.Y < 0 {
		t.Fatalf("clamp failed: %+v", w.Bounds)
	}
}

func TestWindowResizeFromBottomRightCorner(t *testing.T) {
	w := New(vio.R(0, 0, 30, 8), "R", FlagsAll)
	hostedWindow(vio.R(0, 0, 80, 24), w)

	corner := vio.Point{X: w.Bounds.Right() - 1, Y: w.Bounds.Bottom() - 1}
	w.HandleEvent(mouseDown(corner.X, corner.Y))
	if w.drag != dragResize {
		t.Fatalf("expected dragResize, got %v", w.drag)
	}
	w.HandleEvent(mouseDown(corner.X+10, corner.Y+4))
	if w.Bounds.W != 40 || w.Bounds.H != 12 {
		t.Fatalf("resize wrong: %+v", w.Bounds)
	}
	w.HandleEvent(mouseUp(0, 0))
}

func TestWindowResizeRespectsMinimum(t *testing.T) {
	w := New(vio.R(0, 0, 30, 8), "R", FlagsAll)
	hostedWindow(vio.R(0, 0, 80, 24), w)
	corner := vio.Point{X: w.Bounds.Right() - 1, Y: w.Bounds.Bottom() - 1}
	w.HandleEvent(mouseDown(corner.X, corner.Y))
	w.HandleEvent(mouseDown(0, 0)) // try to collapse
	if w.Bounds.W < minWindowW || w.Bounds.H < minWindowH {
		t.Fatalf("min size violated: %+v", w.Bounds)
	}
}

func TestWindowDragIgnoredWithoutFlagMove(t *testing.T) {
	w := New(vio.R(0, 0, 30, 8), "T", FlagClose|FlagZoom)
	hostedWindow(vio.R(0, 0, 80, 24), w)
	w.HandleEvent(mouseDown(15, 0))
	if w.drag != dragNone {
		t.Fatal("FlagMove off should disable drag")
	}
}

func TestWindowResizeIgnoredWithoutFlagGrow(t *testing.T) {
	w := New(vio.R(0, 0, 30, 8), "T", FlagMove|FlagClose|FlagZoom)
	hostedWindow(vio.R(0, 0, 80, 24), w)
	corner := vio.Point{X: w.Bounds.Right() - 1, Y: w.Bounds.Bottom() - 1}
	w.HandleEvent(mouseDown(corner.X, corner.Y))
	if w.drag != dragNone {
		t.Fatal("FlagGrow off should disable resize")
	}
}

func TestWindowMouseDownOutsideIgnored(t *testing.T) {
	w := New(vio.R(10, 5, 20, 6), "T", FlagsAll)
	hostedWindow(vio.R(0, 0, 80, 24), w)
	w.HandleEvent(mouseDown(0, 0))
	if w.drag != dragNone {
		t.Fatal("click outside bounds should not start drag")
	}
}

func TestWindowZoomRoundTrip(t *testing.T) {
	w := New(vio.R(10, 5, 30, 8), "Z", FlagsAll)
	hostedWindow(vio.R(0, 0, 80, 24), w)

	w.HandleEvent(keyEvent(event.KeyF5, event.ModNone))
	if !w.Zoomed() {
		t.Fatal("F5 should zoom")
	}
	if w.Bounds != (vio.R(0, 0, 80, 24)) {
		t.Fatalf("zoomed bounds: %+v", w.Bounds)
	}

	w.HandleEvent(keyEvent(event.KeyF5, event.ModNone))
	if w.Zoomed() {
		t.Fatal("F5 again should restore")
	}
	if w.Bounds != (vio.R(10, 5, 30, 8)) {
		t.Fatalf("restored bounds: %+v", w.Bounds)
	}
}

func TestWindowZoomedCannotMove(t *testing.T) {
	w := New(vio.R(10, 5, 30, 8), "Z", FlagsAll)
	hostedWindow(vio.R(0, 0, 80, 24), w)
	w.Zoom()
	w.HandleEvent(mouseDown(20, 0))
	if w.drag != dragNone {
		t.Fatal("zoomed window must not start drag-move")
	}
}

func TestWindowF5IgnoredWithoutFlagZoom(t *testing.T) {
	w := New(vio.R(10, 5, 30, 8), "Z", FlagMove|FlagGrow|FlagClose)
	hostedWindow(vio.R(0, 0, 80, 24), w)
	w.HandleEvent(keyEvent(event.KeyF5, event.ModNone))
	if w.Zoomed() {
		t.Fatal("F5 should be a no-op without FlagZoom")
	}
}

func TestWindowCtrlF4RemovesFromOwner(t *testing.T) {
	w := New(vio.R(10, 5, 30, 8), "C", FlagsAll)
	host := hostedWindow(vio.R(0, 0, 80, 24), w)

	w.HandleEvent(keyEvent(event.KeyF4, event.ModCtrl))

	for _, c := range host.Children() {
		if c == w {
			t.Fatal("Ctrl-F4 should remove window from owner")
		}
	}
}

func TestWindowCtrlF4IgnoredWithoutFlagClose(t *testing.T) {
	w := New(vio.R(10, 5, 30, 8), "C", FlagMove|FlagGrow|FlagZoom)
	host := hostedWindow(vio.R(0, 0, 80, 24), w)
	w.HandleEvent(keyEvent(event.KeyF4, event.ModCtrl))

	found := false
	for _, c := range host.Children() {
		if c == w {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("FlagClose off should suppress close")
	}
}

func TestWindowCloseDirectRemoves(t *testing.T) {
	w := New(vio.R(10, 5, 30, 8), "C", FlagsAll)
	host := hostedWindow(vio.R(0, 0, 80, 24), w)
	w.Close()
	if w.Owner != nil {
		t.Fatalf("owner should be cleared after Close: %v", w.Owner)
	}
	for _, c := range host.Children() {
		if c == w {
			t.Fatal("Close should remove window from owner")
		}
	}
}

func TestWindowCloseWithoutOwnerNoop(_ *testing.T) {
	w := New(vio.R(0, 0, 30, 8), "C", FlagsAll)
	w.Close() // must not panic
}

func TestWindowDrawDrawsShadow(t *testing.T) {
	host := newPalettedHost()
	w := New(vio.R(10, 5, 20, 6), "S", FlagsAll)
	host.Insert(w)
	s := vio.NewSurface(80, 24)
	w.Draw(s)

	// shadow should have painted the two-cell column to the right of the
	// window and the row immediately below.
	if c := s.At(w.Bounds.Right(), w.Bounds.Y+1); c.Rune == 0 {
		t.Fatal("shadow column should have painted")
	}
	if c := s.At(w.Bounds.X+2, w.Bounds.Bottom()); c.Rune == 0 {
		t.Fatal("shadow row should have painted")
	}
}

func TestWindowFocusedDrawSyncsFrameActive(t *testing.T) {
	host := newPalettedHost()
	w := New(vio.R(0, 0, 20, 6), "F", FlagsAll)
	host.Insert(w)
	// first inserted selectable child became Current; that path goes through
	// Group.SetCurrent which writes to w.View directly. Window.Draw must
	// re-sync the frame's active state from w.HasState(StateFocused).
	if !w.HasState(view.StateFocused) {
		t.Fatal("first selectable child should be focused after Insert")
	}
	if w.Frame().HasState(view.StateActive) {
		t.Fatal("frame should not yet reflect active")
	}
	w.Draw(vio.NewSurface(80, 24))
	if !w.Frame().HasState(view.StateActive) {
		t.Fatal("Draw should sync focused window into active frame")
	}
}
