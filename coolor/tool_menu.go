package coolor

import (
	"fmt"
	"math"

	"github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/samber/lo"
)

type CoolorToolMenu struct {
	*tview.Flex
	list          *Lister
	mc            *MainContainer
	app           *tview.Application
	selectedColor *CoolorColor
	*eventObserver
	*eventNotifier
	menuItems        []*CoolorButtonMenuItem
	visibleItems     []*CoolorButtonMenuItem
	selected         int
	activeActionFlag CoolorColorActionFlag
}

type CoolorPaletteAction interface {
	IsActivated() bool
	Activate(cc *CoolorColor) bool
	Before(cp *CoolorBlendPalette) bool
	Every(cp *CoolorBlendPalette, menu *CoolorToolMenu) bool
	Finalize(cc *CoolorColor, cp *CoolorPaletteMainView, menu *CoolorToolMenu)
}

type CoolorButtonMenuItem struct {
	*tview.Button
	menu     *CoolorToolMenu
	action   *CoolorColorActor
	icon     string
	name     string
	selected bool
}

func NewButtonMenuItem(
	menu *CoolorToolMenu,
	action *CoolorColorActor,
) *CoolorButtonMenuItem {
	action.menu = menu
	mmi := &CoolorButtonMenuItem{
		Button:   tview.NewButton(action.name),
		name:     action.name,
		action:   action,
		selected: false,
		menu:     menu,
		icon:     action.icon,
	}
	mmi.Button.SetLabel(mmi.icon).
		SetBorder(false).
		SetTitleAlign(tview.AlignCenter)
	mmi.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	mmi.SetLabelColor(tview.Styles.InverseTextColor)
	mmi.SetSelectedFunc(func() {
		mmi.UpdateState()
	})
	// mmi.SetBackgroundColor(tcell.Color(0))
	return mmi
}

func NewCoolorColorMainMenu(app *tview.Application) *CoolorToolMenu {
	cmm := &CoolorToolMenu{
		Flex:             tview.NewFlex(),
		list:             &Lister{},
		mc:               MainC,
		app:              app,
		selectedColor:    &CoolorColor{},
		eventObserver:    NewEventObserver("menu"),
		menuItems:        []*CoolorButtonMenuItem{},
		visibleItems:     []*CoolorButtonMenuItem{},
		selected:         0,
		activeActionFlag: 0,
	}
	cmm.SetDirection(tview.FlexRow)
	cmm.list = NewLister()
  // cmm.list.highlightType=ListerHighlightBars
	cmm.list.SetItemLister(cmm.GetListItems)
  cmm.list.SetHandlers(func(idx int, i interface{}, lis []*ListItem) {
    fmt.Println(idx,i)
  }, func(idx int, selected bool, i interface{}, lis []*ListItem) {
    fmt.Println("chg",idx,selected,i)

    })
  cmm.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {

    cmm.Selected().action.Before(cmm.selectedColor.pallette, cmm.selectedColor)
	// cmm.Before()
  // dump.P(idx, selected, fmt.Sprintf("%T %v", i))
  return x,y,width,height
  })
	return cmm
}

func (cc CoolorButtonMenuItem) MainText() string {
	if cc.selected {
    col := fmt.Sprintf("#%06x", cc.menu.selectedColor.GetFgColor().Hex())
		return fmt.Sprintf("[%s:-:b]%s[-:-:-]", col, cc.action.icon)
	}
  col := fmt.Sprintf("#%06x", cc.menu.selectedColor.GetFgColorShade().Hex())
  return fmt.Sprintf("[%s:-:-]%s[-:-:-]", col, cc.action.icon)
}

func (cc CoolorButtonMenuItem) SecondaryText() string {
	return ""
}

func (cc CoolorButtonMenuItem) Shortcut() ScriptShortcut {
	return NewScriptShortcut(rune(0), rune(0))
}

func (cc CoolorButtonMenuItem) Changed(
	idx int,
	selected bool,
	i interface{},
	lis []*ListItem,
) {
  cc.selected=selected
  // fmt.Println(idx,selected,i)
	if cc.menu.selectedColor == nil || cc.menu == nil ||
		cc.menu.selectedColor.Color == nil {
		return
	}
	if i == nil {
		return
	}
	cc.action.Every(cc.menu.selectedColor)
	cc.menu.UpdateActionStatus(&cc)
}

