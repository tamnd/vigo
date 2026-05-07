package widget

import (
	"unicode"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// Cluster is the shared backing for CheckBoxes and RadioButtons. It
// holds the option titles, the parsed mnemonic letters, the focused
// item index, and a Value field whose meaning is owned by the
// concrete type (bitmask for check boxes, index for radio buttons).
type Cluster struct {
	*view.View

	Items []string
	Value uint32

	titles    []string
	mnemonics []rune
	sel       int
}

func newCluster(bounds vio.Rect, items []string) *Cluster {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.Options |= view.OptSelectable
	v.PaletteIndex = SlotClusterNormal

	titles := make([]string, len(items))
	mns := make([]rune, len(items))
	for i, raw := range items {
		titles[i], mns[i] = parseMnemonic(raw)
	}
	return &Cluster{
		View:      v,
		Items:     items,
		titles:    titles,
		mnemonics: mns,
	}
}

// Sel returns the focused item index.
func (c *Cluster) Sel() int { return c.sel }

// SetSel moves the focus to item i, clamped to a valid index.
func (c *Cluster) SetSel(i int) {
	if len(c.titles) == 0 {
		c.sel = 0
		return
	}
	c.sel = clampInt(i, 0, len(c.titles)-1)
}

// drawItem paints one option line. mark says whether the item is
// "set" (toggled checkbox / chosen radio); openGlyph and closeGlyph
// surround the mark.
func (c *Cluster) drawItem(surf *vio.Surface, idx int, openGlyph, closeGlyph, markGlyph rune) {
	body := c.MapColor(SlotClusterNormal)
	focused := c.MapColor(SlotClusterFocused)
	hot := c.MapColor(SlotClusterHotKey)

	row := c.Bounds.Y + idx
	if row < c.Bounds.Y || row >= c.Bounds.Bottom() {
		return
	}
	attr := body
	if c.HasState(view.StateFocused) && idx == c.sel {
		attr = focused
	}

	x := c.Bounds.X
	surf.Set(x, row, openGlyph, attr)
	mark := ' '
	if markGlyph != 0 {
		mark = markGlyph
	}
	surf.Set(x+1, row, mark, attr)
	surf.Set(x+2, row, closeGlyph, attr)
	surf.Set(x+3, row, ' ', attr)

	pos := x + 4
	right := c.Bounds.Right()
	mn := c.mnemonics[idx]
	mnSeen := false
	for _, r := range c.titles[idx] {
		if pos >= right {
			break
		}
		ch := attr
		if !mnSeen && mn != 0 && unicode.ToLower(r) == mn {
			ch = hot
			mnSeen = true
		}
		surf.Set(pos, row, r, ch)
		pos++
	}
	for ; pos < right; pos++ {
		surf.Set(pos, row, ' ', attr)
	}

	if c.HasState(view.StateFocused) && idx == c.sel {
		c.Cursor.X = c.Bounds.X + 1
		c.Cursor.Y = row
	}
}

// handleNav is the shared key handling for cluster motion: arrow keys
// move sel, Home/End jump to ends. Returns true if the event was
// consumed.
func (c *Cluster) handleNav(e *event.Event) bool {
	if e.What != event.ClassKey {
		return false
	}
	switch e.Key.Key {
	case event.KeyArrowUp:
		c.SetSel(c.sel - 1)
	case event.KeyArrowDown:
		c.SetSel(c.sel + 1)
	case event.KeyHome:
		c.SetSel(0)
	case event.KeyEnd:
		c.SetSel(len(c.titles) - 1)
	default:
		return false
	}
	e.Clear()
	return true
}

// findMnemonic returns the index whose mnemonic matches r, or -1.
func (c *Cluster) findMnemonic(r rune) int {
	r = unicode.ToLower(r)
	for i, mn := range c.mnemonics {
		if mn != 0 && mn == r {
			return i
		}
	}
	return -1
}

// hitItem returns the row index in the cluster that contains p, or -1.
func (c *Cluster) hitItem(p vio.Point) int {
	if !c.Bounds.Contains(p) {
		return -1
	}
	idx := p.Y - c.Bounds.Y
	if idx < 0 || idx >= len(c.titles) {
		return -1
	}
	return idx
}

// CheckBoxes is a Cluster where each item is independently toggled.
// Value is a bitmask: bit i is set when item i is checked.
type CheckBoxes struct {
	*Cluster
}

// NewCheckBoxes returns a check-box cluster sized to bounds. Each item
// title may include a '~x~' mnemonic marker for Alt-x dispatch.
func NewCheckBoxes(bounds vio.Rect, items []string) *CheckBoxes {
	return &CheckBoxes{Cluster: newCluster(bounds, items)}
}

// Mark reports whether item i is checked.
func (c *CheckBoxes) Mark(i int) bool { return c.Value&(1<<i) != 0 }

// Press toggles item i.
func (c *CheckBoxes) Press(i int) {
	if i < 0 || i >= len(c.titles) {
		return
	}
	c.Value ^= 1 << i
}

// Draw paints the check-box list.
func (c *CheckBoxes) Draw(s *vio.Surface) {
	body := c.MapColor(SlotClusterNormal)
	s.FillRect(c.Bounds, ' ', body)
	for i := range c.titles {
		var mark rune
		if c.Mark(i) {
			mark = 'X'
		}
		c.drawItem(s, i, '[', ']', mark)
	}
}

// HandleEvent moves focus, toggles on Space/Enter, and dispatches
// Alt-mnemonic.
func (c *CheckBoxes) HandleEvent(e *event.Event) {
	if c.handleNav(e) {
		return
	}
	if e.What == event.ClassKey {
		k := e.Key
		switch {
		case k.Key == event.KeyEnter,
			k.Key == event.KeyRune && k.Rune == ' ' && k.Mod == event.ModNone:
			if c.HasState(view.StateFocused) {
				c.Press(c.sel)
				e.Clear()
				return
			}
		case k.Mod&event.ModAlt != 0 && k.Key == event.KeyRune:
			if idx := c.findMnemonic(k.Rune); idx >= 0 {
				c.Press(idx)
				c.SetSel(idx)
				e.Clear()
				return
			}
		}
	}
	if e.What == event.ClassMouseDown {
		if idx := c.hitItem(vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}); idx >= 0 {
			c.SetSel(idx)
			c.Press(idx)
			e.Clear()
		}
	}
}

