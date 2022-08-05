package coolor

import (
	"fmt"
	"math"
	"strings"

	"github.com/gookit/goutil/dump"
	"github.com/samber/lo"
	// "github.com/lucasb-eyer/go-colorful"
)

type CoolColorMod struct {
	*CoolorColor
	// *colorful.Color
}

func (ccm *CoolColorMod) GetCC() *CoolorColor {
	return ccm.GetCC()
}

type ColorModifier struct {
	name string
	ChannelMod  *ChannelMod
	*ChannelModOptions
}

type ColorModFunction func(float64) *ColorModAction
// len(ColorModNames) 🮋🮒🮑🮐🮆🮔🮕🮖🮗🮟🮞🮝🮜🮘🮙🮚🮱🮴🮽🮾🮿🯄

var ColorModActionStrings = map[string]string{
  "set": "[blue:black:-][black:blue:-]=[blue:black:-][-:-:-]",
  "inc": "[green:black:-][black:green:-]+[green:black:-][-:-:-]",
  "dec": "[red:black:-][black:red:-]-[red:black:-][-:-:-]",
}
type ColorModAction struct {
  Function string
  Action ColorModFunction
  Argument float64
  Result CoolColor
}

func (cma *ColorModAction) Summary() string {
  return fmt.Sprintf("%s %0.2f %s", ColorModActionStrings[cma.Function], cma.Argument, cma.Result.GetCC().TVPreview())
}

func (cma *ColorModAction) String() string {
  return fmt.Sprintf("%s", ColorModActionStrings[cma.Function])
}
func (cmlog ColorModActions) String() string {
  // _ = lo.Range(5)
  // return ""
  // dump.P("undo shts")
  summActions := lo.Map(cmlog, func (x *ColorModAction, n int) string {
    if x != nil {

    return x.Summary()
    }
    return ""
  })

  _ = summActions
  // lo.Reduce
  return strings.Join(summActions, "\n")
}

type ColorModActions []*ColorModAction
// func (cmas *ColorModActions) String() string {
//   // str
//   return fmt.Sprintf("%s", ColorModActionStrings[cma.Function])
// }


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

type ColorModder interface {
  Incr(float64)
  Decr(float64)
  Set(float64)
  Nop(bool)
}

type ColorMod struct {
  last *ColorModAction
  history ColorModActions 
	orig    *CoolorColor
	current CoolorColor
	ring    CoolorColors
	*ColorModifier
}

func NewColorModAction(name string, f ColorModFunction, arg float64, result CoolColor) *ColorModAction {
  return &ColorModAction{
    Function: name,
    Action: f,
    Argument: arg,
    Result: result,
  }
}

func (cm *ColorMod) Log(action *ColorModAction) **ColorModAction {
  cm.history = lo.Subset(cm.history, -19, 19)
  cm.history = append(cm.history, action)
  cm.last = action
  dump.P(cm.history.String())
  return &cm.last
}

func (cm *ColorMod) Set(value float64) *ColorModAction {
	value = cm.ChannelMod.Max() * value
  c, _ := cm.ChannelMod.SetChannelValue(&cm.current, value) // .GetCC()
  cm.current = *c.GetCC()
	cm.updateState(false)
  return *cm.Log(NewColorModAction("set",cm.Set, value, cm.current.Clone()))
}

func (cm *ColorMod) Incr(value float64) *ColorModAction {
	if value == 0 {
		value = cm.increment
	}
  cnew, err := cm.ChannelMod.Mod(cm.current.GetCC(), math.Abs(cm.increment))
  cm.current = *cnew.GetCC()
  // dump.P(value,cm.current.TerminalPreview())
  _ = err
	cm.updateState(false)
  return *cm.Log(NewColorModAction("inc",cm.Incr, value, cm.current.Clone()))
}

func (cm *ColorMod) Decr(value float64) *ColorModAction {
	if value == 0 {
		value = cm.increment
	}
  cnew,err := cm.ChannelMod.Mod(cm.current.GetCC(), -math.Abs(cm.increment))
  cm.current = *cnew.GetCC()
  // dump.P(value,cm.current.TerminalPreview())
  _ = err
	cm.updateState(false)
  return *cm.Log(NewColorModAction("dec",cm.Decr, value, cm.current.Clone()))
}