func (cc CoolorButtonMenuItem) Visibility() ListItemsVisibility {
	// dump.P(cc.menu.activeActionFlag&cc.action.flag, cc.menu.activeActionFlag, cc.action.flag, cc.action.actionSet)
	if cc.menu.activeActionFlag == 0 {
		return ListItemDefault
	}
	if cc.menu.activeActionFlag&cc.action.flag != 0 {
		return ListItemVisible
	}
	return ListItemHidden
}

func (cc CoolorButtonMenuItem) Cancelled(
	idx int,
	i interface{},
	lis []*ListItem,
) {
	MainC.app.QueueUpdateDraw(func() {
		cc.action.Cancel()
	})
}

func (cc CoolorButtonMenuItem) Selected(
	idx int,
	i interface{},
	lis []*ListItem,
) {
	cc.action.Before(cc.menu.selectedColor.pallette, cc.menu.selectedColor)
	MainC.app.QueueUpdate(func() {
		// cc.action.Before(cc.menu.selectedColor.pallette, cc.menu.selectedColor)
		if cc.action.Activate(cc.menu.selectedColor) {
			cc.menu.ResetViews()
		}
		cc.menu.UpdateVisibleActors(cc.action.Actions())
		cc.menu.UpdateActionStatus(&cc)
	})
}

func (cc *CoolorToolMenu) GetMainTextStyle() tcell.Style {
	if cc.selectedColor == nil || cc.selectedColor.Color == nil {
		return tcell.Style{}
	}
	// tcol := cc.selectedColor.color
	fcol := cc.selectedColor.GetFgColorShade()
	return tcell.StyleDefault.Foreground(fcol)
}

func (cc *CoolorToolMenu) GetSecondaryTextStyle() tcell.Style {
	return tcell.Style{}
}

func (cc *CoolorToolMenu) GetShortcutStyle() tcell.Style {
	return tcell.Style{}
}

func (cc *CoolorToolMenu) GetSelectedStyle() tcell.Style {
	if cc.selectedColor == nil || cc.selectedColor.Color == nil {
		return tcell.Style{}
	}
	tcol := cc.selectedColor.Color
	fcol := cc.selectedColor.GetFgColor()
	return tcell.StyleDefault.Foreground(fcol).Background(*tcol)
}

func (f *CoolorToolMenu) GetListItems() []*ListItem {
	lits := make([]*ListItem, 0)
	for _, v := range f.visibleItems {
		var li ListItem = ListItem(v)
		lits = append(lits, &li)
	}
	return lits
}

// func (ctm *CoolorToolMenu) Init() {
// 	for _, v := range acts {
// 		ctm.AddItem(NewButtonMenuItem(ctm, v))
// 	}
// 	ctm.SetSelected(0)
// 	ctm.app.Draw()
// }
// Draw draws this primitive onto the screen.
func (f *CoolorToolMenu) Draw(screen tcell.Screen) {
	if f == nil || f.list == nil {
		// dump.P(f, f.list, f.list.Box)
		return
	}
	f.SetDontClear(true)
	f.list.SetDontClear(true)
	f.list.Box.DrawForSubclass(screen, f.list)
	// f.list.DrawForSubclass(screen, f.list)

	// How much space can we distribute?
	x, y, _, height := f.GetInnerRect()
	distSize := height
	distSize -= 2 * len(f.visibleItems)

	// drawItems := f.visibleItems[:0]
	// itemOrder := f.RotateItems(f.selected)

	// Calculate positions and draw items.
	// pos := y
	if f.list == nil {
		return
	}
	f.list.SetTitleAlign(tview.AlignCenter)
	selstyle := f.GetSelectedStyle()
	selfg, selbg, selattr := selstyle.Decompose()
	_, _, _ = selfg, selbg, selattr
	f.list.SetSelectedStyle(selstyle)
	mainstyle := f.GetMainTextStyle()
	mainfg, mainbg, mainattr := mainstyle.Decompose()
	_, _, _ = mainattr, mainbg, mainfg
	f.list.SetMainTextStyle(mainstyle)
	f.list.SetMainTextColor(mainfg)

	f.list.ShowSecondaryText(true)
	f.list.SetBorder(false)
	f.list.SetBorderPadding(0, 0, 0, 0)
	width := 1
	f.list.SetRect(x, y, width, height-5)
	// f.list.UpdateView()
	// f.list.SetWrapAround(true)
	f.list.Draw(screen)

	//  itemorder := lo.rangefrom[int](0, len(f.visibleitems))
	//  fmt.println(itemorder)
	// for _, item := range f.visibleitems {
	// 	size := 2
	//    // item, _ := f.findact(act)
	// 	// item := f.visibleitems[num]
	//
	// 	if item != nil {
	// 		item.setrect(x, pos, width, size)
	// 	}
	// 	pos += size
	//
	// 	if item != nil {
	// 			item.draw(screen)
	// 		// if item.hasfocus() {
	// 		// 	// defer item.draw(screen)
	// 		// } else {
	// 		// }
	// 	}
	// }
}

