package coolor

import (
	"fmt"
	"math"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/mazznoer/colorgrad"
	"github.com/teacat/noire"
)

type CoolColorMod struct {
	*noire.Color
}

func (ccm *CoolColorMod) GetCC() *CoolorColor {
	cc := NewCoolorColor(noire.RGBToHTML(ccm.Red, ccm.Blue, ccm.Green))
	return cc
}

type ColorMod struct {
	// G func(cc CoolColor, increment, size float64) (CoolColor, CoolColor, CoolColor)
	name                              string
	orig                              *CoolorColor
	current                           CoolorColor
	diff, count, mid, increment, size float64
	surrounds                         float64
	chm                               *ChannelMod
}

type ChannelMod struct {
	GetChannelValue func(cc CoolColor) float64
	Set             func(cc CoolColor, value float64) CoolorColor
	Mod             func(cc CoolColor, value float64) CoolorColor
	Min             func() float64
	Max             func() float64
}


func NewColorMod(name string, chm ChannelMod, increment, size float64) *ColorMod {
	cm := &ColorMod{
		name: name,
		size: size,
		chm:  &chm,
	}
	cm.updateState(true)
	return cm
}

func (cm *ColorMod) updateState(noUpdateColor bool) {
	if !noUpdateColor {
		cm.mid = cm.GetChannelValue(&cm.current)
		// fmt.Println(cm.mid)
	}
	increment := cm.increment
	incrs := float64(1.0 / increment)
	sizeIncrs := float64(1.0 / cm.size)
	if sizeIncrs > increment {
		increment = sizeIncrs
	}
	count := 1.0 / sizeIncrs
	if math.Mod(incrs, 2) != 0 {
		incrs = incrs - 1
	}

	split := incrs / 2
	diff := (split * increment)
	cm.surrounds = math.Floor(diff / increment)
	cm.increment = increment
	cm.count = count
	cm.diff = diff
}

func (cm *ColorMod) makeGrad(diff float64, invert bool) [] colorful.Color {
	endcc := cm.chm.Mod(&cm.current, diff)
  fmt.Println(endcc.GetCC())
	cg := colorgrad.NewGradient()
  if invert {
    cg.HtmlColors(
      endcc.Html(),cm.current.Html(),
    )
  } else {
    cg.HtmlColors(
      cm.current.Html(),endcc.Html(),
    )
  }
	cg.Mode(colorgrad.BlendHsv).Build()
  grad, _ := cg.Build()
	return grad.ColorfulColors(uint(cm.surrounds))

}

func (cm *ColorMod) Below() []colorful.Color {
	num := clamp(cm.diff*cm.chm.Max(), cm.chm.Min(), cm.chm.Max())
  return cm.makeGrad(-math.Abs(num), false)
}

func (cm *ColorMod) Above() []colorful.Color {
	num := clamp(cm.diff*cm.chm.Max(), cm.chm.Min(), cm.chm.Max())
  return cm.makeGrad(math.Abs(num), true)
}

func (cm *ColorMod) SetColor(cc CoolColor) {
	cm.orig = cc.GetCC()
	// cm.current = *cc.GetCC().Clone()
	cm.current = *cm.orig.Clone()
	cm.updateState(false)
}

func (cm *ColorMod) GetCurrentChannelValue() float64 {
	return cm.mid
}

func (cm *ColorMod) String() string {
	return fmt.Sprintf(
		"ColorMod: %s mid: %f count: %f diff: %f increment %f size %f %s",
		cm.name,
		cm.mid,
		cm.count,
		cm.diff,
		cm.increment,
		cm.size,
		cm.current.GetCC(),
	)
}

func (cm *ColorMod) GetChannelValue(cc CoolColor) float64 {
	value := cm.chm.GetChannelValue(&cm.current)
	return float64(value / cm.chm.Max())
}

func (cm *ColorMod) Set(value float64) {
	value = cm.chm.Max() * value
	cm.current = cm.chm.Set(&cm.current, value)
	cm.updateState(false)
}

func (cm *ColorMod) Incr(value float64) {
	if value == 0 {
		value = cm.increment
	}
	cm.current = cm.chm.Mod(&cm.current, math.Abs(value*cm.chm.Max()))
	cm.updateState(false)
}

func (cm *ColorMod) Decr(value float64) {
	if value == 0 {
		value = cm.increment
	}
	cm.current = cm.chm.Mod(&cm.current, -math.Abs(value*cm.chm.Max()))
	cm.updateState(false)
}

func clamp(val, min, max float64) float64 {
	return math.Max(min, math.Min(val, max))
}
var HueMod *ColorMod = NewColorMod("Hue", *hueFunc(0, 360), 0.1, 10)
func hueFunc(min, max float64) *ChannelMod {
	cm := &ChannelMod{
		GetChannelValue: func(cc CoolColor) float64 {
			h, _, _ := cc.GetCC().HSL()
			return h
		},
		Set: func(cc CoolColor, value float64) CoolorColor {
			h, s, l := cc.GetCC().HSL()
			value = clamp(value, min, max)
			h = value
			fmt.Println(h)
			nc := noire.NewHSL(h, s, l)
			ccm := &CoolColorMod{&nc}
			return *ccm.GetCC()
		},
		Mod: func(cc CoolColor, value float64) CoolorColor {
			h, s, l := cc.GetCC().HSL()
			value = clamp(value, min, max)
			h = h + value
			fmt.Println(h)
			nc := noire.NewHSL(h, s, l)
			ccm := &CoolColorMod{&nc}
			return *ccm.GetCC()
		},
		Min: func() float64 {
			return min
		},
		Max: func() float64 {
			return max
		},
	}
	return cm
}

// vim: ts=2 sw=2 et ft=go
