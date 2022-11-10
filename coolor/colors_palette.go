package coolor

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	// "github.com/gookit/goutil/dump"
	"github.com/samber/lo"

	"github.com/digitallyserviced/coolors/coolor/events"
	// . "github.com/digitallyserviced/coolors/coolor/events"
)

type (
	// CoolorColorsPaletteMeta struct {
 //  }
	CoolorColorsPalette struct {
		l              *sync.RWMutex
		// *events.EventObserver `msgpack:"-" clover:"-,omitempty"`
		*events.EventNotifier `msgpack:"-" clover:"-,omitempty"`
		tagType        *TagType `msgpack:"-" clover:"-,omitempty"`
		// Name           string
		ptype          string
		Colors         CoolorColors // `clover:"foreignKey:Color"`
		Hash           uint64       // `gorm:"uniqueIndex"`
		selectedIdx    int
    dirty bool
	}
)

func (cp *CoolorColorsPalette) GetItemCount() int {
	return cp.Len()
}

func (cp *CoolorColorsPalette) GetPalette() *CoolorColorsPalette {
	return cp
}

// func (cp *CoolorColorsPalette) Coolors() *Coolors {
// 	// colors :=
// 	css := &Coolors{
// 		Key:    cp.Name,
// 		Colors: make([]*Coolor, 0),
// 		Saved:  false,
// 	}
// 	for _, v := range cp.Colors {
// 		css.Colors = append(css.Colors, v.Coolor())
// 	}
// 	return css
// }

func (cp *CoolorColorsPalette) Each(f func(*CoolorColor, int)) {
	// MainC.app.QueueUpdate(func() {
	for i, v := range cp.Colors {
		f(v, i)
	}
	// })
}

func (cp *CoolorColorsPalette) Plainify(s bool) {
	cp.Each(func(cc *CoolorColor, _ int) {
		cc.SetStatic(s)
		cc.SetPlain(s)
	})
}

func (cp *CoolorColorsPalette) Staticize(s bool) {
	cp.Each(func(cc *CoolorColor, _ int) {
		cc.SetStatic(s)
	})
}

func (cp *CoolorColorsPalette) GetRef() interface{} {
	return cp
}

func (cp *CoolorColorsPalette) ClearSelected() {
	cp.Each(func(cc *CoolorColor, _ int) {
		cc.SetSelected(false)
	})
}

func (cp *CoolorColorsPalette) SetSelected(idx int) error {
	if len(cp.Colors) == 0 {
		return fmt.Errorf("no valid color at idx: %d", idx)
	}
	// MainC.app.QueueUpdateDraw(func() {
	if idx < 0 {
		idx = len(cp.Colors) - 1
	}
	if idx > len(cp.Colors)-1 {
		idx = 0
	}
	if idx < len(cp.Colors) {
		cp.ClearSelected()
		cp.selectedIdx = idx
		cp.Colors[cp.selectedIdx].SetSelected(true)
		cp.SpawnSelectionEvent(cp.Colors[cp.selectedIdx], cp.selectedIdx)
		// cp.SpawnSelectionEvent(cp.colors[cp.selectedIdx], cp.selectedIdx)
	}
	// })
	return nil
}

func (cc *CoolorColorsPalette) SpawnSelectionEvent(
	c *CoolorColor,
	idx int,
) bool {
	events.Global.Notify(
		*cc.NewObservableEvent(events.PaletteColorSelectedEvent, "palette_selected", c, cc),
	)
	cc.Notify(
		*cc.NewObservableEvent(events.PaletteColorSelectedEvent, "palette_selected", c, cc),
	)
	return true
}
func (cp *CoolorColorsPalette) IconPalette(
	mainFormat, selectedFormat string,
) []string {
	chars := make([]string, len(cp.Colors))

	for i, v := range cp.Colors {
		f := mainFormat
		if i == cp.selectedIdx {
			f = selectedFormat
		}
		chars[i] = fmt.Sprintf(f, v.Html())
	}

	return chars
}

func (cp *CoolorColorsPalette) TagsKeys(random bool) CoolorPaletteTagsMeta {
  // dump.P(cp)
	tagKeys := make(map[string]*Coolor)
	cptm := &CoolorPaletteTagsMeta{
		tagCount:     0,
		TaggedColors: tagKeys,
	}
	tlist := GetTerminalColorsAnsiTags()
	if Base16Tags.tagList == nil {
	}

	for _, v := range tlist.items {
		cptm.tagCount += 1
		// k := v.data[keyfield.name].(string)
		k := v.GetKey()
		if random {
			tagKeys[k] = cp.RandomColor().Coolor()
		} else {
			tagKeys[k] = nil // cp.RandomColor().Coolor()
		}
	}
	// return *cptm

	cp.Each(func(cc *CoolorColor, i int) {
		tags := cc.GetTags()
		for _, t := range tags {
      if t == nil {
        continue
      }
			k := t.GetKey()
			cptm.TaggedColors[k] = cc.Coolor()
			cptm.tagCount += 1
		}
	})

  // dump.P(cptm)
	return *cptm
}
// ▌▐
// 
// ██
// ████████████████████
// ▆ ▆ ▆ ▆ ▆ ▆ ▆ ▆ ▆ ▆ ▆

