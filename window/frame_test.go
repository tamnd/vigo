package window

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func newPalettedHost() *view.Group {
	g := view.NewGroup(vio.R(0, 0, 80, 24))
	g.Palette = vio.BorlandClassic
	return g
}

func TestFrameDrawsDoubleBorderWhenActive(t *testing.T) {
	host := newPalettedHost()
	f := NewFrame(vio.R(0, 0, 10, 5), 0)
	f.SetState(view.StateActive, true)
	host.Insert(f)
	s := vio.NewSurface(10, 5)
	f.Draw(s)

	row := s.Snapshot()[0]
	if !strings.ContainsRune(row, '╔') || !strings.ContainsRune(row, '╗') {
		t.Fatalf("active frame should use double-line border: %q", row)
	}
}

func TestFrameDrawsSingleBorderWhenInactive(t *testing.T) {
	host := newPalettedHost()
	f := NewFrame(vio.R(0, 0, 10, 5), 0)
	host.Insert(f)
	s := vio.NewSurface(10, 5)
	f.Draw(s)

	row := s.Snapshot()[0]
	if !strings.ContainsRune(row, '┌') || !strings.ContainsRune(row, '┐') {
		t.Fatalf("inactive frame should use single-line border: %q", row)
	}
}

func TestFrameDrawsTitle(t *testing.T) {
	host := newPalettedHost()
	f := NewFrame(vio.R(0, 0, 20, 5), FlagClose|FlagZoom)
	f.Title = "Hello"
	host.Insert(f)
	s := vio.NewSurface(20, 5)
	f.Draw(s)
	row := s.Snapshot()[0]
	if !strings.Contains(row, "Hello") {
		t.Fatalf("title missing from frame: %q", row)
	}
}

func TestFrameDrawsCloseAndZoomGlyphs(t *testing.T) {
	host := newPalettedHost()
	f := NewFrame(vio.R(0, 0, 20, 5), FlagClose|FlagZoom)
	host.Insert(f)
	s := vio.NewSurface(20, 5)
	f.Draw(s)
	row := s.Snapshot()[0]
	if !strings.ContainsRune(row, '■') {
		t.Fatalf("close glyph missing: %q", row)
	}
	if !strings.ContainsRune(row, '↕') {
		t.Fatalf("zoom glyph missing: %q", row)
	}
}

func TestFrameSkipsCloseGlyphWhenTooNarrow(t *testing.T) {
	host := newPalettedHost()
	f := NewFrame(vio.R(0, 0, 4, 3), FlagClose|FlagZoom)
	host.Insert(f)
	s := vio.NewSurface(4, 3)
	f.Draw(s) // must not panic, must not write outside bounds
	for _, line := range s.Snapshot() {
		if strings.ContainsRune(line, '■') || strings.ContainsRune(line, '↕') {
			t.Fatalf("glyphs should be suppressed at narrow widths: %q", line)
		}
	}
}

func TestFrameTooSmallDoesNothing(t *testing.T) {
	host := newPalettedHost()
	f := NewFrame(vio.R(0, 0, 1, 1), FlagsAll)
	host.Insert(f)
	s := vio.NewSurface(1, 1)
	f.Draw(s) // box is too small; surface stays zero
	if c := s.At(0, 0); c.Rune != 0 {
		t.Fatalf("undersized frame should not draw: %+v", c)
	}
}

func TestFrameTitleClampsToWidth(t *testing.T) {
	host := newPalettedHost()
	f := NewFrame(vio.R(0, 0, 12, 3), 0)
	f.Title = "AReallyLongTitleThatDoesNotFit"
	host.Insert(f)
	s := vio.NewSurface(12, 3)
	f.Draw(s)
	row := s.Snapshot()[0]
	if n := utf8.RuneCountInString(row); n != 12 {
		t.Fatalf("row rune width: %d", n)
	}
}

func TestFrameHandleEventIgnored(t *testing.T) {
	f := NewFrame(vio.R(0, 0, 10, 5), 0)
	e := &event.Event{What: event.ClassKey}
	f.HandleEvent(e)
	if e.What != event.ClassKey {
		t.Fatal("frame should not consume events")
	}
}

func TestFrameFlagsAccessor(t *testing.T) {
	f := NewFrame(vio.R(0, 0, 10, 5), FlagClose|FlagZoom)
	if f.Flags() != FlagClose|FlagZoom {
		t.Fatalf("flags: %v", f.Flags())
	}
}
