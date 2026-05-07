package app

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/help"
)

func TestIsHelpKeyOnlyMatchesF1(t *testing.T) {
	if !isHelpKey(event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyF1},
	}) {
		t.Fatal("F1 should be a help key")
	}
	if isHelpKey(event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyF2},
	}) {
		t.Fatal("F2 must not be a help key")
	}
	if isHelpKey(event.Event{What: event.ClassMouseDown}) {
		t.Fatal("mouse events must not be help keys")
	}
}

func TestIsHelpCommandMatchesCmdHelp(t *testing.T) {
	if !isHelpCommand(event.Event{
		What: event.ClassCommand,
		Msg:  event.MessageEvent{Command: event.CmdHelp},
	}) {
		t.Fatal("ClassCommand+CmdHelp should match")
	}
	if !isHelpCommand(event.Event{
		What: event.ClassBroadcast,
		Msg:  event.MessageEvent{Command: event.CmdHelp},
	}) {
		t.Fatal("ClassBroadcast+CmdHelp should match")
	}
	if isHelpCommand(event.Event{
		What: event.ClassCommand,
		Msg:  event.MessageEvent{Command: event.CmdQuit},
	}) {
		t.Fatal("CmdQuit must not be a help command")
	}
}

func TestApplicationSeedsHelpRegistry(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	r := a.HelpRegistry()
	if r == nil {
		t.Fatal("HelpRegistry should return the seeded registry")
	}
	if _, ok := r.Lookup(help.CtxAbout); !ok {
		t.Fatal("seeded registry must contain CtxAbout")
	}
}

func TestApplicationF1OpensAboutThenAltXQuits(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	go func() {
		// F1 opens the About dialog (nested ExecView).
		f.PushEvent(event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyF1},
		})
		// CmdOk broadcast closes the dialog via its closeSink.
		f.PushEvent(event.Event{
			What: event.ClassCommand,
			Msg:  event.MessageEvent{Command: event.CmdOk},
		})
		// Alt-X exits the outer Run.
		f.PushEvent(event.Event{
			What: event.ClassKey,
			Key: event.KeyEvent{
				Key: event.KeyRune, Rune: 'x', Mod: event.ModAlt,
			},
		})
	}()
	if err := a.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestApplicationCmdHelpOpensAbout(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	go func() {
		f.PushEvent(event.Event{
			What: event.ClassCommand,
			Msg:  event.MessageEvent{Command: event.CmdHelp},
		})
		f.PushEvent(event.Event{
			What: event.ClassCommand,
			Msg:  event.MessageEvent{Command: event.CmdOk},
		})
		f.PushEvent(event.Event{
			What: event.ClassKey,
			Key: event.KeyEvent{
				Key: event.KeyRune, Rune: 'x', Mod: event.ModAlt,
			},
		})
	}()
	if err := a.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
}
