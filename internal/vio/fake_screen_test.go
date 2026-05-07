package vio

import (
	"testing"

	"github.com/tamnd/vigo/internal/event"
)

func TestFakeScreenRoundtrip(t *testing.T) {
	f := NewFakeScreen(4, 2)
	if err := f.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	defer f.Fini()

	w, h := f.Size()
	if w != 4 || h != 2 {
		t.Fatalf("size: %dx%d", w, h)
	}

	s := NewSurface(4, 2)
	s.Set(0, 0, 'a', NewAttr(White, Blue))
	if err := f.Show(s); err != nil {
		t.Fatalf("show: %v", err)
	}
	if c := f.Surface().At(0, 0); c.Rune != 'a' {
		t.Fatalf("show did not record: %+v", c)
	}

	f.SetCursor(2, 1, true)
	if x, y, vis := f.Cursor(); x != 2 || y != 1 || !vis {
		t.Fatalf("cursor: %d %d %v", x, y, vis)
	}
}

func TestFakeScreenPushEvent(t *testing.T) {
	f := NewFakeScreen(8, 8)
	defer f.Fini()
	f.PushEvent(event.Event{What: event.ClassKey})
	got := <-f.Events()
	if got.What != event.ClassKey {
		t.Fatalf("class: %v", got.What)
	}
}

func TestFakeScreenResizeEmitsEvent(t *testing.T) {
	f := NewFakeScreen(8, 8)
	defer f.Fini()
	f.Resize(20, 10)
	w, h := f.Size()
	if w != 20 || h != 10 {
		t.Fatalf("size: %dx%d", w, h)
	}
	got := <-f.Events()
	if got.What != event.ClassResize {
		t.Fatalf("class: %v", got.What)
	}
	if got.Mouse.X != 20 || got.Mouse.Y != 10 {
		t.Fatalf("payload: %+v", got.Mouse)
	}
}

func TestFakeScreenFiniIsIdempotent(_ *testing.T) {
	f := NewFakeScreen(2, 2)
	f.Fini()
	f.Fini() // must not panic or close twice
	f.PushEvent(event.Event{What: event.ClassKey})
}
