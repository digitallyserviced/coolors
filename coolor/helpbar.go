package coolor

import (
	// "fmt"

	"fmt"
	"strings"

	"github.com/digitallyserviced/tview"
	"github.com/gookit/color"
	"github.com/gookit/goutil/maputil"
)

type HelpBar struct {
	*tview.TextView
	app *tview.Application
}

func NewHelpBar(app *tview.Application) *HelpBar {
	hb := &HelpBar{
		TextView: tview.NewTextView(),
		app:      app,
	}
	hb.Init()
	return hb
}

type shortcut struct {
	key  rune
	mods string
	icon string
	name string
	desc string
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

func (sc shortcut) String() string {
	return tview.TranslateANSI(color.Render(fmt.Sprintf(" <sckey>%c</> <scicon>%s</> <scname>%s</> ", sc.key, sc.icon, sc.name)))
}

//     ﮜ         
// ﰯ ﰬ ﰵ ﰲ                       
//                            
//                            
//                            
//               
// 
const (
	mainTable    = "palette"
	paletteTable = "palette"
	editTable    = "editor"
)

type (
	keyMap    []*shortcut
	keyTables map[string]keyMap
)

var (
	paletteKeys keyMap
	editKeys    keyMap
	table       keyTables
)

func init() {
	table = make(keyTables)
	paletteKeys = make(keyMap, 0)
	editKeys = make(keyMap, 0)

	paletteKeys = append(paletteKeys, NewShortcut('h', "ﰯ", "colors", " colors"))
	paletteKeys = append(paletteKeys, NewShortcut('l', "ﰲ", "colors", " colors"))
	paletteKeys = append(paletteKeys, NewShortcut('j', "ﰬ", "tools", " tools"))
	paletteKeys = append(paletteKeys, NewShortcut('k', "ﰵ", "tools", " tools"))
	paletteKeys = append(paletteKeys, NewShortcut('+', "螺", "add", "螺add color"))
	paletteKeys = append(paletteKeys, NewShortcut('-', "羅", "del", "羅delete color"))
	paletteKeys = append(paletteKeys, NewShortcut('i', "", "info", " info"))
	paletteKeys = append(paletteKeys, NewShortcut('e', "", "edit", " edit color"))
	paletteKeys = append(paletteKeys, NewShortcut('*', "", "randomize", " randomize"))
	// paletteKeys = append(paletteKeys, NewShortcut('Q', "", "quit", " quit"))
	table[paletteTable] = paletteKeys
	// 
	editKeys = append(editKeys, NewShortcut('h', "ﰯ", "channel", " channel"))
	editKeys = append(editKeys, NewShortcut('j', "", "decrease", " decrease"))
	editKeys = append(editKeys, NewShortcut('k', "", "increase", " increase"))
	editKeys = append(editKeys, NewShortcut('l', "ﰲ", "channel", " channel"))
	editKeys = append(editKeys, NewShortcut('>', "", "incr", " increments"))
	editKeys = append(editKeys, NewShortcut('<', "", "decr", " increments"))
	table[editTable] = editKeys
	// keys = append(keys, NewShortcut('<spc>', "", "random", " randomize colors"))
	// keys = append(keys, NewShortcut('e', "", "edit", " edit color"))
	theme := GetTheme()
	_ = theme
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
	helptxt := strings.Join(helpkeys, " | ")
  s.Clear()
	s.SetText(helptxt)
}

func (s *HelpBar) Init() {
	s.SetRegions(true).SetBorder(false).SetBorderPadding(0, 0, 0, 0)
	s.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
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
	s.Highlight(r)
}

// vim: ts=2 sw=2 et ft=go
