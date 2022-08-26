package coolor

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	// "github.com/digitallyserviced/coolors/status"
	// "github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	"github.com/samber/lo"
)

type (
	Severity int
	Status   struct {
		Severity Severity
		Message  string
	}
)

const (
	Unknown Severity = iota
	Refresh
	Healthy
	Warning
	Alert
)

func GenerateRandomColors(count int) []tcell.Color {
	tcols := make([]tcell.Color, count)
	for i := range tcols {
		tcols[i] = *MakeRandomColor()
	}
	return tcols
}

// func generateColorWithinDistance(tcol Color, maxDistance float64, cs *ColorStream) {
// 	count := 0
// 	for {
// 		if cs.Cancel == nil || cs.OutColors == nil {
// 			return
// 		}
// 		select {
// 		case <-cs.Context.Done():
// 			return
// 		default:
// 			count += 1
// 			col := genRandomSeededColor(cs)
// 			if col != nil {
// 				cs.OutColors <- col
// 			}
// 		}
// 		time.Sleep(2 * time.Millisecond)
// 	}
// }

func init() {
	rand.Seed(time.Now().UnixMilli())
}

// func DrainColorStream(colors *ColorStream, count int, status func(s string)) []*Color {
// 	// defer ShutdownStream(colors)
//   originalCount := count
//   defer func(){
// 			status(fmt.Sprintf("Generated %d colors from %d iterations", originalCount, colors.Count))
//   }()
// 	colors.Start <- struct{}{}
// 	outColors := make([]*Color, 0)
// 	for {
//     if colors.OutColors == nil || colors.Status == nil || colors.generatedColors == nil {
//       return outColors
//     }
// 		select {
// 		case <-colors.Context.Done():
//       return outColors
// 		case color, ok := <-colors.OutColors:
// 			if !ok {
//         return outColors
// 			}
// 			outColors = append(outColors, color)
// 			if color != nil {
// 				count--
// 			}
// 		case stat, ok := <-colors.Status:
// 			if !ok {
// 				// colors.Cancel()
//         return outColors
// 			}
// 			status(stat)
// 		default:
// 			time.Sleep(5 * time.Millisecond)
// 		}
//
// 		if count == 0 {
//       return outColors
// 		}
// 	}
// }

// func genRandomSeededColor(cs *ColorStream) *Color {
// 	rand.Seed(time.Now().UnixNano())
// 	tcol2 := MakeColorFromTcell(randomColor())
// 	return &tcol2
// }
func genRandomSeededColor() interface{} {
	rand.Seed(time.Now().UnixNano())
	tcol2 := MakeColorFromTcell(randomColor())
	return tcol2
}

func checkColorDistance(tcol, tcol2 Color, distance float64) bool {
	if tcol2.DistanceCIEDE2000(tcol) <= distance {
		return true
	}
	return false
}

type ColorStream struct {
	OutColors       <-chan interface{}
	Start           chan struct{}
	Status          *ColorStreamProgress
	Cancel          context.CancelFunc
	Generator       func() interface{}
	Validator       func(interface{}) bool
	Context         context.Context
}

func takeN(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

type ColorCount uint32
type ColorStreamIterationProgressHandler interface {
  OnItr(uint32)
}
type ColorStreamValidProgressHandler interface {
  OnValid(uint32)
}
type ColorStreamProgressHandler interface {
  ColorStreamIterationProgressHandler
  ColorStreamValidProgressHandler
}
type NilProgressHandler struct {}
func (NilProgressHandler) OnItr(i uint32) {}
func (NilProgressHandler) OnValid(v uint32) {}


func NewNilProgressHandler() ColorStreamProgressHandler {
  nph := &NilProgressHandler{}
  return nph
}
type FunctionalProgressHandler struct {
	v func(uint32)
  i func(uint32)
}

func (fph FunctionalProgressHandler) OnItr(i uint32) {
  fph.i(i)
}
func (fph FunctionalProgressHandler) OnValid(v uint32) {
  fph.v(v)
}

func NewProgressHandler(v func(uint32), i func(uint32)) *FunctionalProgressHandler {
  nph := &FunctionalProgressHandler{
  	v: v,
  	i: i,
  }
  return nph
}

type ColorStreamProgress struct {
	Valid, Itr uint32
  ProgressHandler ColorStreamProgressHandler
	// onValid    func(v uint32)
	// onItr      func(i uint32)
}


func NewColorStreamProgress() *ColorStreamProgress {
	csp := &ColorStreamProgress{
		Valid:   0,
		Itr:     0,
    ProgressHandler: NewNilProgressHandler(),
	}
	return csp
}

func (csp *ColorStreamProgress) SetProgressHandler(csph ColorStreamProgressHandler) {
  csp.ProgressHandler = csph
}

func (csp *ColorStreamProgress) GetValid() uint32 {
	return atomic.LoadUint32(&csp.Valid)
}

func (csp *ColorStreamProgress) GetItr() uint32 {
	return atomic.LoadUint32(&csp.Itr)
}

func (csp *ColorStreamProgress) Itrd() {
  res := atomic.AddUint32(&csp.Itr, 1)
  csp.ProgressHandler.OnItr(res)
}

func (csp *ColorStreamProgress) Validd() {
  res := atomic.AddUint32(&csp.Itr, 1)
  csp.ProgressHandler.OnValid(res)
}

func takeFn(done <-chan struct{}, csp *ColorStreamProgress, valueStream <-chan interface{}, fn func(interface{}) bool) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for {
			select {
			case <-done:
				return
			case v := <-valueStream:
        csp.Itrd()
				if fn(v) {
          csp.Validd()
					takeStream <- v
				}
			}
		}
	}()
	return takeStream
}

