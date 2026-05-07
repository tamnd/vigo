package widget

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Orientation distinguishes a horizontal scrollbar (W > H) from a
// vertical one (H >= W). Callers don't construct this directly; the
// ScrollBar constructor picks the orientation from the bounds.
type Orientation uint8

// Orientation values.
const (
	Horizontal Orientation = iota
	Vertical
)

// ScrollBar is the strip with arrow caps at each end and a thumb that
// represents the current value inside the [Min, Max] range. PgStep is
// the increment when the user clicks the page area; ArStep is the
// increment for the arrow caps. Whenever Value changes, ScrollBar
// broadcasts a CmdScrollBarChanged command via its owner.
type ScrollBar struct {
	*view.View

	Min, Max int
	Value    int
	PgStep   int
	ArStep   int

	orient Orientation
}

// NewScrollBar returns a scrollbar sized to bounds. Orientation is
// horizontal when bounds.W > bounds.H. Step defaults are 1 for the
// arrow caps and max(1, length-2) for the page region.
func NewScrollBar(bounds vio.Rect) *ScrollBar {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.PaletteIndex = SlotScrollPage
	v.EventMask = event.ClassMouseDown

	orient := Vertical
	if bounds.W > bounds.H {
		orient = Horizontal
	}
	return &ScrollBar{
		View:   v,
		ArStep: 1,
		PgStep: max(1, sbLength(bounds, orient)-2),
		orient: orient,
	}
}

// Orientation returns whether the bar is horizontal or vertical.
func (s *ScrollBar) Orientation() Orientation { return s.orient }

// SetParams updates value, min, max, page step, and arrow step in one
// call and clamps Value to the new range. SetParams broadcasts a
// CmdScrollBarChanged if the clamped Value differs from the prior one.
func (s *ScrollBar) SetParams(value, minVal, maxVal, pgStep, arStep int) {
	if minVal > maxVal {
		minVal, maxVal = maxVal, minVal
	}
	s.Min = minVal
	s.Max = maxVal
	s.PgStep = max(1, pgStep)
	s.ArStep = max(1, arStep)
	s.SetValue(value)
}

// SetValue clamps v to [Min, Max] and stores it. If the clamped value
// differs from the previous one, it broadcasts CmdScrollBarChanged.
func (s *ScrollBar) SetValue(v int) {
	v = clampInt(v, s.Min, s.Max)
	if v == s.Value {
		return
	}
	s.Value = v
	s.notify()
}

// thumbPos returns the on-screen position of the thumb, relative to
// the start of the page region (i.e. position 0 corresponds to the
// first cell after the leading arrow cap). Returns 0 when the range
// is empty.
func (s *ScrollBar) thumbPos() int {
	pageLen := s.pageLen()
	if pageLen <= 0 {
		return 0
	}
	span := s.Max - s.Min
	if span <= 0 {
		return 0
	}
	pos := (s.Value - s.Min) * (pageLen - 1) / span
	return clampInt(pos, 0, pageLen-1)
}

func (s *ScrollBar) pageLen() int {
	return sbLength(s.Bounds, s.orient) - 2
}

// Draw paints the arrows, the page background, and the thumb.
func (s *ScrollBar) Draw(surf *vio.Surface) {
	page := s.MapColor(SlotScrollPage)
	arrow := s.MapColor(SlotScrollArrow)

	pageLen := max(s.pageLen(), 0)
	thumb := s.thumbPos()

	if s.orient == Vertical {
		x := s.Bounds.X
		y := s.Bounds.Y
		surf.Set(x, y, '▲', arrow)
		for i := range pageLen {
			surf.Set(x, y+1+i, '▒', page)
		}
		if pageLen > 0 {
			surf.Set(x, y+1+thumb, '▓', page)
		}
		surf.Set(x, s.Bounds.Bottom()-1, '▼', arrow)
		return
	}

	x := s.Bounds.X
	y := s.Bounds.Y
	surf.Set(x, y, '◄', arrow)
	for i := range pageLen {
		surf.Set(x+1+i, y, '▒', page)
	}
	if pageLen > 0 {
		surf.Set(x+1+thumb, y, '▓', page)
	}
	surf.Set(s.Bounds.Right()-1, y, '►', arrow)
}

// HandleEvent maps mouse-down on the arrow caps to ArStep increments
// and clicks on the page region to PgStep jumps. Releases and drag
// are not implemented in v0.2 phase 4; they ride along in a later
// polish pass.
func (s *ScrollBar) HandleEvent(e *event.Event) {
	if e.What != event.ClassMouseDown {
		return
	}
	p := vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}
	if !s.Bounds.Contains(p) {
		return
	}

	pos, length := s.localPos(p)
	switch pos {
	case 0:
		s.SetValue(s.Value - s.ArStep)
	case length - 1:
		s.SetValue(s.Value + s.ArStep)
	default:
		thumb := s.thumbPos() + 1 // +1 for leading arrow cap
		switch {
		case pos < thumb:
			s.SetValue(s.Value - s.PgStep)
		case pos > thumb:
			s.SetValue(s.Value + s.PgStep)
		}
	}
	e.Clear()
}

func (s *ScrollBar) localPos(p vio.Point) (int, int) {
	if s.orient == Vertical {
		return p.Y - s.Bounds.Y, s.Bounds.H
	}
	return p.X - s.Bounds.X, s.Bounds.W
}

func (s *ScrollBar) notify() {
	if s.Owner == nil {
		return
	}
	ev := &event.Event{
		What: event.ClassBroadcast,
		Msg: event.MessageEvent{
			Command: event.CmdScrollBarChanged,
			Info:    s,
		},
	}
	s.Owner.HandleEvent(ev)
}

func sbLength(b vio.Rect, o Orientation) int {
	if o == Vertical {
		return b.H
	}
	return b.W
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
