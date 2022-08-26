package coolor

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	// "github.com/gookit/goutil/dump"
	// "github.com/gookit/goutil/dump"
)

type CoolorColors []*CoolorColor


type CoolorPalette struct {
	ColorContainer *tview.Flex
	*tview.Flex
	colors             CoolorColors
	paddles            []*PalettePaddle
	colSize, maxColors int
	selectedIdx        int
	l                  *sync.RWMutex
	handlers           map[string]EventHandlers
	menu               *CoolorToolMenu
	ptype              string
}

type CoolorMainPalette struct {
	*CoolorPalette
	name string
}

type CoolorShadePalette struct {
	*CoolorPalette
	base *CoolorColor
	increments float64
}

func BlankCoolorShadePalette(base *CoolorColor, increments float64) *CoolorShadePalette {
	cp := &CoolorPalette{
		ColorContainer: tview.NewFlex(),
		Flex:           tview.NewFlex(),
		colors:         CoolorColors{},
		paddles:        NewPaddles(),
		colSize:        12,
		maxColors:      8,
		selectedIdx:    0,
		l:              &sync.RWMutex{},
		handlers:       make(map[string]EventHandlers),
		menu:           &CoolorToolMenu{},
		ptype:          "shade",
	}

	cp.ColorContainer.SetDirection(tview.FlexColumn)
	cp.Flex.SetDirection(tview.FlexColumn)
	cp.Flex.AddItem(cp.paddles[0], 4, 0, false)
	cp.Flex.AddItem(cp.ColorContainer, 80, 0, true)
	cp.Flex.AddItem(cp.paddles[1], 4, 0, false)
	cbp := &CoolorShadePalette{
		CoolorPalette: cp,
		base:          base,
		increments:    increments,
	}
	cbp.Init()

	return cbp
}
type CoolorBlendPalette struct {
	*CoolorPalette
	start, end *CoolorColor
	increments float64
}

func BlankCoolorBlendPalette(start, end *CoolorColor, increments float64) *CoolorBlendPalette {
	cp := &CoolorPalette{
		ColorContainer: tview.NewFlex(),
		Flex:           tview.NewFlex(),
		colors:         CoolorColors{},
		paddles:        NewPaddles(),
		colSize:        12,
		maxColors:      8,
		selectedIdx:    0,
		l:              &sync.RWMutex{},
		handlers:       make(map[string]EventHandlers),
		menu:           &CoolorToolMenu{},
		ptype:          "blend",
	}

	cp.ColorContainer.SetDirection(tview.FlexColumn)
	cp.Flex.SetDirection(tview.FlexColumn)
	cp.Flex.AddItem(cp.paddles[0], 4, 0, false)
	cp.Flex.AddItem(cp.ColorContainer, 80, 0, true)
	cp.Flex.AddItem(cp.paddles[1], 4, 0, false)
	cbp := &CoolorBlendPalette{
		CoolorPalette: cp,
		start:         start,
		end:           end,
		increments:    increments,
	}
	cbp.Init()

	return cbp
}

func (ccs *CoolorPalette) Swap(a, b int) {
	if col := ccs.colors[a]; col == nil {
		return
	}
	if col := ccs.colors[b]; col == nil {
		return
	}
	if ccs.selectedIdx == a {
		ccs.colors[b].SetSelected(true)
		ccs.menu.UpdateColor(ccs.colors[b].color)
		ccs.colors[a].SetSelected(false)
	} else {
		ccs.colors[a].SetSelected(true)
		ccs.menu.UpdateColor(ccs.colors[a].color)
		ccs.colors[b].SetSelected(false)
	}
	ccs.colors[a], ccs.colors[b] = ccs.colors[b], ccs.colors[a]
	// ccs.selectedIdx = b
	// ccs.SetSelected(a)
}

func (ccs *CoolorPalette) Less(a, b int) bool {
	return ccs.colors[a].color.Hex() < ccs.colors[b].color.Hex()
}

func (ccs *CoolorPalette) Len() int {
	return len(ccs.colors)
}

