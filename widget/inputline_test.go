package widget

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func newFocusedInput(host *view.Group, w int) *InputLine {
	in := NewInputLine(vio.R(0, 0, w, 1), 0)
	host.Insert(in)
	host.SetCurrent(0)
	return in
}

func TestInputLineSetStringClampsToMaxLen(t *testing.T) {
	host := newPalettedHost()
	in := NewInputLine(vio.R(0, 0, 10, 1), 4)
	host.Insert(in)
	in.SetString("abcdef")
	if in.String() != "abcd" {
		t.Fatalf("expected truncation, got %q", in.String())
	}
	if in.CursorIndex() != 4 {
		t.Fatalf("cursor should sit at end of clamped buffer, got %d", in.CursorIndex())
	}
}

func TestInputLineInsertRune(t *testing.T) {
	host := newPalettedHost()
	in := newFocusedInput(host, 10)

	for _, r := range "hi" {
		in.HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyRune, Rune: r},
		})
	}
	if in.String() != "hi" {
		t.Fatalf("expected 'hi', got %q", in.String())
	}
	if in.CursorIndex() != 2 {
		t.Fatalf("cursor: %d", in.CursorIndex())
	}
}

func TestInputLineBackspaceDeletesPrev(t *testing.T) {
	host := newPalettedHost()
	in := newFocusedInput(host, 10)
	in.SetString("abc")

	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyBackspace},
	})
	if in.String() != "ab" {
		t.Fatalf("expected 'ab', got %q", in.String())
	}
}

func TestInputLineDeleteRemovesNext(t *testing.T) {
	host := newPalettedHost()
	in := newFocusedInput(host, 10)
	in.SetString("abc")
	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyHome},
	})

	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyDelete},
	})
	if in.String() != "bc" {
		t.Fatalf("expected 'bc', got %q", in.String())
	}
}

func TestInputLineArrowAndHomeEnd(t *testing.T) {
	host := newPalettedHost()
	in := newFocusedInput(host, 10)
	in.SetString("abc")

	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyHome},
	})
	if in.CursorIndex() != 0 {
		t.Fatalf("home: %d", in.CursorIndex())
	}

	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyArrowRight},
	})
	if in.CursorIndex() != 1 {
		t.Fatalf("right: %d", in.CursorIndex())
	}

	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyEnd},
	})
	if in.CursorIndex() != 3 {
		t.Fatalf("end: %d", in.CursorIndex())
	}
}

func TestInputLineShiftSelectsAndDeleteRemovesRange(t *testing.T) {
	host := newPalettedHost()
	in := newFocusedInput(host, 10)
	in.SetString("abc")
	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyHome},
	})

	for range 2 {
		in.HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyArrowRight, Mod: event.ModShift},
		})
	}
	a, b := in.Selection()
	if a != 0 || b != 2 {
		t.Fatalf("selection (%d,%d) want (0,2)", a, b)
	}

	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyDelete},
	})
	if in.String() != "c" {
		t.Fatalf("expected 'c' after deleting selection, got %q", in.String())
	}
}

func TestInputLineCtrlYClears(t *testing.T) {
	host := newPalettedHost()
	in := newFocusedInput(host, 10)
	in.SetString("abc")

	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'y', Mod: event.ModCtrl},
	})
	if in.String() != "" {
		t.Fatalf("Ctrl-Y should clear, got %q", in.String())
	}
}

func TestInputLineUnfocusedIgnoresKeys(t *testing.T) {
	host := newPalettedHost()
	in := NewInputLine(vio.R(0, 0, 10, 1), 0)
	host.Insert(in)
	host.SetCurrent(-1)

	in.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'a'},
	})
	if in.String() != "" {
		t.Fatalf("unfocused input should not insert: %q", in.String())
	}
}

func TestInputLineMaxLenBlocksInsert(t *testing.T) {
	host := newPalettedHost()
	in := NewInputLine(vio.R(0, 0, 10, 1), 2)
	host.Insert(in)
	host.SetCurrent(0)

	for _, r := range "abc" {
		in.HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyRune, Rune: r},
		})
	}
	if in.String() != "ab" {
		t.Fatalf("MaxLen should cap insertion, got %q", in.String())
	}
}
