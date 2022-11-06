package coolor

import (
	"fmt"
	// log "github.com/sirupsen/logrus"
	"math/rand"
	"regexp"
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"

	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/zzlog"

	// "github.com/gookit/goutil/errorx"
	msgpack "github.com/vmihailenco/msgpack/v5"
	// "github.com/digitallyserviced/coolors/coolor/zzlog"
)

type (
	Severity int
	Status   struct {
		Message  string
		Severity Severity
	}
)

const (
	Unknown Severity = iota
	Refresh
	Healthy
	Warning
	Alert
)

const (
	cssInteger       = "[-\\+]?\\d+%?"
	cssNumber        = "[-\\+]?\\d*\\.\\d+%?"
	cssUnit          = "(?:" + cssNumber + ")|(?:" + cssInteger + ")"
	permissiveMatch3 = "[\\s|\\(]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")\\s*\\)?"
	permissiveMatch4 = "[\\s|\\(]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")\\s*\\)?"
	rgb              = "rgb" + permissiveMatch3
	rgba             = "RGBA" + permissiveMatch4
	hsl              = "hsl" + permissiveMatch3
	hsla             = "hsla" + permissiveMatch4
	hsv              = "hsv" + permissiveMatch3
	hsva             = "hsva" + permissiveMatch4
	// hex3             = `#?([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})`
	// hex6             = `#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})`
	// hex4             = `#?([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})`
	// hex8             = `#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})`
	hex3 = `(#[0-9a-fA-F]{3})\b`
	hex6 = `(#[0-9a-fA-F]{6})\b`
	// 0;#090300;1;#db2d20;2;#01a252;3;#fded02;4;#01a0e4;5;#a16a94;6;#b5e4f4;7;#a5a2a2;8;#5c5855;9;#e8bbd0;10;#3a3432;11;#4a4543;12;#807d7c;13;#d6d5d4;14;#cdab53;15;#f7f7f7
	// printf "\033]10;#4a4543;#f7f7f7;#4a4543\007"
	// printf "\033]17;#a5a2a2\007"
	// printf "\033]19;#4a4543\007"
	// printf "\033]5;0;#4a4543\007"
	set4BitDynamicColors string = "\033]4;%s\007"
	dynamicColorIndex    string = "%d;%s"
	setTextFgBgCursor    string = "\033]10;%s\007"
	setBgColor           string = "\033]17;%s\007"
	setSelectionFgColor  string = "\033]19;%s\007"
	setDynamicColorBold  string = "\033]5;%d;%s\007"
)

type ColorRepParserFunc func(cr *ColorRep) *Color

// Color representations hsl(n,n,n), #7bafcd, rgba(0.6,0.6,0.6,255)
// assume only need at most 4 value to make it in the color model/space
type ColorRep struct {
	found          map[string]int
	parseRegex     *regexp.Regexp
	scanFormat     string
	v1, v2, v3, v4 float64
	parseFunc      ColorRepParserFunc
}

// var CssHex6 ColorRep
var (
	CssHex6 *ColorRep
)
func init() {
	CssHex6 = NewColorRep(hex6, "%s", func(cr *ColorRep) *Color {
		return &Color{cr.v1 / 255.0, cr.v2 / 255.0, cr.v3 / 255.0}
	})
}

func (cr *ColorRep) ParseSingle(s string) *Color {
	n, err := fmt.Sscanf(s, cr.scanFormat, &cr.v1, cr.v2, cr.v3, cr.v4)
	if err != nil {
		panic(err)
	}
	if n > 0 {
		return cr.parseFunc(cr)
	}
	return nil
}

func (cr *ColorRep) Matches(s string) bool {
	found := cr.parseRegex.Match([]byte(s))
	dump.P(s, found)
	return found
}

func NewColorRep(regex, scanFormat string, p ColorRepParserFunc) *ColorRep {
	rx := regexp.MustCompile(regex)
	cr := &ColorRep{
		found:      make(map[string]int),
		parseRegex: rx,
		scanFormat: scanFormat,
		v1:         0,
		v2:         0,
		v3:         0,
		v4:         0,
		parseFunc:  p,
	}
	return cr
}

func StringColorizer(s string) (string, *CoolorColorsPalette) {
	for _, v := range []*ColorRep{CssHex6} {
		if v.Matches(s) {
			content, cols := v.FindAndColorize(s)
			if cols.Len() > 0 {
				return content, cols
			}
		}
	}
	return "", nil
}

