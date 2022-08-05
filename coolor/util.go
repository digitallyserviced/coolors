package coolor

import (
	"fmt"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
)

type Severity int
type Status struct {
	Severity Severity
	Message  string
}

const (
	Unknown Severity = iota
	Refresh
	Healthy
	Warning
	Alert
)

func Init() {
	rand.Seed(time.Now().UnixMilli())
}
func GenerateRandomColors(count int) []tcell.Color {
	tcols := make([]tcell.Color, count)
	for i := range tcols {
		tcols[i] = *MakeRandomColor()
	}
	return tcols
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
	r, g, b := col.RGB()
	if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
		return tcell.ColorBlack
	}
	return tcell.ColorWhite
}

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
		MainC.menu.ActivateSelected(cc)
	}
}

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

var (
	colorRegs = []string{rgb, rgba, hsl, hsla, hsv, hsva, hex3, hex6} //
)

// vim: ts=2 sw=2 et ft=go
