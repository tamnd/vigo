package main

import (
	"reflect"
	"testing"

	"github.com/tamnd/vigo/app"
	"github.com/tamnd/vigo/demos"
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
)

func newDemoApp(t *testing.T) (*app.Application, *vio.FakeScreen) {
	t.Helper()
	f := vio.NewFakeScreen(80, 24)
	if err := f.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	a := app.New(f)
	a.MenuBar().Items = demoMenuItems()
	a.Desktop().Insert(newLauncherSink(a))
	return a, f
}

func TestDemoMenuListsAllFourTools(t *testing.T) {
	items := demoMenuItems()
	if len(items) != 2 {
		t.Fatalf("top-level menus=%d want 2 (system + File)", len(items))
	}
	file := items[1]
	if file.Title != "File" {
		t.Fatalf("second item=%q want File", file.Title)
	}
	cmds := []event.CommandID{}
	for _, ch := range file.Children {
		if ch.Sep {
			continue
		}
		cmds = append(cmds, ch.Cmd)
	}
	want := []event.CommandID{cmdCalculator, cmdCalendar, cmdASCIITable, cmdPuzzle, event.CmdQuit}
	if !reflect.DeepEqual(cmds, want) {
		t.Fatalf("file children cmds=%v want %v", cmds, want)
	}
}

func TestLauncherSinkOpensCalculator(t *testing.T) {
	a, f := newDemoApp(t)
	defer f.Fini()

	before := len(a.Desktop().Children())
	a.Desktop().HandleEvent(&event.Event{
		What: event.ClassCommand,
		Msg:  event.MessageEvent{Command: cmdCalculator},
	})
	after := len(a.Desktop().Children())
	if after != before+1 {
		t.Fatalf("calculator did not insert: %d -> %d", before, after)
	}
	last := a.Desktop().Children()[after-1]
	if _, ok := last.(*demos.Calculator); !ok {
		t.Fatalf("inserted child=%T want *Calculator", last)
	}
}

func TestLauncherSinkOpensEachTool(t *testing.T) {
	cases := []struct {
		cmd  event.CommandID
		kind string
	}{
		{cmdCalculator, "*demos.Calculator"},
		{cmdCalendar, "*demos.Calendar"},
		{cmdASCIITable, "*demos.ASCIITable"},
		{cmdPuzzle, "*demos.Puzzle"},
	}
	a, f := newDemoApp(t)
	defer f.Fini()

	for _, tc := range cases {
		before := len(a.Desktop().Children())
		a.Desktop().HandleEvent(&event.Event{
			What: event.ClassCommand,
			Msg:  event.MessageEvent{Command: tc.cmd},
		})
		after := len(a.Desktop().Children())
		if after != before+1 {
			t.Fatalf("cmd %d did not insert a window", tc.cmd)
		}
	}
}

func TestLauncherSinkIgnoresUnknownCommand(t *testing.T) {
	a, f := newDemoApp(t)
	defer f.Fini()

	before := len(a.Desktop().Children())
	a.Desktop().HandleEvent(&event.Event{
		What: event.ClassCommand,
		Msg:  event.MessageEvent{Command: event.CmdUser + 99},
	})
	if got := len(a.Desktop().Children()); got != before {
		t.Fatalf("unknown command spawned a window: %d -> %d", before, got)
	}
}

func TestCenterInClampsOversizeRect(t *testing.T) {
	outer := vio.R(0, 0, 20, 10)
	inner := vio.R(0, 0, 100, 100)
	got := centerIn(outer, inner)
	if got.W != outer.W || got.H != outer.H {
		t.Fatalf("clamped=%+v want %dx%d", got, outer.W, outer.H)
	}
}
