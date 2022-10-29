package coolor

import (
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/digitallyserviced/coolors/theme"
	// "github.com/digitallyserviced/coolors/theme"
	// "github.com/gookit/goutil/dump"
)

func init(){
  tview.Styles.PrimitiveBackgroundColor = theme.GetTheme().SidebarBackground
}

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
  TableContent *CoolorColorsTable
  cols,rows int
	*tview.Table
	Palette *CoolorColorsPalette
}

func (pt *PaletteTable) UpdateView() {
	x, y, width, height := pt.GetInnerRect()
	pt.TableContent.UpdateView(x, y, width, height)
	// pt.ResetCells()

}

func (pt *PaletteTable) ResetCells() {
	x, y, w, h := pt.GetInnerRect()
	_, _, _, _ = x, y, w, h
  // colw := (w)/ 12
	pt.Palette.Each(func(cc *CoolorColor, idx int) {
		tc := NewPaletteTableCell(cc)
		// tc.SetText(strings.Repeat("▉", colw)).SetTextColor(*cc.Color)
    tc.SetText(cc.TVPreview())
		// tc.SetText("").SetBackgroundColor(*cc.color)
		tc.SetExpansion(1)
    row := idx / pt.cols
		pt.SetCell(row, idx, tc.TableCell)
		// pt.SetCell(1, idx, tc.TableCell)
	})
}

func NewPaletteTableCell(cc *CoolorColor) *PaletteTableCell {
	tc := &PaletteTableCell{
		// Box:   MakeBoxItem("", cc.Html()),
		TableCell: tview.NewTableCell(""),
		Color:     cc,
	}
	tc.SetAlign(tview.AlignCenter)
	tc.SetStyle(
		tcell.StyleDefault.Background(cc.GetFgColor()).Foreground(*cc.Color),
	)
	tc.SetTransparency(true)
	return tc
}

func NewPaletteTable(cp *CoolorColorsPalette) *PaletteTable {
	pt := &PaletteTable{
		Table:   tview.NewTable(),
		Palette: cp,
    TableContent: NewCoolorColorTable(),
	}
  pt.TableContent.CoolorColorsPalette = pt.Palette
  pt.SetContent(pt.TableContent)
  pt.UpdateView()
	pt.Table.SetContent(pt.TableContent)
	pt.Table.SetSelectable(true, true)
	pt.Table.SetBordersColor(tview.Styles.PrimitiveBackgroundColor)
	pt.Table.SetBorders(true).SetBorder(true).SetBorderPadding(0, 0, 1, 1)
  // pt.Table.SetSelectedStyle(tcell.StyleDefault.Foreground(0).Background(tcell.Color238))
  pt.SetSelectionChangedFunc(func(row, column int) {
  // pt.Table.SetSelectedStyle(tcell.StyleDefault.Foreground(0).Background(tcell.Color238))
  })
  pt.SetSelectedFunc(func(row, column int) {
    // c:=pt.GetCell(row,column)
    // c.SetStyle(tcell.StyleDefault.Normal().Reverse(false).Blink(true))
    // c.SetText("*")
    // dump.P(row, column, pt.Palette.GetItem(column).TerminalPreview())
  })
	pt.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			// if pt.Palette.Len() == 0 {
			// 	pt.UpdateItems()
			// }
    // x,y,width,height = pt.GetRect()
    // pt.cols = (width - pt.Palette.Len()) / 14 
    // pt.rows = pt.Palette.Len() / pt.cols
	// colw := (width) / pt.cols
  // px := (width - (colw * 12)) / 2
			pt.UpdateView()
			p := width - (pt.TableContent.cols * 12)
			px := (p / 2) + 1
			return x + px, y, width, height
	})
	return pt
}

func (pt *PaletteTable) Draw(s tcell.Screen) {
	pt.Box.DrawForSubclass(s, pt)
	tview.Borders = InvisBorders
	pt.Table.Draw(s)
	tview.Borders = OrigBorders
}
