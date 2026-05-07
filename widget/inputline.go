package widget

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// InputLine is a single-line text editor. It exposes Data for the
// caller to read or seed; the cursor and selection are kept in private
// fields and updated by HandleEvent.
//
// Phase 3 supports printable insertion, Backspace, Delete, arrow-key
// navigation, Home/End, Ctrl-Y (clear), and shift-extend selection.
// Word-wrap, undo, and the History dropdown ride along in phase 4.
type InputLine struct {
	*view.View

	Data   []rune
	MaxLen int

	cursor int // rune index of the insertion point
	anchor int // rune index of the selection anchor (== cursor when no selection)
	scroll int // first visible rune index
}

// NewInputLine returns an InputLine with the given bounds and length cap.
// maxLen <= 0 disables the cap.
func NewInputLine(bounds vio.Rect, maxLen int) *InputLine {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.Options |= view.OptSelectable
	v.PaletteIndex = SlotInputNormal
	return &InputLine{View: v, MaxLen: maxLen}
}

// String returns the current edit value.
func (i *InputLine) String() string { return string(i.Data) }

// SetString replaces the data. The cursor moves to the end and any
// selection clears.
func (i *InputLine) SetString(s string) {
	i.Data = []rune(s)
	if i.MaxLen > 0 && len(i.Data) > i.MaxLen {
		i.Data = i.Data[:i.MaxLen]
	}
	i.cursor = len(i.Data)
	i.anchor = i.cursor
	i.adjustScroll()
}

// CursorIndex returns the insertion point as a rune index. It does not
// shadow the View.Cursor field, which holds the on-screen cell offset
// for the rendered caret.
func (i *InputLine) CursorIndex() int { return i.cursor }

// Selection returns the selected range [start, end) as rune indices.
// start == end when nothing is selected.
func (i *InputLine) Selection() (int, int) {
	if i.anchor <= i.cursor {
		return i.anchor, i.cursor
	}
	return i.cursor, i.anchor
}

// Draw paints the field. The selection is highlighted in the selection
// slot; otherwise the line uses normal/focused color.
func (i *InputLine) Draw(s *vio.Surface) {
	slot := SlotInputNormal
	if i.HasState(view.StateFocused) {
		slot = SlotInputFocused
	}
	body := i.MapColor(slot)
	sel := i.MapColor(SlotInputSelection)

	s.FillRect(i.Bounds, ' ', body)
	if i.Bounds.W <= 0 || i.Bounds.H <= 0 {
		return
	}

	selStart, selEnd := i.Selection()
	for col := 0; col < i.Bounds.W; col++ {
		idx := i.scroll + col
		if idx >= len(i.Data) {
			break
		}
		attr := body
		if idx >= selStart && idx < selEnd {
			attr = sel
		}
		s.Set(i.Bounds.X+col, i.Bounds.Y, i.Data[idx], attr)
	}

	if i.HasState(view.StateFocused) {
		i.Cursor.X = i.Bounds.X + (i.cursor - i.scroll)
		i.Cursor.Y = i.Bounds.Y
	}
}

// HandleEvent processes one keystroke. Non-key events fall through.
func (i *InputLine) HandleEvent(e *event.Event) {
	if e.What != event.ClassKey {
		return
	}
	if !i.HasState(view.StateFocused) {
		return
	}

	shift := e.Key.Mod&event.ModShift != 0
	switch e.Key.Key {
	case event.KeyRune:
		if e.Key.Mod&event.ModCtrl != 0 {
			i.handleCtrlRune(e)
			return
		}
		if e.Key.Mod&event.ModAlt != 0 {
			return // dialog mnemonic dispatch handles Alt-keys
		}
		i.insertRune(e.Key.Rune)
		e.Clear()
	case event.KeyBackspace:
		i.deletePrev()
		e.Clear()
	case event.KeyDelete:
		i.deleteNext()
		e.Clear()
	case event.KeyArrowLeft:
		i.moveCursor(-1, shift)
		e.Clear()
	case event.KeyArrowRight:
		i.moveCursor(1, shift)
		e.Clear()
	case event.KeyHome:
		i.moveTo(0, shift)
		e.Clear()
	case event.KeyEnd:
		i.moveTo(len(i.Data), shift)
		e.Clear()
	}
}

func (i *InputLine) handleCtrlRune(e *event.Event) {
	if e.Key.Rune == 'y' {
		i.Data = i.Data[:0]
		i.cursor = 0
		i.anchor = 0
		i.scroll = 0
		e.Clear()
	}
}

func (i *InputLine) insertRune(r rune) {
	if i.anchor != i.cursor {
		i.deleteSelection()
	}
	if i.MaxLen > 0 && len(i.Data) >= i.MaxLen {
		return
	}
	i.Data = append(i.Data, 0)
	copy(i.Data[i.cursor+1:], i.Data[i.cursor:])
	i.Data[i.cursor] = r
	i.cursor++
	i.anchor = i.cursor
	i.adjustScroll()
}

func (i *InputLine) deletePrev() {
	if i.anchor != i.cursor {
		i.deleteSelection()
		return
	}
	if i.cursor == 0 {
		return
	}
	i.Data = append(i.Data[:i.cursor-1], i.Data[i.cursor:]...)
	i.cursor--
	i.anchor = i.cursor
	i.adjustScroll()
}

func (i *InputLine) deleteNext() {
	if i.anchor != i.cursor {
		i.deleteSelection()
		return
	}
	if i.cursor >= len(i.Data) {
		return
	}
	i.Data = append(i.Data[:i.cursor], i.Data[i.cursor+1:]...)
}

func (i *InputLine) deleteSelection() {
	a, b := i.Selection()
	i.Data = append(i.Data[:a], i.Data[b:]...)
	i.cursor = a
	i.anchor = a
	i.adjustScroll()
}

func (i *InputLine) moveCursor(delta int, extend bool) {
	i.cursor += delta
	if i.cursor < 0 {
		i.cursor = 0
	}
	if i.cursor > len(i.Data) {
		i.cursor = len(i.Data)
	}
	if !extend {
		i.anchor = i.cursor
	}
	i.adjustScroll()
}

func (i *InputLine) moveTo(idx int, extend bool) {
	if idx < 0 {
		idx = 0
	}
	if idx > len(i.Data) {
		idx = len(i.Data)
	}
	i.cursor = idx
	if !extend {
		i.anchor = i.cursor
	}
	i.adjustScroll()
}

func (i *InputLine) adjustScroll() {
	if i.Bounds.W <= 0 {
		return
	}
	if i.cursor < i.scroll {
		i.scroll = i.cursor
	}
	if i.cursor >= i.scroll+i.Bounds.W {
		i.scroll = i.cursor - i.Bounds.W + 1
	}
}
