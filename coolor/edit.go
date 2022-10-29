package coolor

import (
	"fmt"
	"math"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mazznoer/colorgrad"

	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/util"
	"github.com/digitallyserviced/coolors/status"
)

const (
	SelectedSize          int = 4
	NearestNeighbor       int = 2
	SecondNearestNeighbor int = 1
	LastNearestNeighbors  int = 1
)
var (
	ColorModNames = []string{"Hue", "Chroma", "Light"}
	ColorMods     = map[string]*ColorMod{
		"Hue":    HueMod,
		"Chroma": SatMod,
		"Light":  LightMod,
	}
)

type CoolorColorEditor struct {
	*tview.Flex
	palette  *CoolorColorsPalette
	previews *CoolorColorEditorStrips
	app      *tview.Application
	updateCh chan<- *status.StatusUpdate
	settings EditorSettings
}
type fibWin struct {
	prev, curr, next float64
}


type EditorSettings struct {
	incrementSize float64
	sizeNum       float64
	wrap          bool
	win           fibWin
	selectedMod   int
}

type CoolorColorEditorStrips struct {
	*tview.Flex
	mainColor *CoolorColor
	cce       *CoolorColorEditor
	colorMods map[string]*GradStrip
	increment float64
}


type CoolColor interface {
	GetCC() *CoolorColor
	RGBA() (r, g, b, a uint32)
}
type GradStrip struct {
	*tview.Flex
	cces           *CoolorColorEditorStrips
	footer         *tview.Box
	container      *tview.Flex
	gauge          *tview.Flex
	strip          *tview.Flex
	frame          *tview.Frame
	cm             *ColorMod
	gradOffsets    []float64
	validGrad      []bool
gradVal []float64
	centeredColors []string
	centeredGrad   []*CoolorColor
	sizes          []int
	previewGrad    []string
	height         int
	selected       bool
}


type CoolorColorEdit struct {
	original *CoolorColor
}
func NewEditorStrip(cce *CoolorColorEditor) *CoolorColorEditorStrips {
	ccep := &CoolorColorEditorStrips{
		cce: cce,
	}
	ccep.Flex = tview.NewFlex() // .SetFullScreen(true)
	ccep.SetDirection(tview.FlexColumn)
	ccep.Clear()
	ccep.colorMods = make(map[string]*GradStrip, len(ColorMods))
	col, _ := ccep.cce.palette.GetSelected()
	// ccep.mainColor.original = col
	BlankSpace(ccep.Flex)
	for i, v := range ColorMods {
		ccep.colorMods[i] = NewGradStrip(v, ccep)
	}
	ccep.forColorMods(func(c *GradStrip, n string) {
		ccep.AddItem(c, 0, 8, true)
	})
	BlankSpace(ccep.Flex)
	ccep.UpdateColor(col)
	ccep.updateState()
	return ccep
}
func NewEditorSetting(cce *CoolorColorEditor) {
	cce.settings = EditorSettings{
		incrementSize: 0,
		sizeNum:       20,
		wrap:          true,
		win:           getFibNum(1, 0),
		selectedMod:   0,
	}
	cce.updateState()
}


func NewCoolorEditor(a *tview.Application, cp Palette) *CoolorColorEditor {
	cce := &CoolorColorEditor{}
	NewEditorSetting(cce)
	cce.palette = cp.GetPalette()
	// cce.updateCh = updateCh
	cce.app = a
	cce.previews = NewEditorStrip(cce)
	cce.Init()
	return cce
}


func NewCoolorColorEdit(c *CoolorColor) *CoolorColorEdit {
	cce := &CoolorColorEdit{
		original: c.Clone(),
	}
	cce.original.updateStyle()
	return cce
}