func (cbp *CoolorShadePalette) UpdateColors(base *CoolorColor) {
  cbp.base = base
	cbp.Init()
}
func (cbp *CoolorBlendPalette) UpdateColors(start, end *CoolorColor) {
	cbp.start = start
	cbp.end = end
	cbp.Init()
}

func (cbp *CoolorShadePalette) Init() {
	cbp.ColorContainer.Clear()
	base, _ := MakeColor(cbp.base)
	// tcol, _ := Hex("#be1685")
  done := make(chan struct{})
  defer close(done)
	colors := RandomShadesStream(base, 0.15)
  colors.Status.SetProgressHandler(NewProgressHandler(func(u uint32) {
		status.NewStatusUpdate("action", fmt.Sprintf("Found Shades (%d / %d)", u, colors.Status.GetItr()))
    // dump.P(u,colors.Status.Itr)
  },func(i uint32) {
		status.NewStatusUpdate("action_str", fmt.Sprintf("Iterating Shades (%d)", i))
  }))
	cbp.colors = make(CoolorColors, 0)
  colors.Run(done)
  for _, v := range TakeNColors(done, colors.OutColors, int(cbp.increments)) {
		newcc := NewStaticCoolorColor(v.Hex())
		cbp.CoolorPalette.AddCoolorColor(newcc)
  }
  cbp.GetPalette().Sort()

	MainC.conf.AddPalette("shades", cbp)
	cbp.SetSelected(0)
}
func (cbp *CoolorBlendPalette) Init() {
	cbp.ColorContainer.Clear()
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

func NewPaddles() []*PalettePaddle {
	//  ﰯ  ﰭ  鹿      ﲕ     ﮾      ﰬ ﰳ  ﯀    壟     ﰷ ﰮ     ﬕ ﯁ ﲐ  ﬔ    ﲓ      ﰵ      ﮿    ﰰ ﰴ  ﲔ ﲒ         ﲑ               ﲗ ﲖ ﰲ ﰶ  ﰱ 
	left := NewPalettePaddle("", "")
	right := NewPalettePaddle("", "")
	return []*PalettePaddle{left, right}
}

func BlankCoolorPalette() *CoolorPalette {
	cp := &CoolorPalette{
		ColorContainer: tview.NewFlex(),
		Flex:           tview.NewFlex(),
		colors:         make(CoolorColors, 0),
		colSize:        12,
		maxColors:      8,
		selectedIdx:    0,
		l:              &sync.RWMutex{},
		handlers:       make(map[string]EventHandlers),
		paddles:        NewPaddles(),
		// menu:        MainC.menu,
		ptype: "regular",
	}
	cp.ColorContainer.SetDirection(tview.FlexColumn)
	cp.Flex.SetDirection(tview.FlexColumn)
	cp.Flex.AddItem(cp.paddles[0], 4, 0, false)
	cp.Flex.AddItem(cp.ColorContainer, 80, 0, true)
	cp.Flex.AddItem(cp.paddles[1], 4, 0, false)
	return cp
}

func DefaultCoolorPalette() *CoolorMainPalette {
	tcols := GenerateRandomColors(5)
	cmp := NewCoolorPaletteWithColors(tcols)
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
	// cp.l.Unlock()
	cp.colors = append(cp.colors, color.Clone())
	cp.ColorContainer.AddItem(cp.colors[len(cp.colors)-1], 0, 1, false)
	cp.l.Unlock()
	cp.SetSelected(len(cp.colors) - 1)
	cp.UpdateSize()
	cp.ResetViews()
	MainC.conf.AddPalette(fmt.Sprintf("current_%x", time.Now().Unix()), cp)
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

func (cp *CoolorPalette) RemoveItem(rcc *CoolorColor) {
	newColors := cp.colors[:0]
	cp.Each(func(cc *CoolorColor, i int) {
		if cc != rcc {
			newColors = append(newColors, cc)
		}
	})
	cp.colors = newColors
	cp.ResetViews()
	MainC.conf.AddPalette(fmt.Sprintf("current_%x", time.Now().Unix()), cp)
}

func (cp *CoolorPalette) RandomColor() *CoolorColor {
	return cp.GetColorAt(int(rand.Uint32()))
}

func (cp *CoolorPalette) AddRandomCoolorColor() *CoolorColor {
	newc := NewRandomCoolorColor()
	newc.pallette = cp
	cp.colors = append(cp.colors, newc)
	cp.ColorContainer.AddItem(newc, 0, 1, false)
	cp.UpdateSize()
	MainC.conf.AddPalette(fmt.Sprintf("current_%x", time.Now().Unix()), cp)
	return newc
}

func (cp *CoolorPalette) GetSelected() (*CoolorColor, int) {
	cp.l.RLock()
	defer cp.l.RUnlock()

	if cp.selectedIdx > len(cp.colors)-1 {
		cp.selectedIdx = 0
	}
	if cp.selectedIdx < 0 {
		cp.selectedIdx = len(cp.colors) - 1
	}
	if cp.colors[cp.selectedIdx] != nil {
		return cp.colors[cp.selectedIdx], cp.selectedIdx
	}
	return nil, -1
}

func (cp *CoolorPalette) Sort() {
	MainC.app.QueueUpdateDraw(func() {
		sort.Sort(cp)
	})
}

func (cp *CoolorPalette) UpdateDots(dots []string) {
	status.NewStatusUpdate("dots", strings.Join(dots, " "))
}

func (cp *CoolorPalette) NavSelection(idx int) {
	cp.l.RLock()
	newidx := cp.selectedIdx + idx
	if newidx >= len(cp.colors) {
		newidx = 0
	}
	if newidx < 0 {
		newidx = len(cp.colors) - 1
	}
	cp.l.RUnlock()
	cp.SetSelected(newidx)
}

func (cp *CoolorPalette) Randomize() int {
	changed := 0
	cp.Each(func(cc *CoolorColor, i int) {
		if cc.Random() {
			changed += 1
		}
	})
	MainC.conf.AddPalette("random", cp)
	cp.ResetViews()
	return changed
}

// ResetViews
// Clear palette container and add back currently visible window
// and pagination dots and scroll left and right paddles
func (cp *CoolorPalette) ResetViews() {
	MainC.app.QueueUpdateDraw(func() {
		cp.ColorContainer.Clear()
		max := math.Max(float64(cp.selectedIdx), float64(cp.maxColors-1))
		min := math.Max(0, max-float64(cp.maxColors-1))
		if min > 0 {
			cp.paddles[0].SetStatus("enabled")
		} else {
			cp.paddles[0].SetStatus("disabled")
		}
		if int(max) < len(cp.colors)-1 { // && cp.selectedIdx < len(cp.colors)
			cp.paddles[1].SetStatus("enabled")
		} else {
			cp.paddles[1].SetStatus("disabled")
		}
		//     ﱣ   ﳁ ﭜ   ﳂ  ﱤ     喇    ﴞ           
		dots := make([]string, len(cp.colors))
		//   
		for i, v := range cp.colors {
			dots[i] = fmt.Sprintf("[%s:-:-]ﱤ[-:-:-]", v.Html())
			if i == cp.selectedIdx {
				dots[i] = fmt.Sprintf("[%s:-:b][-:-:-]", v.Html())
			}
			if i < int(min) || i > int(max) {
				continue
			}
			cp.ColorContainer.AddItem(v, cp.colSize, 0, false)
		}
		cp.UpdateDots(dots)
	})
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

var PaddleMinWidth int = 4

func (cp *CoolorPalette) UpdateSize() {
	MainC.app.QueueUpdateDraw(func() {
		x, y, w, h := cp.Flex.GetInnerRect()
		_, _, _, _ = x, y, w, h
		x, y, w, h = cp.Flex.GetRect()
		cp.maxColors = (w - (PaddleMinWidth * 2)) / cp.colSize
		if len(cp.colors) < cp.maxColors {
			cp.colSize = (w - (PaddleMinWidth * 2)) / len(cp.colors)
			cp.maxColors = (w - (PaddleMinWidth * 2)) / cp.colSize
		}
		overflow := (w - (PaddleMinWidth * 2)) % cp.colSize
		left, right := PaddleMinWidth, PaddleMinWidth
		left += overflow / 2
		right += overflow / 2
		if overflow%2 != 0 {
			right += 1
		}
		contW := cp.maxColors * cp.colSize
		cp.Flex.Clear()
		cp.Flex.AddItem(cp.paddles[0], left, 0, false)
		cp.Flex.AddItem(cp.ColorContainer, contW, 0, true)
		cp.Flex.AddItem(cp.paddles[1], right, 0, false)
	})
}

func (cp *CoolorPalette) Draw(screen tcell.Screen) {
	// cp.UpdateSize()
	// num := cp.Flex.GetItemCount()
	// for i := 0; i < num; i++ {
	// 	it := cp.Flex.GetItem(i)
	// 	it.Draw(screen)
	// }
	cp.Flex.Draw(screen)
}

func (cc *CoolorPalette) AddEventHandler(t string, h tcell.EventHandler) {
	cc.l.Lock()
	defer cc.l.Unlock()

	if cc.handlers[t] == nil {
		cc.handlers[t] = make(EventHandlers, 0)
	}
	cc.handlers[t] = append(cc.handlers[t], h)
}

func (cp *CoolorPalette) ClearSelected() {
	cp.Each(func(cc *CoolorColor, i int) {
		cc.SetSelected(false)
	})
}

func (cp *CoolorPalette) SetSelected(idx int) error {
	// sel, _ := cp.GetSelected()
	// sel.SetSelected(false)
  MainC.app.QueueUpdateDraw(func() {
	cp.ClearSelected()
	if idx < 0 {
		idx = len(cp.colors) - 1
	}
	if idx > len(cp.colors)-1 {
		idx = 0
	}
	if idx < len(cp.colors) {
		// cp.l.Lock()
		cp.selectedIdx = idx
		cp.colors[cp.selectedIdx].SetSelected(true)
		cp.menu.UpdateColor(cp.colors[cp.selectedIdx].color)
		// cp.l.Unlock()
		cp.ResetViews()
		cp.SpawnSelectionEvent(cp.colors[cp.selectedIdx], cp.selectedIdx)
	}
  })

	if idx < len(cp.colors) && idx >= 0 {
    return nil
  }
	return fmt.Errorf("No valid color at idx: %d", idx)
}

func (cp *CoolorPalette) Each(f func(*CoolorColor, int)) {
	for i, v := range cp.colors {
		f(v, i)
	}
}

func (cp *CoolorPalette) Plainify(s bool) {
	cp.Each(func(cc *CoolorColor, i int) {
		cc.SetStatic(s)
		cc.SetPlain(s)
	})
	// for _, v := range cp.colors {
	// 	v.SetStatic(s)
	// 	v.SetPlain(s)
	// }
}

func (cp *CoolorPalette) Staticize(s bool) {
	cp.Each(func(cc *CoolorColor, i int) {
		cc.SetStatic(s)
	})
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
	return cp.ColorContainer.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		ch := event.Rune()
		kp := event.Key()
		_ = kp
		switch {
		case ch == '*':
			cp.Randomize()
		case ch == '+': // Add a color
			cp.AddRandomCoolorColor()
			cp.UpdateSize()
		case ch == '=':
			cp.GetPalette().Sort()
			cp.UpdateSize()
		case ch == 'd':
			color, _ := cp.GetSelected()
			cp.AddCoolorColor(color.Clone())
			cp.UpdateSize()
		}
		// if handler := cp.InputHandler(); handler != nil {
		// 	dump.P(fmt.Sprintf("%s input handled", cp.ptype))
		// 	// handler(event, setFocus)
		// 	return
		// }
	})
}

// vim: ts=2 sw=2 et ft=go
