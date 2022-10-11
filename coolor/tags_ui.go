package coolor

import (
	"github.com/digitallyserviced/tview"
	"github.com/gookit/goutil/dump"

	"github.com/digitallyserviced/coolors/theme"
)

type CoolorColorTagsView struct {
	// Frame *tview.Frame
	*tview.Flex
	*CoolorColorsPalette
	infoView *tview.Flex
	gridView *tview.Grid
	details  *CoolorColorDetails
}

func NewCoolorColorTags(cp *CoolorColorsPalette) *CoolorColorTagsView {
	cci := &CoolorColorTagsView{
		// Frame:       &tview.Frame{},
		Flex:        tview.NewFlex(),
		infoView:    tview.NewFlex(),
		gridView:    tview.NewGrid(),
	}
	// cci.infoView.SetDirection(tview.FlexRow)
  cci.Flex.AddItem(cci.gridView, 0, 10, true)

	cci.UpdatePalette(cp)
	return cci
}

func (cct *CoolorColorTagsView) UpdatePalette(cp *CoolorColorsPalette) {
cct.CoolorColorsPalette = cp
	cct.UpdateView()
}

func (cct *CoolorColorTagsView) UpdateStatus(aa string) {

}
func NewTagEditFloater(cp *CoolorColorsPalette) *RootFloatContainer {
        tagV := NewCoolorColorTags(cp)
	f := NewSizedFloater(90, 18, 0)
	f.Item = tagV 
	f.UpdateView()
	return f
}


func (cct *CoolorColorTagsView) UpdateView() {
	MainC.app.QueueUpdateDraw(func() {
    cct.gridView.SetGap(1, 1)
    cct.gridView.SetBorderPadding(0,0,1,1).SetBorder(false)
    tagKeys := cct.CoolorColorsPalette.TagsKeys(false)
    // tagKeys := cct.CoolorColorsPalette.TagsKeys(true)

    if tagKeys.tagCount == 0 {
      cct.UpdateStatus("No tags set. Set tags or use auto tagger.")
    }

    dump.P(tagKeys)
    dynTags := cct.tagType.tagList.items[:0]
    tags := cct.tagType.tagList.items[:0]

    for _, v := range cct.CoolorColorsPalette.tagType.tagList.TagListItems.items {
      f := v.GetFlag(FieldDynamic)
      if v.data[f.name].(bool) == true {
        dynTags = append(dynTags, v)
      } else {
        tags = append(tags, v)
      }
    }

    row := 0
    col := 0

    normalTagsGrid := tview.NewGrid()
    normalTagsGrid.SetGap(1, 1)
    cct.gridView.AddItem(normalTagsGrid, 0, 0, 1, 1, -1,7, false)
    for i, v := range dynTags {
      cctb := NewCoolorColorTagBox(v,0-i)
      color, ok := tagKeys.TaggedColors[v.GetKey()]
      if !ok || color == nil {
        continue
      }
      x,y,width,height := cctb.GetRect()
      _,_ = width,height
      cctb.SetRect(x, y, 7, 5)
      cctb.SetColor(color)
      normalTagsGrid.AddItem(cctb, i, col,1,1,-1,-1,false)
    }

    col = 1
    cols := 8
    row = 0

    normalTags := tview.NewGrid()
    normalTags.SetGap(1, 1)
    cct.gridView.SetBackgroundColor(theme.GetTheme().InfoLabel)
    cct.gridView.AddItem(normalTags, 0, 1, 1, 1, -1,8*10, false)
    cct.gridView.SetColumns(-10,-90)
    colSizes := make([]int, 8)
    rowSizes := make([]int, 2)
    for i, v := range tags {
      cctb := NewCoolorColorTagBox(v, i)
      color, ok := tagKeys.TaggedColors[v.GetKey()]
      if !ok {
        continue
      }
      colSizes[i%8] = 9
      x,y,width,height := cctb.GetRect()
      _,_ = width,height
      cctb.SetRect(x, y, 7, 5)
      cctb.SetColor(color)
      row = i / cols
      rowSizes[row] = 5
      normalTags.AddItem(cctb, (row),i % cols,1,1,-1,-1,true)
    }
    // rowss := lo.RepeatBy[int](3, func(i int) int {return 5})
    // colss := lo.RepeatBy[int](12, func(i int) int {return 7})
    normalTags.SetRows(rowSizes...)
    normalTags.SetColumns(colSizes...)
		// cct.Flex.SetBorder(false).SetBorderPadding(1, 1, 1, 1)
		cct.SetDirection(tview.FlexColumn)
    // cct.gridView.SetSize(3, 12, 5, 7)
		// cct.gridView.SetOffset(0, 0)
	})
}