func NewColorMod(mod *ChannelModOptions, chm ChannelModifier) *ColorMod {
	cc := NewRandomCoolorColor()
	cm := &ColorMod{
		ColorModifier: &ColorModifier{},
	}
  cm.history = make(ColorModActions, 0)
  cm.last = NewColorModAction("nop", cm.Set, 0.0, nil)
	cm.orig = cc
	cm.current = *cc.Clone()
	cm.ColorModifier.name = mod.name
	cm.ColorModifier.ChannelMod = chm(mod)
	cm.ColorModifier.ChannelModOptions = mod
	cm.updateState(true)
	return cm
}
func (cm *ColorModifier) SetSize(size float64) {
	cm.size = size
	cm.updateState(false)
}
func (cm *ColorModifier) updateState(noUpdateColor bool) {
	// if !noUpdateColor {
	// 	cm.mid = cm.GetChannelValue(&cm.current)
	// }
	// increment := cm.increment
	// incrs := float64(1.0 / increment)

  // cm.increment = cm.minIncrValue
  // cm.increment = cm.increment * ((cm.size / cm.minIncrValue) + 2)
  // dump.P("incr: ,", cm.increment)
  // cm.increment = i
  // cm.increment := cm.minIncrValue // / cm.max
  // cm.increment = increment
	// sizeIncrs := float64(1.0 / cm.size)
	// count := 1.0 / sizeIncrs
	// split := incrs / 2
	// diff := (split * increment)
	// cm.above = math.Floor(diff/increment) - 1
	// cm.below = math.Floor(diff/increment) - 1
	// cm.increment = increment
	// cm.count = count
	// cm.diff = diff
}

func (cm *ColorMod) Pop(cc CoolColor) {
	cm.ring = append(cm.ring, cm.current.GetCC())
	cm.current = *cc.GetCC()
}

func (cm *ColorMod) Push(cc CoolColor) {
	cm.ring = append(cm.ring, cm.current.GetCC())
	cm.current = *cc.GetCC()
}

func (cm *ColorMod) Next() CoolColor {
	// fmt.Printf("val: %f\n", cm.increment)
	// fmt.Println(cm.current.TerminalPreview())
	cnew := cm.ChannelMod.ModPct(cm.current.Clone(), cm.increment)
	// fmt.Printf("color: %s\n", cnew.GetCC().TerminalPreview())
	if cnew.GetCC().Html() == cm.current.GetCC().Html() {
		// fmt.Printf("Colors are the same %06x", cnew.color.Hex())
		return cnew.GetCC()
	}
	if len(cm.ring) < int(cm.size) {
		cm.Push(cnew.GetCC())
	}
	return cnew.GetCC()
}

func (cm *ColorMod) makeGrad(above bool) []CoolorColor {
	num := cm.size
	if above {
		num = (-math.Abs(cm.size))
	}
	fmt.Println(cm.size, cm.increment)
	colors := cm.ChannelMod.RangePct(&cm.current, cm.increment, num)
	return colors
}

// func (cm *ColorMod) makeGrad(diff float64, invert bool) []colorful.Color {
// 	endcc := cm.chm.Mod(&cm.current, diff, false)
// 	// fmt.Println(endcc.GetCC())
// 	cg := colorgrad.NewGradient()
// 	if invert {
// 		cg.HtmlColors(
// 			endcc.Html(), cm.current.Html(),
// 		)
// 	} else {
// 		cg.HtmlColors(
// 			cm.current.Html(), endcc.Html(),
// 		)
// 	}
// 	cg.Mode(colorgrad.BlendLinearRgb).Build()
// 	grad, _ := cg.Build()
// 	return grad.ColorfulColors(uint(cm.above))
//
// }

// func (cm *ColorMod) Below() []CoolorColor {
// 	// num := clamp(cm.diff*cm.chm.Max(), cm.chm.Min(), cm.chm.Max())
// 	return cm.makeGrad(false)
// }
//
// func (cm *ColorMod) Above() []CoolorColor {
// 	// num := clamp(cm.diff*cm.chm.Max(), cm.chm.Min(), cm.chm.Max())
// 	return cm.makeGrad(true)
// 	// return cm.makeGrad(math.Abs(num), true)
// }

func (cm *ColorMod) SetColor(cc CoolColor) {
	cm.orig = cc.GetCC()
	cm.current = *cm.orig.Clone()
	cm.updateState(false)
}

func (cm *ColorMod) GetCurrentChannelValue() float64 {
	return cm.GetChannelValue(&cm.current)
}
// len(ColorModNames) 
// func (cm *ChannelMod) GetStatus(cc CoolColor) string {
//   return cm.FormatChannelValue(cc)
// }

func (cm *ColorMod) GetStatus() (string,string) {
  last := ""
  status := ""
  if cm.last != nil {
    last = cm.last.String()
  }
  status = cm.ChannelMod.FormatChannelValue(&cm.current)
	return status, last
}

func (cm *ColorModifier) String() string {
	return fmt.Sprintf(
		"%s +/- %0.2f ",
		cm.name,
		cm.increment,
	)
}

func (cm *ColorMod) GetChannelValue(cc CoolColor) float64 {
	value := cm.ChannelMod.GetChannelValue(&cm.current)
	return float64(value)
}

func clamped(val, min, max float64) (float64, bool) {
	clampd := false
	if val > max {
		clampd = true
	}
	if val < min {
		clampd = true
	}
	return clamp(val, min, max), clampd
}
func clamp(val, min, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}

// vim: ts=2 sw=2 et ft=go
