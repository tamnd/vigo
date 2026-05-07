// Package assets exposes the static files vigo ships in-binary:
// the default menu tree and any future palette or shortcut config.
// The TOML and other source files live next to this package so they
// can be edited by hand and embedded by relative path.
package assets

import _ "embed"

//go:embed menus/main.toml
var mainMenu []byte

// MainMenu returns the bytes of the canonical IDE menu tree
// (assets/menus/main.toml). Hosts feed it to menu.LoadTOML to
// populate the menu bar.
func MainMenu() []byte { return mainMenu }
