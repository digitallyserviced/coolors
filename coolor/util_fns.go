package coolor

import (
	"log"
	"os"
	"runtime"
	"text/template"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
)

func doLog(args... interface{}) {
  log.Printf("%v", args)
}

func doCallers() {
  st := make([]uintptr, 10)
  n := runtime.Callers(1, st)
  st = st[:n]
  f := runtime.CallersFrames(st)
  for {
    frame, more := f.Next()
    dump.P(frame.Line, frame.Func.Name())
    if !more {
      break
    }
  }
}
func setupLogging() func() error {
	f, _ := os.OpenFile("dumps", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	// f, _ := os.OpenFile(os.DevNull, os.O_RDWR|os.O_APPEND, 0666)

	log.SetOutput(f)

	return f.Close
}

func MakeDebugDump(tp tview.Primitive) {
	dump.P(tp)
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
	switch kp {
	case tcell.KeyEnter:
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

func inverseColor(col tcell.Color) tcell.Color {
	r, g, b := col.RGB()
	return tcell.NewRGBColor(255-r, 255-g, 255-b)
}

type HookDrawInfo struct {
	x, y, width, height   int
	centerY, lowerCenterY int
}

type HookDrawFunctions struct {
	Target *tview.Primitive
	Wrap   DrawFunction
	Chain  DrawFunctionChain
}

type (
	DrawFunction      func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int)
	DrawFunctionChain []*DrawFunction
)

func DrawFunctionDispatcher(
	p *tview.Primitive,
	dfc DrawFunctionChain,
) DrawFunction {
	return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		// for _, v := range dfc {
		// 	// if v != nil {
		// 	// }
		// }
		return x, y, width, height
	}
}

func (cc *CoolorColor) DrawHook(df *DrawFunction) {
}

func (hdf *HookDrawFunctions) CoolorColorStatusText(
	p tview.Primitive,
	screen tcell.Screen,
	x int,
	y int,
	width int,
	height int,
) (int, int, int, int) {
	cc, ok := p.(*CoolorColor)
	if !ok {
		return x, y, width, height
	}
	centerY := y + height/2
	lowerCenterY := centerY + centerY/2
	for cx := x + 1; cx < x+width-1; cx++ {
		screen.SetContent(
			cx,
			centerY+(height/3),
			tview.BoxDrawingsLightHorizontal,
			nil,
			tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(*cc.Color),
		)
	}

	status_tpl := MakeTemplate("color_status", `
      {{define "locked"}}{{- if locked -}}  {{- else -}}  {{- end -}}{{- end -}}
      {{define "selected"}}{{- if selected -}}   {{- end -}}{{- end -}}
    `, template.FuncMap{
		"locked":   cc.GetLocked,
		"selected": cc.GetSelected,
		"dirty":    cc.GetDirty,
		"css":      cc.GetColor,
	})
	sel := status_tpl(`{{- template "selected" . -}}`, cc)
	lock := status_tpl(`{{- template "locked" . -}}`, cc)
	txtColor := cc.GetFgColor()
	tview.Print(screen, sel, x+1, centerY, width-2, tview.AlignCenter, txtColor)
	tview.Print(
		screen,
		lock,
		x+1,
		lowerCenterY,
		width-2,
		tview.AlignCenter,
		txtColor,
	)

	return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
}

func CenteredStrikeText() {
}

// vim: ts=2 sw=2 et ft=go
