package vio

import (
	"strings"
	"testing"
)

func TestSurfaceFillAndAt(t *testing.T) {
	s := NewSurface(4, 3)
	attr := NewAttr(White, Blue)
	s.FillRect(R(1, 1, 2, 1), 'x', attr)

	if c := s.At(0, 0); c.Rune != 0 {
		t.Fatalf("untouched cell should be zero, got %q", c.Rune)
	}
	if c := s.At(1, 1); c.Rune != 'x' || c.Attr != attr {
		t.Fatalf("filled cell wrong: %+v", c)
	}
	if c := s.At(2, 1); c.Rune != 'x' {
		t.Fatalf("filled cell wrong: %+v", c)
	}
	if c := s.At(3, 1); c.Rune != 0 {
		t.Fatalf("cell past fill should be zero, got %q", c.Rune)
	}
	if c := s.At(-1, 0); c.Rune != 0 {
		t.Fatal("oob At should return zero")
	}
}

func TestSurfaceClear(t *testing.T) {
	s := NewSurface(3, 2)
	attr := NewAttr(LightGray, Blue)
	s.Clear(attr)
	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			if c := s.At(x, y); c.Rune != ' ' || c.Attr != attr {
				t.Fatalf("cell (%d,%d): %+v", x, y, c)
			}
		}
	}
}

func TestSurfaceSetClipsOutOfBounds(t *testing.T) {
	s := NewSurface(2, 2)
	s.Set(-1, 0, 'a', Attr{})
	s.Set(0, -1, 'a', Attr{})
	s.Set(2, 0, 'a', Attr{})
	s.Set(0, 2, 'a', Attr{})
	for _, line := range s.Snapshot() {
		if strings.ContainsRune(line, 'a') {
			t.Fatalf("oob Set leaked: %q", line)
		}
	}
}

func TestSurfaceFillRectClips(t *testing.T) {
	s := NewSurface(4, 4)
	s.FillRect(R(2, 2, 10, 10), '#', NewAttr(White, Black))
	want := []string{
		"    ",
		"    ",
		"  ##",
		"  ##",
	}
	got := s.Snapshot()
	for i, w := range want {
		if got[i] != w {
			t.Errorf("line %d: %q want %q", i, got[i], w)
		}
	}
}

func TestSurfaceDrawString(t *testing.T) {
	s := NewSurface(6, 1)
	n := s.DrawString(1, 0, "hello", Attr{})
	if n != 5 {
		t.Fatalf("draw count: %d", n)
	}
	if got := s.Snapshot()[0]; got != " hello" {
		t.Fatalf("snapshot: %q", got)
	}

	s2 := NewSurface(3, 1)
	s2.DrawString(0, 0, "abcdef", Attr{})
	if got := s2.Snapshot()[0]; got != "abc" {
		t.Fatalf("clipped: %q", got)
	}

	if n := s.DrawString(0, 5, "xx", Attr{}); n != 0 {
		t.Fatalf("oob row should draw nothing: %d", n)
	}
}

func TestSurfaceBox(t *testing.T) {
	s := NewSurface(5, 4)
	s.Box(R(0, 0, 5, 4), Single, Attr{})
	want := []string{
		"┌───┐",
		"│   │",
		"│   │",
		"└───┘",
	}
	got := s.Snapshot()
	for i, w := range want {
		if got[i] != w {
			t.Errorf("line %d: %q want %q", i, got[i], w)
		}
	}
}

func TestSurfaceBoxTooSmall(t *testing.T) {
	s := NewSurface(2, 2)
	s.Box(R(0, 0, 1, 1), Single, Attr{})
	for _, line := range s.Snapshot() {
		for _, r := range line {
			if r != ' ' {
				t.Fatalf("box should not draw: %q", line)
			}
		}
	}
}

func TestSurfaceResize(t *testing.T) {
	s := NewSurface(2, 2)
	s.Set(0, 0, 'a', Attr{})
	s.Resize(2, 2)
	if c := s.At(0, 0); c.Rune != 0 {
		t.Fatalf("same-size resize should clear: %+v", c)
	}
	s.Set(1, 1, 'b', Attr{})
	s.Resize(4, 4)
	w, h := s.Size()
	if w != 4 || h != 4 {
		t.Fatalf("size: %dx%d", w, h)
	}
	if c := s.At(1, 1); c.Rune != 0 {
		t.Fatal("resize should clear")
	}
	s.Resize(-1, -1)
	w, h = s.Size()
	if w != 0 || h != 0 {
		t.Fatalf("negative resize: %dx%d", w, h)
	}
}

func TestSurfaceShadow(t *testing.T) {
	s := NewSurface(8, 6)
	s.FillRect(R(1, 1, 4, 3), 'X', NewAttr(White, Blue))
	s.DrawShadow(R(1, 1, 4, 3))
	for y := 2; y <= 4; y++ {
		for dx := 0; dx < 2; dx++ {
			x := 5 + dx
			c := s.At(x, y)
			if c.Attr.Bg != Black || c.Attr.Fg != DarkGray {
				t.Errorf("shadow (%d,%d): attr %+v", x, y, c.Attr)
			}
		}
	}
	for x := 2; x < 5; x++ {
		c := s.At(x, 4)
		if c.Attr.Bg != Black || c.Attr.Fg != DarkGray {
			t.Errorf("bottom shadow (%d): attr %+v", x, c.Attr)
		}
	}
}

func TestNewSurfaceClampsNegative(t *testing.T) {
	s := NewSurface(-3, -2)
	w, h := s.Size()
	if w != 0 || h != 0 {
		t.Fatalf("size: %dx%d", w, h)
	}
}
