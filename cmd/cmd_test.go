package cmd

import (
	"testing"

	"github.com/tamnd/vigo/event"
)

func TestSetAddRemoveHas(t *testing.T) {
	s := NewSet(event.CmdOk, event.CmdCancel)
	if !s.Has(event.CmdOk) {
		t.Fatalf("seeded set missing CmdOk")
	}
	if s.Has(event.CmdQuit) {
		t.Fatalf("set should not have CmdQuit")
	}
	s.Add(event.CmdQuit)
	if !s.Has(event.CmdQuit) {
		t.Fatalf("Add did not register CmdQuit")
	}
	s.Remove(event.CmdOk)
	if s.Has(event.CmdOk) {
		t.Fatalf("Remove did not delete CmdOk")
	}
}

func TestEnablerEnableDisable(t *testing.T) {
	en := NewEnabler(NewSet(event.CmdOk))
	if !en.IsEnabled(event.CmdOk) {
		t.Fatalf("seeded CmdOk should be enabled")
	}
	en.Enable(NewSet(event.CmdQuit))
	if !en.IsEnabled(event.CmdQuit) {
		t.Fatalf("Enable did not turn on CmdQuit")
	}
	en.Disable(NewSet(event.CmdOk, event.CmdQuit))
	if en.IsEnabled(event.CmdOk) || en.IsEnabled(event.CmdQuit) {
		t.Fatalf("Disable left commands enabled")
	}
}

func TestEnablerSubscribeFiresOnChange(t *testing.T) {
	en := NewEnabler(nil)
	var hits int
	en.Subscribe(func() { hits++ })

	en.Enable(NewSet(event.CmdHelp))
	if hits != 1 {
		t.Fatalf("Enable should notify once: %d", hits)
	}

	en.Enable(NewSet(event.CmdHelp)) // already enabled, no change
	if hits != 1 {
		t.Fatalf("no-op Enable should not notify: %d", hits)
	}

	en.Disable(NewSet(event.CmdHelp))
	if hits != 2 {
		t.Fatalf("Disable should notify: %d", hits)
	}

	en.Disable(NewSet(event.CmdHelp)) // already disabled
	if hits != 2 {
		t.Fatalf("no-op Disable should not notify: %d", hits)
	}
}

func TestEnablerNilInitial(t *testing.T) {
	en := NewEnabler(nil)
	if en.IsEnabled(event.CmdOk) {
		t.Fatalf("nil-initial enabler should be empty")
	}
}

func TestBindingsLookup(t *testing.T) {
	b := NewBindings()
	b.Bind(event.KeyEvent{Key: event.KeyF1}, event.CmdHelp)
	b.Bind(event.KeyEvent{Key: event.KeyRune, Mod: event.ModAlt, Rune: 'x'}, event.CmdQuit)

	if got := b.Lookup(event.KeyEvent{Key: event.KeyF1}); got != event.CmdHelp {
		t.Fatalf("F1 -> %d, want CmdHelp", got)
	}
	if got := b.Lookup(event.KeyEvent{Key: event.KeyRune, Mod: event.ModAlt, Rune: 'x'}); got != event.CmdQuit {
		t.Fatalf("Alt-x -> %d, want CmdQuit", got)
	}
	// bare 'x' must not collide with Alt-x
	if got := b.Lookup(event.KeyEvent{Key: event.KeyRune, Rune: 'x'}); got != event.CmdNone {
		t.Fatalf("bare x -> %d, want CmdNone", got)
	}
	// unbound key
	if got := b.Lookup(event.KeyEvent{Key: event.KeyF12}); got != event.CmdNone {
		t.Fatalf("unbound F12 -> %d, want CmdNone", got)
	}
}

func TestBindingsRebindAndUnbind(t *testing.T) {
	b := NewBindings()
	k := event.KeyEvent{Key: event.KeyF5}
	b.Bind(k, event.CmdZoom)
	b.Bind(k, event.CmdNext) // rebind
	if got := b.Lookup(k); got != event.CmdNext {
		t.Fatalf("rebind: %d", got)
	}
	b.Unbind(k)
	if got := b.Lookup(k); got != event.CmdNone {
		t.Fatalf("after unbind: %d", got)
	}
}

func TestEnablerSubscribeUnsubscribe(t *testing.T) {
	en := NewEnabler(nil)
	var hits int
	cancel := en.Subscribe(func() { hits++ })

	en.Enable(NewSet(event.CmdOk))
	if hits != 1 {
		t.Fatalf("subscribed: %d", hits)
	}

	cancel()
	en.Disable(NewSet(event.CmdOk))
	if hits != 1 {
		t.Fatalf("after cancel: %d", hits)
	}
}
