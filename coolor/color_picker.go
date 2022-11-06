package coolor

import (
	"fmt"
	"strings"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/gookit/goutil/dump"

	"github.com/digitallyserviced/coolors/coolor/events"
	// "github.com/digitallyserviced/coolors/coolor/util"
	"github.com/digitallyserviced/coolors/theme"
)


type ColorPicker struct {
	*tview.Flex
	acts    *tview.Flex
	current Color
  *ColorSpace
	*events.EventNotifier
  OnChange func(c Color)
	OnAdd func(c Color)
	OnSet func(c Color)
}

func OpenColorPicker(mc *MainContainer, colors... *CoolorColor) bool {
  c := RandomColor()
  if len(colors) > 0 {
    c = *colors[0].GetColorable()
  }
	tv := NewTabbedView()
	rgbPick := NewColorPicker(&RGBChannel, c)
	hslPick := NewColorPicker(&HSLChannel, c)
  setColor := func(c Color) {
    colors[0].SetColor(c.GetCC().Color)
	}
  addColor := func(c Color) {
		MainC.palette.AddCoolorColor(c.GetCC())
	}
  syncColor := func(c Color) {
    rgbPick.SetColor(c.GetCC())
    hslPick.SetColor(c.GetCC())
  }
	rgbPick.OnSet = setColor
	hslPick.OnSet = setColor
	rgbPick.OnAdd = addColor
  hslPick.OnAdd = addColor
  rgbPick.OnChange = syncColor
  hslPick.OnChange = syncColor
  picked := events.NewAnonymousHandlerFunc(func(e events.ObservableEvent) bool {
		mc.Pop()
		return true
	})
	rgbPick.Register(events.PromptedEvents, picked)
	hslPick.Register(events.PromptedEvents, picked)
	rgbpick := NewTabView("  RGB", tv.TakeNext(), rgbPick)
	tv.AddTab(rgbpick)
	hslpick := NewTabView("  HSL", tv.TakeNext(), hslPick)
	tv.AddTab(hslpick)
	tv.UpdateView()
	sb := NewSideBar("colorpick", tv, 0, 40)
	mc.Push("colorpick", sb, true)
	mc.app.SetFocus(tv)
	return false
}

func NewColorPicker(cs *ColorSpace, colors... Color) *ColorPicker {
  c := RandomColor()
  if len(colors) > 0 {
    c = colors[0]
  }
	ttvc := &ColorPicker{
		Flex:          tview.NewFlex(),
		acts:          tview.NewFlex(),
		current:       c,
		ColorSpace:    cs,
		EventNotifier: events.NewEventNotifier("colorpick"),
		OnAdd: func(c Color) {
		},
		OnSet: func(c Color) {
		},
	}
	ttvc.Flex.SetDirection(tview.FlexRow)
	ttvc.Flex.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	ttvc.SetDontClear(false)
	ttvc.acts.SetDirection(tview.FlexColumn)
  ttvc.acts.SetBorder(false).SetBorderPadding(0, 1, 0, 0)
	ttvc.acts.SetDontClear(true)
	ttvc.acts.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	// ﰂ  ﰂ 烙
	add := tview.NewButton("ﰂ  add")
	add.SetBorderPadding(0, 0, 1, 1)
  add.SetBorderSides(false, false, true, false)
	add.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	add.SetBorderColor(theme.GetTheme().Primary)
	add.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	add.SetLabelColor(theme.GetTheme().Primary)
	add.SetBackgroundColorActivated(theme.GetTheme().GrayerBackground)

	set := tview.NewButton("  set")
	set.SetBorderPadding(0, 0, 1, 1)
  set.SetBorderSides(false, false, true, false)
	set.SetBorder(true).SetTitleAlign(tview.AlignCenter)
	set.SetBorderColor(theme.GetTheme().Secondary)
	set.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	set.SetLabelColor(theme.GetTheme().Secondary)
	set.SetBackgroundColorActivated(theme.GetTheme().GrayerBackground)

	ttvc.acts.AddItem(nil, 0, 1, false)
	ttvc.acts.AddItem(set, 0, 3, false)
	ttvc.acts.AddItem(nil, 0, 1, false)
	ttvc.acts.AddItem(add, 0, 3, false)
	ttvc.acts.AddItem(nil, 0, 1, false)
	return ttvc
}

// InputBubbler implements TabViewContent
func (ccs *ColorPicker) InputBubbler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ccs.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		})
}

// hide implements TabViewContent
func (*ColorPicker) hide(t *TabView) {
	// panic("unimplemented")
}

// show implements TabViewContent
func (ttvc *ColorPicker) show(t *TabView) {
	// x, y, w, h := ttvc.GetRect()

	// tview.Borders = OrigBorders

	AppModel.app.QueueUpdateDraw(func() {

		// total := 255.0 / steps
		// 🭣🮃🭘  🭥🬎🭚 🭭
		ttvc.UpdateBars()
	})
}

