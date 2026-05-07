// Package help is vigo's contextual-help slot. Every view carries a
// HelpCtx; pressing F1 should open the topic registered for the
// focused view's context. v0.2 ships only the wiring and an About
// dialog; the real topic browser arrives in v0.7.
package help

import (
	"github.com/tamnd/vigo/dialog"
	"github.com/tamnd/vigo/event"
	"github.com/tamnd/vigo/vio"
	"github.com/tamnd/vigo/widget"
)

// HelpCtx identifies a help topic. Zero (CtxNone) means the view has
// no specific topic; F1 then falls back to the About dialog.
//
//nolint:revive // HelpCtx mirrors the spec name and the Turbo Vision API.
type HelpCtx uint16

// Reserved help contexts. User contexts start at CtxUser.
const (
	CtxNone  HelpCtx = 0
	CtxAbout HelpCtx = 1

	CtxUser HelpCtx = 100
)

// Topic is one entry in the registry: a short title plus the body
// text shown in the help viewer. v0.2 only renders Title in the
// About dialog; v0.7 will paint Body in a scrolling Help window.
type Topic struct {
	Title string
	Body  string
}

// Registry maps HelpCtx values to topics. The zero value is ready to
// use; pass it around by pointer so writes are visible everywhere.
type Registry struct {
	topics map[HelpCtx]Topic
}

// NewRegistry returns an empty registry.
func NewRegistry() *Registry {
	return &Registry{topics: map[HelpCtx]Topic{}}
}

// Register stores t under ctx, replacing any prior topic.
func (r *Registry) Register(ctx HelpCtx, t Topic) {
	if r.topics == nil {
		r.topics = map[HelpCtx]Topic{}
	}
	r.topics[ctx] = t
}

// Lookup returns the topic registered for ctx and whether one was
// found. Callers should fall back to CtxAbout when ok is false.
func (r *Registry) Lookup(ctx HelpCtx) (Topic, bool) {
	t, ok := r.topics[ctx]
	return t, ok
}

// Resolve returns the topic for ctx, or the About topic when ctx is
// not registered. If neither is registered it returns the zero Topic.
func (r *Registry) Resolve(ctx HelpCtx) Topic {
	if t, ok := r.topics[ctx]; ok {
		return t
	}
	if t, ok := r.topics[CtxAbout]; ok {
		return t
	}
	return Topic{}
}

// Default returns a registry seeded with the About topic. Hosts use
// it as a starting point and Register additional contexts on top.
func Default() *Registry {
	r := NewRegistry()
	r.Register(CtxAbout, Topic{
		Title: "About",
		Body: "vigo " + Version + "\n" +
			"A Turbo Vision-style IDE in Go.\n" +
			"https://github.com/tamnd/vigo",
	})
	return r
}

// Version is the human-readable build label shown in the About
// dialog. Hosts override it from main if they want a richer string
// (commit hash, build date, etc.).
//
//nolint:gochecknoglobals // single mutable build label
var Version = "v0.2.0"

const (
	aboutW    = 44
	aboutH    = 9
	aboutBtnW = 10
	aboutBtnH = 1
)

// About returns a small modal Dialog displaying the About topic from
// r. The dialog is not inserted anywhere; pass it to
// Application.ExecView (or similar) to show it modally.
func About(r *Registry) *dialog.Dialog {
	if r == nil {
		r = Default()
	}
	t := r.Resolve(CtxAbout)
	if t.Title == "" {
		t = Topic{Title: "About", Body: "vigo " + Version}
	}

	d := dialog.New(vio.R(0, 0, aboutW, aboutH), t.Title)

	body := widget.NewStaticText(
		vio.R(3, 2, aboutW-6, aboutH-5),
		t.Body,
	)
	d.Insert(body)

	ok := widget.NewButton(
		vio.R(0, 0, aboutBtnW, aboutBtnH),
		"~O~K",
		event.CmdOk,
		widget.BfDefault,
	)
	d.PlaceButtons(ok)
	d.SetDefaultButton(ok)
	d.SetCancelButton(ok)
	return d
}
