package menu

import (
	"github.com/tamnd/vigo/internal/event"
	"github.com/tamnd/vigo/internal/view"
	"github.com/tamnd/vigo/internal/vio"
)

// Hint is a single label on the status line: a key indicator (e.g. "F1") and
// the action it triggers ("Help"). The Cmd field is the command broadcast
// when the key is pressed; v0.1 paints the labels but does not fire commands
// yet (that ships with menu activation in v0.2).
type Hint struct {
	Key    string
	Action string
	Cmd    event.CommandID
}

// DefaultHints is the v0.1 status line: F1 Help, F10 Menu, Alt-X Exit.
//
//nolint:gochecknoglobals // immutable default
var DefaultHints = []Hint{
	{Key: "F1", Action: "Help", Cmd: event.CmdHelp},
	{Key: "F10", Action: "Menu", Cmd: event.CmdMenu},
	{Key: "Alt-X", Action: "Exit", Cmd: event.CmdQuit},
}

// Line is the one-row status bar pinned to the bottom of the screen.
type Line struct {
	*view.View

	Hints []Hint
}

// NewLine returns a status line of width w, anchored at row y, with the
// supplied hints. Pass DefaultHints for the IDE defaults.
func NewLine(w, y int, hints []Hint) *Line {
	v := view.NewView(vio.Rect{X: 0, Y: y, W: w, H: 1})
	v.GrowMode = view.GrowLoY | view.GrowHiX | view.GrowHiY
	v.EventMask = view.EventMaskDefault
	return &Line{View: v, Hints: hints}
}

// Draw paints the bar background and each hint as "<key> <action>".
func (l *Line) Draw(s *vio.Surface) {
	normal := l.MapColor(6)   // status normal
	shortcut := l.MapColor(9) // status shortcut

	s.FillRect(l.Bounds, ' ', normal)

	x := l.Bounds.X + 1
	y := l.Bounds.Y
	for _, h := range l.Hints {
		x = drawText(s, x, y, h.Key, shortcut, l.Bounds.Right())
		if x >= l.Bounds.Right() {
			return
		}
		s.Set(x, y, ' ', normal)
		x++
		x = drawText(s, x, y, h.Action, normal, l.Bounds.Right())
		x += 2
	}
}

// HandleEvent is a no-op for v0.1; the demo binary handles Alt-X directly.
func (l *Line) HandleEvent(e *event.Event) { _ = e }

func drawText(s *vio.Surface, x, y int, str string, attr vio.Attr, right int) int {
	for _, r := range str {
		if x >= right {
			return x
		}
		s.Set(x, y, r, attr)
		x++
	}
	return x
}
