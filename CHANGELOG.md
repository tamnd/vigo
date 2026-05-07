# Changelog

All notable changes to this project are recorded here. Format follows
Keep a Changelog 1.1.0; the project uses Semantic Versioning.

## [Unreleased]

## [v0.2.0] - 2026-05-08

### Added

- `window`: draggable, resizable, closable `Window` with a `Frame`
  child (titled border, close box, zoom box) and Z-order bring-to-front
  on focus.
- `dialog`: modal `Dialog` themed for forms, with default/cancel
  button slots, button row layout, and a `MessageBox` factory for
  the standard kind/button combinations.
- `widget`: the standard leaf views — `StaticText`, `Label`,
  `Button` (mnemonic parsing, default/normal flags, palette slots),
  `InputLine` with selection/clipboard, `Cluster`/`CheckBoxes`/
  `RadioButtons`, `ScrollBar` with arrows/page/thumb,
  `ListViewer`/`ListBox`, multi-line `Memo`, `ParamText`, and a
  `History` ring with a dropdown picker.
- `cmd`: command bus with a typed command set, an `Enabler` that
  notifies subscribers on enable/disable transitions, and a
  `Bindings` table that resolves a `KeyEvent` to a `CommandID`.
- `menu`: pull-down `MenuBox` with separators, dim disabled rows,
  hotkey letters and right-aligned shortcuts, plus submenu arrows.
  The `Bar` runs an open box via an injectable `MenuRunner` and
  cycles between top-level menus on left/right arrow.
- `menu.LoadTOML` / `menu.SaveTOML`: round-tripping reader and
  writer for a minimal `[[menu]] / [[menu.items]]` subset, so the
  IDE's menu tree can live in `assets/menus/main.toml`.
- `assets`: embeds `menus/main.toml` so `cmd/vigo` boots with a
  declarative menu tree at startup.
- `help`: contextual-help slot with a `HelpCtx` registry, a
  `Default()` registry seeded with the About topic, and an
  `About(r)` modal dialog. F1 and `CmdHelp` open About at the
  application loop level.
- `app`: modal sub-loop via `ExecView`, modal slot routing in the
  event dispatcher, and `HelpRegistry` exposed for hosts.
- `demos`: classic Turbo Vision sample tools — Calculator,
  Calendar, ASCIITable, Puzzle — each as a standalone window with
  unit tests for state transitions.
- `cmd/vigo-demos`: launcher binary that wires the four demos to
  menu commands and opens them on the desktop via a hidden
  post-process command sink.

### Changed

- `cmd/vigo` boots its menu bar from `assets/menus/main.toml` via
  embed instead of using the hard-coded default item list.
- `help.Version` now reads `v0.2.0`.

## [v0.1.0] - 2026-05-07

### Added

- Foundation packages laid out at the repo root: `vio`, `event`, `view`,
  `menu`, `app`, with a `cmd/vigo` binary that boots straight into the
  Borland-blue desktop.
- `vio`: tcell-backed `Screen`, in-memory `Surface` compositor,
  `Palette`, `Attr`, single and double border glyph sets, drop-shadow
  helper, ASCII fallback, and a fake screen for tests.
- `event`: typed event union (key, mouse, command, broadcast, idle,
  resize) with class bitmasks.
- `view`: `View` base, `Viewer` interface, `Group` with three-phase
  event dispatch (pre, focused, post) and `GrowMode`-aware resize
  propagation.
- `menu`: `MenuBar` and `StatusLine` with Turbo Vision IDE defaults
  (File..Help, F1/F10/Alt-X).
- `app`: `Application` main loop, thread-safe `PutEvent`, Alt-X and
  `CmdQuit` shutdown, golden snapshot test for the empty desktop.
- GitHub Actions CI: lint with golangci-lint v2.5 against the v2 config,
  build and test on Linux/macOS/Windows with `-race -cover`, `go mod
  tidy` check.
- GitHub Actions release: matrix cross-build (linux, darwin, windows,
  freebsd) staged into tar.gz/zip archives with `LICENSE`,
  `README.md`, `CHANGELOG.md`, plus a SHA-256 checksums manifest.
- Dependabot for Go modules and GitHub Actions.
- MIT license, README, full multi-version specification in `docs/`
  (overview, roadmap, architecture).

[Unreleased]: https://github.com/tamnd/vigo/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/tamnd/vigo/releases/tag/v0.2.0
[v0.1.0]: https://github.com/tamnd/vigo/releases/tag/v0.1.0
