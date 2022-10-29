package coolor

import (
	"fmt"
	"strings"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/digitallyserviced/coolors/coolor/shortcuts"

	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/theme"
)

type TabViewContent interface {
	tview.Primitive
	InputBubble
	show(*TabView)
	hide(*TabView)
}

type TabView struct {
	style   tcell.Style
	Content TabViewContent
	*TabLabel
	*tview.Frame
	*shortcuts.ScriptShortcut
	*EventObserver
	name  string
	color tcell.Color
}

type TabbedView struct {
	*tview.Flex
	tabsView  *tview.Flex
	Container *tview.Flex
	*shortcuts.ScriptShortcuts
	*EventNotifier
	name        string
	tabs        []*TabView
	selectedIdx int
}

type TabLabel struct {
	*tview.TextView
	*shortcuts.ScriptShortcut
	name, label           string
	labelFmt, selectedFmt string
	selected              bool
}

func (tl *TabLabel) SetSelected(sel bool) {
	tl.selected = sel
	tl.UpdateView()
}

func (tl *TabLabel) UpdateView() {
	lFmt := tl.labelFmt
	if tl.selected {
		lFmt = tl.selectedFmt
	}
	shortR := ""
	if tl.ScriptShortcut != nil {
		shortR = string(tl.ScriptShortcut.Script())
	}
	tabTitle := fmt.Sprintf(lFmt, tl.name)
  stringW := tview.TaggedStringWidth(tabTitle)
	// _, _, _, _, _, _, stringW := decomposeString(tabTitle, true, true)
	tabMarker := fmt.Sprintf(lFmt, strings.Repeat("🮂", stringW+3))
	bw := tl.BatchWriter()
	defer bw.Close()
	bw.Clear()
	//
	bw.Write([]byte(fmt.Sprintf(" %s[yellow]%s[-]", tabTitle, shortR)))
	bw.Write([]byte("\n"))
	bw.Write([]byte(tabMarker))
}

func NewTabLabel(name string) *TabLabel {
	tv := tview.NewTextView()
	tl := &TabLabel{
		TextView:    tv,
		name:        name,
		labelFmt:    "[red:-:-]%s[-:-:-]",
		selectedFmt: "[teal:-:b]%s[-:-:-]",
		selected:    false,
	}

	tv.SetDynamicColors(true)
	tv.SetToggleHighlights(false).SetRegions(true)
	tv.
		SetScrollable(false).
		SetTextAlign(tview.AlignCenter).
		SetWordWrap(false).
		SetWrap(false)
	tv.ScrollToBeginning()
	tv.SetMaxLines(2)
	tv.SetBorder(false).SetBorderPadding(1,0,1,1)
	tl.UpdateView()
	return tl
}

func NewTabView(name string, shortcut shortcuts.ScriptShortcut, p TabViewContent) *TabView {
	if p == nil {
		return nil
	}
	tv := &TabView{
		style:    tcell.StyleDefault.Normal(),
		Content:  p,
		TabLabel: NewTabLabel(name),
		ScriptShortcut: &shortcut,
		EventObserver:  NewEventObserver(name),
		name:           name,
		color:          0,
	}
	tv.TabLabel.ScriptShortcut = &shortcut
	tv.Frame = tview.NewFrame(tv.Content)
  tv.Frame.SetBorders(0, 0, 0, 0, 0, 0)
  // tv.Frame.AddText("SHIT", false, AlignCenter, tcell.ColorRebeccaPurple)
  // tv.Frame.AddText("SHIT", false, AlignCenter, tcell.ColorRebeccaPurple)
  tv.Frame.SetBorderPadding(0, 0, 1, 1)
  tv.Frame.SetBackgroundColor(theme.GetTheme().ContentBackground)
  tv.Frame.AddText(" ", false, tview.AlignCenter, theme.GetTheme().InfoLabel)
	return tv
}

