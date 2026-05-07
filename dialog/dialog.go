// Package dialog ships the modal Dialog view: a Window dressed in the
// gray dialog palette, with helpers for laying out a row of buttons
// along the bottom edge.
//
// A Dialog is just a Window with FlagClose|FlagMove (no zoom, no grow)
// and the frame swapped to the dialog palette slots. Insert child
// widgets the same way you would with a Window. Phase 3 wires Enter
// and Esc shortcuts; the modal sub-loop arrives with Application.ExecView.
package dialog

import (
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
	"github.com/tamnd/vigo/window"
)

// buttonRowMargin is the gap between the bottom button row and the
// frame's bottom edge.
const buttonRowMargin = 2

// Dialog is a modal Window themed for forms and prompts. Embed
// *window.Window so all the framing, drag, and event-dispatch logic is
// inherited; Dialog only customizes palette slots and adds a button
// placement helper.
type Dialog struct {
	*window.Window

	defaultBtn *widget.Button
	cancelBtn  *widget.Button
}

// New returns a Dialog sized to bounds. The frame uses dialog slots
// (15 inactive, 16 active), the body fills with dialog text color, and
// FlagClose|FlagMove are enabled so the user can drag and dismiss.
func New(bounds vio.Rect, title string) *Dialog {
	w := window.New(bounds, title, window.FlagMove|window.FlagClose)
	w.PaletteIndex = widget.SlotDialogText
	w.Frame().PassiveSlot = 15
	w.Frame().ActiveSlot = 16
	return &Dialog{Window: w}
}

// SetDefaultButton remembers b as the dialog's default action so
// pressing Enter anywhere in the dialog fires it. Pass nil to clear.
func (d *Dialog) SetDefaultButton(b *widget.Button) { d.defaultBtn = b }

// SetCancelButton remembers b as the dialog's cancel action so
// pressing Esc anywhere in the dialog fires it. Pass nil to clear.
func (d *Dialog) SetCancelButton(b *widget.Button) { d.cancelBtn = b }

// DefaultButton returns the button registered via SetDefaultButton, or nil.
func (d *Dialog) DefaultButton() *widget.Button { return d.defaultBtn }

// CancelButton returns the button registered via SetCancelButton, or nil.
func (d *Dialog) CancelButton() *widget.Button { return d.cancelBtn }

// PlaceButtons inserts the given buttons in a centered row along the
// bottom of the client rectangle. Buttons keep the size set by the
// caller; PlaceButtons only repositions them. Useful for the common
// OK/Cancel cluster a host doesn't want to lay out by hand.
func (d *Dialog) PlaceButtons(buttons ...*widget.Button) {
	if len(buttons) == 0 {
		return
	}
	client := d.ClientRect()
	total := 0
	for i, b := range buttons {
		total += b.Bounds.W
		if i > 0 {
			total++ // single-cell gap between buttons
		}
	}
	x := client.X + (client.W-total)/2
	y := client.Bottom() - buttonRowMargin
	for _, b := range buttons {
		nb := b.Bounds
		nb.X = x
		nb.Y = y
		b.ChangeBounds(nb)
		d.Insert(b)
		x += nb.W + 1
	}
}

// HandleEvent intercepts dialog-level shortcuts before falling back to
// the Window dispatch. Enter fires the registered default button; Esc
// fires the registered cancel button. A disabled button is skipped so
// the rest of the dispatch chain still gets the keystroke.
func (d *Dialog) HandleEvent(e *event.Event) {
	if e.What == event.ClassKey {
		switch e.Key.Key {
		case event.KeyEnter:
			if d.defaultBtn != nil && !d.defaultBtn.HasState(view.StateDisabled) {
				d.defaultBtn.Press()
				e.Clear()
				return
			}
		case event.KeyEsc:
			if d.cancelBtn != nil && !d.cancelBtn.HasState(view.StateDisabled) {
				d.cancelBtn.Press()
				e.Clear()
				return
			}
		}
	}
	d.Window.HandleEvent(e)
}
