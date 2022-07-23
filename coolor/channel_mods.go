package coolor

import (
	"fmt"
	"math"
)

type ChannelModOptions struct {
	name                                 string
	increment, min, max                  float64
	above, below, size, diff, count, mid float64
  minIncrPct, minIncrValue float64
}

func NewChannelModOption(name string, min, max float64, defaulIncr, minIncrPct, minIncrValue float64) *ChannelModOptions {
	cmo := &ChannelModOptions{
		name:         name,
		increment:    defaulIncr,
		min:          min,
		max:          max,
		above:        5,
		below:        5,
		size:         9,
		diff:         0,
		count:        0,
		mid:          0,
		minIncrPct:   minIncrPct,
		minIncrValue: minIncrValue,
	}
	return cmo
}

type ChannelModifier func(hmo *ChannelModOptions) *ChannelMod

type ChannelMod struct {
	cmo             *ChannelModOptions
	// GetName         func() string
	GetChannelValue func(cc CoolColor) float64
	SetChannelValue func(cc CoolColor, value float64) CoolColor
	// ModPct          func(cc CoolColor, value float64, negative bool) CoolorColor
	// Min             func() float64
	// Max             func() float64
}

func (cm *ChannelMod) RangePct(cc CoolColor, increment, num float64) []CoolorColor {
	if increment > 1 {
		increment = increment * 0.01
	}
	incr := math.Floor(increment * cm.cmo.max)
	return cm.Range(cc, incr, num)
}

func (cm *ChannelMod) Range(cc CoolColor, increment, num float64) []CoolorColor {
	colors := make([]CoolorColor, 0)
	start := cm.GetChannelValue(cc)
	increment = math.Floor(increment)
	end := cm.GetChannelValue(cc) + (increment * num)
	i := start
	// j := end
	if end < start {
		i = end
		// j = start
	}
	num = math.Abs(num)
	// fmt.Println(start, end, i, num, increment)
	for {
		if num == 0 {
			break
		}
		// fmt.Println(num, i, j, increment)
		color := cm.Set(cc, i)
		colors = append(colors, *color.GetCC())
		i += increment
		num--
	}
	return colors
}

func (cm *ChannelMod) ModPct(cc CoolColor, value float64) CoolColor {
	num := value
	if num > 1 {
		num = num * 0.01
	}
  incr := num*cm.cmo.max
  ccn := cm.Mod(cc, incr)
  // fmt.Println(incr, num)
  // fmt.Println(ccn.GetCC().TerminalPreview())
	return ccn
}

func (cm *ChannelMod) Set(cc CoolColor, value float64) CoolColor {
	// value = clamp(math.Floor(math.Mod(value, cm.cmo.max)), cm.cmo.min, cm.cmo.max)
	return cm.SetChannelValue(cc, value)
}

func (cm *ChannelMod) Mod(cc CoolColor, value float64) CoolColor {
	val := cm.GetChannelValue(cc)
	val = val + value
  // fmt.Printf("vals: %f %f\n", val, value)
  ccn:= cm.SetChannelValue(cc, val)
  // fmt.Println(ccn.GetCC().HSL())
	return ccn
}

func (cm *ChannelMod) GetName() string {
	return cm.cmo.name
}
func (cm *ChannelMod) Max() float64 {
	return cm.cmo.max
}

func (cm *ChannelMod) Min() float64 {
	return cm.cmo.min
}

var (
	HueModOptions *ChannelModOptions = NewChannelModOption("Hue", 0, 360, 0.003, 0.003, 1)
	SatModOptions *ChannelModOptions = NewChannelModOption("Sat", 0, 1.0, 0.002, 0.002, 0.002)
)

var HueMod *ColorMod = NewColorMod(HueModOptions, hueFunc)
var SatMod *ColorMod = NewColorMod(SatModOptions, satFunc)

func hueFunc(hmo *ChannelModOptions) *ChannelMod {
	cm := &ChannelMod{
		cmo: hmo,
		GetChannelValue: func(cc CoolColor) float64 {
			h, s, l := cc.GetCC().HSL()
      _,_ = s,l
			return h
		},
		SetChannelValue: func(cc CoolColor, value float64) CoolColor {
			h, s, l := cc.GetCC().HSL()
      fmt.Println(h,s,l)
			value = math.Mod(value, hmo.max)
			h = value
      cca := NewHSL(h,s,l)
      // fmt.Println(cca.TerminalPreview())
			return cca
		},
	}
	return cm
}
func satFunc(hmo *ChannelModOptions) *ChannelMod {
	cm := &ChannelMod{
		cmo: hmo,
		GetChannelValue: func(cc CoolColor) float64 {
			_, s, _ := cc.GetCC().HSL()
			return s
		},
		SetChannelValue: func(cc CoolColor, value float64) CoolColor {
			h, s, l := cc.GetCC().HSL()
			// value = clamp(math.Mod(value, hmo.max), hmo.min, hmo.max)
      value = math.Mod(value, hmo.max)
			s = value
      cca := NewHSL(h,s,l)
      // fmt.Println(cca.TerminalPreview())
			return cca
		},
	}
	return cm
}

// var SatMod *ColorMod = NewColorMod("Sat", *satFunc(0.0, 1.0), 0.05, 10)

// func satFunc(min, max float64) *ChannelMod {
// 	cm := &ChannelMod{
// 		GetChannelValue: func(cc CoolColor) float64 {
// 			_, s, _ := cc.GetCC().HSL()
// 			return s
// 		},
// 		SetChannelValue: func(cc CoolColor, value float64) CoolorColor {
// 			h, s, l := cc.GetCC().HSL()
// 			value = clamp(value, min, max)
// 			s = value
// 			hsl := colorful.Hsl(h, s, l)
// 			ccm := &CoolColorMod{&hsl}
// 			return *ccm.GetCC()
// 		},
// 	}
// 	return cm
// }

// vim: ts=2 sw=2 et ft=go
