package widget

import (
	"strings"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Memo is a multi-line text editor. Lines are stored as a slice of
// rune slices; the cursor is a (line, col) pair. v0.2 phase 4 covers
// printable insertion, Backspace/Delete, Enter inserts a newline, Tab
// inserts spaces, arrow/Home/End/PgUp/PgDn navigation, and viewport
// scrolling. Undo, clipboard, word-wrap, and selection ride along in
// later phases.
type Memo struct {
	*view.View

	TabWidth int

	lines   [][]rune
	curLine int
	curCol  int
	topLine int
	leftCol int
}

// NewMemo returns a Memo sized to bounds with one empty line. TabWidth
// defaults to 4; callers can change it before the first edit.
func NewMemo(bounds vio.Rect) *Memo {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.Options |= view.OptSelectable
	v.PaletteIndex = SlotInputNormal
	v.EventMask = event.ClassKey | event.ClassMouseDown
	return &Memo{
		View:     v,
		TabWidth: 4,
		lines:    [][]rune{{}},
	}
}

// String returns the memo content joined with newlines.
func (m *Memo) String() string {
	parts := make([]string, len(m.lines))
	for i, ln := range m.lines {
		parts[i] = string(ln)
	}
	return strings.Join(parts, "\n")
}

// SetString replaces the buffer with s, splitting on newlines, and
// puts the cursor at the end of the last line.
func (m *Memo) SetString(s string) {
	parts := strings.Split(s, "\n")
	m.lines = make([][]rune, len(parts))
	for i, p := range parts {
		m.lines[i] = []rune(p)
	}
	if len(m.lines) == 0 {
		m.lines = [][]rune{{}}
	}
	m.curLine = len(m.lines) - 1
	m.curCol = len(m.lines[m.curLine])
	m.adjustView()
}

// Lines returns the current line slice. The slice and its rune slices
// are shared with the memo; callers must not mutate concurrently.
func (m *Memo) Lines() [][]rune { return m.lines }

// Cursor returns the (line, col) cursor position.
func (m *Memo) CursorPos() (int, int) { return m.curLine, m.curCol }

// TopLine returns the first visible line index.
func (m *Memo) TopLine() int { return m.topLine }

// Draw paints visible lines and positions the on-screen caret.
func (m *Memo) Draw(s *vio.Surface) {
	slot := SlotInputNormal
	if m.HasState(view.StateFocused) {
		slot = SlotInputFocused
	}
	body := m.MapColor(slot)

	s.FillRect(m.Bounds, ' ', body)
	if m.Bounds.W <= 0 || m.Bounds.H <= 0 {
		return
	}

	for row := range m.Bounds.H {
		idx := m.topLine + row
		if idx >= len(m.lines) {
			break
		}
		ln := m.lines[idx]
		x := m.Bounds.X
		for col := range m.Bounds.W {
			runeIdx := m.leftCol + col
			if runeIdx >= len(ln) {
				break
			}
			s.Set(x+col, m.Bounds.Y+row, ln[runeIdx], body)
		}
	}

	if m.HasState(view.StateFocused) {
		m.Cursor.X = m.Bounds.X + (m.curCol - m.leftCol)
		m.Cursor.Y = m.Bounds.Y + (m.curLine - m.topLine)
	}
}

// HandleEvent dispatches one event. Non-key events fall through.
func (m *Memo) HandleEvent(e *event.Event) {
	switch e.What {
	case event.ClassKey:
		m.handleKey(e)
	case event.ClassMouseDown:
		m.handleMouse(e)
	}
}

func (m *Memo) handleKey(e *event.Event) {
	if !m.HasState(view.StateFocused) {
		return
	}
	switch e.Key.Key {
	case event.KeyRune:
		if e.Key.Mod&(event.ModCtrl|event.ModAlt) != 0 {
			return
		}
		m.insertRune(e.Key.Rune)
		e.Clear()
	case event.KeyEnter:
		m.insertNewline()
		e.Clear()
	case event.KeyTab:
		for range m.TabWidth {
			m.insertRune(' ')
		}
		e.Clear()
	case event.KeyBackspace:
		m.deletePrev()
		e.Clear()
	case event.KeyDelete:
		m.deleteNext()
		e.Clear()
	case event.KeyArrowLeft:
		m.moveLeft()
		e.Clear()
	case event.KeyArrowRight:
		m.moveRight()
		e.Clear()
	case event.KeyArrowUp:
		m.moveVertical(-1)
		e.Clear()
	case event.KeyArrowDown:
		m.moveVertical(1)
		e.Clear()
	case event.KeyHome:
		m.curCol = 0
		m.adjustView()
		e.Clear()
	case event.KeyEnd:
		m.curCol = len(m.lines[m.curLine])
		m.adjustView()
		e.Clear()
	case event.KeyPgUp:
		m.moveVertical(-max(m.Bounds.H-1, 1))
		e.Clear()
	case event.KeyPgDn:
		m.moveVertical(max(m.Bounds.H-1, 1))
		e.Clear()
	}
}

func (m *Memo) handleMouse(e *event.Event) {
	p := vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}
	if !m.Bounds.Contains(p) {
		return
	}
	line := m.topLine + (p.Y - m.Bounds.Y)
	if line >= len(m.lines) {
		line = len(m.lines) - 1
	}
	m.curLine = line
	m.curCol = min(m.leftCol+(p.X-m.Bounds.X), len(m.lines[line]))
	m.adjustView()
	e.Clear()
}

