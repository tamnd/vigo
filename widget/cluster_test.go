package widget

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func TestCheckBoxesPressTogglesBit(t *testing.T) {
	host := newPalettedHost()
	c := NewCheckBoxes(vio.R(0, 0, 20, 3), []string{"~A~", "~B~", "~C~"})
	host.Insert(c)
	host.SetCurrent(0)

	c.Press(0)
	c.Press(2)
	if c.Value != (1<<0)|(1<<2) {
		t.Fatalf("expected bits 0,2 set, got %b", c.Value)
	}
	c.Press(0)
	if c.Value != 1<<2 {
		t.Fatalf("expected bit 2 only, got %b", c.Value)
	}
}

func TestCheckBoxesSpaceTogglesFocused(t *testing.T) {
	host := newPalettedHost()
	c := NewCheckBoxes(vio.R(0, 0, 20, 2), []string{"a", "b"})
	host.Insert(c)
	host.SetCurrent(0)
	c.SetSel(1)

	c.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: ' '},
	})
	if c.Value != 1<<1 {
		t.Fatalf("Space should toggle focused item, value=%b", c.Value)
	}
}

func TestCheckBoxesAltMnemonicTogglesAndFocuses(t *testing.T) {
	host := newPalettedHost()
	c := NewCheckBoxes(vio.R(0, 0, 20, 3), []string{"~A~lpha", "~B~eta", "~G~amma"})
	host.Insert(c)
	host.SetCurrent(0)

	c.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'b', Mod: event.ModAlt},
	})
	if c.Value != 1<<1 {
		t.Fatalf("Alt-B should toggle item 1, value=%b", c.Value)
	}
	if c.Sel() != 1 {
		t.Fatalf("Alt-B should focus item 1, sel=%d", c.Sel())
	}
}

func TestCheckBoxesArrowsMoveSel(t *testing.T) {
	host := newPalettedHost()
	c := NewCheckBoxes(vio.R(0, 0, 20, 3), []string{"a", "b", "c"})
	host.Insert(c)
	host.SetCurrent(0)

	c.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	if c.Sel() != 1 {
		t.Fatalf("ArrowDown sel=%d", c.Sel())
	}
	c.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnd}})
	if c.Sel() != 2 {
		t.Fatalf("End sel=%d", c.Sel())
	}
	c.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyHome}})
	if c.Sel() != 0 {
		t.Fatalf("Home sel=%d", c.Sel())
	}
}

func TestCheckBoxesMouseClickTogglesAndFocuses(t *testing.T) {
	host := newPalettedHost()
	c := NewCheckBoxes(vio.R(0, 0, 20, 3), []string{"a", "b", "c"})
	host.Insert(c)

	c.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 1, Y: 1},
	})
	if c.Sel() != 1 || c.Value != 1<<1 {
		t.Fatalf("click on row 1: sel=%d value=%b", c.Sel(), c.Value)
	}
}

func TestCheckBoxesDrawShowsBrackets(t *testing.T) {
	host := newPalettedHost()
	c := NewCheckBoxes(vio.R(0, 0, 12, 2), []string{"a", "b"})
	host.Insert(c)
	c.Press(0)

	s := vio.NewSurface(12, 2)
	c.Draw(s)
	row0 := s.Snapshot()[0]
	if !strings.HasPrefix(row0, "[X] a") {
		t.Fatalf("checked row should start with [X] a, got %q", row0)
	}
	row1 := s.Snapshot()[1]
	if !strings.HasPrefix(row1, "[ ] b") {
		t.Fatalf("unchecked row should start with [ ] b, got %q", row1)
	}
}

func TestRadioButtonsExclusiveSelection(t *testing.T) {
	host := newPalettedHost()
	r := NewRadioButtons(vio.R(0, 0, 20, 3), []string{"a", "b", "c"})
	host.Insert(r)
	host.SetCurrent(0)

	r.Press(2)
	if r.Value != 2 || r.Mark(0) || r.Mark(1) || !r.Mark(2) {
		t.Fatalf("Press(2) should select only item 2: value=%d", r.Value)
	}
	r.Press(0)
	if r.Value != 0 || !r.Mark(0) || r.Mark(2) {
		t.Fatalf("Press(0) should switch selection to 0: value=%d", r.Value)
	}
}

func TestRadioButtonsAltMnemonicSelects(t *testing.T) {
	host := newPalettedHost()
	r := NewRadioButtons(vio.R(0, 0, 20, 3), []string{"~A~", "~B~", "~C~"})
	host.Insert(r)
	host.SetCurrent(0)

	r.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'c', Mod: event.ModAlt},
	})
	if r.Value != 2 {
		t.Fatalf("Alt-C should select item 2, value=%d", r.Value)
	}
}

func TestRadioButtonsDrawShowsParens(t *testing.T) {
	host := newPalettedHost()
	r := NewRadioButtons(vio.R(0, 0, 10, 2), []string{"a", "b"})
	host.Insert(r)
	r.Press(0)

	s := vio.NewSurface(10, 2)
	r.Draw(s)
	row0 := s.Snapshot()[0]
	if !strings.HasPrefix(row0, "(•) a") {
		t.Fatalf("selected row: %q", row0)
	}
	row1 := s.Snapshot()[1]
	if !strings.HasPrefix(row1, "( ) b") {
		t.Fatalf("unselected row: %q", row1)
	}
}

func TestRadioButtonsMouseSelects(t *testing.T) {
	host := newPalettedHost()
	r := NewRadioButtons(vio.R(0, 0, 20, 3), []string{"a", "b", "c"})
	host.Insert(r)
	host.SetCurrent(0)

	r.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 0, Y: 2},
	})
	if r.Value != 2 || r.Sel() != 2 {
		t.Fatalf("mouse-down row 2: sel=%d value=%d", r.Sel(), r.Value)
	}
}

func TestClusterUnfocusedSpaceIgnored(t *testing.T) {
	host := newPalettedHost()
	c := NewCheckBoxes(vio.R(0, 0, 20, 1), []string{"a"})
	host.Insert(c)
	host.SetCurrent(-1)
	c.SetState(view.StateFocused, false)

	c.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: ' '},
	})
	if c.Value != 0 {
		t.Fatal("unfocused cluster should ignore Space")
	}
}
