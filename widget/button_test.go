package widget

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func TestButtonParsesMnemonicAndDrawsTitle(t *testing.T) {
	host := newPalettedHost()
	b := NewButton(vio.R(0, 0, 10, 1), "~O~K", event.CmdOk, BfDefault)
	host.Insert(b)
	if b.Mnemonic() != 'o' {
		t.Fatalf("mnemonic: %q", b.Mnemonic())
	}
	s := vio.NewSurface(10, 1)
	b.Draw(s)
	row := s.Snapshot()[0]
	if !strings.Contains(row, "<OK>") {
		t.Fatalf("default button should render <OK>, got %q", row)
	}
}

func TestButtonDefaultFlag(t *testing.T) {
	b := NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, BfDefault)
	if !b.IsDefault() {
		t.Fatal("BfDefault flag should be reflected in IsDefault")
	}
	plain := NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, BfNormal)
	if plain.IsDefault() {
		t.Fatal("plain button should not be default")
	}
}

func TestButtonEnterFiresCommandWhenFocused(t *testing.T) {
	host := newPalettedHost()
	b := NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, BfDefault)
	host.Insert(b)
	host.SetCurrent(0)

	got := captureCommand(host, func() {
		b.HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyEnter},
		})
	})
	if got != event.CmdOk {
		t.Fatalf("Enter should fire CmdOk, got %v", got)
	}
}

func TestButtonSpaceFiresWhenFocused(t *testing.T) {
	host := newPalettedHost()
	b := NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, BfNormal)
	host.Insert(b)
	host.SetCurrent(0)

	got := captureCommand(host, func() {
		b.HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyRune, Rune: ' '},
		})
	})
	if got != event.CmdOk {
		t.Fatalf("Space should fire CmdOk, got %v", got)
	}
}

func TestButtonAltMnemonicFiresWithoutFocus(t *testing.T) {
	host := newPalettedHost()
	b := NewButton(vio.R(0, 0, 10, 1), "~O~K", event.CmdOk, BfNormal)
	host.Insert(b)
	host.SetCurrent(-1)

	got := captureCommand(host, func() {
		b.HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'o', Mod: event.ModAlt},
		})
	})
	if got != event.CmdOk {
		t.Fatalf("Alt-O should fire CmdOk, got %v", got)
	}
}

func TestButtonMouseDownInsideFires(t *testing.T) {
	host := newPalettedHost()
	b := NewButton(vio.R(2, 1, 6, 1), "OK", event.CmdOk, BfNormal)
	host.Insert(b)

	got := captureCommand(host, func() {
		b.HandleEvent(&event.Event{
			What:  event.ClassMouseDown,
			Mouse: event.MouseEvent{X: 4, Y: 1},
		})
	})
	if got != event.CmdOk {
		t.Fatalf("mouse-down inside should fire CmdOk, got %v", got)
	}
}

func TestButtonMouseDownOutsideIgnored(t *testing.T) {
	host := newPalettedHost()
	b := NewButton(vio.R(2, 1, 6, 1), "OK", event.CmdOk, BfNormal)
	host.Insert(b)

	got := captureCommand(host, func() {
		b.HandleEvent(&event.Event{
			What:  event.ClassMouseDown,
			Mouse: event.MouseEvent{X: 0, Y: 0},
		})
	})
	if got != event.CmdNone {
		t.Fatalf("mouse-down outside should not fire, got %v", got)
	}
}

func TestButtonDisabledIgnoresEvents(t *testing.T) {
	host := newPalettedHost()
	b := NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, BfDefault)
	host.Insert(b)
	host.SetCurrent(0)
	b.SetState(view.StateDisabled, true)

	got := captureCommand(host, func() {
		b.HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyEnter},
		})
	})
	if got != event.CmdNone {
		t.Fatal("disabled button must not fire")
	}
}

func TestButtonBroadcastFlagEmitsBroadcast(t *testing.T) {
	host := newPalettedHost()
	spy := newCommandSpy()
	host.Insert(spy)
	b := NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, BfBroadcast)
	host.Insert(b)
	host.SetCurrent(1)

	b.Press()
	if spy.lastClass != event.ClassBroadcast {
		t.Fatalf("expected ClassBroadcast, got %v", spy.lastClass)
	}
	if spy.lastCommand != event.CmdOk {
		t.Fatalf("expected CmdOk, got %v", spy.lastCommand)
	}
}

// commandSpy is a post-process child that records commands and broadcasts
// dispatched into its owner Group, used by button tests to verify
// fire-and-forget semantics.
type commandSpy struct {
	*view.View

	lastClass   event.Class
	lastCommand event.CommandID
}

func newCommandSpy() *commandSpy {
	v := view.NewView(vio.R(0, 0, 1, 1))
	v.Options |= view.OptPostProcess
	v.EventMask = event.ClassCommand | event.ClassBroadcast
	return &commandSpy{View: v}
}

func (s *commandSpy) HandleEvent(e *event.Event) {
	if e.What != event.ClassCommand && e.What != event.ClassBroadcast {
		return
	}
	s.lastClass = e.What
	s.lastCommand = e.Msg.Command
	e.Clear()
}

// captureCommand inserts a fresh spy into host, runs fn, and returns the
// command observed (or CmdNone if nothing fired).
func captureCommand(host *view.Group, fn func()) event.CommandID {
	spy := newCommandSpy()
	host.Insert(spy)
	defer host.Remove(spy)
	fn()
	return spy.lastCommand
}
