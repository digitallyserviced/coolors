package coolor

import (
	"fmt"
	"math"
	"strings"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"

	"github.com/digitallyserviced/coolors/theme"
)


type SideBar struct {
  *FixedFloater
  posLeft bool
  width, height int
}

func NewSideBar(name string, p tview.Primitive, args... int) *SideBar {
  w := 30
  h := 10
  pos := 1
  if len(args) > 0 {
    pos = args[0]
    if len(args) > 1 {
      w = args[1]

    if len(args) > 2 {
        h = args[2]
      }
    }
  }
	f := &FixedFloater{
		Header:             tview.NewTextView(),
		Footer:             tview.NewTextView(),
		RootFloatContainer: NewFloater(p),
	}
  sb := &SideBar{
  	FixedFloater: f,
  	posLeft:      pos == 0,
  	width:        w,
  	height:       h,
  }
	f.RootFloatContainer.
		SetBorder(false).
		SetBorderPadding(0, 0, 1, 1).
		SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.Container.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	f.RootFloatContainer.Rows.Clear()
	f.RootFloatContainer.Rows.AddItem(f.Container, 0, 10, true)

  if pos == 1 {
    f.RootFloatContainer.Clear()
    f.RootFloatContainer.AddItem(nil, 0, 100-w, false)
    f.RootFloatContainer.AddItem(f.Rows, 0, w, true)
  } else {
    f.RootFloatContainer.Clear()
    f.RootFloatContainer.AddItem(f.Rows, 0, w, true)
    f.RootFloatContainer.AddItem(nil, 0, 100-w, false)
  }

	f.Header.SetDynamicColors(true)
	f.Header.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name)).SetBorderPadding(1, 1, 1, 1)
	f.Header.SetBackgroundColor(theme.GetTheme().SidebarLines).
		SetBorderColor(theme.GetTheme().SidebarBackground)
	bw := f.Header.BatchWriter()
	bw.Close()

	f.Footer.SetDynamicColors(true)
	f.Footer.SetTextAlign(tview.AlignCenter).
		SetText(fmt.Sprintf("[yellow]%s[-]", name))
	f.Footer.SetBackgroundColor(theme.GetTheme().SidebarBackground).
		SetBorderColor(theme.GetTheme().SidebarBackground)

	sb.UpdateView()

	return sb
}

func (sb *SideBar) UpdateView() {
	// f.Container.SetBackgroundColor(theme.GetTheme().TopbarBorder)
	sb.Container.SetBorder(false)
	sb.Flex.SetBorder(false)
	sb.Rows.SetBorder(false)
	sb.Container.SetDirection(tview.FlexRow)
	sb.Container.Clear()
	// f.Container.AddItem(f.Header, 3, 0, false)
	// f.Container.AddItem(f.List, 0, 6, true)

	// f.Lister.UpdateView()

	sb.GetRoot().UpdateView()
	// f.Header.SetBorder(true).SetBorderPadding(1, 1, 1, 1)
	// f.Container.AddItem(f.Footer, 2, 0, false)
}

type TestTabViewContenter struct {
	*tview.Grid
  current Color
}

func NewTestTabViewContenter() *TestTabViewContenter {
  ttvc := &TestTabViewContenter{
  	Grid: tview.NewGrid(),
    current: Color{
    	R: 0.2,
    	G: 0.8,
    	B: 0.4,
    },
  }
  return ttvc
}

// InputBubbler implements TabViewContent
func (ccs *TestTabViewContenter) InputBubbler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ccs.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		})
}

// hide implements TabViewContent
func (*TestTabViewContenter) hide(t *TabView) {
	// panic("unimplemented")
}

// show implements TabViewContent
func (ttvc *TestTabViewContenter) show(t *TabView) {
  dump.P(t.GetRect())
  x,y,w,h := t.GetRect()

  steps := math.Ceil(255.0 / float64(w-4))
  total := math.Floor(255.0 / steps)
  pct := total / 255.0
  scale := 1
  _,_,_,_ = x,y,h,scale

  ticks := make([]string, 0)
  dump.P(steps, total, pct)

  for i := 0; i < 255; i+=int(steps) {
    r := ttvc.current.R
    tc := pct *float64(i) 
    chr := ' '
    if math.Abs(r - tc) < (pct) {
      chr = '🭯'
    }
    // col := Color{ttvc.current.R,  ttvc.current.G, ttvc.current.B}
    col := Color{tc,  ttvc.current.G, ttvc.current.B}
    dump.P(col)
    cc := col.GetCC()
    it := fmt.Sprintf("[#%06x:%s:b]%c[-:-:-]", cc.GetFgColor().Hex(),cc.Html(),chr)
    ticks = append(ticks, it)
  }
  rtv := tview.NewTextView()
  rtv.SetDynamicColors(true)
  rtv.SetText(strings.Join(ticks, ""))
  ttvc.Grid.AddItem(rtv, 0, 0, 1, 1, 1, 1, false)
  ticks = make([]string, 0)

  for i := 0; i < 255; i+=int(steps) {
    g := ttvc.current.G
    tc := pct *float64(i) 
    chr := ' '
    if math.Abs(g - tc) < (steps * pct) {
      chr = '🭯'
    }
    col := Color{ttvc.current.R,  tc, ttvc.current.B}
    cc := col.GetCC()
    it := fmt.Sprintf("[#%06x:%s:b]%c[-:-:-]", cc.GetFgColor().Hex(),cc.Html(),chr)
    ticks = append(ticks, it)
  }
  gtv := tview.NewTextView()
  gtv.SetDynamicColors(true)
  gtv.SetText(strings.Join(ticks, ""))
  ttvc.Grid.AddItem(gtv, 1, 0, 1, 1, 1, 1, false)
  ticks = make([]string, 0)

  for i := 0; i < 255; i+=int(steps) {
    b := ttvc.current.B
    tc := pct *float64(i) 
    chr := ' '
    if math.Abs(b - tc) < (steps * pct) {
      chr = '🭯'
    }
    col := Color{ttvc.current.R,  ttvc.current.G, tc}
    cc := col.GetCC()
    it := fmt.Sprintf("[#%06x:%s:b]%c[-:-:-]", cc.GetFgColor().Hex(),cc.Html(),chr)
    ticks = append(ticks, it)
  }
  btv := tview.NewTextView()
  btv.SetDynamicColors(true)
  btv.SetText(strings.Join(ticks, ""))
  ttvc.Grid.AddItem(btv, 2, 0, 1, 1, 1, 1, false)
}
