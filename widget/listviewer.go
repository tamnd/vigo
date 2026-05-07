package widget

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// GetTextFunc returns the rendered string for item index when at most
// maxLen runes are wanted. It is the abstract entry point ListViewer
// uses to fetch row text; ListBox installs one that reads from a
// string slice. Callers can install any function (e.g. backed by a
// log buffer) without subtyping.
type GetTextFunc func(item, maxLen int) string

// ListViewer is the scrolling-list base. It owns the focused item,
// the top-visible row, optional scrollbars, and a GetText function
// that yields the row text. The view paints rows top-down, highlights
// the focused row when StateFocused is set, and broadcasts no
// commands; concrete lists wire commands themselves.
type ListViewer struct {
	*view.View

	Range   int
	Focused int
	HScroll *ScrollBar
	VScroll *ScrollBar

	GetText GetTextFunc

	topItem int
}

// NewListViewer returns a list viewer sized to bounds. The caller
// must set GetText (and Range) before the first Draw, otherwise rows
// render blank.
func NewListViewer(bounds vio.Rect, vScroll, hScroll *ScrollBar) *ListViewer {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.Options |= view.OptSelectable
	v.PaletteIndex = SlotListNormal
	v.EventMask = event.ClassKey | event.ClassMouseDown | event.ClassBroadcast
	return &ListViewer{
		View:    v,
		HScroll: hScroll,
		VScroll: vScroll,
	}
}

// SetRange updates Range and clamps Focused and topItem to fit. It
// also reconfigures the vertical scrollbar if one is attached.
func (l *ListViewer) SetRange(n int) {
	if n < 0 {
		n = 0
	}
	l.Range = n
	if l.Focused >= n {
		l.Focused = n - 1
	}
	if l.Focused < 0 {
		l.Focused = 0
	}
	l.adjustView()
	l.syncScrollbar()
}

// FocusItem moves the focused row to i (clamped) and pages the view
// if i would be off-screen.
func (l *ListViewer) FocusItem(i int) {
	if l.Range == 0 {
		l.Focused = 0
		return
	}
	l.Focused = clampInt(i, 0, l.Range-1)
	l.adjustView()
	l.syncScrollbar()
}

// TopItem returns the first visible row index.
func (l *ListViewer) TopItem() int { return l.topItem }

func (l *ListViewer) adjustView() {
	if l.Bounds.H <= 0 {
		l.topItem = 0
		return
	}
	if l.Focused < l.topItem {
		l.topItem = l.Focused
	}
	if l.Focused >= l.topItem+l.Bounds.H {
		l.topItem = l.Focused - l.Bounds.H + 1
	}
	if l.topItem < 0 {
		l.topItem = 0
	}
}

func (l *ListViewer) syncScrollbar() {
	if l.VScroll == nil {
		return
	}
	maxVal := max(l.Range-1, 0)
	pgStep := max(l.Bounds.H-1, 1)
	l.VScroll.SetParams(l.Focused, 0, maxVal, pgStep, 1)
}

// Draw paints all visible rows. Rows beyond Range are blanked. The
// focused row is highlighted when the view holds StateFocused.
func (l *ListViewer) Draw(s *vio.Surface) {
	body := l.MapColor(SlotListNormal)
	focused := l.MapColor(SlotListFocused)
	selected := l.MapColor(SlotListSelected)

	s.FillRect(l.Bounds, ' ', body)
	if l.Bounds.H <= 0 || l.Bounds.W <= 0 {
		return
	}

	for row := range l.Bounds.H {
		idx := l.topItem + row
		if idx >= l.Range {
			break
		}
		attr := body
		switch {
		case idx == l.Focused && l.HasState(view.StateFocused):
			attr = focused
		case idx == l.Focused:
			attr = selected
		}
		s.FillRect(vio.R(l.Bounds.X, l.Bounds.Y+row, l.Bounds.W, 1), ' ', attr)
		if l.GetText == nil {
			continue
		}
		text := l.GetText(idx, l.Bounds.W)
		x := l.Bounds.X
		for _, r := range text {
			if x >= l.Bounds.Right() {
				break
			}
			s.Set(x, l.Bounds.Y+row, r, attr)
			x++
		}
	}
}

// HandleEvent moves Focused with the arrow keys, pages with PgUp/Dn,
// jumps with Home/End, and selects on mouse-down inside the bounds.
// It also listens for CmdScrollBarChanged broadcasts from its own
// vertical scrollbar and snaps Focused to the new value.
func (l *ListViewer) HandleEvent(e *event.Event) {
	switch e.What {
	case event.ClassKey:
		l.handleKey(e)
	case event.ClassMouseDown:
		l.handleMouse(e)
	case event.ClassBroadcast:
		if e.Msg.Command == event.CmdScrollBarChanged && e.Msg.Info == l.VScroll {
			l.FocusItem(l.VScroll.Value)
			e.Clear()
		}
	}
}

func (l *ListViewer) handleKey(e *event.Event) {
	if !l.HasState(view.StateFocused) {
		return
	}
	switch e.Key.Key {
	case event.KeyArrowUp:
		l.FocusItem(l.Focused - 1)
		e.Clear()
	case event.KeyArrowDown:
		l.FocusItem(l.Focused + 1)
		e.Clear()
	case event.KeyPgUp:
		l.FocusItem(l.Focused - max(l.Bounds.H-1, 1))
		e.Clear()
	case event.KeyPgDn:
		l.FocusItem(l.Focused + max(l.Bounds.H-1, 1))
		e.Clear()
	case event.KeyHome:
		l.FocusItem(0)
		e.Clear()
	case event.KeyEnd:
		l.FocusItem(l.Range - 1)
		e.Clear()
	}
}

func (l *ListViewer) handleMouse(e *event.Event) {
	p := vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}
	if !l.Bounds.Contains(p) {
		return
	}
	idx := l.topItem + (p.Y - l.Bounds.Y)
	if idx >= 0 && idx < l.Range {
		l.FocusItem(idx)
		e.Clear()
	}
}
