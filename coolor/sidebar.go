package coolor

import (
	"fmt"
	"math"
	"strings"

	// "math"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/gookit/goutil/dump"

	"github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/shortcuts"
	"github.com/digitallyserviced/coolors/theme"
)

type PanelMenu struct {
	// *tview.Flex
	*tview.Grid
	maxWidth, selIdx int
	Items            []PanelItem
	title            string
}

type PanelMenuItem struct {
	Menu     *PanelMenu
	Item     tview.Primitive
	selected bool
	OnSelect func()
}
type PanelTextMenuItem struct {
	*PanelMenuItem
	name, wrap string
	*tview.TextView
}

type PanelItem interface {
	GetItem() tview.Primitive
	UpdateItem() tview.Primitive
}

type Panel struct {
	*FixedFloater
	Sibling *SideBarFloater
	*shortcuts.Scope
}

type SideBar struct {
	*FixedFloater
	Sibling       *SideBarFloater
	posLeft       bool
	width, height int
	*shortcuts.Scope
}

type SideBarFloater struct {
	*tview.Flex
	Padding *tview.Flex
	Item    tview.Primitive
	Sibling *SideBar
}

func NewPanel(name string, tvp tview.Primitive, args ...int) *Panel {
	f := &FixedFloater{
		Header:             tview.NewTextView(),
		Footer:             tview.NewTextView(),
		RootFloatContainer: NewFloater(tvp),
	}
	p := &Panel{
		FixedFloater: f,
	}
	f.RootFloatContainer.
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.Container.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.RootFloatContainer.Rows.Clear()
	f.RootFloatContainer.Rows.AddItem(f.Container, 0, 10, true)

	p.PositionPanel(args...)

	f.Header.SetDynamicColors(true)
	f.Header.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name)) // .SetBorderPadding(1, 1, 1, 1)
	f.Header.SetBackgroundColor(theme.GetTheme().SidebarLines).
		SetBorderColor(theme.GetTheme().SidebarBackground)
	bw := f.Header.BatchWriter()
	bw.Close()

	f.Footer.SetDynamicColors(true)
	f.Footer.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name))
	f.Footer.SetBackgroundColor(theme.GetTheme().SidebarBackground).
		SetBorderColor(theme.GetTheme().SidebarBackground)

	p.UpdateView()

	return p
}

