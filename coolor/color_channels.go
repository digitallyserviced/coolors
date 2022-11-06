package coolor

import (
	"fmt"
	"math"

	// "github.com/gookit/goutil/dump"

	"github.com/digitallyserviced/coolors/coolor/util"
)

type ChannelRange struct {
	Min, Max float64
	Step     float64
}

type ColorChannel struct {
	Get GetChannelFunc
	Set ModChannelFunc
	*ChannelRange
	name   string
	format string
	scale  float64
}

type ColorSpace struct {
	name     string
	Channels []ColorChannel
}

type GetChannelFunc func(c Color) float64
type ModChannelFunc func(v float64, c Color) Color

var (
	SetHue ModChannelFunc = func(tc float64, v Color) Color {
		// dump.P(tc, func(args... float64) []float64{return args}(v.Hsl()))
		_, b, c := v.Hsl()
		return Hsl(tc, b, c)
	}
	SetSaturation ModChannelFunc = func(tc float64, v Color) Color {
		// dump.P(tc, func(args... float64) []float64{return args}(v.Hsl()))
		a, _, c := v.Hsl()
		return Hsl(a, tc, c)
	}
	SetLightness ModChannelFunc = func(tc float64, v Color) Color {
		// dump.P(tc, func(args... float64) []float64{return args}(v.Hsl()))
		a, b, _ := v.Hsl()
		return Hsl(a, b, tc)
	}
	GetHue GetChannelFunc = func(c Color) float64 {
		v, _, _ := c.Hsl()
		return v
	}
	GetSaturation GetChannelFunc = func(c Color) float64 {
		_, v, _ := c.Hsl()
		return v
	}
	GetLightness GetChannelFunc = func(c Color) float64 {
		_, _, v := c.Hsl()
		return v
	}
	GetBlue GetChannelFunc = func(c Color) float64 {
		r, g, b := c.RGB255()
		_, _, _ = r, g, b
		return float64(b) 
	}
	GetGreen GetChannelFunc = func(c Color) float64 {
		r, g, b := c.RGB255()
		_, _, _ = r, g, b
		return float64(g) 
	}
	GetRed GetChannelFunc = func(c Color) float64 {
		r, g, b := c.RGB255()
		_, _, _ = r, g, b
		return float64(r) 
	}
	SetRed ModChannelFunc = func(tc float64, c Color) Color {
    var v uint8 = uint8(math.RoundToEven(tc))
		r, g, b := c.RGB255()
		_, _, _ = r, g, b
		return RGB255(v, g, b)
	}
	SetGreen ModChannelFunc = func(tc float64, c Color) Color {
    var v uint8 = uint8(math.RoundToEven(tc))
		r, g, b := c.RGB255()
		_, _, _ = r, g, b
		return RGB255(r, v, b)
	}
	SetBlue ModChannelFunc = func(tc float64, c Color) Color {
    var v uint8 = uint8(math.RoundToEven(tc))
		r, g, b := c.RGB255()
		_, _, _ = r, g, b
		return RGB255(r, g, v)
	}
	ChannelRed = ColorChannel{
		Get:          GetRed,
		Set:          SetRed,
		ChannelRange: &ChannelRange{0, 255.0, 1},
		name:         "R",
		scale:        1,
	}
	ChannelBlue = ColorChannel{
		name:         "B",
		Get:          GetBlue,
		Set:          SetBlue,
		ChannelRange: &ChannelRange{0, 255.0, 1},
		scale:        1,
	}
	ChannelGreen = ColorChannel{
		name:         "G",
		Get:          GetGreen,
		Set:          SetGreen,
		ChannelRange: &ChannelRange{0, 255.0, 1},
		scale:        1,
	}
	RGBChannel = ColorSpace{
		name:     "RGB",
		Channels: []ColorChannel{0: ChannelRed, 1: ChannelGreen, 2: ChannelBlue},
	}
	ChannelHue = ColorChannel{
		name:         "Hue",
		Get:          GetHue,
		Set:          SetHue,
		ChannelRange: &ChannelRange{0.0, 360.0, 0.5},
		format:       "% 3.2f",
		scale:        1 / 360.0,
	}
	ChannelSaturation = ColorChannel{
		name:         "Saturation",
		Get:          GetSaturation,
		Set:          SetSaturation,
		ChannelRange: &ChannelRange{0.0, 1.0, 0.001},
		format:       "%0.2f",
		scale:        1 / 1.0,
	}
	ChannelLightness = ColorChannel{
		name:         "Lightness",
		Get:          GetLightness,
		Set:          SetLightness,
		ChannelRange: &ChannelRange{0.0, 1.0, 0.001},
		format:       "%0.2f",
		scale:        1 / 1.0,
	}
	HSLChannel = ColorSpace{
		name: "HSL",
		Channels: []ColorChannel{
			0: ChannelHue,
			1: ChannelSaturation,
			2: ChannelLightness,
		},
	}
)

func (mod ColorChannel) Decr(v float64, c Color) Color {
	return mod.ModColor(-v, c)
}
func (mod ColorChannel) Incr(v float64, c Color) Color {
  fmt.Println(v, c)
	return mod.ModColor(v, c)
}
func (mod ColorChannel) ModColor(v float64, c Color) Color {
	o := mod.Get(c)
	o += v
	o = util.Clamp(o, mod.Min, mod.Max)
  fmt.Println(mod.Max, mod.Min, o)
	return mod.Set(o, c)
}
func (mod ColorChannel) SetColor(v float64, c Color) Color {
	return mod.Set(v, c)
}

func (mod ColorChannel) ValueFormat(c Color) string {
	val := mod.format
	if val == "" {
		val = "%3.0f"
	}

	return fmt.Sprintf(val, mod.Get(c))
}
func (mod ColorChannel) Display(c Color) string {
	return fmt.Sprintf("%s (%s)", mod.name, mod.ValueFormat(c))
}
