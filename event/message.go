package event

// CommandID is a numeric identifier for a UI action: file open, save, quit,
// menu activate, and so on. Commands let views send actions without knowing
// who handles them.
type CommandID uint16

// System commands. User-defined commands start at CmdUser.
const (
	CmdNone CommandID = iota
	CmdQuit
	CmdError
	CmdMenu
	CmdClose
	CmdZoom
	CmdResize
	CmdNext
	CmdPrev
	CmdHelp
	CmdOk
	CmdCancel
	CmdYes
	CmdNo

	CmdUser CommandID = 100
)

// MessageEvent is the payload of ClassCommand and ClassBroadcast events.
// Info is an optional typed payload; consumers type-assert as needed.
type MessageEvent struct {
	Command CommandID
	Info    any
}
