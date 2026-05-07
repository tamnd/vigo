// Package widget collects the small leaf views that live inside a Window
// or Dialog: static text, labels, buttons, input lines, clusters,
// scrollbars, and so on.
//
// v0.2 phase 3 ships the first batch (StaticText, Label, Button,
// InputLine). Lists, scrollbars, and clusters arrive in phase 4.
package widget

// Palette slot constants used by widgets. They line up with the slots
// added to vio.BorlandClassic for dialog content. Widgets default to
// these but callers may override PaletteIndex if a host theme adds
// extra slots.
const (
	SlotDialogText     byte = 14
	SlotDialogShortcut byte = 17
	SlotButtonNormal   byte = 18
	SlotButtonFocused  byte = 19
	SlotButtonDisabled byte = 20
	SlotButtonShortcut byte = 21
	SlotInputNormal    byte = 22
	SlotInputFocused   byte = 23
	SlotInputSelection byte = 24
	SlotClusterNormal  byte = 25
	SlotClusterFocused byte = 26
	SlotClusterSelect  byte = 27
	SlotClusterHotKey  byte = 28
	SlotScrollPage     byte = 29
	SlotScrollArrow    byte = 30
	SlotListNormal     byte = 31
	SlotListFocused    byte = 32
	SlotListSelected   byte = 33
	SlotListDisabled   byte = 34
)
