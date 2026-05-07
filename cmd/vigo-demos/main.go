// Command vigo-demos boots the vigo desktop with the four classic
// Turbo Vision sample tools (Calculator, Calendar, ASCIITable,
// Puzzle) reachable from the menu bar. It exists to exercise the
// framework end-to-end and to give the v0.2 release something the
// reader can actually click around in.
package main

import (
	"fmt"
	"os"

	"github.com/tamnd/vigo/app"
	"github.com/tamnd/vigo/demos"
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/menu"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Demo command IDs. They sit in the user range so they cannot clash
// with framework-level commands.
const (
	cmdCalculator = event.CmdUser + 1
	cmdCalendar   = event.CmdUser + 2
	cmdASCIITable = event.CmdUser + 3
	cmdPuzzle     = event.CmdUser + 4
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "vigo-demos: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	screen := vio.NewTcellScreen()
	if err := screen.Init(); err != nil {
		return fmt.Errorf("init screen: %w", err)
	}
	defer screen.Fini()

	a := app.New(screen)
	a.MenuBar().Items = demoMenuItems()
	a.Desktop().Insert(newLauncherSink(a))
	return a.Run()
}

// demoMenuItems builds the launcher's menu tree: a File menu with
// the four demo entries plus Exit.
func demoMenuItems() []menu.Item {
	demos := []menu.Item{
		{Title: "~C~alculator", Cmd: cmdCalculator, Hotkey: 0},
		{Title: "C~a~lendar", Cmd: cmdCalendar, Hotkey: 1},
		{Title: "~A~SCII Table", Cmd: cmdASCIITable, Hotkey: 0},
		{Title: "~P~uzzle", Cmd: cmdPuzzle, Hotkey: 0},
		{Sep: true},
		{Title: "E~x~it", Hotkey: 1, Cmd: event.CmdQuit, Shortcut: "Alt-X"},
	}
	return []menu.Item{
		{Title: "≡", Hotkey: -1},
		{Title: "File", Hotkey: 0, Children: demos},
	}
}

// launcherSink is a hidden post-process desktop child that watches
// for demo commands and inserts the corresponding window.
type launcherSink struct {
	*view.View

	app *app.Application
}

func newLauncherSink(a *app.Application) *launcherSink {
	v := view.NewView(vio.Rect{})
	v.Options |= view.OptPostProcess
	v.EventMask = event.ClassCommand | event.ClassBroadcast
	return &launcherSink{View: v, app: a}
}

func (s *launcherSink) Draw(*vio.Surface) {}

func (s *launcherSink) HandleEvent(e *event.Event) {
	if e.What != event.ClassCommand && e.What != event.ClassBroadcast {
		return
	}
	switch e.Msg.Command {
	case cmdCalculator:
		s.app.Desktop().Insert(demos.NewCalculator(centerIn(s.app.Desktop().Bounds, demos.CalcDefaultBounds)))
	case cmdCalendar:
		s.app.Desktop().Insert(demos.NewCalendar(centerIn(s.app.Desktop().Bounds, demos.CalDefaultBounds)))
	case cmdASCIITable:
		s.app.Desktop().Insert(demos.NewASCIITable(centerIn(s.app.Desktop().Bounds, demos.ASCIIDefaultBounds)))
	case cmdPuzzle:
		s.app.Desktop().Insert(demos.NewPuzzle(centerIn(s.app.Desktop().Bounds, demos.PuzzleDefaultBounds)))
	default:
		return
	}
	e.Clear()
}

// centerIn places the inner rectangle at the center of outer,
// preserving inner's width and height.
func centerIn(outer, inner vio.Rect) vio.Rect {
	w, h := inner.W, inner.H
	if w > outer.W {
		w = outer.W
	}
	if h > outer.H {
		h = outer.H
	}
	return vio.Rect{
		X: outer.X + (outer.W-w)/2,
		Y: outer.Y + (outer.H-h)/2,
		W: w,
		H: h,
	}
}