func asStream(done <-chan struct{}, fn func() interface{}) <-chan interface{} {
	s := make(chan interface{})
	go func() {
		defer close(s)

		for {
			select {
			case <-done:
				return
			case s <- fn():
			}
		}
	}()
	return s
}

func fanIn(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chans))

		for _, c := range chans {
			go func(c <-chan interface{}) {
				for v := range c {
					out <- v
				}
				wg.Done()
			}(c)
		}

		wg.Wait()
		close(out)
	}()
	return out
}

func (cs *ColorStream) Run(done <-chan struct{}) {
	numRoutines := 4
	generators := make([]<-chan interface{}, 0)
	for i := 0; i < numRoutines; i++ {
		generators = append(generators, takeFn(done, cs.Status, asStream(done, genRandomSeededColor), cs.Validator))
	}
	cs.OutColors = fanIn(generators...)
}

func TakeNColors(done <-chan struct{}, valueStream <-chan interface{}, num int) []Color {
	colors := make([]Color, 0)
	for cv := range takeN(done, valueStream, num) {
		col := cv.(Color)
		colors = append(colors, col)
	}
	return colors
}

// func (cs *ColorStream) Run() {
//   defer func(){
//     if err := recover(); err != nil {
//       // fmt.Println(fmt.Errorf("%v", err))
//     }
//   }()
// 	numRoutines := 4
// 	go func() {
//   defer close(cs.generatedColors)
//   defer close(cs.OutColors)
//   defer close(cs.Status)
// 		<-cs.Start
//     go func(){
//       <-time.After(5 * time.Second)
//       cs.Cancel()
//     }()
// 		go cs.Validate()
// 		for i := 0; i < numRoutines; i++ {
// 			go cs.Generate()
// 		}
// 		<-cs.Context.Done()
// 	}()
// }

// func (cs *ColorStream) Validate() {
//   defer func(){
//     if err := recover(); err != nil {
//       // fmt.Println(fmt.Errorf("%v", err))
//     }
//   }()
// 	ticker := time.NewTicker(50 * time.Millisecond)
//   defer ticker.Stop()
// 	colorCount := 0
// 	for {
// 		select {
// 		case <-cs.Context.Done():
//       return
// 		case color := <-cs.generatedColors:
// 			cs.Count++
// 			if cs.Validator(color, cs) {
// 				colorCount++
// 				cs.OutColors <- color
// 			}
// 		case <-ticker.C:
//       select {
//       case <-cs.Context.Done():
//         return
//       case cs.Status <- fmt.Sprintf("%d / %d", colorCount, cs.Count):
//       }
//
// 		}
// 	}
// }

// func (cs *ColorStream) Generate() {
//   defer func(){
//     if err := recover(); err != nil {
//       // fmt.Println(fmt.Errorf("%v", err))
//     }
//   }()
// 	for {
// 		select {
// 		case <-cs.Context.Done():
//       return
// 		// case cs.generatedColors <- cs.Generator():
// 			// color := cs.Generator(cs)
// 			// cs.generatedColors <- color
// 		}
// 			time.Sleep(5 * time.Millisecond)
// 	}
// }

func StartColorStream(g func() interface{}, v func(interface{}) bool) *ColorStream {
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	cs := &ColorStream{
		OutColors: make(<-chan interface{}),
		Start:     make(chan struct{}),
		Status:    NewColorStreamProgress(),
		Cancel:    cancel,
		Generator: func() interface{} {
			return &Color{0, 0, 0}
		},
		Validator: func(interface{}) bool {
			return true
		},
		Context: ctx,
	}

	if g != nil {
		cs.Generator = g
	}

	if v != nil {
		cs.Validator = v
	}

	return cs
}

