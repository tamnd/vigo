# 1599. vigo architecture

This is the load-bearing technical reference. It describes layers,
packages, the object model, the event and command system, the drawing
pipeline, and how external processes (gopls, dlv) plug in. The
per-version specs (1600..1607) refine this doc; they do not contradict
it.

## 1. Layered design

```
                                  cmd/vigo
                                     |
   +---------------+---------------+--------------+----------------+
   |               |               |              |                |
   |  app          |  edit         |  project     |  debug (dap)   |
   |  desktop      |  buffer       |  run         |  lsp           |
   |  menu         |  syntax       |  output      |  watch (fs)    |
   |  widget       |               |              |                |
   +-------+-------+-------+-------+------+-------+-------+--------+
           |               |              |               |
           v               v              v               v
   +---------------------------------------------------------+
   |   view  -- View / Group / focus / palette               |
   |   event -- events, commands, dispatch                   |
   |   cmd   -- command bus + key bindings                   |
   |   help  -- help context registry                        |
   +---------------------------------------------------------+
                              |
                              v
                 +-------------------------+
                 |  vio                    |
                 |  - Screen wrapper       |
                 |  - DrawBuffer / Cell    |
                 |  - Palette / Attr       |
                 |  - Glyphs (CP437/UCS)   |
                 +------------+------------+
                              v
                       gdamore/tcell/v2
```

Rules:

- Packages above only depend on packages strictly below.
- `pkg/vigo` re-exports the v1.0-stable surface (frozen at 1.0).
- Pre-1.0, the rest of the tree may break across minor versions freely.
- Only `vio` imports `tcell`. The screen layer is swappable
  for tests, headless rendering, or future GUI targets.

## 2. Object model

Turbo Vision is class-based with single inheritance and (some)
multiple-inheritance flavour. Go has neither, so the translation uses
three primitives:

1. A concrete embed-friendly base struct `View` with all common state
   (bounds, owner, options, state, growmode, eventmask, helpCtx,
   palette indirection).
2. Small interfaces that consumers test for at run-time:
   ```go
   type Drawable interface { Draw(*Surface) }
   type Eventer  interface { HandleEvent(*Event) }
   type Sizer    interface { ChangeBounds(Rect) }
   ```
   The default methods on `*View` satisfy each of these; concrete
   views embed `*View` and override what they need.
3. A `Group` is a `View` that owns child views, dispatches events
   down, composes a Surface, and tracks the focus chain.

Hierarchy mirroring Turbo Vision:

```
View                           - all visual primitives
  +-- Group                    - container of Views
        +-- Window             - framed, movable, resizable
        |     +-- Dialog       - modal Window with default button
        |     +-- EditWindow   - (v0.3) editor host
        +-- Desktop            - screen-spanning background group
        +-- MenuBar            - top-row pull-down host
        +-- StatusLine         - bottom-row hint host
        +-- Application        - root group, owns Desktop+MenuBar+StatusLine
View
  +-- Frame, ScrollBar, Background        (decorations)
  +-- Label, InputLine, Button            (controls)
  +-- ListView, ListBox, ListViewer       (list family)
  +-- CheckBoxes, RadioButtons, Cluster   (cluster family)
  +-- Memo, StaticText, ParamText         (text family)
  +-- ColorSelector, MonoSelector         (palette tools)
```

State and options bitmasks reproduce Turbo Vision exactly:

```go
const (
    StateVisible    State = 1 << iota
    StateCursorVis
    StateCursorIns
    StateShadow
    StateActive
    StateSelected
    StateFocused
    StateDragging
    StateDisabled
    StateModal
    StateDefault
    StateExposed
)

const (
    OptSelectable Options = 1 << iota
    OptTopSelect
    OptFirstClick
    OptFramed
    OptPreProcess
    OptPostProcess
    OptBuffered
    OptTileable
    OptCenterX
    OptCenterY
)
```

## 3. Event model

Every event is a value type, not an interface, copied by value (small
enough, no escape analysis surprises):

```go
type Event struct {
    What  Class
    Mouse MouseEvent
    Key   KeyEvent
    Msg   MessageEvent
}
```

Event flow on a Group:

1. Pre-processing. Walk children top-of-Z first; call HandleEvent on
   any child with `OptPreProcess` (e.g. menu accelerators).
2. Focused dispatch. Forward to the current child.
3. Post-processing. Walk again; call HandleEvent on any child with
   `OptPostProcess` (e.g. status-line hot keys).

After handling, a view calls `ev.Clear()` to mark it consumed.
Unconsumed events bubble up to the owning group, then to the
Application.

The command system decouples actions from views. A view emits
`Event{What: ClassCommand, Msg: MessageEvent{Command: cmSave}}`. Any
view in the chain may handle it. Application owns global handlers
(cmQuit, cmMenu, cmHelp). A `CommandSet` bitset gates enabled
commands; MenuBar and StatusLine re-render automatically when
`enableCommands` or `disableCommands` is called.

Key bindings live in `cmd/keymap.go`. Default Borland map:

