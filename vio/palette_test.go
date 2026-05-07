package vio

import "testing"

func TestPaletteMap(t *testing.T) {
	p := Palette{
		{Fg: White, Bg: Blue},
		{Fg: Yellow, Bg: Red},
	}
	if got := p.Map(1); got != (Attr{Fg: White, Bg: Blue}) {
		t.Fatalf("slot 1: %+v", got)
	}
	if got := p.Map(2); got != (Attr{Fg: Yellow, Bg: Red}) {
		t.Fatalf("slot 2: %+v", got)
	}
	if got := p.Map(0); got != p[0] {
		t.Fatalf("slot 0 should fall back to first: %+v", got)
	}
	if got := p.Map(99); got != p[0] {
		t.Fatalf("oob slot: %+v", got)
	}
}

func TestEmptyPaletteMap(t *testing.T) {
	var p Palette
	if got := p.Map(1); got != (Attr{}) {
		t.Fatalf("empty: %+v", got)
	}
}

func TestBorlandClassicHasIDEColors(t *testing.T) {
	desktop := BorlandClassic.Map(1)
	if desktop.Bg != Blue {
		t.Fatalf("desktop bg: %v", desktop.Bg)
	}
	frame := BorlandClassic.Map(11)
	if frame.Fg != Yellow || frame.Bg != Blue {
		t.Fatalf("active frame: %+v", frame)
	}
}
