package dialog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
)

func TestMessageBoxOKHasOnlyOKButton(t *testing.T) {
	d := MessageBox("hello", MbOK)
	if d.DefaultButton() == nil {
		t.Fatal("OK should be default")
	}
	if d.DefaultButton().Command != event.CmdOk {
		t.Fatalf("default command: %v", d.DefaultButton().Command)
	}
	if d.CancelButton() != nil {
		t.Fatal("MbOK should not register a cancel button")
	}
}

func TestMessageBoxOKCancelHasBoth(t *testing.T) {
	d := MessageBox("are you sure?", MbOKCancel)
	if d.DefaultButton() == nil || d.DefaultButton().Command != event.CmdOk {
		t.Fatalf("default: %+v", d.DefaultButton())
	}
	if d.CancelButton() == nil || d.CancelButton().Command != event.CmdCancel {
		t.Fatalf("cancel: %+v", d.CancelButton())
	}
}

func TestMessageBoxYesNoMakesYesDefaultAndNoCancel(t *testing.T) {
	d := MessageBox("delete?", MbYesNo)
	if d.DefaultButton() == nil || d.DefaultButton().Command != event.CmdYes {
		t.Fatalf("yes default: %+v", d.DefaultButton())
	}
	if d.CancelButton() == nil || d.CancelButton().Command != event.CmdNo {
		t.Fatalf("no cancel: %+v", d.CancelButton())
	}
}

func TestMessageBoxYesNoCancelPrefersCancelOverNo(t *testing.T) {
	d := MessageBox("save?", MbYesNoCancel)
	if d.CancelButton() == nil || d.CancelButton().Command != event.CmdCancel {
		t.Fatalf("cancel button: %+v", d.CancelButton())
	}
}

func TestMessageBoxTitleByKind(t *testing.T) {
	cases := []struct {
		flags Flags
		title string
	}{
		{KindInformation | BtnOK, "Information"},
		{KindWarning | BtnOK, "Warning"},
		{KindError | BtnOK, "Error"},
		{KindConfirmation | BtnOK, "Confirm"},
	}
	for _, tc := range cases {
		got := MessageBox("x", tc.flags).Title
		if got != tc.title {
			t.Errorf("flags %v: got title %q, want %q", tc.flags, got, tc.title)
		}
	}
}

func TestMessageBoxSizingAccommodatesText(t *testing.T) {
	d := MessageBox("a longer message that should set the width", MbOK)
	if d.Bounds.W < 40 {
		t.Fatalf("dialog should clamp to min width: %d", d.Bounds.W)
	}
}

func TestMessageBoxConfirmDialogSnapshot(t *testing.T) {
	host := newPalettedHost()
	d := MessageBox("Discard unsaved\nchanges?", MbYesNoCancel)
	d.SetState(0, false) // make sure it does not look "active" without focus
	host.Insert(d)
	host.SetCurrent(0)

	s := vio.NewSurface(d.Bounds.W, d.Bounds.H)
	d.Draw(s)

	got := strings.Join(s.Snapshot(), "\n")
	goldenPath := filepath.Join("testdata", "confirm_dialog.golden")

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
	}

	raw, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden: %v (re-run with UPDATE_GOLDEN=1 to seed)", err)
	}
	want := strings.ReplaceAll(string(raw), "\r\n", "\n")
	if want != got {
		t.Fatalf("snapshot drift:\nwant:\n%s\ngot:\n%s", want, got)
	}
}
