package coolor

import (
	"fmt"
	"os"
	"sync"
	"text/template"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/gookit/color"
	// "github.com/lucasb-eyer/go-colorful"
	"golang.org/x/term"
)

type SelectionEvent struct {
  *tcell.EventTime
  color *CoolorColor
  idx int8
}

type OnCoolorColorSelected interface {
  Selected(ev SelectionEvent) bool
}

type SelectedEventHandler interface {
  tcell.EventHandler
  OnCoolorColorSelected
}

type EventHandlers []tcell.EventHandler

type CoolorColor struct {
	*tview.Box
	idx                     int8
	color                   *tcell.Color
	static, locked, selected, dirty bool
	pallette                *CoolorPalette
	l                       *sync.RWMutex
  handlers      map[string]EventHandlers
}

func NewCoolorColor(col string) *CoolorColor {
	cc := NewDefaultCoolorColor()
	cc.SetColorCss(col)
	return cc
}
// 	
func NewStaticCoolorColor(col string) *CoolorColor {
	cc := NewDefaultCoolorColor()
	cc.SetColorCss(col)
  cc.static = true
	return cc
}
// 	
 func (cc *CoolorColor) RGBA() (r, g, b, a uint32) {
   ri, gi, bi := cc.color.RGB()
   fmt.Println(ri,gi,bi)
   return uint32(ri),uint32(gi),uint32(bi),0xffff
 }

 // func (cc *CoolorColor) Noire() noire.Color {
 //   pcol := noire.NewRGB(float64(R),float64(G),float64(B))
 //   return pcol
 // }

 func (cc *CoolorColor) HSL() (float64,float64,float64) {
   // R,G,B,A := cc.RGBA()
   // cf, _ := colorful.MakeColor(cc)
   hsla := NewHSLA(cc)

   return hsla.H, hsla.S, hsla.L
 }

 func (cc *CoolorColor) GetCC() *CoolorColor {
   return cc.Clone()
 }
 func (cc *CoolorColor) Clone() *CoolorColor {
   ccc := NewDefaultCoolorColor()
   c := tcell.GetColor(cc.Html())
   ccc.SetColorInt(c.Hex())
   fmt.Println(ccc.TerminalPreview())
   // ccc.color = &c
   ccc.locked = false
   ccc.selected = cc.selected
   ccc.dirty = false
   ccc.pallette = cc.pallette
   return ccc
 }

func NewCoolorBox() *tview.Box {
	return tview.NewBox()
}

func NewDefaultCoolorColor() *CoolorColor {
	box := NewCoolorBox()
	cc := &CoolorColor{
		// Box:      tview.NewBox(),
		idx:      0,
    static: false,
		color:    nil,
		locked:   false,
		dirty:    false,
		selected: false,
		pallette: nil,
		l:        &sync.RWMutex{},
    handlers: make(map[string]EventHandlers),
	}
	cc.Box = box

	box.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
    if cc.static {
      return x,y,width,height
    }
		// Draw a horizontal line across the middle of the box.
		centerY := y + height/2
    lowerCenterY := centerY + centerY / 2
		// stripTop := centerY - 2
		// stripBottom := centerY + 1
		// for ypos := stripTop; ypos < stripBottom; ypos++ {
		// 	tview.Print(
		// 		screen,
		// 		fmt.Sprintf("[red:black:b]%s", strings.Repeat(" ", width-2)),
		// 		x+1,
		// 		ypos,
		// 		width-2,
		// 		tview.AlignCenter,
		// 		tcell.ColorRed,
		// 	)
		// }
		for cx := x + 1; cx < x+width-1; cx++ {
		  screen.SetContent(cx, centerY + (height/ 3), tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(*cc.color))
		}


		//"" Write some text along the horizontal line.
    
    status_tpl := MakeTemplate("color_status", `
      {{define "locked"}}{{- if locked -}}{{- end -}}{{- end -}}
      {{define "selected"}}{{- if selected -}}{{- end -}}{{- end -}}
    `, template.FuncMap{
      "locked": cc.GetLocked,
      "selected": cc.GetSelected,
      "dirty": cc.GetDirty,
      "css": cc.GetColor,
    })
    // sel := &strings.Builder{}
    // lock := &strings.Builder{}
    sel := status_tpl(`{{- template "selected" . -}}`, cc)
    lock := status_tpl(`{{- template "locked" . -}}`, cc)
    // ok = template.Must(status_tpl.Clone()).Parse(`{{- template "selected" . -}}`).Execute(sel, cc)
    // if ok != nil {
    //   fmt.Println(fmt.Errorf("%s", ok))
    // }
    // txt := output.String()
    txtColor := cc.GetFgColor()
		tview.Print(screen, sel, x+1, centerY, width-2, tview.AlignCenter, txtColor)
		tview.Print(screen, lock, x+1, lowerCenterY, width-2, tview.AlignCenter, txtColor)

		// Space for other content.
		return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
	})
	return cc
}

