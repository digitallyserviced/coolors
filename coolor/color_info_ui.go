package coolor

import (
	"fmt"
	"math"
	"strconv"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"

	"github.com/digitallyserviced/coolors/theme"
)





type CoolorColorInfo struct {
	Frame *tview.Frame
	*tview.Flex
	*CoolorColor
	infoView *tview.Flex
	gridView *tview.Grid
	details  *CoolorColorDetails
}

type CoolorColorDetails struct {
	*CoolorColor
	cssValues map[string]string
	selected  string
}
type CoolorColorFloater struct {
	*tview.Flex
	Color *CoolorColorInfo
	Items []*tview.Primitive
}

func NewCoolorColorFloater(ncc *CoolorColor) *CoolorColorFloater {
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

	cci.UpdateColor(cc)
	return cci
}

func NewCoolorColorDetails(cc *CoolorColor) *CoolorColorDetails {
	ccd := &CoolorColorDetails{
		CoolorColor: cc,
		cssValues:   make(map[string]string),
		selected:    primaryCssValues[0],
	}
	return ccd
}
func NewInfoTextView(text string, bg bool) *tview.TextView {
	infoTextView := tview.NewTextView()
	infoTextView.SetWrap(false)
	infoTextView.ScrollTo(0, 0)
	infoTextView.SetDynamicColors(true)
	infoTextView.SetRegions(true)
	infoTextView.SetBorder(true)
	if bg {
		infoTextView.SetBorderColor(theme.GetTheme().SidebarBackground)
		infoTextView.SetBorderColor(theme.GetTheme().SidebarLines)
	} else {
		infoTextView.SetBorder(true)
	}
	infoTextView.SetTextAlign(tview.AlignCenter)
	infoTextView.SetText(text)
	infoTextView.ScrollTo(0, 0)
	return infoTextView
}
// GetSelected implements CoolorSelectable
func (*CoolorColorFloater) GetSelected() (*CoolorColor, int) {
	return nil, -1
}

// NavSelection implements VimNav
func (*CoolorColorFloater) NavSelection(int) {
}

func (ccd *CoolorColorDetails) AddCssValue(name, value string) {
	ccd.cssValues[name] = value
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
	cci.UpdateView()
}

func (cci *CoolorColorInfo) UpdateView() {
	MainC.app.QueueUpdateDraw(func() {
		newFrame := tview.NewFrame(cci.CoolorColor)
		newFrame.SetBorder(true).SetBorderPadding(0, 0, 0, 0)
		newFrame.SetBorders(0, 0, 1, 1, 1, 1)
		newFrame.SetBorderColor(tview.Styles.PrimitiveBackgroundColor)
		cci.Frame = newFrame
		topbar := tview.NewTextView()
		topbar.SetBorderPadding(0, 0, 0, 0)
		topbar.SetTextAlign(tview.AlignCenter)
		topbar.ScrollTo(0, 0)
		topbar.SetDynamicColors(true)
		ccci := NewCoolorColorClusterInfo(cci.Clone())
		ccci.FindClusters()
		ccci.Sort()
		topbar.SetText(tview.TranslateANSI(color.Render(ccci.String())))

		cci.Clear()
		cci.Flex.SetBorder(true)
		cci.SetDirection(tview.FlexColumn)
		cci.AddItem(cci.Frame, 0, 1, false)
		if cci.gridView == nil {
			cci.gridView = tview.NewGrid()
		} else {
			cci.gridView.Clear()
		}
		cci.AddItem(cci.gridView, 0, 5, true)
		cci.gridView.SetBackgroundColor(theme.GetTheme().ContentBackground)
		topbar.SetBorderPadding(0, 0, 0, 0)
		topbar.SetBackgroundColor(theme.GetTheme().SidebarBackground)
		cci.gridView.SetRows(1)
		cci.gridView.AddItem(topbar, 0, 0, 1, 4, 0, 0, false)
		cci.ColorInfoTable()
		cci.ColorInfoRows()
		cci.gridView.SetOffset(0, 0)
	})
}

