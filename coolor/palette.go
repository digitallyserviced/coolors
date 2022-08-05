package coolor

import (
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
)

type CoolorColors []*CoolorColor

type CoolorPalette struct {
	*tview.Flex
	colors      CoolorColors
	selectedIdx int
	l           *sync.RWMutex
	handlers    map[string]EventHandlers
	menu        *CoolorToolMenu
	ptype       string
}

type CoolorMainPalette struct {
	*CoolorPalette
	name string
}

type CoolorBlendPalette struct {
	*CoolorPalette
	start, end *CoolorColor
	increments float64
}

func BlankCoolorBlendPalette(start, end *CoolorColor, increments float64) *CoolorBlendPalette {
	cp := &CoolorPalette{
		Flex:        tview.NewFlex(),
		colors:      CoolorColors{},
		selectedIdx: 0,
		l:           &sync.RWMutex{},
		handlers:    make(map[string]EventHandlers),
		ptype:       "blend",
	}

	cp.Flex.SetDirection(tview.FlexColumn)
	cbp := &CoolorBlendPalette{
		CoolorPalette: cp,
		start:         start,
		end:           end,
		increments:    increments,
	}
	cbp.Init()

	return cbp
}

func (cbp *CoolorBlendPalette) UpdateColors(start, end *CoolorColor) {
	cbp.start = start
	cbp.end = end
	cbp.Init()
}

func (cbp *CoolorBlendPalette) Init() {
	cbp.Clear()
	cbp.colors = make(CoolorColors, 0)
	incrSizes := 1.0 / cbp.increments
	start, _ := MakeColor(cbp.start)
	end, _ := MakeColor(cbp.end)
	for i := 0; i <= int(cbp.increments); i++ {
		newc := start.BlendLuv(end, float64(i)*float64(incrSizes))
		newcc := NewStaticCoolorColor(newc.Hex())
		cbp.CoolorPalette.AddCoolorColor(newcc)
	}

	MainC.conf.AddPalette("blend", cbp)
	// cbp.SetSelected(0)
}

func BlankCoolorPalette() *CoolorPalette {
	cp := &CoolorPalette{
		Flex:        tview.NewFlex(),
		colors:      CoolorColors{},
		selectedIdx: 0,
		l:           &sync.RWMutex{},
		handlers:    make(map[string]EventHandlers),
		// menu:        MainC.menu,
		ptype: "regular",
	}
	cp.Flex.SetDirection(tview.FlexColumn)

	return cp
}

func DefaultCoolorPalette() *CoolorMainPalette {
	tcols := GenerateRandomColors(5)
	cmp := NewCoolorPaletteWithColors(tcols)
	// cmp := &CoolorMainPalette{
	// 	CoolorPalette: cp,
	// 	name:          "",
	// }
	return cmp
}

func NewCoolorPaletteFromMap(cols map[string]string) *CoolorMainPalette {
	cp := BlankCoolorPalette()
	for n, v := range cols {
		col := cp.AddCssCoolorColor(v)
		col.SetName(n)
	}
	cp.SetSelected(0)
	cmp := &CoolorMainPalette{
		CoolorPalette: cp,
		name:          "random untitled",
	}
	return cmp
}

func NewCoolorPaletteFromCssStrings(cols []string) *CoolorMainPalette {
	cp := BlankCoolorPalette()
	for _, v := range cols {
		cp.AddCssCoolorColor(v)
	}
	cp.SetSelected(0)
	cmp := &CoolorMainPalette{
		CoolorPalette: cp,
		name:          "random untitled",
	}
	return cmp
}

func NewCoolorPaletteWithColors(tcols []tcell.Color) *CoolorMainPalette {
	cp := BlankCoolorPalette()
	for _, v := range tcols {
		// fmt.Printf("%06x", v.Hex())
		cp.AddCssCoolorColor(fmt.Sprintf("#%06x", v.Hex()))
	}
	cp.SetSelected(0)
	cmp := &CoolorMainPalette{
		CoolorPalette: cp,
		name:          "defined untitled",
	}
	return cmp
}

func (cp *CoolorPalette) AddCoolorColor(color *CoolorColor) *CoolorColor {
	color.pallette = cp
	cp.l.Lock()
	defer cp.l.Unlock()
	cp.colors = append(cp.colors, color)
	cp.Flex.AddItem(cp.colors[len(cp.colors)-1], 0, 1, true)
	return cp.colors[len(cp.colors)-1]
}

func (cp *CoolorPalette) AddCssCoolorColor(c string) *CoolorColor {
	color := cp.AddCoolorColor(NewCoolorColor(c))
	return color
}

func (cp *CoolorPalette) GetColorAt(idx int) *CoolorColor {
	return cp.colors[int(math.Mod(float64(idx), float64(len(cp.colors))))]
}
func (cp *CoolorPalette) SetMenu(menu *CoolorToolMenu) {
	cp.menu = menu
	cc, _ := cp.GetSelected()
	cp.menu.UpdateColor(cc.color)
}

