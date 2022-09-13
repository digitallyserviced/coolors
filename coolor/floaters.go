package coolor

import (
	// "container/list"
	"fmt"
	"math"
	"strconv"

	// "strings"

	"github.com/digitallyserviced/coolors/theme"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
	"github.com/gookit/goutil/dump"
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

type RootFloatContainer struct {
	Item tview.Primitive
	*tview.Flex
	Rows          *tview.Flex
	Container     *tview.Flex
	cancel        func()
	finish        func()
	escapeCapture tcell.Key
	captureInput  bool
}

type Floater interface {
	GetRoot() *RootFloatContainer
}

type (
	ListSelectedHandler func(idx int, i interface{}, lis []ListItem)
	ListChangedHandler  func(idx int, selected bool, i interface{}, lis []ListItem)
)

type CoolorColorFloater struct {
	*tview.Flex
	Color *CoolorColorInfo
	Items []*tview.Primitive
}

type FixedFloater struct {
	Header *tview.TextView
	Footer *tview.TextView
	*RootFloatContainer
}
type ListFloater struct {
	Header *tview.TextView
	Footer *tview.TextView
	*RootFloatContainer
	*Lister
	// listItems []ListItem
}

type ListStyle interface {
	GetMainTextStyle() tcell.Style
	GetSecondaryTextStyle() tcell.Style
	GetShortcutStyle() tcell.Style
	GetSelectedStyle() tcell.Style
}

// GetSelected implements CoolorSelectable
func (*CoolorColorFloater) GetSelected() (*CoolorColor, int) {
	return nil, -1
}

// NavSelection implements VimNav
func (*CoolorColorFloater) NavSelection(int) {
}

func NewCoolorColorInfoFloater(ncc *CoolorColor) *CoolorColorFloater {
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
	centerFlex.AddItem(spf.Color.Flex, 13, 0, true)
	centerFlex.AddItem(nil, 0, 1, false)

	spf.SetDirection(tview.FlexColumn)
	spf.AddItem(nil, 0, 2, false)
	spf.AddItem(centerFlex, 0, 5, true)
	spf.AddItem(nil, 0, 2, false)
	return spf
}

func (f *RootFloatContainer) SetFinish(fin func()) *RootFloatContainer {
	f.finish = fin
	return f
}

func (f *RootFloatContainer) SetCancel(c func()) *RootFloatContainer {
	f.cancel = c
	return f
}

func (f *RootFloatContainer) GetRoot() *RootFloatContainer {
	return f
}

func (f *ListFloater) GetRoot() *RootFloatContainer {
	return f.RootFloatContainer
}

func NewListStyle() ListStyles {
	lis := ListStyles{
		main:  "",
		sec:   "",
		short: "",
		sel:   "",
	}

	return lis
}

func (f ListStyles) GetSelectedStyle() tcell.Style {
	if f.sel != "" {
		return *theme.GetTheme().Get(f.sel)
	}
	if theme.GetTheme().Get("list_sel") != nil {
		return *theme.GetTheme().Get("list_sel")
	}
	return tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor)
}

func (f ListStyles) GetShortcutStyle() tcell.Style {
	if f.short != "" {
		return *theme.GetTheme().Get(f.short)
	}
	if theme.GetTheme().Get("list_short") != nil {
		return *theme.GetTheme().Get("list_short")
	}
	return tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor)
}

func (f ListStyles) GetSecondaryTextStyle() tcell.Style {
	if f.sec != "" {
		return *theme.GetTheme().Get(f.sec)
	}
	if theme.GetTheme().Get("list_second") != nil {
		return *theme.GetTheme().Get("list_second")
	}
	return tcell.StyleDefault.Foreground(tcell.ColorBlue)
}

func (f ListStyles) GetMainTextStyle() tcell.Style {
	if f.main != "" {
		return *theme.GetTheme().Get(f.main)
	}
	if theme.GetTheme().Get("list_main") != nil {
		return *theme.GetTheme().Get("list_main")
	}
	return tcell.StyleDefault.Foreground(tcell.ColorGreen)
}

func (f *ListFloater) Selected() {
}

func (f *ListFloater) UpdateView() {
	f.Container.SetBorder(true)
	f.Container.SetDirection(tview.FlexRow)
	f.Container.Clear()
	f.Container.AddItem(f.Header, 1, 0, false)
	// f.Container.AddItem(f.List, 0, 6, true)

	// f.Lister.UpdateView()

	f.GetRoot().UpdateView()
	f.Container.AddItem(f.Footer, 1, 0, false)
}

