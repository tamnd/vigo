# vigo

[![CI](https://github.com/tamnd/vigo/actions/workflows/ci.yml/badge.svg)](https://github.com/tamnd/vigo/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tamnd/vigo.svg)](https://pkg.go.dev/github.com/tamnd/vigo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A Turbo Vision style IDE for Go, written in Go, run in a terminal.

vigo brings back the Borland Turbo Pascal 6.0 / Borland C++ 3.x look:
blue desktop, double-line window borders, drop shadows, F-key shortcuts,
pull-down menus, modal dialogs. Underneath it uses gopls for code
intelligence and dlv for debugging. One static binary, no CGO.

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

## Status

Pre-alpha. The project is being bootstrapped toward v0.1 (foundation).
The full plan, version by version, is in `docs/roadmap.md` and in the
matching specs at `~/notes/Spec/1500/1597..1607`.

| Version | Theme               | Status      |
|---------|---------------------|-------------|
| 0.1     | Foundation          | in progress |
| 0.2     | TUI core            | planned     |
| 0.3     | Editor              | planned     |
| 0.4     | Project & toolchain | planned     |
| 0.5     | gopls / LSP         | planned     |
| 0.6     | Delve / DAP         | planned     |
| 0.7     | Polish              | planned     |
| 1.0     | Release             | planned     |

## Why

There is no Go port of Turbo Vision. There are good Rust and C++ ports,
and there are good Go terminal UI libraries (tcell, tview, bubbletea),
but nothing combines the Turbo Vision object model with Go's editor
tooling. vigo is that combination, scoped narrowly to a Go IDE that fits
in a terminal.

Goals:

- Reproduce the Turbo Vision look and feel honestly. Same palette, same
  glyphs, same keys, same modal dialogs, same drop shadows.
- Stay in pure Go. No CGO, one binary per platform, fast cold start.
- Use gopls and dlv as supervised subprocesses. Don't reimplement
  parsers, type checkers, or debuggers.
- Keep the UI single-threaded. Background work posts events back to
  the loop on a channel, never mutates view state directly.

Non-goals:

- Not a clone of VS Code, GoLand, or Neovim. Mouse is supported but the
  keyboard is the point.
- Not a multi-language IDE. Go is the supported target. Other languages
  may work via gopls's neighbours, but they are not on the roadmap.
- No bundled package manager and no AI features in 1.0. Both are
  reasonable plugin candidates after 1.0.

## Quickstart

Once published:

```sh
go install github.com/tamnd/vigo/cmd/vigo@latest
vigo .
```

For now, build from source:

```sh
git clone https://github.com/tamnd/vigo
cd vigo
go run ./cmd/vigo
```

Press `Alt-X` (or `F10`, then `File`, then `Exit`) to quit.

## Project layout

```
vigo/
  cmd/vigo/        binary entrypoint
  vio/             tcell-backed screen, palette, surface, glyphs
  event/           typed events, keyboard, mouse, commands
  view/            View, Group, focus, palette indirection
  app/             Application, Desktop, Background, main loop
  menu/            MenuBar, StatusLine
  docs/            roadmap, architecture, overview
  .github/         CI, release workflows
  CHANGELOG.md
  LICENSE          MIT
  README.md
```

## Architecture

vigo is a thin layer over [`gdamore/tcell/v2`](https://github.com/gdamore/tcell)
that implements Turbo Vision's object model in Go. The translation uses
composition with embedded `*View` and small interfaces (`Drawable`,
`Eventer`, `Sizer`) instead of class-based inheritance. Drawing happens
once per frame onto an in-memory `Surface`; tcell handles the diff to
the terminal. The event loop is single-threaded for UI; gopls and dlv
run as subprocesses and post results back through a channel that the
main loop reads with `select`.

Long version: `docs/architecture.md`.

## References

- Borland, Turbo Vision Programmer's Guide, 1990
  ([Internet Archive](https://archive.org/details/bitsavers_borlandturVersion6.0TurboVision1990_16007263)).
- [magiblot/tvision](https://github.com/magiblot/tvision), modern C++ port.
- [aovestdipaperino/turbo-vision-4-rust](https://github.com/aovestdipaperino/turbo-vision-4-rust), Rust port.
- [`gdamore/tcell`](https://github.com/gdamore/tcell), terminal cell library.
- [`golang.org/x/tools/gopls`](https://pkg.go.dev/golang.org/x/tools/gopls).
- [`go-delve/delve`](https://github.com/go-delve/delve) (`dlv dap`).

## Contributing

Patches and issues are welcome. Discuss bigger design changes in an
issue first; the framework is small and opinionated, and a wrong move
in the foundation propagates everywhere. `CONTRIBUTING.md` will land
with v0.2.

## License

[MIT](LICENSE), Copyright (c) 2026 Duc-Tam Nguyen.
