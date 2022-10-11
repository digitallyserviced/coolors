package coolor

import (
	"github.com/digitallyserviced/tview"

	"github.com/digitallyserviced/coolors/theme"
)


type PaletteFloater struct {
	*tview.Flex
	Palette *CoolorPaletteContainer
}

func NewScratchPaletteFloater(cp *CoolorColorsPalette) *PaletteFloater {
	spf := &PaletteFloater{
		Flex:    tview.NewFlex(),
		Palette: NewCoolorPaletteContainer(cp),
	}

	spf.SetDirection(tview.FlexRow)
	spf.AddItem(nil, 0, 2, false)
	spf.AddItem(spf.Palette, 0, 4, true)
	spf.AddItem(nil, 0, 2, false)
	return spf
}
func NewCoolorPaletteContainer(
	cp *CoolorColorsPalette,
) *CoolorPaletteContainer {
	p := cp.GetPalette()
	p.Plainify(true)
	p.Sort()
	pt := NewPaletteTable(p)
	cpc := &CoolorPaletteContainer{
		Frame:   tview.NewFrame(pt),
		Palette: pt,
	}
	cpc.SetBorders(1, 1, 0, 0, 0, 0)
	cpc.SetTitle("")
	cpc.Frame.SetBorder(true).
		SetBorderPadding(0, 0, 1, 1).
		SetBorderColor(theme.GetTheme().TopbarBorder)
	return cpc
}

