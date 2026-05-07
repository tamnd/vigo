package widget

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func TestHistoryRingDedupesAndOrders(t *testing.T) {
	HistoryReset()
	const id HistoryID = 1

	HistoryAdd(id, "alpha")
	HistoryAdd(id, "beta")
	HistoryAdd(id, "gamma")
	HistoryAdd(id, "alpha") // moves alpha to front

	if got, want := HistoryCount(id), 3; got != want {
		t.Fatalf("count=%d want %d", got, want)
	}
	if got := HistoryStr(id, 0); got != "alpha" {
		t.Fatalf("front=%q", got)
	}
	if got := HistoryStr(id, 1); got != "gamma" {
		t.Fatalf("second=%q", got)
	}
	if got := HistoryStr(id, 2); got != "beta" {
		t.Fatalf("third=%q", got)
	}
}

func TestHistoryRingCapsAtLimit(t *testing.T) {
	HistoryReset()
	const id HistoryID = 2

	for i := range HistoryLimit + 5 {
		HistoryAdd(id, string(rune('a'+i)))
	}
	if HistoryCount(id) != HistoryLimit {
		t.Fatalf("count=%d want %d", HistoryCount(id), HistoryLimit)
	}
}

func TestHistoryRingIgnoresEmptyAndZeroID(t *testing.T) {
	HistoryReset()
	HistoryAdd(0, "x")
	HistoryAdd(7, "")
	if HistoryCount(0) != 0 {
		t.Fatalf("zero id should be ignored")
	}
	if HistoryCount(7) != 0 {
		t.Fatalf("empty string should be ignored")
	}
}

func TestHistoryStrOutOfRange(t *testing.T) {
	HistoryReset()
	HistoryAdd(3, "only")
	if got := HistoryStr(3, -1); got != "" {
		t.Fatalf("negative index: %q", got)
	}
	if got := HistoryStr(3, 5); got != "" {
		t.Fatalf("past end: %q", got)
	}
}

func TestHistoryButtonOpensOnClickAndEnter(t *testing.T) {
	HistoryReset()
	host := newPalettedHost()
	h := NewHistory(vio.R(10, 0, 3, 1), 4)
	host.Insert(h)
	host.SetCurrent(0)

	h.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 11, Y: 0},
	})
	if !h.Open {
		t.Fatalf("click should open dropdown")
	}

	h.Open = false
	h.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyEnter},
	})
	if !h.Open {
		t.Fatalf("Enter should open dropdown")
	}
}

func TestHistoryWindowSeedsListBox(t *testing.T) {
	HistoryReset()
	HistoryAdd(5, "first")
	HistoryAdd(5, "second")
	HistoryAdd(5, "third")

	w := NewHistoryWindow(vio.R(0, 0, 20, 5), 5)
	if got := len(w.List.Items()); got != 3 {
		t.Fatalf("seeded items=%d", got)
	}
	if got := w.List.Items()[0]; got != "third" {
		t.Fatalf("front=%q", got)
	}
}

func TestHistoryWindowEnterPicksFocused(t *testing.T) {
	HistoryReset()
	HistoryAdd(6, "alpha")
	HistoryAdd(6, "beta")
	HistoryAdd(6, "gamma")

	w := NewHistoryWindow(vio.R(0, 0, 20, 5), 6)
	w.List.SetState(view.StateFocused, true)
	w.List.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyArrowDown}})
	w.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEnter}})

	if w.Result() != event.CmdOk {
		t.Fatalf("result=%d", w.Result())
	}
	if w.Selected != "beta" {
		t.Fatalf("selected=%q", w.Selected)
	}
}

func TestHistoryWindowEscCancels(t *testing.T) {
	HistoryReset()
	HistoryAdd(7, "x")
	w := NewHistoryWindow(vio.R(0, 0, 20, 5), 7)

	w.HandleEvent(&event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyEsc}})
	if w.Result() != event.CmdCancel {
		t.Fatalf("Esc result=%d", w.Result())
	}
	if w.Selected != "" {
		t.Fatalf("Esc should not store selection: %q", w.Selected)
	}
}