func (cp *CoolorPalette) RandomColor() *CoolorColor {
	return cp.GetColorAt(int(rand.Uint32()))
}

func (cp *CoolorPalette) AddRandomCoolorColor() *CoolorColor {
	newc := NewRandomCoolorColor()
	newc.pallette = cp
	cp.colors = append(cp.colors, newc)
	cp.Flex.AddItem(newc, 0, 1, false)
	return newc
}

func (cp *CoolorPalette) GetSelected() (*CoolorColor, int) {
	cp.l.RLock()
	defer cp.l.RUnlock()

	if cp.colors[cp.selectedIdx] != nil {
		return cp.colors[cp.selectedIdx], cp.selectedIdx
	}
	return nil, -1
}

func (cp *CoolorPalette) NavSelection(idx int) {
	cp.l.RLock()
	newidx := cp.selectedIdx + idx
	if newidx >= len(cp.colors) {
		newidx = len(cp.colors) - 1
	}
	if newidx < 0 {
		newidx = 0
	}
	cp.l.RUnlock()
	cp.SetSelected(newidx)
}

func (cp *CoolorPalette) Randomize() int {
	changed := 0
	for _, v := range cp.colors {
		if v.Random() {
			changed += 1
		}
	}
	MainC.conf.AddPalette("random", cp)
	cp.ResetViews()
	return changed
}

func (cp *CoolorPalette) ResetViews() {
	cp.SetSelected(0)
}

func (cc *CoolorPalette) SpawnSelectionEvent(c *CoolorColor, idx int) bool {
	if len(cc.handlers["selected"]) > 0 {
		ev := &SelectionEvent{
			color: c,
			idx:   int8(idx),
		}
		for _, v := range cc.handlers["selected"] {
			if v != nil {
				eh, ok := v.(tcell.EventHandler)
				if !ok {
					panic(ok)
				}
				eh.HandleEvent(ev)
			}
		}
	}
	return true
}

func (cc *CoolorPalette) AddEventHandler(t string, h tcell.EventHandler) {
	cc.l.Lock()
	defer cc.l.Unlock()

	if cc.handlers[t] == nil {
		cc.handlers[t] = make(EventHandlers, 0)
	}
	cc.handlers[t] = append(cc.handlers[t], h)
}

func (cp *CoolorPalette) SetSelected(idx int) error {
	sel, _ := cp.GetSelected()
	sel.SetSelected(false)
	if idx < len(cp.colors) {
		if idx < 0 {
			idx = 0
		}
		cp.l.Lock()
		cp.selectedIdx = idx
		cp.colors[cp.selectedIdx].SetSelected(true)
		cp.menu.UpdateColor(cp.colors[cp.selectedIdx].color)
		cp.l.Unlock()
		cp.SpawnSelectionEvent(cp.colors[cp.selectedIdx], cp.selectedIdx)
		return nil
	}
	return fmt.Errorf("No valid color at idx: %d", idx)
}

func (cp *CoolorPalette) Plainify(s bool) {
	for _, v := range cp.colors {
		v.SetStatic(s)
		v.SetPlain(s)
	}
}

func (cp *CoolorPalette) Staticize(s bool) {
	for _, v := range cp.colors {
		v.SetStatic(s)
	}
}

func (cp *CoolorPalette) ToggleLockSelected() (*CoolorColor, int) {
	cc, _ := cp.GetSelected()
	cc.ToggleLocked()
	return cp.colors[cp.selectedIdx], cp.selectedIdx
}

type Palette interface {
	GetPalette() *CoolorPalette
}

func (cp *CoolorPalette) GetPalette() *CoolorPalette {
	return cp
}
func (cbp *CoolorBlendPalette) GetPalette() *CoolorPalette {
	return cbp.CoolorPalette
}

func (cp *CoolorMainPalette) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cp.Flex.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		ch := event.Rune()
		kp := event.Key()
		_ = kp
		switch {
		case ch == ' ' || ch == 'R':
			cp.Randomize()
		case ch == 'w':
			cp.ToggleLockSelected()
		case ch == 'r':
			color, _ := cp.GetSelected()
			color.Random()
		case ch == '+': // Add a color
			cp.AddRandomCoolorColor()
		case ch == '-': // Remove a color
			remcolor, idx := cp.GetSelected()
			cp.SetSelected(idx - 1)
			remcolor.Remove()
		}
		// if handler := cp.InputHandler(); handler != nil {
		// 	dump.P(fmt.Sprintf("%s input handled", cp.ptype))
		// 	// handler(event, setFocus)
		// 	return
		// }
	})
}

// vim: ts=2 sw=2 et ft=go
