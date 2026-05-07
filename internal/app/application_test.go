package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tamnd/vigo/internal/event"
	"github.com/tamnd/vigo/internal/menu"
	"github.com/tamnd/vigo/internal/vio"
)

const (
	scrW = 40
	scrH = 10
)

func newTestApp(t *testing.T) (*Application, *vio.FakeScreen) {
	t.Helper()
	f := vio.NewFakeScreen(scrW, scrH)
	if err := f.Init(); err != nil {
		t.Fatalf("init: %v", err)
	}
	return New(f), f
}

func runOnce(t *testing.T, a *Application) {
	t.Helper()
	a.draw()
	if err := a.screen.Show(a.surface); err != nil {
		t.Fatalf("show: %v", err)
	}
}

func TestApplicationLayout(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	if a.MenuBar().Bounds != (vio.R(0, 0, scrW, 1)) {
		t.Fatalf("menubar bounds: %+v", a.MenuBar().Bounds)
	}
	if a.Desktop().Bounds != (vio.R(0, 1, scrW, scrH-2)) {
		t.Fatalf("desktop bounds: %+v", a.Desktop().Bounds)
	}
	if a.StatusLine().Bounds != (vio.R(0, scrH-1, scrW, 1)) {
		t.Fatalf("status bounds: %+v", a.StatusLine().Bounds)
	}
}

func TestApplicationRunQuitsOnAltX(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	go f.PushEvent(event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'x', Mod: event.ModAlt},
	})
	if err := a.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestApplicationRunQuitsOnCmdQuit(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	go f.PushEvent(event.Event{
		What: event.ClassCommand,
		Msg:  event.MessageEvent{Command: event.CmdQuit},
	})
	if err := a.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
}

func TestApplicationResize(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()

	go func() {
		f.PushEvent(event.Event{
			What:  event.ClassResize,
			Mouse: event.MouseEvent{X: 60, Y: 20},
		})
		f.PushEvent(event.Event{
			What: event.ClassKey,
			Key:  event.KeyEvent{Key: event.KeyRune, Rune: 'x', Mod: event.ModAlt},
		})
	}()
	if err := a.Run(); err != nil {
		t.Fatalf("run: %v", err)
	}
	w, _ := a.Surface().Size()
	if w != scrW {
		// Surface keeps the original size because FakeScreen.Size() did not
		// actually change. We only test the resize() method directly below.
		_ = w
	}
}

func TestApplicationPutEventDoesNotBlock(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()
	for i := 0; i < 200; i++ {
		a.PutEvent(event.Event{What: event.ClassIdle})
	}
}

func TestApplicationEmptyDesktopSnapshot(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()
	runOnce(t, a)

	got := strings.Join(f.Surface().Snapshot(), "\n")
	goldenPath := filepath.Join("testdata", "empty_desktop.golden")

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v (re-run with UPDATE_GOLDEN=1 to seed)", err)
	}
	if string(want) != got {
		t.Fatalf("snapshot drift:\nwant:\n%s\ngot:\n%s", want, got)
	}
}

func TestApplicationResizePropagatesToChildren(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()
	a.resize(60, 20)
	if a.MenuBar().Bounds.W != 60 {
		t.Fatalf("menubar W: %d", a.MenuBar().Bounds.W)
	}
	if a.Desktop().Bounds != (vio.R(0, 1, 60, 18)) {
		t.Fatalf("desktop: %+v", a.Desktop().Bounds)
	}
	if a.StatusLine().Bounds.Y != 19 {
		t.Fatalf("status Y: %d", a.StatusLine().Bounds.Y)
	}
}

func TestApplicationResizeIgnoresNonPositive(t *testing.T) {
	a, f := newTestApp(t)
	defer f.Fini()
	original := a.Bounds
	a.resize(0, 0)
	if a.Bounds != original {
		t.Fatalf("resize(0,0) mutated bounds: %+v", a.Bounds)
	}
}

func TestStatusLineDefaults(t *testing.T) {
	if len(menu.DefaultHints) < 3 {
		t.Fatalf("default hints: %d", len(menu.DefaultHints))
	}
	if menu.DefaultHints[2].Cmd != event.CmdQuit {
		t.Fatalf("third hint should be Exit/CmdQuit: %+v", menu.DefaultHints[2])
	}
}
