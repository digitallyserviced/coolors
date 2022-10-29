package coolor

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/status"
)

var PaddleMinWidth int = 0

type (
	ColorHash   uint64
	PaletteHash struct {
		ColorHash
		Colors []uint64
	}
	CoolorColors []*CoolorColor
	Palette      interface {
		GetPalette() *CoolorColorsPalette
	}
	CoolorColorsPalette struct {
	l              *sync.RWMutex
	*EventObserver `gorm:"-"`
	*EventNotifier `gorm:"-"`
	tagType        *TagType
	Name           string
	ptype          string
    Colors         CoolorColors `gorm:"foreignKey:Color"`
	Hash           uint64 `gorm:"uniqueIndex"`
	selectedIdx    int
}
	CoolorPaletteMainView struct {
		ColorContainer *tview.Flex
		// ScrollLine     *ScrollLine
		*tview.Flex
		handlers map[string]EventHandlers
		menu     *CoolorToolMenu
		*CoolorColorsPalette
		paddles   []*PalettePaddle
		maxColors int
		colSize   int
	}
	CoolorMainPalette struct {
		*CoolorPaletteMainView
		name string
	}
	CoolorShadePalette struct {
		*CoolorPaletteMainView
		base       *CoolorColor
		increments float64
	}
	CoolorBlendPalette struct {
		*CoolorPaletteMainView
		start, end *CoolorColor
		increments float64
	}
)

func NewPaddles() []*PalettePaddle {
	//  ﰯ  ﰭ  鹿      ﲕ     ﮾      ﰬ ﰳ  ﯀    壟     ﰷ ﰮ     ﬕ ﯁ ﲐ  ﬔ    ﲓ      ﰵ      ﮿    ﰰ ﰴ  ﲔ ﲒ         ﲑ               ﲗ ﲖ ﰲ ﰶ  ﰱ 
	// left := NewPalettePaddle("", "")
	// right := NewPalettePaddle("", "") ﲑ  
	// left := NewPalettePaddle("", "")
	// right := NewPalettePaddle("", "ﰲ")
	left := NewPalettePaddle("", "")
	right := NewPalettePaddle("", "")
	return []*PalettePaddle{left, right}
}

func NewCoolorColorsPaletteFromMeta(
	ccm *CoolorColorsPaletteMeta,
) *CoolorMainPalette {
	colors := make([]string, len(ccm.Current.Colors))
	for i, v := range ccm.Current.Colors {
		// colors = append(colors, v.Clone())
		colors[i] = v.Escalate().Html()
	}
	ccp := NewCoolorPaletteFromCssStrings(colors)
	ccp.CoolorColorsPalette.Name = ccm.Named
	ccp.UpdateHash()
	// ccp := ccm.Current
	// ccp.Colors = make(CoolorColors, 0)
	// ccp.l = &sync.RWMutex{}
	// ccp.EventObserver = NewEventObserver("palette")
	// ccp.EventNotifier = NewEventNotifier("palette")
	// cmp := &CoolorMainPalette{
	// 	CoolorPaletteMainView: LoadedCoolorPalette(ccp),
	// 	name:                  "random untitled",
	// }
	//  // cmp.UpdateSize()
	//  // cmp.ResetViews()
	// for _, v := range colors {
	//   cmp.AddCssCoolorColor(v)
	// }
	// ccp.SetSelected(0)
	// ccp.UpdateHash()
	MainC.conf.Meta = append(MainC.conf.Meta, ccm)
	return ccp
}

func NewCoolorColorsPalette() *CoolorColorsPalette {
	ccp := &CoolorColorsPalette{
		l:             &sync.RWMutex{},
		EventObserver: NewEventObserver("palette"),
		EventNotifier: NewEventNotifier("palette"),
		Name:          Generator().GenerateName(2),
		ptype:         "",
		Colors:        make(CoolorColors, 0),
		Hash:          0,
		selectedIdx:   0,
		tagType:       &Base16Tags,
	}
	ccm := NewCoolorColorsPaletteMeta(Generator().GenerateName(2), ccp)
	ccm.Update()
	// MainC.conf.Meta = append(MainC.conf.Meta, ccm)
	return ccp
}

const maxColSize int = 12

