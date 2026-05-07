package demos

import (
	"fmt"
	"time"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
	"github.com/tamnd/vigo/window"
)

// Calendar is a month-view window. The cursor day defaults to today;
// PgUp/PgDn step months, arrow keys move the cursor day inside the
// current month with overflow into adjacent months.
type Calendar struct {
	*window.Window

	header *widget.StaticText
	grid   *calendarGrid
}

// CalDefaultBounds is the calendar's default footprint: 25x10 leaves
// room for the seven-column day grid (7*3 cells) plus a header and
// frame.
//
//nolint:gochecknoglobals // immutable default
var CalDefaultBounds = vio.R(2, 2, 25, 10)

// NewCalendar returns a Calendar window opened to today's month.
func NewCalendar(bounds vio.Rect) *Calendar {
	now := time.Now()
	return NewCalendarAt(bounds, now)
}

// NewCalendarAt returns a Calendar showing on's month with on's day
// as the cursor. Tests use this for deterministic dates.
func NewCalendarAt(bounds vio.Rect, on time.Time) *Calendar {
	w := window.New(bounds, "Calendar", window.FlagMove|window.FlagClose)
	c := &Calendar{Window: w}

	client := w.ClientRect()
	c.header = widget.NewStaticText(
		vio.R(client.X+1, client.Y, client.W-2, 1),
		"",
	)
	w.Insert(c.header)

	c.grid = newCalendarGrid(
		vio.R(client.X+1, client.Y+1, client.W-2, client.H-1),
		on,
	)
	w.Insert(c.grid)

	c.refreshHeader()
	return c
}

// Cursor returns the currently selected date.
func (c *Calendar) Cursor() time.Time { return c.grid.cursor }

// HandleEvent intercepts navigation keys and falls back to the
// Window dispatch for everything else.
func (c *Calendar) HandleEvent(e *event.Event) {
	if e.What == event.ClassKey && c.handleKey(e) {
		c.refreshHeader()
		e.Clear()
		return
	}
	c.Window.HandleEvent(e)
}

func (c *Calendar) handleKey(e *event.Event) bool {
	switch e.Key.Key {
	case event.KeyPgUp:
		c.grid.cursor = addMonths(c.grid.cursor, -1)
	case event.KeyPgDn:
		c.grid.cursor = addMonths(c.grid.cursor, 1)
	case event.KeyArrowLeft:
		c.grid.cursor = c.grid.cursor.AddDate(0, 0, -1)
	case event.KeyArrowRight:
		c.grid.cursor = c.grid.cursor.AddDate(0, 0, 1)
	case event.KeyArrowUp:
		c.grid.cursor = c.grid.cursor.AddDate(0, 0, -7)
	case event.KeyArrowDown:
		c.grid.cursor = c.grid.cursor.AddDate(0, 0, 7)
	case event.KeyHome:
		c.grid.cursor = time.Date(c.grid.cursor.Year(), c.grid.cursor.Month(), 1, 0, 0, 0, 0, c.grid.cursor.Location())
	default:
		return false
	}
	return true
}

func (c *Calendar) refreshHeader() {
	t := c.grid.cursor
	c.header.Text = fmt.Sprintf("   %s %d", t.Month().String(), t.Year())
}

// addMonths returns t shifted by n months, clamped so the day-of-
// month stays valid in the destination (Jan 31 + 1mo -> Feb 28/29).
func addMonths(t time.Time, n int) time.Time {
	y, m, d := t.Date()
	target := time.Date(y, m+time.Month(n), 1, 0, 0, 0, 0, t.Location())
	last := lastDayOfMonth(target)
	if d > last {
		d = last
	}
	return time.Date(target.Year(), target.Month(), d, 0, 0, 0, 0, t.Location())
}

func lastDayOfMonth(t time.Time) int {
	first := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	return first.AddDate(0, 1, -1).Day()
}

// calendarGrid paints a 7-column day grid for the cursor's month.
type calendarGrid struct {
	*view.View

	cursor time.Time
}

func newCalendarGrid(bounds vio.Rect, on time.Time) *calendarGrid {
	v := view.NewView(bounds)
	v.GrowMode = view.GrowFixed
	v.EventMask = 0 // parent window dispatches keys
	v.PaletteIndex = widget.SlotDialogText
	return &calendarGrid{View: v, cursor: on}
}

// Draw paints the day-of-week header row and the dates of the
// cursor's month, highlighting the cursor day in the selected slot.
func (g *calendarGrid) Draw(s *vio.Surface) {
	normal := g.MapColor(g.PaletteIndex)
	selected := g.MapColor(widget.SlotListSelected)

	s.FillRect(g.Bounds, ' ', normal)

	// Day-of-week strip.
	dow := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	for i, d := range dow {
		x := g.Bounds.X + i*3
		if x+2 > g.Bounds.Right() {
			break
		}
		s.DrawString(x, g.Bounds.Y, d, normal)
	}

	first := time.Date(g.cursor.Year(), g.cursor.Month(), 1, 0, 0, 0, 0, g.cursor.Location())
	startCol := int(first.Weekday()) // 0..6, Sunday=0
	last := lastDayOfMonth(g.cursor)

	for d := 1; d <= last; d++ {
		idx := startCol + d - 1
		row := idx / 7
		col := idx % 7
		x := g.Bounds.X + col*3
		y := g.Bounds.Y + 1 + row
		if y >= g.Bounds.Bottom() {
			break
		}
		if x+2 > g.Bounds.Right() {
			continue
		}
		attr := normal
		if d == g.cursor.Day() {
			attr = selected
		}
		s.DrawString(x, y, fmt.Sprintf("%2d", d), attr)
	}
}
