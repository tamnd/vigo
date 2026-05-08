package demos

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
)

func TestASCIITableCursorStartsAtZero(t *testing.T) {
	a := NewASCIITable(ASCIIDefaultBounds)
	hostFor(a)
	if a.Cursor() != 0 {
		t.Fatalf("cursor=%d want 0", a.Cursor())
	}
}

func TestASCIITableArrowsMoveCursor(t *testing.T) {
	a := NewASCIITable(ASCIIDefaultBounds)
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

func TestASCIITableLeftWrapsToEnd(t *testing.T) {
	a := NewASCIITable(ASCIIDefaultBounds)
	hostFor(a)

	a.HandleEvent(keyOf(event.KeyArrowLeft))
	if a.Cursor() != 255 {
		t.Fatalf("wrap left cursor=%d want 255", a.Cursor())
	}
}

func TestASCIITableEndJumpsToLast(t *testing.T) {
	a := NewASCIITable(ASCIIDefaultBounds)
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

func TestASCIITableStatusFormatsCodepoint(t *testing.T) {
	a := NewASCIITable(ASCIIDefaultBounds)
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

func TestASCIITableNonPrintableShowsDot(t *testing.T) {
	a := NewASCIITable(ASCIIDefaultBounds)
	hostFor(a)
	// Cursor starts at 0 (NUL); non-printable -> '.'.
	if !strings.Contains(a.status.Text, "char .") {
		t.Fatalf("status %q want '.' for NUL", a.status.Text)
	}
}

func TestASCIITableClickSetsCursor(t *testing.T) {
	a := NewASCIITable(ASCIIDefaultBounds)
	hostFor(a)

	// Click cell at (col=5, row=3) -> codepoint 3*16+5 = 53.
	c := a.grid.Bounds.X + 5
	r := a.grid.Bounds.Y + 3
	a.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: c, Y: r, Buttons: event.MouseLeft},
	})
	if got := a.Cursor(); got != 53 {
		t.Fatalf("click cursor=%d want 53", got)
	}
}