func (cr *ColorRep) FindAndColorize(sc string) (string, *CoolorColorsPalette) {
	cols := NewCoolorColorsPalette()
	newsc := cr.parseRegex.ReplaceAllStringFunc(sc, func(s string) string {
		col := NewCoolorColor(s)
		colnum, exists := cr.found[s]
		if exists {
			cr.found[s] = colnum + 1
		} else {
			cols.AddCoolorColor(col)
			cr.found[s] = 1
		}
		return col.TVCSSString(false)
	})
	return newsc, cols
	// if matchIdxs := cr.parseRegex.FindAllStringSubmatchIndex(sc, -1); len(matchIdxs) > 0 {
	// 	for _, c := range matchIdxs {
	//      var cuint int32 = 0
	//      str := sc[c[0]:c[1]]
	//      // fmt.Printf("%q", str)
	//      // n, err := fmt.Sscanf(strings.TrimSpace(str), "%s", &cuint)
	//      // if err != nil || n == 0{
	//      //   panic(fmt.Errorf("%v %v %v %v", c, n, err, str))
	//      // }
	//      // tcol := tcell.NewHexColor(cuint)
	//      // str = NewIntCoolorColor(cuint).TVPreview()
	//      newsc =
	// 		if len(c) == 2 {
	// 			colors = append(colors, str)
	// 		}
	// 	}
	//    dump.P(sc)
	//    return colors
	// }
	// if match := cr.parseRegex.FindAllStringSubmatch(sc, -1); match != nil {
	// 	colors := make([]string, 0)
	// 	for _, c := range match {
	// 		if len(c) == 2 {
	// 			colors = append(colors, c[1])
	// 		}
	//      dump.P(c)
	// 	}
	//    return colors
	// }
	// regexp.MustCompile(reg).FindAllSubmatch()
}



var colorRegexes = []string{rgb, rgba, hsl, hsla, hsv, hsva, hex3, hex6} //
func genSimilarHslColor(
	tcol Color,
	f func(r, h, s, l float64) Color,
) interface{} {
	rand.Seed(time.Now().UnixNano())
	h, s, l := tcol.Hsl()
	adjust := rand.Float64() * 360
	return f(adjust, h, s, l)
}

func checkHslColorDistance(tcol, tcol2 Color, distance float64) bool {
	return tcol2.DistanceRgb(tcol) <= distance
}

func GetColorName(col tcell.Color) string {
	for n, v := range tcell.ColorNames {
		if col == v {
			return n
		}
	}
	return ""
}

func ErrorAssert[R any](v R, e error) R {
	iserr := func(v interface{}) {
		if v == nil {
			return
		}
		e, ok := v.(error)
		if ok {
      zlog.Error(fmt.Sprintf("%T %v", v, e), zzlog.String("msg", e.Error()))
			// doLog(errorx.WithPrevf(e, "%V", v))
			panic(e)
		}
		// return
	}
	iserr(e)
	return v
}

func checkErrX(err error, vars ...interface{}) bool {
	if err != nil {
    // zlog.Error(msg string, fields ...log.Field)
    zlog.Error(fmt.Sprintf("checkErr: %T", err), zzlog.String("msg", err.Error()))
		// doLog(errorx.WithPrevf(errorx.Traced(err), "%T %s", err, vars))
		return false
	}
	return true
}

func checkErr(err error) {
	if err != nil {
    zlog.Error(fmt.Sprintf("checkErr: %T", err), zzlog.String("msg", err.Error()))
		panic(err)
	}
}

func startBoltStats() {
	go func() {
		// Grab the initial stats.
		// prev := Store.Bolt().Stats()
		tick := time.NewTicker(10 * time.Second)
		tick.Reset(1000 * time.Millisecond)

		for {
			// Wait for 10s.
			select {
			case <-tick.C:
				TrimSeentCoolors(RecentCoolorsMax)
				// stats := Store.Bolt().Stats()
				// diff := stats.Sub(&prev)

				// Encode stats to JSON and print to STDERR.
				// json.NewEncoder(os.Stderr).Encode(diff)

				// Save stats for the next loop.
				// prev = stats
			}
		}
	}()
}

func SeentColor(from string, cc *CoolorColor, src Referenced) {
	// fmt.Println(from, cc)
	if MainC == nil || MainC.EventNotifier == nil {
		return
	}
	MainC.EventNotifier.Notify(
		*MainC.EventNotifier.NewObservableEvent(ColorSeentEvent, from, cc, src),
	)
}

func IfElse[T any](b bool, a, c T) T {
	if b {
		return a
	}
	return c
}