func LoadedCoolorPalette(ccp *CoolorColorsPalette) *CoolorPaletteMainView {
	cp := &CoolorPaletteMainView{
		CoolorColorsPalette: ccp,
		ColorContainer:      tview.NewFlex(),
		Flex:                tview.NewFlex(),
		colSize:             maxColSize,
		maxColors:           8,
		// selectedIdx:    0,
		// l:              &sync.RWMutex{},
		handlers: make(map[string]EventHandlers),
		paddles:  NewPaddles(),
		// menu:        MainC.menu,
		// ptype: "regular",
	}
	cp.ColorContainer.SetDirection(tview.FlexColumn)
	cp.SetDirection(tview.FlexColumn)
	// cp.SetBorder(true).SetBorderPadding(0, 0, 0, 0)
	cp.SetTitle("[black:purple:b] Palette [-:-:-]")
	cp.AddItem(cp.paddles[0], PaddleMinWidth, 0, false)
	cp.AddItem(cp.ColorContainer, 80, 0, true)
	cp.AddItem(cp.paddles[1], PaddleMinWidth, 0, false)
	return cp
}


func BlankCoolorPalette() *CoolorPaletteMainView {
	cp := &CoolorPaletteMainView{
		CoolorColorsPalette: NewCoolorColorsPalette(),
		ColorContainer:      tview.NewFlex(),
		// ScrollLine: btmLine,
		Flex:      tview.NewFlex(),
		colSize:   maxColSize,
		maxColors: 8,
		// selectedIdx:    0,
		// l:              &sync.RWMutex{},
		handlers: make(map[string]EventHandlers),
		paddles:  NewPaddles(),
		// menu:        MainC.menu,
		// ptype: "regular",
	}
	// cp.ScrollLine = NewScrollLine(cp.CoolorColorsPalette)
	cp.ColorContainer.SetDirection(tview.FlexColumn)
	cp.SetDirection(tview.FlexColumn)
	// cp.SetBorder(true).SetBorderPadding(0, 0, 0, 0)
	cp.SetTitle("[black:purple:b] Palette [-:-:-]")
	cp.AddItem(cp.paddles[0], PaddleMinWidth, 0, false)
	cp.AddItem(cp.ColorContainer, 80, 0, true)
	cp.AddItem(cp.paddles[1], PaddleMinWidth, 0, false)
	return cp
}

func BlankCoolorShadePalette(
	base *CoolorColor,
	increments float64,
) *CoolorShadePalette {
	cp := &CoolorPaletteMainView{
		ColorContainer:      tview.NewFlex(),
		Flex:                tview.NewFlex(),
		CoolorColorsPalette: NewCoolorColorsPalette(),
		paddles:             NewPaddles(),
		colSize:             maxColSize,
		maxColors:           8,
		// selectedIdx:    0,
		// l:              &sync.RWMutex{},
		handlers: make(map[string]EventHandlers),
		// menu:           &CoolorToolMenu{},
		// ptype: "shade",
	}

	cp.AddCoolorColor(base)
	cp.ColorContainer.SetDirection(tview.FlexColumn)
	cp.SetDirection(tview.FlexColumn)
	cp.AddItem(cp.paddles[0], PaddleMinWidth, 0, false)
	cp.AddItem(cp.ColorContainer, 80, 0, true)
	cp.AddItem(cp.paddles[1], PaddleMinWidth, 0, false)
	cbp := &CoolorShadePalette{
		CoolorPaletteMainView: cp,
		base:                  base,
		increments:            increments,
	}
	go cbp.Init()

	return cbp
}

func BlankCoolorBlendPalette(
	start, end *CoolorColor,
	increments float64,
) *CoolorBlendPalette {
	cp := &CoolorPaletteMainView{
		ColorContainer:      tview.NewFlex(),
		Flex:                tview.NewFlex(),
		CoolorColorsPalette: NewCoolorColorsPalette(),
		paddles:             NewPaddles(),
		colSize:             maxColSize,
		maxColors:           8,
		// selectedIdx:    0,
		// l:              &sync.RWMutex{},
		handlers: make(map[string]EventHandlers),
		menu:     &CoolorToolMenu{},
		// ptype:          "blend",
	}

	cp.ColorContainer.SetDirection(tview.FlexColumn)
	cp.SetDirection(tview.FlexColumn)
	cp.AddItem(cp.paddles[0], PaddleMinWidth, 0, false)
	cp.AddItem(cp.ColorContainer, 80, 0, true)
	cp.AddItem(cp.paddles[1], PaddleMinWidth, 0, false)
	cbp := &CoolorBlendPalette{
		CoolorPaletteMainView: cp,
		start:                 start,
		end:                   end,
		increments:            increments,
	}
	go cbp.Init()

	return cbp
}

