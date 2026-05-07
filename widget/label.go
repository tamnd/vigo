package widget

import (
	"unicode"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Label is a non-focusable text element bound to a peer view. Pressing
// Alt-<mnemonic> on the dialog moves focus to the peer. The mnemonic is
// the rune surrounded by '~' in the title (e.g. "~N~ame:" advertises
// Alt-N and underlines the N).
type Label struct {
	*view.View

	Title    string
	Peer     view.Viewer
	mnemonic rune
}

// NewLabel returns a Label sized to bounds. peer may be nil for a
// caption that just announces text without claiming a hot-key.
func NewLabel(bounds vio.Rect, title string, peer view.Viewer) *Label {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.PaletteIndex = SlotDialogText
	// Labels react to Alt-key broadcasts only; everything else is dialog
	// territory.
	v.EventMask = event.ClassKey | event.ClassBroadcast

	plain, mn := parseMnemonic(title)
	return &Label{View: v, Title: plain, Peer: peer, mnemonic: mn}
}

// Mnemonic returns the parsed hot-key rune (always lowercase), or 0 if
// the title had no '~x~' marker.
func (l *Label) Mnemonic() rune { return l.mnemonic }

// Draw paints the label. The mnemonic rune is drawn in the dialog
// shortcut color so it stands out.
func (l *Label) Draw(s *vio.Surface) {
	plain := l.MapColor(l.PaletteIndex)
	hot := l.MapColor(SlotDialogShortcut)

	s.FillRect(l.Bounds, ' ', plain)
	x := l.Bounds.X
	y := l.Bounds.Y
	if y < 0 || y >= l.Bounds.Y+l.Bounds.H {
		return
	}

	right := l.Bounds.X + l.Bounds.W
	mnSeen := false
	for _, r := range l.Title {
		if x >= right {
			break
		}
		attr := plain
		if !mnSeen && l.mnemonic != 0 && unicode.ToLower(r) == l.mnemonic {
			attr = hot
			mnSeen = true
		}
		s.Set(x, y, r, attr)
		x++
	}
}

// HandleEvent moves focus to the peer when the user presses Alt-<mnemonic>.
// Without a peer, or without a mnemonic, the label is inert.
func (l *Label) HandleEvent(e *event.Event) {
	if e.What != event.ClassKey || l.Peer == nil || l.mnemonic == 0 {
		return
	}
	if e.Key.Mod&event.ModAlt == 0 {
		return
	}
	if e.Key.Key != event.KeyRune {
		return
	}
	if unicode.ToLower(e.Key.Rune) != l.mnemonic {
		return
	}
	if l.Owner == nil {
		return
	}
	for i, c := range l.Owner.Children() {
		if c == l.Peer {
			l.Owner.SetCurrent(i)
			e.Clear()
			return
		}
	}
}

// parseMnemonic extracts the '~x~' rune from s. It returns the plain
// title (with the markers stripped) and the lowercase mnemonic rune, or
// 0 if no marker is present.
func parseMnemonic(s string) (string, rune) {
	out := make([]rune, 0, len(s))
	var mn rune
	inMarker := false
	for _, r := range s {
		if r == '~' {
			inMarker = !inMarker
			continue
		}
		if inMarker && mn == 0 {
			mn = unicode.ToLower(r)
		}
		out = append(out, r)
	}
	return string(out), mn
}