func NewIntCoolorColor(h int32) *CoolorColor {
	cc := NewDefaultCoolorColor()
	cc.SetColorInt(h)
	return cc
}

func NewRandomCoolorColor() *CoolorColor {
	c := MakeRandomColor()
  
	// cc := NewDefaultCoolorColor()
	// cc.SetColor()
	return NewIntCoolorColor(c.Hex())
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
	cc.l.RLock()
	defer cc.l.RUnlock()
	return cc.locked
}

func (c *CoolorColor) Random() bool {
	c.l.Lock()
	defer c.l.Unlock()
	if c.locked {
		return false
	}
	c.dirty = true
	c.color = MakeRandomColor()
	c.Box.SetBackgroundColor(*c.color)
	return true
}

func (c *CoolorColor) Remove() {
	c.pallette.RemoveItem(c)
}

func (c *CoolorColor) SetColor(col *tcell.Color) {
	c.l.Lock()
	defer c.l.Unlock()
  hex := fmt.Sprintf("#%06x", col.Hex())
  // fmt.Println(hex)
  colo := tcell.GetColor(hex)
	c.color = &colo
	c.Box.SetBackgroundColor(*c.color)
}

func (c *CoolorColor) GetColor() string {
	c.l.RLock()
	defer c.l.RUnlock()
  return fmt.Sprintf("#%06x", c.color.Hex())
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

func (cc *CoolorColor) SpawnSelectionEvent(t string, ev tcell.Event) bool {
	cc.l.Lock()
	defer cc.l.Unlock()
  if len(cc.handlers[t]) > 0 {
    for _, v := range cc.handlers[t] {
      if v != nil {
        eh, ok := v.(tcell.EventHandler)
        if !ok {
          panic(ok)
        }
        eh.HandleEvent(ev)
      }
    }
  }
  return true
}

func (cc *CoolorColor) AddEventHandler(t string, h *tcell.EventHandler) {
	cc.l.Lock()
	defer cc.l.Unlock()

  if cc.handlers[t] == nil {
    cc.handlers[t] = make(EventHandlers, 0)
  }
  cc.handlers[t] = append(cc.handlers[t], *h)
}

func (cc *CoolorColor) SetSelected(s bool) {
	cc.l.Lock()
	defer cc.l.Unlock()

	cc.selected = s
	cc.updateStyle()
}

func (cc *CoolorColor) GetFgColor() tcell.Color {
	r, g, b := cc.color.RGB()
	if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
		return tcell.ColorBlack
	}
	return tcell.ColorWhite
}

func (cc *CoolorColor) Html() string {
	return fmt.Sprintf("#%06x", cc.color.Hex())
}

func (cc *CoolorColor) TerminalPreview() string {
	r, g, b := cc.color.RGB()
	br, bg, bb := getFGColor(*cc.color).RGB()
return fmt.Sprintf(
			"\033[48;2;%d;%d;%d;38;2;%d;%d;%dm#%06x\033[0m",
			r, g, b, br, bg, bb, cc.color.Hex(),
		)
}
func (cc *CoolorColor) String() string {
	r, g, b := cc.color.RGB()
	br, bg, bb := getFGColor(*cc.color).RGB()
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return fmt.Sprintf(
			"\033[48;2;%d;%d;%d;38;2;%d;%d;%dm #%06x \033[0m\n",
			r, g, b, br, bg, bb, cc.color.Hex(),
		)
	}
	return fmt.Sprintf(" #%06x \n", cc.color.Hex())
}

func (cc *CoolorColor) ToggleLocked() {
	cc.SetLocked(!cc.GetLocked())
}

func (cc *CoolorColor) updateStyle() {
  if cc.static {
    return
  }
	if cc.selected {
		inverse := getFGColor(*cc.color)
		cc.Box.
			SetBorder(true).
			SetBorderAttributes(tcell.AttrBold).
			SetBorderPadding(4, 0, 0, 0).
			SetBorderColor(inverse).
			SetTitleColor(inverse)
		cc.Box.Focus(nil)
	} else {
		cc.Box.SetBorderPadding(0, 0, 0, 0)
		cc.Box.SetBorder(false)
		cc.Box.Blur()
	}
}

// vim: ts=2 sw=2 et ft=go
