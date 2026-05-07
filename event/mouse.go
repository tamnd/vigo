package event

// MouseButton is a bitmask of mouse buttons currently pressed.
type MouseButton uint8

// Mouse button bits.
const (
	MouseLeft MouseButton = 1 << iota
	MouseMiddle
	MouseRight
	MouseWheelUp
	MouseWheelDown
)

// MouseEvent is the payload of any Mouse* event. For Resize events, X and Y
// carry the new screen width and height (no other field is meaningful).
type MouseEvent struct {
	X, Y    int
	Buttons MouseButton
	Mod     Modifier
}
