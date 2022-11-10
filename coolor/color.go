package coolor

import (
	// "bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	// "text/template"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/samber/lo"
	_ "github.com/samber/lo"
	"golang.org/x/term"
)

type CoolorColor struct {
	*tview.Box `msgpack:"-" clover:"-,omitempty"`
	handlers   map[string]EventHandlers
	Color      *tcell.Color `msgpack:"-" clover:",omitempty"`
	l          *sync.RWMutex
	pallette   *CoolorPaletteMainView
	*Tagged    `msgpack:"-" clover:"-,omitempty"`
	name       string
	infoline   string
	Favorite   bool `clover:"favorite"`
	static     bool
	selected   bool
	dirty      bool
	plain      bool
	centered   bool
	idx        int8
	valid      bool
	locked     bool
}

// UnmarshalJSON implements json.Unmarshaler
func (cc *CoolorColor) UnmarshalJSON(b []byte) error {
  // fmt.Println(b)
  m := make(map[string]interface{})
  cc = NewDefaultCoolorColor()
  err := json.Unmarshal(b, &m)
  if err != nil {
    return err
  }
// fmt.Printf("#%06x", int32(m["Color"].(float64)))
  cc.SetColorCss(fmt.Sprintf("#%06x", int32(m["Color"].(float64))))
  cc.Favorite = m["Favorite"].(bool)
  return nil
}

func NewCoolorColor(col string) *CoolorColor {
	cc := NewDefaultCoolorColor()
	cc.SetColorCss(col)
	return cc
}

func NewStaticCenteredCoolorColor(col string) *CoolorColor {
	cc := NewDefaultCoolorColor()
	cc.SetColorCss(col)
	cc.static = true
	cc.centered = true
	return cc
}

func NewStaticCoolorColor(col string) *CoolorColor {
	cc := NewDefaultCoolorColor()
	cc.SetColorCss(col)
	cc.static = true
	return cc
}

func NewCoolorBox() *tview.Box {
	return tview.NewBox()
}

