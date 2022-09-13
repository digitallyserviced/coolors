package coolor

import (
	"fmt"
	"log"
	"math"

	"github.com/digitallyserviced/coolors/theme"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	// "github.com/gookit/goutil/dump"
)

type CoolorColorsTable struct {
	tview.TableContentReadOnly
	*CoolorColorsPalette
	startIndex int
	rows       int
	cols       int
	ch         int
	cw         int
}

func NewCoolorColorTable() *CoolorColorsTable {
	cct := &CoolorColorsTable{
		TableContentReadOnly: tview.TableContentReadOnly{},
		// CoolorColorsPalette:  cp,
		startIndex:           0,
		rows:                 0,
		cols:                 0,
		ch:                   3,
		cw:                   8,
	}
	return cct
}

func (ccp *CoolorColorsPalette) GetColorIndex(row, rows, column, columns int) *CoolorColor {
	colorIdx := column + (row * (columns))
  if colorIdx < 0 || colorIdx > ccp.Len() - 1 {
    return nil
  }
	col := ccp.Colors[colorIdx]
  return col
}

func NewColorCell(col *CoolorColor) *tview.TableCell {
  if col == nil {
    return nil
  }

	tc := tview.NewTableCell(fmt.Sprintf(" %s ", col.TVPreview()))
	tc.SetAlign(AlignCenter)
	tc.SetStyle(
		tcell.StyleDefault.Background(col.GetFgColor()).Foreground(*col.Color),
	)
	tc.SetTransparency(true)
  return tc
}

func (cct *CoolorColorsTable) GetCell(row, column int) *tview.TableCell {
	// ▀ ▉░▒ ▒
  col := cct.GetColorIndex(row, cct.rows, column, cct.cols)
  if col == nil {
    return nil
  }
  tc := NewColorCell(col)
  // tc.SetReference(col.GetMeta())
	return tc
}

func (d *CoolorColorsTable) GetRowCount() int {
	return d.rows
}

func (d *CoolorColorsTable) GetColumnCount() int {
	return d.cols
}

type CoolorColorSwatch struct {
	*tview.Table
	TableContent *CoolorColorsTable
	*CoolorColorsPalette
	name      string
	getColors func(cs *CoolorColorSwatch) *CoolorColorsPalette
}

func NewColorSwatch() {
	if MainC.sidebar == nil {
		tv := NewCoolorColorSwatch(func(cs *CoolorColorSwatch) *CoolorColorsPalette { return NewCoolorColorsPalette() })
		tv.SetBackgroundColor(theme.GetTheme().SidebarBackground)
		tv.SetBorder(true).SetBorderPadding(1, 1, 1, 1)
		MainC.sidebar = NewFixedFloater(" Color Stash", tv)
		MainC.pages.AddPage("sidebar", MainC.sidebar.GetRoot(), true, true)
		MainC.pages.ShowPage("sidebar")
		MainC.app.SetFocus(MainC.sidebar.GetRoot().Item)
	} else {
		name, page := MainC.pages.GetFrontPage()
		if name == "sidebar" {
			MainC.pages.HidePage("sidebar")
			page.Blur()
			MainC.pages.RemovePage("sidebar")
			MainC.sidebar = nil
		} else {
			MainC.pages.ShowPage("sidebar")
			MainC.app.SetFocus(MainC.sidebar.GetRoot().Item)
		}
		AppModel.helpbar.SetTable("sidebar")
	}
}

func NewCoolorColorSwatch(f func(cs *CoolorColorSwatch) *CoolorColorsPalette) *CoolorColorSwatch {
	if f == nil {
    return nil
		// return &CoolorColorSwatch{
		// 	Table: tview.NewTable(),
		// 	TableContent: &CoolorColorsTable{
		// 		TableContentReadOnly: tview.TableContentReadOnly{},
		// 		CoolorColorsPalette:  &CoolorColorsPalette{},
		// 	},
		// 	getColors: func(cs *CoolorColorSwatch) *CoolorColorsPalette {
		// 		return NewCoolorColorsPalette()
		// 	},
		// }
	}
	ccs := &CoolorColorSwatch{
		Table:     tview.NewTable(),
		getColors: f,
		name: "",
	}

	// p := f(ccs)
	ccs.CoolorColorsPalette = f(ccs)
	ccs.TableContent = NewCoolorColorTable()
  ccs.UpdateItems()

	ccs.Table.SetContent(ccs.TableContent)
	ccs.Table.SetSelectable(true, true)
	ccs.Table.SetBordersColor(tview.Styles.PrimitiveBackgroundColor)
	ccs.Table.SetBorders(true).SetBorder(true).SetBorderPadding(0, 0, 1, 1)
  ccs.Table.SetSelectedFunc(func(row, column int) {
    col := ccs.GetColorIndex(row, ccs.TableContent.rows, column, ccs.TableContent.cols)
    log.Println(col)
    ccs.Notify(*ccs.NewObservableEvent(ColorSelectionEvent, "color_swatch", col, ccs))
  })
  ccs.Table.SetSelectionChangedFunc(func(row, column int) {
    col := ccs.GetColorIndex(row, ccs.TableContent.rows, column, ccs.TableContent.cols)
    log.Println(col)
    ccs.Notify(*ccs.NewObservableEvent(ColorSelectedEvent, "color_swatch", col, ccs))
  })
	ccs.SetDrawFunc(
		func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
      if ccs.CoolorColorsPalette.Len() == 0 {
        ccs.UpdateItems()
      }
			// dump.P(x, y, width, height)
			// x, y, width, height = ccs.GetInnerRect()
			// dump.P(x, y, width, height)
			x, y, width, height = ccs.GetRect()
			// dump.P(x, y, width, height)
			ccs.UpdateView()
			p := width - (ccs.TableContent.cols * 12)
			px := (p / 2) + 1
			return x + px, y, width, height
		},
	)
	return ccs
}

