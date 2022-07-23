package coolor

import (
	"fmt"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mazznoer/colorgrad"
)

// type CoolorColorRetard struct {
// 	*CoolorColor
// }
//
// func (cc *CoolorColorRetard) GetCC() *CoolorColor {
// 	return cc.Clone()
// }
//
// func NewCoolorColorRetard(c *CoolorColor) *CoolorColorRetard {
// 	ccr := &CoolorColorRetard{
// 		CoolorColor: c.Clone(),
// 	}
// 	return ccr
// }

// func (ccr *CoolorColorRetard) Lighten(amount float64) *CoolorColorRetard {
// 	r, g, b := ccr.color.RGB()
// 	n := noire.NewRGB(float64(r), float64(g), float64(b))
// 	_ = n.Lighten(amount)
// 	cful, _ := colorful.MakeColor(ccr)
// 	_ = cful
// 	h := n.Lighten(0.5).Hex()
// 	return NewCoolorColorRetard(NewCoolorColor(fmt.Sprintf("#%s", h)))
// }

type CoolorColorEdit struct {
	original *CoolorColor
	// *CoolorColorRetard
}

func NewCoolorColorEdit(c *CoolorColor) *CoolorColorEdit {
	cce := &CoolorColorEdit{
		original:          c.Clone(),
		// CoolorColorRetard: NewCoolorColorRetard(c),
	}
	cce.original.updateStyle()
	// cce.updateStyle()
	return cce
}

func (cce *CoolorColorEdit) Reset(c *CoolorColor) {
	cce.original = c.Clone()
  // cce.
	// cce.CoolorColorRetard = NewCoolorColorRetard(c.Clone())
	// cce = NewCoolorColorEdit(c)
	// cce.original.updateStyle()
	// cce.updateStyle()
}

// type Retards []*CoolorColorRetard

type CoolColor interface {
	GetCC() *CoolorColor
}
type Gradiater interface {
	Above() []CoolColor
	Below() []CoolColor
	At(value float64) CoolColor
	Set(value float64)
	Incr(value float64)
	Decr(value float64)
	GetChannelValue(cc CoolColor) float64
	GetCurrentChannelValue() float64
}

type GradStrip struct {
	*tview.Frame
  Flex *tview.Flex
	cces      *CoolorColorEditorStrips
  cm *ColorMod
}

func NewGradStrip(cm *ColorMod, cces *CoolorColorEditorStrips) *GradStrip {
	gs := &GradStrip{
		// Frame:            ,
    Flex: tview.NewFlex(),
		cces:            cces,
    cm: cm,
	}
  // gs.Flex.SetSize(int(cm.size), 1, 0, 0)
  // gs.Frame.Clear
  gs.Flex.SetDirection(tview.FlexRow)
  gs.Flex.SetBorder(false)
  gs.Flex.SetBorderPadding(0,0,0,0)
  gs.Frame = tview.NewFrame(gs.Flex)
  gs.Frame.SetBorderPadding(0,0,0,0)
  // gs.updateStrip()
	return gs
}

func (gs *GradStrip) updateStrip() {
  // gs.Flex.SetSize(int(gs.cm.size*10), 3, 2, 0)
  gs.Flex.SetBorder(false)
  gs.Flex.SetBorderPadding(0,0,0,0)
  // gs.Flex.SetFullScreen(true)
  gs.Flex.Clear()
  above := gs.cm.Above()
  count := 0
  for c, v := range above[0:len(above) - 2] {
    gs.Flex.AddItem(NewStaticCoolorColor(v.Html()).SetBorder(false).SetBorderPadding(0,0,0,0).SetTitleAlign(tview.AlignCenter).SetTitle(fmt.Sprintf("%d", count+c)), 0, 1, false)
    count = c
  }
    gs.Flex.AddItem(NewStaticCoolorColor(gs.cm.orig.Html()).SetBorder(true).SetBorderPadding(0,0,0,0).SetTitleAlign(tview.AlignCenter).SetTitle(fmt.Sprintf("%d", count+1)), 0, 1, false)
    count++
  below := gs.cm.Below()
  for c, v := range below[1:len(below) - 1] {
    // if count + c == count {
    //   continue
    // }
    gs.Flex.AddItem(NewStaticCoolorColor(v.Html()).SetBorder(false).SetBorderPadding(0,0,0,0).SetTitleAlign(tview.AlignCenter).SetTitle(fmt.Sprintf("%d", count+c)), 0, 1, false)
    count = c
  }
  // gs.cces.Flex.SetFullScreen(true)
}

