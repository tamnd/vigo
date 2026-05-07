package widget

import (
	"fmt"

	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// ParamText is a read-only text view whose payload is rendered through
// fmt.Sprintf. Use it for status fields with live values: "Line %d, Col
// %d" and so on. The format string is fixed at construction; SetParams
// updates the slice of values, redraw paints the formatted result.
type ParamText struct {
	*view.View

	Format string
	Params []any
}

// NewParamText returns a ParamText sized to bounds with the given
// format string. The initial parameter list is empty.
func NewParamText(bounds vio.Rect, format string) *ParamText {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.EventMask = 0
	v.PaletteIndex = SlotDialogText
	return &ParamText{View: v, Format: format}
}

// SetParams replaces the parameter slice. The slice is shared with the
// view; callers must not mutate concurrently with redraws.
func (p *ParamText) SetParams(params ...any) {
	p.Params = params
}

// Text returns the formatted string. It is the value Draw paints.
func (p *ParamText) Text() string {
	if p.Format == "" {
		return ""
	}
	return fmt.Sprintf(p.Format, p.Params...)
}

// Draw paints the formatted text. Lines longer than the bounds width
// are clipped; only the first row is used.
func (p *ParamText) Draw(s *vio.Surface) {
	attr := p.MapColor(p.PaletteIndex)
	s.FillRect(p.Bounds, ' ', attr)
	if p.Bounds.H <= 0 || p.Bounds.W <= 0 {
		return
	}
	s.DrawString(p.Bounds.X, p.Bounds.Y, clipLine(p.Text(), p.Bounds.W), attr)
}

var _ view.Viewer = (*ParamText)(nil)
