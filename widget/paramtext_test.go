package widget

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/vio"
)

func TestParamTextFormatsParams(t *testing.T) {
	host := newPalettedHost()
	p := NewParamText(vio.R(0, 0, 20, 1), "Line %d, Col %d")
	host.Insert(p)
	p.SetParams(12, 34)

	if got, want := p.Text(), "Line 12, Col 34"; got != want {
		t.Fatalf("Text=%q want %q", got, want)
	}
}

func TestParamTextSetParamsReplacesSlice(t *testing.T) {
	host := newPalettedHost()
	p := NewParamText(vio.R(0, 0, 20, 1), "%s=%d")
	host.Insert(p)

	p.SetParams("ans", 42)
	if got := p.Text(); got != "ans=42" {
		t.Fatalf("first format: %q", got)
	}

	p.SetParams("pi", 3)
	if got := p.Text(); got != "pi=3" {
		t.Fatalf("second format: %q", got)
	}
}

func TestParamTextDrawClipsToWidth(t *testing.T) {
	host := newPalettedHost()
	p := NewParamText(vio.R(0, 0, 5, 1), "abcdefgh")
	host.Insert(p)

	s := vio.NewSurface(5, 1)
	p.Draw(s)
	if got := s.Snapshot()[0]; !strings.HasPrefix(got, "abcde") {
		t.Fatalf("clipped row: %q", got)
	}
}

func TestParamTextEmptyFormatBlankRow(t *testing.T) {
	host := newPalettedHost()
	p := NewParamText(vio.R(0, 0, 8, 1), "")
	host.Insert(p)

	s := vio.NewSurface(8, 1)
	p.Draw(s)
	if got := s.Snapshot()[0]; got != "        " {
		t.Fatalf("empty format should leave blanks: %q", got)
	}
}
