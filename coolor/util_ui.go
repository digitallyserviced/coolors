package coolor

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
)

type PalettePaddle struct {
	*tview.Box
	icon, iconActive string
	status           string
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

func AddFlexItem(fl *tview.Flex, tp tview.Primitive, f, p int) {
	fl.AddItem(tp, f, p, false)
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
		switch pp.status {
		case "enabled":
			iconColor = tview.Styles.ContrastBackgroundColor
			icon = pp.iconActive
		case "disabled":
			iconColor = tview.Styles.PrimitiveBackgroundColor
		default:
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
    fg := NewIntCoolorColor(tview.Styles.ContrastBackgroundColor.Hex())
	for cx := x + 1; cx < x+width-1; cx++ {
		screen.SetContent(cx, lowerCenterY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(*fg.Color))
	}
	tview.Print(screen, fmt.Sprintf(" %s ", txt), x+1, lowerCenterY, width-2, tview.AlignCenter, *fg.Color)
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
// DotPalette
