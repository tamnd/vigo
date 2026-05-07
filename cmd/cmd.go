// Package cmd implements vigo's command bus: a CommandSet for grouping
// command IDs, an Enabler that tracks which commands are currently
// active, and a key-binding table that maps KeyEvents to commands.
//
// Views consult cmd.IsEnabled to decide their dim state; bindings
// resolve framework keys (F1 = CmdHelp) and host-defined ones (Ctrl-S
// = save). Hosts toggle commands via Enable/Disable; the Enabler
// notifies subscribers so menus can repaint.
package cmd

import (
	"github.com/tamnd/vigo/event"
)

// Set is a collection of command IDs. Use it to flip many commands
// on or off in one call. The zero value is an empty set.
type Set map[event.CommandID]struct{}

// NewSet returns a set seeded with ids.
func NewSet(ids ...event.CommandID) Set {
	s := make(Set, len(ids))
	for _, id := range ids {
		s[id] = struct{}{}
	}
	return s
}

// Add adds id to the set.
func (s Set) Add(id event.CommandID) { s[id] = struct{}{} }

// Remove removes id from the set.
func (s Set) Remove(id event.CommandID) { delete(s, id) }

// Has reports whether id is in the set.
func (s Set) Has(id event.CommandID) bool {
	_, ok := s[id]
	return ok
}

// Enabler tracks which commands are enabled. It also keeps a list of
// subscribers that get notified when the enabled set changes; the
// menu uses that signal to dim or undim items without polling.
type Enabler struct {
	enabled     Set
	subscribers []func()
}

// NewEnabler returns an Enabler with the given initial set enabled.
// Pass nil for an empty starting set.
func NewEnabler(initial Set) *Enabler {
	en := &Enabler{enabled: make(Set)}
	for id := range initial {
		en.enabled[id] = struct{}{}
	}
	return en
}

// IsEnabled reports whether id is currently enabled.
func (e *Enabler) IsEnabled(id event.CommandID) bool { return e.enabled.Has(id) }

// Enable turns each id in s on. Subscribers are notified once after
// the bulk update if any state actually changed.
func (e *Enabler) Enable(s Set) {
	changed := false
	for id := range s {
		if !e.enabled.Has(id) {
			e.enabled.Add(id)
			changed = true
		}
	}
	if changed {
		e.notify()
	}
}

// Disable turns each id in s off. Subscribers are notified once after
// the bulk update if any state actually changed.
func (e *Enabler) Disable(s Set) {
	changed := false
	for id := range s {
		if e.enabled.Has(id) {
			e.enabled.Remove(id)
			changed = true
		}
	}
	if changed {
		e.notify()
	}
}

// Subscribe registers fn to be called after every Enable/Disable that
// changes state. The returned function removes the subscription.
func (e *Enabler) Subscribe(fn func()) func() {
	e.subscribers = append(e.subscribers, fn)
	idx := len(e.subscribers) - 1
	return func() { e.subscribers[idx] = nil }
}

func (e *Enabler) notify() {
	for _, fn := range e.subscribers {
		if fn != nil {
			fn()
		}
	}
}

// Bindings maps a KeyEvent to a command ID. The map key compares all
// three KeyEvent fields (Key, Mod, Rune) so Alt-X and a bare X do not
// collide.
type Bindings map[event.KeyEvent]event.CommandID

// NewBindings returns an empty binding table.
func NewBindings() Bindings { return make(Bindings) }

// Bind installs k -> id, replacing any prior binding for k.
func (b Bindings) Bind(k event.KeyEvent, id event.CommandID) { b[k] = id }

// Unbind removes the binding for k, if any.
func (b Bindings) Unbind(k event.KeyEvent) { delete(b, k) }

// Lookup returns the command bound to k, or CmdNone if none.
func (b Bindings) Lookup(k event.KeyEvent) event.CommandID {
	if id, ok := b[k]; ok {
		return id
	}
	return event.CmdNone
}
