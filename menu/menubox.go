package menu

import (
	"unicode"

	"github.com/tamnd/vigo/cmd"
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// MenuBox is the vertical pull-down that drops below a Bar item or
// nests under another MenuBox. It paints a single-line border, lays
// out items top to bottom, and exposes Result so a host can drive it
// like a modal sub-loop.
//
// Disabled rows are dimmed, separators paint as a horizontal rule,
// submenu rows show a "▶" cap. Arrow keys move focus skipping
// separators; Enter fires the row's command (or sets ChildOpen on a
// submenu); Esc closes; a hot-key letter jumps to the matching row
// and fires it.
//
//nolint:revive // MenuBox matches the spec name; menu.Box would conflict with the conceptual model
type MenuBox struct {
	*view.View

	Items   []Item
	Enabler *cmd.Enabler

	sel       int
	result    event.CommandID
	chosenCmd event.CommandID
	openChild bool
}

// NewMenuBox returns a pull-down anchored at (x, y) sized to fit
// items. The width is the longest title plus shortcut plus padding;
// the height is len(items) plus 2 for the border.
func NewMenuBox(x, y int, items []Item, enabler *cmd.Enabler) *MenuBox {
	w, h := menuBoxSize(items)
	v := view.NewView(vio.R(x, y, w, h))
	v.GrowMode = view.GrowFixed
	v.Options |= view.OptSelectable
	v.PaletteIndex = 2
	v.EventMask = event.ClassKey | event.ClassMouseDown

	m := &MenuBox{
		View:    v,
		Items:   items,
		Enabler: enabler,
	}
	m.sel = m.firstSelectable()
	return m
}

// Result returns CmdNone while the menu is still open, CmdCancel when
// the user closed it without picking, or CmdMenu when a row fired a
// command (which is then available via Chosen).
func (m *MenuBox) Result() event.CommandID { return m.result }

// Chosen returns the command attached to the activated row, or
// CmdNone if the user canceled.
func (m *MenuBox) Chosen() event.CommandID { return m.chosenCmd }

// Selected returns the focused item.
func (m *MenuBox) Selected() Item {
	if m.sel < 0 || m.sel >= len(m.Items) {
		return Item{}
	}
	return m.Items[m.sel]
}

// SelectedIndex returns the focused row index.
func (m *MenuBox) SelectedIndex() int { return m.sel }

// SetSelectedIndex moves focus to i, skipping separators and disabled
// rows in the supplied direction (1 forward, -1 backward, 0 forward).
func (m *MenuBox) SetSelectedIndex(i int) {
	if len(m.Items) == 0 {
		m.sel = -1
		return
	}
	if i < 0 || i >= len(m.Items) {
		return
	}
	m.sel = i
}

// ChildOpen reports whether the focused row is a submenu trigger that
// the user just activated. The Bar/host clears this after spawning a
// nested MenuBox.
func (m *MenuBox) ChildOpen() bool { return m.openChild }

// AckChild clears the ChildOpen flag.
func (m *MenuBox) AckChild() { m.openChild = false }

// IsEnabled reports whether row i can be activated. Separators,
// out-of-range indices, and rows whose Cmd is in a disabled state
// return false.
func (m *MenuBox) IsEnabled(i int) bool {
	if i < 0 || i >= len(m.Items) {
		return false
	}
	it := m.Items[i]
	if it.Sep {
		return false
	}
	if m.Enabler == nil || it.Cmd == event.CmdNone {
		return true
	}
	return m.Enabler.IsEnabled(it.Cmd)
}

// Draw paints the box.
func (m *MenuBox) Draw(s *vio.Surface) {
	normal := m.MapColor(2)
	disabled := m.MapColor(3)
	selected := m.MapColor(4)
	shortcut := m.MapColor(5)

	s.FillRect(m.Bounds, ' ', normal)
	m.drawBorder(s, normal)

	innerW := m.Bounds.W - 2
	for i, it := range m.Items {
		row := m.Bounds.Y + 1 + i
		if it.Sep {
			s.Set(m.Bounds.X, row, '├', normal)
			for x := range innerW {
				s.Set(m.Bounds.X+1+x, row, '─', normal)
			}
			s.Set(m.Bounds.Right()-1, row, '┤', normal)
			continue
		}
		attr := normal
		if !m.IsEnabled(i) {
			attr = disabled
		}
		if i == m.sel {
			attr = selected
		}
		// fill row interior
		for x := range innerW {
			s.Set(m.Bounds.X+1+x, row, ' ', attr)
		}
		// Title with hotkey
		x := m.Bounds.X + 2
		hotAttr := shortcut
		if i == m.sel {
			hotAttr = selected
		}
		for j, r := range []rune(it.Title) {
			ch := attr
			if j == it.Hotkey && m.IsEnabled(i) {
				ch = hotAttr
			}
			s.Set(x, row, r, ch)
			x++
		}
		// Right-side: shortcut or submenu arrow
		right := m.Bounds.Right() - 2
		if len(it.Children) > 0 {
			s.Set(right, row, '▶', attr)
			continue
		}
		if it.Shortcut == "" {
			continue
		}
		runes := []rune(it.Shortcut)
		startX := right - len(runes) + 1
		for j, r := range runes {
			ch := shortcut
			if i == m.sel || !m.IsEnabled(i) {
				ch = attr
			}
			s.Set(startX+j, row, r, ch)
		}
	}
}

func (m *MenuBox) drawBorder(s *vio.Surface, attr vio.Attr) {
	x0, y0 := m.Bounds.X, m.Bounds.Y
	x1, y1 := m.Bounds.Right()-1, m.Bounds.Bottom()-1
	s.Set(x0, y0, '┌', attr)
	s.Set(x1, y0, '┐', attr)
	s.Set(x0, y1, '└', attr)
	s.Set(x1, y1, '┘', attr)
	for x := x0 + 1; x < x1; x++ {
		s.Set(x, y0, '─', attr)
		s.Set(x, y1, '─', attr)
	}
	for y := y0 + 1; y < y1; y++ {
		s.Set(x0, y, '│', attr)
		s.Set(x1, y, '│', attr)
	}
}

// HandleEvent moves focus and dispatches.
func (m *MenuBox) HandleEvent(e *event.Event) {
	switch e.What {
	case event.ClassKey:
		m.handleKey(e)
	case event.ClassMouseDown:
		m.handleMouse(e)
	}
}

func (m *MenuBox) handleKey(e *event.Event) {
	switch e.Key.Key {
	case event.KeyArrowDown:
		m.advance(1)
		e.Clear()
	case event.KeyArrowUp:
		m.advance(-1)
		e.Clear()
	case event.KeyHome:
		m.sel = m.firstSelectable()
		e.Clear()
	case event.KeyEnd:
		m.sel = m.lastSelectable()
		e.Clear()
	case event.KeyEsc:
		m.result = event.CmdCancel
		e.Clear()
	case event.KeyEnter:
		m.activate(m.sel)
		e.Clear()
	case event.KeyRune:
		if idx := m.findHotkey(e.Key.Rune); idx >= 0 {
			m.sel = idx
			m.activate(idx)
			e.Clear()
		}
	}
}

func (m *MenuBox) handleMouse(e *event.Event) {
	p := vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}
	if !m.Bounds.Contains(p) {
		m.result = event.CmdCancel
		e.Clear()
		return
	}
	row := p.Y - m.Bounds.Y - 1
	if row < 0 || row >= len(m.Items) {
		return
	}
	if m.Items[row].Sep {
		return
	}
	m.sel = row
	m.activate(row)
	e.Clear()
}

