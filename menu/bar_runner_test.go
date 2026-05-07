package menu

import (
	"testing"

	"github.com/tamnd/vigo/event"
)

func barWithSubmenus() []Item {
	file := []Item{
		{Title: "New", Cmd: event.CmdUser + 1, Hotkey: 0},
		{Title: "Open", Cmd: event.CmdUser + 2, Hotkey: 0},
	}
	edit := []Item{
		{Title: "Cut", Cmd: event.CmdUser + 10, Hotkey: 0},
	}
	return []Item{
		{Title: "≡", Hotkey: -1},
		{Title: "File", Hotkey: 0, Children: file},
		{Title: "Edit", Hotkey: 0, Children: edit},
	}
}

// fakeRunner records every box it runs and returns deterministically.
type fakeRunner struct {
	calls   int
	resolve func(*MenuBox)
}

func (f *fakeRunner) Run(box *MenuBox) event.CommandID {
	f.calls++
	f.resolve(box)
	return box.Result()
}

func TestBarF10OpensFirstSubmenuItem(t *testing.T) {
	host := newPalettedHost()
	b := NewBar(40, barWithSubmenus())
	host.Insert(b)

	resolved := false
	r := &fakeRunner{resolve: func(box *MenuBox) {
		resolved = true
		// pick first row
		box.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnter}})
	}}
	b.Runner = r.Run

	e := &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyF10}}
	b.HandleEvent(e)
	if !resolved {
		t.Fatalf("F10 should open first menu with children")
	}
	if b.LastChosen != event.CmdUser+1 {
		t.Fatalf("LastChosen=%d want CmdUser+1", b.LastChosen)
	}
	if e.What != event.ClassNothing {
		t.Fatalf("F10 should be consumed")
	}
}

func TestBarAltLetterOpensMatchingMenu(t *testing.T) {
	host := newPalettedHost()
	b := NewBar(40, barWithSubmenus())
	host.Insert(b)

	r := &fakeRunner{resolve: func(box *MenuBox) {
		box.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnter}})
	}}
	b.Runner = r.Run

	e := &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyRune, Rune: 'e', Mod: event.ModAlt}}
	b.HandleEvent(e)
	if r.calls != 1 {
		t.Fatalf("Alt-E should run runner once: %d", r.calls)
	}
	if b.LastChosen != event.CmdUser+10 {
		t.Fatalf("LastChosen=%d want CmdUser+10", b.LastChosen)
	}
}

func TestBarMouseClickOpensMenu(t *testing.T) {
	host := newPalettedHost()
	b := NewBar(40, barWithSubmenus())
	host.Insert(b)

	r := &fakeRunner{resolve: func(box *MenuBox) {
		box.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnter}})
	}}
	b.Runner = r.Run

	// File starts at column itemX(1)
	x := b.itemX(1)
	e := &event.Event{What: event.ClassMouseDown, Mouse: event.MouseEvent{X: x + 1, Y: 0}}
	b.HandleEvent(e)
	if r.calls != 1 {
		t.Fatalf("click should open: %d", r.calls)
	}
	if b.LastChosen != event.CmdUser+1 {
		t.Fatalf("clicked menu LastChosen=%d", b.LastChosen)
	}
}

func TestBarRightCyclesToNextMenu(t *testing.T) {
	host := newPalettedHost()
	b := NewBar(40, barWithSubmenus())
	host.Insert(b)

	step := 0
	r := &fakeRunner{resolve: func(box *MenuBox) {
		step++
		switch step {
		case 1:
			// open File menu, press Right to cycle
			box.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowRight}})
		case 2:
			// now in Edit menu, pick first row
			box.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnter}})
		}
	}}
	b.Runner = r.Run

	e := &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyF10}}
	b.HandleEvent(e)

	if r.calls != 2 {
		t.Fatalf("Right should reopen runner: %d", r.calls)
	}
	if b.LastChosen != event.CmdUser+10 {
		t.Fatalf("after right cycle, LastChosen=%d", b.LastChosen)
	}
}

func TestBarEscDoesNotChooseAnything(t *testing.T) {
	host := newPalettedHost()
	b := NewBar(40, barWithSubmenus())
	host.Insert(b)

	r := &fakeRunner{resolve: func(box *MenuBox) {
		box.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEsc}})
	}}
	b.Runner = r.Run

	e := &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyF10}}
	b.HandleEvent(e)
	if b.LastChosen != event.CmdNone {
		t.Fatalf("Esc must not pick: %d", b.LastChosen)
	}
}

func TestBarNoRunnerNoOp(t *testing.T) {
	host := newPalettedHost()
	b := NewBar(40, barWithSubmenus())
	host.Insert(b)
	// Runner is nil

	e := &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyF10}}
	b.HandleEvent(e)
	if b.LastChosen != event.CmdNone {
		t.Fatalf("nil runner should not pick: %d", b.LastChosen)
	}
}