func NewSelectionFloater(
	name string,
	il func() []*ListItem,
	sel func(lis ListItem, hdr *tview.TextView, ftr *tview.TextView),
	chg func(lis ListItem, hdr *tview.TextView, ftr *tview.TextView),
) *ListFloater {
	ler := NewLister()
	ler.SetItemLister(il)

	ler.UpdateListItems()
	f := &ListFloater{
		Header:             tview.NewTextView(),
		Footer:             tview.NewTextView(),
		RootFloatContainer: NewFloater(ler),
		Lister:             ler,
	}
	ler.SetHandlers(func(idx int, i interface{}, lis []*ListItem) {
		sel(*lis[idx], f.Header, f.Footer)
	}, func(idx int, selected bool, i interface{}, lis []*ListItem) {
		chg(*lis[idx], f.Header, f.Footer)
	})

	ler.SetBorderPadding(1, 1, 1, 1)

	f.Header.SetDynamicColors(true)
	f.Header.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name))
	f.Header.SetBackgroundColor(theme.GetTheme().SidebarBackground).
		SetBorderColor(theme.GetTheme().SidebarBackground)

	f.Footer.SetDynamicColors(true)
	f.Footer.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name))
	f.Footer.SetBackgroundColor(theme.GetTheme().SidebarBackground).
		SetBorderColor(theme.GetTheme().SidebarBackground)

	f.UpdateView()

	return f
}

func (f *FixedFloater) UpdateView() {
  // f.Container.SetBackgroundColor(theme.GetTheme().TopbarBorder)
  f.Container.SetBorder(false)
  f.Flex.SetBorder(false)
  f.Rows.SetBorder(false)
	f.Container.SetDirection(tview.FlexRow)
	f.Container.Clear()
	// f.Container.AddItem(f.Header, 3, 0, false)
	// f.Container.AddItem(f.List, 0, 6, true)

	// f.Lister.UpdateView()

	f.GetRoot().UpdateView()
  // f.Header.SetBorder(true).SetBorderPadding(1, 1, 1, 1)
	// f.Container.AddItem(f.Footer, 2, 0, false)
}

func NewFixedFloater(name string, p tview.Primitive) *FixedFloater {
	f := &FixedFloater{
		Header:             tview.NewTextView(),
		Footer:             tview.NewTextView(),
		RootFloatContainer: NewFloater(p),
	}
	f.RootFloatContainer.
  SetBorder(false).
		SetBorderPadding(0,0,1,1).
		SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.Container.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.RootFloatContainer.Rows.Clear()
	f.RootFloatContainer.Rows.AddItem(f.Container, 0, 10, true)
	f.RootFloatContainer.Clear()
	f.RootFloatContainer.AddItem(nil, 0, 70, false)
	f.RootFloatContainer.AddItem(f.Rows, 0, 30, true)

	f.Header.SetDynamicColors(true)
	f.Header.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name)).SetBorderPadding(1, 1, 1, 1)
	f.Header.SetBackgroundColor(theme.GetTheme().SidebarLines).
		SetBorderColor(theme.GetTheme().SidebarBackground)
	bw := f.Header.BatchWriter()
	bw.Close()

	f.Footer.SetDynamicColors(true)
	f.Footer.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name))
	f.Footer.SetBackgroundColor(theme.GetTheme().SidebarBackground).
		SetBorderColor(theme.GetTheme().SidebarBackground)

	f.UpdateView()

	return f
}

// func NewFixedFloater(nw, nh int, prop int) *RootFloatContainer {
//   r := NewSizedFloater(0,0,0)
//   r.Clear()
//   r.Rows.Clear()
//   r.AddItem(item tview.Primitive, fixedSize int, proportion int, focus bool)
//   return r
// }

func NewSizedFloater(nw, nh int, prop int) *RootFloatContainer {
	f := &RootFloatContainer{
		Flex:          tview.NewFlex(),
		Rows:          tview.NewFlex(),
		Container:     tview.NewFlex(),
		Item:          nil,
		captureInput:  true,
		escapeCapture: tcell.KeyEscape,
		cancel: func() {
		},
		finish: func() {
		},
	}

  f.Flex.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.SetDirection(tview.FlexColumn)
	f.Rows.SetDirection(tview.FlexRow)
	f.Container.SetBorder(true)
	f.Container.SetDirection(tview.FlexRow)

	f.Center(nw, nh, prop)

	return f
}