func DefaultCoolorPalette() *CoolorMainPalette {
	tcols := GenerateRandomColors(5)
	cmp := NewCoolorPaletteWithColors(tcols)
	return cmp
}

func NewCoolorPaletteFromMap(cols map[string]string) *CoolorColorsPalette {
	// cp := BlankCoolorPalette()
	cp := NewCoolorColorsPalette()
	for n, v := range cols {
		col := cp.AddCoolorColor(NewCoolorColor(v))
		col.SetName(n)
	}
	cp.SetSelected(0)
	return cp
}

func NewCoolorColorsPaletteFromCssStrings(cols []string) *CoolorColorsPalette {
	ccp := NewCoolorColorsPalette()
	for _, v := range cols {
		ccp.AddCoolorColor(NewCoolorColor(v))
	}
	return ccp
}
func NewCoolorPaletteFromCssStrings(cols []string) *CoolorMainPalette {
	cp := BlankCoolorPalette()
	cp.Name = Generator().GenerateName(2)
	for _, v := range cols {
		cp.AddCssCoolorColor(v)
	}
	cp.SetSelected(0)
	cmp := &CoolorMainPalette{
		CoolorPaletteMainView: cp,
		name:                  "random untitled",
	}
	return cmp
}

func NewCoolorPaletteWithColors(tcols []tcell.Color) *CoolorMainPalette {
	cp := BlankCoolorPalette()
	for _, v := range tcols {
		// fmt.Printf("%06x", v.Hex())
		cp.AddCssCoolorColor(fmt.Sprintf("#%06x", v.Hex()))
		SeentColor("startup", NewIntCoolorColor(v.Hex()), cp)
	}

	cp.Name = Generator().WithSeed(int64(cp.UpdateHash())).GenerateName(2)
	cp.SetSelected(0)
	cmp := &CoolorMainPalette{
		CoolorPaletteMainView: cp,
		name:                  "defined untitled",
	}
	return cmp
}

func K(from string, cc *CoolorColor, src Referenced) {
	// fmt.Println(from, cc)
	if MainC == nil || MainC.EventNotifier == nil {
		return
	}
	MainC.EventNotifier.Notify(
		*MainC.EventNotifier.NewObservableEvent(ColorSeentEvent, from, cc, src),
	)
}
func (cp *CoolorPaletteMainView) AddCoolorColor(
	color *CoolorColor,
) *CoolorColor {
	color = cp.CoolorColorsPalette.AddCoolorColor(color.Clone())
	color.pallette = cp
	cp.ColorContainer.AddItem(cp.Colors[len(cp.Colors)-1], 0, 1, false)
	cp.SetSelected(len(cp.Colors) - 1)
	return cp.Colors[len(cp.Colors)-1]
}

func (cp *CoolorPaletteMainView) AddCssCoolorColor(c string) *CoolorColor {
	color := cp.AddCoolorColor(NewCoolorColor(c))
	return color
}

func (cp *CoolorPaletteMainView) SetMenu(menu *CoolorToolMenu) {
	cp.menu = menu
	cc, i := cp.GetSelected()
	if cc == nil {
		return
	}

	cp.CoolorColorsPalette.Register(PaletteColorSelectedEvent, cp.menu)
	cp.SpawnSelectionEvent(cc, i)
}

func (cp *CoolorPaletteMainView) RemoveItem(rcc *CoolorColor) {
	newColors := cp.Colors[:0]
	cp.Each(func(cc *CoolorColor, _ int) {
		if cc != rcc {
			newColors = append(newColors, cc)
		}
	})
	cp.Colors = newColors
	cp.PaletteEvent(PaletteColorRemovedEvent, rcc)
	cp.ResetViews()
	// MainC.conf.AddPalette(fmt.Sprintf("current_%x", time.Now().Unix()), cp)
}

func (cp *CoolorPaletteMainView) AddRandomCoolorColor() *CoolorColor {
	newc := NewRandomCoolorColor()
	newc.pallette = cp
	cp.Colors = append(cp.Colors, newc)
	cp.ColorContainer.AddItem(newc, 0, 1, false)
	cp.UpdateSize()
	cp.ResetViews()
	// MainC.conf.AddPalette(fmt.Sprintf("current_%x", time.Now().Unix()), cp)
	return newc
}