func (ctm *CoolorToolMenu) Init() {
	// ctm.list.List.SetSelectedFunc(func(i int, s1, s2 string, r rune) {
	//
	//
	// })
	ctm.SetBorder(false)
	ctm.SetBorderPadding(0, 0, 0, 0)
	for _, v := range actOrder {
		act, ok := acts[v]
		if !ok {
			continue
		}
		ctm.AddItem(NewButtonMenuItem(ctm, act))
	}
	// ctm.list.ListerHandler.SetChangedFunc(func(idx int, selected bool, i interface{}, lis []ListItem) {
	//   dump.P(idx, selected, i, lis)
	// })
	// ctm.list.ListerHandler.SetSelectedFunc(func(idx int, i interface{}, lis []ListItem) {
	//   dump.P(idx, i, lis)
	// })
	ctm.Flex.AddItem(ctm.list, 0, 1, true)
	ctm.list.ShowSecondaryText(true)
	ctm.list.SetBorder(false)
	ctm.list.SetBorderPadding(0, 0, 0, 0)
	// ctm.list.SetRect(x, y, width, height - 5)
	ctm.list.ShowSecondaryText(false)
	ctm.list.SetWrapAround(true)
	ctm.list.SetWrapAround(true).SetHighlightFullLine(false)
	// selCol := ctm.selectedColor.GetFgColor()
	// ctm.list.SetSelectedStyle(tcell.StyleDefault.Foreground(selCol))
	// mainCol := ctm.selectedColor.GetFgColor()
	// ctm.list.SetMainTextStyle(tcell.StyleDefault.Foreground(mainCol))

	ctm.SetSelected(0)
	// ctm.app.Draw()
}

func (mmi *CoolorButtonMenuItem) UpdateState() {
	if mmi.selected {
		mmi.SetLabel(fmt.Sprintf(`[::b]%s[::-]`, mmi.action.icon))
	} else {
		mmi.SetLabel(fmt.Sprintf(`%s`, mmi.action.icon))
	}
}

func (ctm *CoolorToolMenu) GetMenuItem(m string) (cbmi *CoolorButtonMenuItem) {
	ctm.forMenuItems(false, func(c *CoolorButtonMenuItem, idx int) (err error) {
		if c.name == m {
			cbmi = c
		}
		return
	})
	if cbmi != nil {
		return
	}
	return nil
}

func (ctm *CoolorToolMenu) forMenuItems(
	visible bool,
	f func(c *CoolorButtonMenuItem, idx int) (err error),
) {
	if ctm == nil || ctm.Flex == nil {
		return
	}
	items := ctm.menuItems[0:]
	if visible {
		items = ctm.visibleItems[0:]
	}
	if len(items) > 0 {
		for i := len(items) - 1; i >= 0; i-- {
			v := items[i]
			if v == nil {
				continue
			}
			e := f(v, i)
			if e != nil {
				break
			}
		}
	}
}

func (ctm *CoolorToolMenu) HandleEvent(e ObservableEvent) bool {
	switch e.Type {
	case PaletteColorSelectedEvent:
		var cc *CoolorColor = e.Ref.(*CoolorColor)
		ctm.UpdateColor(cc.Color)
	}
	return true
}

