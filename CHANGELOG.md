# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- v0.1 foundation: tcell-backed screen abstraction, event loop, view tree,
  Borland Classic palette, double-line frame helpers, drop shadows.
- `Application`, `Desktop`, `Background`, `MenuBar`, `StatusLine` baseline
  views — enough to render an authentic Turbo Vision-style empty desktop.
- `cmd/vigo` demo binary; `Alt-X` quits cleanly.
- GitHub Actions CI matrix (Linux / macOS / Windows) with `golangci-lint`,
  `go vet`, `go test -race -cover`.
- MIT license, README, full multi-version specification (1597–1607).
