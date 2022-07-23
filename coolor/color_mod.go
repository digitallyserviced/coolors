package coolor

import (
	"fmt"
	"math"
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
	chm  *ChannelMod
	*ChannelModOptions
}

type ColorMod struct {
	orig    *CoolorColor
	current CoolorColor
	ring    CoolorColors
	*ColorModifier
}

func NewColorMod(mod *ChannelModOptions, chm ChannelModifier) *ColorMod {
	cc := NewRandomCoolorColor()
	cm := &ColorMod{
		ColorModifier: &ColorModifier{},
	}
	cm.orig = cc
	cm.current = *cc.Clone()
	cm.ColorModifier.name = ""
	cm.ColorModifier.chm = chm(mod)
	cm.ColorModifier.ChannelModOptions = mod
	cm.updateState(true)
	return cm
}

func (cm *ColorModifier) updateState(noUpdateColor bool) {
	// if !noUpdateColor {
	// 	cm.mid = cm.GetChannelValue(&cm.current)
	// }
	increment := cm.increment
	incrs := float64(1.0 / increment)
	sizeIncrs := float64(1.0 / cm.size)
	if sizeIncrs < increment {
		increment = sizeIncrs
	}
	count := 1.0 / sizeIncrs
	if math.Mod(incrs, 2) != 0 {
		incrs = incrs - 1
	}

	split := incrs / 2
	diff := (split * increment)
	// cm.above = math.Floor(diff/increment) - 1
	// cm.below = math.Floor(diff/increment) - 1
	// cm.increment = increment
	cm.count = count
	cm.diff = diff
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
	cnew := cm.chm.ModPct(cm.current.Clone(), cm.increment)
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
	// fmt.Println(cm)

	// dump.P(cm)
	// dump.
	// dump.Dump(cm)
	num := cm.size
	if above {
		num = (-math.Abs(cm.size))
	}
	// fmt.Println(cm.below, cm.above, cm.increment)
	colors := cm.chm.RangePct(&cm.current, cm.increment, num)
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

func (cm *ColorMod) Below() []CoolorColor {
	// num := clamp(cm.diff*cm.chm.Max(), cm.chm.Min(), cm.chm.Max())
	return cm.makeGrad(false)
}

func (cm *ColorMod) Above() []CoolorColor {
	// num := clamp(cm.diff*cm.chm.Max(), cm.chm.Min(), cm.chm.Max())
	return cm.makeGrad(true)
	// return cm.makeGrad(math.Abs(num), true)
}

func (cm *ColorMod) SetColor(cc CoolColor) {
	cm.orig = cc.GetCC()
	cm.current = *cm.orig.Clone()
	cm.updateState(false)
}

func (cm *ColorMod) GetCurrentChannelValue() float64 {
  return cm.GetChannelValue(&cm.current)
	// return cm.mid
}

func (cm *ColorModifier) String() string {
	return fmt.Sprintf(
		"ColorMod: %s mid: %f count: %f diff: %f increment %f size %f",
		cm.name,
		cm.mid,
		cm.count,
		cm.diff,
		cm.increment,
		cm.size,
	)
}

func (cm *ColorMod) GetChannelValue(cc CoolColor) float64 {
	value := cm.chm.GetChannelValue(&cm.current)
	return float64(value)
}

func (cm *ColorMod) Set(value float64) {
	value = cm.chm.Max() * value
	cm.current = *cm.chm.SetChannelValue(&cm.current, value).GetCC()
	cm.updateState(false)
}

func (cm *ColorMod) Incr(value float64) {
	if value == 0 {
		value = cm.increment
	}
	cm.current = *cm.chm.Mod(&cm.current, math.Abs(value*cm.chm.Max())).GetCC()
	cm.updateState(false)
}

func (cm *ColorMod) Decr(value float64) {
	if value == 0 {
		value = cm.increment
	}
	cm.current = *cm.chm.Mod(&cm.current, -math.Abs(value*cm.chm.Max())).GetCC()
	cm.updateState(false)
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