var IfElseStr = IfElse[string]
var IfElseUint64 = IfElse[uint64]
var IfElseTCol = IfElse[*tcell.Color]
var IfElseCCol = IfElse[*CoolorColor]

// custom msgpack decoding function for bolthold (faster than gobs)
func dec(data []byte, value interface{}) error {
	return msgpack.Unmarshal(data, value)
}

// custom msgpack encoding function for bolthold (faster than gobs)
func enc(value interface{}) ([]byte, error) {
	return msgpack.Marshal(value)
}

// var _ msgpack.CustomEncoder = (*CoolorColorPalette)(nil)
// var _ msgpack.CustomDecoder = (*CoolorColorPalette)(nil)
//
// Debug message categories.
// const (
// 	LogConn  = LogMask(1 << iota)   // Connection events
// 	LogState                        // State changes
// 	LogCmd                          // Command execution
// 	LogRaw                          // Raw data stream excluding literals
// 	LogGo                           // Goroutine execution
// 	LogAll   = LogMask(1<<iota - 1) // All messages
// 	LogNone  = LogMask(0)           // No messages
// )

// var logMasks = []enumName{
// 	{uint32(LogAll), "LogAll"},
// 	{uint32(LogConn), "LogConn"},
// 	{uint32(LogState), "LogState"},
// 	{uint32(LogCmd), "LogCmd"},
// 	{uint32(LogRaw), "LogRaw"},
// 	{uint32(LogNone), "LogNone"},
// }
//
// func (v LogMask) String() string   { return enumString(uint32(v), logMasks, false) }
// func (v LogMask) GoString() string { return enumString(uint32(v), logMasks, true) }


var Shortcuts []*Shortcut
var scopes []*Scope

func addShortcut(
	identifier, name string,
	scope *Scope,
	event *tcell.EventKey,
) *Shortcut {
	shortcut := &Shortcut{
		Identifier:   identifier,
		Name:         name,
		Scope:        scope,
		Event:        event,
		defaultEvent: event,
	}

	Shortcuts = append(Shortcuts, shortcut)

	return shortcut
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
	// Parent is this scopes upper Scope, which may be null, in case this is a
	// root scope.
	Parent *Scope

	// Identifier will be used for persistence and should never change
	Identifier string

	// Name will be shown on the UI
	Name string
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
	// Identifier will be used for persistence and should never change
	Identifier string

	// Name will be shown on the UI
	Name string

	// The Scope will be omitted, as this needed be persisted anyway.
	Scope *Scope

	// Event is the shortcut expressed as it's resulting tcell Event.
	Event *tcell.EventKey

	//This shortcuts default, in order to be able to reset it.
	defaultEvent *tcell.EventKey
}

// Equals compares the given EventKey with the Shortcuts Event.
func (shortcut *Shortcut) Equals(event *tcell.EventKey) bool {
	return EventsEqual(shortcut.Event, event)
}

var (
	globalScope = addScope("global", "Application wide", nil)
	FocusUp     = addShortcut("focus_up", "Focus the next widget above",
		globalScope, tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModShift))
	FocusDown = addShortcut("focus_down", "Focus the next widget below",
		globalScope, tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModShift))
	FocusLeft = addShortcut("focus_left", "Focus the next widget to the left",
		globalScope, tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModShift))
	FocusRight = addShortcut("focus_right", "Focus the next widget to the right",
		globalScope, tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModShift))
)

func DirectionalFocusHandling(
	event *tcell.EventKey,
	app *tview.Application,
) *tcell.EventKey {
	focused := app.GetFocus()

	if FocusUp.Equals(event) {
		FocusNextIfPossible(tview.Up, app, focused)
	} else if FocusDown.Equals(event) {
		FocusNextIfPossible(tview.Down, app, focused)
	} else if FocusLeft.Equals(event) {
		FocusNextIfPossible(tview.Left, app, focused)
	} else if FocusRight.Equals(event) {
		FocusNextIfPossible(tview.Right, app, focused)
	} else {
		return event
	}
	return nil
}

func addScope(identifier, name string, parent *Scope) *Scope {
	scope := &Scope{
		Parent:     parent,
		Identifier: identifier,
		Name:       name,
	}

	scopes = append(scopes, scope)

	return scope
}

