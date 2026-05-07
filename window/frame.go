// Package window implements the framed Window view that hosts content on
// the desktop. v0.2 phase 1 delivers the minimum: a Frame border (title,
// optional close and zoom glyphs) and a Window that owns the frame plus
// a client rectangle. Drag, resize, zoom, modal dispatch, and the close
// button arrive in phase 2.
package window

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Flags controls which decorations the frame draws and which actions a
// Window will accept once phase 2 lands. Mirrors the flags Turbo Vision
// passes to Window.Init: wfMove, wfGrow, wfClose, wfZoom.
type Flags uint8

// Flag bits.
const (
	FlagMove Flags = 1 << iota
	FlagGrow
	FlagClose
	FlagZoom
)

// FlagsAll enables every decoration. Stock IDE windows use this set.
const FlagsAll = FlagMove | FlagGrow | FlagClose | FlagZoom

// Frame is the boxed border drawn around a Window. The frame paints a
// title centered along the top edge, plus optional close and zoom
// glyphs at the corners. Active windows use a different palette slot
// for the line glyphs so the focused window stands out.
//
// A Frame's bounds are the same rectangle as its owning Window; the
// frame paints only the perimeter and leaves the interior untouched
// for the window's content children.
//
// PassiveSlot and ActiveSlot let a Dialog point the frame at a
// different palette pair (e.g. 15/16 for the gray-on-gray dialog
// frame) without touching the regular window theme.
type Frame struct {
	*view.View

	Title       string
	flags       Flags
	PassiveSlot byte
	ActiveSlot  byte
}

// NewFrame returns a frame sized to bounds, with the given decoration
// flags. Frames are not selectable on their own; the parent Window owns
// focus.
func NewFrame(bounds vio.Rect, flags Flags) *Frame {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowHiX | view.GrowHiY
	v.EventMask = 0 // frames do not consume events directly
	return &Frame{
		View:        v,
		flags:       flags,
		PassiveSlot: 10,
		ActiveSlot:  11,
	}
}

// Flags returns the frame's decoration flags.
func (f *Frame) Flags() Flags { return f.flags }

// Draw paints the box, title, and corner glyphs.
func (f *Frame) Draw(s *vio.Surface) {
	active := f.HasState(view.StateActive)
	slot := f.PassiveSlot
	if active {
		slot = f.ActiveSlot
	}
	attr := f.MapColor(slot)
	border := vio.Double
	if !active {
		border = vio.Single
	}
	s.Box(f.Bounds, border, attr)

	if f.flags&FlagClose != 0 && f.Bounds.W >= 5 {
		f.drawClose(s, attr)
	}
	if f.flags&FlagZoom != 0 && f.Bounds.W >= 7 {
		f.drawZoom(s, attr)
	}
	if f.Title != "" {
		f.drawTitle(s, f.Title, attr)
	}
}

func (f *Frame) drawClose(s *vio.Surface, attr vio.Attr) {
	x := f.Bounds.X + 2
	y := f.Bounds.Y
	s.Set(x-1, y, '[', attr)
	s.Set(x, y, '■', attr)
	s.Set(x+1, y, ']', attr)
}

func (f *Frame) drawZoom(s *vio.Surface, attr vio.Attr) {
	x := f.Bounds.Right() - 3
	y := f.Bounds.Y
	glyph := '↕'
	if f.HasState(view.StateModal) {
		glyph = '↑'
	}
	s.Set(x-1, y, '[', attr)
	s.Set(x, y, glyph, attr)
	s.Set(x+1, y, ']', attr)
}

func (f *Frame) drawTitle(s *vio.Surface, title string, attr vio.Attr) {
	runes := []rune(title)
	maxLen := f.Bounds.W - 6
	if maxLen <= 0 {
		return
	}
	if len(runes) > maxLen {
		runes = runes[:maxLen]
	}
	x := f.Bounds.X + (f.Bounds.W-len(runes))/2
	y := f.Bounds.Y
	leftPad := x - 1
	rightPad := x + len(runes)
	closeRight := -1
	if f.flags&FlagClose != 0 && f.Bounds.W >= 5 {
		closeRight = f.Bounds.X + 3 // last column of "[■]"
	}
	zoomLeft := -1
	if f.flags&FlagZoom != 0 && f.Bounds.W >= 7 {
		zoomLeft = f.Bounds.Right() - 4 // first column of "[↕]"
	}
	if leftPad > closeRight {
		s.Set(leftPad, y, ' ', attr)
	}
	for i, r := range runes {
		s.Set(x+i, y, r, attr)
	}
	if zoomLeft < 0 || rightPad < zoomLeft {
		s.Set(rightPad, y, ' ', attr)
	}
}

// HandleEvent ignores everything; frames are decoration only in v0.2.1.
func (f *Frame) HandleEvent(e *event.Event) { _ = e }