func RandomShadesStream(tcol Color, maxDistance float64) *ColorStream { // , cs *ColorStream
	cs := StartColorStream(genRandomSeededColor, func(c interface{}) bool {
    defer func(){
      if err := recover(); err != nil {
        // fmt.Println("%v %v", c, err)
      }
    }()
		return checkColorDistance(tcol, c.(Color), maxDistance)
	})
	return cs
}

func ShutdownStream(cs *ColorStream) {
	cs.Cancel()
}

func RandomNamedAnsiClusterShade(tcol Color, distance float64) Color {
	clusterColor := lo.Sample[*CoolorColorCluster](getNamedAnsiColors())
	return RandomShadeFromCluster(tcol, clusterColor, distance)
}

func RandomAnsiClusterShade(tcol Color, distance float64) Color {
	clusterColor := lo.Sample[*CoolorColorCluster](getBaseAnsiClusterColors())
	return RandomShadeFromCluster(tcol, clusterColor, distance)
}

func RandomShadeFromCluster(tcol Color, cluster *CoolorColorCluster, distance float64) Color {
	col2 := cluster.leadColor
	return RandomShadeFromColors(tcol, col2, distance)
}

func RandomShadeFromColors(tcol, tcol2 Color, distance float64) Color {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := (distance * r.Float64()) + 0.01
	return tcol.BlendLuv(tcol2, max)
}

func MakeRandomColor() *tcell.Color {
	col := tcell.NewRGBColor(int32(randRange(0, 255)), int32(randRange(0, 255)), int32(randRange(0, 255)))
	return &col
}

func randRange(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

func randomColor() tcell.Color {
	r := int32(randRange(0, 255))
	g := int32(randRange(0, 255))
	b := int32(randRange(0, 255))
	return tcell.NewRGBColor(r, g, b)
}

func getFGColor(col tcell.Color) tcell.Color {
	cc := NewIntCoolorColor(inverseColor(col).Hex())
	c, ok := MakeColor(cc)
	// dump.P(cc.TerminalPreview())
	if ok {
		r, g, b := c.LinearRgb()
		if (255*float64(r)*0.299 + 255*float64(g)*0.587 + 255*float64(b)*0.114) > 150 {
			// if (255*float64(r)*0.2926 + 255*float64(g)*0.5152 + 255*float64(b)*0.1722) > 150 {
			// if (float64(r)*0.2926 + float64(g)*0.5152 + float64(b)*0.1722) > 150 {
			return tcell.ColorBlack
		}
		return tcell.ColorWhite
	}
	return tcell.ColorBlack
	// r, g, b := cc.RGB
	// if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
}

// func getFGColor(col tcell.Color) tcell.Color {
// 	r, g, b := col.RGB()
// 	if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
// 		return tcell.ColorBlack
// 	}
// 	return tcell.ColorWhite
// }

func inverseColor(col tcell.Color) tcell.Color {
	r, g, b := col.RGB()
	return tcell.NewRGBColor(255-r, 255-g, 255-b)
}

func MakeTemplate(name, tpl string, funcMap template.FuncMap) func(s string, data interface{}) string {
	status_tpl := template.New(name)
	status_tpl.Funcs(funcMap)

	status_tpl.Parse(tpl)

	return func(s string, data interface{}) string {
		out := &strings.Builder{}
		ntpl, ok := template.Must(status_tpl.Clone()).Parse(s)
		if ok != nil {
			fmt.Println(fmt.Errorf("%s", ok))
		}
		ntpl.Execute(out, data)
		return out.String()
	}
}

func MakeDebugDump(tp tview.Primitive) {
	dump.P(tp)
}

func AddFlexItem(fl *tview.Flex, tp tview.Primitive, f, p int) {
	fl.AddItem(tp, f, p, false)
}

type PalettePaddle struct {
	*tview.Box
	icon, iconActive string
	status           string
}

func NewPalettePaddle(icon, iconActive string) *PalettePaddle {
	nb := tview.NewBox()
	pp := &PalettePaddle{
		Box:        nb,
		icon:       icon,
		iconActive: iconActive,
		status:     "active",
	}
	nb.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	nb.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		iconColor := tview.Styles.ContrastBackgroundColor
		icon := pp.icon
		if pp.status == "enabled" {
			iconColor = tview.Styles.ContrastBackgroundColor
			icon = pp.iconActive
		} else if pp.status == "disabled" {
			iconColor = tview.Styles.PrimitiveBackgroundColor
		} else {
			iconColor = tview.Styles.MoreContrastBackgroundColor
		}
		centerX := x + (width / 2)
		centerY := y + (height / 2)
		tview.Print(screen, icon, centerX-1, centerY-1, 1, tview.AlignCenter, iconColor)
		return x, y, width, height
	})
	return pp
}

func (pp *PalettePaddle) SetStatus(status string) {
	pp.status = status
	// MainC.app.Sync()
}

