# Changelog

All notable changes to this project are recorded here. Format follows
Keep a Changelog 1.1.0; the project uses Semantic Versioning.

## [Unreleased]

### Added

- v0.1 foundation: tcell-backed screen abstraction, event types and
  dispatch primitives, View and Group base, Application main loop.
- Borland classic palette (16 entries), single and double border glyph
  sets, drop-shadow helper, ASCII fallback border set.
- Application boots into a blue desktop with a stub menu bar and status
  line. Alt-X quits cleanly.
- GitHub Actions CI: lint with golangci-lint v1.62, build and test on
  Linux/macOS/Windows with `-race -cover`, and a `go mod tidy` check.
- GitHub Actions release: cross-builds Linux/macOS/Windows/FreeBSD on
  `v*` tags, with SHA-256 checksums.
- Dependabot for Go modules and GitHub Actions.
- MIT license, README, full multi-version specification in `docs/`
  (overview, roadmap, architecture).
