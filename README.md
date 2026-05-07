# vigo

[![CI](https://github.com/tamnd/vigo/actions/workflows/ci.yml/badge.svg)](https://github.com/tamnd/vigo/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/tamnd/vigo.svg)](https://pkg.go.dev/github.com/tamnd/vigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/tamnd/vigo)](https://goreportcard.com/report/github.com/tamnd/vigo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> A 100% faithful, modernized port of Borland's Turbo Vision IDE,
> reimagined as a self-hosted **Go IDE written in Go**, in your terminal.

`vigo` boots into the iconic **Borland Classic blue desktop** with double-line
borders, F-key shortcuts, mouse-aware menus and dialogs — the look and feel
of Turbo Pascal 6.0 / Borland C++ 3.x, faithfully reproduced — and ships a
complete edit/build/debug workflow for Go programmers, powered by `gopls` and
`dlv` underneath.

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

## Status

**Pre-alpha** — actively bootstrapping toward `v0.1` (foundation). See
[`docs/roadmap.md`](docs/roadmap.md) and the per-version specs (1597–1607)
for the full plan.

| Version | Theme               | Status        |
|---------|---------------------|---------------|
| 0.1     | Foundation          | 🚧 in progress |
| 0.2     | TUI core            | ⏳ planned     |
| 0.3     | Editor              | ⏳ planned     |
| 0.4     | Project & toolchain | ⏳ planned     |
| 0.5     | Code intelligence   | ⏳ planned     |
| 0.6     | Debugger            | ⏳ planned     |
| 0.7     | Polish              | ⏳ planned     |
| 1.0     | Release             | ⏳ planned     |

## Why?

- **Authentic.** Pixel-faithful to the 1990 Borland TV look — palette,
  glyphs, drop shadows, F-keys, modal dialogs, command system.
- **Pure Go.** No CGO. One static binary on Linux, macOS, Windows, FreeBSD.
- **Full IDE.** Editor, gopls (LSP), Delve (DAP), build/run/test, project
  tree, find-in-files, refactor — everything in the terminal.
- **Hackable.** Plain-text menus, themes, key-maps. Embedded Starlark plugin
  host (planned for v0.7).
- **Single-threaded UI.** No data races, no spinners, sub-millisecond input
  latency on a warm event loop.

## Quickstart

```sh
# Once published:
go install github.com/tamnd/vigo/cmd/vigo@latest
vigo .
```

For now, build from source:

```sh
git clone https://github.com/tamnd/vigo
cd vigo
go run ./cmd/vigo
```

Press **Alt-X** (or **F10 → File → Exit**) to quit.

## Project layout

```
vigo/
├── cmd/vigo/        # binary entrypoint
├── internal/
│   ├── vio/         # tcell-backed screen, palette, surface, glyphs
│   ├── event/       # typed events, queue, dispatch primitives
│   ├── view/        # View / Group / focus / palette indirection
│   ├── app/         # Application / Desktop / Background / main loop
│   └── menu/        # MenuBar / StatusLine
├── docs/            # generated docs and specs (roadmap, architecture)
├── .github/         # CI, release workflows
├── CHANGELOG.md
├── LICENSE          # MIT
└── README.md
```

## Architecture (tl;dr)

`vigo` is a thin layer over [`gdamore/tcell/v2`](https://github.com/gdamore/tcell)
implementing Borland's Turbo Vision object model in idiomatic Go (composition
+ small interfaces, never embedded inheritance trees). The framework is
**single-threaded for UI**; gopls and Delve run as supervised subprocesses
that post events back to the event loop via a thread-safe channel.

For the load-bearing reference, see
[`docs/architecture.md`](docs/architecture.md).

## References

- Borland, *Turbo Vision Programmer's Guide*, 1990
  ([Internet Archive](https://archive.org/details/bitsavers_borlandturVersion6.0TurboVision1990_16007263)).
- [magiblot/tvision](https://github.com/magiblot/tvision) — modern C++ port.
- [aovestdipaperino/turbo-vision-4-rust](https://github.com/aovestdipaperino/turbo-vision-4-rust) — Rust port.
- [`gdamore/tcell`](https://github.com/gdamore/tcell) — terminal cell library.
- [`golang.org/x/tools/gopls`](https://pkg.go.dev/golang.org/x/tools/gopls).
- [`go-delve/delve`](https://github.com/go-delve/delve) (`dlv dap`).

## Contributing

Patches, issues, and feature ideas welcome. Please read
[`CONTRIBUTING.md`](CONTRIBUTING.md) (TBA) and open a discussion before
landing big design changes — the framework is small and opinionated by design.

## License

[MIT](LICENSE) © 2026 Tam Nguyen Duy.
