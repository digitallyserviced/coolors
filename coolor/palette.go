package coolor

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"

	"github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/samber/lo"
)

var PaddleMinWidth int = 4

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
		l *sync.RWMutex
		*eventObserver
		*eventNotifier
		Name        string
		ptype       string
		Colors      CoolorColors
		Hash        uint64 `boltholdKey:"Hash"`
		selectedIdx int
	}
	CoolorPaletteMainView struct {
		ColorContainer *tview.Flex
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
	left := NewPalettePaddle("", "")
	right := NewPalettePaddle("", "")
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
	// ccp.eventObserver = NewEventObserver("palette")
	// ccp.eventNotifier = NewEventNotifier("palette")
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
		eventObserver: NewEventObserver("palette"),
		eventNotifier: NewEventNotifier("palette"),
		Name:          Generator().GenerateName(2),
		ptype:         "",
		Colors:        make(CoolorColors, 0),
		Hash:          0,
		selectedIdx:   0,
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
	cp.AddItem(cp.paddles[0], 4, 0, false)
	cp.AddItem(cp.ColorContainer, 80, 0, true)
	cp.AddItem(cp.paddles[1], 4, 0, false)
	return cp
}

func BlankCoolorPalette() *CoolorPaletteMainView {
	cp := &CoolorPaletteMainView{
		CoolorColorsPalette: NewCoolorColorsPalette(),
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
	cp.AddItem(cp.paddles[0], 4, 0, false)
	cp.AddItem(cp.ColorContainer, 80, 0, true)
	cp.AddItem(cp.paddles[1], 4, 0, false)
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
	cp.AddItem(cp.paddles[0], 4, 0, false)
	cp.AddItem(cp.ColorContainer, 80, 0, true)
	cp.AddItem(cp.paddles[1], 4, 0, false)
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
	cp.AddItem(cp.paddles[0], 4, 0, false)
	cp.AddItem(cp.ColorContainer, 80, 0, true)
	cp.AddItem(cp.paddles[1], 4, 0, false)
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

func (ccs *CoolorColorsPalette) Swap(a, b int) {
	if col := ccs.Colors[a]; col == nil {
		return
	}
	if col := ccs.Colors[b]; col == nil {
		return
	}
	if ccs.selectedIdx == a {
		ccs.Colors[b].SetSelected(true)
		ccs.Colors[a].SetSelected(false)
	} else {
		ccs.Colors[a].SetSelected(true)
		ccs.Colors[b].SetSelected(false)
	}
	ccs.Colors[a], ccs.Colors[b] = ccs.Colors[b], ccs.Colors[a]
}

func (ccs *CoolorColorsPalette) Less(a, b int) bool {
	return ccs.Colors[a].Color.Hex() < ccs.Colors[b].Color.Hex()
}

func (ccs *CoolorColorsPalette) Len() int {
	return len(ccs.Colors)
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
	cbp.Colors = make(CoolorColors, 0)
	base, _ := MakeColor(cbp.base)
	done := make(chan struct{})
	// cbp.colors = make(CoolorColors, 0)
	defer close(done)
	colors := RandomShadesStream(base, 0.2)
	colors.Status.SetProgressHandler(NewProgressHandler(func(u uint32) {
		status.NewStatusUpdate(
			"action_str",
			fmt.Sprintf("Found Shades (%d / %d)", u, colors.Status.GetItr()),
		)
	}, func(i uint32) {
		status.NewStatusUpdate(
			"action_str",
			fmt.Sprintf("Found Shades (%d / %d)", colors.Status.GetValid(), i),
		)
	}))
	colors.Run(done)
	for _, v := range TakeNColors(done, colors.OutColors, int(cbp.increments)) {
		newcc := NewStaticCoolorColor(v.Hex())
		cbp.AddCoolorColor(newcc)
		// SeentColor("stream_random_shade", newcc, newcc.pallette)
	}
	// cbp.UpdateSize()
 //  cbp.ResetViews()
	// cbp.SetSelected(0)
}

func (cbp *CoolorBlendPalette) Init() {
	cbp.ColorContainer.Clear()
	cbp.Colors = make(CoolorColors, 0)
	incrSizes := 1.0 / cbp.increments
	start, _ := MakeColor(cbp.start)
	end, _ := MakeColor(cbp.end)
	for i := 0; i <= int(cbp.increments); i++ {
		newc := start.BlendLab(end, float64(i)*float64(incrSizes))
		newcc := NewStaticCoolorColor(newc.Hex())
		cbp.AddCoolorColor(newcc)
		// SeentColor("mixed_colors_gradient", newcc, newcc.pallette)
	}

	MainC.conf.AddPalette("blend", cbp)
}

func K(from string, cc *CoolorColor, src Referenced) {
	// fmt.Println(from, cc)
	if MainC == nil || MainC.eventNotifier == nil {
		return
	}
	MainC.eventNotifier.Notify(
		*MainC.eventNotifier.NewObservableEvent(ColorSeentEvent, from, cc, src),
	)
}

func (cp *CoolorColorsPalette) PaletteEvent(
	t ObservableEventType,
	color *CoolorColor,
) {
	cp.Notify(*MainC.NewObservableEvent(t, "palette_event", color, cp))
	MainC.eventNotifier.Notify(
		*MainC.NewObservableEvent(t, "palette_event", color, cp),
	)
}

func (cp *CoolorColorsPalette) SetColors(
	cols []tcell.Color,
) *CoolorColorsPalette {
	cp.Colors = make(CoolorColors, 0)
	for _, v := range cols {
		cp.AddCoolorColor(NewIntCoolorColor(v.Hex()))
		// cp.Colors = append(cp.Colors, )
	}

	return cp
}

func (cp *CoolorColorsPalette) AddCoolorColor(color *CoolorColor) *CoolorColor {
	cp.l.Lock()
	cp.Colors = append(cp.Colors, color)
	cp.l.Unlock()
	cp.PaletteEvent(PaletteColorModifiedEvent, color)
	return cp.Colors[len(cp.Colors)-1]
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

func (cp *CoolorColorsPalette) GetItem(idx uint) *CoolorColor {
	id := int(math.Mod(float64(idx), float64(len(cp.Colors))))
	// fmt.Println(id, float64(len(cp.Colors)))
	return cp.Colors[id%len(cp.Colors)]
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

func (cp *CoolorColorsPalette) RandomColor() *CoolorColor {
	return cp.GetItem(uint(rand.Uint32()))
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

func (cp *CoolorColorsPalette) GetSelected() (*CoolorColor, int) {
	cp.l.RLock()
	defer cp.l.RUnlock()

	if len(cp.Colors) == 0 {
		return nil, -1
	}

	if cp.selectedIdx > len(cp.Colors)-1 {
		cp.selectedIdx = 0
	}
	if cp.selectedIdx < 0 {
		cp.selectedIdx = len(cp.Colors) - 1
	}

	if cp.Colors[cp.selectedIdx] != nil {
		return cp.Colors[cp.selectedIdx], cp.selectedIdx
	}
	return nil, -1
}

func (cp *CoolorColorsPalette) String() string {
	return strings.Join(
		lo.Map[*CoolorColor, string](
			cp.Colors,
			func(cc *CoolorColor, i int) string {
				return cc.TerminalPreview()
			},
		),
		" ",
	)
}

func (cp *CoolorColorsPalette) Sort() {
	sort.Sort(cp)
}

func (cp *CoolorPaletteMainView) Sort() {
	MainC.app.QueueUpdateDraw(func() {
		cp.CoolorColorsPalette.Sort()
    cp.UpdateSize()
    cp.ResetViews()
	})
}

func (cp *CoolorPaletteMainView) UpdateDots(dots []string) {
	status.NewStatusUpdate("dots", strings.Join(dots, " "))
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
    fmt.Println(min,max,cp.colSize,cp.maxColors)
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
      if numColors >= cp.ColorContainer.GetItemCount()  {
        cp.ColorContainer.AddItem(v, cp.colSize, 0, false)
      } else {
        cp.ColorContainer.SetItem(numColors, v, cp.colSize, 0)
        
      }
      numColors+=1
		}
		// cp.UpdateDots(cp.IconPalette("[%s:-:-]ﱤ[-:-:-]", "[%s:-:b][-:-:-]"))
		cp.UpdateDots(cp.MakeSquarePalette(true))
	})
}

func (cp *CoolorColorsPalette) MakeSquarePalette(showSelected bool) []string {
	// main, sel := "[%s:-:-]ﱢ[-:-:-]", "[%s:-:b]ﱡ[-:-:-]"
	main, sel := "[%s:-:b]ﱡ[-:-:-]", "[%s:-:-]ﱢ[-:-:-]" 
	// main, sel := "[%s:-:-]▉▉[-:-:-]", "[%s:-:b]▉▉[-:-:-]"
	if !showSelected {
		main = sel
	}
	return cp.IconPalette(main, sel)
}

func (cp *CoolorColorsPalette) MakeDotPalette() []string {
	return cp.IconPalette("[%s:-:-]ﱤ[-:-:-]", "[%s:-:b][-:-:-]")
}

func (cp *CoolorColorsPalette) IconPalette(
	mainFormat, selectedFormat string,
) []string {
	chars := make([]string, len(cp.Colors))

	for i, v := range cp.Colors {
		f := mainFormat
		if i == cp.selectedIdx {
			f = selectedFormat
		}
		chars[i] = fmt.Sprintf(f, v.Html())
	}

	return chars
}

func (cp *CoolorColorsPalette) TagsKeys(random bool) CoolorPaletteTagsMeta {
	tagKeys := make(map[string]*Coolor)
	cptm := &CoolorPaletteTagsMeta{
		TaggedColors: tagKeys,
	}
	tlist := GetTerminalColorsAnsiTags()
	if Base16Tags.tagList == nil {
	}

	for _, v := range tlist.items {
		// k := v.data[keyfield.name].(string)
		k := v.GetKey()
		tagKeys[k] = cp.RandomColor().Coolor()
	}

	cp.Each(func(cc *CoolorColor, i int) {
		tags := cc.GetTags()
		for _, t := range tags {
			k := t.GetKey()
			cptm.TaggedColors[k] = cc.Coolor()
		}
	})

	return *cptm
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

func (cc *CoolorColorsPalette) SpawnSelectionEvent(
	c *CoolorColor,
	idx int,
) bool {
	cc.Notify(
		*cc.NewObservableEvent(PaletteColorSelectedEvent, "palette_selected", c, cc),
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
		// w -= 2
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
		left += overflow / 2
		right += overflow / 2
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
func (cc *CoolorColor) GetRef() interface{} {
	return cc
}

func (cp *CoolorColorsPalette) GetRef() interface{} {
	return cp
}

func (cp *CoolorColorsPalette) ClearSelected() {
	cp.Each(func(cc *CoolorColor, _ int) {
		cc.SetSelected(false)
	})
}

func (cp *CoolorColorsPalette) SetSelected(idx int) error {
	if len(cp.Colors) == 0 {
		return fmt.Errorf("no valid color at idx: %d", idx)
	}
	// MainC.app.QueueUpdateDraw(func() {
	if idx < 0 {
		idx = len(cp.Colors) - 1
	}
	if idx > len(cp.Colors)-1 {
		idx = 0
	}
	if idx < len(cp.Colors) {
		cp.ClearSelected()
		cp.selectedIdx = idx
		cp.Colors[cp.selectedIdx].SetSelected(true)
		cp.SpawnSelectionEvent(cp.Colors[cp.selectedIdx], cp.selectedIdx)
		// cp.SpawnSelectionEvent(cp.colors[cp.selectedIdx], cp.selectedIdx)
	}
	// })
	return nil
}

func (cp *CoolorPaletteMainView) SetSelected(idx int) error {
	dirty := "*"
	// if cp.CoolorColorsPalette.GetMeta().Saved {
	//   dirty = ""
	// }
	status.NewStatusUpdate(
		"name",
		fmt.Sprintf("%s %s", cp.CoolorColorsPalette.Name, dirty),
	)
	cp.UpdateSize()
	cp.ResetViews()
	cp.CoolorColorsPalette.SetSelected(idx)
	// cp.SpawnSelectionEvent(cp.colors[cp.selectedIdx], cp.selectedIdx)
	// MainC.app.QueueUpdateDraw(func() {
	// })
	return nil
}

func (cp *CoolorColorsPalette) Coolors() *Coolors {
	// colors :=
	css := &Coolors{
		Key:    cp.Name,
		Colors: make([]*Coolor, 0),
		Saved:  false,
	}
	for _, v := range cp.Colors {
		css.Colors = append(css.Colors, v.Coolor())
	}
	return css
}

func (cp *CoolorColorsPalette) Each(f func(*CoolorColor, int)) {
	// MainC.app.QueueUpdate(func() {
	for i, v := range cp.Colors {
		f(v, i)
	}
	// })
}

func (cp *CoolorColorsPalette) Plainify(s bool) {
	cp.Each(func(cc *CoolorColor, _ int) {
		cc.SetStatic(s)
		cc.SetPlain(s)
	})
}

func (cp *CoolorColorsPalette) Staticize(s bool) {
	cp.Each(func(cc *CoolorColor, _ int) {
		cc.SetStatic(s)
	})
}

func (cp *CoolorPaletteMainView) ToggleLockSelected() (*CoolorColor, int) {
	cc, _ := cp.GetSelected()
	cc.ToggleLocked()
	return cp.GetSelected()
}

func (cp *CoolorColorsPalette) GetItemCount() int {
	return cp.Len()
}

func (cp *CoolorColorsPalette) GetPalette() *CoolorColorsPalette {
	return cp
}

func (cp *CoolorPaletteMainView) GetPalette() *CoolorColorsPalette {
	return cp.CoolorColorsPalette
}

func (cbp *CoolorBlendPalette) GetPalette() *CoolorColorsPalette {
	return cbp.CoolorColorsPalette
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
