package menu

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/cmd"
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
)

func sampleItems() []Item {
	return []Item{
		{Title: "New", Hotkey: 0, Cmd: event.CmdUser + 1, Shortcut: "F3"},
		{Title: "Open", Hotkey: 0, Cmd: event.CmdUser + 2, Shortcut: "F4"},
		{Sep: true},
		{Title: "Save", Hotkey: 0, Cmd: event.CmdUser + 3, Shortcut: "F2"},
		{Title: "Exit", Hotkey: 1, Cmd: event.CmdQuit, Shortcut: "Alt-X"},
	}
}

func TestMenuBoxArrowsSkipSeparators(t *testing.T) {
	host := newPalettedHost()
	mb := NewMenuBox(0, 0, sampleItems(), nil)
	host.Insert(mb)

	if mb.SelectedIndex() != 0 {
		t.Fatalf("initial sel=%d", mb.SelectedIndex())
	}
	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	if mb.SelectedIndex() != 3 {
		t.Fatalf("ArrowDown should skip separator: %d", mb.SelectedIndex())
	}
	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowUp}})
	if mb.SelectedIndex() != 1 {
		t.Fatalf("ArrowUp should skip separator: %d", mb.SelectedIndex())
	}
}

func TestMenuBoxEnterFiresChosen(t *testing.T) {
	host := newPalettedHost()
	mb := NewMenuBox(0, 0, sampleItems(), nil)
	host.Insert(mb)

	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnter}})
	if mb.Result() != event.CmdMenu {
		t.Fatalf("Result=%d", mb.Result())
	}
	if mb.Chosen() != event.CmdUser+2 {
		t.Fatalf("Chosen=%d", mb.Chosen())
	}
}

func TestMenuBoxEscCancels(t *testing.T) {
	host := newPalettedHost()
	mb := NewMenuBox(0, 0, sampleItems(), nil)
	host.Insert(mb)

	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEsc}})
	if mb.Result() != event.CmdCancel {
		t.Fatalf("Esc result=%d", mb.Result())
	}
	if mb.Chosen() != event.CmdNone {
		t.Fatalf("cancel should not fire chosen, got %d", mb.Chosen())
	}
}

func TestMenuBoxHotkeyDispatches(t *testing.T) {
	host := newPalettedHost()
	mb := NewMenuBox(0, 0, sampleItems(), nil)
	host.Insert(mb)

	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyRune, Rune: 'x'}})
	if mb.Result() != event.CmdMenu {
		t.Fatalf("hotkey x result=%d", mb.Result())
	}
	if mb.Chosen() != event.CmdQuit {
		t.Fatalf("hotkey x chosen=%d", mb.Chosen())
	}
}

func TestMenuBoxDisabledRowDimAndIgnored(t *testing.T) {
	host := newPalettedHost()
	en := cmd.NewEnabler(cmd.NewSet(event.CmdUser+1, event.CmdUser+3, event.CmdQuit))
	// CmdUser+2 (Open) is NOT enabled
	mb := NewMenuBox(0, 0, sampleItems(), en)
	host.Insert(mb)

	if mb.IsEnabled(1) {
		t.Fatalf("row 1 (Open) should be disabled")
	}
	if !mb.IsEnabled(0) {
		t.Fatalf("row 0 (New) should be enabled")
	}
	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnter}})
	if mb.Result() != event.CmdNone {
		t.Fatalf("disabled row should not fire: %d", mb.Result())
	}
}

func TestMenuBoxEnableFlipsDimState(t *testing.T) {
	host := newPalettedHost()
	en := cmd.NewEnabler(cmd.NewSet(event.CmdUser + 1))
	mb := NewMenuBox(0, 0, sampleItems(), en)
	host.Insert(mb)

	if mb.IsEnabled(1) {
		t.Fatalf("row 1 starts disabled")
	}
	en.Enable(cmd.NewSet(event.CmdUser + 2))
	if !mb.IsEnabled(1) {
		t.Fatalf("Enable should make row 1 active")
	}
	en.Disable(cmd.NewSet(event.CmdUser + 2))
	if mb.IsEnabled(1) {
		t.Fatalf("Disable should make row 1 dim again")
	}
}

func TestMenuBoxSubmenuFlagsChildOpen(t *testing.T) {
	host := newPalettedHost()
	items := []Item{
		{Title: "More", Hotkey: 0, Children: []Item{
			{Title: "Sub1", Cmd: event.CmdUser + 10},
		}},
	}
	mb := NewMenuBox(0, 0, items, nil)
	host.Insert(mb)

	mb.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnter}})
	if !mb.ChildOpen() {
		t.Fatalf("submenu Enter should set ChildOpen")
	}
	mb.AckChild()
	if mb.ChildOpen() {
		t.Fatalf("AckChild should clear")
	}
}

func TestMenuBoxMouseClickFiresRow(t *testing.T) {
	host := newPalettedHost()
	mb := NewMenuBox(0, 0, sampleItems(), nil)
	host.Insert(mb)

	// row 0 sits at Y=1 (after top border)
	mb.HandleEvent(&event.Event{What: event.ClassMouseDown, Mouse: event.MouseEvent{X: 2, Y: 1}})
	if mb.Result() != event.CmdMenu {
		t.Fatalf("mouse row 0 result=%d", mb.Result())
	}
	if mb.Chosen() != event.CmdUser+1 {
		t.Fatalf("mouse row 0 chosen=%d", mb.Chosen())
	}
}

func TestMenuBoxClickOutsideCancels(t *testing.T) {
	host := newPalettedHost()
	mb := NewMenuBox(0, 0, sampleItems(), nil)
	host.Insert(mb)

	mb.HandleEvent(&event.Event{What: event.ClassMouseDown, Mouse: event.MouseEvent{X: 50, Y: 20}})
	if mb.Result() != event.CmdCancel {
		t.Fatalf("click outside result=%d", mb.Result())
	}
}

func TestMenuBoxDrawsTitleAndShortcut(t *testing.T) {
	host := newPalettedHost()
	mb := NewMenuBox(0, 0, sampleItems(), nil)
	host.Insert(mb)

	w := mb.Bounds.W
	h := mb.Bounds.H
	s := vio.NewSurface(w, h)
	mb.Draw(s)
	rows := s.Snapshot()
	// row 1 = first item ("New", shortcut "F3")
	row1 := rows[1]
	if !strings.Contains(row1, "New") {
		t.Fatalf("row1 missing title: %q", row1)
	}
	if !strings.Contains(row1, "F3") {
		t.Fatalf("row1 missing shortcut: %q", row1)
	}
}
