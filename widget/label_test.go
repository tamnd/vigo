package widget

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func TestLabelParsesMnemonic(t *testing.T) {
	l := NewLabel(vio.R(0, 0, 10, 1), "~N~ame:", nil)
	if l.Mnemonic() != 'n' {
		t.Fatalf("mnemonic: %q", l.Mnemonic())
	}
	if l.Title != "Name:" {
		t.Fatalf("title not stripped: %q", l.Title)
	}
}

func TestLabelDrawTextRendersTitle(t *testing.T) {
	host := newPalettedHost()
	l := NewLabel(vio.R(0, 0, 10, 1), "~N~ame:", nil)
	host.Insert(l)
	s := vio.NewSurface(10, 1)
	l.Draw(s)
	row := s.Snapshot()[0]
	if row[:5] != "Name:" {
		t.Fatalf("expected Name: prefix, got %q", row)
	}
}

func TestLabelAltMnemonicMovesFocusToPeer(t *testing.T) {
	host := newPalettedHost()
	peer := NewInputLine(vio.R(0, 1, 10, 1), 0)
	l := NewLabel(vio.R(0, 0, 10, 1), "~N~ame:", peer)
	host.Insert(l)
	host.Insert(peer)
	host.SetCurrent(0) // label is not selectable but force focus elsewhere

	e := &event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'n', Mod: event.ModAlt},
	}
	l.HandleEvent(e)

	if e.What != event.ClassNothing {
		t.Fatalf("Alt-mnemonic should be consumed, got class %v", e.What)
	}
	if host.Current() != peer {
		t.Fatal("peer should now be focused")
	}
}

func TestLabelWithoutPeerIsInert(t *testing.T) {
	host := newPalettedHost()
	l := NewLabel(vio.R(0, 0, 10, 1), "~N~ame:", nil)
	host.Insert(l)

	e := &event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'n', Mod: event.ModAlt},
	}
	l.HandleEvent(e)

	if e.What != event.ClassKey {
		t.Fatal("label without peer should ignore Alt-mnemonic")
	}
}

func TestLabelIgnoresNonAltKeys(t *testing.T) {
	host := newPalettedHost()
	peer := NewInputLine(vio.R(0, 1, 10, 1), 0)
	l := NewLabel(vio.R(0, 0, 10, 1), "~N~ame:", peer)
	host.Insert(l)
	host.Insert(peer)

	e := &event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'n'},
	}
	l.HandleEvent(e)
	if e.What != event.ClassKey {
		t.Fatal("plain rune should not trigger label dispatch")
	}
	_ = view.StateFocused
}