func (cc *CoolorColor) DrawFunc() func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	return func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		centerY := y + height/2
		lowerCenterY := centerY + height/3
		txtColor := cc.GetFgColor()
		// marker := fmt.Sprintf
		markers := []string{"", ""}
		// fill := lo.Repeat(width, " ")
		spcs := strings.Join(make([]string, width-1), " ")
		needles := fmt.Sprintf("%s%s%s", markers[0], spcs, markers[1])

		if cc.static || cc.centered {
			if cc.infoline == " " {
				tview.Print(
					screen,
					needles,
					x,
					centerY,
					width,
					tview.AlignCenter,
					tcell.ColorDarkRed,
				)
			}
			for cx := x + 1; cx < x+width-2; cx++ {
				col := tcell.StyleDefault.Foreground(tcell.ColorWhite).
					Background(*cc.Color)
				lw := tview.BoxDrawingsLightHorizontal
				if !cc.valid {
					lw = tview.BoxDrawingsHeavyHorizontal
					col = tcell.StyleDefault.Foreground(tcell.ColorRed).
						Background(*cc.Color)
				}
				// screen.SetContent(cx, centerY, lw, nil, col)
				if !cc.plain {
					screen.SetContent(cx, centerY, lw, nil, col)
				}
			}
			if cc.infoline != "" {
				tview.Print(
					screen,
					cc.infoline,
					x+1,
					centerY,
					width-1,
					tview.AlignCenter,
					txtColor,
				)
			}
		} else {
			// cc.InRect(lowerCenterY +2)
			if width-2 >= 8 && lowerCenterY+1 <= y+height-2 {
				// for cx := x + 1; cx < x+width-1; cx++ {
				// tcol, _ := MakeColor(cc)
				// r, g, b := tcol.RGB255()
				//        yiq := rgbToYIQ(uint(r), uint(g), uint(b))
				// fmt.Println()
				tview.Print(screen, fmt.Sprintf("[#%06x:-:b]%s[-:-:-]", cc.GetFgColor().Hex(), strings.ToUpper(fmt.Sprintf("%06x", cc.Color.Hex()))), x+1, lowerCenterY+2, width-2, tview.AlignCenter, txtColor)
				// tview.Print(screen, fmt.Sprintf("[#%06x:-:b]%0.2f[-:-:-]", cc.GetFgColor().Hex(), yiq), x+1, lowerCenterY+3, width-2, tview.AlignCenter, txtColor)
				// screen.SetContent(cx, lowerCenterY + 4, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(*cc.color))
				// }
			}
			// status_tpl := MakeTemplate("color_status", `
   //    {{define "locked"}}{{- if not locked }}  {{ else }}  {{ end -}}{{- end -}}
   //    {{define "selected"}}{{- if selected }}   {{ else }}   {{ end -}}{{- end -}}
   //  `, template.FuncMap{
			// 	"locked":   cc.GetLocked,
			// 	"selected": cc.GetSelected,
			// 	"dirty":    cc.GetDirty,
			// 	"css":      cc.GetColor,
			// })
			// sel := status_tpl(`{{- template "selected" . -}}`, cc)
			// lock := status_tpl(`{{- template "locked" . -}}`, cc)
			// tview.Print(screen, sel, x+1, (lowerCenterY - centerY) / 2, width-2, tview.AlignCenter, txtColor)
      lock := IfElseStr(!cc.locked, " ", " ")
			tview.Print(screen, lock, x+1, lowerCenterY, width-2, tview.AlignCenter, txtColor)
		}
    if cc.selected && (!cc.plain && !cc.static){
      if cc.pallette == nil || cc.pallette.menu == nil || cc.pallette.menu.list == nil {
        return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
      }
      cc.GetMenuPosition(cc.pallette.menu.list)
      cc.pallette.menu.list.Draw(screen)
    }
    // if cc.selected {
    //   cc.GetMenuPosition()
    //   cc.pallette.menu.list.SetRect(x int, y int, width int, height int)
    // }

		// Space for other content.
		return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
		// return x + 1, y + 1, width - 2, height
	}
}

func NewDefaultCoolorColor() *CoolorColor {
	box := NewCoolorBox()
	cc := &CoolorColor{
		Box:      box,
		handlers: make(map[string]EventHandlers),
		Color:    nil,
		l:        &sync.RWMutex{},
		pallette: nil,
		Tagged:   NewTaggable(&Base16Tags),
		name:     "",
		infoline: "",
		static:   false,
		selected: false,
		dirty:    false,
		plain:    false,
		centered: false,
		idx:      0,
		valid:    false,
		locked:   false,
	}
	cc.Tagged.Item = &cc
	box.SetDrawFunc(cc.DrawFunc())
	return cc
}

func NewIntCoolorColor(h int32) *CoolorColor {
	cc := NewDefaultCoolorColor()
	cc.SetColorInt(h)
	return cc
}

func NewRandomCoolorColor() *CoolorColor {
	c := MakeRandomColor()
	return NewIntCoolorColor(c.Hex())
}

func (cc *CoolorColor) RGBA() (r, g, b, a uint32) {
	// ri, gi, bi := cc.color.RGB()
	// fmt.Println(ri,gi,bi)
	ccn, _ := Hex(cc.Html())
	// ccn.RGBA()
	return ccn.RGBA()
	// return uint32(ri),uint32(gi),uint32(bi),0xffff
}

// func (cc *CoolorColor) Noire() noire.Color {
//   pcol := noire.NewRGB(float64(R),float64(G),float64(B))
//   return pcol
// }

func (cc *CoolorColor) HSL() (float64, float64, float64) {
	hsla, _ := MakeColor(cc)
	return hsla.Hsl()
}

func (cc *CoolorColor) Unstatic() *CoolorColor {
	cc.static = false
	cc.centered = false
	return cc
}

func (cc *CoolorColor) Coolor() *Coolor {
	if cc == nil || cc.Color == nil {

		return &Coolor{
			Color: 0,
		}
	}
	return &Coolor{
		Color: *cc.Color,
	}
}

