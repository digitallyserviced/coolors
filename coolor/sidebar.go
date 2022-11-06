package coolor

import (
	"fmt"
	// "math"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/gookit/goutil/dump"

	"github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/theme"
)

type SideBar struct {
	*FixedFloater
	posLeft       bool
	width, height int
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
		SetBorder(false).
		SetBorderPadding(0, 0, 1, 1).
		SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.Container.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.RootFloatContainer.Rows.Clear()
	f.RootFloatContainer.Rows.AddItem(f.Container, 0, 10, true)

	if pos == 1 {
		f.RootFloatContainer.Clear()
		f.RootFloatContainer.AddItem(nil, 0, 100-w, false)
		f.RootFloatContainer.AddItem(f.Rows, 0, w, true)
	} else {
		f.RootFloatContainer.Clear()
		f.RootFloatContainer.AddItem(f.Rows, 0, w, true)
		f.RootFloatContainer.AddItem(nil, 0, 100-w, false)
	}

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

func OpenFavoritesView() *SideBar {
	tv := NewTabbedView()
	sb := NewSideBar("QuickColors", tv, 1, 30)
	ccs := NewCoolorColorSwatch(
		func(cs *CoolorColorSwatch) *CoolorColorsPalette {
			cp := GetStore().FavoriteColors.GetPalette()
      cp.Sort()
			return cp
		},
	)
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
					return true
				}
			case e.Type&events.ColorSelectionEvent != 0:
				col, ok := e.Ref.(*CoolorColor)
				if ok {
					MainC.palette.AddCoolorColor(col)
					return true
				}
			case e.Type&events.PaletteColorSelectionEvent != 0 || e.Type&events.PaletteColorSelectedEvent != 0:
				// pal, ok := e.Ref.(*CoolorColorsPalette)

				// sqs := fmt.Sprintf(
				// 	" %s ",
				// 	strings.Join(
				// 		pal.GetPalette().MakeSquarePalette(false),
				// 		" ",
				// 	),
				// )
				// if ok {
				// ccpstv.Frame.Clear().
				// 	AddText(sqs, false, tview.AlignCenter, tcell.ColorRed)
				// }
			default:
				fmt.Println(e)

			}
			return true
		},
	)
	// ccps.Register(events.ColorSelectionEvent, f)
	ccs.Register(events.ColorSelectedEvent, f)
	// tv.AddTab(ccpstv)
	tv.AddTab(ccstv)
	tv.UpdateView()
	return sb
}

func (sb *SideBar) UpdateView() {
	// f.Container.SetBackgroundColor(theme.GetTheme().TopbarBorder)
	sb.Container.SetBorder(false)
	sb.Flex.SetBorder(false)
	sb.Rows.SetBorder(false)
	sb.Container.SetDirection(tview.FlexRow)
	sb.Container.Clear()
	// f.Container.AddItem(f.Header, 3, 0, false)
	// f.Container.AddItem(f.List, 0, 6, true)

	// f.Lister.UpdateView()

	sb.GetRoot().UpdateView()
	// f.Header.SetBorder(true).SetBorderPadding(1, 1, 1, 1)
	// f.Container.AddItem(f.Footer, 2, 0, false)
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
