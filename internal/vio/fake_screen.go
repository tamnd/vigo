package vio

import (
	"sync"

	"github.com/tamnd/vigo/internal/event"
)

// FakeScreen is an in-memory Screen for tests. It keeps a Surface that
// reflects the latest Show() call, exposes a queueable event channel via
// PushEvent, and never blocks on terminal I/O.
type FakeScreen struct {
	w, h   int
	mu     sync.Mutex
	last   *Surface
	cx, cy int
	cvis   bool
	events chan event.Event
	done   chan struct{}
	closed bool
}

// NewFakeScreen returns a fake screen of the given size.
func NewFakeScreen(w, h int) *FakeScreen {
	return &FakeScreen{
		w: w, h: h,
		events: make(chan event.Event, 64),
		done:   make(chan struct{}),
	}
}

// Init is a no-op; FakeScreen is ready to use as soon as it is constructed.
func (f *FakeScreen) Init() error { return nil }

// Fini closes the event channel, signaling any consumer to exit.
func (f *FakeScreen) Fini() {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return
	}
	f.closed = true
	close(f.done)
	close(f.events)
}

// Size returns the configured dimensions.
func (f *FakeScreen) Size() (int, int) { return f.w, f.h }

// Show records the supplied surface for inspection by tests.
func (f *FakeScreen) Show(s *Surface) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	w, h := s.Size()
	if f.last == nil {
		f.last = NewSurface(w, h)
	} else if lw, lh := f.last.Size(); lw != w || lh != h {
		f.last.Resize(w, h)
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := s.At(x, y)
			f.last.Set(x, y, c.Rune, c.Attr)
		}
	}
	return nil
}

// Events returns the queue of pushed events.
func (f *FakeScreen) Events() <-chan event.Event { return f.events }

// SetCursor records the cursor position for inspection.
func (f *FakeScreen) SetCursor(x, y int, visible bool) {
	f.mu.Lock()
	f.cx, f.cy, f.cvis = x, y, visible
	f.mu.Unlock()
}

// Cursor returns the last cursor state set by Application.
func (f *FakeScreen) Cursor() (x, y int, visible bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.cx, f.cy, f.cvis
}

// PushEvent queues an event to be delivered to the consumer. Tests can drive
// an Application by pushing key/mouse events.
func (f *FakeScreen) PushEvent(e event.Event) {
	f.mu.Lock()
	closed := f.closed
	f.mu.Unlock()
	if closed {
		return
	}
	f.events <- e
}

// Surface returns the surface as of the last Show call. The returned surface
// is shared with the screen and callers must not mutate it while tests are
// running, but for read-only snapshots (Snapshot, At) it is safe.
func (f *FakeScreen) Surface() *Surface {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.last
}

// Resize changes the screen dimensions and emits a synthetic resize event.
func (f *FakeScreen) Resize(w, h int) {
	f.mu.Lock()
	f.w, f.h = w, h
	f.mu.Unlock()
	f.PushEvent(event.Event{
		What:  event.ClassResize,
		Mouse: event.MouseEvent{X: w, Y: h},
	})
}