// func (ccs *CoolorColorSwatch) Draw(s tcell.Screen) {
// 	ccs.Box.SetDontClear(false)
// 	ccs.Box.DrawForSubclass(s, ccs)
//   ccs.UpdateView()
	// cols = (width - cols*2) / cct.cw
	// cct.cols = int(clamp(float64(cols), 1, float64(cols)-1))
// 	ccs.Grid.Draw(s)
// }
//
func (cct *CoolorColorsTable) UpdateView(x, y, width, height int) {
  if cct.CoolorColorsPalette == nil {
    return
  }
  items := len(cct.CoolorColorsPalette.Colors)
	cols := math.Floor(float64(width) / 12.0)
	rows := float64(items) / float64(cols)
	cct.rows = (int(math.Ceil(rows)))
  cct.cols = (int(cols))
  // dump.P(cols,cct.cols,rows,cct.rows)
}

func (ccs *CoolorColorSwatch) UpdateItems() {
	ccs.CoolorColorsPalette = ccs.getColors(ccs)
  ccs.TableContent.CoolorColorsPalette = ccs.CoolorColorsPalette
  ccs.Table.Select(0, 0)
  ccs.Table.ScrollToBeginning()
}

func (ccs *CoolorColorSwatch) hide() {
}

func (ccs *CoolorColorSwatch) show() {
  ccs.UpdateItems()
}

func (ccs *CoolorColorSwatch) Draw(s tcell.Screen) {
	ccs.Box.DrawForSubclass(s, ccs)
	tview.Borders = InvisBorders
	ccs.Table.Draw(s)
	tview.Borders = OrigBorders
}

func (ccs *CoolorColorSwatch) UpdateView() {
	x, y, width, height := ccs.GetInnerRect()
	ccs.TableContent.UpdateView(x, y, width, height)
}

type InputBubble interface {
	InputBubbler() func(event *tcell.EventKey, setFocus func(p tview.Primitive))
	InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive))
}

func (ccs *CoolorColorSwatch) InputBubbler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ccs.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		})
}

func (ccs *CoolorColorSwatch) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ccs.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			ccs.UpdateView()
      key := event.Key()
      switch key {
      case tcell.KeyRune:
        ch := event.Rune()
        switch ch {
        case 'f':
          r, c := ccs.Table.GetSelection()
          col := ccs.GetColorIndex(r, ccs.TableContent.rows, c, ccs.TableContent.cols)
          ccs.Notify(*ccs.NewObservableEvent(ColorSelectedEvent, "favorited", col, ccs))
          GetStore().MetaService.ToggleFavorite(col)
          return
        }
      }
      ccs.Table.InputHandler()(event, setFocus)
		})
}
	// main, sel := "[%s:-:-]▉▉[-:-:-]", "[%s:-:b]▉▉[-:-:-]"
  // dump.P(len(cct.CoolorColorsPalette.Colors))
	// colorIdx = int(clamp(float64(colorIdx), 0, float64(len(cct.CoolorColorsPalette.Colors)-1)))
  // dump.P(colorIdx)
  // fmt.Println(len(cct.Colors))
	// tc.SetTextColor()
	// row = int(float64(i / ccs.cols))
	// col = i % ccs.cols
//       row, col := ccs.Table.GetSelection()
// 		key := event.Key()
// 		// ch := event.Rune()
// 		// prevTab := ccs.selectedIdx
// 		switch {
// 		case key == tcell.KeyRight:
//       if ccs.TableContent.cols == col {
//
//
//       }
// 		case key == tcell.KeyLeft:
// 			ccs.selectedIdx+=1
// 		}
//     })
// }
