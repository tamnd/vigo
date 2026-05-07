package view

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
)

func TestNewViewDefaults(t *testing.T) {
	v := NewView(vio.R(0, 0, 4, 4))
	if !v.HasState(StateVisible) {
		t.Fatal("NewView should set StateVisible")
	}
	if v.EventMask&event.ClassKey == 0 {
		t.Fatal("default mask should include keys")
	}
	if v.Base() != v {
		t.Fatal("Base should return self")
	}
}

func TestSetHasState(t *testing.T) {
	v := NewView(vio.Rect{})
	v.SetState(StateFocused, true)
	if !v.HasState(StateFocused) {
		t.Fatal("SetState true did not set bit")
	}
	v.SetState(StateFocused, false)
	if v.HasState(StateFocused) {
		t.Fatal("SetState false did not clear bit")
	}
}

func TestChangeBoundsDefault(t *testing.T) {
	v := NewView(vio.R(0, 0, 1, 1))
	v.ChangeBounds(vio.R(2, 3, 4, 5))
	if v.Bounds != (vio.R(2, 3, 4, 5)) {
		t.Fatalf("bounds: %+v", v.Bounds)
	}
}

func TestHandleEventDefaultIsNoOp(t *testing.T) {
	v := NewView(vio.Rect{})
	e := &event.Event{What: event.ClassKey}
	v.HandleEvent(e)
	if e.What != event.ClassKey {
		t.Fatal("default HandleEvent should not consume")
	}
}

func TestDrawUsesPaletteIndex(t *testing.T) {
	g := NewGroup(vio.R(0, 0, 4, 1))
	g.Palette = vio.Palette{
		{Fg: vio.White, Bg: vio.Blue},
		{Fg: vio.Yellow, Bg: vio.Red},
	}
	v := NewView(vio.R(0, 0, 4, 1))
	v.Owner = g
	v.PaletteIndex = 2
	s := vio.NewSurface(4, 1)
	v.Draw(s)
	cell := s.At(0, 0)
	if cell.Attr.Fg != vio.Yellow || cell.Attr.Bg != vio.Red {
		t.Fatalf("draw used wrong palette slot: %+v", cell.Attr)
	}
}

func TestMapColorWalksOwnerChain(t *testing.T) {
	root := NewGroup(vio.R(0, 0, 10, 10))
	root.Palette = vio.Palette{{Fg: vio.White, Bg: vio.Blue}}
	mid := NewGroup(vio.R(0, 0, 10, 10))
	root.Insert(mid)
	leaf := NewView(vio.R(0, 0, 1, 1))
	mid.Insert(leaf)
	if got := leaf.MapColor(1); got != (vio.Attr{Fg: vio.White, Bg: vio.Blue}) {
		t.Fatalf("walk failed: %+v", got)
	}
}

func TestMapColorOrphanFallback(t *testing.T) {
	v := NewView(vio.Rect{})
	got := v.MapColor(1)
	want := vio.BorlandClassic.Map(1)
	if got != want {
		t.Fatalf("orphan fallback: %+v", got)
	}
}
