package demos

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// keyRune builds a printable-key event for tests.
func keyRune(r rune) *event.Event {
	return &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: event.KeyRune, Rune: r}}
}

// keyOf builds a non-rune key event for tests.
func keyOf(k event.Key) *event.Event {
	return &event.Event{What: event.ClassKey, Key: event.KeyEvent{Key: k}}
}

// hostFor returns a paletted Group hosting v so MapColor works.
func hostFor(v view.Viewer) *view.Group {
	g := view.NewGroup(vio.R(0, 0, 80, 24))
	g.Palette = vio.BorlandClassic
	g.Insert(v)
	return g
}

func sendRunes(c *Calculator, runes string) {
	for _, r := range runes {
		c.HandleEvent(keyRune(r))
	}
}

func TestCalculatorAddTwoIntegers(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "2+3=")
	if got := c.Display(); got != "5" {
		t.Fatalf("2+3 display=%q want 5", got)
	}
}

func TestCalculatorChainOps(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	// 2 + 3 * 4 evaluates left-to-right: (2+3)*4 = 20
	sendRunes(c, "2+3*4=")
	if got := c.Display(); got != "20" {
		t.Fatalf("2+3*4 display=%q want 20", got)
	}
}

func TestCalculatorDivisionAndDecimal(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "5/2=")
	if got := c.Display(); got != "2.5" {
		t.Fatalf("5/2 display=%q want 2.5", got)
	}
}

func TestCalculatorDivideByZeroClears(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "8/0=")
	if got := c.Display(); got != "0" {
		t.Fatalf("8/0 display=%q want 0 after clear", got)
	}
}

func TestCalculatorClearResets(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "12+7")
	sendRunes(c, "C")
	if got := c.Display(); got != "0" {
		t.Fatalf("after C display=%q want 0", got)
	}
}

func TestCalculatorBackspaceTrimsOperand(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "12")
	c.HandleEvent(keyOf(event.KeyBackspace))
	if got := c.Display(); got != "1" {
		t.Fatalf("backspace display=%q want 1", got)
	}
}

func TestCalculatorPendingOpShownInDisplay(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "9+")
	if got := c.Display(); got != "9 +" {
		t.Fatalf("after op display=%q want '9 +'", got)
	}
}

func TestCalculatorEnterEqualsRunsPending(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "10-4")
	c.HandleEvent(keyOf(event.KeyEnter))
	if got := c.Display(); got != "6" {
		t.Fatalf("Enter as = display=%q want 6", got)
	}
}

func TestCalculatorDoubleDotIgnored(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "1.2.3")
	if got := c.Display(); got != "1.23" {
		t.Fatalf("double dot display=%q want 1.23", got)
	}
}