func (cc *CoolorColor) GetColorable() *Color {
	c := MakeColorFromTcell(*cc.Color)
	return &c
}

func (cc *CoolorColor) U64() uint64 {
	return uint64(cc.Color.Hex())
}

func (cc *CoolorColor) GetCC() *CoolorColor {
	ccn := NewStaticCoolorColor(cc.Html())
	return ccn
}

func (cc *CoolorColor) Clone() *CoolorColor {
	ccc := NewDefaultCoolorColor()
	c := tcell.GetColor(cc.Html())
	ccc.SetColorInt(c.Hex())
	ccc.static = cc.static
	ccc.centered = cc.centered
	ccc.locked = false
	ccc.selected = cc.selected
	ccc.dirty = false
	ccc.pallette = cc.pallette
	return ccc
}

func (cc *CoolorColor) GetSelected() bool {
	cc.l.RLock()
	defer cc.l.RUnlock()
	return cc.selected
}

func (cc *CoolorColor) GetDirty() bool {
	cc.l.RLock()
	defer cc.l.RUnlock()
	return cc.dirty
}

func (cc *CoolorColor) GetLocked() bool {
	// dump.P(cc.Html(), cc.locked, cc.selected, cc.plain, cc.static)
	cc.l.RLock()
	defer cc.l.RUnlock()
	return cc.locked
}

func (c *CoolorColor) Random() bool {
	if c.locked {
		return false
	}
	c.dirty = true
	col := MakeRandomColor()
	c.SetColor(col)
	SeentColor("setcolor", c, c)
	return true
}

func (c *CoolorColor) Remove() {
	c.pallette.RemoveItem(c)
}

func (c *CoolorColor) SetName(n string) {
	c.l.Lock()
	defer c.l.Unlock()
	c.name = n
}

func (c *CoolorColor) GetName() string {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.name
}

func (c *CoolorColor) SetColor(col *tcell.Color) {
	c.l.Lock()
	defer c.l.Unlock()
	hex := fmt.Sprintf("#%06x", col.Hex())
	colo := tcell.GetColor(hex)
	c.Color = &colo
	c.SetBackgroundColor(*c.Color)
}

func (c *CoolorColor) GetColor() string {
	c.l.RLock()
	defer c.l.RUnlock()
	return fmt.Sprintf("#%06x", c.Color.Hex())
}

func (c *CoolorColor) SetInfoLine(str string, valid bool) {
	c.l.Lock()
	c.infoline = str
	c.valid = valid
	c.l.Unlock()
	c.updateStyle()
}

func (c *CoolorColor) SetColorCss(str string) {
	col := tcell.GetColor(str)
	c.SetColor(&col)
}

func (c *CoolorColor) SetColorInt(h int32) {
	c.SetColorCss(fmt.Sprintf("#%06x", h))
}

func (cc *CoolorColor) SetLocked(s bool) {
	cc.l.Lock()
	defer cc.l.Unlock()
	cc.locked = s
	cc.updateStyle()
}

// func (cc *CoolorColor) SpawnSelectionEvent(t string, ev tcell.Event) bool {
// 	cc.l.Lock()
// 	defer cc.l.Unlock()
// 	if len(cc.handlers[t]) > 0 {
// 		for _, v := range cc.handlers[t] {
// 			if v != nil {
// 				// eh, ok := v.(tcell.EventHandler)
// 				// if !ok {
// 				// 	panic(ok)
// 				// }
// 				v.HandleEvent(ev)
// 			}
// 		}
// 	}
// 	return true
// }
//
// func (cc *CoolorColor) AddEventHandler(t string, h *tcell.EventHandler) {
// 	cc.l.Lock()
// 	defer cc.l.Unlock()
//
// 	if cc.handlers[t] == nil {
// 		cc.handlers[t] = make(EventHandlers, 0)
// 	}
// 	cc.handlers[t] = append(cc.handlers[t], *h)
// }

