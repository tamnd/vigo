package help

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func TestRegistryLookupReturnsRegisteredTopic(t *testing.T) {
	r := NewRegistry()
	r.Register(CtxUser+5, Topic{Title: "Editor", Body: "edit stuff"})

	got, ok := r.Lookup(CtxUser + 5)
	if !ok {
		t.Fatal("registered topic not found")
	}
	if got.Title != "Editor" || got.Body != "edit stuff" {
		t.Fatalf("got %+v", got)
	}
}

func TestRegistryLookupUnknownReturnsFalse(t *testing.T) {
	r := NewRegistry()
	if _, ok := r.Lookup(CtxUser + 99); ok {
		t.Fatal("unknown ctx should not return ok")
	}
}

func TestRegistryResolveFallsBackToAbout(t *testing.T) {
	r := Default()
	got := r.Resolve(CtxUser + 42) // unregistered
	if !strings.Contains(got.Body, "vigo") {
		t.Fatalf("fallback should be About, got %+v", got)
	}
}

func TestRegistryResolveZeroWhenEmpty(t *testing.T) {
	r := NewRegistry()
	got := r.Resolve(CtxUser + 1)
	if got.Title != "" || got.Body != "" {
		t.Fatalf("empty registry should resolve to zero Topic, got %+v", got)
	}
}

func TestDefaultSeedsAboutTopic(t *testing.T) {
	r := Default()
	t1, ok := r.Lookup(CtxAbout)
	if !ok {
		t.Fatal("Default should register CtxAbout")
	}
	if t1.Title != "About" {
		t.Fatalf("about title=%q", t1.Title)
	}
	if !strings.Contains(t1.Body, Version) {
		t.Fatalf("about body should mention version, got %q", t1.Body)
	}
}

func TestAboutHasOKDefaultAndCancelButton(t *testing.T) {
	d := About(Default())
	if d.Title != "About" {
		t.Fatalf("dialog title=%q", d.Title)
	}
	if d.DefaultButton() == nil || d.DefaultButton().Command != event.CmdOk {
		t.Fatalf("default button: %+v", d.DefaultButton())
	}
	if d.CancelButton() == nil || d.CancelButton().Command != event.CmdOk {
		t.Fatalf("cancel button should also be OK so Esc closes: %+v", d.CancelButton())
	}
}

func TestAboutNilRegistryFallsBackToDefault(t *testing.T) {
	d := About(nil)
	if d == nil || d.Title != "About" {
		t.Fatalf("nil registry should yield About dialog, got %+v", d)
	}
}

// Confirm About renders without panic and the body contains the
// version string, so we know the StaticText payload is wired through.
func TestAboutRendersBodyText(t *testing.T) {
	host := view.NewGroup(vio.R(0, 0, 80, 24))
	host.Palette = vio.BorlandClassic
	d := About(Default())
	host.Insert(d)
	host.SetCurrent(0)

	s := vio.NewSurface(d.Bounds.W, d.Bounds.H)
	d.Draw(s)

	joined := strings.Join(s.Snapshot(), "\n")
	if !strings.Contains(joined, "vigo") {
		t.Fatalf("About body should mention vigo:\n%s", joined)
	}
}