func NewPanelMenu(title string) *PanelMenu {
	fv := tview.NewGrid()

	fv.SetDontClear(false)
	// fv.SetDirection(tview.FlexRow)
	fv.SetBackgroundColor(theme.GetTheme().SidebarBackground)

	mm := &PanelMenu{
		// Flex:     fv,
    Grid: fv,
		maxWidth: 0,
		selIdx:   0,
		Items:    make([]PanelItem, 0),
		title:    title,
	}
	return mm
}
func NewSideBar(name string, p tview.Primitive, args ...int) *SideBar {
	w := 30
	h := 10
	pos := 1
	if len(args) > 0 {
		pos = args[0]
		if len(args) > 1 {
			w = args[1]

			if len(args) > 2 {
				h = args[2]
			}
		}
	}
	f := &FixedFloater{
		Header:             tview.NewTextView(),
		Footer:             tview.NewTextView(),
		RootFloatContainer: NewFloater(p),
	}
	sb := &SideBar{
		FixedFloater: f,
		posLeft:      pos == 0,
		width:        w,
		height:       h,
	}
	f.RootFloatContainer.
		SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.Container.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.RootFloatContainer.Rows.Clear()
	f.RootFloatContainer.Rows.AddItem(f.Container, 0, 10, true)

	sb.PositionSidebar()

	f.Header.SetDynamicColors(true)
	f.Header.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name)) // .SetBorderPadding(1, 1, 1, 1)
	f.Header.SetBackgroundColor(theme.GetTheme().SidebarLines).
		SetBorderColor(theme.GetTheme().SidebarBackground)
	bw := f.Header.BatchWriter()
	bw.Close()

	f.Footer.SetDynamicColors(true)
	f.Footer.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name))
	f.Footer.SetBackgroundColor(theme.GetTheme().SidebarBackground).
		SetBorderColor(theme.GetTheme().SidebarBackground)

	sb.UpdateView()

	return sb
}
// ◧ ◧◨ ◧  
// ⎧  ⎫
//◀⎨  ⎬
// ⎩  ⎭
// ↉⌌⌍⌎⌏
// ⌌ ⌍ ⌎ ⌏ ⌐ ⌌        ⌍  ⌞ ⌟ ⭨
// ◀            MENU 
// ⌌ ⌍ ⌎ ⌏ ⌐ ⌎        ⌏   ⭨
// ⌌ ⌍ ⌎ ⌏ ⌐ ⌜       ⌝  ⌞ ⌟ ⭨
// ◀            MENU
// ⌌ ⌍ ⌎ ⌏ ⌐ ⌞       ⌟   ⭨
// ⎤
// ⎧▕▏⎫
// ⎪🮇▎⎪
// ⎪🮈▍⎪
// ⎪▐▌⎪
// ⎨🮉▋⎬
// ⎪🮊▊⎪
// ⎪🮋▉⎪
// ⎪██⎪
// ⎪🭨🭪⎪
// ⎩🭪🭨⎭
// ▁▂▃▄▅▆▇█⎥
// ▔🮂🮃▀🮄🮅🮆█⎥
//                              ⎦
// ⎞
// ⎠
func OpenPalettesHistory() {
	mm := NewPanelMenu("Palettes")

  pals := GetStore().LoadPalettes()

  for _, v := range pals {
    // i := MakeBoxItem(v.Name, v.Palette.GetItem(0).GetColor())
    i := tview.NewTextView()
    i.SetBorderPadding(0, 0, 0, 0).SetBorder(false)
    i.SetTextAlign(tview.AlignCenter)
    i.Clear().SetDynamicColors(true)
    v.Palette.Sort()
    i.SetText(strings.Join(v.Palette.MakeMenuPalette(false),"\n"))
    mm.NewPanelMenuItem(i)
  }

	p := NewPanel("main", mm, -30, 40, -30)
	p.Scope = shortcuts.NewScope(
		"main",
		"Main Menu",
		shortcuts.GlobalScope,
	)
	p.Scope.NewShortcut(
		"remove",
		"unfavorite",
		tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone),
		func(i ...interface{}) bool {
			return true
		},
	)
	mm.ResetView()
	MainC.Push("main", p, true)
	MainC.UpdateScope(p)
	MainC.app.SetFocus(p.Item)

}
func OpenMainMenu(title string) {
	mm := NewPanelMenu(title)

	mmi := mm.NewPanelTextMenuItem("NEW PALETTE")
	mmibd := mm.NewPanelTextMenuItem("LOAD PALETTE")
	_, _ = mmi, mmibd
	p := NewPanel("main", mm, -30, 40, -30)
	p.Scope = shortcuts.NewScope(
		"main",
		"Main Menu",
		shortcuts.GlobalScope,
	)
	p.Scope.NewShortcut(
		"remove",
		"unfavorite",
		tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone),
		func(i ...interface{}) bool {
			return true
		},
	)
	mm.ResetView()
	MainC.Push("main", p, true)
	MainC.UpdateScope(p)
	MainC.app.SetFocus(p.Item)

}

