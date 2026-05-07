package app

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/window"
)

func newDesktopHost(t *testing.T) *Desktop {
	t.Helper()
	g := view.NewGroup(vio.R(0, 0, 80, 24))
	g.Palette = vio.BorlandClassic
	d := NewDesktop(vio.R(0, 0, 80, 24))
	g.Insert(d)
	return d
}

func TestDesktopF6CyclesWindows(t *testing.T) {
	d := newDesktopHost(t)
	a := window.New(vio.R(2, 2, 30, 8), "A", window.FlagsAll)
	b := window.New(vio.R(5, 5, 30, 8), "B", window.FlagsAll)
	c := window.New(vio.R(8, 8, 30, 8), "C", window.FlagsAll)
	d.Insert(a)
	d.Insert(b)
	d.Insert(c)

	if d.Current() != a {
		t.Fatalf("first window should be focused after Insert: %v", d.Current())
	}

	f6 := func() {
		d.HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyF6},
		})
	}

	f6()
	if d.Current() != b {
		t.Fatalf("after first F6 want b, got %v", d.Current())
	}
	f6()
	if d.Current() != c {
		t.Fatalf("after second F6 want c, got %v", d.Current())
	}
	f6()
	if d.Current() != a {
		t.Fatalf("after third F6 want a (wrap), got %v", d.Current())
	}
}

func TestDesktopF6SkipsBackground(t *testing.T) {
	d := newDesktopHost(t)
	a := window.New(vio.R(2, 2, 30, 8), "A", window.FlagsAll)
	d.Insert(a)
	// Only one selectable child: F6 should land on it again (wrap).
	d.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyF6},
	})
	if d.Current() != a {
		t.Fatal("F6 with one window should leave it focused")
	}
}

func TestDesktopMouseDownFocusesClickedWindow(t *testing.T) {
	d := newDesktopHost(t)
	a := window.New(vio.R(2, 2, 30, 8), "A", window.FlagsAll)
	b := window.New(vio.R(40, 5, 30, 8), "B", window.FlagsAll)
	d.Insert(a)
	d.Insert(b)

	// Click into B's area while A is focused.
	d.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 50, Y: 8, Buttons: event.MouseLeft},
	})
	if d.Current() != b {
		t.Fatalf("clicking B should focus B, got %v", d.Current())
	}
}

func TestDesktopMouseDownOnBackgroundKeepsFocus(t *testing.T) {
	d := newDesktopHost(t)
	a := window.New(vio.R(2, 2, 10, 5), "A", window.FlagsAll)
	d.Insert(a)
	// Click into the empty desktop, far from the window.
	d.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 70, Y: 20, Buttons: event.MouseLeft},
	})
	if d.Current() != a {
		t.Fatal("click on background should leave focus alone")
	}
}

func TestDesktopCmdCloseRemovesFocused(t *testing.T) {
	d := newDesktopHost(t)
	a := window.New(vio.R(2, 2, 30, 8), "A", window.FlagsAll)
	b := window.New(vio.R(5, 5, 30, 8), "B", window.FlagsAll)
	d.Insert(a)
	d.Insert(b)

	d.HandleEvent(&event.Event{
		What: event.ClassCommand,
		Msg:  event.MessageEvent{Command: event.CmdClose},
	})
	for _, c := range d.Children() {
		if c == a {
			t.Fatal("focused window should be removed by CmdClose")
		}
	}
	if d.Current() != b {
		t.Fatalf("after removing a, focus should fall to b, got %v", d.Current())
	}
}

func TestDesktopCmdCloseWithoutWindowsIsNoop(t *testing.T) {
	d := newDesktopHost(t)
	d.HandleEvent(&event.Event{
		What: event.ClassCommand,
		Msg:  event.MessageEvent{Command: event.CmdClose},
	})
	// just verify nothing panicked and the background is still there
	if len(d.Children()) != 1 {
		t.Fatalf("background should remain: %d children", len(d.Children()))
	}
}
