package coolor

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/digitallyserviced/coolors/theme"
	// "github.com/gookit/goutil/dump"
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

func MakeTemplate(
	name, tpl string,
	funcMap template.FuncMap,
) func(s string, data interface{}) string {
	status_tpl := template.New(name)
	status_tpl.Funcs(funcMap)

	status_tpl.Parse(tpl)

	return func(s string, data interface{}) string {
		out := &strings.Builder{}
		ntpl, ok := template.Must(status_tpl.Clone()).Parse(s)
		if ok != nil {
			// fmt.Println(fmt.Errorf("%s", ok))
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
	nb.SetDrawFunc(
		func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
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
      halfY := (height / 2) / 2
			tview.Print(
				screen,
				icon,
				centerX,
				centerY-1-halfY,
				2,
				tview.AlignCenter,
				iconColor,
			)
			tview.Print(
				screen,
				icon,
				centerX,
				centerY-1,
				2,
				tview.AlignCenter,
				iconColor,
			)
			tview.Print(
				screen,
				icon,
				centerX,
				centerY-1+halfY,
				2,
				tview.AlignCenter,
				iconColor,
			)
			return x, y, width, height
		},
	)
	return pp
}

func (pp *PalettePaddle) SetStatus(status string) {
	pp.status = status
	// MainC.app.Sync()
}

type CoolorColorTag struct {
	Tag   *TagItem
	Color *CoolorColor
	*tview.Box
	fm *tview.FocusManager
	// *tview.Frame
}

func (cct *CoolorColorTag) SetColor(cc *Coolor) *CoolorColorTag {
  if cc == nil {
    cct.Box.SetBackgroundColor(theme.GetTheme().GrayerBackground)
    cct.Color = FromTcell(theme.GetTheme().GrayerBackground).Escalate()
    // cct.Color = MakeColor
    return cct
  }
	cct.Color = cc.Escalate()
	cct.Box.SetBorderColor(cct.Color.GetFgColorShade())
	cct.Box.SetBorderFocusColor(cct.Color.GetFgColor())
	cct.Box.SetBackgroundColor(*cct.Color.Color)
	return cct
}

var bigNums []string = []string{
	`█▀█`,
	`▌█▐`,
	`▌▀▐`,
	`▀█`,
	` █`,
	` █`,
	`▀▀`,
	`▀ `,
	` ▀`,
	`▀▀`,
	`▀ `,
	`▀ `,
	`▀██`,
	` ▀▐`,
	`██▐`,
	`▀▀`,
	` ▀`,
	`▀ `,
	`▀▀█`,
	` ▀█`,
	` ▀▐`,
	`▀▀`,
	`█ `,
	`▌▐`,
	`▀▀█`,
	` ▀▐`,
	` ▀▐`,
	`█▀▀`,
	`▌▀ `,
	`██ `,

	// `┌┐`,
	// `││`,
	// `└┘`,
	// `┐ `,
	// `│ `,
	// `┴ `,
	// `┌┐`,
	// `┌┘`,
	// `└┘`,
	// `┌┐`,
	// ` ┤`,
	// `└┘`,
	// `┬ `,
	// `└┤`,
	// ` ┴`,
	// `┌┐`,
	// `└┐`,
	// `└┘`,
	// `┌┐`,
	// `├┐`,
	// `└┘`,
	// `┌┐`,
	// ` ┤`,
	// ` ┴`,
	// `┌┐`,
	// `├┤`,
	// `└┘`,
	// `┌┐`,
	// `└┤`,
	// ` ┴`,
}

func NewCoolorColorTagBox(ti *TagItem, dynIdx int) *CoolorColorTag {
	tb := tview.NewBox()
	// tb.SetDontClear(true)
	cct := &CoolorColorTag{
		Tag: ti,
		Box: tb,
		fm: tview.NewFocusManager(func(p tview.Primitive) {
			MainC.app.SetFocus(p)
		}),
		// Frame: tview.NewFrame(tb),
	}
	tb.SetDrawFunc(
		func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			if dynIdx < 0 {
				return x, y, width, height
			}
			strNum := fmt.Sprintf("%d", dynIdx)
			lines := make([]string, 0)

			for i := 0; i < 3; i++ {
				line := make([]string, 0)
				for _, v := range strNum {
					idx, err := strconv.ParseInt(string(v), 10, 8)
					// fmt.Println(idx)
					checkErr(err)
					line = append(line, bigNums[(int(idx)*3)+i])
				}
				if len(line) > 0 {
					lines = append(lines, strings.Join(line, "█"))
				}
			}
			// nums := strings.Join(lines, "\n")
			// fmt.Println(nums)
			x, y, width, height = tb.GetInnerRect()
			// dump.P(x, y, width, height, lines, strNum, len(lines[0]))
			centerY := y // height/2
			// xOff := (width) - (len(lines[0]) / 2)
			// centerX :=  clamp(float64(x + xOff), float64(x), float64((x + width) - len(lines[0])))
			fmtr := "[#%06x:#%06x:rb]%s[-::-]"
      fade := 0.45
      // fmt.Printf(fmtr, cct.Color.GetFgColorFade(fade).Hex(),cct.Color.Color.Hex(),"SHIT BOY")
			for i, v := range lines {
      colString := fmt.Sprintf(fmtr, cct.Color.GetFgColorFade(fade).Hex(),cct.Color.Color.Hex(),v) // ,cct.Color.GetFgColorShade().Hex()
				tview.Print(
					screen,
          colString,
					x, // x + (len(v) / 2),
					centerY+i,
					width,
					tview.AlignCenter,
					cct.Color.GetFgColorFade(fade),
				)

			}
			return x, y, width, height
		},
	)
  cct.Box.SetBorder(false).SetBorderPadding(0,0,0,0)
	// cct.Frame.SetBorder(true).SetBorderPadding(1,1,1,1)
	cct.Box.SetBorderVisible(true)
	// cct.Frame.SetBorders(0,0,0,0,0,0)
	// cc.SetPlain(true)
	return cct
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
	spc := MakeBoxItem(
		" ",
		fmt.Sprintf("%06x", rc.Hex()),
	).SetBackgroundColor(tcell.ColorBlack)
	AddFlexItem(fl, spc, 0, 1)
	return spc
}

func DrawCenteredLine(
	txt string,
	screen tcell.Screen,
	x int,
	y int,
	width int,
	height int,
) (int, int, int, int) {
	centerY := y + height/2
	lowerCenterY := centerY + centerY/3
	fg := NewIntCoolorColor(tview.Styles.ContrastBackgroundColor.Hex())
	for cx := x + 1; cx < x+width-1; cx++ {
		screen.SetContent(
			cx,
			lowerCenterY,
			tview.BoxDrawingsLightHorizontal,
			nil,
			tcell.StyleDefault.Foreground(*fg.Color),
		)
	}
	tview.Print(
		screen,
		fmt.Sprintf(" %s ", txt),
		x+1,
		lowerCenterY,
		width-2,
		tview.AlignCenter,
		*fg.Color,
	)
	return x, y, width, height
}

func MakeCenterLineSpacer(fl *tview.Flex) (*tview.Box, func(string)) {
	spc := MakeSpace(fl, "", "", 0, 1).SetBackgroundColor(tcell.ColorBlack)
	AddFlexItem(fl, spc, 0, 1)
	ctrtxt := ""
	spc.SetDrawFunc(
		func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			return DrawCenteredLine(ctrtxt, screen, x, y, width, height)
		},
	)
	return spc, func(txt string) {
		ctrtxt = txt
	}
}

// DotPalette
