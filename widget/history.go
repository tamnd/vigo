package widget

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// HistoryID groups history entries by purpose: each numeric ID owns its
// own ring of recent strings. Two views that share an ID share their
// history. Zero means "no history" and HistoryAdd ignores it.
type HistoryID uint16

// HistoryLimit caps the number of strings kept per ID. Older entries
// fall off the back when the ring fills up.
const HistoryLimit = 16

var histTable = make(map[HistoryID][]string)

// HistoryAdd pushes s onto the front of the ring for id. If s already
// exists in the ring it is moved to the front (no duplicates). The
// ring is capped at HistoryLimit; the oldest entry is discarded when
// it overflows. Empty strings and id == 0 are ignored.
func HistoryAdd(id HistoryID, s string) {
	if id == 0 || s == "" {
		return
	}
	ring := histTable[id]
	for i, v := range ring {
		if v == s {
			ring = append(ring[:i], ring[i+1:]...)
			break
		}
	}
	ring = append([]string{s}, ring...)
	if len(ring) > HistoryLimit {
		ring = ring[:HistoryLimit]
	}
	histTable[id] = ring
}

// HistoryCount returns the number of stored strings for id.
func HistoryCount(id HistoryID) int { return len(histTable[id]) }

// HistoryStr returns the n-th string in the ring for id, with 0 the
// most recent. Out-of-range indices return "".
func HistoryStr(id HistoryID, n int) string {
	ring := histTable[id]
	if n < 0 || n >= len(ring) {
		return ""
	}
	return ring[n]
}

// HistoryReset clears all rings. Tests use this to start clean.
func HistoryReset() { histTable = make(map[HistoryID][]string) }

// History is the small dropdown-arrow button that sits next to an
// InputLine. A click or a focused Enter sets Open to true; the host
// dialog watches Open and pushes a HistoryWindow when it flips. Phase
// 4 keeps the trigger simple; richer integration with ExecView lands
// alongside the menu/cmd work in phase 5.
type History struct {
	*view.View

	ID   HistoryID
	Open bool
}

// NewHistory returns a History trigger attached to id.
func NewHistory(bounds vio.Rect, id HistoryID) *History {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.PaletteIndex = SlotButtonNormal
	v.Options |= view.OptSelectable
	v.EventMask = event.ClassKey | event.ClassMouseDown
	return &History{View: v, ID: id}
}

// Draw paints the trigger as a "[▼]" cap.
func (h *History) Draw(s *vio.Surface) {
	attr := h.MapColor(SlotButtonNormal)
	if h.HasState(view.StateFocused) {
		attr = h.MapColor(SlotButtonFocused)
	}
	s.FillRect(h.Bounds, ' ', attr)
	if h.Bounds.W >= 3 {
		s.Set(h.Bounds.X, h.Bounds.Y, '[', attr)
		s.Set(h.Bounds.X+1, h.Bounds.Y, '▼', attr)
		s.Set(h.Bounds.X+2, h.Bounds.Y, ']', attr)
	}
}

// HandleEvent flips Open on click or Enter. Hosts reset Open after
// they spawn the dropdown window.
func (h *History) HandleEvent(e *event.Event) {
	switch e.What {
	case event.ClassMouseDown:
		if h.Bounds.Contains(vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}) {
			h.Open = true
			e.Clear()
		}
	case event.ClassKey:
		if h.HasState(view.StateFocused) && e.Key.Key == event.KeyEnter {
			h.Open = true
			e.Clear()
		}
	}
}

// HistoryWindow is the dropdown list that appears under a History
// trigger. It hosts a ListBox seeded from the ring for ID. Enter on
// the focused row stores Selected and sets Result to CmdOk; Esc sets
// Result to CmdCancel. Hosts insert a HistoryWindow into a Group and
// poll Result like any other modal view.
type HistoryWindow struct {
	*view.Group

	ID       HistoryID
	List     *ListBox
	Selected string

	result event.CommandID
}

// NewHistoryWindow returns a dropdown sized to bounds, populated from
// the ring for id. The internal ListBox is inserted as the first
// child.
func NewHistoryWindow(bounds vio.Rect, id HistoryID) *HistoryWindow {
	g := view.NewGroup(bounds)
	g.Palette = vio.BorlandClassic
	g.PaletteIndex = SlotListNormal

	items := make([]string, HistoryCount(id))
	for i := range items {
		items[i] = HistoryStr(id, i)
	}
	listBounds := vio.R(0, 0, bounds.W, bounds.H)
	list := NewListBox(listBounds, items, nil)
	g.Insert(list)
	g.SetCurrent(0)

	return &HistoryWindow{
		Group: g,
		ID:    id,
		List:  list,
	}
}

// Result returns the modal command, or CmdNone while still open.
func (w *HistoryWindow) Result() event.CommandID { return w.result }

// HandleEvent intercepts Enter and Esc, then defers everything else
// to the embedded Group (which feeds the ListBox).
func (w *HistoryWindow) HandleEvent(e *event.Event) {
	if e.What == event.ClassKey {
		switch e.Key.Key {
		case event.KeyEnter:
			w.Selected = w.List.Selected()
			w.result = event.CmdOk
			e.Clear()
			return
		case event.KeyEsc:
			w.result = event.CmdCancel
			e.Clear()
			return
		}
	}
	w.Group.HandleEvent(e)
}
