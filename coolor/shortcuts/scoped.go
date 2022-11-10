package shortcuts

import (
	"fmt"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	// "github.com/gookit/goutil/dump"
)

var (
	Shortcuts []*Shortcut
	scopes    []*Scope
)

func NewScope(identifier, name string, parent *Scope, args ...string) *Scope {
	scope := &Scope{
		Parent:     parent,
		Identifier: identifier,
		Name:       name,
		Shortcuts:  make([]*Shortcut, 0),
	}

  if len(args) > 0 {
    scope.Icon = args[0]
  }

	scopes = append(scopes, scope)

	return scope
}

// EventsEqual compares the given events, respecting everything except for the
// When field.
func EventsEqual(eventOne, eventTwo *tcell.EventKey) bool {
	if (eventOne == nil && eventTwo != nil) ||
		(eventOne != nil && eventTwo == nil) {
		return false
	}

	return eventOne.Rune() == eventTwo.Rune() &&
		eventOne.Modifiers() == eventTwo.Modifiers() &&
		eventOne.Key() == eventTwo.Key()
}

// Scope is what describes a shortcuts scope within the application. Usually
// a scope can only have a specific shortcut once and a children scope will
// overwrite that shortcut, since that lower scope has the upper hand.
type Scope struct {
	Handler    ShortcutsHandler
	Parent     *Scope
	Identifier string
	Name       string
	Icon       string
	Shortcuts  []*Shortcut
}

// ShortcutDataRepresentation represents a shortcut configured by the user.
// This prevents redundancy of scopes.
type ShortcutDataRepresentation struct {
	Identifier      string
	ScopeIdentifier string
	EventKey        tcell.Key
	EventMod        tcell.ModMask
	EventRune       rune
}

// Shortcut defines a shortcut within the application. The scope might for
// example be a widget or situation in which the user is.
type Shortcut struct {
	Handler      ShortcutsHandler
	Scope        *Scope
	Event        *tcell.EventKey
	defaultEvent *tcell.EventKey
	Callback     ShortcutCallback
	Identifier   string
	Name         string
	Icon         string
}

type ShortcutCallback func(i ...interface{}) bool

type ShortcutsHandler interface {
	// SetupKeys(sh *ShortcutsHandler)
	GetScope() *Scope
	// HandleShortcuts(e *tcell.EventKey)
}

func (scope *Scope) FindScope(name string) *Scope {
  for _, v := range scopes {
    if v.Name == name {
      return v
    }
  }
  return nil
}

func (scope *Scope) NewShortcut(
	identifier, name string,
	event *tcell.EventKey,
	cb ShortcutCallback,
  args ...string,
) *Shortcut {
	shortcut := &Shortcut{
		Identifier:   identifier,
		Name:         name,
		Callback:     cb,
		Scope:        scope,
		Event:        event,
		defaultEvent: event,
		Handler:      scope.Handler,
	}

  if len(args) > 0 {
    shortcut.Icon = args[0]
  }

	Shortcuts = append(Shortcuts, shortcut)
	scope.Shortcuts = append(scope.Shortcuts, shortcut)

	return shortcut
}

type DefaultInputHandler func(event *tcell.EventKey, setFocus func(p tview.Primitive))

func (sh *Scope) HandleShortcuts(e *tcell.EventKey, setFocus func(p tview.Primitive)) {
  // dump.P(sh.Shortcuts)
  fmt.Println(e)
	for _, v := range sh.Shortcuts {
    // fmt.Println(v.Name)
		if v.Equals(e) {
			if !v.Callback() {
				return
			}
		}
	}

	if sh.Parent != nil {
		for _, v := range sh.Parent.Shortcuts {
    // fmt.Println(v.Name)
			if v.Equals(e) {
				if !v.Callback() {
					return
				}
			}
		}
	}
}

// Equals compares the given EventKey with the Shortcuts Event.
func (shortcut *Shortcut) Equals(event *tcell.EventKey) bool {
	return EventsEqual(shortcut.Event, event)
}

var GlobalScope = NewScope("global", "Application wide", nil) // favScope.NewSho
// FocusUp     = NewShortcut("focus_up", "Focus the next widget above",
// 	globalScope, tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModShift))
// FocusDown = NewShortcut("focus_down", "Focus the next widget below",
// 	globalScope, tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModShift))
// FocusLeft = NewShortcut("focus_left", "Focus the next widget to the left",
// 	globalScope, tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModShift))
// FocusRight = NewShortcut("focus_right", "Focus the next widget to the right",
// 	globalScope, tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModShift))

func DirectionalFocusHandling(
	event *tcell.EventKey,
	app *tview.Application,
) *tcell.EventKey {
	// focused := app.GetFocus()

	// if FocusUp.Equals(event) {
	// 	FocusNextIfPossible(tview.Up, app, focused)
	// } else if FocusDown.Equals(event) {
	// 	FocusNextIfPossible(tview.Down, app, focused)
	// } else if FocusLeft.Equals(event) {
	// 	FocusNextIfPossible(tview.Left, app, focused)
	// } else if FocusRight.Equals(event) {
	// 	FocusNextIfPossible(tview.Right, app, focused)
	// } else {
	// 	return event
	// }
	return nil
}

func FocusNextIfPossible(
	direction tview.FocusDirection,
	app *tview.Application,
	focused tview.Primitive,
) {
	if focused == nil {
		return
	}

	// zlog.Debug("focus_next", zzlog.Int("dir", int(direction)))
	focusNext := focused.NextFocusableComponent(direction)
	// zlog.Debug("focus_next_orig", log.Int("dir", int(direction)))
	// zlog.Printf("%v %T -> %T", direction, focused, focusNext)
	if focusNext != nil {
		app.SetFocus(focusNext)
	}
}
