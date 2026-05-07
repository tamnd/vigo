package widget

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func TestListBoxSelectedDefault(t *testing.T) {
	host := newPalettedHost()
	l := NewListBox(vio.R(0, 0, 20, 5), []string{"alpha", "beta", "gamma"}, nil)
	host.Insert(l)
	if l.Selected() != "alpha" {
		t.Fatalf("default selection: %q", l.Selected())
	}
}

func TestListBoxArrowsMoveFocus(t *testing.T) {
	host := newPalettedHost()
	l := NewListBox(vio.R(0, 0, 20, 3), []string{"a", "b", "c", "d"}, nil)
	host.Insert(l)
	host.SetCurrent(0)

	l.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	if l.Focused != 1 {
		t.Fatalf("ArrowDown: focused=%d", l.Focused)
	}
	l.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnd}})
	if l.Focused != 3 {
		t.Fatalf("End: focused=%d", l.Focused)
	}
	l.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyHome}})
	if l.Focused != 0 {
		t.Fatalf("Home: focused=%d", l.Focused)
	}
}

func TestListBoxScrollsViewportWhenFocusedExits(t *testing.T) {
	host := newPalettedHost()
	items := make([]string, 10)
	for i := range items {
		items[i] = string(rune('a' + i))
	}
	l := NewListBox(vio.R(0, 0, 20, 3), items, nil)
	host.Insert(l)
	host.SetCurrent(0)

	for range 5 {
		l.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	}
	if l.Focused != 5 {
		t.Fatalf("focused=%d", l.Focused)
	}
	if l.TopItem() != l.Focused-2 {
		t.Fatalf("topItem should slide so focused is at bottom: top=%d focused=%d", l.TopItem(), l.Focused)
	}
}

func TestListBoxPageDownPagesByViewport(t *testing.T) {
	host := newPalettedHost()
	items := make([]string, 20)
	for i := range items {
		items[i] = string(rune('a' + i))
	}
	l := NewListBox(vio.R(0, 0, 20, 5), items, nil)
	host.Insert(l)
	host.SetCurrent(0)

	l.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyPgDn}})
	if l.Focused != 4 {
		t.Fatalf("PgDn focused=%d", l.Focused)
	}
}

func TestListBoxMouseSelectsRow(t *testing.T) {
	host := newPalettedHost()
	l := NewListBox(vio.R(0, 0, 20, 5), []string{"a", "b", "c"}, nil)
	host.Insert(l)
	host.SetCurrent(0)

	l.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 0, Y: 2},
	})
	if l.Focused != 2 {
		t.Fatalf("mouse focused=%d", l.Focused)
	}
}

func TestListBoxScrollBarSyncsWithFocus(t *testing.T) {
	host := newPalettedHost()
	sb := NewScrollBar(vio.R(20, 0, 1, 5))
	host.Insert(sb)

	l := NewListBox(vio.R(0, 0, 20, 5), []string{"a", "b", "c", "d", "e"}, sb)
	host.Insert(l)
	host.SetCurrent(1)

	l.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	if sb.Value != 1 {
		t.Fatalf("scrollbar should track focus: value=%d", sb.Value)
	}
}

func TestListBoxBroadcastFromScrollBarUpdatesFocus(t *testing.T) {
	host := newPalettedHost()
	sb := NewScrollBar(vio.R(20, 0, 1, 5))
	host.Insert(sb)

	l := NewListBox(vio.R(0, 0, 20, 5), []string{"a", "b", "c", "d", "e"}, sb)
	host.Insert(l)
	host.SetCurrent(1)

	// Fake the broadcast that ScrollBar would post when its value changes.
	sb.Value = 3
	l.HandleEvent(&event.Event{
		What: event.ClassBroadcast,
		Msg:  event.MessageEvent{Command: event.CmdScrollBarChanged, Info: sb},
	})
	if l.Focused != 3 {
		t.Fatalf("scrollbar broadcast should set focus: %d", l.Focused)
	}
}

func TestListBoxDrawHighlightsFocusedRow(t *testing.T) {
	host := newPalettedHost()
	l := NewListBox(vio.R(0, 0, 10, 3), []string{"alpha", "beta", "gamma"}, nil)
	host.Insert(l)
	l.SetState(view.StateFocused, true)
	l.FocusItem(1)

	s := vio.NewSurface(10, 3)
	l.Draw(s)
	rows := s.Snapshot()
	if !strings.HasPrefix(rows[1], "beta") {
		t.Fatalf("row 1 should contain beta, got %q", rows[1])
	}
}

func TestListBoxSetItemsResetsFocus(t *testing.T) {
	host := newPalettedHost()
	l := NewListBox(vio.R(0, 0, 20, 5), []string{"a", "b"}, nil)
	host.Insert(l)
	l.FocusItem(1)

	l.SetItems([]string{"x", "y", "z"})
	if l.Focused != 0 {
		t.Fatalf("SetItems should reset focus, got %d", l.Focused)
	}
	if l.Range != 3 {
		t.Fatalf("SetItems should resize range, got %d", l.Range)
	}
}
