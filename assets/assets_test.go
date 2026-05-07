package assets

import (
	"bytes"
	"testing"

	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/menu"
)

func TestMainMenuLoadsAndContainsExpectedTopLevels(t *testing.T) {
	items, err := menu.LoadTOML(bytes.NewReader(MainMenu()))
	if err != nil {
		t.Fatalf("LoadTOML: %v", err)
	}
	want := []string{"≡", "File", "Edit", "Window", "Help"}
	if len(items) != len(want) {
		t.Fatalf("top-level menus=%d want %d (%v)", len(items), len(want), titlesOf(items))
	}
	for i, w := range want {
		if items[i].Title != w {
			t.Fatalf("menu[%d]=%q want %q", i, items[i].Title, w)
		}
	}
}

func TestMainMenuFileExitFiresCmdQuit(t *testing.T) {
	items, err := menu.LoadTOML(bytes.NewReader(MainMenu()))
	if err != nil {
		t.Fatalf("LoadTOML: %v", err)
	}
	file := items[1] // File
	last := file.Children[len(file.Children)-1]
	if last.Title != "Exit" {
		t.Fatalf("last File row=%q want Exit", last.Title)
	}
	if last.Cmd != event.CmdQuit {
		t.Fatalf("Exit cmd=%d want CmdQuit (%d)", last.Cmd, event.CmdQuit)
	}
}

func TestMainMenuHelpAboutFiresCmdHelp(t *testing.T) {
	items, err := menu.LoadTOML(bytes.NewReader(MainMenu()))
	if err != nil {
		t.Fatalf("LoadTOML: %v", err)
	}
	help := items[len(items)-1] // Help is the last menu
	if help.Title != "Help" {
		t.Fatalf("last menu=%q want Help", help.Title)
	}
	if len(help.Children) == 0 || help.Children[0].Cmd != event.CmdHelp {
		t.Fatalf("first Help row=%+v want CmdHelp", help.Children)
	}
}

func titlesOf(items []menu.Item) []string {
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.Title
	}
	return out
}