func (f *RootFloatContainer) Center(nw, nh int, prop int) {
	if nw == 0 && nh == 0 && prop == 1 {
		nw = 14
		nh = 10
		prop = 16
	}

	nwp := math.Abs(float64(nw)) / 100.0
	nhp := math.Abs(float64(nh)) / 100.0
	w, h := AppModel.scr.Size()
	rat := float64(float64(float64(h)/float64(w))) * 1.2
	padw := (prop - nw) / 2
	padh := (prop - nh) / 2
	pw, ph := 0, 0
	fw, fh := 0, 0

	if prop == 0 {
		if nw < 0 {
			nw = int(nwp * float64(w))
		}
		if nh == 0 {
			nh = int(float64(nw) * rat)
		}
		if nh < 0 {
			nh = int(nhp * float64(h))
		}
		if nw == 0 {
			nw = int(float64(nh) * rat)
		}
		padw = (w - nw) / 2
		padh = (h - nh) / 2
		fw = nw
		fh = nh

		dump.P(nw, nh, nwp, nhp, prop, w, h, rat, padw, padh, pw, ph, fw, fh)

		f.AddItem(nil, padw, 0, false)
		f.AddItem(f.Rows, fw, 0, true)
		f.AddItem(nil, padw, 0, false)

		f.Rows.AddItem(nil, padh, 0, false)
		f.Rows.AddItem(f.Container, fh, 0, true)
		f.Rows.AddItem(nil, padh, 0, false)
	}

	if prop > 1 {
		fw = 0
		fh = 0
		pw = prop - nw
		ph = prop - nh

		f.AddItem(nil, 0, padw, false)
		f.AddItem(f.Rows, 0, pw, true)
		f.AddItem(nil, 0, padw, false)

		f.Rows.AddItem(nil, 0, padh, false)
		f.Rows.AddItem(f.Container, 0, ph, true)
		f.Rows.AddItem(nil, 0, padh, false)
	}
}

func NewFloater(i tview.Primitive) *RootFloatContainer {
	f := NewSizedFloater(40, 16, 0)
	f.Item = i
	f.UpdateView()
	return f
}

func (f *RootFloatContainer) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return f.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			handler := f.Item.InputHandler()
			handler(event, setFocus)
			if event.Key() == tcell.KeyEscape {
				if f.cancel != nil {
					f.cancel()
				}
				name, _ := MainC.pages.GetFrontPage()
				MainC.pages.HidePage(name)
			}
			// if f.captureInput {
			//   event = nil
			// }
		},
	)
}

func (f *RootFloatContainer) UpdateView() {
	f.Container.AddItem(f.Item, 0, 8, true)
}

//
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

func (ccd *CoolorColorDetails) AddCssValue(name, value string) {
	ccd.cssValues[name] = value
}

func NewCoolorColorDetails(cc *CoolorColor) *CoolorColorDetails {
	ccd := &CoolorColorDetails{
		CoolorColor: cc,
		cssValues:   make(map[string]string),
		selected:    primaryCssValues[0],
	}
	return ccd
}

type PaletteFloater struct {
	*tview.Flex
	Palette *CoolorPaletteContainer
}

func NewScratchPaletteFloater(cp *CoolorColorsPalette) *PaletteFloater {
	spf := &PaletteFloater{
		Flex:    tview.NewFlex(),
		Palette: NewCoolorPaletteContainer(cp),
	}

	spf.SetDirection(tview.FlexRow)
	spf.AddItem(nil, 0, 2, false)
	spf.AddItem(spf.Palette, 0, 4, true)
	spf.AddItem(nil, 0, 2, false)
	return spf
}

func NewCoolorPaletteContainer(
	cp *CoolorColorsPalette,
) *CoolorPaletteContainer {
	p := cp.GetPalette()
	p.Plainify(true)
	p.Sort()
	pt := NewPaletteTable(p)
	cpc := &CoolorPaletteContainer{
		Frame:   tview.NewFrame(pt),
		Palette: pt,
	}
	cpc.SetBorders(1, 1, 0, 0, 0, 0)
	cpc.SetTitle("")
	cpc.Frame.SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(theme.GetTheme().TopbarBorder)
	return cpc
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

// vim: ts=2 sw=2 et ft=go
