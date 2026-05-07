package main

import (
	"testing"

	"github.com/tamnd/vigo/demos"
	"github.com/tamnd/vigo/event"
)

// Verifies that after the launcher inserts a Calculator window,
// keypresses delivered to the Desktop reach the calculator's handler
// and update its display. This is the end-to-end smoke test for the
// "apps feel like mockups" complaint.
func TestCalculatorReceivesKeysAfterLauncher(t *testing.T) {
	a, f := newDemoApp(t)
	defer f.Fini()

	a.Desktop().HandleEvent(&event.Event{
		What: event.ClassCommand,
		Msg:  event.MessageEvent{Command: cmdCalculator},
	})

	children := a.Desktop().Children()
	calc, ok := children[len(children)-1].(*demos.Calculator)
	if !ok {
		t.Fatalf("last child=%T want *Calculator", children[len(children)-1])
	}

	for _, r := range "2+3=" {
		a.Desktop().HandleEvent(&event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyRune, Rune: r},
		})
	}
	if got := calc.Display(); got != "5" {
		t.Fatalf("display=%q want 5 after 2+3=", got)
	}
}
