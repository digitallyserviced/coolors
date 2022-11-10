package coolor

import (
	"fmt"
	"math"
	"strconv"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/gookit/color"
	"golang.org/x/exp/constraints"

	"github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/shortcuts"
	"github.com/digitallyserviced/coolors/theme"
)

type CoolorColorInfo struct {
	Frame *tview.Frame
	*tview.Flex
	*CoolorColor
	titleBar *tview.TextView
	infoView *tview.Flex
	gridView *tview.Grid
	details  *CoolorColorDetails
}

type CSSValue struct {
	*tview.TextView
	text string
	idx  int
}

type CoolorColorDetails struct {
	*CoolorColor
	cssValues map[string]*CSSValue
	count     int
	selected  string
  *shortcuts.Scope
}

// GetScope implements shortcuts.ShortcutsHandler
func (ccd *CoolorColorDetails) GetScope() *shortcuts.Scope {
  if ccd.Scope != nil {
    return ccd.Scope
  }
  return nil
  // return ccd.Scope
}

type CoolorColorFloater struct {
	*tview.Flex
	Color *CoolorColorInfo
	Items []*tview.Primitive
}

func NewCoolorColorModal(ncc *CoolorColor) *CoolorColorFloater {
	cc := ncc.Clone()
	cc.SetStatic(true)
	cc.SetPlain(true)
	spf := &CoolorColorFloater{
		Flex:  tview.NewFlex(),
		Color: NewCoolorColorInfo(cc),
	}

	centerFlex := tview.NewFlex()
	centerFlex.SetDirection(tview.FlexRow)
	centerFlex.AddItem(nil, 0, 1, false)
	centerFlex.AddItem(spf.Color.Flex, 12, 0, true)
	centerFlex.AddItem(nil, 0, 1, false)

	spf.SetDirection(tview.FlexColumn)
	spf.AddItem(nil, 0, 2, false)
	spf.AddItem(centerFlex, 0, 5, true)
	spf.AddItem(nil, 0, 2, false)
	return spf
}

func NewCoolorColorInfo(cc *CoolorColor) *CoolorColorInfo {
	cci := &CoolorColorInfo{
		Frame:       &tview.Frame{},
		Flex:        tview.NewFlex(),
		CoolorColor: cc,
		infoView:    tview.NewFlex(),
		gridView:    tview.NewGrid(),
		details:     NewCoolorColorDetails(cc),
	}
	cci.infoView.SetDirection(tview.FlexRow)
	cci.Clear()
	cci.Flex.SetBorder(true)
	cci.AddItem(cci.gridView, 0, 5, true)
	// newFrame.SetBorderColor(tview.Styles.PrimitiveBackgroundColor)
	cci.AddColorFrame()

	cci.AddTitle()
	cci.ColorValues()
	cci.AddColorInfoRows()
	cci.UpdateColor(cc)
	return cci
}

func (cci *CoolorColorInfo) AddColorFrame() {
	newFrame := tview.NewFrame(cci.CoolorColor)
	newFrame.SetBorder(false).SetBorderPadding(0, 1, 1, 1)
	newFrame.SetBorders(0, 0, 1, 1, 0, 0)

	cci.Frame = newFrame
	cci.gridView.AddItem(cci.Frame, 0, 0, 3, 4, 0, 0, false)
}

func NewCoolorColorDetails(cc *CoolorColor) *CoolorColorDetails {
	ccd := &CoolorColorDetails{
		CoolorColor: cc,
		count:       0,
		// cssValues:   make(map[string]string),
		cssValues: make(map[string]*CSSValue),
		// selected:    primaryCssValues[0],
	}
    ccd.Scope = shortcuts.NewScope("color_info", "info panel", nil)
	return ccd
}

func NewCSSValue(text string, idx int) *CSSValue {
	csv := &CSSValue{
		TextView: &tview.TextView{},
		text:     text,
		idx:      idx,
	}
	infoTextView := tview.NewTextView()
	infoTextView.SetWrap(false)
	infoTextView.ScrollTo(0, 0)
	infoTextView.SetDynamicColors(true)
	infoTextView.SetRegions(true)
	infoTextView.SetBorderPadding(0, 0, 0, 0)
	infoTextView.SetBackgroundColor(theme.GetTheme().HeaderBackground)
	infoTextView.SetBorderColor(theme.GetTheme().ContentBackground)
	infoTextView.SetTextAlign(tview.AlignCenter)
	infoTextView.SetText(text)
	infoTextView.ScrollTo(0, 0)
	csv.TextView = infoTextView
	// label := fmt.Sprintf("[green:-:b] %s[-:-:-][red:-:-]%c[-:-:-]", cv, keyLabel[count%(len(keyLabel)-1)])
	return csv
}