// RadioButtons is a Cluster with single selection: Value is the index
// of the chosen item.
type RadioButtons struct {
	*Cluster
}

// NewRadioButtons returns a radio-button cluster sized to bounds.
func NewRadioButtons(bounds vio.Rect, items []string) *RadioButtons {
	return &RadioButtons{Cluster: newCluster(bounds, items)}
}

// Mark reports whether item i is the chosen one.
func (r *RadioButtons) Mark(i int) bool {
	if i < 0 {
		return false
	}
	return uint32(i) == r.Value //nolint:gosec // i bounded by len(titles)
}

// Press makes item i the chosen one.
func (r *RadioButtons) Press(i int) {
	if i < 0 || i >= len(r.titles) {
		return
	}
	r.Value = uint32(i) //nolint:gosec // bounded by titles length
}

// Draw paints the radio-button list.
func (r *RadioButtons) Draw(s *vio.Surface) {
	body := r.MapColor(SlotClusterNormal)
	s.FillRect(r.Bounds, ' ', body)
	for i := range r.titles {
		var mark rune
		if r.Mark(i) {
			mark = '•'
		}
		r.drawItem(s, i, '(', ')', mark)
	}
}

// HandleEvent moves focus, selects on Space/Enter, and dispatches
// Alt-mnemonic. Selecting an item also moves focus to that item.
func (r *RadioButtons) HandleEvent(e *event.Event) {
	if r.handleNav(e) {
		return
	}
	if e.What == event.ClassKey {
		k := e.Key
		switch {
		case k.Key == event.KeyEnter,
			k.Key == event.KeyRune && k.Rune == ' ' && k.Mod == event.ModNone:
			if r.HasState(view.StateFocused) {
				r.Press(r.sel)
				e.Clear()
				return
			}
		case k.Mod&event.ModAlt != 0 && k.Key == event.KeyRune:
			if idx := r.findMnemonic(k.Rune); idx >= 0 {
				r.SetSel(idx)
				r.Press(idx)
				e.Clear()
				return
			}
		}
	}
	if e.What == event.ClassMouseDown {
		if idx := r.hitItem(vio.Point{X: e.Mouse.X, Y: e.Mouse.Y}); idx >= 0 {
			r.SetSel(idx)
			r.Press(idx)
			e.Clear()
		}
	}
}
