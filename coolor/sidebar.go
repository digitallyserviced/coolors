package coolor

import (
	"fmt"
	// "math"

	"github.com/digitallyserviced/tview"

	// "github.com/gookit/goutil/dump"

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
