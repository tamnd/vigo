package vio

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gdamore/tcell/v2"

	"github.com/tamnd/vigo/event"
)

// NewTcellScreen returns a Screen backed by github.com/gdamore/tcell/v2.
func NewTcellScreen() Screen { return &tcellScreen{} }

type tcellScreen struct {
	s        tcell.Screen
	events   chan event.Event
	done     chan struct{}
	pumpDone chan struct{}
	closed   atomic.Bool

	cursorMu sync.Mutex
	cx, cy   int
	cvis     bool
}

func (t *tcellScreen) Init() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("vio: new screen: %w", err)
	}
	if err := s.Init(); err != nil {
		return fmt.Errorf("vio: init screen: %w", err)
	}
	s.SetStyle(tcell.StyleDefault)
	s.EnableMouse()
	s.EnablePaste()
	s.HideCursor()
	s.Clear()

	t.s = s
	t.events = make(chan event.Event, 64)
	t.done = make(chan struct{})
	t.pumpDone = make(chan struct{})
	go t.pump()
	return nil
}

func (t *tcellScreen) Fini() {
	if !t.closed.CompareAndSwap(false, true) {
		return
	}
	if t.done != nil {
		close(t.done)
	}
	if t.s != nil {
		// PostEvent with nil interrupts the blocked PollEvent.
		_ = t.s.PostEvent(nil)
	}
	if t.pumpDone != nil {
		<-t.pumpDone
	}
	if t.s != nil {
		t.s.Fini()
	}
}

func (t *tcellScreen) Size() (int, int) {
	if t.s == nil {
		return 0, 0
	}
	return t.s.Size()
}

func (t *tcellScreen) Events() <-chan event.Event { return t.events }

func (t *tcellScreen) SetCursor(x, y int, visible bool) {
	t.cursorMu.Lock()
	t.cx, t.cy, t.cvis = x, y, visible
	t.cursorMu.Unlock()
}

func (t *tcellScreen) Show(surf *Surface) error {
	if t.s == nil {
		return errScreenNotInit
	}
	w, h := surf.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := surf.At(x, y)
			r := c.Rune
			if r == 0 {
				r = ' '
			}
			t.s.SetContent(x, y, r, nil, c.Attr.Style())
		}
	}
	t.cursorMu.Lock()
	cx, cy, cvis := t.cx, t.cy, t.cvis
	t.cursorMu.Unlock()
	if cvis {
		t.s.ShowCursor(cx, cy)
	} else {
		t.s.HideCursor()
	}
	t.s.Show()
	return nil
}

func (t *tcellScreen) pump() {
	defer close(t.pumpDone)
	defer close(t.events)
	for {
		ev := t.s.PollEvent()
		if ev == nil {
			return
		}
		out, ok := translate(ev)
		if !ok {
			continue
		}
		select {
		case t.events <- out:
		case <-t.done:
			return
		}
	}
}

// translate converts a tcell.Event into a vigo event.Event.
func translate(ev tcell.Event) (event.Event, bool) {
	switch e := ev.(type) {
	case *tcell.EventKey:
		return translateKey(e), true
	case *tcell.EventMouse:
		return translateMouse(e), true
	case *tcell.EventResize:
		w, h := e.Size()
		return event.Event{
			What: event.ClassResize,
			Mouse: event.MouseEvent{
				X: w, Y: h,
			},
		}, true
	case *tcell.EventPaste:
		// In v0.1 we ignore paste boundaries; full handling lands with the
		// editor in v0.3.
		return event.Event{}, false
	default:
		return event.Event{}, false
	}
}

func translateKey(e *tcell.EventKey) event.Event {
	mod := event.ModNone
	if e.Modifiers()&tcell.ModShift != 0 {
		mod |= event.ModShift
	}
	if e.Modifiers()&tcell.ModCtrl != 0 {
		mod |= event.ModCtrl
	}
	if e.Modifiers()&tcell.ModAlt != 0 {
		mod |= event.ModAlt
	}
	if e.Modifiers()&tcell.ModMeta != 0 {
		mod |= event.ModMeta
	}

	k, r := translateTcellKey(e)
	return event.Event{
		What: event.ClassKey,
		Key: event.KeyEvent{
			Key:  k,
			Mod:  mod,
			Rune: r,
		},
	}
}

