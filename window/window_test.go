package window

import (
	"testing"

	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func newHostedWindow(bounds vio.Rect, title string, flags Flags) (*view.Group, *Window) {
	host := newPalettedHost()
	w := New(bounds, title, flags)
	host.Insert(w)
	return host, w
}

func TestWindowNewWiresFrame(t *testing.T) {
	_, w := newHostedWindow(vio.R(0, 0, 20, 6), "Welcome", FlagsAll)

	if w.Frame() == nil {
		t.Fatal("frame missing")
	}
	if w.Frame().Title != "Welcome" {
		t.Fatalf("frame title: %q", w.Frame().Title)
	}
	if w.Frame().Flags() != FlagsAll {
		t.Fatalf("frame flags: %v", w.Frame().Flags())
	}
	if w.Flags() != FlagsAll {
		t.Fatalf("window flags: %v", w.Flags())
	}
	if !w.HasState(view.StateVisible) {
		t.Fatal("window should start visible")
	}
	if w.Options&view.OptSelectable == 0 {
		t.Fatal("window must be selectable so the desktop can focus it")
	}
}

func TestWindowClientRectInsetByOne(t *testing.T) {
	_, w := newHostedWindow(vio.R(2, 3, 20, 8), "", 0)
	got := w.ClientRect()
	want := vio.R(3, 4, 18, 6)
	if got != want {
		t.Fatalf("client rect %+v want %+v", got, want)
	}
}

func TestWindowSetTitlePropagates(t *testing.T) {
	_, w := newHostedWindow(vio.R(0, 0, 20, 5), "Old", FlagsAll)
	w.SetTitle("New")
	if w.Title != "New" {
		t.Fatalf("window title: %q", w.Title)
	}
	if w.Frame().Title != "New" {
		t.Fatalf("frame title: %q", w.Frame().Title)
	}
}

func TestWindowSetStateFlipsFrameActive(t *testing.T) {
	_, w := newHostedWindow(vio.R(0, 0, 20, 5), "T", FlagsAll)

	if w.Frame().HasState(view.StateActive) {
		t.Fatal("frame should start inactive")
	}
	w.SetState(view.StateFocused, true)
	if !w.Frame().HasState(view.StateActive) {
		t.Fatal("focusing window should activate frame")
	}
	w.SetState(view.StateFocused, false)
	if w.Frame().HasState(view.StateActive) {
		t.Fatal("losing focus should deactivate frame")
	}
}

func TestWindowChangeBoundsResizesFrame(t *testing.T) {
	_, w := newHostedWindow(vio.R(0, 0, 20, 5), "", FlagsAll)
	w.ChangeBounds(vio.R(5, 5, 30, 10))
	if w.Bounds != (vio.R(5, 5, 30, 10)) {
		t.Fatalf("window bounds: %+v", w.Bounds)
	}
	if w.Frame().Bounds != (vio.R(5, 5, 30, 10)) {
		t.Fatalf("frame bounds: %+v", w.Frame().Bounds)
	}
}

func TestWindowDrawFillsClientRect(t *testing.T) {
	_, w := newHostedWindow(vio.R(0, 0, 10, 5), "Hi", FlagsAll)
	s := vio.NewSurface(10, 5)
	w.Draw(s)

	// interior cells (inset by 1) should be filled with spaces in the body
	// color, not left as zero cells.
	for y := 1; y < 4; y++ {
		for x := 1; x < 9; x++ {
			c := s.At(x, y)
			if c.Rune != ' ' {
				t.Fatalf("client cell (%d,%d) rune=%q expected space", x, y, c.Rune)
			}
		}
	}

	// the top-left corner is part of the frame border, not the client; it
	// must carry a box-drawing glyph.
	tl := s.At(0, 0).Rune
	if tl != '╔' && tl != '┌' {
		t.Fatalf("top-left should be a frame corner, got %q", tl)
	}
}

func TestWindowFlagsAccessor(t *testing.T) {
	_, w := newHostedWindow(vio.R(0, 0, 10, 5), "", FlagClose|FlagMove)
	if w.Flags() != FlagClose|FlagMove {
		t.Fatalf("flags: %v", w.Flags())
	}
}
