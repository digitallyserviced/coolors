package coolor

import (
	// "container/list"
	"fmt"
	"math"

	// "strings"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/samber/lo"

	"github.com/digitallyserviced/coolors/theme"
	"github.com/digitallyserviced/coolors/coolor/lister"

	// "github.com/gdamore/tcell/v2"

	"github.com/gookit/goutil/dump"
)

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
	ListSelectedHandler func(idx int, i interface{}, lis []lister.ListItem)
	ListChangedHandler  func(idx int, selected bool, i interface{}, lis []lister.ListItem)
)

type FixedFloater struct {
	Header *tview.TextView
	Footer *tview.TextView
	*RootFloatContainer
}
type ListFloater struct {
	Header *tview.TextView
	Footer *tview.TextView
	*RootFloatContainer
	*lister.Lister
	// listItems []ListItem
}

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

func NewFloater(i tview.Primitive) *RootFloatContainer {
	f := NewSizedFloater(40, 16, 0)
	f.Item = i
	f.UpdateView()
	return f
}

func NewSelectionFloater(
	name string,
	il func() []*lister.ListItem,
	sel func(lis lister.ListItem, hdr *tview.TextView, ftr *tview.TextView),
	chg func(lis lister.ListItem, hdr *tview.TextView, ftr *tview.TextView),
) *ListFloater {
	ler := lister.NewLister()
	ler.SetItemLister(il)

	ler.UpdateListItems()
	f := &ListFloater{
		Header:             tview.NewTextView(),
		Footer:             tview.NewTextView(),
		RootFloatContainer: NewFloater(ler),
		Lister:             ler,
	}
	ler.SetHandlers(func(idx int, i interface{}, lis []*lister.ListItem) {
		sel(*lis[idx], f.Header, f.Footer)
	}, func(idx int, selected bool, i interface{}, lis []*lister.ListItem) {
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

func NewFixedFloater(name string, p tview.Primitive) *FixedFloater {
	f := &FixedFloater{
		Header:             tview.NewTextView(),
		Footer:             tview.NewTextView(),
		RootFloatContainer: NewFloater(p),
	}
	f.RootFloatContainer.
		SetBorder(false).
		SetBorderPadding(0, 0, 1, 1).
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
func (f *ListFloater) GetRoot() *RootFloatContainer {
	return f.RootFloatContainer
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

// vim: ts=2 sw=2 et ft=go
