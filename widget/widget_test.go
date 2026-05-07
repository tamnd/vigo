package widget

import (
	"github.com/tamnd/vigo/view"
	"github.com/tamnd/vigo/vio"
)

// newPalettedHost returns a Group with the BorlandClassic palette so widget
// tests resolve MapColor without panicking.
func newPalettedHost() *view.Group {
	g := view.NewGroup(vio.R(0, 0, 80, 24))
	g.Palette = vio.BorlandClassic
	return g
}
