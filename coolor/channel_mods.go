package coolor

import (
	"fmt"
	// "github.com/gookit/goutil/dump"
)

var (
	HueModOptions   *ChannelModOptions = NewChannelModOption("Hue", 0.001, 360, 1, 1, 0.1, 1)
	SatModOptions   *ChannelModOptions = NewChannelModOption("Chroma", 0.001, 1.0, 0.001, 0.002, 0.002, 0.25)
	LightModOptions *ChannelModOptions = NewChannelModOption("Light", 0.001, 1.0, 0.001, 0.002, 0.002, 0.25)
)

var HueMod *ColorMod = NewChannelMod(HueModOptions, hueFunc)
var SatMod *ColorMod = NewChannelMod(SatModOptions, satFunc)
var LightMod *ColorMod = NewChannelMod(LightModOptions, lightFunc)

func hueFunc(hmo *ChannelModOptions) *ChannelMod {
	cm := &ChannelMod{
		cmo: hmo,
		FormatValue: func(cc CoolColor) string {
			hsl, _ := MakeColor(cc)
			l, s, h := hsl.LuvLCh()
			_, _ = l, s
			return fmt.Sprintf(" h = %0.2f° ", h)
		},
		GetValue: func(cc CoolColor) float64 {
			hsl, ok := MakeColor(cc)
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = s, l
			return h
		},
		SetValue: func(cc CoolColor, value float64) (CoolColor, bool) {
			hsl, ok := MakeColor(cc)
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = s, l
			h = value
			hsl = LuvLCh(l, s, h)
			return hsl, hsl.IsValid()
		},
	}
	return cm
}
func satFunc(hmo *ChannelModOptions) *ChannelMod {
	cm := &ChannelMod{
		cmo: hmo,
		FormatValue: func(cc CoolColor) string {
			hsl, _ := MakeColor(cc)
			l, s, h := hsl.LuvLCh()
			_, _ = l, h
			return fmt.Sprintf(" ch = %0.2f ", s)
		},
		GetValue: func(cc CoolColor) float64 {
			hsl, ok := MakeColor(cc)
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = h, l
			return s
		},
		SetValue: func(cc CoolColor, value float64) (CoolColor, bool) {
			hsl, ok := MakeColor(cc)
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = h, l
			s = value
			hsl = LuvLCh(l, s, h)
			return hsl, hsl.IsValid()
		},
	}
	return cm
}
func lightFunc(hmo *ChannelModOptions) *ChannelMod {
	cm := &ChannelMod{
		cmo: hmo,
		FormatValue: func(cc CoolColor) string {
			hsl, _ := MakeColor(cc)
			l, s, h := hsl.LuvLCh()
			_, _ = s, h
			return fmt.Sprintf(" l = %0.2f ", l)
		},
		GetValue: func(cc CoolColor) float64 {
			hsl, ok := MakeColor(cc)
			if !ok {
				fmt.Println("Error making color")
			}
			l, s, h := hsl.LuvLCh()
			_, _ = h, s
			return l
		},
		SetValue: func(cc CoolColor, value float64) (CoolColor, bool) {
			hsl, ok := MakeColor(cc)
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

// vim: ts=2 sw=2 et ft=go
