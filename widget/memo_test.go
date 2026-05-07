package widget

import (
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

func keyEvent(k event.Key) *event.Event {
	return &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: k}}
}

func runeEvent(r rune) *event.Event {
	return &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyRune, Rune: r}}
}

func newFocusedMemo(t *testing.T, w, h int) *Memo {
	t.Helper()
	host := newPalettedHost()
	m := NewMemo(vio.R(0, 0, w, h))
	host.Insert(m)
	m.SetState(view.StateFocused, true)
	return m
}

func TestMemoInitialEmpty(t *testing.T) {
	m := newFocusedMemo(t, 20, 5)
	if m.String() != "" {
		t.Fatalf("empty memo: %q", m.String())
	}
	if l, c := m.CursorPos(); l != 0 || c != 0 {
		t.Fatalf("cursor at origin: line=%d col=%d", l, c)
	}
}

func TestMemoInsertPrintable(t *testing.T) {
	m := newFocusedMemo(t, 20, 5)
	for _, r := range "hi" {
		m.HandleEvent(runeEvent(r))
	}
	if m.String() != "hi" {
		t.Fatalf("after insert: %q", m.String())
	}
	if _, c := m.CursorPos(); c != 2 {
		t.Fatalf("cursor col=%d", c)
	}
}

func TestMemoEnterSplitsLine(t *testing.T) {
	m := newFocusedMemo(t, 20, 5)
	m.SetString("abcdef")
	m.HandleEvent(keyEvent(event.KeyHome))
	for range 3 {
		m.HandleEvent(keyEvent(event.KeyArrowRight))
	}
	m.HandleEvent(keyEvent(event.KeyEnter))
	if m.String() != "abc\ndef" {
		t.Fatalf("split: %q", m.String())
	}
	if l, c := m.CursorPos(); l != 1 || c != 0 {
		t.Fatalf("cursor after Enter: line=%d col=%d", l, c)
	}
}

func TestMemoBackspaceJoinsLines(t *testing.T) {
	m := newFocusedMemo(t, 20, 5)
	m.SetString("abc\ndef")
	m.HandleEvent(keyEvent(event.KeyHome))
	m.HandleEvent(keyEvent(event.KeyArrowDown))
	m.HandleEvent(keyEvent(event.KeyHome))
	m.HandleEvent(keyEvent(event.KeyBackspace))
	if m.String() != "abcdef" {
		t.Fatalf("backspace join: %q", m.String())
	}
	if l, c := m.CursorPos(); l != 0 || c != 3 {
		t.Fatalf("cursor after join: line=%d col=%d", l, c)
	}
}

func TestMemoDeleteJoinsLines(t *testing.T) {
	m := newFocusedMemo(t, 20, 5)
	m.SetString("abc\ndef")
	m.HandleEvent(keyEvent(event.KeyHome))
	m.HandleEvent(keyEvent(event.KeyArrowUp))
	m.HandleEvent(keyEvent(event.KeyEnd))
	m.HandleEvent(keyEvent(event.KeyDelete))
	if m.String() != "abcdef" {
		t.Fatalf("delete join: %q", m.String())
	}
}

func TestMemoTabInsertsSpaces(t *testing.T) {
	m := newFocusedMemo(t, 20, 5)
	m.TabWidth = 4
	m.HandleEvent(keyEvent(event.KeyTab))
	if m.String() != "    " {
		t.Fatalf("tab inserts spaces: %q", m.String())
	}
}

func TestMemoArrowsCrossLines(t *testing.T) {
	m := newFocusedMemo(t, 20, 5)
	m.SetString("abc\ndef")
	m.HandleEvent(keyEvent(event.KeyHome))
	m.HandleEvent(keyEvent(event.KeyArrowUp))
	m.HandleEvent(keyEvent(event.KeyHome))
	m.HandleEvent(keyEvent(event.KeyEnd))
	m.HandleEvent(keyEvent(event.KeyArrowRight))
	if l, c := m.CursorPos(); l != 1 || c != 0 {
		t.Fatalf("right at EOL crossed to next line: line=%d col=%d", l, c)
	}
	m.HandleEvent(keyEvent(event.KeyArrowLeft))
	if l, c := m.CursorPos(); l != 0 || c != 3 {
		t.Fatalf("left at BOL crossed to prev line end: line=%d col=%d", l, c)
	}
}

func TestMemoPgDnPagesByViewport(t *testing.T) {
	m := newFocusedMemo(t, 20, 3)
	lines := make([]string, 10)
	for i := range lines {
		lines[i] = "x"
	}
	m.SetString(strings.Join(lines, "\n"))
	for range len(lines) {
		m.HandleEvent(keyEvent(event.KeyArrowUp))
	}
	if l, _ := m.CursorPos(); l != 0 {
		t.Fatalf("after up to top: line=%d", l)
	}
	m.HandleEvent(keyEvent(event.KeyPgDn))
	if l, _ := m.CursorPos(); l != 2 {
		t.Fatalf("PgDn line=%d", l)
	}
}

func TestMemoScrollsViewportWhenCursorExits(t *testing.T) {
	m := newFocusedMemo(t, 20, 3)
	lines := make([]string, 6)
	for i := range lines {
		lines[i] = "row"
	}
	m.SetString(strings.Join(lines, "\n"))
	for range 5 {
		m.HandleEvent(keyEvent(event.KeyArrowUp))
	}
	for range 5 {
		m.HandleEvent(keyEvent(event.KeyArrowDown))
	}
	if m.TopLine() != 3 {
		t.Fatalf("topLine should slide: %d", m.TopLine())
	}
}

func TestMemoMouseSetsCursor(t *testing.T) {
	m := newFocusedMemo(t, 20, 5)
	m.SetString("abc\ndefgh")
	m.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: 2, Y: 1},
	})
	if l, c := m.CursorPos(); l != 1 || c != 2 {
		t.Fatalf("mouse cursor: line=%d col=%d", l, c)
	}
}

func TestMemoUnfocusedIgnoresKeys(t *testing.T) {
	host := newPalettedHost()
	m := NewMemo(vio.R(0, 0, 20, 5))
	host.Insert(m)
	host.SetCurrent(-1)
	m.SetState(view.StateFocused, false)

	m.HandleEvent(runeEvent('a'))
	if m.String() != "" {
		t.Fatalf("unfocused memo should ignore keys, got %q", m.String())
	}
}

func TestMemoDrawShowsContent(t *testing.T) {
	m := newFocusedMemo(t, 10, 3)
	m.SetString("hello\nworld")

	s := vio.NewSurface(10, 3)
	m.Draw(s)
	rows := s.Snapshot()
	if !strings.HasPrefix(rows[0], "hello") {
		t.Fatalf("row 0: %q", rows[0])
	}
	if !strings.HasPrefix(rows[1], "world") {
		t.Fatalf("row 1: %q", rows[1])
	}
}