func (ctm *CoolorToolMenu) UpdateColor(col *tcell.Color) {
	if ctm == nil || col == nil {
		return
	}
	pcol, i := ctm.mc.palette.GetSelected()
	ctm.selectedColor = NewIntCoolorColor(col.Hex())
	ctm.forMenuItems(false, func(c *CoolorButtonMenuItem, idx int) (err error) {
		c.action.Every(MainC.palette.CoolorColorsPalette.Colors[i])
		c.SetBackgroundColor(*pcol.Color)
		c.SetLabelColor(pcol.GetFgColor())
		// c.action.Every(c.menu.selectedColor)
		return
	})

	ctm.updateState()
}

func (ctm *CoolorToolMenu) UpdateVisibleActors(c CoolorColorActionFlag) bool {
	ctm.activeActionFlag = c
	// MainC.app.QueueUpdateDraw(func() {
	// 	slen := len(ctm.visibleItems)
	// 	if len(ctm.visibleItems) != slen {
	// 		ctm.selected = 0
	// 	}
	// 	ctm.visibleItems = ctm.visibleItems[:0]
	// 	if c == -1 {
	// 		c = 1
	// 	}
	// 	for _, v := range ctm.menuItems {
	//      if v == nil {
	//        continue
	//      }
	// 		if c&v.action.flag != 0 {
	// 			ctm.visibleItems = append(ctm.visibleItems, v)
	// 		} else {
	// 			ctm.visibleItems = append(ctm.visibleItems, nil)
	// 		}
	// 	}
	// })
	// MainC.app.Draw()
	return true
}

func (ctm *CoolorToolMenu) Selected() *CoolorButtonMenuItem {
	if len(ctm.visibleItems) > 0 {
		if ctm.selected < 0 {
			ctm.selected = 0
		}
		if ctm.selected > len(ctm.visibleItems)-1 {
			ctm.selected = len(ctm.visibleItems) - 1
		}
		mmi := ctm.visibleItems[ctm.selected]
		return mmi
	}
	return nil
}

func (ctm *CoolorToolMenu) Activated() *CoolorColorActor {
	for _, a := range acts {
		if a.IsActivated() {
			return a
		}
	}
	return nil
}

func (ctm *CoolorToolMenu) UpdateActionStatus(mmi *CoolorButtonMenuItem) {
	if mmi == nil {
		return
	}

	status.NewStatusUpdate(
		"action",
		fmt.Sprintf(
			"[black:yellow:b] %s %s [-:-:-]",
			mmi.action.icon,
			mmi.action.name,
		),
	)
}

func (ctm *CoolorToolMenu) updateState() {
	// ctm.list.List.GetCurrentItem()
	MainC.app.QueueUpdateDraw(func() {
    if ctm.list != nil && MainC.menu != nil && MainC.menu.selectedColor != nil && MainC.menu.selectedColor.Color != nil {
      ctm.list.SetMainTextStyle(tcell.StyleDefault.Background(*MainC.menu.selectedColor.Color).Foreground(MainC.menu.selectedColor.GetFgColorShade()))
      ctm.list.SetSelectedStyle(tcell.StyleDefault.Background(*MainC.menu.selectedColor.Color).Foreground(MainC.menu.selectedColor.GetFgColor()))
    }
		ctm.forMenuItems(false, func(c *CoolorButtonMenuItem, idx int) (err error) {
			c.selected = false
			return
		})

		// if act := ctm.Activated(); act != nil {
		// 	ctm.UpdateVisibleActors(act.Actions())
		// } else {
		// 	ctm.UpdateVisibleActors(MainPaletteActionsFlag)
		// }

		ctm.forMenuItems(true, func(c *CoolorButtonMenuItem, idx int) (err error) {
			if idx == ctm.selected {
				c.selected = true
				if MainC.menu != nil && MainC.menu.selectedColor != nil {
					col := MainC.menu.selectedColor.GetFgColor()
					c.SetLabelColor(col)
				}
			} else {
				c.SetBorder(false)
				if MainC.menu != nil && MainC.menu.selectedColor != nil {
					col := MainC.menu.selectedColor.GetFgColorShade()
					c.SetLabelColor(col)
				}
			}
			c.UpdateState()
			return
		})

		mmi := ctm.Selected()
		if mmi == nil {
			return
		}
		ctm.UpdateActionStatus(nil)
		ctm.ResetViews()
	})
}

