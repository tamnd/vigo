package vio

import (
	"testing"

	"github.com/gdamore/tcell/v2"

	"github.com/tamnd/vigo/event"
)

func TestTranslateKeyRune(t *testing.T) {
	ev := tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModShift|tcell.ModAlt)
	out, ok := translate(ev)
	if !ok {
		t.Fatal("translate refused key event")
	}
	if out.What != event.ClassKey {
		t.Fatalf("class: %v", out.What)
	}
	if out.Key.Key != event.KeyRune || out.Key.Rune != 'a' {
		t.Fatalf("payload: %+v", out.Key)
	}
	if out.Key.Mod&event.ModShift == 0 || out.Key.Mod&event.ModAlt == 0 {
		t.Fatalf("modifiers: %v", out.Key.Mod)
	}
}

func TestTranslateNamedKeys(t *testing.T) {
	cases := []struct {
		in   tcell.Key
		want event.Key
	}{
		{tcell.KeyEnter, event.KeyEnter},
		{tcell.KeyEscape, event.KeyEsc},
		{tcell.KeyTab, event.KeyTab},
		{tcell.KeyBacktab, event.KeyShiftTab},
		{tcell.KeyBackspace, event.KeyBackspace},
		{tcell.KeyBackspace2, event.KeyBackspace},
		{tcell.KeyDelete, event.KeyDelete},
		{tcell.KeyInsert, event.KeyInsert},
		{tcell.KeyHome, event.KeyHome},
		{tcell.KeyEnd, event.KeyEnd},
		{tcell.KeyPgUp, event.KeyPgUp},
		{tcell.KeyPgDn, event.KeyPgDn},
		{tcell.KeyUp, event.KeyArrowUp},
		{tcell.KeyDown, event.KeyArrowDown},
		{tcell.KeyLeft, event.KeyArrowLeft},
		{tcell.KeyRight, event.KeyArrowRight},
		{tcell.KeyF1, event.KeyF1},
		{tcell.KeyF2, event.KeyF2},
		{tcell.KeyF3, event.KeyF3},
		{tcell.KeyF4, event.KeyF4},
		{tcell.KeyF5, event.KeyF5},
		{tcell.KeyF6, event.KeyF6},
		{tcell.KeyF7, event.KeyF7},
		{tcell.KeyF8, event.KeyF8},
		{tcell.KeyF9, event.KeyF9},
		{tcell.KeyF10, event.KeyF10},
		{tcell.KeyF11, event.KeyF11},
		{tcell.KeyF12, event.KeyF12},
	}
	for _, c := range cases {
		ev := tcell.NewEventKey(c.in, 0, tcell.ModNone)
		out, ok := translate(ev)
		if !ok || out.Key.Key != c.want {
			t.Errorf("%v: got %v ok=%v", c.in, out.Key.Key, ok)
		}
	}
}

func TestTranslateCtrlLetter(t *testing.T) {
	ev := tcell.NewEventKey(tcell.KeyCtrlA, 0, tcell.ModCtrl)
	out, ok := translate(ev)
	if !ok {
		t.Fatal("translate refused")
	}
	if out.Key.Key != event.KeyRune || out.Key.Rune != 'a' {
		t.Fatalf("ctrl-a: %+v", out.Key)
	}
}

func TestTranslateUnknownKey(t *testing.T) {
	ev := tcell.NewEventKey(tcell.KeyHelp, 0, tcell.ModNone)
	out, _ := translate(ev)
	if out.Key.Key != event.KeyNone {
		t.Fatalf("unknown should map to KeyNone: %v", out.Key.Key)
	}
}

func TestTranslateMouseButtons(t *testing.T) {
	cases := []struct {
		btn      tcell.ButtonMask
		wantBits event.MouseButton
		wantWhat event.Class
	}{
		{tcell.Button1, event.MouseLeft, event.ClassMouseDown},
		{tcell.Button2, event.MouseMiddle, event.ClassMouseDown},
		{tcell.Button3, event.MouseRight, event.ClassMouseDown},
		{tcell.WheelUp, event.MouseWheelUp, event.ClassMouseWheel},
		{tcell.WheelDown, event.MouseWheelDown, event.ClassMouseWheel},
		{tcell.ButtonNone, 0, event.ClassMouseUp},
	}
	for _, c := range cases {
		ev := tcell.NewEventMouse(3, 4, c.btn, tcell.ModShift|tcell.ModCtrl|tcell.ModAlt)
		out, ok := translate(ev)
		if !ok {
			t.Fatalf("%v refused", c.btn)
		}
		if out.What != c.wantWhat {
			t.Errorf("%v: class %v want %v", c.btn, out.What, c.wantWhat)
		}
		if out.Mouse.Buttons&c.wantBits != c.wantBits {
			t.Errorf("%v: bits %v want %v", c.btn, out.Mouse.Buttons, c.wantBits)
		}
		if out.Mouse.X != 3 || out.Mouse.Y != 4 {
			t.Errorf("%v: pos %v", c.btn, out.Mouse)
		}
		wantMod := event.ModShift | event.ModCtrl | event.ModAlt
		if out.Mouse.Mod&wantMod != wantMod {
			t.Errorf("%v: mod %v", c.btn, out.Mouse.Mod)
		}
	}
}

func TestTranslateResize(t *testing.T) {
	ev := tcell.NewEventResize(80, 24)
	out, ok := translate(ev)
	if !ok {
		t.Fatal("resize refused")
	}
	if out.What != event.ClassResize {
		t.Fatalf("class: %v", out.What)
	}
	if out.Mouse.X != 80 || out.Mouse.Y != 24 {
		t.Fatalf("size: %+v", out.Mouse)
	}
}

func TestUninitTcellScreen(t *testing.T) {
	s := NewTcellScreen()
	w, h := s.Size()
	if w != 0 || h != 0 {
		t.Fatalf("size before init: %dx%d", w, h)
	}
	if err := s.Show(NewSurface(1, 1)); err == nil {
		t.Fatal("Show should fail before Init")
	}
	s.Fini() // safe to call without Init
}
