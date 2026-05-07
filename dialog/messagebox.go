package dialog

import (
	"strings"
	"unicode/utf8"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
)

// Flags selects the title and the buttons of a MessageBox. Combine one
// kind bit with one or more button bits, or pass a stock combo such
// as MbOK or MbYesNo.
type Flags uint16

// Kind bits.
const (
	KindInformation Flags = 1 << iota
	KindWarning
	KindError
	KindConfirmation
)

// Button bits.
const (
	BtnOK Flags = 1 << (iota + 8)
	BtnCancel
	BtnYes
	BtnNo
)

// Stock button combinations.
const (
	MbOK          = KindInformation | BtnOK
	MbOKCancel    = KindConfirmation | BtnOK | BtnCancel
	MbYesNo       = KindConfirmation | BtnYes | BtnNo
	MbYesNoCancel = KindConfirmation | BtnYes | BtnNo | BtnCancel
)

const (
	mbButtonW = 10
	mbButtonH = 1
	mbMinW    = 40
	mbMaxW    = 60
	mbPadX    = 3
	mbPadY    = 1
	mbBtnRow  = 2 // matches buttonRowMargin
)

// MessageBox returns a Dialog that displays text with the title and
// buttons selected by flags. The dialog is not yet inserted anywhere;
// pass it to Application.ExecView (or similar) to show it modally.
//
// Lines in text are split on '\n'. Long lines are clipped at mbMaxW;
// the helper does not word-wrap. The first present button in
// (OK, Yes, No, Cancel) becomes the default; Cancel (or No when no
// Cancel exists) becomes the cancel button.
func MessageBox(text string, flags Flags) *Dialog {
	lines := strings.Split(text, "\n")
	textW := 0
	for _, line := range lines {
		textW = max(textW, utf8.RuneCountInString(line))
	}

	width := textW + 2*mbPadX + 2 // +2 for the frame
	width = max(width, mbMinW)
	width = min(width, mbMaxW)
	height := len(lines) + 2*mbPadY + mbBtnRow + 2 // +2 for frame
	bounds := vio.Rect{X: 0, Y: 0, W: width, H: height}

	d := New(bounds, messageBoxTitle(flags))

	body := widget.NewStaticText(
		vio.Rect{
			X: mbPadX,
			Y: mbPadY + 1, // +1 to skip the top border
			W: width - 2*mbPadX,
			H: len(lines),
		},
		text,
	)
	d.Insert(body)

	buttons := messageBoxButtons(flags)
	if len(buttons) > 0 {
		d.PlaceButtons(buttons...)
		applyDefaultAndCancel(d, buttons, flags)
	}

	return d
}

func messageBoxTitle(flags Flags) string {
	switch {
	case flags&KindError != 0:
		return "Error"
	case flags&KindWarning != 0:
		return "Warning"
	case flags&KindConfirmation != 0:
		return "Confirm"
	default:
		return "Information"
	}
}

func messageBoxButtons(flags Flags) []*widget.Button {
	var bs []*widget.Button
	bw := vio.Rect{W: mbButtonW, H: mbButtonH}
	if flags&BtnOK != 0 {
		bs = append(bs, widget.NewButton(bw, "~O~K", event.CmdOk, widget.BfDefault))
	}
	if flags&BtnYes != 0 {
		fl := widget.BfNormal
		if flags&BtnOK == 0 {
			fl = widget.BfDefault
		}
		bs = append(bs, widget.NewButton(bw, "~Y~es", event.CmdYes, fl))
	}
	if flags&BtnNo != 0 {
		bs = append(bs, widget.NewButton(bw, "~N~o", event.CmdNo, widget.BfNormal))
	}
	if flags&BtnCancel != 0 {
		bs = append(bs, widget.NewButton(bw, "Cancel", event.CmdCancel, widget.BfNormal))
	}
	return bs
}

func applyDefaultAndCancel(d *Dialog, buttons []*widget.Button, flags Flags) {
	for _, b := range buttons {
		if b.Flags&widget.BfDefault != 0 {
			d.SetDefaultButton(b)
			break
		}
	}
	switch {
	case flags&BtnCancel != 0:
		for _, b := range buttons {
			if b.Command == event.CmdCancel {
				d.SetCancelButton(b)
				return
			}
		}
	case flags&BtnNo != 0:
		for _, b := range buttons {
			if b.Command == event.CmdNo {
				d.SetCancelButton(b)
				return
			}
		}
	}
}
