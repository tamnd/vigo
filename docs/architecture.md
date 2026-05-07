# 1599 — Vigo Architecture

This document is the load-bearing technical reference for Vigo. It describes
layers, packages, the object model, the event/command system, the drawing
pipeline, and the way external processes (gopls, dlv) are integrated. All
per-version specs (1600..1607) refine this doc but do not contradict it.

## 1. Layered design

```
                                  cmd/vigo
                                     │
   ┌─────────────┬─────────────┬─────┴──────┬──────────────┐
   │  internal/  │  internal/  │  internal/ │  internal/   │
   │  app        │  edit       │  project   │  debug (dap) │
   │  desktop    │  buffer     │  run       │  lsp         │
   │  menu       │  syntax     │  output    │  watch (fs)  │
   │  widget     │             │            │              │
   └─────┬───────┴──────┬──────┴──────┬─────┴───────┬──────┘
         │              │             │             │
         ▼              ▼             ▼             ▼
   ┌────────────────────────────────────────────────────┐
   │   internal/view  ── object model: View/Group/...    │
   │   internal/event ── events, commands, dispatch      │
   │   internal/cmd   ── command bus + key bindings      │
   │   internal/help  ── help context registry           │
   └────────────────────────────────────────────────────┘
                              │
                              ▼
                 ┌─────────────────────────┐
                 │  internal/vio           │
                 │  - Screen wrapper       │
                 │  - DrawBuffer / Cell    │
                 │  - Palette / Attr       │
                 │  - Glyphs (CP437↔UCS)   │
                 └────────────┬────────────┘
                              ▼
                       gdamore/tcell/v2
```

Rules:

- Every package above only depends on packages strictly below it.
- `pkg/vigo` re-exports the v1.0-stable surface (frozen at 1.0).
- `internal/*` may break across minor versions pre-1.0 freely.
- No direct `tcell` imports outside `internal/vio` — keeps the screen layer
  swappable (e.g. for tests, headless rendering, future GUI target).

## 2. Object model

Turbo Vision is class-based with multiple-inheritance-flavoured single
inheritance. Go has neither, so we translate the model with three primitives:

1. **A concrete embed-friendly base struct `View`** with all common state
   (bounds, owner, options, state, growmode, eventmask, helpCtx, palette
   indirection).
2. **Small interfaces** that consumers test for at run-time:
   ```go
   type Drawable interface { Draw(*Surface) }
   type Eventer  interface { HandleEvent(*Event) }
   type Sizer    interface { ChangeBounds(Rect) }
   ```
   The default methods on `*View` satisfy each of these; concrete views embed
   `*View` and override what they need.
3. **`Group`** is a `View` that owns child views, dispatches events down,
   composes a final `Surface`, and tracks the focus chain (`current`/`first`
   doubly linked list as in TV).

Hierarchy mirroring Turbo Vision:

```
View                           ── all visual primitives
  └── Group                    ── container of Views
        ├── Window             ── framed, movable, resizable
        │     ├── Dialog       ── modal Window with default button
        │     └── EditWindow   ── (v0.3) editor host
        ├── Desktop            ── screen-spanning background group
        ├── MenuBar            ── top-row pull-down host
        ├── StatusLine         ── bottom-row hint host
        └── Application        ── root group, owns Desktop+MenuBar+StatusLine
View
  ├── Frame, ScrollBar, Background        (decorations)
  ├── Label, InputLine, Button            (controls)
  ├── ListView, ListBox, ListViewer       (list family)
  ├── CheckBoxes, RadioButtons, Cluster   (cluster family)
  ├── Memo, StaticText, ParamText         (text family)
  └── ColorSelector, MonoSelector         (palette tools)
```

State and options bitmasks reproduce Turbo Vision exactly:

```go
const (
    sfVisible    State = 1 << iota
    sfCursorVis
    sfCursorIns
    sfShadow
    sfActive
    sfSelected
    sfFocused
    sfDragging
    sfDisabled
    sfModal
    sfDefault
    sfExposed
)

const (
    ofSelectable Options = 1 << iota
    ofTopSelect
    ofFirstClick
    ofFramed
    ofPreProcess
    ofPostProcess
    ofBuffered
    ofTileable
    ofCenterX
    ofCenterY
)
```

## 3. Event model

Every event is a value type, not an interface, copied by value (small enough,
escape-free):

```go
type Event struct {
    What  EventClass        // evNothing|evMouse*|evKey*|evCommand|evBroadcast|evIdle
    Mouse MouseEvent
    Key   KeyEvent
    Msg   MessageEvent      // Command int + InfoPtr any
}
```

Event flow on a `Group`:

1. **Phase 1 — pre-processing.** Walk children top-of-Z down; call HandleEvent
   on any with `ofPreProcess` (e.g. menu accelerators).
2. **Phase 2 — focused dispatch.** Forward to `current` child.
3. **Phase 3 — post-processing.** Walk again; call HandleEvent on any with
   `ofPostProcess` (e.g. status-line hot keys).

After handling, a view calls `ev.Clear()` to mark it consumed; unconsumed
events bubble up to the owning group, then to the Application.

The **command system** decouples actions from views. A view emits
`Event{What: evCommand, Msg: MessageEvent{Command: cmSave}}`; any view in the
chain may handle it. Application owns global handlers (cmQuit, cmMenu, cmHelp).
A `CommandSet` bitset gates enabled commands; `MenuBar` and `StatusLine`
re-render automatically on `enableCommands`/`disableCommands` calls.

