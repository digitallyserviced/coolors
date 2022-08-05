package coolor

import (
	"fmt"
	"math"
	// "github.com/gookit/goutil/dump"
)

type ChannelModOptions struct {
	name                     string
	increment, min, max      float64
	size, mid, scale         float64
	minIncrPct, minIncrValue float64
}

func NewChannelModOption(name string, min, max float64, defaulIncr, minIncrPct, minIncrValue, scale float64) *ChannelModOptions {
	cmo := &ChannelModOptions{
		name:         name,
		increment:    defaulIncr,
		min:          min,
		max:          max,
		size:         20,
		mid:          0,
		scale:        scale,
		minIncrPct:   minIncrPct,
		minIncrValue: minIncrValue,
	}
	return cmo
}

type ChannelModifier func(hmo *ChannelModOptions) *ChannelMod

type ChannelMod struct {
	cmo *ChannelModOptions
	// GetName         func() string
	FormatChannelValue func(cc CoolColor) string
	GetChannelValue    func(cc CoolColor) float64
	SetChannelValue    func(cc CoolColor, value float64) (CoolColor, bool)
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
	orig := cm.GetChannelValue(cc)
	start := orig
	rng := make([]int, int(math.Abs(num)))
	for i := range rng {
		newv := start + (increment * float64(i))
		color, _ := cm.SetChannelValue(cc, newv)
		cc := *color.GetCC()
		colors = append(colors, cc)
	}
	return colors
}

func (cm *ChannelMod) ModPct(cc CoolColor, value float64) CoolColor {
	num := value
	if num > 1 {
		num = num * 0.01
	}
	incr := num * cm.cmo.max
	ccn, _ := cm.Mod(cc, incr)
	return ccn
}

func (cm *ChannelMod) Set(cc CoolColor, value float64) (CoolColor, bool) {
  c, v := cm.SetChannelValue(cc, value)
	return c,v
}

func (cm *ChannelMod) Mod(cc CoolColor, value float64) (CoolColor, bool) {
	val := cm.GetChannelValue(cc)
	val = val + value
	ccn, v := cm.SetChannelValue(cc.GetCC().Clone(), val)
	return ccn,v
}

// func (cm *ChannelMod) GetStatus(cc CoolColor) string,string {
//   return cm.FormatChannelValue(cc)
// }
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
	HueModOptions   *ChannelModOptions = NewChannelModOption("Hue", 0.001, 360, 1, 1, 0.1, 1)
	SatModOptions   *ChannelModOptions = NewChannelModOption("Chroma", 0.001, 1.0, 0.001, 0.002, 0.002, 0.25)
	LightModOptions *ChannelModOptions = NewChannelModOption("Light", 0.001, 1.0, 0.001, 0.002, 0.002, 0.25)
)

var HueMod *ColorMod = NewColorMod(HueModOptions, hueFunc)
var SatMod *ColorMod = NewColorMod(SatModOptions, satFunc)
var LightMod *ColorMod = NewColorMod(LightModOptions, lightFunc)

func hueFunc(hmo *ChannelModOptions) *ChannelMod {
	cm := &ChannelMod{
		cmo: hmo,
		FormatChannelValue: func(cc CoolColor) string {
			hsl, _ := MakeColor(cc)
			l, s, h := hsl.LuvLCh()
			_, _ = l, s
			return fmt.Sprintf(" h = %0.2f° ", h)
		},
		GetChannelValue: func(cc CoolColor) float64 {
			hsl, ok := MakeColor(cc)
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = s, l
			return h
		},
		SetChannelValue: func(cc CoolColor, value float64) (CoolColor, bool) {
			hsl, ok := MakeColor(cc)
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = s, l
			h = value
			hsl = LuvLCh(l, s, h)
			return hsl,hsl.IsValid()
		},
	}
	return cm
}
func satFunc(hmo *ChannelModOptions) *ChannelMod {
	cm := &ChannelMod{
		cmo: hmo,
		FormatChannelValue: func(cc CoolColor) string {
			hsl, _ := MakeColor(cc)
			l, s, h := hsl.LuvLCh()
			_, _ = l, h
			return fmt.Sprintf(" ch = %0.1f ", s)
		},
		GetChannelValue: func(cc CoolColor) float64 {
			hsl, ok := MakeColor(cc)
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = h, l
			return s
		},
		SetChannelValue: func(cc CoolColor, value float64) (CoolColor, bool) {
			// value = math.Mod(value, hmo.max)
			hsl, ok := MakeColor(cc)
			// fmt.Println(hsl.GetCC().TerminalPreview())
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = h, l
			s = value
			hsl = LuvLCh(l, s, h)
			return hsl,hsl.IsValid()
		},
	}
	return cm
}
func lightFunc(hmo *ChannelModOptions) *ChannelMod {
	cm := &ChannelMod{
		cmo: hmo,
		FormatChannelValue: func(cc CoolColor) string {
			hsl, _ := MakeColor(cc)
			l, s, h := hsl.LuvLCh()
			_, _ = s, h
			return fmt.Sprintf(" l = %0.1f ", l)
		},
		GetChannelValue: func(cc CoolColor) float64 {
			hsl, ok := MakeColor(cc)
			// fmt.Println(hsl.GetCC().TerminalPreview())
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = h, s
			return l
		},
		SetChannelValue: func(cc CoolColor, value float64) (CoolColor, bool) {
			// value = math.Mod(value, hmo.max)
			hsl, ok := MakeColor(cc)
			// fmt.Println(hsl.GetCC().TerminalPreview())
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = h, s
			l = value
			hsl = LuvLCh(l, s, h)
			return hsl, hsl.IsValid()
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