func (cci *CoolorColorInfo) ColorInfoTable() {
	tcol := MakeColorFromTcell(*cci.Color)
	cci.details.AddCssValue("hex", cci.Html())
	h, s, l := tcol.Hsl()
	cci.details.AddCssValue(
		"hsl",
		fmt.Sprintf("[yellow:-:b] hsl(%0.2f, %0.2f, %0.2f) [-:-:-]", h, s, l),
	)
	r, g, b := tcol.RGB255()
	lr, lg, lb := tcol.LinearRgb()
	cci.details.AddCssValue(
		"rgb",
		fmt.Sprintf("[yellow:-:b] rgb(%d, %d, %d) [-:-:-]", r, g, b),
	)
	cci.details.AddCssValue(
		"srgb",
		fmt.Sprintf("[yellow:-:b] rgb(%0.2f, %0.2f, %0.2f) [-:-:-]", lr, lg, lb),
	)
	// cci.details.AddCssValue("rgba", fmt.Sprintf("[yellow:-:b] rgba(%d, %d, %d, %d) [-:-:-]", r, g, b, a))
	h, s, v := tcol.Hsv()
	cci.details.AddCssValue(
		"hsv",
		fmt.Sprintf("[yellow:-:b] hsl(%0.2f, %0.2f, %0.2f) [-:-:-]", h, s, v),
	)
	l, c, h := tcol.LuvLCh()
	cci.details.AddCssValue(
		"LuvLCh",
		fmt.Sprintf(
			"[yellow:-:b] Light: %0.2f Chroma: %0.2f Hue: %0.2f) [-:-:-]",
			l,
			c,
			h,
		),
	)
	x, y, z := tcol.Xyz()
	cci.details.AddCssValue(
		"XYZ",
		fmt.Sprintf("[yellow:-:b] X: %0.2f Y: %0.2f Z: %0.2f) [-:-:-]", x, y, z),
	)
	ciex, ciey, ciey2 := tcol.Xyy()
	cci.details.AddCssValue(
		"xyY",
		fmt.Sprintf(
			"[yellow:-:b] x: %0.2f y: %0.2f Y: %0.2f) [-:-:-]",
			ciex,
			ciey,
			ciey2,
		),
	)
	ciel, ciea, cieb := tcol.Lab()
	cci.details.AddCssValue(
		"L*a*b",
		fmt.Sprintf(
			"[yellow:-:b] L: %0.2f a: %0.2f b: %0.2f) [-:-:-]",
			ciel,
			ciea,
			cieb,
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
func (cci *CoolorColorInfo) ColorInfoRows() {
	count := 0
	baseRow := 2
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
		row := math.Floor(float64(count) / 4)
		col := count % 4
		value := fmt.Sprintf("[yellow:-:b] %s [-:-:-]", ci)
		labelCol := 0
		valueCol := 1
		if cci.details.selected == cv {
			label := fmt.Sprintf("[green:-:b] %s[-:-:-]", cv)
			labelView := NewInfoTextView(label, true)
			valueView := NewInfoTextView(value, false)
			cci.gridView.AddItem(labelView, 1, labelCol, 1, 1, 1, 1, false)
			cci.gridView.AddItem(valueView, 1, valueCol, 1, 3, 1, 1, false)
			continue
		} else {
			selectableCssValues = append(selectableCssValues, cv)
			label := fmt.Sprintf("[green:-:b] %s[-:-:-][red:-:-]%c[-:-:-]", cv, keyLabel[count%(len(keyLabel)-1)])
			if col != 0 {
				labelCol = col*1 + labelCol
				// valueCol = col*1 + valueCol
			}
			labelView := NewInfoTextView(label, true)
			cci.gridView.AddItem(labelView, baseRow+int(row), labelCol, 1, 1, 1, 1, false)
			count++
		}
	}
	// cci.gridView.SetRows(rowSizes...)
}
