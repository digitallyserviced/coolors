package coolor

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	_ "github.com/samber/lo"

	// "github.com/gookit/color"
	// "github.com/lucasb-eyer/go-colorful"
	"golang.org/x/term"
)

type SelectionEvent struct {
	*tcell.EventTime
	color *CoolorColor
	idx   int8
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
	idx                                                     int8
	color                                                   *tcell.Color
	name                                                    string
	plain, valid, centered, static, locked, selected, dirty bool
	pallette                                                *CoolorPalette
	l                                                       *sync.RWMutex
	handlers                                                map[string]EventHandlers
	infoline                                                string
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

//
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

func NewCoolorBox() *tview.Box {
	return tview.NewBox()
}

type HookDrawInfo struct {
	x, y, width, height   int
	centerY, lowerCenterY int
}

type HookDrawFunctions struct {
	Target *tview.Primitive
	Chain  DrawFunctionChain
	Wrap   DrawFunction
	// Before func(HookDrawInfo)(int, int, int, int)
	// Draw func(HookDrawInfo)(int, int, int, int)
	// After func(HookDrawInfo)(int, int, int, int)
}

type (
	DrawFunction      func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int)
	DrawFunctionChain []*DrawFunction
)

func DrawFunctionDispatcher(p *tview.Primitive, dfc DrawFunctionChain) DrawFunction {
	return func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		for _, v := range dfc {
			if v != nil {
			}
		}
		return x, y, width, height
	}
}

func (cc *CoolorColor) DrawHook(df *DrawFunction) {
}

func (hdf *HookDrawFunctions) CoolorColorStatusText(p tview.Primitive, screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	cc, ok := p.(*CoolorColor)
	if !ok {
		return x, y, width, height
	}
	centerY := y + height/2
	lowerCenterY := centerY + centerY/2
	for cx := x + 1; cx < x+width-1; cx++ {
		screen.SetContent(cx, centerY+(height/3), tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(*cc.color))
	}

	status_tpl := MakeTemplate("color_status", `
      {{define "locked"}}{{- if locked -}}  {{- else -}}  {{- end -}}{{- end -}}
      {{define "selected"}}{{- if selected -}}   {{- end -}}{{- end -}}
    `, template.FuncMap{
		"locked":   cc.GetLocked,
		"selected": cc.GetSelected,
		"dirty":    cc.GetDirty,
		"css":      cc.GetColor,
	})
	sel := status_tpl(`{{- template "selected" . -}}`, cc)
	lock := status_tpl(`{{- template "locked" . -}}`, cc)
	txtColor := cc.GetFgColor()
	tview.Print(screen, sel, x+1, centerY, width-2, tview.AlignCenter, txtColor)
	tview.Print(screen, lock, x+1, lowerCenterY, width-2, tview.AlignCenter, txtColor)

	return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
}

func CenteredStrikeText() {
}

func NewDefaultCoolorColor() *CoolorColor {
	box := NewCoolorBox()
	cc := &CoolorColor{
		Box:      box,
		idx:      0,
		color:    nil,
		name:     "",
		plain:    false,
		valid:    false,
		centered: false,
		static:   false,
		locked:   false,
		selected: false,
		dirty:    false,
		pallette: nil,
		l:        &sync.RWMutex{},
		handlers: make(map[string]EventHandlers),
		infoline: "",
	}
	box.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
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
				tview.Print(screen, needles, x, centerY, width, tview.AlignCenter, tcell.ColorDarkRed)
			}
			for cx := x + 1; cx < x+width-2; cx++ {
				col := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(*cc.color)
				lw := tview.BoxDrawingsLightHorizontal
				if !cc.valid {
					lw = tview.BoxDrawingsHeavyHorizontal
					col = tcell.StyleDefault.Foreground(tcell.ColorRed).Background(*cc.color)
				}
				// screen.SetContent(cx, centerY, lw, nil, col)
				if !cc.plain {
					screen.SetContent(cx, centerY, lw, nil, col)
				}
			}
			if cc.infoline != "" {
				tview.Print(screen, cc.infoline, x+1, centerY, width-1, tview.AlignCenter, txtColor)
			}
		} else {
			// cc.InRect(lowerCenterY +2)
			if width-2 >= 8 && lowerCenterY+2 <= y+height-2 {
				// for cx := x + 1; cx < x+width-1; cx++ {
				tview.Print(screen, fmt.Sprintf("[#%06x:-:b]%s[-:-:-]", cc.GetFgColor().Hex(), strings.ToUpper(fmt.Sprintf("%06x",cc.color.Hex()))), x+1, lowerCenterY+2, width-2, tview.AlignCenter, txtColor)
				// screen.SetContent(cx, lowerCenterY + 4, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(*cc.color))
				// }
			}
			// {{define "locked"}}{{- if locked -}}{{- end -}}{{- end -}}
			// {{define "selected"}}{{- if selected -}}{{- end -}}{{- end -}}
			status_tpl := MakeTemplate("color_status", `
      {{define "locked"}}{{- if not locked }}  {{ else }}  {{ end -}}{{- end -}}
      {{define "selected"}}{{- if selected }}   {{ else }}   {{ end -}}{{- end -}}
    `, template.FuncMap{
				"locked":   cc.GetLocked,
				"selected": cc.GetSelected,
				"dirty":    cc.GetDirty,
				"css":      cc.GetColor,
			})
			// sel := status_tpl(`{{- template "selected" . -}}`, cc)
			lock := status_tpl(`{{- template "locked" . -}}`, cc)
			// tview.Print(screen, sel, x+1, (lowerCenterY - centerY) / 2, width-2, tview.AlignCenter, txtColor)
			tview.Print(screen, lock, x+1, lowerCenterY, width-2, tview.AlignCenter, txtColor)
		}

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
	defer cc.l.Unlock()

	cc.selected = s
	cc.updateStyle()
}

