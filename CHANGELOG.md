# Changelog

All notable changes to this project are recorded here. Format follows
Keep a Changelog 1.1.0; the project uses Semantic Versioning.

## [Unreleased]

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

[Unreleased]: https://github.com/tamnd/vigo/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/tamnd/vigo/releases/tag/v0.1.0
