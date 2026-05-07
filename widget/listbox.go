package widget

import (
	"github.com/tamnd/vigo/vio"
)

// ListBox is a ListViewer backed by a string slice. Items can be
// updated at any time with SetItems; the list resizes accordingly.
type ListBox struct {
	*ListViewer

	items []string
}

// NewListBox returns a list box sized to bounds with the given items
// and an optional vertical scrollbar. Pass a nil scrollbar to omit it.
func NewListBox(bounds vio.Rect, items []string, vScroll *ScrollBar) *ListBox {
	l := &ListBox{
		ListViewer: NewListViewer(bounds, vScroll, nil),
		items:      items,
	}
	l.GetText = l.itemText
	l.SetRange(len(items))
	return l
}

// Items returns the current item slice. The slice is shared with the
// list box; callers must not modify it concurrently with redraws.
func (l *ListBox) Items() []string { return l.items }

// SetItems replaces the items and resets focus to the first row.
func (l *ListBox) SetItems(items []string) {
	l.items = items
	l.Focused = 0
	l.SetRange(len(items))
}

// Selected returns the focused item's string, or "" when the list is
// empty.
func (l *ListBox) Selected() string {
	if l.Focused < 0 || l.Focused >= len(l.items) {
		return ""
	}
	return l.items[l.Focused]
}

func (l *ListBox) itemText(item, _ int) string {
	if item < 0 || item >= len(l.items) {
		return ""
	}
	return l.items[item]
}
