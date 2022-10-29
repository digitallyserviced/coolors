package coolor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	// "os"

	// "log"

	"reflect"
	"runtime"
	"text/template"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	"github.com/samber/lo"

	// "golang.org/x/exp/constraints"

	// "github.com/digitallyserviced/coolors/coolor/log"
	"github.com/digitallyserviced/coolors/coolor/zzlog"
	xxp "github.com/digitallyserviced/coolors/coolor/xp"
)

var xp = xxp.Xp
var zlog *zzlog.Logger

func doLog(args ...interface{}) {
  
  zlog.Info(fmt.Sprintf("%v", args),lo.Map[interface{}, zzlog.Field](args, func(i1 interface{}, i2 int) zzlog.Field {
    return zzlog.String(fmt.Sprintf("arg_%d", i2),reflect.TypeOf(i1).Name())
  })...)
	// log.Printf("%v", args)
}
func doCallers() {
	st := make([]uintptr, 10)
	n := runtime.Callers(1, st)
	st = st[:n]
	f := runtime.CallersFrames(st)
	for {
		f, more := f.Next()
    dump.P(f.File, f.Function)
		if !more {
			break
		}
	}
}

func setupLogger(){
  var tops = []zzlog.TeeOption{
		{
			Filename: "out.log",
			Ropt: zzlog.RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 1,
			},
			Lef: func(lvl zzlog.Level) bool {
        return true
				// return lvl <= zzlog.InfoLevel
			},
		},
	}

	logger := zzlog.NewTeeWithRotate(tops, zzlog.WithCaller(true),zzlog.AddStacktrace(zzlog.WarnLevel),zzlog.AddCallerSkip(1))
  zlog = logger
	zzlog.ResetDefault(logger)
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

func getColorsFromArray(m []interface{}) (strs []string) {
	for _, v := range m {
		switch v := v.(type) {
		case string:
			strs = append(strs, v)
		case map[string]interface{}:
			mstrs := getColorsFromMap(v)
			strs = append(strs, mstrs...)
		}
	}
	return strs
}

func getColorsFromMap(mapd map[string]interface{}) (strs []string) {
	for _, v := range mapd {
		switch m := v.(type) {
		case string:
			strs = append(strs, m)
		case []string:
			strs = append(strs, m...)
		case []interface{}:
			mstrs := getColorsFromArray(m)
			strs = append(strs, mstrs...)
			// for _, vv := range m {
			// 	str, ok := vv.(string)
			// 	if ok {
			// 		strs = append(strs, str)
			// 	}
			// }
		case map[string]interface{}:
			mstrs := getColorsFromMap(m)
			strs = append(strs, mstrs...)
		}
	}
	// switch m := v.(type) {
	// case map[string]interface{}:
	// 	for _, v := range m {
	// 		mstrs := getColorsFromMap(v)
	// 		strs = append(strs, mstrs...)
	// 	}
	// case []string:
	// 	strs = append(strs, m...)
	// case []interface{}:
	// 	for _, v := range m {
	//      str, ok := v.(string)
	//      if ok {
	//        strs = append(strs, str)
	//      }
	// 	}
	// }
	strs = lo.FilterMap[string, string](strs, func(s string, i int) (string, bool) {
    if len(s) == 0 {
      return "", false
    }
		if []rune(s)[0] == '#' {
			return s, true
		}
		return "", false
	})
	strs = lo.Uniq[string](strs)
	return strs
}

func formatPath(p string) string {
	dir := filepath.Dir(p)
	base := filepath.Base(p)

	home := os.Getenv("HOME")

	if strings.HasPrefix(dir, home) {
		dir = strings.Replace(dir, home, "~", 1)
	}

	if dir == "/" {
		return fmt.Sprintf("[blue]/[normal]%s", base)
	}

	return fmt.Sprintf("[blue]%s/[normal]%s", dir, base)
}

// vim: ts=2 sw=2 et ft=go
