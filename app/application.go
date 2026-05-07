// Package app wires the view tree to a backend Screen and runs the main
// event loop: read events, dispatch them down the tree, draw, flush. v0.1
// hosts an Application with a MenuBar, Desktop, and StatusLine and quits on
// Alt-X or the CmdQuit broadcast.
package app

import (
	"sync"
	"sync/atomic"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/help"
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

	help *help.Registry

	posted chan event.Event
	quit   atomic.Bool

	modal view.Viewer

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
		help:       help.Default(),
		posted:     make(chan event.Event, 64),
		finishedCh: make(chan struct{}),
	}
	a.bar = menu.NewBar(w, menu.DefaultItems)
	a.desktop = NewDesktop(vio.Rect{X: 0, Y: 1, W: w, H: h - 2})
	a.status = menu.NewLine(w, h-1, menu.DefaultHints)
	a.Insert(a.bar)
	a.Insert(a.desktop)
	a.Insert(a.status)
	a.bar.Runner = func(box *menu.MenuBox) event.CommandID { return a.ExecView(box) }
	a.bar.OnCommand = func(c event.CommandID) {
		a.PutEvent(event.Event{
			What: event.ClassCommand,
			Msg:  event.MessageEvent{Command: c},
		})
	}
	return a
}

// Desktop returns the desktop child for tests and host applications that
// want to insert windows into it.
func (a *Application) Desktop() *Desktop { return a.desktop }

// MenuBar returns the menu bar child.
func (a *Application) MenuBar() *menu.Bar { return a.bar }

// StatusLine returns the status line child.
func (a *Application) StatusLine() *menu.Line { return a.status }

// HelpRegistry returns the help topic registry. Hosts register their
// own contexts on it before Run; v0.2 only seeds CtxAbout.
func (a *Application) HelpRegistry() *help.Registry { return a.help }

// ShowAbout opens the About dialog as a modal sub-loop.
func (a *Application) ShowAbout() event.CommandID {
	return a.ExecView(help.About(a.help))
}

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

// SetModal pushes v as the currently active modal view. While set, only v
// receives non-system events; the rest of the tree (menu bar, desktop,
// status line) sees ClassNothing. Pass nil to clear the modal.
//
// The host is responsible for inserting v into a Group (typically the
// desktop) before calling SetModal, and removing it after ClearModal.
// Phase 2 ships this minimal slot; Dialog/MessageBox in phase 3 will wrap
// it in an ExecView convenience.
func (a *Application) SetModal(v view.Viewer) { a.modal = v }

// ClearModal removes any active modal slot.
func (a *Application) ClearModal() { a.modal = nil }

// Modal returns the currently active modal view, or nil.
func (a *Application) Modal() view.Viewer { return a.modal }

// ModalView is implemented by views ExecView can drive: each iteration
// the loop checks Result; a non-zero value ends the sub-loop and is
// returned. Dialog satisfies this interface.
type ModalView interface {
	view.Viewer
	Result() event.CommandID
}

// ExecView shows v as a modal sub-loop and blocks until it closes.
// The view is inserted into the desktop, set as the modal slot, and
// drawn on every iteration. The loop exits when v.Result() returns a
// non-zero command, when the screen channel closes, or when Quit is
// called; it returns the close command (CmdCancel as a fall-back).
func (a *Application) ExecView(v ModalView) event.CommandID {
	if v == nil {
		return event.CmdCancel
	}
	a.desktop.Insert(v)
	a.SetModal(v)
	defer func() {
		a.ClearModal()
		a.desktop.Remove(v)
	}()

	a.draw()
	if err := a.screen.Show(a.surface); err != nil {
		return event.CmdCancel
	}

	events := a.screen.Events()
	for !a.quit.Load() {
		if r := v.Result(); r != event.CmdNone {
			return r
		}
		select {
		case e, ok := <-events:
			if !ok {
				return event.CmdCancel
			}
			a.process(e)
		case e := <-a.posted:
			a.process(e)
		}
		if r := v.Result(); r != event.CmdNone {
			return r
		}
		a.draw()
		if err := a.screen.Show(a.surface); err != nil {
			return event.CmdCancel
		}
	}
	return event.CmdCancel
}

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
// (resize, Alt-X, CmdQuit), routes to the modal slot if one is active, and
// otherwise delegates to the Group dispatch.
func (a *Application) process(e event.Event) {
	if e.What == event.ClassResize {
		a.resize(e.Mouse.X, e.Mouse.Y)
		return
	}
	if isQuitKey(e) || isQuitCommand(e) {
		a.Quit()
		return
	}
	if isHelpKey(e) || isHelpCommand(e) {
		a.ShowAbout()
		return
	}
	ev := e
	if a.modal != nil {
		a.modal.HandleEvent(&ev)
		return
	}
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

func isHelpKey(e event.Event) bool {
	if e.What != event.ClassKey {
		return false
	}
	return e.Key.Key == event.KeyF1
}

func isHelpCommand(e event.Event) bool {
	if e.What != event.ClassCommand && e.What != event.ClassBroadcast {
		return false
	}
	return e.Msg.Command == event.CmdHelp
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