func (ctm *CoolorToolMenu) ResetViews() {
	MainC.app.QueueUpdateDraw(func() {
		ctm.list.UpdateListItems()
		ctm.Clear()
		itemOrder := ctm.RotateItems(ctm.selected)
		for _, v := range itemOrder {
			act, _ := ctm.FindAct(v)
			ctm.Flex.AddItem(act, 2, 0, false)
			// ctm.list.AddItem(mainText string, secondaryText string, shortcut rune, selected func())
		}
		// ctm.forMenuItems(true, func(c *CoolorButtonMenuItem, idx int) {
		// })
	})
	// MainC.app.Sync()
}

func (ctm *CoolorToolMenu) ActivateSelected(cc *CoolorColor) {
	sel := ctm.Selected()
	if sel == nil {
		sel = ctm.visibleItems[0]
		if sel == nil {
			return
		}
	}
	sel.action.Before(ctm.selectedColor.pallette, cc)
	// sel.action.Activate(cc)
	if ctm.Activated() == nil {
		ctm.UpdateVisibleActors(MainPaletteActionsFlag)
	} else {
		ctm.UpdateVisibleActors(ctm.Activated().Actions())
	}

	ctm.updateState()
}

func (ctm *CoolorToolMenu) FindAct(n string) (*CoolorButtonMenuItem, int) {
	var cbmi *CoolorButtonMenuItem
	i := 0
	ctm.forMenuItems(true, func(c *CoolorButtonMenuItem, idx int) (err error) {
		if c.name == n {
			cbmi = c
			i = idx
		}
		return
	})
	return cbmi, i
}

func (ctm *CoolorToolMenu) RotateItems(idx int) []string {
	if len(ctm.visibleItems) == 0 {
		return []string{}
	}
	itemOrder := MakeCenteredCircleInts(idx, len(actOrder))
	actsOrder := make([]string, len(actOrder))
	for i := range itemOrder {
		act := actOrder[i]
		actsOrder = append(actsOrder, act)
		// action := ctm.FindAct(act)
	}
	return actsOrder
}

func MakeCenteredCircleInts(center, num int) []int {
	nums := lo.RangeFrom(0, num)
	center = center % num
	nums = reverse(nums, 0, len(nums)-1)
	nums = reverse(nums, 0, center-1)
	nums = reverse(nums, center, len(nums)-1)
	return nums
	// nums = lo.Reverse[int](nums)
}

func reverse(n []int, start, end int) []int {
	// nums := make([]int, len(n))
	// nums := copy([]int{}, n)
	for start >= 0 && end >= 0 {
		if start >= end {
			break
		}
		t := n[start]
		n[start] = n[end]
		n[end] = t
		start++
		end--
	}
	return n
}

/*

function rotate3(nums, k) {
 k = k % nums.length;
 reverse(nums, 0, nums.length - 1);
 reverse(nums, 0, k - 1);
 reverse(nums, k, nums.length - 1);
 return nums;

 function reverse(nums, start, end) {
   while (start < end) {
     let temporary = nums[start];
     nums[start] = nums[end];
     nums[end] = temporary;
     start++;
     end -- ;
   }
   return nums;
 }
}

  var proportionSum int
*/

func (ctm *CoolorToolMenu) SetSelected(idx int) {
	if idx < 0 {
		idx = len(ctm.visibleItems) - 1
	}
	if idx >= len(ctm.visibleItems) {
		idx = 0
	}
	ctm.selected = int(math.Mod(float64(idx), float64(len(ctm.visibleItems))))
	ctm.updateState()
}

// NavSelection(idx int) error
func (ctm *CoolorToolMenu) NavSelection(idx int) {
	ctm.app.QueueUpdateDraw(func() {
		ctm.list.SetSelectedFocusOnly(false)
		ctm.list.SetCurrentItem(ctm.list.GetCurrentItem() + idx)
		// ctm.SetSelected(ctm.selected + idx)
	})
}

func (ctm *CoolorToolMenu) ActorDispatchor() {
}

func (cmm *CoolorToolMenu) AddItem(b *CoolorButtonMenuItem) {
	cmm.menuItems = append(cmm.menuItems, b)
	cmm.visibleItems = cmm.menuItems[0:]
	litem := ListItem(b)
	var li *ListItem = &litem
	cmm.list.AddItem(li)
	// cmm.Flex.AddItem(b, 2, 0, false)
	cmm.updateState()
}

// vim: ts=2 sw=2 et ft=go