func NewTabbedView() *TabbedView {
	tv := &TabbedView{
		Flex:            tview.NewFlex(),
		Container:       tview.NewFlex(),
		tabsView:        tview.NewFlex(),
		name:            "tabview",
		tabs:            make([]*TabView, 0),
		selectedIdx:     0,
		ScriptShortcuts: shortcuts.NewSuperScriptShortcuts(),
		EventNotifier:   NewEventNotifier("tabs"),
	}

	tv.SetDirection(tview.FlexRow)

	tv.Flex.AddItem(tv.tabsView, 3, 0, true)
	tv.Flex.SetBorder(false).SetBorderPadding(0, 0, 1, 1)
  // tv.Flex.SetBackgroundColor(theme.GetTheme().GrayerBackground)

	tv.SetDrawFunc(
		func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			tv.UpdateView()
			return x, y, width, height
		},
	)
	return tv
}

func (o *TabView) HandleEvent(e ObservableEvent) bool {
	var tabs *TabbedView = e.Src.(*TabbedView)
	var tab *TabView = e.Ref.(*TabView)
	if o == tab {
    tab.Content.show(tab)
	} else {
    tab.Content.hide(tab)
  }
	tabs.UpdateView()
  tabs.UpdateTabSelector()
	return true
}

func (tv *TabView) GetRef() interface{} {
	return tv
}

func (tv *TabbedView) GetRef() interface{} {
	return tv
}

func (tv *TabbedView) SetCurrentTab(prevTab, newTab int) {
	if len(tv.tabs) == 0 {
		return
	}
	if newTab < 0 {
		newTab = 0
	}
	tv.selectedIdx = newTab % len(tv.tabs)
	tb := tv.tabs[tv.selectedIdx]
	if tb == nil {
		return
	}
	oe := tv.NewObservableEvent(SelectedEvent, "tab_selected", tb, tv)
	tv.Notify(*oe)
	tv.UpdateView()
}

func (tv *TabbedView) GetCurrentTab() (*TabView, int) {
	// dump.P(tv.selectedIdx, len(tv.tabs))
	if tv.selectedIdx >= 0 && tv.selectedIdx < len(tv.tabs) {
		return tv.tabs[tv.selectedIdx], tv.selectedIdx
	}
	return nil, 0
}

func (tv *TabbedView) UpdateTabSelector() {
	_, i := tv.GetCurrentTab()
	tv.tabsView.Clear()
	tv.Each(func(tab *TabView, idx int) {
		tab.SetSelected(i == idx)
		tv.tabsView.AddItem(tab.TabLabel, 0, 1, false)
	})
}

func (tv *TabbedView) Each(f func(tab *TabView, idx int)) {
	for i, v := range tv.tabs {
		f(v, i)
	}
}

func (tv *TabbedView) Draw(s tcell.Screen) {
  tv.Box.DrawForSubclass(s, tv)
  tview.Borders = InvisBorders
  tv.Flex.Draw(s)
  tview.Borders = OrigBorders
}
func (tv *TabbedView) UpdateView() {
	curr, _ := tv.GetCurrentTab()
	if curr == nil {
		return
	}
  // dump.P(tv.GetRect())
	if tv.Flex.GetItemCount() == 2 {
		it := tv.Flex.GetItem(1)
		tv.Flex.RemoveItem(it)
	}
	tv.Flex.AddItem(curr, 0, 10, true)
}

func (tv *TabbedView) AddTab(tab *TabView) *TabbedView {
	if tab == nil {
		return nil
	}
	tv.tabs = append(tv.tabs, tab)
	MainC.app.QueueUpdateDraw(func() {
		tv.SetCurrentTab(0, 0)
		tv.Register(SelectedEvent, tab)
		tv.UpdateView()
		tv.UpdateTabSelector()
	})
	return tv
}

func (tv *TabbedView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return tv.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
    // log.Println(event)
		current, _ := tv.GetCurrentTab()
		if current != nil {
			// bubble := current.InputHandler()
		}
		key := event.Key()
		ch := event.Rune()
		prevTab := tv.selectedIdx
		if key == tcell.KeyRune {
			switch key {
			default:
				tv.Each(func(tab *TabView, idx int) {
					if ch == tab.ScriptShortcut.Text() {
						tv.SetCurrentTab(prevTab, idx)
						return
					}
				})
			}
		}
		tab, _ := tv.GetCurrentTab()
		ih := tab.Content.InputHandler()
		ih(event, setFocus)
		// MainC.app.QueueUpdateDraw(func() {
		// 	tv.SetCurrentTab(prevTab, tv.selectedIdx)
		// 	tv.UpdateTabSelector()
		// 	// ccs.UpdateView()
		// })
	})
}
