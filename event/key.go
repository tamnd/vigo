package event

// Modifier is a bitmask of keyboard modifier keys.
type Modifier uint8

// Modifier bits.
const (
	ModNone  Modifier = 0
	ModShift Modifier = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
)

// Key identifies a non-printable key, or KeyRune for printable input where
// the rune lives in KeyEvent.Rune.
type Key uint16

// Key constants. Function keys, navigation, and KeyRune for printable text.
const (
	KeyNone Key = iota
	KeyRune
	KeyEnter
	KeyEsc
	KeyTab
	KeyShiftTab
	KeyBackspace
	KeyDelete
	KeyInsert
	KeyHome
	KeyEnd
	KeyPgUp
	KeyPgDn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)

// KeyEvent is the payload of a KeyDown event.
type KeyEvent struct {
	Key  Key
	Mod  Modifier
	Rune rune
}
