// Package view defines vigo's base view types: View (the visual
// primitive) and Group (a container of views that owns event dispatch
// and composition). Concrete views embed *View and override the
// methods they need.
package view

import "github.com/tamnd/vigo/internal/event"

// State is a bitmask of view state flags. The set roughly matches
// Borland Turbo Vision's sfXxx constants.
type State uint16

// State flags.
const (
	StateVisible State = 1 << iota
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

// Options is a bitmask of options that affect how the framework
// dispatches events to a view and how the view participates in its
// owning group.
type Options uint16

// Options flags.
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

// GrowMode describes how a view's bounds change when its owner is
// resized. It mirrors Turbo Vision's gfXxx constants.
type GrowMode uint8

// GrowMode flags.
const (
	GrowLoX GrowMode = 1 << iota
	GrowLoY
	GrowHiX
	GrowHiY
)

// GrowAll follows the parent in all four directions; GrowFixed sticks
// at the original bounds.
const (
	GrowFixed GrowMode = 0
	GrowAll            = GrowLoX | GrowLoY | GrowHiX | GrowHiY
)

// EventMaskDefault is the default class mask: keys, primary mouse
// transitions, commands, and broadcasts.
const EventMaskDefault = event.ClassKey |
	event.ClassMouseDown |
	event.ClassMouseUp |
	event.ClassCommand |
	event.ClassBroadcast
