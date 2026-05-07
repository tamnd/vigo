package widget

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/vio"
)

func TestStaticTextDrawsSingleLine(t *testing.T) {
	host := newPalettedHost()
	st := NewStaticText(vio.R(0, 0, 10, 1), "hello")
	host.Insert(st)
	s := vio.NewSurface(10, 1)
	st.Draw(s)

	row := s.Snapshot()[0]
	if !strings.HasPrefix(row, "hello") {
		t.Fatalf("expected text at start of row, got %q", row)
	}
}

func TestStaticTextSplitsOnNewline(t *testing.T) {
	host := newPalettedHost()
	st := NewStaticText(vio.R(0, 0, 10, 3), "first\nsecond")
	host.Insert(st)
	s := vio.NewSurface(10, 3)
	st.Draw(s)

	rows := s.Snapshot()
	if !strings.HasPrefix(rows[0], "first") {
		t.Fatalf("row 0: %q", rows[0])
	}
	if !strings.HasPrefix(rows[1], "second") {
		t.Fatalf("row 1: %q", rows[1])
	}
}

func TestStaticTextClipsLongLine(t *testing.T) {
	host := newPalettedHost()
	st := NewStaticText(vio.R(0, 0, 4, 1), "abcdefgh")
	host.Insert(st)
	s := vio.NewSurface(4, 1)
	st.Draw(s)

	row := s.Snapshot()[0]
	if row != "abcd" {
		t.Fatalf("expected clipped row %q, got %q", "abcd", row)
	}
}

func TestStaticTextEmptyStaysBlank(t *testing.T) {
	host := newPalettedHost()
	st := NewStaticText(vio.R(0, 0, 4, 1), "")
	host.Insert(st)
	s := vio.NewSurface(4, 1)
	st.Draw(s)

	row := s.Snapshot()[0]
	if strings.TrimSpace(row) != "" {
		t.Fatalf("empty static text should leave blanks, got %q", row)
	}
}