// type Retards []*CoolorColorRetard
func NewGradStrip(cm *ColorMod, cces *CoolorColorEditorStrips) *GradStrip {
	gs := &GradStrip{
		Flex:           tview.NewFlex(),
		height:         24,
		selected:       false,
		sizes:          []int{},
		centeredGrad:   []*CoolorColor{},
		centeredColors: []string{},
		gradOffsets:    []float64{},
		previewGrad:    []string{},
		// frame:          &tview.Frame{},
		strip:          tview.NewFlex(),
		gauge:          tview.NewFlex(),
		container:      tview.NewFlex(),
		footer:         nil,
		cces:           cces,
		cm:             cm,
	}
	gs.SetDirection(tview.FlexRow)
	gs.container.SetDirection(tview.FlexColumn)
	gs.gauge.SetDirection(tview.FlexRow)
	gs.strip.SetDirection(tview.FlexRow)
	gs.strip.SetBorder(false)
	gs.strip.SetBorderPadding(0, 0, 0, 0)
	gs.strip.Clear()
	gs.frame = tview.NewFrame(gs.container)
  gs.frame.SetBorders(0, 0, 1, 0, 0, 0)
  gs.frame.SetBorderPadding(0, 0, 0, 0)
	gs.AddItem(gs.frame, 0, 10, false)
	gs.container.AddItem(gs.gauge, 0, 1, false)
	gs.container.AddItem(gs.strip, 0, 10, false)
	// spc, settxt := MakeCenterLineSpacer(gs.Flex)
	// gs.footer = spc
	MakeSpacer(gs.strip)
	// MakeSpacer(gs.gauge)
	BlankSpace(gs.gauge)
	// gs.Frame = tview.NewFrame(gs.strip)
	//  gs.Frame.SetBorder(false)
	// gs.Frame.SetBorderPadding(2,2, 0, 0)
	// gs.Frame.AddText(gs.cm.name, true, tview.AlignCenter, getFGColor(*gs.cm.current.color))
	gs.updateStrip()
	return gs
}

func (gs *GradStrip) NavSelection(idx int) {
	if idx > 0 {
		gs.cm.Incr(gs.cm.increment)
	}
	if idx < 0 {
		gs.cm.Decr(gs.cm.increment)
	}
}

func (cce *CoolorColorEdit) Reset(c *CoolorColor) {
	cce.original = c.Clone()
	// cce.updateStyle()
}
func (gs *GradStrip) MakeSelectionGradient() {
	centerIdx := 0
	centerColor := NewStaticCenteredCoolorColor(gs.cm.current.GetColor())
	gs.centeredGrad = make([]*CoolorColor, len(gs.sizes))
	gs.centeredColors = make([]string, len(gs.sizes))
	gs.gradOffsets = make([]float64, len(gs.sizes))
	gs.previewGrad = make([]string, len(gs.sizes))
	gs.validGrad = make([]bool, len(gs.sizes))
	gs.gradVal = make([]float64, len(gs.sizes))
	mirrorIdx := 0.0
	for i, v := range gs.sizes {
		if v == SelectedSize {
			centerColor := gs.cm.current.Clone()
      gs.gradVal[i] = gs.cm.ChannelMod.Get(centerColor)
			gs.centeredGrad[i] = centerColor
			gs.centeredColors[i] = centerColor.GetColor()
			gs.previewGrad[i] = centerColor.GetCC().TerminalPreview()
			centerIdx = i
		}
		if centerIdx > 1 {
			mirrorIdx = math.Mod(float64(i), float64(centerIdx)+1)
			leftIdx := centerIdx + int(mirrorIdx)
			rightIdx := centerIdx - int(mirrorIdx)
			ldiff := -math.Abs(gs.cm.increment * (mirrorIdx))
			lCol, v := gs.cm.ChannelMod.Mod(centerColor.Clone(), ldiff)
      gs.gradVal[leftIdx] = gs.cm.ChannelMod.Get(lCol)
			gs.validGrad[leftIdx] = v
			gs.gradOffsets[leftIdx] = ldiff
			gs.centeredGrad[leftIdx] = lCol.GetCC()
			gs.centeredColors[leftIdx] = lCol.GetCC().GetColor()
			gs.previewGrad[leftIdx] = lCol.GetCC().TerminalPreview()
			rdiff := math.Abs(gs.cm.increment * (mirrorIdx))
			rCol, v := gs.cm.ChannelMod.Mod(centerColor.Clone(), rdiff)
      gs.gradVal[rightIdx] = gs.cm.ChannelMod.Get(rCol)
			gs.validGrad[rightIdx] = v
			gs.gradOffsets[rightIdx] = rdiff
			gs.centeredGrad[rightIdx] = rCol.GetCC()
			gs.centeredColors[rightIdx] = rCol.GetCC().GetColor()
			gs.previewGrad[rightIdx] = rCol.GetCC().TerminalPreview()
			if leftIdx < 0 || rightIdx > len(gs.sizes) {
				break
			}
		}
	}
}