func OpenFavoritesView() *SideBar {
	tv := NewTabbedView()
	sb := NewSideBar("Favorites", tv, 1, 30)
	ccs := NewCoolorColorSwatch(
		func(cs *CoolorColorSwatch) *CoolorColorsPalette {
			cp := GetStore().FavoriteColors.GetPalette()
			cp.Sort()
			return cp
		},
	)
	var cci *CoolorColorInfo
	ccstv := NewTabView(" Favorites", tv.TakeNext(), ccs)
	ccstv.SetBackgroundColor(theme.GetTheme().GrayerBackground)
	f := events.NewAnonymousHandlerFunc(
		func(e events.ObservableEvent) bool {
			switch {
			case e.Type&events.ColorSelectedEvent != 0:
				col, ok := e.Ref.(*CoolorColor)
				if ok && col != nil {
					ccstv.Frame.Clear().
						AddText(col.TVPreview(), false, tview.AlignCenter, tcell.ColorRed)
					if cci != nil {
						cci.UpdateColor(col)
					}
					return true
				}
			case e.Type&events.ColorSelectionEvent != 0:
				col, ok := e.Ref.(*CoolorColor)
				if ok {
					MainC.palette.AddCoolorColor(col)
					return true
				}
			case e.Type&events.PaletteColorSelectionEvent != 0 || e.Type&events.PaletteColorSelectedEvent != 0:
			default:
				// fmt.Println(e)

			}
			return true
		},
	)
	// ccps.Register(events.ColorSelectionEvent, f)
	ccs.Register(events.ColorSelectedEvent, f)
	// tv.AddTab(ccpstv)
	tv.AddTab(ccstv)
	tv.UpdateView()
	sb.Scope = shortcuts.NewScope(
		"favorites",
		"Favorites View",
		shortcuts.GlobalScope,
	)
	sb.Scope.NewShortcut(
		"remove",
		"unfavorite",
		tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone),
		func(i ...interface{}) bool {
			r, c := ccs.Table.GetSelection()
			col := ccs.GetColorIndex(
				r,
				ccs.TableContent.rows,
				c,
				ccs.TableContent.cols,
			)
			ccs.Notify(
				*ccs.NewObservableEvent(events.ColorUnfavoriteEvent, "unfavorited", col, ccs),
			)
			GetStore().MetaService.ToggleFavorite(col)
			ccs.UpdateView()
			return true
		},
	)
	sb.Scope.NewShortcut(
		"info",
		"show color info",
		tcell.NewEventKey(tcell.KeyRune, 'i', tcell.ModNone),
		func(i ...interface{}) bool {
			r, c := ccs.Table.GetSelection()
			col := ccs.GetColorIndex(
				r,
				ccs.TableContent.rows,
				c,
				ccs.TableContent.cols,
			)
			if cci == nil {
				cci = NewCoolorColorInfo(col)
				cci.details.Scope.Parent = sb.Scope
				sb.NewSideBarFloater("color_info", cci.Flex)
				MainC.UpdateScope(cci.details)
			} else {
				sb.ClearSideBarFloater()
				MainC.UpdateScope(sb)
				cci = nil
			}
			ccs.Notify(
				*ccs.NewObservableEvent(events.ColorEvent, "color_info", col, ccs),
			)
			return true
		},
	)
	MainC.Push("favorites", sb, true)
	MainC.UpdateScope(sb)
	MainC.app.SetFocus(sb.Item)
	return sb
}
func (ccs *SideBar) GetScope() *shortcuts.Scope {
	return ccs.Scope
}

func (sb *SideBar) UpdateView() {
	sb.Container.SetBorder(false)
	sb.Flex.SetBorder(false)
	sb.Rows.SetBorder(false)
	sb.Container.SetDirection(tview.FlexRow)
	sb.Container.Clear()
	sb.GetRoot().UpdateView()
	col := theme.GetTheme().Border
	cc := NewDefaultCoolorColor()
	cc.SetColor(&col)
	sb.Container.SetBorder(true).SetBorderColor(cc.GetFgColor())
	sb.Container.SetBorderSides(false, true, false, false)
}
func (sb *SideBar) ClearSideBarFloater() {
	// idx := IfElse[int](!sb.posLeft, 0, 1)
	// sb.RootFloatContainer.SetItem(idx, nil, 0, 100 - sb.width - 2)
	sb.PositionSidebar()
}

func (sb *SideBar) NewSideBarFloater(name string, p tview.Primitive) *SideBarFloater {
	f := &SideBarFloater{
		Flex:    tview.NewFlex(),
		Padding: tview.NewFlex(),
		Item:    p,
		Sibling: sb,
	}
	sb.Sibling = f

	f.Flex.SetDirection(tview.FlexColumn)
	f.Flex.AddItem(nil, 0, 3, false)
	f.Flex.AddItem(f.Item, 40, 0, false)
	f.Flex.AddItem(nil, 0, 3, false)

	idx := IfElse[int](!sb.posLeft, 0, 1)
	sb.RootFloatContainer.SetItem(idx, f, 0, 100-sb.width-2)

	return f
}

