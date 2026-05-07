// Package app wires the view tree to a backend Screen and runs the main
// event loop: read events, dispatch them down the tree, draw, flush. v0.1
// hosts an Application with a MenuBar, Desktop, and StatusLine and quits on
// Alt-X or the CmdQuit broadcast.
package app

import (
	"sync"
	"sync/atomic"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/menu"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Application is the root Group of the view tree. It holds the screen, the
// composition surface, and the channels that feed the main loop.
type Application struct {
	*view.Group

	screen  vio.Screen
	surface *vio.Surface

	bar     *menu.Bar
	desktop *Desktop
	status  *menu.Line

	posted chan event.Event
	quit   atomic.Bool

	// finishedCh is closed by Run when the loop exits. Tests use it to
	// synchronize on shutdown.
	finishedCh chan struct{}
	once       sync.Once
}

// New constructs an Application backed by screen. It sizes itself to the
// screen, attaches the standard MenuBar/Desktop/StatusLine triplet, and
// installs the BorlandClassic palette.
func New(screen vio.Screen) *Application {
	w, h := screen.Size()
	g := view.NewGroup(vio.Rect{X: 0, Y: 0, W: w, H: h})
	g.GrowMode = view.GrowAll
	g.Palette = vio.BorlandClassic

	a := &Application{
		Group:      g,
		screen:     screen,
		surface:    vio.NewSurface(w, h),
		posted:     make(chan event.Event, 64),
		finishedCh: make(chan struct{}),
	}
	a.bar = menu.NewBar(w, menu.DefaultItems)
	a.desktop = NewDesktop(vio.Rect{X: 0, Y: 1, W: w, H: h - 2})
	a.status = menu.NewLine(w, h-1, menu.DefaultHints)
	a.Insert(a.bar)
	a.Insert(a.desktop)
	a.Insert(a.status)
	return a
}

// Desktop returns the desktop child for tests and host applications that
// want to insert windows into it.
func (a *Application) Desktop() *Desktop { return a.desktop }

// MenuBar returns the menu bar child.
func (a *Application) MenuBar() *menu.Bar { return a.bar }

// StatusLine returns the status line child.
func (a *Application) StatusLine() *menu.Line { return a.status }

// Surface returns the composition surface, primarily for tests.
func (a *Application) Surface() *vio.Surface { return a.surface }

// PutEvent posts an event into the application from any goroutine. It is
// the only thread-safe entry point into the UI; everything else expects to
// run on the main loop goroutine.
func (a *Application) PutEvent(e event.Event) {
	if a.quit.Load() {
		return
	}
	select {
	case a.posted <- e:
	default:
		// Drop on overflow rather than blocking the producer. Posted events
		// are best-effort signaling; a busy UI is the loud failure mode.
	}
}

// Quit asks the loop to stop after the current iteration.
func (a *Application) Quit() { a.quit.Store(true) }

// Run drives the event loop until Quit is called or the screen channel
// closes. It returns the error returned by Screen.Show, or nil.
func (a *Application) Run() error {
	a.draw()
	if err := a.screen.Show(a.surface); err != nil {
		return err
	}

	events := a.screen.Events()
	for !a.quit.Load() {
		select {
		case e, ok := <-events:
			if !ok {
				return nil
			}
			a.process(e)
		case e := <-a.posted:
			a.process(e)
		}

		if a.quit.Load() {
			break
		}

		a.draw()
		if err := a.screen.Show(a.surface); err != nil {
			return err
		}
	}
	a.once.Do(func() { close(a.finishedCh) })
	return nil
}

// Finished returns a channel that is closed when Run exits.
func (a *Application) Finished() <-chan struct{} { return a.finishedCh }

// process is the per-event handler. It applies framework-level shortcuts
// (resize, Alt-X, CmdQuit) before delegating to the Group dispatch.
func (a *Application) process(e event.Event) {
	if e.What == event.ClassResize {
		a.resize(e.Mouse.X, e.Mouse.Y)
		return
	}
	if isQuitKey(e) || isQuitCommand(e) {
		a.Quit()
		return
	}
	ev := e
	a.HandleEvent(&ev)
}

func isQuitKey(e event.Event) bool {
	if e.What != event.ClassKey {
		return false
	}
	if e.Key.Mod&event.ModAlt == 0 {
		return false
	}
	if e.Key.Key != event.KeyRune {
		return false
	}
	return e.Key.Rune == 'x' || e.Key.Rune == 'X'
}

func isQuitCommand(e event.Event) bool {
	if e.What != event.ClassCommand && e.What != event.ClassBroadcast {
		return false
	}
	return e.Msg.Command == event.CmdQuit
}

func (a *Application) resize(w, h int) {
	if w <= 0 || h <= 0 {
		return
	}
	a.surface.Resize(w, h)
	a.ChangeBounds(vio.Rect{X: 0, Y: 0, W: w, H: h})
	a.bar.ChangeBounds(vio.Rect{X: 0, Y: 0, W: w, H: 1})
	a.desktop.ChangeBounds(vio.Rect{X: 0, Y: 1, W: w, H: h - 2})
	a.status.ChangeBounds(vio.Rect{X: 0, Y: h - 1, W: w, H: 1})
}

func (a *Application) draw() {
	w, h := a.screen.Size()
	if sw, sh := a.surface.Size(); sw != w || sh != h {
		a.surface.Resize(w, h)
	}
	a.surface.Clear(a.Palette.Map(1))
	a.Draw(a.surface)
}