func (gs *GradStrip) MakeSelections() {
	gs.MakeSelectionGradient()
	gs.gauge.Clear()
	gc := gs.strip.GetItemCount()
	for c, v := range gs.centeredColors {
		size := gs.sizes[c]
		if gc != len(gs.centeredColors) {
      sign := ""
      if gs.gradOffsets[c] < 0 {
        sign = ""
      }
			if size == SelectedSize {
				gs.centeredGrad[c].SetInfoLine(fmt.Sprintf(" %s% 6.2f ", sign, gs.gradVal[c]), gs.validGrad[c])
			} else {
				gs.centeredGrad[c].SetInfoLine(fmt.Sprintf(" %s% 6.2f (% 6.2f)", sign, math.Abs(gs.gradOffsets[c]), gs.gradVal[c]), gs.validGrad[c])
			}
			gs.strip.AddItem(gs.centeredGrad[c], size, 0, true)
			continue
		}
		if fItem := gs.strip.GetItem(c); fItem != nil {
			cc, ok := fItem.(*CoolorColor)
			if ok {
				cc.SetColorCss(v)

				// cc.SetInfoLine(fmt.Sprintf("%0.2f", gs.gradOffsets[c]), gs.validGrad[c])
				gs.strip.ResizeItem(fItem, size, 0)
			}
			continue
		}
	}
}

func (gs *GradStrip) updateIncrement(incr float64) {
	gs.cm.increment = incr * gs.cm.scale
	gs.updateStrip()
}

func (gs *GradStrip) SetStatus() {
	gs.frame.Clear()
	// len(ColorModNames) 🮋🮒🮑🮐🮆🮔🮕🮖🮗🮟🮞🮝🮜🮘🮙🮚🮱🮴🮽🮾🮿🯄
	if gs.selected {
		gs.frame.AddText(fmt.Sprintf("[blue:black:-][black:blue:-] %s [blue:black:-][-:-:-]", gs.cm.name), true, tview.AlignCenter, tcell.ColorWhite)
	} else {
		gs.frame.AddText(fmt.Sprintf(" %s ", gs.cm.name), true, tview.AlignCenter, tcell.ColorWhite)
	}
	// status, last := gs.cm.ChannelMod.GetStatus(&gs.cm.current)
	status, last := gs.cm.GetStatus()
	gs.frame.AddText(fmt.Sprintf("  %s %s", status, last), false, tview.AlignCenter, tcell.ColorWhite)
}

func (gs *GradStrip) updateStrip() {
	// MainC.app.QueueUpdateDraw(func() {
	gs.SetStatus()
	sizes := []int{SecondNearestNeighbor, NearestNeighbor, SelectedSize, NearestNeighbor, SecondNearestNeighbor}
	// gs.strip.Clear()
	// MakeSpace(gs.strip, "", "#000000", 0, 100)
	// MakeSpace(gs.strip, "", "#000000", 0, 100)
	x, y, w, h := gs.container.GetRect()
	_, _, _, _ = x, y, w, h
	gs.strip.Clear()

	mainReserve := 0
	for _, v := range sizes {
		mainReserve = mainReserve + v
	}
	// h := gs.height
	leftovers := h - mainReserve
	wrapNum := leftovers / 2
	outsize := len(sizes) + (wrapNum * 2)
	out := make([]int, int(math.Abs(float64(outsize))))
	for i := 0; i < len(out); i++ {
		if i >= wrapNum && i < len(out)-wrapNum {
			out[i] = sizes[i-wrapNum]
		} else {
			out[i] = LastNearestNeighbors
		}
	}
	// dump.P(out, len(out), leftovers, wrapNum)
	gs.sizes = out
	gs.MakeSelections()
	// })
}

func (ccep *CoolorColorEditorStrips) updateState() {
	ccep.forColorMods(func(c *GradStrip, n string) {
		c.updateStrip()
	})
}