func (ttvc *ColorPicker) SetColor(cc *CoolorColor) {
	ttvc.current = *cc.GetColorable()
  ttvc.UpdateBars()
}
func (ttvc *ColorPicker) GetColorTicks(
	c Color,
  channel ColorChannel,
) []string {
	_,_,w,_ := ttvc.GetRect()
	ticks := make([]string, 0)
  current := channel.Get(c)

  aw := float64(w-10)
  pcts := 1.0 / aw
	// steps := channel.Max / aw 
  // tickC := channel.Max / steps
	// pct := channel.scale * channel.Max
	for i := 0.0; i < aw; i++ {
		tc := i * pcts * channel.Max
		chr := ' '
		// if tc-((steps/2)*pct) < current && current < tc+((steps/2)*pct) {
		if ((i)*pcts*channel.Max) <= current && current < ((i+1)*pcts*channel.Max) {
			chr = '🭭'
		}
		// col := Color{ttvc.current.R, ttvc.current.G, tc}
		col := channel.Set(tc, ttvc.current)
		cc := col.GetCC()
		it := fmt.Sprintf(
			"[#%06x:%s:b]%c[-:-:-]",
			cc.GetFgColorShade().Hex(),
			cc.Html(),
			chr,
		)
		ticks = append(ticks, it)
	}
	return ticks
}
func (ttvc *ColorPicker) AddChannel(cc *ColorChannel) {
  ticks := ttvc.GetColorTicks(ttvc.current, *cc)
	ttvc.NewChannelBar(cc.Display(ttvc.current), ticks)
}
func (ttvc *ColorPicker) UpdateBars() {
	ttvc.Flex.Clear()

	tview.Styles.BorderColor = tcell.Color101
	prevCol := tview.NewTextView()
	prevCol.SetDynamicColors(true).
		SetText(ttvc.current.GetCC().TVPreview()).
		SetTextAlign(tview.AlignCenter)
	prevCol.SetBorderPadding(0, 0, 0, 0)
	ttvc.Flex.AddItem(prevCol, 0, 2, false)

  for _, ch := range ttvc.Channels {
    ttvc.AddChannel(&ch)
  }

	// ticks = make([]string, 0)
 //  ticks := ttvc.GetColorTicks(
	// 	ttvc.current.G,
	// 	ChannelGreen.Set,
	// )
	// ttvc.NewChannelBar(ChannelGreen.Display(ttvc.current), ticks)
	//
	// ticks = make([]string, 0)
	// ticks = ttvc.GetColorTicks(steps, pct, ticks, ttvc.current.B, ChannelBlue.Set)
	// ttvc.NewChannelBar(ChannelBlue.Display(ttvc.current), ticks)
	// // colActs := tview.NewTextView()
	// colActs.SetDynamicColors(true).SetText(ttvc.current.GetCC().TVPreview()).SetTextAlign(tview.AlignCenter)
	// colActs.SetBorderPadding(0,0,0,0)
	ttvc.Flex.AddItem(ttvc.acts, 3, 3, false)
}

func (ttvc *ColorPicker) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ttvc.WrapInputHandler(
		func(ek *tcell.EventKey, f func(p tview.Primitive)) {
			// var current Color
			fmt.Println(ttvc.current)
        switch ek.Key() {
        case tcell.KeyEnter:
					ttvc.Notify(
						*ttvc.NewObservableEvent(events.SecondaryEvent, "set", ttvc.current.GetCC(), nil),
					)
					if ttvc.OnSet != nil {
						ttvc.OnSet(ttvc.current)
					}
      return
      }

			if ek.Key() == tcell.KeyRune {
				incr := 1.0
				if ek.Rune()&' ' == 0 {
					incr = 5.0
				}
        var cchan *ColorChannel
        var inc bool
				switch ek.Rune() {
				case 'e', 'E':
          cchan = &ttvc.Channels[0]
          inc = true
				case 'd', 'D':
          cchan = &ttvc.Channels[1]
          inc = true
				case 'c', 'C':
          cchan = &ttvc.Channels[2]
          inc = true

				case 'q', 'Q':
          cchan = &ttvc.Channels[0]
				case 'a', 'A':
          cchan = &ttvc.Channels[1]
				case 'z', 'Z':
          cchan = &ttvc.Channels[2]
				case '+':
					ttvc.Notify(
						*ttvc.NewObservableEvent(events.PrimaryEvent, "add", ttvc.current.GetCC(), nil),
					)
					if ttvc.OnAdd != nil {
						ttvc.OnAdd(ttvc.current)
					}
				}
        if cchan != nil {
        var x Color
        if inc {
        x = ttvc.Incr(cchan, incr)
        } else {
          x = ttvc.Decr(cchan, incr)
        }
					ttvc.current = x // ttvc.Incr(ttvc.Channels[0], incr)
          if ttvc.OnChange != nil {
            ttvc.OnChange(x)
          }

        }

			}

		},
	)
}

func (ttvc *ColorPicker) Decr(c *ColorChannel, incr float64) Color {
	x := c.Decr(incr*c.Step, ttvc.current)
	return x
}
func (ttvc *ColorPicker) Incr(c *ColorChannel, incr float64) Color {
	x := c.Incr(incr*c.Step, ttvc.current)
	return x
}

func (ttvc *ColorPicker) NewChannelBar(t string, ticks []string) {
	tv := tview.NewTextView()
	tv.SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetBorder(true).SetBorderVisible(false).
		SetBorderColor(theme.GetTheme().InfoLabel)
	tv.SetTitle(t)
	ttvc.Flex.AddItem(tv, 0, 2, false)
	tv.SetText(
		fmt.Sprintf(
			"[white:gray:db] - [-:-:-]%s[white:gray:db] + [-:-:-]",
			strings.Join(ticks, ""),
		),
	)
}


func (ttvc *ColorPicker) I() {

}
