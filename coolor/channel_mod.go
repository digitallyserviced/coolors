package coolor

import (
	"math"
	// "github.com/gookit/goutil/dump"
)

type ModOptions struct {
	name                     string
	increment, min, max      float64
	size, mid, scale         float64
	minIncrPct, minIncrValue float64
}

type ChannelModOptions struct {
	*ModOptions
}
type ColorModOptions struct {
	*ModOptions
}

func NewChannelModOption(name string, min, max float64, defaulIncr, minIncrPct, minIncrValue, scale float64) *ChannelModOptions {
	mo := &ModOptions{
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
	cmo := &ChannelModOptions{
		ModOptions: mo,
	}
	return cmo
}

func NewColorModOption(name string, min, max float64, defaulIncr, minIncrPct, minIncrValue, scale float64) *ColorModOptions {
	mo := &ModOptions{
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
	cmo := &ColorModOptions{
		ModOptions: mo,
	}
	return cmo
}

type ChannelModifier func(hmo *ChannelModOptions) *ChannelMod

type ChannelMod struct {
	cmo         *ChannelModOptions
	FormatValue func(cc CoolColor) string
	GetValue    func(cc CoolColor) float64
	SetValue    func(cc CoolColor, value float64) (CoolColor, bool)
}
type CoolorMod struct {
	cmo         *ColorModOptions
	FormatValue func(cc CoolColor) string
	GetValue    func(cc CoolColor) float64
	SetValue    func(cc CoolColor, value float64) (CoolColor, bool)
}

type ColorModificator interface {
	GetOptions() *ModOptions
	FormatValue(cc CoolColor) string
	GetValue(cc CoolColor) float64
	SetValue(cc CoolColor, value float64) (CoolColor, bool)
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
	orig := cm.GetValue(cc)
	start := orig
	rng := make([]int, int(math.Abs(num)))
	for i := range rng {
		newv := start + (increment * float64(i))
		color, _ := cm.SetValue(cc, newv)
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
	c, v := cm.SetValue(cc, value)
	return c, v
}

func (cm *ChannelMod) Get(cc CoolColor) (float64) {
	val := cm.GetValue(cc)
	// val = val + value
	// ccn, v := cm.SetValue(cc.GetCC().Clone(), val)
	return val
}
func (cm *ChannelMod) Mod(cc CoolColor, value float64) (CoolColor, bool) {
	val := cm.GetValue(cc)
	val = val + value
	ccn, v := cm.SetValue(cc.GetCC().Clone(), val)
	return ccn, v
}

func (cm *ChannelMod) GetOptions() *ModOptions {
	return cm.cmo.ModOptions
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

// vim: ts=2 sw=2 et ft=go