func (ccep *CoolorColorEditorStrips) SaveColor() {
	if ccep.mainColor != nil {
		// var cc *CoolorColor
		// ccep.forColorMods(func(c *GradStrip, n string) {
		//    cc = &c.cm.current
		// })
		// ccep.mainColor.SetColor(cc.color)
		col, _ := ccep.cce.palette.GetSelected()
		col.SetColor(ccep.mainColor.Color)
	}
}

func (ccep *CoolorColorEditorStrips) UpdateColor(cc *CoolorColor) {
	// MainC.app.QueueUpdateDraw(func() {
	if ccep == nil {
		return
	}
	ccep.mainColor = cc
	ccep.forColorMods(func(c *GradStrip, n string) {
		c.cm.SetColor(cc.GetCC())
	})
	status.NewStatusUpdate("color", ccep.cce.String())
	// })
	ccep.updateState()
}

func (ccep *GradStrip) Lerp(startc, endc CoolColor) *colorgrad.GradientBuilder {
	start, _ := colorful.MakeColor(startc.GetCC())
	end, ok := colorful.MakeColor(endc.GetCC())
	_ = ok
	cg := colorgrad.NewGradient().Colors(
		start, end,
	)
	return cg
}


func (ccep *CoolorColorEditorStrips) forColorMods(f func(c *GradStrip, n string)) {
	if ccep == nil || ccep.colorMods == nil || len(ccep.colorMods) == 0 {
		return
	}
	for n, cm := range ccep.colorMods {
		if cm == nil {
			continue
		}
		f(cm, n)
	}
}

func (ccep *CoolorColorEditorStrips) updateColorMods() {
	ccep.forColorMods(func(c *GradStrip, n string) {
		c.SetBorder(false)
		c.SetBorderPadding(0, 0, 0, 0)
		c.Clear()
		MakeSpace(c.Flex, "", "#000000", 0, 1)
		x, y, w, h := c.GetInnerRect()
		_, _, _ = x, y, w
		c.height = h
		c.cm.SetSize(float64(h))
		c.updateStrip()
	})
}
func (es EditorSettings) GetSelected() int {
	return es.selectedMod
}

func (es *EditorSettings) GetSelectedMod() string {
	return ColorModNames[int(math.Mod(float64(es.selectedMod), float64(len(ColorModNames))))]
}
func getFibNum(base, num float64) fibWin {
	start := base
	prev := start
	curr := prev + start
	next := prev + curr
	tprev, tcurr, tnext := prev, curr, next
	for i := 0.0; i < num; i++ {
		tprev, tcurr, tnext = (curr), (prev + curr), curr+next
		prev, curr, next = tprev, tcurr, tnext
	}
	return fibWin{prev, curr, next}
}

func (cce *CoolorColorEditor) CurrentFibValue(base float64) float64 {
	win := getFibNum(base, cce.settings.win.curr)
	return win.curr
}

func (cce *CoolorColorEditor) updateState() {
	cce.settings.win = getFibNum(1, cce.settings.sizeNum)
	status.NewStatusUpdate("action", fmt.Sprintf("%s %s", "", "edit"))
	cce.previews.forColorMods(func(c *GradStrip, n string) {
		// currIncr := (c.cm.minIncrValue * cce.CurrentFibValue(1))*c.cm.scale
		currIncr := cce.settings.sizeNum*c.cm.minIncrValue + (c.cm.minIncrValue * cce.settings.sizeNum)
		c.updateIncrement(currIncr)
		if cce != nil {
			status.NewStatusUpdate("color", cce.String())
		}
	})
	// MainC.app.Draw()
}

func (cce *CoolorColorEditor) String() string {
	if cce == nil || cce.previews == nil || cce.previews.mainColor == nil {
		return ""
	}
	// col, _ := cce.palette.GetSelected()
	// gs.Frame.AddText(gs.cm.String(), true, tview.AlignCenter, getFGColor(*gs.cm.current.color))
	col := cce.previews.mainColor
	is := ""
	// cce.previews.forColoMods(func(c *GradStrip, n string) {
	// 	is = fmt.Sprintf("%s %s ", is, c.cm.String())
	// })
	is = fmt.Sprintf("%s %s", col.TVPreview(), is)
	return tview.TranslateANSI(is)
	// return tview.TranslateANSI(is)
}

