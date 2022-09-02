package coolor

import (
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
)

type CoolorPaletteContainer struct {
	*tview.Frame
	Palette *PaletteTable
}

type PaletteTableCell struct {
	*tview.Box
	Color *CoolorColor
}
type PaletteTable struct {
	*tview.Flex
	Palette *CoolorPalette
}

func NewPaletteTableCell(cc *CoolorColor) *PaletteTableCell {
	ptc := &PaletteTableCell{
		// TableCell: tview.NewTableCell(fmt.Sprintf(" \n %s \n ", cc.TVPreview())),
		Box:   MakeBoxItem("", cc.Html()),
		Color: cc,
	}
	ptc.SetBackgroundColor(*cc.color)
	// ptc.TableCell.SetAlign(tview.AlignCenter).SetTransparency(false)
	return ptc
}

func NewPaletteTable(cp *CoolorPalette) *PaletteTable {
	pt := &PaletteTable{
		// Table:   tview.NewTable(),
		Flex:    tview.NewFlex(),
		Palette: cp,
	}
	return pt
}

func (pt *PaletteTable) Draw(s tcell.Screen) {
	pt.Clear()
	pt.SetDirection(tview.FlexColumn)
	x, y, w, h := pt.GetInnerRect()
	_, _, _, _ = x, y, w, h
	pt.Palette.Each(func(cc *CoolorColor, _ int) {
		pt.AddItem(NewPaletteTableCell(cc), 0, 1, false)
	})
	pt.Flex.Draw(s)
}

