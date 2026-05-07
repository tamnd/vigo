// Package demos collects the classic Turbo Vision sample tools the
// IDE shipped with: Calculator, Calendar, AsciiTable, and Puzzle.
// They exist to exercise the framework end-to-end and to pin a few
// useful golden screens.
package demos

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
	"github.com/tamnd/vigo/window"
)

// Calculator is a four-function calculator window: an accumulator,
// a pending operation, and the operand currently being typed.
// Activation: digits append to the operand, +-*/ apply the pending
// op (if any) and queue a new one, = / Enter runs the pending op
// and clears it, C clears everything.
type Calculator struct {
	*window.Window

	display *widget.StaticText

	acc     float64 // accumulator (left side of pending op)
	op      rune    // pending op: '+' '-' '*' '/' or 0
	operand string  // operand currently being typed
	primed  bool    // true after = so next digit starts fresh
}

// CalcDefaultBounds is a reasonable default size for the calculator
// window. Hosts can resize after construction via ChangeBounds.
//
//nolint:gochecknoglobals // immutable default
var CalcDefaultBounds = vio.R(2, 2, 24, 6)

// NewCalculator returns a Calculator window sized to bounds.
func NewCalculator(bounds vio.Rect) *Calculator {
	w := window.New(bounds, "Calculator", window.FlagMove|window.FlagClose)
	c := &Calculator{Window: w}

	client := w.ClientRect()
	c.display = widget.NewStaticText(
		vio.R(client.X+1, client.Y+1, client.W-2, 1),
		"0",
	)
	w.Insert(c.display)
	return c
}

// Display returns the current display string. Useful in tests.
func (c *Calculator) Display() string { return c.display.Text }

// HandleEvent intercepts key events for digits, operators, =, and C
// before delegating the rest to the Window dispatch.
func (c *Calculator) HandleEvent(e *event.Event) {
	if e.What == event.ClassKey && c.handleKey(e) {
		e.Clear()
		return
	}
	c.Window.HandleEvent(e)
}

func (c *Calculator) handleKey(e *event.Event) bool {
	k := e.Key
	if k.Key == event.KeyEnter {
		c.equals()
		return true
	}
	if k.Key == event.KeyBackspace {
		c.backspace()
		return true
	}
	if k.Key != event.KeyRune {
		return false
	}
	r := k.Rune
	switch {
	case unicode.IsDigit(r):
		c.appendDigit(r)
	case r == '.':
		c.appendDot()
	case r == '+' || r == '-' || r == '*' || r == '/':
		c.applyOp(r)
	case r == '=':
		c.equals()
	case r == 'c' || r == 'C':
		c.clear()
	default:
		return false
	}
	return true
}

func (c *Calculator) appendDigit(r rune) {
	if c.primed {
		c.operand = ""
		c.primed = false
	}
	c.operand += string(r)
	c.refresh()
}

func (c *Calculator) appendDot() {
	if c.primed {
		c.operand = ""
		c.primed = false
	}
	if c.operand == "" {
		c.operand = "0"
	}
	for _, r := range c.operand {
		if r == '.' {
			return
		}
	}
	c.operand += "."
	c.refresh()
}

func (c *Calculator) backspace() {
	if c.operand == "" {
		return
	}
	c.operand = c.operand[:len(c.operand)-1]
	c.refresh()
}

func (c *Calculator) applyOp(r rune) {
	c.commit() // fold any pending op
	c.op = r
	c.operand = ""
	c.primed = false
	c.refresh()
}

func (c *Calculator) equals() {
	c.commit()
	c.op = 0
	c.primed = true
	c.refresh()
}

func (c *Calculator) commit() {
	if c.operand == "" {
		// no new operand: just adopt acc as operand
		return
	}
	v, err := strconv.ParseFloat(c.operand, 64)
	if err != nil {
		c.clear()
		return
	}
	if c.op == 0 {
		c.acc = v
		c.operand = ""
		return
	}
	switch c.op {
	case '+':
		c.acc += v
	case '-':
		c.acc -= v
	case '*':
		c.acc *= v
	case '/':
		if v == 0 {
			c.clear()
			return
		}
		c.acc /= v
	}
	c.operand = ""
}

func (c *Calculator) clear() {
	c.acc = 0
	c.op = 0
	c.operand = ""
	c.primed = false
	c.refresh()
}

// refresh syncs the display text from the calculator's state. The
// rule: if the user is typing an operand, show that; otherwise show
// the accumulator. The pending op is appended as a hint.
func (c *Calculator) refresh() {
	var text string
	if c.operand != "" {
		text = c.operand
	} else {
		text = formatNumber(c.acc)
	}
	if c.op != 0 && !c.primed {
		text = fmt.Sprintf("%s %c", text, c.op)
	}
	c.display.Text = text
}

func formatNumber(v float64) string {
	if v == float64(int64(v)) {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}