Key bindings live in `internal/cmd/keymap.go`. Default Borland map:

```
F1=cmHelp  F2=cmSave  F3=cmOpen  F4=cmRunUntilCursor
F5=cmZoom  F6=cmNextWin F7=cmStepInto F8=cmStepOver
F9=cmMake  F10=cmMenu Alt-X=cmQuit Ctrl-Ins=cmCopy ...
```

## 4. Drawing pipeline

```
Application.draw()
  └── Group.draw() recursively
        └── for each child where sfVisible & sfExposed:
              child.Draw(surface) → DrawBuffer
        └── Group composes children onto own buffer
  └── vio.Screen.Sync()  ── flushes diff to tcell
```

Key invariants:

- Drawing is **idempotent and pure**: `Draw` reads view state, writes a buffer.
  No side effects, no I/O.
- Drawing happens at most once per event-loop iteration, after dispatch.
  Dirtiness propagates up via `drawView()`/`drawSubViews()` which mark the
  group dirty.
- Damage tracking: in v0.1 we redraw the full Application each frame (tcell
  diffs to terminal); in v0.2+ we add a coarse damage-rect cache.
- Cursor is drawn last, after all views, by the focused leaf.
- A 60 fps cap is enforced via a render token; the token only fires a redraw
  if the dirty flag is set, so an idle Vigo uses ~0% CPU.

## 5. Color and palette

Each `View` has a 1-byte `palette index` referencing its parent group's
palette string; the group translates `(child.palette → bytes → Application
palette → Attr)`. This is the Borland indirection scheme verbatim, and lets
themes be swapped at runtime without touching any view code.

`vio.Attr` is a `uint64` packing `(fg, bg, mod)` where:

- fg, bg: `tcell.Color` (24-bit + named-default)
- mod: bold/italic/underline/reverse/blink/strike

Bundled themes are TOML files in `assets/themes/*.toml`, loaded into the
application palette on startup or via `cmTheme`.

## 6. Concurrency model

The TUI is **single-threaded**: event loop, drawing, and view state mutation
all run on the main goroutine. This is a hard rule — there is no shared
mutable view state between goroutines.

Background work runs on **worker goroutines** (`run`, `lsp`, `dap`, `watch`):

```
worker goroutine ── completes work ──► chan PostedEvent ── posts to Application
                                              │
                                              ▼
                                  Application.PutEvent(...)
                                              │
                                              ▼
                                       main event loop
```

`Application.PutEvent` is the only thread-safe entry point into the UI.
Workers post `evBroadcast` events with typed payloads; the responsible owner
view consumes them on the main thread.

Cancellation: every worker takes a `context.Context`; the app's root context
is cancelled on quit, propagating shutdown to gopls/dlv subprocesses with a
2-second grace period before SIGKILL.

## 7. External process integration

### 7.1 LSP / gopls (v0.5)

```
Editor → internal/lsp.Client
            │  JSON-RPC 2.0 over stdio
            ▼
         gopls subprocess (managed lifecycle)
```

- One gopls per project root (`go.mod` dir).
- Methods used: `initialize`, `textDocument/didOpen|didChange|didClose`,
  `completion`, `hover`, `signatureHelp`, `definition`, `references`,
  `documentSymbol`, `workspaceSymbol`, `rename`, `formatting`, `codeAction`,
  `publishDiagnostics`, `semanticTokens/full`.
- Document sync = Incremental.
- Heavy responses (completions, hovers) are cancellable via `$/cancelRequest`
  on every keystroke that supersedes the request.

### 7.2 DAP / dlv (v0.6)

```
Debug UI → internal/dap.Client → dlv dap subprocess → debuggee
```

- Use Delve's native DAP server: `dlv dap --listen=...`.
- Launch configurations from `.vigo/launch.toml`.
- DAP messages: `initialize`, `launch`, `setBreakpoints`,
  `setExceptionBreakpoints`, `configurationDone`, `threads`, `stackTrace`,
  `scopes`, `variables`, `evaluate`, `continue`, `next`, `stepIn`, `stepOut`,
  `pause`, `terminate`.

## 8. Persistence

- User config: `~/.config/vigo/config.toml`.
- Project config: `<project>/.vigo/{targets,launch,recent}.toml`.
- Window layout: `<project>/.vigo/desktop.json` (Turbo Vision–style stream
  serialization, but JSON not the original binary format).
- Undo/clipboard history: process memory only.

## 9. Testing strategy

- `internal/vio` tests use a **fake screen** (`vio/fake.go`) that records
  cells; assertions are made against snapshot strings.
- `internal/view` and `internal/widget` tests build a tiny app, post events
  via `Application.PutEvent`, and snapshot-assert the resulting screen text.
- `internal/edit` runs against a buffer fixture set with golden expected
  edits.
- LSP/DAP packages have integration tests gated behind `-tags=integration`,
  which start a real `gopls`/`dlv`.

## 10. Build & release

- Single Go module `github.com/tamnd/vigo`, Go 1.23+ (no CGO).
- `cmd/vigo/main.go` is the only entrypoint.
- GitHub Actions matrix `{ubuntu, macos, windows} × {stable}` for CI; release
  builds via `goreleaser` for 1.0.
- `golangci-lint v2.x` enforces style; pre-commit via `lefthook`.
