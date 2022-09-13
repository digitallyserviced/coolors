package coolor

import (
	"github.com/digitallyserviced/coolors/theme"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	// "github.com/gookit/goutil/dump"
)

type CoolorPaletteContainer struct {
	*tview.Frame
	Palette *PaletteTable
}

type PaletteTableCell struct {
	// *tview.Box
	*tview.TableCell
	Color *CoolorColor
}

type PaletteTable struct {
	*tview.Table
	Palette *CoolorColorsPalette
}

func (pt *PaletteTable) UpdateView(w, h int) {
	// pt.SetCell(0, 0, tview.NewTableCell(" ").SetTransparency(true))
	// pt.SetCell(0, cols+1, tview.NewTableCell(" ").SetTransparency(true))
	pt.ResetCells()
}

func (pt *PaletteTable) ResetCells() {
	x, y, w, h := pt.GetInnerRect()
	_, _, _, _ = x, y, w, h
	// cols := pt.Palette.Len()
	// colw := (w) / cols
	pt.Palette.Each(func(cc *CoolorColor, idx int) {
		tc := NewPaletteTableCell(cc)
		// tc.SetText(strings.Repeat("▉", colw))
		// tc.SetText("").SetBackgroundColor(*cc.color)
		tc.SetExpansion(1)
		pt.SetCell(0, idx, tc.TableCell)
		pt.SetCell(1, idx, tc.TableCell)
	})
}

func NewPaletteTableCell(cc *CoolorColor) *PaletteTableCell {
	tc := &PaletteTableCell{
		// Box:   MakeBoxItem("", cc.Html()),
		TableCell: tview.NewTableCell(""),
		Color:     cc,
	}
	// tc.SetSelectable(true)
	tc.SetAlign(AlignCenter)
	tc.SetTransparency(true) //.SetTextColor(cc.GetFgColor())
	// pt.SetBackgroundColor(*cc.color)
	// ptc.TableCell.SetAlign(tview.AlignCenter).SetTransparency(false)
	return tc
}

func NewPaletteTable(cp *CoolorColorsPalette) *PaletteTable {
	pt := &PaletteTable{
		Table:   tview.NewTable(),
		Palette: cp,
	}
  pt.SetSelectedStyle(tcell.StyleDefault.Normal())
  pt.SetSelectionChangedFunc(func(row, column int) {
    pt.GetCell(row,column).SetStyle(tcell.StyleDefault.Normal().Reverse(false).Blink(true))
  })
  pt.SetSelectedFunc(func(row, column int) {
    c:=pt.GetCell(row,column)
    c.SetStyle(tcell.StyleDefault.Normal().Reverse(false).Blink(true))
    c.SetText("*")
    // dump.P(row, column, pt.Palette.GetItem(column).TerminalPreview())
  })
	pt.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		pt.UpdateView(width, height)
		return x, y, width, height
	})
	pt.SetBackgroundColor(theme.GetTheme().ContentBackground)
	pt.SetSelectable(false, true)
  
	pt.SetOffset(0, 0)
	pt.SetSeparator(' ')
	return pt
}

func (pt *PaletteTable) Draw(s tcell.Screen) {
	pt.SetOffset(0, 0)
	pt.Table.Draw(s)
}