func (cp *CoolorColorsPalette) MakeMenuPalette(showSelected bool) []string {
	// main, sel := "[%s:-:-]ﱢ[-:-2:-]", "[%s:-:b]ﱡ[-:-:-]"
  // 
// ﱢ  ﱡ            ﱣ  ﱤ 
  // ██  ﯟ ⬢ ⬡
	// main, sel := "[%s:-:b]🮐[-:-:-]", "[%s:-:-]▉[-:-:-]"
	// main, sel := "[%s:-:b]▆[-:-:-]", "[%s:-:-]▆[-:-:-]"
	// main, sel := "[%s:-]⬢ [-:-:-]", "[%s:-:-]ﯟ [-:-:-]"
	main, sel := "[%s:-] [-:-:-]", "[%s:-:-] [-:-:-]"
	// main, sel := "[%s:-:b][-:-:-]", "[%s:-:-][-:-:-]"
	// main, sel := "[%s:-:b]ﱡ[-:-:-]", "[%s:-:-]ﱢ[-:-:-]"
	// main, sel := "[%s:-:-]▉▉[-:-:-]", "[%s:-:b]▉▉[-:-:-]"
	if !showSelected {
		main = sel
	}

  cols := cp.IconPalette(main, sel)
  chunks := lo.Chunk[string](cols, 8)
  lines := make([]string, len(chunks)+1)
  for i, v := range chunks {
    // str := fmt.Sprintf("%s%s", IfElseStr(i % 2 == 1, "", " "),strings.Join(v, " "))
    str := fmt.Sprintf("%s%s", "",strings.Join(v, " "))
    lines[i] = str
  }
	return lines
}
func (cp *CoolorColorsPalette) MakeSquarePalette(showSelected bool) []string {
	// main, sel := "[%s:-:-]ﱢ[-:-2:-]", "[%s:-:b]ﱡ[-:-:-]"
  // 
// ﱢ ﱡ      ﱣ ﱤ 
  // ██
	// main, sel := "[%s:-:b]🮐[-:-:-]", "[%s:-:-]▉[-:-:-]"
	main, sel := "[%s:-:b][-:-:-]", "[%s:-:-]█[-:-:-]"
	// main, sel := "[%s:-:b]ﱡ[-:-:-]", "[%s:-:-]ﱢ[-:-:-]"
	// main, sel := "[%s:-:-]▉▉[-:-:-]", "[%s:-:b]▉▉[-:-:-]"
	if !showSelected {
		main = sel
	}
	return cp.IconPalette(main, sel)
}

func (cp *CoolorColorsPalette) MakeDotPalette() []string {
	return cp.IconPalette("[%s:-:-]ﱤ [-:-:-]", "[%s:-:b] [-:-:-]")
}
 // █▉▊▋▌▍▎
func (cp *CoolorColorsPalette) SquishPaletteBar(
) []string {
	// chars := make([]string, 0)
  // cp.Colors.Contains(c *CoolorColor)
  csss := cp.Colors.Strings()
  cols := csss[0:]
  chunks := lo.Chunk[string](cols, 2)
  format := "[%s:%s:]▍[-:-:-]"

  bars := lo.Map[[]string, string](chunks, func(s []string, i int) string {
    if len(s) < 2 {
      s = append(s, "")
    }
    return fmt.Sprintf(format, s[0], s[1])
  })
	return bars
}
func (cp *CoolorColorsPalette) GetSelected() (*CoolorColor, int) {
  if cp == nil {
    return nil, 0
  }
	cp.l.RLock()
	defer cp.l.RUnlock()

	if len(cp.Colors) == 0 {
		return nil, -1
	}

	if cp.selectedIdx > len(cp.Colors)-1 {
		cp.selectedIdx = 0
	}
	if cp.selectedIdx < 0 {
		cp.selectedIdx = len(cp.Colors) - 1
	}

	if cp.Colors[cp.selectedIdx] != nil {
		return cp.Colors[cp.selectedIdx], cp.selectedIdx
	}
	return nil, -1
}

func (cp *CoolorColorsPalette) String() string {
	return strings.Join(
		lo.Map[*CoolorColor, string](
			cp.Colors,
			func(cc *CoolorColor, i int) string {
				return cc.TerminalPreview()
			},
		),
		" ",
	)
}

func (cp *CoolorColorsPalette) Sort() {
	sort.Sort(cp)
}

