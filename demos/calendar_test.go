package demos

import (
	"testing"
	"time"

	"github.com/tamnd/vigo/event"
)

func dateAt(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestCalendarStartsOnGivenDate(t *testing.T) {
	c := NewCalendarAt(CalDefaultBounds, dateAt(2026, time.May, 7))
	hostFor(c)

	if got := c.Cursor(); !got.Equal(dateAt(2026, time.May, 7)) {
		t.Fatalf("cursor=%v want 2026-05-07", got)
	}
}

func TestCalendarPgDnAdvancesMonth(t *testing.T) {
	c := NewCalendarAt(CalDefaultBounds, dateAt(2026, time.May, 7))
	hostFor(c)

	c.HandleEvent(keyOf(event.KeyPgDn))
	if got := c.Cursor(); got.Month() != time.June || got.Day() != 7 {
		t.Fatalf("PgDn cursor=%v want June 7", got)
	}
}

func TestCalendarPgUpStepsBack(t *testing.T) {
	c := NewCalendarAt(CalDefaultBounds, dateAt(2026, time.January, 15))
	hostFor(c)

	c.HandleEvent(keyOf(event.KeyPgUp))
	if got := c.Cursor(); got.Year() != 2025 || got.Month() != time.December {
		t.Fatalf("PgUp cursor=%v want 2025-12-15", got)
	}
}

func TestCalendarMonthOverflowClampsDay(t *testing.T) {
	c := NewCalendarAt(CalDefaultBounds, dateAt(2026, time.January, 31))
	hostFor(c)

	c.HandleEvent(keyOf(event.KeyPgDn))
	got := c.Cursor()
	if got.Month() != time.February {
		t.Fatalf("month=%v want February", got.Month())
	}
	if got.Day() != 28 {
		t.Fatalf("day=%d want 28 (clamped)", got.Day())
	}
}

func TestCalendarArrowsMoveDay(t *testing.T) {
	c := NewCalendarAt(CalDefaultBounds, dateAt(2026, time.May, 7))
	hostFor(c)

	c.HandleEvent(keyOf(event.KeyArrowRight))
	if c.Cursor().Day() != 8 {
		t.Fatalf("right day=%d want 8", c.Cursor().Day())
	}
	c.HandleEvent(keyOf(event.KeyArrowDown))
	if c.Cursor().Day() != 15 {
		t.Fatalf("down day=%d want 15", c.Cursor().Day())
	}
	c.HandleEvent(keyOf(event.KeyArrowLeft))
	c.HandleEvent(keyOf(event.KeyArrowUp))
	if c.Cursor().Day() != 7 {
		t.Fatalf("back to start day=%d want 7", c.Cursor().Day())
	}
}

func TestCalendarHomeJumpsToFirstOfMonth(t *testing.T) {
	c := NewCalendarAt(CalDefaultBounds, dateAt(2026, time.May, 17))
	hostFor(c)

	c.HandleEvent(keyOf(event.KeyHome))
	if got := c.Cursor(); got.Day() != 1 {
		t.Fatalf("Home day=%d want 1", got.Day())
	}
}

func TestCalendarUnknownKeyFallsThrough(t *testing.T) {
	c := NewCalendarAt(CalDefaultBounds, dateAt(2026, time.May, 7))
	hostFor(c)

	before := c.Cursor()
	c.HandleEvent(keyRune('q'))
	if !c.Cursor().Equal(before) {
		t.Fatalf("unrelated rune changed cursor")
	}
}
