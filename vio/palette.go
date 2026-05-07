package vio

// Palette is an ordered list of attributes. Views look up attributes by a
// 1-indexed slot number to decouple their drawing code from concrete colors;
// themes can swap palettes at runtime without touching any view's Draw.
//
// Slot 0 is conventionally reserved for an "unset" attribute (transparent
// fall-through). The first usable slot is 1.
type Palette []Attr

// Map returns the attribute at slot `i` (1-indexed). Out-of-range slots fall
// back to the first attribute in the palette, or to Attr{} on an empty
// palette, silently, to match Turbo Vision's tolerant indexing scheme.
func (p Palette) Map(i byte) Attr {
	if int(i) <= 0 || int(i) > len(p) {
		if len(p) == 0 {
			return Attr{}
		}
		return p[0]
	}
	return p[int(i)-1]
}

// BorlandClassic is the default Turbo Vision palette (subset; full set is
// fleshed out as widgets ship in v0.2). Indices map 1..N as Turbo Vision
// did, leaving room for theme overlays.
//
//nolint:gochecknoglobals // theme tables are immutable singletons.
var BorlandClassic = Palette{
	{Fg: LightGray, Bg: Blue},     // 1  desktop fill
	{Fg: Black, Bg: LightGray},    // 2  menu normal
	{Fg: DarkGray, Bg: LightGray}, // 3  menu disabled
	{Fg: Black, Bg: Green},        // 4  menu selected
	{Fg: LightRed, Bg: LightGray}, // 5  menu shortcut
	{Fg: Black, Bg: LightGray},    // 6  status normal
	{Fg: DarkGray, Bg: LightGray}, // 7  status disabled
	{Fg: White, Bg: Black},        // 8  status selected
	{Fg: LightRed, Bg: LightGray}, // 9  status shortcut
	{Fg: LightGray, Bg: Blue},     // 10 window frame inactive
	{Fg: Yellow, Bg: Blue},        // 11 window frame active
	{Fg: White, Bg: Blue},         // 12 window text
	{Fg: DarkGray, Bg: Black},     // 13 shadow fill
	{Fg: Black, Bg: LightGray},    // 14 dialog text (static, label)
	{Fg: Black, Bg: LightGray},    // 15 dialog frame inactive
	{Fg: White, Bg: LightGray},    // 16 dialog frame active
	{Fg: LightRed, Bg: LightGray}, // 17 dialog shortcut
	{Fg: Black, Bg: Green},        // 18 button normal
	{Fg: White, Bg: Green},        // 19 button focused
	{Fg: DarkGray, Bg: LightGray}, // 20 button disabled
	{Fg: Yellow, Bg: Green},       // 21 button shortcut
	{Fg: White, Bg: Blue},         // 22 input line normal
	{Fg: White, Bg: Cyan},         // 23 input line focused
	{Fg: Black, Bg: LightGray},    // 24 input line selection
	{Fg: Black, Bg: Cyan},         // 25 cluster normal
	{Fg: Black, Bg: Green},        // 26 cluster focused
	{Fg: Yellow, Bg: Green},       // 27 cluster selected
	{Fg: LightRed, Bg: Cyan},      // 28 cluster shortcut
	{Fg: Yellow, Bg: Cyan},        // 29 scrollbar page
	{Fg: Black, Bg: Green},        // 30 scrollbar arrows
	{Fg: Black, Bg: Cyan},         // 31 list normal
	{Fg: White, Bg: Green},        // 32 list focused
	{Fg: White, Bg: Cyan},         // 33 list selected
	{Fg: DarkGray, Bg: Cyan},      // 34 list disabled
}