func translateTcellKey(e *tcell.EventKey) (event.Key, rune) {
	switch e.Key() {
	case tcell.KeyRune:
		return event.KeyRune, e.Rune()
	case tcell.KeyEnter:
		return event.KeyEnter, 0
	case tcell.KeyEscape:
		return event.KeyEsc, 0
	case tcell.KeyTab:
		return event.KeyTab, 0
	case tcell.KeyBacktab:
		return event.KeyShiftTab, 0
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		return event.KeyBackspace, 0
	case tcell.KeyDelete:
		return event.KeyDelete, 0
	case tcell.KeyInsert:
		return event.KeyInsert, 0
	case tcell.KeyHome:
		return event.KeyHome, 0
	case tcell.KeyEnd:
		return event.KeyEnd, 0
	case tcell.KeyPgUp:
		return event.KeyPgUp, 0
	case tcell.KeyPgDn:
		return event.KeyPgDn, 0
	case tcell.KeyUp:
		return event.KeyArrowUp, 0
	case tcell.KeyDown:
		return event.KeyArrowDown, 0
	case tcell.KeyLeft:
		return event.KeyArrowLeft, 0
	case tcell.KeyRight:
		return event.KeyArrowRight, 0
	case tcell.KeyF1:
		return event.KeyF1, 0
	case tcell.KeyF2:
		return event.KeyF2, 0
	case tcell.KeyF3:
		return event.KeyF3, 0
	case tcell.KeyF4:
		return event.KeyF4, 0
	case tcell.KeyF5:
		return event.KeyF5, 0
	case tcell.KeyF6:
		return event.KeyF6, 0
	case tcell.KeyF7:
		return event.KeyF7, 0
	case tcell.KeyF8:
		return event.KeyF8, 0
	case tcell.KeyF9:
		return event.KeyF9, 0
	case tcell.KeyF10:
		return event.KeyF10, 0
	case tcell.KeyF11:
		return event.KeyF11, 0
	case tcell.KeyF12:
		return event.KeyF12, 0
	}
	if e.Key() >= tcell.KeyCtrlA && e.Key() <= tcell.KeyCtrlZ {
		// Translate Ctrl-letter to KeyRune with ModCtrl folded in by the
		// caller. The rune is the lowercase letter for ergonomics.
		return event.KeyRune, rune(e.Key()-tcell.KeyCtrlA) + 'a'
	}
	return event.KeyNone, 0
}

func translateMouse(e *tcell.EventMouse) event.Event {
	x, y := e.Position()
	btn := e.Buttons()
	mod := event.ModNone
	if e.Modifiers()&tcell.ModShift != 0 {
		mod |= event.ModShift
	}
	if e.Modifiers()&tcell.ModCtrl != 0 {
		mod |= event.ModCtrl
	}
	if e.Modifiers()&tcell.ModAlt != 0 {
		mod |= event.ModAlt
	}
	var buttons event.MouseButton
	if btn&tcell.Button1 != 0 {
		buttons |= event.MouseLeft
	}
	if btn&tcell.Button2 != 0 {
		buttons |= event.MouseMiddle
	}
	if btn&tcell.Button3 != 0 {
		buttons |= event.MouseRight
	}
	if btn&tcell.WheelUp != 0 {
		buttons |= event.MouseWheelUp
	}
	if btn&tcell.WheelDown != 0 {
		buttons |= event.MouseWheelDown
	}

	var what event.Class
	switch {
	case buttons&(event.MouseWheelUp|event.MouseWheelDown) != 0:
		what = event.ClassMouseWheel
	case buttons != 0:
		what = event.ClassMouseDown
	default:
		what = event.ClassMouseUp
	}
	return event.Event{
		What: what,
		Mouse: event.MouseEvent{
			X: x, Y: y, Buttons: buttons, Mod: mod,
		},
	}
}
