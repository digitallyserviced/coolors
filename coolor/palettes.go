package coolor

import (
	"fmt"

	"github.com/digitallyserviced/coolors/status"
)

func (cbp *CoolorBlendPalette) GetPalette() *CoolorColorsPalette {
	return cbp.CoolorColorsPalette
}

func (cbp *CoolorShadePalette) UpdateColors(base *CoolorColor) {
	cbp.base = base
	cbp.Init()
}

func (cbp *CoolorBlendPalette) UpdateColors(start, end *CoolorColor) {
	cbp.start = start
	cbp.end = end
	cbp.Init()
}

func (cbp *CoolorShadePalette) Init() {
	cbp.ColorContainer.Clear()
	cbp.Colors = make(CoolorColors, 0)
	base, _ := MakeColor(cbp.base)
	done := make(chan struct{})
	// cbp.colors = make(CoolorColors, 0)
	defer close(done)
	colors := RandomShadesStream(base, 0.2)
	colors.Status.SetProgressHandler(NewProgressHandler(func(u uint32) {
		status.NewStatusUpdate(
			"action_str",
			fmt.Sprintf("Found Shades (%d / %d)", u, colors.Status.GetItr()),
		)
	}, func(i uint32) {
		status.NewStatusUpdate(
			"action_str",
			fmt.Sprintf("Found Shades (%d / %d)", colors.Status.GetValid(), i),
		)
	}))
	colors.Run(done)
	for _, v := range TakeNColors(done, colors.OutColors, int(cbp.increments)) {
		newcc := NewStaticCoolorColor(v.Hex())
		cbp.AddCoolorColor(newcc)
		// SeentColor("stream_random_shade", newcc, newcc.pallette)
	}
	// cbp.UpdateSize()
	//  cbp.ResetViews()
	// cbp.SetSelected(0)
}

func (cbp *CoolorBlendPalette) Init() {
	cbp.ColorContainer.Clear()
	cbp.Colors = make(CoolorColors, 0)
	incrSizes := 1.0 / cbp.increments
	start, _ := MakeColor(cbp.start)
	end, _ := MakeColor(cbp.end)
	for i := 0; i <= int(cbp.increments); i++ {
		newc := start.BlendLab(end, float64(i)*float64(incrSizes))
		newcc := NewStaticCoolorColor(newc.Hex())
		cbp.AddCoolorColor(newcc)
		// SeentColor("mixed_colors_gradient", newcc, newcc.pallette)
	}

	MainC.conf.AddPalette("blend", cbp)
}