func (cc *CoolorColor) SetCentered(s bool) {
	cc.l.Lock()
	defer cc.l.Unlock()

	cc.centered = s
	cc.updateStyle()
}

func (cc *CoolorColor) SetPlain(s bool) {
	cc.l.Lock()
	defer cc.l.Unlock()

	cc.plain = s
	cc.updateStyle()
}

func (cc *CoolorColor) SetStatic(s bool) {
	cc.l.Lock()
	defer cc.l.Unlock()

	cc.static = s
	cc.updateStyle()
}

func (cc *CoolorColor) SetSelected(s bool) {
	cc.l.Lock()
	cc.selected = s
	cc.updateStyle()
	cc.l.Unlock()
}

func (cc *CoolorColor) GetFgColorFade(a float64) tcell.Color {
	tcol, _ := MakeColor(cc)
	ncc := tcol.BlendLuvLCh(MakeColorFromTcell(cc.GetFgColorShade()), a)
	return *ncc.GetCC().Color
}

func (cc *CoolorColor) GetFgColorShade() tcell.Color {
	tcol, _ := MakeColor(cc)
	r, g, b := tcol.RGB255()
	if rgbToYIQ(uint(r), uint(g), uint(b)) >= 128 {
		return *NewCoolorColor("#505050").Color
	} else {
		return *NewCoolorColor("#b0b0b0").Color
	}
}

func (cc *CoolorColor) GetFgColor() tcell.Color {
	tcol, _ := MakeColor(cc)
	r, g, b := tcol.RGB255()
	if rgbToYIQ(uint(r), uint(g), uint(b)) >= 128 {
		return *NewCoolorColor("#101010").Color
	} else {
		return *NewCoolorColor("#f0f0f0").Color
	}
}

// func (cc *CoolorColor) GetFgColor() tcell.Color {
// 	c, ok := MakeColor(cc)
// 	// dump.P(cc.TerminalPreview())
// 	if ok {
// 		r, g, b := c.LinearRgb()
// 		if (255*float64(r)*0.299 + 255*float64(g)*0.587 + 255*float64(b)*0.114) > 150 {
// 			// if (255*float64(r)*0.2926 + 255*float64(g)*0.5152 + 255*float64(b)*0.1722) > 150 {
// 			// if (float64(r)*0.2126 + float64(g)*0.7152 + float64(b)*0.0722) > 140 {
// 			return tcell.ColorBlack
// 		}
// 		return tcell.ColorWhite
// 	}
// 	return tcell.ColorWhite
// 	// r, g, b := cc.RGB
// 	// if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
// }

func (cc *CoolorColor) Html() string {
	if cc == nil || cc.Color == nil {
		return "#000000"
	}
	return strings.ToUpper(fmt.Sprintf("#%06x", cc.Color.Hex()))
}

func (cc *CoolorColor) TVCSSString(spaces bool) string {
	if cc == nil || cc.Color == nil {
		return "#000000"
	}
	space := " "
	if !spaces {
		space = ""
	}
	return strings.ToUpper(fmt.Sprintf(
		"[#%06x:#%06x:-]%s#%06x%s[-:-:-]",
		cc.GetFgColor().Hex(),
		cc.Color.Hex(),
		space,
		cc.Color.Hex(),
		space,
	))
}

func (cc *CoolorColor) TVPreview() string {
	return cc.TVCSSString(true)
}

func (cc *CoolorColor) TerminalPreview() string {
	r, g, b := cc.Color.RGB()
	br, bg, bb := cc.GetFgColor().RGB()
	return fmt.Sprintf(
		"\033[48;2;%d;%d;%d;38;2;%d;%d;%dm #%06x \033[0m",
		r, g, b, br, bg, bb, cc.Color.Hex(),
	)
}

