// Package event defines the typed events that flow through the vigo view
// tree: keyboard, mouse, commands, broadcasts, idle ticks, and resize.
package event

// Class is the kind of an event. The zero value, ClassNothing, marks an
// event as consumed; views call (*Event).Clear when they handle one.
type Class uint16

// ClassNothing is the consumed marker. All other classes are bit flags so
// callers can build masks (e.g. an event filter for mouse events only).
const ClassNothing Class = 0

// Event class flags.
const (
	ClassKey Class = 1 << iota
	ClassMouseDown
	ClassMouseUp
	ClassMouseMove
	ClassMouseAuto
	ClassMouseWheel
	ClassCommand
	ClassBroadcast
	ClassIdle
	ClassResize
)

// Event is the union type carried through the event loop. Only the field
// indicated by What is meaningful; the others are zero. Events are passed by
// pointer so handlers can mark them consumed via Clear.
type Event struct {
	What  Class
	Key   KeyEvent
	Mouse MouseEvent
	Msg   MessageEvent
}

// Clear marks the event as consumed. The convention matches Turbo Vision: a
// view that handles an event sets What back to ClassNothing, and dispatch
// loops stop at that point.
func (e *Event) Clear() { e.What = ClassNothing }

// Consumed reports whether the event has already been handled.
func (e *Event) Consumed() bool { return e.What == ClassNothing }
