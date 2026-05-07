package widget

import (
	"strings"

	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// StaticText is a non-selectable text view that paints its Text into its
// bounds. Newlines start new rows; lines longer than the bounds width are
// clipped (no word-wrap in v0.2 phase 3).
type StaticText struct {
	*view.View

	Text string
}

// NewStaticText returns a StaticText sized to bounds with the given text.
func NewStaticText(bounds vio.Rect, text string) *StaticText {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.EventMask = 0 // never receives events
	v.PaletteIndex = SlotDialogText
	return &StaticText{View: v, Text: text}
}

// Draw paints the static text. Empty text leaves the bounds blank.
func (s *StaticText) Draw(surf *vio.Surface) {
	attr := s.MapColor(s.PaletteIndex)
	surf.FillRect(s.Bounds, ' ', attr)
	if s.Text == "" {
		return
	}
	for i, line := range strings.Split(s.Text, "\n") {
		if i >= s.Bounds.H {
			break
		}
		surf.DrawString(s.Bounds.X, s.Bounds.Y+i, clipLine(line, s.Bounds.W), attr)
	}
}

// clipLine returns the prefix of s whose rune count is at most n. It does
// not allocate when no clipping is needed.
func clipLine(s string, n int) string {
	if n <= 0 {
		return ""
	}
	count := 0
	for i := range s {
		if count == n {
			return s[:i]
		}
		count++
	}
	return s
}
