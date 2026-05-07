package app

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// recordingView counts the events it receives so tests can assert that the
// modal slot redirects events away from the rest of the tree.
type recordingView struct {
	*view.View
	count int
}

func newRecordingView() *recordingView {
	v := view.NewView(vio.R(0, 0, 1, 1))
	v.Options |= view.OptSelectable
	return &recordingView{View: v}
}

func (r *recordingView) HandleEvent(e *event.Event) {
	r.count++
	e.Clear()
}

func TestApplicationSetModalRoutesEventsToModal(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	modal := newRecordingView()
	a.SetModal(modal)
	if a.Modal() != modal {
		t.Fatal("Modal accessor did not return the slot")
	}

	a.process(event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyEnter},
	})
	if modal.count != 1 {
		t.Fatalf("modal should have received the event: %d", modal.count)
	}
}

func TestApplicationModalSuppressesDesktopDispatch(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	desktopSpy := newRecordingView()
	a.Desktop().Insert(desktopSpy)

	modal := newRecordingView()
	a.SetModal(modal)

	a.process(event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyEnter},
	})

	if desktopSpy.count != 0 {
		t.Fatalf("desktop child should not see events while modal active: %d", desktopSpy.count)
	}
	if modal.count != 1 {
		t.Fatalf("modal should have seen the event: %d", modal.count)
	}
}

func TestApplicationClearModalRestoresDispatch(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	desktopSpy := newRecordingView()
	a.Desktop().Insert(desktopSpy)

	modal := newRecordingView()
	a.SetModal(modal)
	a.ClearModal()
	if a.Modal() != nil {
		t.Fatal("ClearModal did not clear the slot")
	}

	a.process(event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyEnter},
	})
	if desktopSpy.count != 1 {
		t.Fatalf("after ClearModal the desktop child should see events: %d", desktopSpy.count)
	}
	if modal.count != 0 {
		t.Fatalf("modal should not see events after ClearModal: %d", modal.count)
	}
}

func TestApplicationModalStillHandlesQuit(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	modal := newRecordingView()
	a.SetModal(modal)

	// Alt-X should still quit even with a modal active; the system shortcut
	// runs before modal routing.
	a.process(event.Event{
		What: event.ClassKey,
		Key: event.KeyEvent{
			Key: event.KeyRune, Rune: 'x', Mod: event.ModAlt,
		},
	})
	if !a.quit.Load() {
		t.Fatal("Alt-X should quit even while modal is active")
	}
	if modal.count != 0 {
		t.Fatal("Alt-X must not be delivered to modal")
	}
}