func (sb *SideBar) PositionSidebar() {
	w := sb.width
	sb.FixedFloater.RootFloatContainer.Clear()
	if !sb.posLeft {
		sb.FixedFloater.RootFloatContainer.AddItem(nil, 0, 100-w, false)
		sb.FixedFloater.RootFloatContainer.AddItem(sb.FixedFloater.Rows, 0, w, true)
	} else {
		sb.FixedFloater.RootFloatContainer.AddItem(sb.FixedFloater.Rows, 0, w, true)
		sb.FixedFloater.RootFloatContainer.AddItem(nil, 0, 100-w, false)
	}
}

func (p *Panel) PositionPanel(args ...int) {
	p.FixedFloater.RootFloatContainer.Clear()
	w := AppModel.w
	for _, v := range args {
		abs := (math.Abs(float64(v)) / 100.0) * float64(w)
		if v < 0 {
			p.FixedFloater.RootFloatContainer.AddItem(nil, int(abs), 0, false)
			// p.FixedFloater.RootFloatContainer.AddItem(MakeBoxItem(fmt.Sprintf("%d", v),"#6e8493"), int(abs), 0, false)
		} else {
			p.FixedFloater.RootFloatContainer.AddItem(p.FixedFloater.Rows, int(abs), 0, true)
		}
	}
}

// GetItem implements PanelItem
func (pmi *PanelTextMenuItem) UpdateItem() tview.Primitive {
	colwrap := "[blue]╍[-][green]╺╸[-][magenta]╺━[-][red]╺━╸[-]⎸[yellow]%s[-]⎹[red]╺━╸[-][magenta]━╸[-][green]╺╸[-][blue]╍[-]"
	dimwrap := "╍╺╸╺━╺━╸⎸%s⎹╺━╸━╸╺╸╍"
		n := theme.Jcenter(pmi.name, pmi.Menu.maxWidth)
		pmi.SetDynamicColors(true)
		pmi.SetText(fmt.Sprintf(colwrap, n))
		if !pmi.selected {
			pmi.SetText(fmt.Sprintf(dimwrap, n))
		}
  return pmi
}
func (pmi *PanelTextMenuItem) GetItem() tview.Primitive {
  return pmi
}

// GetItem implements PanelItem
func (pmi *PanelMenuItem) UpdateItem() tview.Primitive {
  return pmi.Item
}
func (pmi *PanelMenuItem) GetItem() tview.Primitive {
  // pmi.Menu.SetBorder(false)

  return pmi.Item
}

func (mm *PanelMenu) ResetView() {
	mm.Clear()
	tv := tview.NewTextView()
	tv.SetBorderPadding(0, 0, 0, 0)
	tv.SetText(mm.title)
	tv.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
  tv.SetMaxLines(1)

	// mm.AddItem(tv, 0, 1, false)
		mm.AddItem(tv, 0, 0, 1, 1, 1, 1, false)
	//
	//
	// ╍╺╸╺━╺━╸⎸LOAD PALETTE⎹╺━╸━╸╺╸╍

  rows := make([]int, len(mm.Items))

    rows = append(rows, 3)
	for n, i := range mm.Items {
    rows = append(rows, 3)
    i.UpdateItem()
    // mm.SetBo
		mm.AddItem(i.GetItem(), 1+n, 0, 1, 1, 3, 1, false)
		// mm.AddItem(MakeBoxItem("SHITE", RandomColor().Hex()), 1+n*2+1, 0, 1, 1, -4, 1, false)
		// mm.Flex.AddItem(nil, 3, 0, false)
	}
  mm.SetMinSize(3, 8)
  mm.SetRows(rows...)
}

