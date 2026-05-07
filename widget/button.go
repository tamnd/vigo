package widget

import (
	"unicode"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// ButtonFlags is a bitmask of behavior flags for a Button. The names
// match Turbo Vision's bfXxx constants.
type ButtonFlags uint8

// Button flag bits.
const (
	BfNormal    ButtonFlags = 0
	BfDefault   ButtonFlags = 1 << iota // fired on Enter regardless of focus
	BfLeftJust                          // left-justify the title (default is centered)
	BfBroadcast                         // emit ClassBroadcast instead of ClassCommand
)

// Button is a clickable control that emits a command. Mouse-down inside
// the bounds, Space or Enter while focused, or Alt-<mnemonic> all fire
// the same command.
type Button struct {
	*view.View

	Title    string
	Command  event.CommandID
	Flags    ButtonFlags
	mnemonic rune
}

// NewButton returns a Button sized to bounds. The title may include a
// '~x~' marker to advertise an Alt-x hot-key.
func NewButton(bounds vio.Rect, title string, cmd event.CommandID, flags ButtonFlags) *Button {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.Options |= view.OptSelectable
	v.PaletteIndex = SlotButtonNormal

	plain, mn := parseMnemonic(title)
	return &Button{
		View:     v,
		Title:    plain,
		Command:  cmd,
		Flags:    flags,
		mnemonic: mn,
	}
}

// Mnemonic returns the parsed hot-key rune (always lowercase), or 0.
func (b *Button) Mnemonic() rune { return b.mnemonic }

// IsDefault reports whether the button fires on Enter without focus.
func (b *Button) IsDefault() bool { return b.Flags&BfDefault != 0 }

// Draw paints the button. The title sits inside square brackets when
// the button is the default; otherwise plain. The focused button uses
// the focused palette slot.
func (b *Button) Draw(s *vio.Surface) {
	slot := SlotButtonNormal
	if b.HasState(view.StateFocused) {
		slot = SlotButtonFocused
	}
	if b.HasState(view.StateDisabled) {
		slot = SlotButtonDisabled
	}
	body := b.MapColor(slot)
	hot := b.MapColor(SlotButtonShortcut)

	s.FillRect(b.Bounds, ' ', body)
	if b.Bounds.W <= 0 || b.Bounds.H <= 0 {
		return
	}

	title := b.Title
	if b.IsDefault() {
		title = "<" + title + ">"
	} else {
		title = " " + title + " "
	}

	titleRunes := []rune(title)
	if len(titleRunes) > b.Bounds.W {
		titleRunes = titleRunes[:b.Bounds.W]
	}

	x := b.Bounds.X
	if b.Flags&BfLeftJust == 0 {
		x += (b.Bounds.W - len(titleRunes)) / 2
	}
	y := b.Bounds.Y + b.Bounds.H/2

	mnSeen := false
	for _, r := range titleRunes {
		attr := body
		if !mnSeen && b.mnemonic != 0 && unicode.ToLower(r) == b.mnemonic {
			attr = hot
			mnSeen = true
		}
		s.Set(x, y, r, attr)
		x++
	}
}

// HandleEvent fires the command on Space, Enter, mouse-down inside the
// bounds, or Alt-<mnemonic>.
func (b *Button) HandleEvent(e *event.Event) {
	if b.HasState(view.StateDisabled) {
		return
	}
	switch e.What {
	case event.ClassKey:
		b.handleKey(e)
	case event.ClassMouseDown:
		if b.Bounds.Contains(vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}) {
			b.fire()
			e.Clear()
		}
	}
}

func (b *Button) handleKey(e *event.Event) {
	switch {
	case e.Key.Key == event.KeyEnter,
		e.Key.Key == event.KeyRune && e.Key.Rune == ' ' && e.Key.Mod == event.ModNone:
		if b.HasState(view.StateFocused) {
			b.fire()
			e.Clear()
		}
	case e.Key.Mod&event.ModAlt != 0 && e.Key.Key == event.KeyRune && b.mnemonic != 0:
		if unicode.ToLower(e.Key.Rune) == b.mnemonic {
			b.fire()
			e.Clear()
		}
	}
}

// Press programmatically fires the button. Useful for Dialog Enter
// dispatch and for tests.
func (b *Button) Press() { b.fire() }

func (b *Button) fire() {
	if b.Owner == nil {
		return
	}
	what := event.ClassCommand
	if b.Flags&BfBroadcast != 0 {
		what = event.ClassBroadcast
	}
	ev := &event.Event{
		What: what,
		Msg:  event.MessageEvent{Command: b.Command, Info: b},
	}
	b.Owner.HandleEvent(ev)
}