func (ccd *CoolorColorDetails) UpdateCssValue(name, fmtStr string) {
	if tv, ok := ccd.cssValues[name]; ok {
		tv.SetText(fmt.Sprintf("[green:-:b]%s[yellow:-:b]([-:-:-] %s [yellow:-:b])[-:-:-][red:-:-]%c[-:-:-]", name, fmtStr, keyLabel[tv.idx%(len(keyLabel)-1)]))
	}
	// ccd.cssValues[name] = NewCSSValue(fmtStr)
}

func (ccd *CoolorColorDetails) AddCssValue(name, fmtStr string) {
	// var sh shortcuts.ShortcutsHandler = ccd
	ccd.cssValues[name] = NewCSSValue(fmtStr, ccd.count)
    ccd.Scope.NewShortcut(
      fmt.Sprintf("copy_%s", name),
      "copy color values",
    NewAltEventKey(tcell.KeyRune, rune(fmt.Sprintf("%d", ccd.count)[0])),
      func(i ...interface{}) bool {
      ev := *events.Global.NewObservableEvent(events.CopyEvent, ccd.cssValues[name].GetText(true), nil, nil)
      ev.Ref = fmt.Sprintf("[yellow]Copied [red]%s[-] [yellow]color values[-]", name)
      // ev.Src = "  " 
      ev.Src = " " 
			events.Global.Notify(ev)
      fmt.Println("You DIDD IT! ", ccd.cssValues[name].idx)
        return true
      },
    )
	ccd.count += 1
}

func (ccf *CoolorColorFloater) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ccf.WrapInputHandler(
		func(event *tcell.EventKey, _ func(p tview.Primitive)) {
			ch := event.Rune()
			// kp := event.Key()

			num, err := strconv.ParseInt(fmt.Sprintf("%c", ch), 10, 8)
			if err != nil {
				return
			}
			if num >= 0 && int(num) < len(selectableCssValues) {
				name := selectableCssValues[num]
				ccf.Color.details.selected = name
				ccf.Color.UpdateView()
			}
			// re , _ := regexp.Compile("[0-9]")
			// if re.MatchString(fmt.Sprintf("%c", ch)) {
			//
			// }
		},
	)
}

func (cci *CoolorColorInfo) UpdateColor(ncc *CoolorColor) {
	fcc := ncc.Clone()
	fcc.SetStatic(true)
	fcc.SetPlain(true)
	cci.CoolorColor = fcc
	cci.Color = fcc.Color
	ccci := NewCoolorColorClusterInfo(cci.Clone())
	ccci.FindClusters()
	ccci.Sort()
	cci.titleBar.SetText(ccci.String())
	cci.UpdateView()
}

func (cci *CoolorColorInfo) UpdateView() {
	MainC.app.QueueUpdateDraw(func() {
		cci.Frame.SetFramed(cci.CoolorColor)
		cci.UpdateColorInfoViews()
		// cci.AddColorInfoRows()
	})
}

func (cci *CoolorColorInfo) OldUpdateView() {
	MainC.app.QueueUpdateDraw(func() {
		// titleBar := cci.GetTitle()

		// cci.gridView.AddItem(titleBar, 0, 0, 1, 4, 0, 0, false)
		// cci.gridView.SetRows(1)
		cci.Clear()
		cci.Flex.SetBorder(true)
		cci.Flex.SetBorderVisible(false)
		cci.SetDirection(tview.FlexColumn)
		cci.AddItem(cci.Frame, 0, 1, false)
		if cci.gridView == nil {
			cci.gridView = tview.NewGrid()
		} else {
			cci.gridView.Clear()
		}
		cci.gridView.SetBackgroundColor(theme.GetTheme().ContentBackground)
		cci.UpdateColorInfoViews()
		// cci.AddColorInfoRows()
		cci.gridView.SetOffset(0, 0)
	})
}

func (cci *CoolorColorInfo) AddTitle() *tview.TextView {
	titleBar := tview.NewTextView()
	titleBar.SetBorderPadding(0, 0, 1, 1)
	titleBar.SetTextAlign(tview.AlignCenter)
	titleBar.ScrollTo(0, 0)
	titleBar.SetDynamicColors(true)
	titleBar.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	cci.titleBar = titleBar
	cci.gridView.AddItem(cci.titleBar, 3, 0, 1, 4, 0, 0, false)
	return titleBar
}

func cf[T constraints.Ordered](i T, f string) string {
	return fmt.Sprintf(f, i)
}

var cffl = cf[float64]
var cfint = cf[int]
var cfui8 = cf[uint8]

func (cci *CoolorColorInfo) ColorValues() {
	cci.details.AddCssValue(
		"hex",
		"%s",
	)
	cci.details.AddCssValue(
		"rgb",
		"%s, %s, %s",
	)
	cci.details.AddCssValue(
		"srgb",
		"%s, %s, %s",
	)
	cci.details.AddCssValue(
		"hsl",
		"%s, %s, %s",
	)
	cci.details.AddCssValue(
		"hsv",
		"%s, %s, %s",
	)
	cci.details.AddCssValue(
		"LuvLCh",
		"%s, %s, %s",
	)
	cci.details.AddCssValue(
		"XYZ",
		"%s, %s, %s",
	)
	cci.details.AddCssValue(
		"xyY",
		"%s, %s, %s",
	)
	cci.details.AddCssValue(
		"L*a*b",
		"%s, %s, %s",
	)
}