func (mm *PanelMenu) NewPanelTextMenuItem(name string) *PanelTextMenuItem {
	mmi := &PanelMenuItem{
		Menu:     mm,
		selected: false,
		OnSelect: func() {
		},
	}

	ptmi := &PanelTextMenuItem{
		PanelMenuItem: mmi,
		TextView:      &tview.TextView{},
		name:          name,
		wrap:          "[blue]╍[-][green]╺╸[-][magenta]╺━[-][red]╺━╸[-]⎸[yellow]%s[-]⎹[red]╺━╸[-][magenta]━╸[-][green]╺╸[-][blue]╍[-]",
		// wrap:     "[#ecb31e]╍[-][#28e589]╺╸[-][#4E6798]╺━[-][#9F175B]╺━╸[-]⏽ [yellow]%s[-]⏽ [#9F175B]╺━╸[-][#4E6798]━╸[-][#28e589]╺╸[-][#ecb31e]╍[-]",
	}

	if mm.maxWidth == 0 {
		mmi.selected = true
	}
	if mm.maxWidth < len(name) {
		mm.maxWidth = len(name)
		fmt.Println(mm.maxWidth)
	}

	tv := tview.NewTextView()
	tv.SetBorderPadding(1, 1, 0, 0)
	tv.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)

	ptmi.TextView = tv

	mm.Items = append(mm.Items, ptmi)

	return ptmi
}

func (mm *PanelMenu) NewPanelMenuItem(p tview.Primitive) *PanelMenuItem {
	mmi := &PanelMenuItem{
		Menu:     mm,
		Item:     p,
		selected: false,
		OnSelect: func() {
		},
	}

	mm.Items = append(mm.Items, mmi)

	return mmi
}
func (p *Panel) GetScope() *shortcuts.Scope {
	return p.Scope
}

func (p *Panel) UpdateView() {
	p.Container.SetBorder(false)
	p.Flex.SetBorder(false)
	p.Rows.SetBorder(false)
	p.Container.SetDirection(tview.FlexRow)
	p.Container.Clear()
	// p.GetRoot().UpdateView()
	p.Container.AddItem(p.Item, 0, 8, true)
	col := theme.GetTheme().Border
	cc := NewDefaultCoolorColor()
	cc.SetColor(&col)
	p.Container.SetBorder(true).SetBorderColor(cc.GetFgColor())
	p.Container.SetBorderSides(false, true, false, true)
}