func (cc *CoolorColor) GetFgColor() tcell.Color {
	c, ok := MakeColor(cc)
	// dump.P(cc.TerminalPreview())
	if ok {
		r, g, b := c.LinearRgb()
		if (255*float64(r)*0.299 + 255*float64(g)*0.587 + 255*float64(b)*0.114) > 150 {
			// if (255*float64(r)*0.2926 + 255*float64(g)*0.5152 + 255*float64(b)*0.1722) > 150 {
			// if (float64(r)*0.2126 + float64(g)*0.7152 + float64(b)*0.0722) > 140 {
			return tcell.ColorBlack
		}
		return tcell.ColorWhite
	}
	return tcell.ColorWhite
	// r, g, b := cc.RGB
	// if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
}

func (cc *CoolorColor) Html() string {
	return strings.ToUpper(fmt.Sprintf("#%06x", cc.color.Hex()))
}

func (cc *CoolorColor) TVPreview() string {
	return strings.ToUpper(fmt.Sprintf(
		"[#%06x:#%06x:-]#%06x[-:-:-]",
		getFGColor(*cc.color).Hex(),
		cc.color.Hex(),
		cc.color.Hex(),
	))
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
	// MainC.app.QueueUpdateDraw(func() {
		if cc.plain {
			cc.Box.SetBorderPadding(0, 0, 0, 0)
			cc.Box.SetBorder(false)
			cc.Box.Blur()
			return
		}
		if cc.selected || cc.centered {
			inverse := cc.GetFgColor()
			cc.Box.
				SetBorder(true).
				SetBorderAttributes(tcell.AttrBold).
				SetBorderPadding(0, 0, 0, 0).
				SetBorderColor(inverse).
				SetTitleColor(inverse)
			cc.Box.Focus(nil)
		} else {
			cc.Box.SetBorderPadding(0, 0, 0, 0)
			cc.Box.SetBorder(false)
			cc.Box.Blur()
		}
	// })
}

func (cp *CoolorColor) Draw(screen tcell.Screen) {
	cp.Box.DrawForSubclass(screen, cp)
	x, y, w, h := cp.Box.GetRect()
	_, _, _, _ = x, y, w, h
	// dump.P(x,y,w,h)
	if cp.plain || cp.pallette == nil || cp.pallette.menu == nil || !cp.selected {
		return
	}
	cx, cy, cw, ch := cp.GetInnerRect()
	cw = ((cw * 2) / 2) - 1
	cp.SetRect(cx, cy, cw, ch)
	_, _, _, _ = cx, cy, cw, ch
	cx, cy, cw, ch = cp.GetInnerRect()
	_, _, _, _ = cx, cy, cw, ch
	mx, my, mw, mh := cp.pallette.menu.GetRect()
	imx, imy, imw, imh := cp.pallette.menu.GetInnerRect()
	// mw = (((imw * 3) / 3))
	mw = 5
	mh = (cp.pallette.menu.GetItemCount() * 3) + 2
	// dump.P(imx, imy, imw, imh)
	_, _, _, _ = mx, my, mw, mh
	_, _, _, _ = imx, imy, imw, imh
	centerX := cx + ((cw / 2) - (mw / 2))
	// centerY := cy + ((ch/2 + (ch / 3) - mh/2) - (mh / 2))
	// centerY := cy + ((ch/2 + (ch / 2)) - (ch / 3))
  centerY := y + ((mh / cp.pallette.menu.GetItemCount()))
	cp.pallette.menu.SetRect(centerX, centerY, mw, mh)
	cp.pallette.menu.Draw(screen)
}

// vim: ts=2 sw=2 et ft=go
