package menu

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/tamnd/vigo/event"
)

const sampleTOML = `# top-level menus
[[menu]]
title = "File"
hotkey = 0

  [[menu.items]]
  title = "New"
  hotkey = 0
  cmd = 101
  shortcut = "F3"

  [[menu.items]]
  title = "Open"
  hotkey = 0
  cmd = 102
  shortcut = "F4"

  [[menu.items]]
  sep = true

  [[menu.items]]
  title = "Exit"
  hotkey = 1
  cmd = 1
  shortcut = "Alt-X"

[[menu]]
title = "Edit"
hotkey = 0

  [[menu.items]]
  title = "Undo"
  hotkey = 0
  cmd = 200
  shortcut = "Alt-Bksp"
`

func TestLoadTOMLParsesMenus(t *testing.T) {
	menus, err := LoadTOML(strings.NewReader(sampleTOML))
	if err != nil {
		t.Fatalf("LoadTOML: %v", err)
	}
	if len(menus) != 2 {
		t.Fatalf("got %d menus", len(menus))
	}
	if menus[0].Title != "File" {
		t.Fatalf("menu0 title=%q", menus[0].Title)
	}
	if len(menus[0].Children) != 4 {
		t.Fatalf("File children=%d", len(menus[0].Children))
	}
	if menus[0].Children[0].Cmd != 101 {
		t.Fatalf("New cmd=%d", menus[0].Children[0].Cmd)
	}
	if !menus[0].Children[2].Sep {
		t.Fatalf("third row should be sep")
	}
	if menus[0].Children[3].Cmd != event.CmdQuit {
		t.Fatalf("Exit cmd=%d", menus[0].Children[3].Cmd)
	}
	if menus[1].Children[0].Shortcut != "Alt-Bksp" {
		t.Fatalf("Undo shortcut=%q", menus[1].Children[0].Shortcut)
	}
}

func TestLoadTOMLIgnoresCommentsAndUnknownKeys(t *testing.T) {
	src := `
# leading comment
[[menu]]
title = "Test" # inline
hotkey = 0
unknown = "ignored"

  [[menu.items]]
  title = "Row"
  cmd = 5
  weirdkey = 99
`
	menus, err := LoadTOML(strings.NewReader(src))
	if err != nil {
		t.Fatalf("LoadTOML: %v", err)
	}
	if len(menus) != 1 || menus[0].Title != "Test" {
		t.Fatalf("comment/unknown handling broke parse: %+v", menus)
	}
}

func TestLoadTOMLErrors(t *testing.T) {
	cases := []struct {
		name, src string
	}{
		{"items before menu", "[[menu.items]]\ntitle = \"x\"\n"},
		{"unquoted string", "[[menu]]\ntitle = unquoted\n"},
		{"bad cmd", "[[menu]]\ntitle = \"X\"\n[[menu.items]]\ncmd = abc\n"},
		{"missing equals", "[[menu]]\ntitle\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := LoadTOML(strings.NewReader(tc.src)); err == nil {
				t.Fatalf("expected error for %q", tc.src)
			}
		})
	}
}

func TestTOMLRoundTrip(t *testing.T) {
	original := []Item{
		{Title: "File", Hotkey: 0, Children: []Item{
			{Title: "New", Hotkey: 0, Cmd: event.CommandID(101), Shortcut: "F3"},
			{Sep: true},
			{Title: "Exit", Hotkey: 1, Cmd: event.CmdQuit, Shortcut: "Alt-X"},
		}},
		{Title: "Edit", Hotkey: 0, Children: []Item{
			{Title: "Undo", Hotkey: 0, Cmd: event.CommandID(200), Shortcut: "Alt-Bksp"},
		}},
	}

	var buf bytes.Buffer
	if err := SaveTOML(&buf, original); err != nil {
		t.Fatalf("SaveTOML: %v", err)
	}
	round, err := LoadTOML(&buf)
	if err != nil {
		t.Fatalf("LoadTOML round-trip: %v\n--- input ---\n%s", err, buf.String())
	}
	if !reflect.DeepEqual(round, original) {
		t.Fatalf("round-trip mismatch:\nwant %+v\ngot  %+v", original, round)
	}
}