// ccps := NewCoolorColorsPaletteSwatch(
// 	func(cs *CoolorColorsPaletteSwatch) []*CoolorColorsPalette {
// 		cps := make([]*CoolorColorsPalette, 0)
// 		ccps := GetStore().PaletteHistory(false)
// 		for _, v := range ccps {
// 			if v.Current == nil || v.Current.Colors == nil ||
// 				len(v.Current.Colors) == 0 {
// 				continue
// 			}
// 			cps = append(cps, v.Current.GetPalette())
// 		}
// 		return cps
// 	},
// )
// ccpstv := NewTabView(" Palettes", tv.TakeNext(), ccps)
// ccpstv.SetBackgroundColor(ct.GetTheme().GrayerBackground)
// ▶️ ⃤
// ■ □ ▢ ▣ ▤ ▥ ▦ ▧ ▨ ▩ ▪ ▫ ▬ ▭ ▮ ▯ ▰ ▱ ▲ △ ▴ ▵ ▶ ▷ ▸ ▹ ► ▻ ▼ ▽ ▾ ▿ ◀ ◁ ◂ ◃ ◄ ◅ ◆ ◇ ◈ ◉ ◊ ○ ◌ ◍ ◎ ● ◐ ◑ ◒ ◓ ◔ ◕ ◖ ◗ ◘ ◙ ◚ ◛ ◜ ◝ ◞ ◟ ◠ ◡ ◢ ◣ ◤ ◥ ◦ ◧ ◨ ◩ ◪ ◫ ◬ ◭ ◮ ◯ ◰ ◱ ◲ ◳ ◴ ◵ ◶ ◷ ◸ ◹ ◺ ◻ ◼ ◽ ◾ ◿       🞀 🞁 🞂 🞃 🞁 🞄 🞅 🞆 🞇 🞈 🞉 🞊 🞋 🞗 🞘 🞙 🞚 🞛 🞜 🞝 🞞 🞟 🞠 🞌 🞍 🞎 🞏 🞐 🞑 🞒 🞓 🞔 🞕 🞖 ‐ ‑ ‒ – — ― ❙ ❚ ❛ 
//   ■ □ ▢ ▣ ▤ ▥ ▦ ▧ ▨ ▩ ▪ ▫ ▬ ▭ ▮ ▯ ▰ ▱ ▲ △ ▴ ▵ ▶ ▷ ▸ ▹ ► ▻ ▼ ▽ ▾ ▿ ◀ ◁ ◂ ◃ ◄ ◅ ◆ ◇ ◈ ◉ ◊ ○ ◌ ◍ ◎ ● ◐ ◑ ◒ ◓ ◔ ◕ ◖ ◗ ◘ ◙ ◚ ◛ ◜ ◝ ◞ ◟ ◠ ◡ ◢ ◣ ◤ ◥ ◦ ◧ ◨ ◩ ◪ ◫ ◬ ◭ ◮ ◯ ◰ ◱ ◲ ◳ ◴ ◵ ◶ ◷ ◸ ◹ ◺ ◻ ◼ ◽ ◾ ◿       🞀 🞁 🞂 🞃 🞁 🞄 🞅 🞆 🞇 🞈 🞉 🞊 🞋 🞗 🞘 🞙 🞚 🞛 🞜 🞝 🞞 🞟 🞠 🞌 🞍 🞎 🞏 🞐 🞑 🞒 🞓 🞔 🞕 🞖 ‐ ‑ ‒ – — ― ❙ ❚ ❛ 
// ⌠⌡⌃⌁⌇⌈⌉⌊⌌⌍⌎⌏⌐⌜⌝⌞⌟⌢⌣⌤⌥⌦⌧⌨〈〉⌫⌬⌭⌰⌴⌷⌸⌹⌺⌻⌼⌽⌾⎅⎇⎌⎍⎎⎏⎐⎑⎒⎔⎕⎖⎗⎘⎙⎚⎛⎜⎝⎞⎟⎠⎡⎢⎣⎤⎥⎦⎧⎨⎩⎪⎫⎬⎭⎮⎯⎰⎱⎲⎳⎴⎵⎶⎷⎸⎹⎺⎻⎼⎽⎾⎿⏀⏁⏂⏃⏄⏅⏆⏇⏈⏉⏊⏋⏌⏍⏎⏏⏐⏑⏒⏓⏔⏕⏖⏗⏘⏙⏚⏛⏜⏝⏞⏟⏠⏡⏢⏣⏤⏥⏦⏧⏨⏿⏾⏽⏼⏻🙼🙽🙾🙿
// ⌠ ⌡ ⌃ ⌁ ⌇ ⌈ ⌉ ⌊ ⌌ ⌍ ⌎ ⌏ ⌐ ⌜⌝ ⌞ ⌟ ⌢ ⌣ ⌤ ⌥ ⌦ ⌧ ⌨ 〈 〉 ⌫ ⌬ ⌭ ⌰ ⌴ 
// ⌷ ⌸ ⌹ ⌺ ⌻ ⌼ ⌽ ⌾ ⎅ ⎇ ⎌ ⎍ ⎎ ⎏ ⎐ ⎑ ⎒ ⎔ ⎕ ⎖ ⎗ ⎘ ⎙ ⎚ ⎛ ⎜ ⎝ ⎞ ⎟ ⎠ 
// ⎡ ⎢ ⎣ ⎤ ⎥ ⎦ ⎧ ⎨ ⎩ ⎪ ⎫ ⎬ ⎭ ⎮ ⎯ ⎰ ⎱ ⎲ ⎳ ⎴ ⎵ ⎶ ⎷ ⎸ ⎹ ⎺ ⎻ ⎼ ⎽ 
// ⎾ ⎿ ⏀ ⏁ ⏂ ⏃ ⏄ ⏅ ⏆ ⏇ ⏈ ⏉ ⏊ ⏋ ⏌ ⏍ ⏎ ⏏ ⏐ ⏑ ⏒ ⏓ ⏔ ⏕ ⏖ ⏗ ⏘ ⏙ ⏚ 
// ⏛ ⏜ ⏝ ⏞ ⏟ ⏠ ⏡ ⏢ ⏣ ⏤ ⏥ ⏦ ⏧ ⏨ ⏿ ⏾ ⏽ ⏼ ⏻ 🙼 🙽 🙾 🙿                
//             
//                               
//                       
//                          
//                      
//                         
//              