func (cc *CoolorColor) String() string {
	r, g, b := cc.Color.RGB()
	br, bg, bb := cc.GetFgColor().RGB()
	// br, bg, bb := getFGColor(*cc.color).RGB()
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return fmt.Sprintf(
			"\033[48;2;%d;%d;%d;38;2;%d;%d;%dm #%06x \033[0m\n",
			r, g, b, br, bg, bb, cc.Color.Hex(),
		)
	}
	return fmt.Sprintf(" #%06x \n", cc.Color.Hex())
}

func (cc *CoolorColor) ToggleLocked() {
	cc.SetLocked(!cc.GetLocked())
}

func (cc *CoolorColor) updateStyle() {
	// MainC.app.QueueUpdateDraw(func() {
	// dump.P(cc.Html(), cc.selected, cc.plain, cc.static)
	if cc.plain {
		cc.SetBorderPadding(0, 0, 0, 0)
		cc.SetBorder(false)
		cc.Blur()
		return
	}
	if cc.selected || cc.centered {
		// cc.GetFgColorShade()
		inverse := cc.GetFgColor()
		// fmt.Println(MakeColorFromTcell(inverse).GetCC().TerminalPreview())
		cc.Box.
			SetBorderFocusColor(inverse).
			SetBorder(true).
			SetBorderPadding(0, 0, 0, 0).
			SetBorderColor(tcell.GetColor("#101010"))
		cc.Focus(nil)
	} else {
		cc.SetBorderPadding(0, 0, 0, 0)
		cc.SetBorder(false)
		cc.Blur()
	}
}

func (cc *CoolorColor) Draw(screen tcell.Screen) {
  // cc.GetMenuPosition(cc.pallette.menu)
	cc.DrawForSubclass(screen, cc)
	if cc == nil || cc.pallette == nil || cc.pallette.menu == nil || !cc.selected {
		return
	}
	// mw = (((imw * 3) / 3))
	// centerY := cy + ((ch/2 + (ch / 3) - mh/2) - (mh / 2))
	// centerY := cy + ((ch/2 + (ch / 2)) - (ch / 3))
	// shouldReturn := 
	// if shouldReturn {
	// 	return
	// }
	// cc.pallette.menu.Draw(screen)
	// cc.pallette.menu.updateState()
}

func (cc *CoolorColor) GetMenuPosition(p tview.Primitive) {
	if cc == nil || cc.pallette == nil || cc.pallette.menu == nil {
    return
		// return cc.GetRect()
	}
  // dump.P(cc.GetRect())
  // dump.P(cc.GetInnerRect())
	x, y, w, h := cc.GetRect()
	_, _, _, _ = x, y, w, h
	cx, cy, cw, ch := cc.GetInnerRect()
	cw = ((cw * 2) / 2) - 1
	cc.SetRect(cx, cy, cw, ch)
	_, _, _, _ = cx, cy, cw, ch
	cx, cy, cw, ch = cc.GetInnerRect()
	_, _, _, _ = cx, cy, cw, ch
	mx, my, _, _ := cc.pallette.menu.GetRect()
	imx, imy, imw, imh := cc.pallette.menu.GetInnerRect()
	// mic := cc.pallette.menu.GetItemCount()
  mic := len(cc.pallette.menu.GetListItems())
  if mic == 0 {
    mic = 1
  }

	mw := 1
	itemH := mic * 3
	mh := (itemH) + 2
	endY := (y + h)
	centerX := cx + ((cw / 2) - (mw / 2))
	centerY := y + (mh / mic)
    // dump.P(centerX, centerY, itemH, mh)

	if centerY+mh > cy+ch-5 {
		mh = (endY) - centerY
	}
	_, _, _, _ = mx, my, mw, mh
	_, _, _, _ = imx, imy, imw, imh

  if p != nil {
	p.SetRect(centerX, centerY, mw, mh)

  }
}

func (cc CoolorColors) Strings() []string {
	cssStrings := lo.Map[*CoolorColor, string](
		cc,
		func(cc *CoolorColor, i int) string {
			return cc.Html()
		},
	)
	return cssStrings
}
func (cc *CoolorColor) GetRef() interface{} {
	return cc
}

// vim: ts=2 sw=2 et ft=go