func (cp *CoolorColorsPalette) RandomColor() *CoolorColor {
	return cp.GetItem(uint(rand.Uint32()))
}
func (cp *CoolorColorsPalette) GetItem(idx uint) *CoolorColor {
	id := int(math.Mod(float64(idx), float64(len(cp.Colors))))
	// fmt.Println(id, float64(len(cp.Colors)))
	return cp.Colors[id%len(cp.Colors)]
}

func (cp *CoolorColorsPalette) RemoveItem(rcc *CoolorColor) {
	newColors := cp.Colors[:0]
  remIdx := -1
	cp.Each(func(cc *CoolorColor, i int) {
		if cc != rcc {
			newColors = append(newColors, cc)
		} else {
      remIdx = i
    }
	})
  if remIdx > 0 {
    cp.SetSelected(remIdx-1)
  }
	cp.Colors = newColors
	cp.PaletteEvent(events.PaletteColorRemovedEvent, rcc)
}

func (cp *CoolorColorsPalette) PaletteEvent(
	t events.ObservableEventType,
	color *CoolorColor,
) {
  if cp == nil || MainC == nil {
    return
  }
	cp.Notify(*MainC.NewObservableEvent(t, "palette_event", color, cp))
	MainC.EventNotifier.Notify(
		*MainC.NewObservableEvent(t, "palette_event", color, cp),
	)
  hash := cp.HashColors()
  if cp.Hash != hash {
    if cp == GetStore().Current.Palette {

    }

  }
}

func (cp *CoolorColorsPalette) SetColors(
	cols []tcell.Color,
) *CoolorColorsPalette {
	cp.Colors = make(CoolorColors, 0)
	for _, v := range cols {
		cp.AddCoolorColor(NewIntCoolorColor(v.Hex()))
		// cp.Colors = append(cp.Colors, )
	}

	return cp
}
func (cp *CoolorColorsPalette) HandleEvent(o events.ObservableEvent) bool {
	// var cp *CoolorColorsPalette
	// if cp, ok := o.Src.(*CoolorColorsPalette); ok {
	// 	if cp != nil {
	// 		var selcol *CoolorColor
	// 		selcol, sl.selectedIndex = cp.GetSelected()
	// 		sl.selectedColor = selcol.Color
	// 		// AppModel.app.QueueUpdateDraw(func() {
	// 		sl.UpdatePos(cp)
	// 		// })
	// 	}
	// }
	return true

}

// Name implements events.Observer
func (cp *CoolorColorsPalette) Name() string {
  return Generator().WithSeed(int64(cp.UpdateHash())).GenerateName(2)
	// return "coolors_color_palette"
}

func (cp *CoolorColorsPalette) Lock() {
  cp.l.Lock()
}

func (cp *CoolorColorsPalette) Unlock() {
  cp.l.Unlock()
}

func (cp *CoolorColorsPalette) AddCoolorColor(color *CoolorColor) *CoolorColor {
	cp.l.Lock()
	cp.Colors = append(cp.Colors, color)
	cp.l.Unlock()
	cp.PaletteEvent(events.PaletteColorModifiedEvent, color)
	return cp.Colors[len(cp.Colors)-1]
}

func (ccs *CoolorColorsPalette) Swap(a, b int) {
	if col := ccs.Colors[a]; col == nil {
		return
	}
	if col := ccs.Colors[b]; col == nil {
		return
	}
	if ccs.selectedIdx == a {
		ccs.Colors[b].SetSelected(true)
		ccs.Colors[a].SetSelected(false)
	} else {
		ccs.Colors[a].SetSelected(true)
		ccs.Colors[b].SetSelected(false)
	}
	ccs.Colors[a], ccs.Colors[b] = ccs.Colors[b], ccs.Colors[a]
}

func (ccs *CoolorColorsPalette) Less(a, b int) bool {
	return ccs.Colors[a].Color.Hex() < ccs.Colors[b].Color.Hex()
}

func (ccs *CoolorColorsPalette) Len() int {
	return len(ccs.Colors)
}

func (cp *CoolorColorsPalette) UpdateHash() uint64 {
  cp.Hash = cp.HashColors()
  return cp.Hash
}
func (cp *CoolorColorsPalette) HashColors() uint64 {
	// var hash uint64 = 0
  hashed := lo.Reduce[*CoolorColor, uint64](cp.Colors, func(h uint64, c *CoolorColor, i int) uint64 {
    return h + uint64(c.Color.Hex())
  }, 0)
	// for _, v := range cp.Colors {
	// 	hash += uint64(v.Color.Hex())
	// }
  cp.Hash = hashed
	return hashed
}

func (cp *CoolorColorsPalette) ToMap() map[string]string {
	outcols := make(map[string]string)
	for i, v := range cp.Colors {
		k := fmt.Sprintf("color%d", i)
		outcols[k] = v.Html()
	}
	return outcols
}
