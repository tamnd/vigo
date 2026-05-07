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
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
	"github.com/tamnd/vigo/window"
)

// Calculator is a four-function calculator window: an accumulator,
// a pending operation, and the operand currently being typed.
// Activation: digits append to the operand, +-*/ apply the pending
// op (if any) and queue a new one, = / Enter runs the pending op
// and clears it, C clears everything. The 4x5 button grid lets the
// user click each action with the mouse; every button is also wired
// to the matching keyboard binding.
type Calculator struct {
	*window.Window

	display *widget.StaticText

	acc     float64 // accumulator (left side of pending op)
	op      rune    // pending op: '+' '-' '*' '/' or 0
	operand string  // operand currently being typed
	primed  bool    // true after = so next digit starts fresh
}

// CalcDefaultBounds is the size that fits the display row plus the
// 4-column, 5-row button grid below it.
//
//nolint:gochecknoglobals // immutable default
var CalcDefaultBounds = vio.R(2, 2, 24, 10)

// calcCmdBase is the offset used to encode each button's action as a
// CommandID. The low bits of the command carry the action rune so
// the Calculator's command sink can decode it without a lookup table.
const calcCmdBase event.CommandID = event.CmdUser + 200

// calcButtonAction returns the CommandID that fires the action keyed
// by r when the button is pressed. Action runes always sit in the
// 7-bit ASCII range, so the cast cannot wrap.
func calcButtonAction(r rune) event.CommandID {
	return calcCmdBase + event.CommandID(r&0x7f)
}

// calcButtons describes the 4-column, 5-row button grid. Each entry
// pairs the rendered glyph with the action rune fed to the calculator
// engine. Action runes mirror the keyboard bindings ('+', '-', '*',
// '/', '=', '.', digits, 'C') with a few extras: '\b' for backspace,
// '%' for percent, '~' for negate.
//
//nolint:gochecknoglobals // immutable layout
var calcButtons = [5][4]struct {
	label  string
	action rune
}{
	{{"C", 'C'}, {"←", '\b'}, {"%", '%'}, {"±", '~'}},
	{{"7", '7'}, {"8", '8'}, {"9", '9'}, {"/", '/'}},
	{{"4", '4'}, {"5", '5'}, {"6", '6'}, {"*", '*'}},
	{{"1", '1'}, {"2", '2'}, {"3", '3'}, {"-", '-'}},
	{{"0", '0'}, {".", '.'}, {"=", '='}, {"+", '+'}},
}

// NewCalculator returns a Calculator window sized to bounds. The
// button grid is laid out from CalcDefaultBounds; if bounds is taller
// or wider, the grid stays anchored to the top-left of the client.
func NewCalculator(bounds vio.Rect) *Calculator {
	w := window.New(bounds, "Calculator", window.FlagMove|window.FlagClose)
	c := &Calculator{Window: w}

	client := w.ClientRect()
	c.display = widget.NewStaticText(
		vio.R(client.X+1, client.Y+1, client.W-2, 1),
		"",
	)
	w.Insert(c.display)
	c.refresh()

	for r := range calcButtons {
		for col := range calcButtons[r] {
			spec := calcButtons[r][col]
			b := widget.NewButton(
				vio.R(client.X+1+col*5, client.Y+3+r, 4, 1),
				spec.label,
				calcButtonAction(spec.action),
				widget.BfNormal,
			)
			w.Insert(b)
		}
	}

	w.Insert(newCalcSink(c))
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
	return c.applyAction(k.Rune)
}

// applyAction runs the calculator action keyed by r. It is the shared
// entry point for keystrokes and button clicks.
func (c *Calculator) applyAction(r rune) bool {
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
	case r == '\b':
		c.backspace()
	case r == '%':
		c.percent()
	case r == '~':
		c.negate()
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

// percent divides the operand (or the accumulator if no operand is
// being typed) by 100, matching the classic calculator behavior.
func (c *Calculator) percent() {
	if c.operand != "" {
		v, err := strconv.ParseFloat(c.operand, 64)
		if err != nil {
			return
		}
		c.operand = formatNumber(v / 100)
	} else {
		c.acc /= 100
	}
	c.refresh()
}

// negate flips the sign of the operand or the accumulator.
func (c *Calculator) negate() {
	if c.operand != "" {
		if c.operand[0] == '-' {
			c.operand = c.operand[1:]
		} else {
			c.operand = "-" + c.operand
		}
	} else {
		c.acc = -c.acc
	}
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

// calcSink is a hidden post-process child that captures the command
// fired by each Calculator button. The button's owner is the embedded
// Group, not Calculator, so without a sink the command would walk
// past Calculator.HandleEvent's override and leak up to the desktop.
type calcSink struct {
	*view.View
	calc *Calculator
}

func newCalcSink(c *Calculator) *calcSink {
	v := view.NewView(vio.Rect{})
	v.Options |= view.OptPostProcess
	v.EventMask = event.ClassCommand
	return &calcSink{View: v, calc: c}
}

func (s *calcSink) Draw(*vio.Surface) {}

func (s *calcSink) HandleEvent(e *event.Event) {
	if e.What != event.ClassCommand {
		return
	}
	cmd := e.Msg.Command
	if cmd < calcCmdBase || cmd > calcCmdBase+0x7f {
		return
	}
	if s.calc.applyAction(rune(cmd - calcCmdBase)) {
		e.Clear()
	}
}