func (cci *CoolorColorInfo) UpdateColorInfoViews() {
	tcol := MakeColorFromTcell(*cci.Color)
	cci.details.UpdateCssValue("hex", cci.Html())
	flts := "[magenta]%0.2f[::]"
	decs := "[magenta]%d[::]"
	r, g, b := tcol.RGB255()
	lr, lg, lb := tcol.LinearRgb()
	cci.details.UpdateCssValue(
		"rgb",
		fmt.Sprintf("%s, %s, %s", cfui8(r, decs), cfui8(g, decs), cfui8(b, decs)),
	)
	cci.details.UpdateCssValue(
		"srgb",
		fmt.Sprintf("%s, %s, %s", cffl(lr, flts), cffl(lg, flts), cffl(lb, flts)),
	)
	// cci.details.UpdateCssValue("rgba", fmt.Sprintf("[yellow:-:b] rgba(%d, %d, %d, %d) [-:-:-]", r, g, b, a))
	h, s, l := tcol.Hsl()
	cci.details.UpdateCssValue(
		"hsl",
		fmt.Sprintf("%s, %s, %s", cffl(h, flts), cffl(s, flts), cffl(l, flts)),
	)
	h, s, v := tcol.Hsv()
	cci.details.UpdateCssValue(
		"hsv",
		fmt.Sprintf("%s, %s, %s", cffl(h, flts), cffl(s, flts), cffl(v, flts)),
	)
	l, c, h := tcol.LuvLCh()
	cci.details.UpdateCssValue(
		"LuvLCh",
		fmt.Sprintf(
			"%s, %s, %s",
			cffl(l, flts), cffl(c, flts), cffl(h, flts),
		),
	)
	x, y, z := tcol.Xyz()
	cci.details.UpdateCssValue(
		"XYZ",
		fmt.Sprintf("%s, %s, %s", cffl(x, flts), cffl(y, flts), cffl(z, flts)),
	)
	ciex, ciey, ciey2 := tcol.Xyy()
	cci.details.UpdateCssValue(
		"xyY",
		fmt.Sprintf(
			"%s, %s, %s",
			cffl(ciex, flts),
			cffl(ciey, flts),
			cffl(ciey2, flts),
		),
	)
	ciel, ciea, cieb := tcol.Lab()
	cci.details.UpdateCssValue(
		"L*a*b",
		fmt.Sprintf(
			"%s, %s, %s",
			cffl(ciel, flts),
			cffl(ciea, flts),
			cffl(cieb, flts),
		),
	)
}

var (
	primaryCssValues []string = []string{
		"hex",
		"rgb",
		"srgb",
		"hsl",
		"hsv",
		"LuvLCh",
		"XYZ",
		"xyY",
		"L*a*b",
	}
	selectableCssValues []string
)
var (
	keyLabel = []rune{'⁰', '¹', '²', '³', '⁴', '⁵', '⁶', '⁷', '⁸', '⁹'}
)

func (cci *CoolorColorInfo) AddColorInfoRows() {
	count := 0
	baseRow := 4
	// valueSpan := 3
	// rowSizes := make([]int, 0)
	// rowSizes = append(rowSizes, 1)
	// rowSizes = append(rowSizes, 1)
	selectableCssValues = make([]string, 0)
	for _, cv := range primaryCssValues {
		ci, ok := cci.details.cssValues[cv]
		if !ok {
			panic(fmt.Errorf("no color info type named %s", cv))
		}
		row := float64(baseRow) + math.Floor(float64(count)/1)
		// col := count % 1
		// value := fmt.Sprintf("[green:-:b]%s[yellow:-:b]([-:-:-] %s [yellow:-:b])[-:-:-]", cv, ci)
		// labelCol := 0
		valueCol := 0
		// if cci.details.selected == cv {
		// label := fmt.Sprintf("[green:-:b] %s[-:-:-]", cv)
		// labelView := NewInfoTextView(label, true)
		// valueView := NewInfoTextView(value, false)
		// cci.gridView.
		// cci.gridView.AddItem(labelView, int(row), labelCol, 1, 1, 1, 1, false)
		cci.gridView.AddItem(ci.TextView, int(row), valueCol, 1, 4, 1, 1, false)
		// continue
		// } else {
		// 	selectableCssValues = append(selectableCssValues, cv)
		// 	if col != 0 {
		// 		labelCol = col*1 + labelCol
		// 		// valueCol = col*1 + valueCol
		// 	}
		// labelView := NewInfoTextView(label, true)
		// cci.gridView.AddItem(labelView, baseRow+int(row), labelCol, 1, 1, 1, 1, false)
		count++
		// }
	}
	// cci.gridView.SetRows(rowSizes...)
}
