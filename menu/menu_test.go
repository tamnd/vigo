package menu

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func newPalettedHost() *view.Group {
	g := view.NewGroup(vio.R(0, 0, 80, 24))
	g.Palette = vio.BorlandClassic
	return g
}

func TestBarDrawsItems(t *testing.T) {
	host := newPalettedHost()
	b := NewBar(40, []Item{{Title: "File", Hotkey: 0}, {Title: "Edit", Hotkey: 0}})
	host.Insert(b)
	s := vio.NewSurface(40, 1)
	b.Draw(s)
	row := s.Snapshot()[0]
	if !strings.Contains(row, "File") || !strings.Contains(row, "Edit") {
		t.Fatalf("bar row: %q", row)
	}
}

func TestBarTruncatesAtRightEdge(t *testing.T) {
	host := newPalettedHost()
	b := NewBar(6, []Item{{Title: "FileMenu", Hotkey: 0}})
	host.Insert(b)
	s := vio.NewSurface(6, 1)
	b.Draw(s)
	row := s.Snapshot()[0]
	if len(row) != 6 {
		t.Fatalf("row width: %d", len(row))
	}
}

func TestBarHandleEventNoOp(t *testing.T) {
	b := NewBar(20, DefaultItems)
	e := &event.Event{What: event.ClassKey}
	b.HandleEvent(e)
	if e.What != event.ClassKey {
		t.Fatal("bar should not consume keys in v0.1")
	}
}

func TestLineDrawsHints(t *testing.T) {
	host := newPalettedHost()
	l := NewLine(40, 23, DefaultHints)
	host.Insert(l)
	s := vio.NewSurface(40, 24)
	l.Draw(s)
	row := s.Snapshot()[23]
	for _, want := range []string{"F1", "Help", "F10", "Menu", "Alt-X", "Exit"} {
		if !strings.Contains(row, want) {
			t.Errorf("status row missing %q: %q", want, row)
		}
	}
}

func TestLineHandleEventNoOp(t *testing.T) {
	l := NewLine(40, 0, DefaultHints)
	e := &event.Event{What: event.ClassKey}
	l.HandleEvent(e)
	if e.What != event.ClassKey {
		t.Fatal("line should not consume in v0.1")
	}
}

func TestLineTruncatesAtRightEdge(t *testing.T) {
	host := newPalettedHost()
	l := NewLine(8, 0, []Hint{{Key: "Alt-X", Action: "ExitNow", Cmd: event.CmdQuit}})
	host.Insert(l)
	s := vio.NewSurface(8, 1)
	l.Draw(s)
	row := s.Snapshot()[0]
	if len(row) != 8 {
		t.Fatalf("row width: %d", len(row))
	}
}
