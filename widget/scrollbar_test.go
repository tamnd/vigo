package widget

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
)

func TestScrollBarOrientation(t *testing.T) {
	v := NewScrollBar(vio.R(0, 0, 1, 10))
	if v.Orientation() != Vertical {
		t.Fatal("tall bar should be vertical")
	}
	h := NewScrollBar(vio.R(0, 0, 10, 1))
	if h.Orientation() != Horizontal {
		t.Fatal("wide bar should be horizontal")
	}
}

func TestScrollBarSetParamsClampsValue(t *testing.T) {
	s := NewScrollBar(vio.R(0, 0, 1, 10))
	s.SetParams(50, 0, 10, 4, 1)
	if s.Value != 10 {
		t.Fatalf("value should clamp to Max, got %d", s.Value)
	}
	s.SetParams(-5, 0, 10, 4, 1)
	if s.Value != 0 {
		t.Fatalf("value should clamp to Min, got %d", s.Value)
	}
}

func TestScrollBarThumbAtExtremes(t *testing.T) {
	s := NewScrollBar(vio.R(0, 0, 1, 10))
	s.SetParams(0, 0, 100, 1, 1)

	if got := s.thumbPos(); got != 0 {
		t.Fatalf("thumb at Min should be 0, got %d", got)
	}
	s.SetValue(100)
	pageLen := s.pageLen()
	if got := s.thumbPos(); got != pageLen-1 {
		t.Fatalf("thumb at Max should be pageLen-1=%d, got %d", pageLen-1, got)
	}
}

func TestScrollBarBroadcastsOnChange(t *testing.T) {
	host := newPalettedHost()
	spy := newCommandSpy()
	host.Insert(spy)
	s := NewScrollBar(vio.R(0, 0, 1, 10))
	host.Insert(s)
	s.SetParams(0, 0, 10, 4, 1)

	spy.lastCommand = event.CmdNone
	s.SetValue(5)
	if spy.lastClass != event.ClassBroadcast {
		t.Fatalf("expected Broadcast, got %v", spy.lastClass)
	}
	if spy.lastCommand != event.CmdScrollBarChanged {
		t.Fatalf("expected CmdScrollBarChanged, got %v", spy.lastCommand)
	}
}

func TestScrollBarSameValueDoesNotNotify(t *testing.T) {
	host := newPalettedHost()
	spy := newCommandSpy()
	host.Insert(spy)
	s := NewScrollBar(vio.R(0, 0, 1, 10))
	host.Insert(s)
	s.SetParams(5, 0, 10, 4, 1)
	spy.lastCommand = event.CmdNone

	s.SetValue(5)
	if spy.lastCommand != event.CmdNone {
		t.Fatal("setting the same value should not broadcast")
	}
}

func TestScrollBarMouseOnArrowsSteps(t *testing.T) {
	host := newPalettedHost()
	s := NewScrollBar(vio.R(0, 0, 1, 10))
	host.Insert(s)
	s.SetParams(5, 0, 10, 4, 1)

	s.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 0, Y: 0}, // top arrow
	})
	if s.Value != 4 {
		t.Fatalf("top arrow should decrement, got %d", s.Value)
	}

	s.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 0, Y: 9}, // bottom arrow
	})
	if s.Value != 5 {
		t.Fatalf("bottom arrow should increment, got %d", s.Value)
	}
}

func TestScrollBarMouseOnPagePagesByPgStep(t *testing.T) {
	host := newPalettedHost()
	s := NewScrollBar(vio.R(0, 0, 1, 10))
	host.Insert(s)
	s.SetParams(5, 0, 100, 10, 1)

	s.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 0, Y: 8}, // below thumb (thumb at value 5 of 100)
	})
	if s.Value != 15 {
		t.Fatalf("page-down should add PgStep, got %d", s.Value)
	}
}

func TestScrollBarMouseOutsideIgnored(t *testing.T) {
	host := newPalettedHost()
	s := NewScrollBar(vio.R(0, 0, 1, 10))
	host.Insert(s)
	s.SetParams(5, 0, 10, 4, 1)
	before := s.Value

	s.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 5, Y: 5},
	})
	if s.Value != before {
		t.Fatalf("click outside should not change value")
	}
}

func TestScrollBarDrawVerticalGlyphs(t *testing.T) {
	host := newPalettedHost()
	s := NewScrollBar(vio.R(0, 0, 1, 5))
	host.Insert(s)
	s.SetParams(0, 0, 10, 1, 1)

	surf := vio.NewSurface(1, 5)
	s.Draw(surf)
	rows := surf.Snapshot()
	if rows[0] != "▲" || rows[len(rows)-1] != "▼" {
		t.Fatalf("expected arrow caps, got top=%q bot=%q", rows[0], rows[len(rows)-1])
	}
	body := strings.Join(rows[1:len(rows)-1], "")
	if !strings.ContainsRune(body, '▓') {
		t.Fatalf("expected thumb glyph in body, got %q", body)
	}
}

func TestScrollBarDrawHorizontalGlyphs(t *testing.T) {
	host := newPalettedHost()
	s := NewScrollBar(vio.R(0, 0, 5, 1))
	host.Insert(s)
	s.SetParams(0, 0, 10, 1, 1)

	surf := vio.NewSurface(5, 1)
	s.Draw(surf)
	row := surf.Snapshot()[0]
	if !strings.HasPrefix(row, "◄") {
		t.Fatalf("missing left arrow: %q", row)
	}
	if !strings.HasSuffix(row, "►") {
		t.Fatalf("missing right arrow: %q", row)
	}
	if !strings.ContainsRune(row, '▓') {
		t.Fatalf("missing thumb: %q", row)
	}
}
