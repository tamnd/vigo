# 1598. vigo roadmap

Phased delivery from empty repo (v0.1) to 1.0. Each phase has its own
spec. A phase is done only when its acceptance gates pass and a binary
tagged `vN.M.0` ships via GitHub Actions. Phases are sequential; each
builds on the previous one's data structures and APIs. No skipped
versions; no breaking changes mid-phase.

## Phases

| Ver  | Spec | Theme               | Headline                                              |
|------|------|---------------------|-------------------------------------------------------|
| 0.1  | 1600 | Foundation          | tcell event loop + view tree + Borland palette        |
| 0.2  | 1601 | TUI core            | Windows, dialogs, menus, status line, commands        |
| 0.3  | 1602 | Editor              | Multi-buffer text editor with undo, find, replace     |
| 0.4  | 1603 | Project & toolchain | Project tree, build, run, test integration           |
| 0.5  | 1604 | Code intelligence   | gopls LSP: complete, hover, goto, diagnostics, rename |
| 0.6  | 1605 | Debugger            | Delve DAP: breakpoints, stepping, watches, eval       |
| 0.7  | 1606 | Polish              | Themes, key-maps, help, plugin host (Starlark)        |
| 1.0  | 1607 | Release             | Stability, docs, packaged signed binaries             |

Each version is a 2 to 4 week milestone.

## Acceptance gates that apply to every version

Before tagging:

1. Builds clean on Linux, macOS, Windows via the GitHub Actions matrix.
2. `golangci-lint` passes with the project config; no `nolint` without
   a justification comment.
3. Tests pass with at least 70% line coverage on new packages, 80% on
   `internal/vio`, `internal/event`, `internal/view`, `internal/edit`.
4. Prior version's demo programs still run.
5. `CHANGELOG.md` has an entry under `## [vN.M.0] - YYYY-MM-DD` with
   `Added`, `Changed`, `Fixed`, `Removed` headings.
6. A short asciinema cast goes in `docs/demos/vN.M/`.

## v0.1, foundation, spec 1600

- `internal/vio`: tcell screen, `Cell`, `Attr`, `Surface`, CP437 to
  Unicode glyph table, single and double border helpers.
- `internal/event`: typed events (key, mouse, broadcast, command,
  idle, resize), queue, dispatch primitives.
- `internal/view`: View, Group, focus chain, options bitmask, growMode,
  state bitmask, palette indirection.
- `internal/app`: Application, Desktop, Background, idle loop, main
  loop, suspend/resume, signal handling.
- `internal/menu`: minimal MenuBar and StatusLine (functional in v0.2).
- `cmd/vigo`: launches an empty Borland-blue desktop with menu bar
  and status line; Alt-X quits.
- 16-color BIOS palette, double-line borders, drop shadows.

## v0.2, TUI core, spec 1601

- `internal/window`: Window, Dialog, frame with title, close, zoom,
  number; drag and resize via mouse and Ctrl-F5; modal `ExecView`.
- `internal/widget`: Label, InputLine, Button, CheckBoxes, RadioButtons,
  ListBox with scrollbar, Memo, StaticText.
- `internal/menu`: full pull-downs, submenus, accelerators, dim state.
- `internal/cmd`: command bus with `cmCut`, `cmCopy`, `cmPaste`,
  `cmQuit`, etc.; per-view enable/disable masks.
- `internal/help`: F1 contextual help index, About dialog scaffold.
- Demo: file-open dialog plus the four classic Turbo Vision tools
  (calculator, calendar, ASCII table, puzzle).

## v0.3, editor, spec 1602

- `internal/buffer`: gap buffer with line index, UTF-8 boundaries,
  efficient column/row math.
- `internal/edit`: EditWindow, Editor, multi-cursor selections,
  clipboard, undo/redo (linear plus grouped), find/replace
  (regexp via Go's `regexp`), word-wrap toggle, tabs vs spaces, EOL
  detection.
- `internal/syntax`: highlighter framework plus a Go highlighter built
  on `go/scanner`.
- File ops: open, save, save-as, autosave, recent files, reload-on-disk-
  change.
- `vigo .` opens the project tree window and lets you open and edit
  Go files.

## v0.4, project and toolchain, spec 1603

- `internal/project`: `go.mod` discovery, module graph, file tree
  backed by `os.ReadDir` and fsnotify.
- `internal/run`: `go build|run|test|vet`, `gofmt`, `goimports`. Streams
  output into a tail-log window with click-to-jump on `file:line:col`.
- `internal/output`: scrollback, ANSI passthrough, error parser.
- Run targets in `.vigo/targets.toml`.
- F9 build, Ctrl-F9 run, F11 test, Shift-F9 clean.

## v0.5, code intelligence, spec 1604

- `internal/lsp`: JSON-RPC 2.0 over stdio, LSP client, gopls
  supervisor.
- Editor integrations: completion popup, signature help, hover,
  goto-definition / -declaration / -references, document symbols,
  workspace symbols, rename, organize-imports, format-on-save,
  diagnostics squiggles plus Problems window, code lens, semantic
  tokens.
- Symbol palette (Ctrl-Shift-O), global search (Ctrl-Shift-F).

## v0.6, debugger, spec 1605

- `internal/dap`: DAP client, `dlv dap` supervisor.
- Run/Debug configuration in `.vigo/launch.toml`.
- Debug UI: breakpoint gutter, call-stack, variables/watch, registers,
  goroutines, threads, REPL/eval prompt.
- Step over/into/out (F8 / F7 / Shift-F7), continue (F9 in debug),
  pause, conditional and log breakpoints, exception breakpoints.

## v0.7, polish, spec 1606

- Theming: hot-reloadable palettes; bundled themes (Borland Classic,
  Twilight, Solarized, Monochrome, High-Contrast).
- Key-map remapping via `~/.config/vigo/keymap.toml`; vim-mode opt-in.
- F1 help system with hyperlinked Markdown topics.
- Starlark plugin host; expose `vigo.cmd`, `vigo.buffer`, `vigo.editor`
  namespaces.
- Mouse polish: wheel, smooth-scroll, OSC 8 hyperlinks, bracketed paste.

## v1.0, release, spec 1607

- API freeze on `pkg/vigo` re-exports.
- 90%+ docs coverage on exported identifiers.
- `goreleaser` workflow: macOS (universal), Linux (amd64, arm64),
  Windows (amd64, arm64), FreeBSD; Homebrew tap, scoop bucket,
  AUR PKGBUILD.
- Signed checksums, SBOM via `syft`, provenance attestation.
- A "Build hello-world end to end in vigo" tutorial plus a 90-second
  asciinema.

## Out-of-roadmap candidates (post-1.0)

- Remote dev mode (SSH attach), shared session.
- Tree-sitter highlighter for non-Go languages.
- AI assistant pane (LSP-style protocol to a local model).
- Adaptive GUI rendering target.

## Versioning

SemVer. Pre-1.0 minor bumps may break APIs. After 1.0, only major
bumps may break the exported `pkg/vigo` surface. The `internal/` tree
is never API-stable.
