package coolor

import (
	"fmt"
	"math"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mazznoer/colorgrad"
	"github.com/teacat/noire"
)

type CoolorColorRetard struct {
	*CoolorColor
}

func (cc *CoolorColorRetard) GetCC() *CoolorColor {
	return cc.Clone()
}

func NewCoolorColorRetard(c *CoolorColor) *CoolorColorRetard {
	ccr := &CoolorColorRetard{
		CoolorColor: c.Clone(),
	}
	return ccr
}

func (ccr *CoolorColorRetard) Lighten(amount float64) *CoolorColorRetard {
	r, g, b := ccr.color.RGB()
	n := noire.NewRGB(float64(r), float64(g), float64(b))
	_ = n.Lighten(amount)
	cful, _ := colorful.MakeColor(ccr)
	_ = cful
	// newc := color.Lighten(color.NewHSLA(cful), amount)
	// color.Warmer(newc)
	// ccr.SetColorCss(newc.ToRGBA().ToHex())
	h := n.Lighten(0.5).Hex()
	return NewCoolorColorRetard(NewCoolorColor(fmt.Sprintf("#%s", h)))
}

type CoolorColorEdit struct {
	original *CoolorColor
	*CoolorColorRetard
}

func NewCoolorColorEdit(c *CoolorColor) *CoolorColorEdit {
	cce := &CoolorColorEdit{
		original:          c.Clone(),
		CoolorColorRetard: NewCoolorColorRetard(c),
	}
	cce.original.updateStyle()
	cce.updateStyle()
	return cce
}

func (cce *CoolorColorEdit) Reset(c *CoolorColor) {
	cce.original = c.Clone()
	cce.CoolorColorRetard = NewCoolorColorRetard(c.Clone())
	// cce = NewCoolorColorEdit(c)
	cce.original.updateStyle()
	cce.updateStyle()
}

type Retards []*CoolorColorRetard

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
	*tview.Grid
	cces      *CoolorColorEditorStrip
	increment float32
	size      int
	value     float32
	*colorgrad.GradientBuilder
}

func NewGradStrip(value float32, cces *CoolorColorEditorStrip) *GradStrip {
	gs := &GradStrip{
		Grid:            tview.NewGrid(),
		cces:            cces,
		size:            10,
		value:           value,
		increment:       0.1,
		GradientBuilder: colorgrad.NewGradient(),
	}
	return gs
}

func (gs *GradStrip) updateStrip() {
	_, _, _, h := gs.Grid.GetRect()
	gs.SetSize(gs.size, 1, -1, -1)
	gs.SetMinSize(int(math.Floor(float64(h/int(gs.size)))), 1)
	gs.GradientBuilder.Colors()
}

type CoolorColorEditorStrip struct {
	Name string
	*tview.Flex
	mainColor *CoolorColorEdit
	// shades    *CoolorColorRetard
	increment       float64
	hue, sat, light *GradStrip
}

func (ccep *CoolorColorEditorStrip) UpdateColor(c *CoolorColor) {
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

func NewEditorStrip(name string, middle *CoolorColor) *CoolorColorEditorStrip {
	ccep := &CoolorColorEditorStrip{
	}
	return ccep
}

type CoolorColorEditor struct {
	palette  *CoolorPalette
	cce      *CoolorColorEdit
	previews *CoolorColorEditorStrip
	*tview.Flex
	app *tview.Application
}

func NewCoolorEditor(a *tview.Application, cp *CoolorPalette) *CoolorColorEditor {
	cce := &CoolorColorEditor{}
	cce.palette = cp
	cce.app = a
	cce.Init()
	return cce
}

func (cce *CoolorColorEditor) Init() {
	cce.Flex = tview.NewFlex().SetDirection(tview.FlexColumn)
	cce.palette.AddEventHandler("selected", cce)
	cce.Flex.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 1, false), 0, 1, false)
	cce.Flex.AddItem(cce.cce, 0, 1, true)
	cce.Flex.AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(tview.NewBox(), 0, 1, false), 0, 5, false)
}

func (cce *CoolorColorEditor) HandleEvent(e tcell.Event) bool {
	switch e.(type) {
	case *SelectionEvent:
		se := e.(*SelectionEvent)
		cce.cce.Reset(se.color.Clone())
		return true
	}
	return false
}

// vim: ts=2 sw=2 et ft=go