func (m *MenuBox) activate(i int) {
	if !m.IsEnabled(i) {
		return
	}
	it := m.Items[i]
	if len(it.Children) > 0 {
		m.openChild = true
		return
	}
	m.chosenCmd = it.Cmd
	m.result = event.CmdMenu
}

func (m *MenuBox) advance(step int) {
	n := len(m.Items)
	if n == 0 {
		return
	}
	for k := 1; k <= n; k++ {
		idx := (m.sel + step*k) % n
		if idx < 0 {
			idx += n
		}
		if !m.Items[idx].Sep {
			m.sel = idx
			return
		}
	}
}

func (m *MenuBox) firstSelectable() int {
	for i, it := range m.Items {
		if !it.Sep {
			return i
		}
	}
	return -1
}

func (m *MenuBox) lastSelectable() int {
	for i := len(m.Items) - 1; i >= 0; i-- {
		if !m.Items[i].Sep {
			return i
		}
	}
	return -1
}

func (m *MenuBox) findHotkey(r rune) int {
	target := unicode.ToLower(r)
	for i, it := range m.Items {
		if it.Sep || it.Hotkey < 0 {
			continue
		}
		runes := []rune(it.Title)
		if it.Hotkey >= len(runes) {
			continue
		}
		if unicode.ToLower(runes[it.Hotkey]) == target {
			return i
		}
	}
	return -1
}

func menuBoxSize(items []Item) (int, int) {
	maxTitle := 0
	maxShortcut := 0
	for _, it := range items {
		if it.Sep {
			continue
		}
		titleLen := len([]rune(it.Title))
		if len(it.Children) > 0 || it.Shortcut != "" {
			titleLen += 2 // gutter before shortcut/arrow
		}
		maxTitle = max(maxTitle, titleLen)
		shortcutLen := len([]rune(it.Shortcut))
		if len(it.Children) > 0 {
			shortcutLen = max(shortcutLen, 1)
		}
		maxShortcut = max(maxShortcut, shortcutLen)
	}
	w := max(4+maxTitle+maxShortcut, 8) // 2 padding + content + 2 border
	h := max(len(items)+2, 3)
	return w, h
}
