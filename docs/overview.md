# 1597 — Vigo Overview

`vigo` is a 100% faithful, modernized port of Borland's Turbo Vision integrated
development environment, reimagined as a self-hosted Go IDE written in Go. It
runs in any modern terminal, looks and feels like the Turbo Pascal 6.0 / Borland
C++ 3.x IDE — blue desktop, double-line borders, F-key shortcuts, mouse-aware
menus and dialogs — and provides a complete edit/build/debug workflow for Go
programmers using `gopls` and `dlv` underneath.

## Vision

> The fastest, most beautiful, most keyboard-driven Go IDE that fits in a
> terminal — pixel-faithful to Turbo Vision, batteries-included, single binary,
> zero ceremony.

A developer should be able to:

```
go install github.com/tamnd/vigo/cmd/vigo@latest
vigo .
```

…and instantly land in a familiar Borland-style desktop with a working editor,
project tree, build/run integration, gopls-powered code intelligence, and
Delve-powered debugging — all keyboard-first, mouse-aware, themeable, and
scriptable.

## Goals

1. **100% Turbo Vision look & feel.** Authentic palette, double/single-line
   borders, drop-shadow windows, F1 contextual help, F10 menu activation, Alt-X
   exit, modal dialogs with focus rings, status line with command hints.
2. **Pure Go.** No CGO; one static binary. Built on `gdamore/tcell/v2`.
3. **Faithful object model.** `View`, `Group`, `Window`, `Dialog`, `Application`,
   `Desktop`, `MenuBar`, `StatusLine`, command system, palettes, streams — the
   exact concepts from Turbo Vision 2.0, idiomatically rewritten for Go using
   composition + small interfaces instead of inheritance.
4. **First-class Go IDE.** Editor with syntax highlighting, gopls LSP for
   navigation/completion/diagnostics/refactor, `go build|test|run` integration,
   `dlv dap` debugger with breakpoints, watches, and stepping.
5. **Cross-platform.** Linux, macOS, Windows (Windows Console + Terminal),
   FreeBSD. Unicode by default, with CP437 fall-back glyph table for retro
   terminals.
6. **Hackable.** Plain text resource format for menus, palettes, and key maps.
   Embedded scripting hook (Starlark) for extensions in v1.0+.

## Non-goals

- Not a graphical (pixel) IDE; Vigo lives in a character grid.
- Not a clone of VS Code / GoLand. Mouse is supported, but the keyboard is
  authoritative; UI density and aesthetics follow Borland 1990, not modern
  flat-design conventions.
- Not a general LSP client framework. gopls is the primary, supported backend;
  other languages may work but are not on the roadmap.
- No bundled package manager, no AI features in v1.0 (left to plugins).

## Why now

- Terminal renaissance: `tcell`, ratatui, BubbleTea, magiblot/tvision, Zellij
  prove that high-density TUIs are back in fashion.
- Go's gopls + Delve DAP make a self-hosted Go IDE genuinely feasible without
  reimplementing parsers or debuggers.
- Existing Turbo Vision ports (magiblot/tvision in C++, turbo-vision in Rust)
  have proven the design ports cleanly to modern memory-safe languages — but
  there is no canonical Go port.

## High-level shape

```
+--[≡]--File--Edit--Search--Run--Compile--Debug--Project--Options--Window--Help--+
|                                                                                |
|   ╔══[•]═ main.go ═══════════════════════════════════════════════════════ 1 ╗  |
|   ║ package main                                                            ║  |
|   ║                                                                         ║  |
|   ║ import "fmt"                                                            ║  |
|   ║                                                                         ║  |
|   ║ func main() {                                                           ║  |
|   ║     fmt.Println("hello, vigo")                                          ║  |
|   ║ }                                                                       ║  |
|   ║                                                                         ║  |
|   ╚═════════════════════════════════════════════ 6:1 ════════════ INS ══════╝  |
|                                                                                |
+--F1 Help--F2 Save--F3 Open--F5 Zoom--F6 Switch--F9 Make--F10 Menu--Alt-X Exit--+
```

## References

- Borland, *Turbo Vision Programmer's Guide*, 1990 (bitsavers PDF).
- magiblot/tvision — modern C++ port with Unicode.
- aovestdipaperino/turbo-vision-4-rust — Rust port (1.0 in 2025).
- gdamore/tcell, rivo/tview — Go terminal foundation.
- gopls (`golang.org/x/tools/gopls`) and Delve `dlv dap`.

See:
- 1598_vigo_roadmap.md — phased delivery plan.
- 1599_vigo_architecture.md — layered design.
- 1600..1607 — per-version specifications.