func (cp *CoolorPaletteMainView) Sort() {
	MainC.app.QueueUpdateDraw(func() {
		cp.CoolorColorsPalette.Sort()
		cp.UpdateSize()
		cp.ResetViews()
	})
}

func (cp *CoolorPaletteMainView) UpdateDots(dots []string) {
	status.NewStatusUpdate("dots", strings.Join(dots, ""))
}

func (cp *CoolorPaletteMainView) NavSelection(idx int) {
	cp.l.RLock()
	newidx := cp.selectedIdx + idx
	if newidx >= len(cp.Colors) {
		newidx = 0
	}
	if newidx < 0 {
		newidx = len(cp.Colors) - 1
	}
	cp.l.RUnlock()
	cp.SetSelected(newidx)
}

func (cp *CoolorPaletteMainView) Randomize() int {
	changed := 0
	cp.Each(func(cc *CoolorColor, _ int) {
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
func (cp *CoolorPaletteMainView) ResetViews() {
	MainC.app.QueueUpdateDraw(func() {
		cp.ColorContainer.Clear()
		max := math.Max(float64(cp.selectedIdx), float64(cp.maxColors-1))
		min := math.Max(0, max-float64(cp.maxColors-1))
		if min > 0 {
			cp.paddles[0].SetStatus("enabled")
		} else {
			cp.paddles[0].SetStatus("disabled")
		}
		if int(max) < len(cp.Colors)-1 { // && cp.selectedIdx < len(cp.colors)
			cp.paddles[1].SetStatus("enabled")
		} else {
			cp.paddles[1].SetStatus("disabled")
		}
		fmt.Println(min, max, cp.colSize, cp.maxColors)
		// dots := make([]string, len(cp.colors))
		//   
		//     ﱣ   ﳁ ﭜ   ﳂ  ﱤ     喇    ﴞ           
		numColors := 0
		for i, v := range cp.CoolorColorsPalette.Colors {
			// dots[i] = fmt.Sprintf("[%s:-:-]ﱤ[-:-:-]", v.Html())
			// if i == cp.selectedIdx {
			// 	dots[i] = fmt.Sprintf("[%s:-:b][-:-:-]", v.Html())
			// }
			if i < int(min) || i > int(max) {
				continue
			}
			if numColors >= cp.ColorContainer.GetItemCount() {
				cp.ColorContainer.AddItem(v, cp.colSize, 0, false)
			} else {
				cp.ColorContainer.SetItem(numColors, v, cp.colSize, 0)

			}
			numColors += 1
		}
		// cp.UpdateDots(cp.IconPalette("[%s:-:-]ﱤ[-:-:-]", "[%s:-:b][-:-:-]"))
		cp.UpdateDots(cp.MakeSquarePalette(true))
	})
}

func (cp *CoolorPaletteMainView) SetSelected(idx int) error {
	// dirty := "*"
	// if cp.CoolorColorsPalette.GetMeta().Saved {
	//   dirty = ""
	// } %s
  cc := cp.Colors[cp.selectedIdx]
	status.NewStatusUpdate(
		"name",
    fmt.Sprintf("[red:gray:-][-:-:-][-:gray:-]  [-:-:-][gray:red:-][-:-:-][black:red:b] %s [-:-:-][red:gray:-][-:gray:-]  [gray:pink:-][black::b] %d [pink:gray:-][-:gray:-]  [gray:%s:-][black::b]%s[%s:-:-][-:-:-]", "untitled", cp.Len(), cc.Html(),cc.TVPreview(),cc.Html()/* , cp.CoolorColorsPalette.Name, dirty */),
		// fmt.Sprintf("%s %s", cp.CoolorColorsPalette.Name, dirty),
	)
	cp.UpdateSize()
	cp.ResetViews()
	cp.CoolorColorsPalette.SetSelected(idx)
	// cp.SpawnSelectionEvent(cp.colors[cp.selectedIdx], cp.selectedIdx)
	// MainC.app.QueueUpdateDraw(func() {
	// })
	return nil
}
func (cc *CoolorPaletteMainView) SpawnSelectionEvent(
	c *CoolorColor,
	idx int,
) bool {
	cc.Notify(
		*cc.NewObservableEvent(PaletteColorSelectedEvent, "palette_view_selected", c, cc),
	)
	return true
}

// if len(cc.handlers["selected"]) > 0 {
// 	ev := &ObservableEvent{
// 		// color: c,
// 		// idx:   int8(idx),
// 	}
// 	for _, v := range cc.handlers["selected"] {
// 		if v != nil {
// 			// eh, ok := v.(tcell.EventHandler)
// 			// if !ok {
// 			// 	panic(ok)
// 			// }
// 			v.HandleEvent(ev)
// 		}
// 	}
// }

func (cp *CoolorPaletteMainView) UpdateSize() {
	// cp.colSize = maxColSize
	MainC.app.QueueUpdateDraw(func() {
		if len(cp.Colors) == 0 {
			return
		}
		x, y, w, h := cp.GetInnerRect()
		_, _, _, _ = x, y, w, h
		_, _, w, _ = cp.GetRect()
		cp.maxColors = (w - (PaddleMinWidth * 2)) / cp.colSize
		if len(cp.Colors) < cp.maxColors {
			cp.colSize = (w - (PaddleMinWidth * 2)) / len(cp.Colors)
			if cp.colSize < maxColSize {
				cp.colSize = maxColSize
			}
			cp.maxColors = (w - (PaddleMinWidth * 2)) / cp.colSize
		}
		overflow := (w - (PaddleMinWidth * 2)) % cp.colSize
		left, right := PaddleMinWidth, PaddleMinWidth
		left = overflow / 2
		right = overflow / 2
		if overflow%2 != 0 {
			right += 1
		}
		contW := cp.maxColors * cp.colSize
		cp.Clear()
		cp.AddItem(cp.paddles[0], left, 0, false)
		cp.AddItem(cp.ColorContainer, contW, 0, true)
		cp.AddItem(cp.paddles[1], right, 0, false)
	})
}

func (cp *CoolorPaletteMainView) Draw(screen tcell.Screen) {
	// cp.UpdateSize()
	// num := cp.Flex.GetItemCount()
	// for i := 0; i < num; i++ {
	// 	it := cp.Flex.GetItem(i)
	// 	it.Draw(screen)
	// }
	cp.Flex.Draw(screen)
	// cp.ScrollLine.Draw(screen)
}

// func (cc *CoolorColorsPalette) Register(o Observer) {
// }
// func (cp *CoolorColorsPalette) Register(_ Observer,) {
// 	panic("not implemented") // TODO: Implement
// }
//
// func (cp *CoolorColorsPalette) Deregister(_ Observer) {
// 	panic("not implemented") // TODO: Implement
// }
//
// func (cp *CoolorColorsPalette) Notify(oe ObservableEvent) {
//
// }

// func (cc *CoolorPaletteMainView) AddEventHandler(t string, h tcell.EventHandler) {
// 	// var o Notifier = cc.CoolorColorsPalette
// 	cc.l.Lock()
// 	defer cc.l.Unlock()
// 	if cc.handlers[t] == nil {
// 		cc.handlers[t] = make(EventHandlers, 0)
// 	}
// 	cc.handlers[t] = append(cc.handlers[t], h)
// }

func (cc *CoolorPaletteMainView) GetRef() interface{} {
	return cc
}
func (cp *CoolorPaletteMainView) ToggleLockSelected() (*CoolorColor, int) {
	cc, _ := cp.GetSelected()
	cc.ToggleLocked()
	return cp.GetSelected()
}

func (cp *CoolorPaletteMainView) GetPalette() *CoolorColorsPalette {
	return cp.CoolorColorsPalette
}

func (cp *CoolorMainPalette) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cp.ColorContainer.WrapInputHandler(
		func(event *tcell.EventKey, _ func(p tview.Primitive)) {
			MainC.app.QueueUpdateDraw(func() {
				ch := event.Rune()
				kp := event.Key()
				_ = kp
				switch ch {
				case '*':
					cp.Randomize()
				case '+': // Add a color
					cp.AddRandomCoolorColor()
				case '=':
					cp.GetPalette().Sort()
				case 'd':
					color, _ := cp.GetSelected()
					cp.AddCoolorColor(color.Clone())
					SeentColor("duped", color, color.pallette)
				}
				cp.UpdateSize()
			})
			// if handler := cp.InputHandler(); handler != nil {
			// 	dump.P(fmt.Sprintf("%s input handled", cp.ptype))
			// 	handler(event, setFocus)
			// 	return
			// }
		},
	)
}

// vim: ts=2 sw=2 et ft=go