func FocusNextIfPossible(
	direction tview.FocusDirection,
	app *tview.Application,
	focused tview.Primitive,
) {
	if focused == nil {
		return
	}

  zlog.Debug("focus_next", zzlog.Int("dir", int(direction)))
	focusNext := focused.NextFocusableComponent(direction)
  // zlog.Debug("focus_next_orig", log.Int("dir", int(direction)))
	// zlog.Printf("%v %T -> %T", direction, focused, focusNext)
	if focusNext != nil {
		app.SetFocus(focusNext)
	}
}

func debounce[T any](min time.Duration, max time.Duration, input <-chan T) chan T {
	output := make(chan T)

	go func() {
		var (
			buffer   T
			ok       bool
			minTimer <-chan time.Time
			maxTimer <-chan time.Time
		)

		// Start debouncing
		for {
			select {
			case buffer, ok = <-input:
				if !ok {
					return
				}
				minTimer = time.After(min)
				if maxTimer == nil {
					maxTimer = time.After(max)
				}
			case <-minTimer:
				minTimer, maxTimer = nil, nil
				output <- buffer
			case <-maxTimer:
				minTimer, maxTimer = nil, nil
				output <- buffer
			}
		}
	}()

	return output
}
//
// type PrimitiveFrameHandler func(t time.Duration, step int)
// type AnimatedUpdater func(p PrimitiveActor)

// type Animated interface {
//   Animate()
// Update()
// Start()
// Begin()
// OnTick(delta time.Duration)
// SetDuration(t time.Duration)
// SetTickCount(t int)
// SetFrameHandler(f PrimitiveFrameHandler)
// Stop()
// End()
// }
//
// type PrimitiveActor interface {
// 	tview.Primitive
// 	// Animated
// }
// type animationTimer func(timeStep float64) float64
// type animationDraw func(frameIdx int, timeStep float64, scr tcell.Screen, p interface{})
//
// type Animator struct {
// 	animating bool
// 	frameIdx  int
// 	start     time.Time
// 	duration  time.Duration
// 	timeStep  float64
// 	timeFunc  animationTimer
// 	draw      animationDraw
// }
//
// type Animation interface {
// 	WithOptions(f func(o interface{})) *Animation
// 	WithDuration(d time.Duration) *Animation
// 	Start(p tview.Primitive)
// 	Update(frameIdx int, timeStep float64, scr tcell.Screen, p interface{})
// }
//
// type Rect struct {
// 	x, y, width, height int
// }
//
// type BoxRectMover struct {
// 	originalPos, currentPos Rect
// 	targetPos               Rect
//   Animator
// }
//
// func OffsetMoverX(offset int) *Animator {
//   brm := &BoxRectMover{
//   	originalPos: Rect{
//   		x:      0,
//   		y:      0,
//   		width:  0,
//   		height: 0,
//   	},
//   	currentPos:  Rect{},
//   	targetPos:   Rect{},
//   }
//   timer := func(frameIdx int, timeStep float64, scr tcell.Screen, p interface{}) {
// 	//
// 	}
//
//
//   a := NewAnimator(timer, drawFn)
// 	// return 
// }
//
//
// // Start implements Animation
// func (brm *BoxRectMover) Start(p tview.Primitive, offset int) {
//
// }
//
// // Update implements Animation
// func (brm *BoxRectMover) Update(frameIdx int, timeStep float64, scr tcell.Screen, p interface{}) {
// 	panic("unimplemented")
// }
//
// // WithDuration implements Animation
// func (brm *BoxRectMover) WithDuration(d time.Duration) *Animation {
// 	panic("unimplemented")
// }
//
// // WithOptions implements Animation
// func (brm *BoxRectMover) WithOptions(f func(o interface{})) *Animation {
// 	panic("unimplemented")
// }
//
// func (a *Animator) Start(p interface{}, frameIdx int) {
//
// 	var brm Animation = BoxRectMover{}
//
// }
//
// func AnimateTestBox() {
// 	a := NewAnimator(MakeBoxItem("WHEEE", "#7f81aa"), func(timeStep float64) float64 {
// 		return timeStep
// 	}, OffsetMoverX(20))
// }
//
// func NewAnimator(duration time.Duration, ease animationTimer, drawFunc animationDraw) *Animator {
// 	a := &Animator{
// 		start:    time.Now(),
// 		duration: duration,
// 		timeStep: 0,
// 		timeFunc: ease,
// 		draw:     drawFunc,
// 	}
//
// 	return a
// }
//
// type AnimatorModifier func(p interface{})
//
// type AnimatedBox struct {
// 	*tview.Box
// 	// Animators []AnimateCallback
// }
//
// vim: ts=2 sw=2 et ft=go
