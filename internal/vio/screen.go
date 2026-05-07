package vio

import "github.com/tamnd/vigo/internal/event"

// Screen is the contract every backend implements: a sized cell grid, an
// event source, and a flush primitive. The vio package ships two
// implementations: a tcell-backed real screen and an in-memory fake screen
// suitable for tests.
type Screen interface {
	// Init opens the terminal, switches to the alternate screen, hides the
	// cursor, and starts the event pump. It is illegal to call any other
	// method before Init returns nil.
	Init() error

	// Fini restores the terminal state. It is safe to call multiple times.
	Fini()

	// Size returns the current screen dimensions.
	Size() (w, h int)

	// Show flushes the supplied surface to the terminal. The surface size
	// must match Screen.Size().
	Show(s *Surface) error

	// Events returns the channel of input events. The channel is closed when
	// the screen is finalized.
	Events() <-chan event.Event

	// SetCursor positions the hardware cursor. visible == false hides it.
	SetCursor(x, y int, visible bool)
}
