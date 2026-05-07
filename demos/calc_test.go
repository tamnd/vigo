package demos

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
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

// findCalcButton returns the button labeled label from the
// calculator's children, or nil.
func findCalcButton(c *Calculator, label string) *widget.Button {
	for _, ch := range c.Children() {
		b, ok := ch.(*widget.Button)
		if ok && b.Title == label {
			return b
		}
	}
	return nil
}

func TestCalculatorHasButtonGrid(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	want := []string{"C", "←", "%", "±", "7", "8", "9", "/", "4", "5", "6", "*", "1", "2", "3", "-", "0", ".", "=", "+"}
	for _, lbl := range want {
		if findCalcButton(c, lbl) == nil {
			t.Errorf("missing button %q", lbl)
		}
	}
}

func TestCalculatorButtonClickFiresAction(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	for _, lbl := range []string{"2", "+", "3", "="} {
		b := findCalcButton(c, lbl)
		if b == nil {
			t.Fatalf("missing button %q", lbl)
		}
		b.Press()
	}
	if got := c.Display(); got != "5" {
		t.Fatalf("display=%q want 5 after button-press 2+3=", got)
	}
}

func TestCalculatorClearButtonResets(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "42")
	findCalcButton(c, "C").Press()
	if got := c.Display(); got != "0" {
		t.Fatalf("after C button display=%q want 0", got)
	}
}

func TestCalculatorNegateButtonFlipsSign(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "5")
	findCalcButton(c, "±").Press()
	if got := c.Display(); got != "-5" {
		t.Fatalf("after ± display=%q want -5", got)
	}
	findCalcButton(c, "±").Press()
	if got := c.Display(); got != "5" {
		t.Fatalf("after second ± display=%q want 5", got)
	}
}

func TestCalculatorPercentButtonDivides(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "50")
	findCalcButton(c, "%").Press()
	if got := c.Display(); got != "0.5" {
		t.Fatalf("after %% display=%q want 0.5", got)
	}
}

func TestCalculatorBackspaceButtonTrimsOperand(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	sendRunes(c, "12")
	findCalcButton(c, "←").Press()
	if got := c.Display(); got != "1" {
		t.Fatalf("after backspace display=%q want 1", got)
	}
}

// TestCalculatorButtonClickFromMouse verifies the full click chain:
// a ClassMouseDown event delivered to the calculator's Group routes
// through to the button under the cursor and fires its action via
// the calc sink.
func TestCalculatorButtonClickFromMouse(t *testing.T) {
	c := NewCalculator(CalcDefaultBounds)
	hostFor(c)

	b := findCalcButton(c, "7")
	if b == nil {
		t.Fatal("missing button 7")
	}
	c.HandleEvent(&event.Event{
		What:  event.ClassMouseDown,
		Mouse: event.MouseEvent{X: b.Bounds.X + 1, Y: b.Bounds.Y, Buttons: event.MouseLeft},
	})
	if got := c.Display(); got != "7" {
		t.Fatalf("after click 7 display=%q want 7", got)
	}
}