func (m *Memo) insertRune(r rune) {
	ln := m.lines[m.curLine]
	ln = append(ln, 0)
	copy(ln[m.curCol+1:], ln[m.curCol:])
	ln[m.curCol] = r
	m.lines[m.curLine] = ln
	m.curCol++
	m.adjustView()
}

func (m *Memo) insertNewline() {
	ln := m.lines[m.curLine]
	left := append([]rune{}, ln[:m.curCol]...)
	right := append([]rune{}, ln[m.curCol:]...)
	m.lines[m.curLine] = left
	m.lines = append(m.lines, nil)
	copy(m.lines[m.curLine+2:], m.lines[m.curLine+1:])
	m.lines[m.curLine+1] = right
	m.curLine++
	m.curCol = 0
	m.adjustView()
}

func (m *Memo) deletePrev() {
	if m.curCol > 0 {
		ln := m.lines[m.curLine]
		m.lines[m.curLine] = append(ln[:m.curCol-1], ln[m.curCol:]...)
		m.curCol--
		m.adjustView()
		return
	}
	if m.curLine == 0 {
		return
	}
	prev := m.lines[m.curLine-1]
	cur := m.lines[m.curLine]
	m.curCol = len(prev)
	m.lines[m.curLine-1] = append(prev, cur...)
	m.lines = append(m.lines[:m.curLine], m.lines[m.curLine+1:]...)
	m.curLine--
	m.adjustView()
}

func (m *Memo) deleteNext() {
	ln := m.lines[m.curLine]
	if m.curCol < len(ln) {
		m.lines[m.curLine] = append(ln[:m.curCol], ln[m.curCol+1:]...)
		return
	}
	if m.curLine == len(m.lines)-1 {
		return
	}
	next := m.lines[m.curLine+1]
	m.lines[m.curLine] = append(ln, next...)
	m.lines = append(m.lines[:m.curLine+1], m.lines[m.curLine+2:]...)
}

func (m *Memo) moveLeft() {
	if m.curCol > 0 {
		m.curCol--
	} else if m.curLine > 0 {
		m.curLine--
		m.curCol = len(m.lines[m.curLine])
	}
	m.adjustView()
}

func (m *Memo) moveRight() {
	if m.curCol < len(m.lines[m.curLine]) {
		m.curCol++
	} else if m.curLine < len(m.lines)-1 {
		m.curLine++
		m.curCol = 0
	}
	m.adjustView()
}

func (m *Memo) moveVertical(delta int) {
	m.curLine = clampInt(m.curLine+delta, 0, len(m.lines)-1)
	if m.curCol > len(m.lines[m.curLine]) {
		m.curCol = len(m.lines[m.curLine])
	}
	m.adjustView()
}

func (m *Memo) adjustView() {
	if m.Bounds.H > 0 {
		if m.curLine < m.topLine {
			m.topLine = m.curLine
		}
		if m.curLine >= m.topLine+m.Bounds.H {
			m.topLine = m.curLine - m.Bounds.H + 1
		}
	}
	if m.Bounds.W > 0 {
		if m.curCol < m.leftCol {
			m.leftCol = m.curCol
		}
		if m.curCol >= m.leftCol+m.Bounds.W {
			m.leftCol = m.curCol - m.Bounds.W + 1
		}
	}
	if m.topLine < 0 {
		m.topLine = 0
	}
	if m.leftCol < 0 {
		m.leftCol = 0
	}
}
