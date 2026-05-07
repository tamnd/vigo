# 1597. vigo overview

vigo is a port of Borland's Turbo Vision IDE, written in Go and run in
a terminal. It targets Go programmers who want a self-hosted IDE with
the look and feel of Turbo Pascal 6.0 / Borland C++ 3.x: blue desktop,
double-line borders, drop shadows, F-key shortcuts, pull-down menus,
modal dialogs. Code intelligence comes from gopls; debugging comes
from dlv. The whole thing ships as a single static binary.

The intended workflow is small:

```
go install github.com/tamnd/vigo/cmd/vigo@latest
vigo .
```

That should land the user in a Borland-style desktop, with an editor,
a project tree, build/run/test wired up, gopls-driven completion and
diagnostics, and dlv-driven debugging.

## Goals

1. Match Turbo Vision honestly. Same palette, same glyphs, same keys,
   same modal dialog mechanics, same drop shadows. Mouse works but
   the keyboard is the primary interface.
2. Pure Go. No CGO. One binary per platform. Fast cold start.
3. Faithful object model. View, Group, Window, Dialog, Application,
   Desktop, MenuBar, StatusLine, command system, palettes, streams.
   Class-based inheritance becomes embedding plus small interfaces in
   Go, but the user-visible model is identical.
4. First-class Go IDE. Editor with syntax highlighting, gopls for
   navigation/completion/diagnostics/refactor, `go build|test|run`
   integration, dlv DAP debugger with breakpoints, watches, stepping.
5. Cross-platform. Linux, macOS, Windows (Windows Terminal),
   FreeBSD. Unicode by default, with a CP437 fallback table for
   terminals that can't render Unicode box-drawing glyphs.
6. Hackable. Plain-text TOML for menus, palettes, key maps. Starlark
   plugin host arrives in v0.7.

## Non-goals

- Not a graphical IDE. vigo lives in a character grid.
- Not a clone of VS Code, GoLand, or Neovim. Mouse is supported,
  density and aesthetics follow Borland 1990.
- Not a general LSP client framework. gopls is the supported backend.
- No bundled package manager, no AI features in 1.0. Both are
  reasonable plugin candidates after 1.0.

## Why now

- Terminal UIs are healthy again: tcell, ratatui, bubbletea, magiblot's
  tvision, Zellij. The toolchain is good.
- gopls and `dlv dap` make a self-hosted Go IDE feasible without
  reinventing parsers or debuggers.
- Existing Turbo Vision ports cover C++ (magiblot/tvision) and Rust
  (turbo-vision-4-rust). There is no Go port.

## Reference layout

```
+--[=]--File--Edit--Search--Run--Compile--Debug--Project--Options--Window--Help--+
|                                                                                |
|   ##[*]= main.go ====================================================== 1 ##   |
|   #  package main                                                          #   |
|   #                                                                        #   |
|   #  import "fmt"                                                          #   |
|   #                                                                        #   |
|   #  func main() {                                                         #   |
|   #      fmt.Println("hello, vigo")                                        #   |
|   #  }                                                                     #   |
|   #                                                                        #   |
|   ##======================================== 6:1 =============== INS =====##   |
|                                                                                |
+--F1 Help--F2 Save--F3 Open--F5 Zoom--F6 Switch--F9 Make--F10 Menu--Alt-X Exit--+
```

## References

- Borland, *Turbo Vision Programmer's Guide*, 1990 (bitsavers PDF).
- magiblot/tvision, modern C++ port with Unicode.
- aovestdipaperino/turbo-vision-4-rust, Rust port (1.0 in 2025).
- gdamore/tcell and rivo/tview, Go terminal foundation.
- gopls (`golang.org/x/tools/gopls`) and Delve (`dlv dap`).

See also:

- `1598_vigo_roadmap.md`, phased delivery plan.
- `1599_vigo_architecture.md`, layered design.
- `1600`..`1607`, per-version specifications.
