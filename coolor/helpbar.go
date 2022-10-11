package coolor

import (
	// "fmt"

	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	// "github.com/gookit/color"
	"github.com/gookit/goutil/maputil"
	// "github.com/samber/lo"

	"github.com/digitallyserviced/coolors/theme"
)

type HelpBar struct {
  *tview.Flex
	helpTextView *tview.TextView
	app *tview.Application
}

func NewHelpBar(app *tview.Application) *HelpBar {
	hb := &HelpBar{
    Flex: tview.NewFlex(),
		helpTextView: tview.NewTextView(),
		app:      app,
	}
	hb.Init()
	return hb
}

type shortcut struct {
	mods string
	icon string
	name string
	desc string
	key  rune
}

func NewShortcut(key rune, icon, name, desc string) *shortcut {
	sc := &shortcut{
		key:  key,
		mods: "",
		icon: icon,
		name: name,
		desc: "",
	}

	return sc
}

const (
	mainTable    = "main"
	paletteTable = "palette"
	editTable    = "editor"
)

type (
	keyMap    []*shortcut
	keyTables map[string]keyMap
)

var (
	paletteKeys keyMap
	mainKeys keyMap
	editKeys    keyMap
	table       keyTables
)

func init() {
	table = make(keyTables)
  mainKeys = make(keyMap, 0)
	paletteKeys = make(keyMap, 0)
	editKeys = make(keyMap, 0)

	mainKeys = append(mainKeys, NewShortcut('', "", "toggle help", " toggle help"))
	mainKeys = append(mainKeys, NewShortcut('q', "", "quit", " quit"))
	table[mainTable] = mainKeys

	paletteKeys = append(paletteKeys, NewShortcut('h', "ﰯ", "colors", " colors"))
	paletteKeys = append(paletteKeys, NewShortcut('l', "ﰲ", "colors", " colors"))
	paletteKeys = append(paletteKeys, NewShortcut('j', "ﰬ", "tools", " tools"))
	paletteKeys = append(paletteKeys, NewShortcut('k', "ﰵ", "tools", " tools"))
	paletteKeys = append(paletteKeys, NewShortcut('+', "螺", "add", "螺add color"))
	paletteKeys = append(paletteKeys, NewShortcut('-', "羅", "del", "羅delete color"))
	paletteKeys = append(paletteKeys, NewShortcut('i', "", "info", " info"))
	paletteKeys = append(paletteKeys, NewShortcut('e', "", "edit", " edit color"))
	paletteKeys = append(paletteKeys, NewShortcut('*', "", "randomize", " randomize"))
	table[paletteTable] = paletteKeys
	// 
	editKeys = append(editKeys, NewShortcut('h', "ﰯ", "channel", " channel"))
	editKeys = append(editKeys, NewShortcut('j', "", "decrease", " decrease"))
	editKeys = append(editKeys, NewShortcut('k', "", "increase", " increase"))
	editKeys = append(editKeys, NewShortcut('l', "ﰲ", "channel", " channel"))
	editKeys = append(editKeys, NewShortcut('>', "", "incr", " increments"))
	editKeys = append(editKeys, NewShortcut('<', "", "decr", " increments"))
	table[editTable] = editKeys
	theme := theme.GetTheme()
	_ = theme
}

//     ﮜ         
// ﰯ ﰬ ﰵ ﰲ                       
//                            
//                            
//                            
//               
// 
var testTxt string = `[-:-:-][:#ababab:][black::][:#744241:r]|  |[:#8b4f4f:]| [black::r][#9e5a59::b][-:-:-]|[black::][-:-:-][-:-:-] [::r][#8b4f4f:-:-]|[#744241::] [black::][-:-:-][-:-:-] [::r][#744241::] [black::b][-:-:-][-:-:-][:#744241:] [black::r][#8b4f4f:-:-][-:-:-]|[#744241::] [black::][-:-:-][-:-:-] [::r][#9e5a59::b]|[black::][-:-:-][-:-:-] [1m[:#8b4f4f:]|[:#744241:]|  |
`
func (sc shortcut) String() string {
  return fmt.Sprintf("[gray::]🮤[-:-:-] [blue::db]%s[-:-:-] [yellow::d]%s[-:-:-] [green::b]%c[-:-:-] [gray::]🮥[-:-:-]", sc.icon, sc.name, sc.key)
}

func (s *HelpBar) SetTable(t string) {
	if !maputil.HasKey(table, t) {
		return
	}
	helpkeys := make([]string, 0)
	for _, v := range table[t] {
		helpstr := v.String()
		helpkeys = append(helpkeys, helpstr)
	}
  s.helpTextView.SetDynamicColors(true)
  s.helpTextView.SetTextColor(tcell.GetColor("#010101"))
  // 🭲🭱🭱🭱🭰🭱🭲🭴🭵🭳🬦
	helptxt := strings.Join(helpkeys, "   ")
  s.helpTextView.Clear()
	s.helpTextView.SetText(helptxt)
}

func coder() string {
  rand.Seed(time.Now().UnixNano())
  strs := `   [#744241::]|  |[#9e5a59::b]|[:black:] [black:#9e5a59:]|[:black:] [black:#9e5a59:]|[black:#744241:] [::-]`
	// paletteKeys = append(paletteKeys, NewShortcut('b', "羅", coder(), "羅delete color"))
	// paletteKeys = append(paletteKeys, NewShortcut('Q', "", "quit", " quit"))
	// keys = append(keys, NewShortcut('<spc>', "", "random", " randomize colors"))
	// keys = append(keys, NewShortcut('e', "", "edit", " edit color"))
  // fmtstr := []string{"[#744241::] | [-:-:-]","[#744241::] | [-:-:-]","[#744241::r]|[-:-:-]","[#744241::r] |[-:-:-]","[::r] |[-:-:-]"}
  // str := make([]string, 5)
  // str := append(lo.Shuffle[string](fmtstr),lo.Shuffle[string](fmtstr)...)

  // return strings.Join(str, " ")
  return strs
}

func (s *HelpBar) Init() {
  s.Flex.AddItem(s.helpTextView, 0, 1, false)
  s.Flex.AddItem(nil, 0, 1, false)
	s.helpTextView.SetRegions(true).SetBorder(false).SetBorderPadding(0, 0, 0, 0)
	s.helpTextView.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	s.SetTable(mainTable)
	// , "/ ", "mixer"

	// 
	// 
	// grays := []string{"#5e484c","#685459","#736166","#948790","#a0959e","#7e6d74","#897a82","#aba2ac","#b7b0bb"}
	// keys :=
	// s.SetText(fmt.Sprintf(`[red:-:b]%s[-:-:-] | [%s:-:-]%s[-:-:-]  [red:-:b]%s[-:-:-] | [yellow:-:-]%s[-:-:-]`, " l,","  colors", ",j k,", " tools")).SetTextAlign(tview.AlignCenter)
	// s.UpdateRegion("editor")
}

func (s *HelpBar) UpdateRegion(r string) {
	s.helpTextView.Highlight(r)
}

// vim: ts=2 sw=2 et ft=go
