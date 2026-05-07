package dialog

import (
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
	"github.com/tamnd/vigo/window"
)

func newPalettedHost() *view.Group {
	g := view.NewGroup(vio.R(0, 0, 80, 24))
	g.Palette = vio.BorlandClassic
	return g
}

func TestNewSetsDialogFrameSlots(t *testing.T) {
	d := New(vio.R(0, 0, 30, 8), "Confirm")
	f := d.Frame()
	if f.PassiveSlot != 15 || f.ActiveSlot != 16 {
		t.Fatalf("frame slots want 15/16, got %d/%d", f.PassiveSlot, f.ActiveSlot)
	}
	if d.Flags()&window.FlagZoom != 0 {
		t.Fatal("dialog should not enable FlagZoom")
	}
	if d.Flags()&window.FlagGrow != 0 {
		t.Fatal("dialog should not enable FlagGrow")
	}
}

func TestPlaceButtonsCentersRow(t *testing.T) {
	host := newPalettedHost()
	d := New(vio.R(0, 0, 30, 8), "Confirm")
	host.Insert(d)

	ok := widget.NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, widget.BfDefault)
	cancel := widget.NewButton(vio.R(0, 0, 10, 1), "Cancel", event.CmdCancel, widget.BfNormal)
	d.PlaceButtons(ok, cancel)

	client := d.ClientRect()
	if ok.Bounds.Y != client.Bottom()-2 {
		t.Fatalf("ok Y: %d want %d", ok.Bounds.Y, client.Bottom()-2)
	}
	if cancel.Bounds.X != ok.Bounds.X+ok.Bounds.W+1 {
		t.Fatalf("buttons not adjacent: ok=%+v cancel=%+v", ok.Bounds, cancel.Bounds)
	}
	leftPad := ok.Bounds.X - client.X
	rightPad := client.Right() - cancel.Bounds.Right()
	if abs(leftPad-rightPad) > 1 {
		t.Fatalf("button row not centered: leftPad=%d rightPad=%d", leftPad, rightPad)
	}
}

func TestEnterPressesDefaultButton(t *testing.T) {
	host := newPalettedHost()
	d := New(vio.R(0, 0, 30, 8), "Confirm")
	host.Insert(d)
	ok := widget.NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, widget.BfDefault)
	d.PlaceButtons(ok)
	d.SetDefaultButton(ok)

	d.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyEnter},
	})
	if d.Result() != event.CmdOk {
		t.Fatalf("Enter should set Result=CmdOk, got %v", d.Result())
	}
}

func TestEscPressesCancelButton(t *testing.T) {
	host := newPalettedHost()
	d := New(vio.R(0, 0, 30, 8), "Confirm")
	host.Insert(d)
	cancel := widget.NewButton(vio.R(0, 0, 10, 1), "Cancel", event.CmdCancel, widget.BfNormal)
	d.PlaceButtons(cancel)
	d.SetCancelButton(cancel)

	d.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyEsc},
	})
	if d.Result() != event.CmdCancel {
		t.Fatalf("Esc should set Result=CmdCancel, got %v", d.Result())
	}
}

func TestEndModalAndResetResult(t *testing.T) {
	d := New(vio.R(0, 0, 30, 8), "Confirm")
	d.EndModal(event.CmdYes)
	if d.Result() != event.CmdYes {
		t.Fatalf("EndModal should set result, got %v", d.Result())
	}
	d.Reset()
	if d.Result() != event.CmdNone {
		t.Fatalf("Reset should clear result, got %v", d.Result())
	}
}

func TestDisabledDefaultButtonIgnoresEnter(t *testing.T) {
	host := newPalettedHost()
	d := New(vio.R(0, 0, 30, 8), "Confirm")
	host.Insert(d)
	ok := widget.NewButton(vio.R(0, 0, 10, 1), "OK", event.CmdOk, widget.BfDefault)
	ok.SetState(view.StateDisabled, true)
	d.PlaceButtons(ok)
	d.SetDefaultButton(ok)

	d.HandleEvent(&event.Event{
		What: event.ClassKey,
		Key:  event.KeyEvent{Key: event.KeyEnter},
	})
	if d.Result() != event.CmdNone {
		t.Fatalf("disabled default button should not fire on Enter, result=%v", d.Result())
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