```
F1=cmHelp  F2=cmSave  F3=cmOpen  F4=cmRunUntilCursor
F5=cmZoom  F6=cmNextWin F7=cmStepInto F8=cmStepOver
F9=cmMake  F10=cmMenu Alt-X=cmQuit Ctrl-Ins=cmCopy ...
```

## 4. Drawing pipeline

```
Application.draw()
  +-- Group.draw() recursively
        +-- for each child where StateVisible & StateExposed:
              child.Draw(surface) -> Surface
        +-- Group composes children onto its own buffer
  +-- vio.Screen.Show(surface)  -- flushes diff to tcell
```

Invariants:

- Drawing is idempotent and pure. `Draw` reads view state and writes
  a buffer. No side effects, no I/O.
- Drawing happens at most once per event loop iteration, after
  dispatch. Dirtiness propagates up via `drawView` / `drawSubViews`
  which mark the owning group dirty.
- Damage tracking is deferred. v0.1 redraws the full Application each
  frame; tcell diffs to the terminal. v0.2 onwards adds a coarse
  damage-rect cache.
- Cursor is drawn last, after all views, by the focused leaf.
- A 60 fps cap is enforced by a render token. The token only fires a
  redraw when the dirty flag is set, so an idle vigo uses near-zero
  CPU.

## 5. Color and palette

Each View has a 1-byte `palette index` referencing its parent group's
palette string; the group translates `(child.palette -> bytes ->
Application palette -> Attr)`. This is the Borland indirection scheme
verbatim. It lets themes swap at runtime without touching any view's
draw code.

`vio.Attr` packs `(fg, bg, mod)` where:

- fg, bg: `tcell.Color` (24-bit plus named-default)
- mod: bold/italic/underline/reverse/blink/strike

Bundled themes are TOML files in `assets/themes/*.toml`, loaded into
the application palette on startup or via `cmTheme`.

## 6. Concurrency model

The TUI is single-threaded. Event loop, drawing, and view state
mutation all run on the main goroutine. There is no shared mutable
view state between goroutines.

Background work runs on worker goroutines (`run`, `lsp`, `dap`,
`watch`):

```
worker goroutine -- completes work --> chan PostedEvent
                                              |
                                              v
                                  Application.PutEvent(...)
                                              |
                                              v
                                       main event loop
```

`Application.PutEvent` is the only thread-safe entry point into the
UI. Workers post `ClassBroadcast` events with typed payloads; the
responsible owner view consumes them on the main thread.

Cancellation: every worker takes a `context.Context`. The app's root
context is cancelled on quit, propagating shutdown to gopls and dlv
subprocesses with a 2-second grace period before SIGKILL.

## 7. External process integration

### 7.1 LSP / gopls (v0.5)

```
Editor -> lsp.Client
            |  JSON-RPC 2.0 over stdio
            v
         gopls subprocess (managed lifecycle)
```

- One gopls per project root (`go.mod` dir).
- Methods used: `initialize`, `textDocument/didOpen|didChange|didClose`,
  `completion`, `hover`, `signatureHelp`, `definition`, `references`,
  `documentSymbol`, `workspaceSymbol`, `rename`, `formatting`,
  `codeAction`, `publishDiagnostics`, `semanticTokens/full`.
- Document sync = Incremental.
- Heavy responses (completions, hovers) are cancellable via
  `$/cancelRequest` on every keystroke that supersedes the request.

### 7.2 DAP / dlv (v0.6)

```
Debug UI -> dap.Client -> dlv dap subprocess -> debuggee
```

- Use Delve's native DAP server (`dlv dap --listen=...`).
- Launch configurations from `.vigo/launch.toml`.
- DAP messages: `initialize`, `launch`, `setBreakpoints`,
  `setExceptionBreakpoints`, `configurationDone`, `threads`,
  `stackTrace`, `scopes`, `variables`, `evaluate`, `continue`,
  `next`, `stepIn`, `stepOut`, `pause`, `terminate`.

## 8. Persistence

- User config: `~/.config/vigo/config.toml`.
- Project config: `<project>/.vigo/{targets,launch,recent}.toml`.
- Window layout: `<project>/.vigo/desktop.json` (Turbo Vision-style
  stream serialization, but JSON instead of the binary original).
- Undo and clipboard history: process memory only.

## 9. Testing strategy

- `vio` tests use a fake screen (`vio/fake_screen.go`) that
  records cells; assertions are made against snapshot strings.
- `view` and `widget` tests build a tiny app, post
  events via `Application.PutEvent`, and snapshot-assert the rendered
  surface.
- `edit` runs against a buffer fixture set with golden
  expected edits.
- LSP and DAP packages have integration tests gated by
  `-tags=integration`, which start a real gopls or dlv.

## 10. Build and release

- Single Go module `github.com/tamnd/vigo`. Go 1.23+, no CGO.
- `cmd/vigo/main.go` is the only entrypoint.
- GitHub Actions matrix `{ubuntu, macos, windows} x {stable}` for CI.
  Release builds via `goreleaser` for 1.0.
- `golangci-lint` enforces style; pre-commit via `lefthook`.
