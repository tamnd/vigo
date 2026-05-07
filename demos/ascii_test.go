package demos

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
)

func TestAsciiTableCursorStartsAtZero(t *testing.T) {
	a := NewAsciiTable(AsciiDefaultBounds)
	hostFor(a)
	if a.Cursor() != 0 {
		t.Fatalf("cursor=%d want 0", a.Cursor())
	}
}

func TestAsciiTableArrowsMoveCursor(t *testing.T) {
	a := NewAsciiTable(AsciiDefaultBounds)
	hostFor(a)

	a.HandleEvent(keyOf(event.KeyArrowRight))
	if a.Cursor() != 1 {
		t.Fatalf("right cursor=%d want 1", a.Cursor())
	}
	a.HandleEvent(keyOf(event.KeyArrowDown))
	if a.Cursor() != 17 {
		t.Fatalf("down cursor=%d want 17", a.Cursor())
	}
}

func TestAsciiTableLeftWrapsToEnd(t *testing.T) {
	a := NewAsciiTable(AsciiDefaultBounds)
	hostFor(a)

	a.HandleEvent(keyOf(event.KeyArrowLeft))
	if a.Cursor() != 255 {
		t.Fatalf("wrap left cursor=%d want 255", a.Cursor())
	}
}

func TestAsciiTableEndJumpsToLast(t *testing.T) {
	a := NewAsciiTable(AsciiDefaultBounds)
	hostFor(a)

	a.HandleEvent(keyOf(event.KeyEnd))
	if a.Cursor() != 255 {
		t.Fatalf("End cursor=%d want 255", a.Cursor())
	}
	a.HandleEvent(keyOf(event.KeyHome))
	if a.Cursor() != 0 {
		t.Fatalf("Home cursor=%d want 0", a.Cursor())
	}
}

func TestAsciiTableStatusFormatsCodepoint(t *testing.T) {
	a := NewAsciiTable(AsciiDefaultBounds)
	hostFor(a)

	// Move cursor to 'A' (65 = 0x41).
	for range 65 {
		a.HandleEvent(keyOf(event.KeyArrowRight))
	}
	got := a.status.Text
	if !strings.Contains(got, "dec  65") {
		t.Fatalf("status %q missing dec 65", got)
	}
	if !strings.Contains(got, "hex 41") {
		t.Fatalf("status %q missing hex 41", got)
	}
	if !strings.Contains(got, "char A") {
		t.Fatalf("status %q missing char A", got)
	}
}

func TestAsciiTableNonPrintableShowsDot(t *testing.T) {
	a := NewAsciiTable(AsciiDefaultBounds)
	hostFor(a)
	// Cursor starts at 0 (NUL); non-printable -> '.'.
	if !strings.Contains(a.status.Text, "char .") {
		t.Fatalf("status %q want '.' for NUL", a.status.Text)
	}
}