type CoolorColorEditorStrips struct {
	Name string
	*tview.Flex
	mainColor *CoolorColorEdit
  cce *CoolorColorEditor
	increment       float64
  colorMods     map[string]*ColorMod
  gs *GradStrip
  gsat *GradStrip
}

func (ccep *CoolorColorEditorStrips) updateState() {
  ccep.gs.updateStrip()
  ccep.gsat.updateStrip()
  ccep.Flex.Clear()
  ccep.Flex.AddItem(tview.NewBox(), 0, 1, false)
  ccep.Flex.AddItem(ccep.gs, 0, 10, false)
  ccep.Flex.AddItem(tview.NewBox(), 0, 1, false)
  ccep.Flex.AddItem(ccep.gsat, 0, 10, false)
  ccep.Flex.AddItem(tview.NewBox(), 0, 1, false)
}

func (ccep *CoolorColorEditorStrips) UpdateColor(c *CoolorColor) {
  ccep.gs.cm.SetColor(c.GetCC())
  ccep.gsat.cm.SetColor(c.GetCC())
  ccep.updateState()
}

func (ccep *GradStrip) Lerp(startc, endc CoolColor) *colorgrad.GradientBuilder {
	start, ok := colorful.MakeColor(startc.GetCC())
	end, ok := colorful.MakeColor(endc.GetCC())
	_ = ok
	cg := colorgrad.NewGradient().Colors(
		start, end,
	)
	return cg
}

func NewEditorStrip(name string, cce *CoolorColorEditor) *CoolorColorEditorStrips {
	ccep := &CoolorColorEditorStrips{
    Name: name,
    cce: cce,
	}
  ccep.Flex = tview.NewFlex() // .SetFullScreen(true)
  ccep.Flex.SetDirection(tview.FlexColumn)
  ccep.gs = NewGradStrip(HueMod, ccep)
  ccep.gsat = NewGradStrip(SatMod, ccep)
  // ccep.gs.updateStrip()
  // ccep.Flex.AddItem(tview.NewBox(), 0, 1, false)
  // ccep.Flex.AddItem(ccep.gs, 0, 10, true)
  // ccep.Flex.AddItem(tview.NewBox(), 0, 1, false)
  // ccep.Flex.AddItem(ccep.gsat, 0, 10, false)
  // ccep.Flex.AddItem(tview.NewBox(), 0, 1, false)
	return ccep
}

type CoolorColorEditor struct {
	*tview.Flex
	palette  *CoolorPalette
	previews *CoolorColorEditorStrips
	app *tview.Application
}

func NewCoolorEditor(a *tview.Application, cp *CoolorPalette) *CoolorColorEditor {
	cce := &CoolorColorEditor{}
	cce.palette = cp
	cce.app = a
  cce.previews = NewEditorStrip("Hue", cce)
  cce.previews.colorMods = make(map[string]*ColorMod)
  cce.previews.colorMods["hue"] = HueMod
  // cce.previews.hue = NewGradStrip(HueMod)
  cce.Init()
  // runtime.Breakpoint()
	return cce
}

func (cce *CoolorColorEditor) Init() {
	cce.Flex = tview.NewFlex().SetDirection(tview.FlexColumn)
  cce.Flex.SetBorder(false).SetBorderPadding(0,0,0,0)
	cce.palette.AddEventHandler("selected", cce)
	cce.Flex.AddItem(tview.NewBox(), 0, 1, false)
	cce.Flex.AddItem(cce.previews, 0, 5, false)
	cce.Flex.AddItem(tview.NewBox(), 0, 1, false)
}

func (cce *CoolorColorEditor) HandleEvent(e tcell.Event) bool {
	switch e.(type) {
	case *SelectionEvent:
		se := e.(*SelectionEvent)
    // cce.previews.gs.cm.SetColor(se.color)
    cce.previews.UpdateColor(se.color)
    // cce.previews
		// cce.cce.Reset(se.color.Clone())
		return true
	}
	return false
}

// vim: ts=2 sw=2 et ft=go
