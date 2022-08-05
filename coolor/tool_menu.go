package coolor

import (
	"fmt"
	"math"

	"github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
)

type CoolorToolMenu struct {
	*tview.Flex
	menuItems     []*CoolorButtonMenuItem
	visibleItems  []*CoolorButtonMenuItem
	mc            *MainContainer
	app           *tview.Application
	selectedColor *CoolorColor
	selected      int
}
type CoolorPaletteAction interface {
	IsActivated() bool
	Activate(cc *CoolorColor) bool
	Before(cp *CoolorBlendPalette) bool
	Every(cp *CoolorBlendPalette, menu *CoolorToolMenu) bool
	Finalize(cc *CoolorColor, cp *CoolorPalette, menu *CoolorToolMenu)
}

type CoolorButtonMenuItem struct {
	*tview.Button
	icon, name string
	selected   bool
	menu       *CoolorToolMenu
	action     *CoolorColorActor
}

func NewButtonMenuItem(menu *CoolorToolMenu, action *CoolorColorActor) *CoolorButtonMenuItem {
	action.menu = menu
	mmi := &CoolorButtonMenuItem{
		Button:   tview.NewButton(action.name),
		name:     action.name,
		action:   action,
		selected: false,
		menu:     menu,
		icon:     action.icon,
	}
	mmi.Button.SetLabel(mmi.icon).SetBorder(false).SetTitleAlign(tview.AlignCenter)
	mmi.Button.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	mmi.Button.SetLabelColor(tview.Styles.InverseTextColor)
	mmi.Button.SetSelectedFunc(func() {
		mmi.UpdateState()
	})
	// mmi.SetBackgroundColor(tcell.Color(0))
	return mmi
}
func NewCoolorColorMainMenu(app *tview.Application) *CoolorToolMenu {
	cmm := &CoolorToolMenu{
		app:      app,
		mc:       MainC,
		Flex:     tview.NewFlex(),
		selected: 0,
	}
	// '烙'
	cmm.Flex.SetDirection(tview.FlexRow)
	return cmm
}
func (ctm *CoolorToolMenu) Init() {
	for _, v := range acts {
		ctm.AddItem(NewButtonMenuItem(ctm, v))
	}
	ctm.SetSelected(0)
	ctm.app.Draw()
}

func (mmi *CoolorButtonMenuItem) UpdateState() {
	// if mmi.action.name != mmi.name {
	//	mmi.name = mmi.action.name
	//	mmi.icon = mmi.action.icon
	// }
	if mmi.selected {
		mmi.Button.SetLabel(fmt.Sprintf(`[::b]%s[::-]`, mmi.action.icon))
	} else {
		mmi.Button.SetLabel(fmt.Sprintf(`%s`, mmi.action.icon))
	}
}
func (ctm *CoolorToolMenu) forMenuItems(visible bool, f func(c *CoolorButtonMenuItem, idx int)) {
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
			f(v, i)
		}
	}
}

func (ctm *CoolorToolMenu) UpdateColor(col *tcell.Color) {
	if ctm == nil || col == nil {
		return
	}
	ctm.selectedColor = NewIntCoolorColor(col.Hex())
	ctm.forMenuItems(false, func(c *CoolorButtonMenuItem, idx int) {
		c.SetBackgroundColor(*ctm.selectedColor.color)
		c.SetLabelColor(getFGColor(*ctm.selectedColor.color))
	})
	ctm.updateState()
}

func (ctm *CoolorToolMenu) UpdateVisibleActors(c CoolorColorActionFlag) bool {
	dump.P("should vis", c)
	ctm.Flex.Clear()
	slen := len(ctm.visibleItems)
	ctm.visibleItems = ctm.visibleItems[:0]
	if c == -1 {
		c = 1
	}
	for _, v := range ctm.menuItems {
		dump.P(sets, c, v.name, c&v.action.flag, len(ctm.visibleItems))
		if c&v.action.flag != 0 {
			ctm.visibleItems = append(ctm.visibleItems, v)
			ctm.Flex.AddItem(v, 2, 0, false)
		} else {
			ctm.Flex.AddItem(nil, 2, 0, false)
		}
	}
	if len(ctm.visibleItems) != slen {
		ctm.selected = 0
	}
	MainC.app.Draw()
	return true
}

func (ctm *CoolorToolMenu) Selected() *CoolorButtonMenuItem {
	if len(ctm.visibleItems) > 0 {
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
		if mmi = ctm.Selected(); mmi == nil {
			return
		}
	}

	status.NewStatusUpdate("action", fmt.Sprintf("[black:yellow:b] 🗲 %s %s [-:-:-]", mmi.action.icon, mmi.action.name))
}

func (ctm *CoolorToolMenu) updateState() {
	dump.P(ActionOptions.Data())
	ctm.forMenuItems(false, func(c *CoolorButtonMenuItem, idx int) {
		c.action.Every(ctm.selectedColor)
		c.selected = false
	})
	if act := ctm.Activated(); act != nil {
		dump.P(act.name)
		ctm.UpdateVisibleActors(act.Actions())
	} else {
		ctm.UpdateVisibleActors(MainPaletteActionsFlag)
		dump.P("no active")
	}
	ctm.forMenuItems(true, func(c *CoolorButtonMenuItem, idx int) {
		if idx == ctm.selected {
			c.selected = true
			if MainC.menu != nil && MainC.menu.selectedColor != nil {
				col := getFGColor(*MainC.menu.selectedColor.color)
				c.SetLabelColor(col)
			}
		} else {
			c.SetBorder(false)
			if MainC.menu != nil && MainC.menu.selectedColor != nil {
				col := getFGColor(*MainC.menu.selectedColor.color)
				c.SetLabelColor(inverseColor(col))
			}
		}
		c.UpdateState()
	})

	mmi := ctm.Selected()
	if mmi == nil {
		return
	}
	ctm.UpdateActionStatus(nil)
}

func (ctm *CoolorToolMenu) ActivateSelected(cc *CoolorColor) {
	sel := ctm.Selected()
	if sel == nil {
		sel = ctm.visibleItems[0]
	}
	sel.action.Before(ctm.selectedColor.pallette, cc)
	sel.action.Activate(cc)
	if ctm.Activated() == nil {
		ctm.UpdateVisibleActors(MainPaletteActionsFlag)
	} else {
		ctm.UpdateVisibleActors(ctm.Activated().Actions())
	}

	ctm.updateState()
}

func (ctm *CoolorToolMenu) SetSelected(idx int) {
	if idx < 0 {
		idx = len(ctm.visibleItems) - 1
	}
	ctm.selected = int(math.Mod(float64(idx), float64(len(ctm.visibleItems))))
	ctm.updateState()
}

// NavSelection(idx int) error
func (ctm *CoolorToolMenu) NavSelection(idx int) {
	ctm.SetSelected(ctm.selected + idx)
}

func (ctm *CoolorToolMenu) ActorDispatchor() {

}

func (cmm *CoolorToolMenu) AddItem(b *CoolorButtonMenuItem) {
	cmm.menuItems = append(cmm.menuItems, b)
	// cmm.visibleItems = cmm.menuItems[0:]
	cmm.Flex.AddItem(b, 2, 0, false)
	cmm.updateState()
}

// vim: ts=2 sw=2 et ft=go