func MakeBoxItem(title, col string) *tview.Box {
	nb := tview.NewBox().SetBorder(true)
	if title == "" {
		nb.SetBorder(false)
	} else {
		nb.SetBorder(true).SetTitle(title)
	}
	if col != "" {
		return nb.SetBackgroundColor(tcell.GetColor(col))
	}
	return nb.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
}

func MakeSpace(fl *tview.Flex, title, col string, f, p int) *tview.Box {
	// rc := randomColor()
	spc := MakeBoxItem(title, col)
	AddFlexItem(fl, spc, f, p)
	return spc
}

func BlankSpace(fl *tview.Flex) *tview.Box {
	return MakeSpace(fl, "", "", 0, 1)
}

func MakeSpacer(fl *tview.Flex) *tview.Box {
	rc := randomColor()
	spc := MakeBoxItem(" ", fmt.Sprintf("%06x", rc.Hex())).SetBackgroundColor(tcell.ColorBlack)
	AddFlexItem(fl, spc, 0, 1)
	return spc
}

func DrawCenteredLine(txt string, screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	centerY := y + height/2
	lowerCenterY := centerY + centerY/3
	for cx := x + 1; cx < x+width-1; cx++ {
		screen.SetContent(cx, lowerCenterY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(getFGColor(tview.Styles.ContrastBackgroundColor)))
	}
	tview.Print(screen, fmt.Sprintf(" %s ", txt), x+1, lowerCenterY, width-2, tview.AlignCenter, getFGColor(tview.Styles.ContrastBackgroundColor))
	return x, y, width, height
}

func MakeCenterLineSpacer(fl *tview.Flex) (*tview.Box, func(string)) {
	spc := MakeSpace(fl, "", "", 0, 1).SetBackgroundColor(tcell.ColorBlack)
	AddFlexItem(fl, spc, 0, 1)
	ctrtxt := ""
	spc.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		return DrawCenteredLine(ctrtxt, screen, x, y, width, height)
	})
	return spc, func(txt string) {
		ctrtxt = txt
	}
}

type Navigable interface {
	NavSelection(int)
}

type Activatable interface {
	ActivateSelected()
}

type Selectable interface {
	GetSelected() int
}

type CoolorSelectable interface {
	GetSelected() (*CoolorColor, int)
}
type VimNavSelectable interface {
	GetSelectedVimNav() VimNav
}

type VimNav interface {
	Navigable
}

func HandleVimNavigableHorizontal(vm VimNav, ch rune, kp tcell.Key) {
	switch {
	case ch == 'h' || kp == tcell.KeyLeft:
		vm.NavSelection(-1)
	case ch == 'l' || kp == tcell.KeyRight:
		vm.NavSelection(1)
	}
}

func HandleVimNavigableVertical(vm VimNav, ch rune, kp tcell.Key) {
	switch {
	case ch == 'j' || kp == tcell.KeyDown:
		vm.NavSelection(-1)
	case ch == 'k' || kp == tcell.KeyUp:
		vm.NavSelection(1)
	}
}

func HandleVimNavSelectable(s VimNavSelectable) VimNav {
	return s.GetSelectedVimNav()
}

func HandleSelectable(s Selectable) int {
	return s.GetSelected()
}

// HandleVimNavSelectable
func HandleCoolorSelectable(s CoolorSelectable, ch rune, kp tcell.Key) {
	_ = kp
	switch {
	case kp == tcell.KeyEnter:
		cc, _ := s.GetSelected()
    if cc == nil {
      return
    }
    if MainC.menu == nil {
      return
    }
		MainC.menu.ActivateSelected(cc)
	}
}

const (
	// 0;#090300;1;#db2d20;2;#01a252;3;#fded02;4;#01a0e4;5;#a16a94;6;#b5e4f4;7;#a5a2a2;8;#5c5855;9;#e8bbd0;10;#3a3432;11;#4a4543;12;#807d7c;13;#d6d5d4;14;#cdab53;15;#f7f7f7
	set4BitDynamicColors string = "\033]4;%s\007"
	dynamicColorIndex    string = "%d;%s"
	// printf "\033]10;#4a4543;#f7f7f7;#4a4543\007"
	setTextFgBgCursor string = "\033]10;%s\007"
	// printf "\033]17;#a5a2a2\007"
	setBgColor string = "\033]17;%s\007"
	// printf "\033]19;#4a4543\007"
	setSelectionFgColor string = "\033]19;%s\007"
	// printf "\033]5;0;#4a4543\007"
	setDynamicColorBold string = "\033]5;%d;%s\007"
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
)

var colorRegs = []string{rgb, rgba, hsl, hsla, hsv, hsva, hex3, hex6} //

// vim: ts=2 sw=2 et ft=go