func (cce *CoolorColorEditor) KeyEvent(ev *tcell.EventKey) {
}

func (cce *CoolorColorEditor) Init() {
	cce.Flex = tview.NewFlex().SetDirection(tview.FlexRow)
	cce.Flex.SetBorder(false).SetBorderPadding(0, 0, 0, 0)
  cce.palette.Register(PaletteColorSelectedEvent,cce)
	cce.Clear()
	cce.AddItem(cce.previews, 0, 1, false)
	cce.updateState()
}

func (cce *CoolorColorEditor) SetSelected(idx int) {
	cce.settings.selectedMod = int(util.Clamp(float64(idx), 0, float64(len(ColorModNames))))
}

// NavSelection(idx int) error
func (cce *CoolorColorEditor) NavSelection(idx int) {
	idx = cce.settings.selectedMod + idx
	if idx < 0 {
		idx = len(ColorModNames) - 1
	}
	if idx > len(ColorModNames)-1 {
		idx = 0
	}
	cce.SetSelected(idx)
}

func (cce *CoolorColorEditor) GetSelectedVimNav() VimNav {
	dump.P(cce.settings.GetSelectedMod())
	id := cce.settings.GetSelectedMod()

	selCm := cce.previews.colorMods[id]
	selCm.selected = true
	return selCm
}

func (cce *CoolorColorEditor) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cce.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		ch := event.Rune()
		kp := event.Key()
		if kp == tcell.KeyEnter || kp == tcell.KeyEscape {
			if kp == tcell.KeyEnter {
				cce.previews.SaveColor()
				MainC.pages.SwitchToPage("palette")
			}
			if kp == tcell.KeyEscape {
				MainC.pages.SwitchToPage("palette")
			}
			AppModel.helpbar.SetTable("palette")
			return
		}
		MainC.app.QueueUpdateDraw(func() {
			cce.previews.forColorMods(func(c *GradStrip, n string) {
				c.cm.last = nil
				c.selected = false
			})
			HandleVimNavigableHorizontal(cce, ch, kp)
			vn := HandleVimNavSelectable(cce)
			HandleVimNavigableVertical(vn, ch, kp)
			cce.previews.forColorMods(func(c *GradStrip, n string) {
				cce.previews.UpdateColor(&c.cm.current)
				c.SetStatus()
			})

			if ch == '<' {
				cce.settings.sizeNum = util.Clamp(cce.settings.sizeNum-1, 1, 30)
			}
			if ch == '>' {
				cce.settings.sizeNum = util.Clamp(cce.settings.sizeNum+1, 1, 30)
			}
			cce.updateState()
		})
	})
}

func (cce *CoolorColorEditor) Name() string {
return fmt.Sprintf("%s", "editor")
}
func (cce *CoolorColorEditor) HandleEvent(oe ObservableEvent) bool {
  switch oe.Type {
  case PaletteColorSelectedEvent:
    var cc *CoolorColor = oe.Ref.(*CoolorColor)
    cce.previews.UpdateColor(cc)
  }
	// switch e := e.(type) {
	// case *ObservableEvent:
	// 	// se := e
	// 	// cce.previews.gs.cm.SetColor(se.color)
 //    // var color *CoolorColor = se.Ref
	// 	// cce.previews.UpdateColor()
	// 	// cce.previews
	// 	// cce.cce.Reset(se.color.Clone())
	// 	return true
	// }
	return true
}
// func (cce *CoolorColorEditor) HandleEvent(e tcell.Event) bool {
// 	// switch e := e.(type) {
// 	// case *ObservableEvent:
// 	// 	// se := e
// 	// 	// cce.previews.gs.cm.SetColor(se.color)
//  //    // var color *CoolorColor = se.Ref
// 	// 	// cce.previews.UpdateColor()
// 	// 	// cce.previews
// 	// 	// cce.cce.Reset(se.color.Clone())
// 	// 	return true
// 	// }
// 	return false
// }

// vim: ts=2 sw=2 et ft=go
