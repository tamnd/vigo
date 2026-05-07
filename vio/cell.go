package vio

// Cell is a single screen position: one rune plus its visual attributes.
// Wide characters (e.g. East Asian double-width or emoji) occupy two
// consecutive cells; the trailing cell carries Rune == 0 to mark continuation.
type Cell struct {
	Rune rune
	Attr Attr
}
