# 1598 — Vigo Roadmap

Phased delivery from empty repo (v0.1) to production-ready 1.0. Each phase is
its own milestone with a single corresponding spec. A phase is "done" only when
its acceptance gates in the spec are met *and* a binary tagged `vN.M.0` ships
via GitHub Actions. No skipping versions; no breaking changes mid-phase.

## Phases at a glance

| Ver  | Spec | Theme                           | Headline deliverable                                  |
|------|------|---------------------------------|-------------------------------------------------------|
| 0.1  | 1600 | Foundation                      | tcell-backed event loop + view tree + Borland palette |
| 0.2  | 1601 | TUI core                        | Windows, dialogs, menus, status line, commands        |
| 0.3  | 1602 | Editor                          | Multi-buffer text editor with undo/find/replace       |
| 0.4  | 1603 | Project & toolchain             | Project tree, build/run/test integration             |
| 0.5  | 1604 | Code intelligence               | gopls LSP: complete, hover, goto, diagnostics, rename |
| 0.6  | 1605 | Debugger                        | Delve DAP: breakpoints, stepping, watches, eval       |
| 0.7  | 1606 | Polish                          | Themes, key-maps, help system, plugin host (Starlark) |
| 1.0  | 1607 | Release                         | Stability, docs, packaged binaries, signed releases   |

Each version is a 2–4 week milestone. The roadmap is sequential — each phase
builds on the previous one's data structures and APIs.

## Acceptance gates (every version)

Every version must, before tagging, satisfy:

1. **Builds clean** on Linux/macOS/Windows via GitHub Actions matrix.
2. **`golangci-lint`** passes with the project's `.golangci.yml` (no `nolint`
   without justification).
3. **Tests** ≥ 70% coverage on new packages, ≥ 80% on `internal/vio`,
   `internal/event`, `internal/view`, `internal/edit`.
4. **No known regressions**: prior version's demo programs still run.
5. **CHANGELOG.md** entry under `## [vN.M.0] - YYYY-MM-DD` with `Added`,
   `Changed`, `Fixed`, `Removed` headings.
6. **Demo recording** in `docs/demos/vN.M/` (asciinema cast) showing the new
   functionality.

## v0.1 — Foundation (spec 1600)

- `internal/vio`: tcell screen abstraction, `Cell`, `Attr`, `DrawBuffer`,
  CP437↔Unicode glyph table, single/double border helpers.
- `internal/event`: event types (key, mouse, broadcast, command, idle), queue,
  dispatch.
- `internal/view`: `View`, `Group`, focus chain, options bitmask, growMode,
  state bitmask, palette indirection.
- `internal/app`: `Application`, `Desktop`, `BackgroundView`, idle loop, main
  loop, suspend/resume, signal handling.
- `internal/menu`: minimal `MenuBar` + `StatusLine` shells (functional in v0.2).
- `cmd/vigo`: launches an empty Borland-blue desktop with menubar and status
  line; Alt-X quits.
- Authentic 16-color BIOS palette, double-line borders, drop shadows.

## v0.2 — TUI core (spec 1601)

- `internal/window`: `Window`, `Dialog`, frame with title, close, zoom, number;
  drag/resize via mouse and Ctrl-F5; modal `ExecView`.
- `internal/widget`: `Label`, `InputLine`, `Button`, `CheckBoxes`,
  `RadioButtons`, `ListBox` with scrollbar, `Memo`, `StaticText`.
- `internal/menu`: full pull-down menus, sub-menus, accelerators, dim state.
- `internal/cmd`: command bus with `cmCut`/`cmCopy`/`cmPaste`/`cmQuit`/...,
  enable/disable masks per active view.
- `internal/help`: F1 contextual help index → about dialog scaffold.
- Demo: file-open dialog, calculator window, calendar window, ASCII-table
  window — three classic Turbo Vision tools, all working.

## v0.3 — Editor (spec 1602)

- `internal/buffer`: gap-buffer text storage with line index, UTF-8 boundaries,
  efficient column/row math.
- `internal/edit`: `EditWindow`, `Editor`, multi-cursor selections, clipboard,
  undo/redo (linear + grouped), find/replace (regex via `regexp`),
  word-wrap toggle, tabs vs spaces, EOL detection.
- `internal/syntax`: syntax highlighter framework + Go highlighter (re-using
  `go/scanner`).
- File operations: open, save, save-as, autosave, recent files, file→buffer
  reload on disk-change.
- `cmd/vigo .` opens the project tree window + lets you open and edit Go files.

## v0.4 — Project & toolchain (spec 1603)

- `internal/project`: `go.mod` discovery, module graph, file tree backed by
  `os.ReadDir` cache + fsnotify.
- `internal/run`: `go build`, `go run`, `go test`, `go vet`, `gofmt`,
  `goimports` integration; output captured to a tail-log window with
  click-to-jump on error patterns (`file:line:col`).
- `internal/output`: scrollback buffer, ANSI passthrough, error parser.
- Run targets defined in `.vigo/targets.toml`.
- F9 = build, Ctrl-F9 = run, F11 = test, Shift-F9 = clean.

## v0.5 — Code intelligence (spec 1604)

- `internal/lsp`: JSON-RPC 2.0 over stdio, LSP client implementation, gopls
  process supervisor.
- Editor integrations: completion popup, signature help, hover, goto-definition
  / -declaration / -references, document symbols, workspace symbols, rename,
  organize-imports, format-on-save, diagnostics squiggles + Problems window,
  code lens, semantic tokens for highlighting.
- Symbol palette (Ctrl-Shift-O) and global search (Ctrl-Shift-F).

## v0.6 — Debugger (spec 1605)

- `internal/dap`: DAP client, `dlv dap` supervisor.
- Run/Debug configurations in `.vigo/launch.toml`.
- Debug UI: breakpoint gutter, call-stack window, variables/watch window,
  registers, goroutines window, threads, REPL/eval prompt.
- Step over/into/out (F8/F7/Shift-F7), continue (F9 in debug mode), pause,
  conditional and log breakpoints, exception breakpoints.

## v0.7 — Polish (spec 1606)

- Theming engine: hot-reloadable palettes; bundled themes (Borland Classic,
  Twilight, Solarized, Monochrome, High-Contrast).
- Key-map remapping via `~/.config/vigo/keymap.toml`; vim-mode opt-in.
- F1 help system with hyperlinked `.hlp` files (modern `.md` re-imagined).
- Plugin host via Starlark scripts; expose `vigo.cmd`, `vigo.buffer`,
  `vigo.editor` namespaces.
- Mouse improvements: wheel, smooth-scroll, OSC 8 hyperlinks, bracketed paste.

## v1.0 — Release (spec 1607)

- API freeze on `pkg/vigo` re-exports.
- 90%+ docs coverage on exported identifiers.
- `goreleaser` release workflow: macOS (universal), Linux (amd64+arm64),
  Windows (amd64+arm64), FreeBSD; Homebrew tap, scoop bucket, AUR PKGBUILD.
- Signed checksums, SBOM (`syft`), provenance attestation.
- Tutorial: "Build Hello-World end-to-end in Vigo" + 90-second asciinema.

## Out-of-roadmap candidates (post-1.0)

- Remote dev mode (SSH attach), pair programming via shared session.
- Treesitter-based highlighter for non-Go languages.
- AI assistant pane (LSP-style protocol to local model).
- GUI rendering target (using `bubbletea`-style adaptive layout).

## Versioning

SemVer; pre-1.0 minor bumps may break APIs. After 1.0, only major bumps may
break exported `pkg/vigo` surface. The `internal/` tree is never API-stable.
