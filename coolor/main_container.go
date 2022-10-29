package coolor

import (
	"flag"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/structs"

	"github.com/digitallyserviced/coolors/coolor/anim"
	. "github.com/digitallyserviced/coolors/coolor/anim"
	"github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/util"

	"github.com/digitallyserviced/coolors/theme"

	"github.com/digitallyserviced/coolors/status"
	ct "github.com/digitallyserviced/coolors/theme"
)

type MainContainer struct {
	floater Floater
	current tview.Primitive
	preview *Square
	pages   *tview.Pages
	palette *CoolorMainPalette
	mixer   *CoolorBlendPalette
	shades  *CoolorShadePalette
	scratch *PaletteFloater
	info    *CoolorColorFloater
	editor  *CoolorColorEditor
	*tview.Flex
	fileviewer *CoolorFileView
	filetree   *PaletteFileTree
	menu       *CoolorToolMenu
	ScrollLine *ScrollLine
	sidebar    *FixedFloater
	main       *tview.Flex
	app        *tview.Application
	options    *structs.Data
	conf       *HistoryDataConfig
	screen     *tcell.Screen
	*events.EventNotifier
	inputSwallowed bool
	*PageStack
}

type PageStack struct {
	Stack
	Pages *tview.Pages
}

func NewPageStack(p *tview.Pages) *PageStack {
	ps := &PageStack{
		Stack: *NewStack(),
		Pages: p,
	}
	return ps
}

func (ps *PageStack) Pop() bool {
	ipg := ps.Stack.Pop()
	if pg, ok := ipg.(*tview.Page); ok {
		ps.Pages.RemovePage(pg.Name)
		return true
	}
	return false
}

func (ps *PageStack) Push(name string, p tview.Primitive, rz bool) {
	pg := ps.Pages.NewPage(name, p, rz, true)
	ps.Stack.Push(pg)
	ps.Pages.Addpage(pg)
	ps.Pages.ShowPage(pg.Name)
}

var MainC *MainContainer
var cpName = flag.String("palette", "", "Palette name to open")

func CreateStartupPalette() *CoolorMainPalette {
	flag.Parse()
	values := flag.Args()

	colsize := 0
	var pal *CoolorColorsPaletteMeta
	if len(*cpName) > 0 {
		pal = Store.FindNamedPalette(*cpName)
	}
	if pal != nil {
		return NewCoolorColorsPaletteFromMeta(pal)
	}
	if len(values) > 0 {
		colsize = len(values)
	} else {
		colsize = 16
	}

	var colors *CoolorMainPalette

	if len(values) > 0 {
		colors = NewCoolorPaletteFromCssStrings(values)
	} else {
		colors = NewCoolorPaletteWithColors(GenerateRandomColors(colsize))
	}
	return colors
}

func NewMainContainer(app *tview.Application) *MainContainer {
	pgs := tview.NewPages()
	MainC = &MainContainer{
		Flex:    tview.NewFlex(),
		floater: nil,
		pages:   pgs,
		// anims:         tview.NewPages(),
		main:          tview.NewFlex(),
		app:           app,
		options:       ActionOptions,
		conf:          NewPaletteHistoryFile(),
		EventNotifier: events.NewEventNotifier("main"),
		current:       nil,
		PageStack:     NewPageStack(pgs),
	}

	MainC.menu = NewCoolorColorMainMenu(app)
	MainC.palette = CreateStartupPalette()
	MainC.conf.AddPalette("random", MainC.palette)
	MainC.menu.Init()
	// events.Global.Register(events.PaletteColorSelectedEvent, MainC.ScrollLine)
	MainC.palette.SetMenu(MainC.menu)
	MainC.palette.Register(events.PaletteColorSelectedEvent, MainC.menu)
	MainC.editor = NewCoolorEditor(app, MainC.palette)
	MainC.preview = NewRecursiveSquare(MainC.palette.GetPalette(), 5)
	// MainC.fileviewer = NewFileViewer()
	MainC.filetree = NewPaletteFileTree()
	MainC.Init()
	return MainC
}
func (mc *MainContainer) Init() {
	InitPlugins()
	mc.SetBackgroundColor(ct.GetTheme().SidebarBackground)
	mc.SetDirection(tview.FlexRow)
	MainC.ScrollLine = NewScrollLine()
	mc.AddItem(mc.pages, 0, 80, false)
	mc.AddItem(MainC.ScrollLine, 1, 0, false)
	mc.pages.AddPage("editor", mc.editor, true, false)
	mc.pages.AddPage("preview", mc.preview, true, true)
	mc.pages.AddPage("filetree", mc.filetree, false, false)
	mc.pages.AddAndSwitchToPage("palette", mc.palette, true)
	mc.pages.SetChangedFunc(func() {
		name, page := mc.pages.GetFrontPage()
		mc.current = page
		status.NewStatusUpdate("action", name)
	})
	AppModel.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyF40 {
			mc.inputSwallowed = true
		} else if event.Key() == tcell.KeyF41 {
			mc.inputSwallowed = false
		}

		if mc.inputSwallowed {
			fmt.Printf("swallow:%v", mc.inputSwallowed, event)
			p := MainC.app.GetFocus()
			p.InputHandler()(event, func(tp tview.Primitive) {
				MainC.app.SetFocus(tp)
			})
			return nil
		}

		ch := event.Rune()
		if event.Modifiers() == tcell.ModShift {
		}
		switch ch {
		case 'Q':
			AppModel.app.Stop()
			return nil
		}
		return event
	})
}

func (mc *MainContainer) OpenTagView(cp *CoolorColorsPalette) {
	name, _ := mc.pages.GetFrontPage()
	if name == "tagView" {
	} else {
		mc.pages.RemovePage("tagView")
	}
	tf := NewTagEditFloater(cp)
	mc.pages.AddAndSwitchToPage("tagView", tf, true)
}
func (mc *MainContainer) CloseConfig() {
	mc.conf.ref.Close()
}

// half_bar_right = "╸"
// half_bar_left = "╺━╸━ ╺"
// bar = "━"
//
// width = self.width or options.max_width
// start, end = self.highlight_range
//
// start = max(start, 0)
// end = min(end, width)
//
// output_bar = Text("", end="")
//
// if start == end == 0 or end < 0 or start > end:
//
//	output_bar.append(Text(bar * width, style=background_style, end=""))
//	yield output_bar
//	return
//
// # Round start and end to nearest half
// start = round(start * 2) / 2
// end = round(end * 2) / 2
//
// # Check if we start/end on a number that rounds to a .5
// half_start = start - int(start) > 0
// half_end = end - int(end) > 0

// InputHandler returns the handler for this primitive.
func (mc *MainContainer) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return mc.pages.WrapInputHandler(
		func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
			ch := event.Rune()
			kp := event.Key()
			name, page := mc.pages.GetFrontPage()
			if kp == tcell.KeyEscape {
				if mc.menu.Activated() != nil {
					mc.menu.Activated().Cancel()
				}
				if name == "floater" {
					mc.pages.HidePage("floater")
					page.Blur()
					mc.pages.RemovePage("floater")
					mc.floater = nil
				}
				return
			}
			switch ch {
			case 'P':
				mc.Pop()
				tv := NewTabbedView()
				tst := NewTestTabViewContenter()
        hsvpick := NewTabView(" HSV", tv.TakeNext(), tst)
				tv.AddTab(hsvpick)
				tv.UpdateView()
				sb := NewSideBar("colorpick", tv, 0, 40)
				mc.Push("colorpick", sb, true)
			case 'L':
				mc.app.QueueUpdateDraw(func() {
					if mc.sidebar == nil {
						// ccsTwo := NewCoolorColorSwatch(
						// 	func(cs *CoolorColorSwatch) *CoolorColorsPalette {
						// 		cp := NewCoolorPaletteWithColors(
						// 			GenerateRandomColors(20),
						// tv.AddTab(NewTabView(" Favorites", tv.TakeNext(), ccsTwo))
						// tv.AddTab(NewTabView("ﱬ Stream", tv.TakeNext(), ccsTwo))
						// tv.AddTab(NewTabView("Alter Fuck",tv.TakeNext(), ccs111))
						tv := NewTabbedView()
						mc.sidebar = NewFixedFloater("QuickColors", tv)
						ccps := NewCoolorColorsPaletteSwatch(
							func(cs *CoolorColorsPaletteSwatch) []*CoolorColorsPalette {
								cps := make([]*CoolorColorsPalette, 0)
								ccps := GetStore().PaletteHistory(false)
								for _, v := range ccps {
									if v.Current == nil || v.Current.Colors == nil ||
										len(v.Current.Colors) == 0 {
										continue
									}
									cps = append(cps, v.Current.GetPalette())
								}
								return cps
							},
						)
						ccs := NewCoolorColorSwatch(
							func(cs *CoolorColorSwatch) *CoolorColorsPalette {
								cp := GetStore().FavoriteColors.GetPalette()
								return cp
							},
						)
						ccpstv := NewTabView(" Palettes", tv.TakeNext(), ccps)
						ccpstv.SetBackgroundColor(ct.GetTheme().GrayerBackground)
						ccstv := NewTabView(" Favorites", tv.TakeNext(), ccs)
						ccstv.SetBackgroundColor(ct.GetTheme().GrayerBackground)
						f := events.NewAnonymousHandlerFunc(
							func(e events.ObservableEvent) bool {
								switch {
								case e.Type&events.ColorSelectedEvent != 0:
									col, ok := e.Ref.(*CoolorColor)
									if ok && col != nil {
										ccstv.Frame.Clear().
											AddText(col.GetMeta().String(), false, tview.AlignCenter, tcell.ColorRed)
										return true
									}
								case e.Type&events.ColorSelectionEvent != 0:
									col, ok := e.Ref.(*CoolorColor)
									if ok {
										MainC.palette.AddCoolorColor(col)
										return true
									}
								case e.Type&events.PaletteColorSelectionEvent != 0 || e.Type&events.PaletteColorSelectedEvent != 0:
									pal, ok := e.Ref.(*CoolorColorsPalette)

									sqs := fmt.Sprintf(
										" %s ",
										strings.Join(
											pal.GetPalette().MakeSquarePalette(false),
											" ",
										),
									)
									if ok {
										ccpstv.Frame.Clear().
											AddText(sqs, false, tview.AlignCenter, tcell.ColorRed)
									}
								default:
									fmt.Println(e)

								}
								return true
							},
						)
						ccps.Register(events.ColorSelectionEvent, f)
						ccs.Register(events.ColorSelectedEvent, f)
						tv.AddTab(ccpstv)
						tv.AddTab(ccstv)
						tv.UpdateView()
						mc.pages.AddPage("sidebar", mc.sidebar.GetRoot(), true, true)
						mc.pages.ShowPage("sidebar")
						mc.app.SetFocus(mc.sidebar.GetRoot().Item)
					} else {
						// name, page := mc.pages.GetFrontPage()
						if name == "sidebar" {
							mc.pages.HidePage("sidebar")
							page.Blur()
							mc.pages.RemovePage("sidebar")
							mc.sidebar = nil
						} else {
							mc.pages.ShowPage("sidebar")
							mc.app.SetFocus(mc.sidebar.GetRoot().Item)
						}
						AppModel.helpbar.SetTable("sidebar")
					}
				})
			case 't':
				//  
			case 'i':
				cc, _ := mc.palette.GetSelected()
				if mc.info == nil {
					mc.info = NewCoolorColorFloater(cc)
					mc.pages.AddPage("info", mc.info, true, false)
					mc.pages.ShowPage("info")
				} else {
					// name, page := mc.pages.GetFrontPage()
					if name == "info" {
						mc.pages.HidePage("info")
						page.Blur()
					} else {
						mc.pages.ShowPage("info")
						mc.app.SetFocus(mc.info.Flex)
						mc.info.Color.UpdateColor(cc)
					}
					AppModel.helpbar.SetTable("info")
				}
			case '$':
				mc.OpenTagView(mc.palette.CoolorColorsPalette)
			case '#':
				if mc.scratch == nil {
					p := NewCoolorPaletteWithColors(GenerateRandomColors(8))
					cpc := NewScratchPaletteFloater(p.GetPalette())
					mc.scratch = cpc
					mc.pages.AddPage("scratch", mc.scratch, true, false)
					mc.pages.ShowPage("scratch")
				} else {
					// name, page := mc.pages.GetFrontPage()
					if name == "scratch" {
						mc.pages.HidePage("scratch")
						page.Blur()
					} else {
						mc.pages.ShowPage("scratch")
						mc.app.SetFocus(mc.scratch.Palette)
					}
					AppModel.helpbar.SetTable("scratch")
				}
			case 'f':
				if name != "filetree" {
					mc.pages.ShowPage("filetree")    // .HidePage("editor")
					mc.pages.SendToFront("filetree") // .HidePage("editor")
					mc.app.SetFocus(mc.filetree)
					AppModel.helpbar.SetTable("filetree")
				} else {
					mc.pages.HidePage("filetree") // .HidePage("editor")
					mc.pages.SendToBack("filetree")
					mc.app.SetFocus(mc.palette)
					AppModel.helpbar.SetTable("palette")
				}
			// case 'F':
			// 	mc.pages.SwitchToPage("fileviewer") // .HidePage("editor")
			// 	mc.app.SetFocus(mc.fileviewer.treeView)
			// 	AppModel.helpbar.SetTable("fileviewer")
			case 'S':
				mc.pages.SwitchToPage("shades") // .HidePage("editor")
				AppModel.helpbar.SetTable("shades")
			case 'M':
				mc.pages.SwitchToPage("mixer") // .HidePage("editor")
				AppModel.helpbar.SetTable("mixer")
			case 'p':
				mc.pages.SwitchToPage("palette") // .HidePage("editor")
				AppModel.helpbar.SetTable("palette")
			case 'Y':
				mc.pages.ShowPage("preview").
					SendToFront("preview")
					//.HidePage("palette")
				AppModel.helpbar.SetTable("preview")
				x, y, w, h := mc.pages.GetInnerRect()
				_, _, _, _ = x, y, w, h
				mc.preview.SetRect(x, y, w/2, h/2)
				mc.preview.TopInit(8)
			case 'e':
				mc.pages.SwitchToPage("editor") //.HidePage("palette")
				AppModel.helpbar.SetTable("editor")
			case 'V':
				ani := NewFrameAnimator()
				ani.Start()
				ani.Control.Play()
			// case 'v':
			// 	bFrequency += 0.1
			// case 'G':
			// 	bDamping += 0.1
			// case 'g':
			// 	bDamping -= 0.1
			case 'u':
				noti := anim.GetAnimator().GetAnimation("notif")
				dump.P(noti)
				if noti == nil {
					noti = NewNotification(
						"SHIT",
						"This is a test notification of the notification broadcastsing systems.\n!!",
					)
				} else {
					noti.Next()
				}
				noti.Control.Play()

				// })
				// MainC.pages.AddPage("anim", tview.NewP, resize bool, visible bool)
			}

			if page == nil {
				return
			}
			switch name {
			case "shades":
				HandleVimNavigableHorizontal(mc.shades, ch, kp)
				HandleCoolorSelectable(mc.shades, ch, kp)
				// dump.P(fmt.Sprintf("%s horiz input handled", name))
			case "mixer":
				HandleVimNavigableHorizontal(mc.mixer, ch, kp)
				HandleCoolorSelectable(mc.mixer, ch, kp)
				// dump.P(fmt.Sprintf("%s horiz input handled", name))
			case "palette":
				HandleVimNavigableHorizontal(mc.palette, ch, kp)
				// HandleCoolorSelectable(mc.palette, ch, kp)
				// dump.P(fmt.Sprintf("%s horiz input handled", name))
			case "info":
				HandleVimNavigableHorizontal(mc.info, ch, kp)
				HandleVimNavigableVertical(mc.info, ch, kp)
				// HandleCoolorSelectable(mc.info, ch, kp)
			case "scratch":
				// HandleVimNavigableHorizontal(mc.scratch.Palette.Palette.Palette, ch, kp)
				// HandleCoolorSelectable(mc.scratch.Palette.Palette.Palette, ch, kp)
			}

			if name == "mixer" || name == "palette" || name == "shades" {
				menuHandler := MainC.menu.list.InputHandler()
				menuHandler(event, setFocus)
				// HandleVimNavigableVertical(MainC.menu, ch, kp)
			}

			if handler := page.InputHandler(); handler != nil {
				dump.P(fmt.Sprintf("%s input handled", name))
				handler(event, setFocus)
			}
		},
	)
}

func (mc *MainContainer) Draw(s tcell.Screen) {
	mc.Box.DrawForSubclass(s, mc)
	mc.Flex.Draw(s)
	mc.ScrollLine.Box.DrawForSubclass(s, mc.ScrollLine)
	mc.ScrollLine.Draw(s)
}

func (mc *MainContainer) NewShades(base *CoolorColor) {
	// mc.app.QueueUpdateDraw(func() {
	if mc.pages.HasPage("shades") && mc.shades != nil {
		mc.shades.UpdateColors(base)
	} else {
		mc.shades = BlankCoolorShadePalette(base, 8)
		mc.shades.SetMenu(MainC.menu)
		mc.pages.AddPage("shades", mc.shades, true, false)
	}
	mc.pages.SwitchToPage("shades") // .HidePage("editor")
	mc.pages.HidePage("palette")
	mc.palette.Blur()
	mc.palette.ColorContainer.Blur()
	mc.app.SetFocus(mc.shades)
	// })
}

func (mc *MainContainer) NewMixer(start, end *CoolorColor) {
	if mc.pages.HasPage("mixer") && mc.mixer != nil {
		mc.mixer.UpdateColors(start, end)
	} else {
		mc.mixer = BlankCoolorBlendPalette(start, end, 8)
		mc.mixer.SetMenu(MainC.menu)
		mc.pages.AddPage("mixer", mc.mixer, true, false)
	}
	// name, page := mc.pages.GetFrontPage()
	// if page != nil {
	// 	mc.pages.HidePage(name)
	// }
	mc.pages.SwitchToPage("mixer") // .HidePage("editor")
	mc.pages.HidePage("palette")
	mc.palette.Blur()
	mc.palette.ColorContainer.Blur()
	mc.app.SetFocus(mc.mixer)
}

func (mc *MainContainer) GetPalette() *CoolorColorsPalette {
	return mc.palette.GetPalette()
}

type ScrollLine struct {
	*tview.Box
	*Animation
	// cp                                           *CoolorColorsPalette
	selectedIndex                                int
	lineWidth, lineSegments                      int
	indicatorPad, indicatorBlur                  int
	selectedColor                                *tcell.Color
	indicatorTgtPos, indicatorPos, indicatorSize int
	x, vel                                       float64
}

// HandleEvent implements events.Observer
func (sl *ScrollLine) HandleEvent(o events.ObservableEvent) bool {
	// var cp *CoolorColorsPalette
	if cp, ok := o.Src.(*CoolorColorsPalette); ok {
		if cp != nil {
			var selcol *CoolorColor
			selcol, sl.selectedIndex = cp.GetSelected()
			sl.selectedColor = selcol.Color
			// AppModel.app.QueueUpdateDraw(func() {
			sl.UpdatePos(cp)
			// })
		}
	}
	return true

}

// Name implements events.Observer
func (sl *ScrollLine) Name() string {
	return "scroll_line"
}

func (sl *ScrollLine) UpdatePos(cp *CoolorColorsPalette) {
	x, y, width, height := sl.GetRect()
	_, _, _, _ = x, y, width, height
	sl.lineWidth = width - 2 - 4
	total := cp.GetItemCount()
	sl.indicatorSize = (sl.lineWidth) / total
	pos := (sl.selectedIndex + 1) * sl.indicatorSize
	sl.indicatorPad = sl.indicatorSize / 2
	sl.indicatorTgtPos = pos
	fmt.Println(sl.indicatorTgtPos, sl.indicatorPos)
	sl.Animation.GetCurrentFrame().Motions[0].UpdateTween(
		float64(sl.indicatorPos),
		float64(sl.indicatorTgtPos),
	)
	fmt.Println("delta", math.Abs(float64(sl.indicatorTgtPos)-float64(sl.indicatorPos)), sl.Animation.State.String())
	if math.Abs(float64(sl.indicatorTgtPos)-float64(sl.indicatorPos)) > 1 {
		if sl.Animation.State.Is(events.AnimationNext) || sl.Animation.State.Is(events.AnimationPaused) || sl.Animation.State.Is(events.AnimationFinished) {
			sl.Animation.Next()
			sl.Animation.Control.Play()
		}
	} else {
		sl.Animation.Control.Idle()
	}
}

func NewScrollLine() *ScrollLine {
	btmLine := MakeBoxItem("", "")
	btmLine.SetBorder(false).SetBorderPadding(0, 0, 0, 0)
	btmLine.SetDontClear(true)
	btmLine.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	sl := &ScrollLine{
		Box:             btmLine,
		Animation:       nil,
		selectedIndex:   0,
		lineWidth:       1,
		lineSegments:    1,
		indicatorPad:    1,
		indicatorBlur:   1,
		selectedColor:   nil,
		indicatorTgtPos: 1,
		indicatorPos:    1,
		indicatorSize:   1,
		x:               1,
		vel:             1,
	}
	events.Global.Register(events.PaletteColorSelectedEvent, sl)
	btmLine.SetDrawFunc(
		func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
			x, y, width, height = sl.GetRect()
			sl.lineWidth = width - 2 - 4
			btmLine.SetBorderPadding(0, 0, 2, 2)
			x, y, width, height = btmLine.GetInnerRect()
			_, _, _, _ = x, y, width, height
			width = width - 6
			fullBar, leftHalfBar, rightHalfBar := '━', '╸', '╺'
			blur := util.Clamp(int(MapVal(math.Abs(sl.vel), 0, 100, 0, 10)), 1, 10)
			sl.indicatorBlur = ((sl.indicatorBlur) + blur) / 2
			// sl.indicatorBlur =
			// ((sl.vel * 3) + m.Xvelocity) / 4
			mid := sl.indicatorPos
			pad := sl.indicatorSize / 2
			padl, padr := pad, pad
			if sl.vel > 0 {
				padr += sl.indicatorBlur
			} else {
				padl += sl.indicatorBlur
			}

			if sl.selectedColor == nil {
				return x, y, width, height
			}
			for i := x; i < x+sl.lineWidth-3; i++ {
				glyph := fullBar
				col := tcell.Color238
				//  ━━━━━━━━━━━━━━━━━╸╺━━━━━━━━━━━━━━━━━━━━━━
				if int(i) == mid-padl {
					col = *sl.selectedColor
					glyph = rightHalfBar
				}
				if int(i) >= mid-padl && int(i) <= mid+padr {
					col = *sl.selectedColor
				}
				if int(i) == mid+padr {
					glyph = leftHalfBar
					col = *sl.selectedColor
				}
				AppModel.scr.SetContent(
					x+int(i),
					y,
					glyph,
					nil,
					tcell.StyleDefault.Foreground(col).
						Background(theme.GetTheme().SidebarBackground),
				)
			}
			return x, y, width, height
		},
	)
	var scrollanim *Animation
	scrollanim = NewDynamicTargetAnimation(
		"scroll_line",
		btmLine,
		0,
		100,
		NewCallbackMutator(func(m MotionValues, i interface{}) bool {
			sl.x = m.X
			// if vel != math.Abs(m.Xvelocity) {
			//   vel = math.Abs(((vel * 2) + m.Xvelocity) / 3)
			//   zlog.Debug("velchange", zzlog.Float64("x", m.X), zzlog.Float64("xvel", m.Xvelocity))
			// }
			sl.vel = m.Xvelocity
			if sl.selectedColor == nil {
				return true
			}
			sl.indicatorPos = int(m.X)
			// AppModel.app.QueueUpdateDraw(func() {
			// 	AppModel.app.Draw(btmLine)
			// })
			return true
		}),
		func(a *Animation) (int, int, bool) {
			// sl.UpdatePos(cp *CoolorColorsPalette)
			AppModel.app.QueueUpdateDraw(func() {
				AppModel.app.Draw(btmLine)
			})
			if sl.indicatorPos != sl.indicatorTgtPos {
				scrollanim.Control.Play()
			}
			return sl.indicatorPos, sl.indicatorTgtPos, true
		},
	)
	// scrollanim.AutoNext = true
	// scrollanim.Loop = true
	scrollanim.Frames.IdleTime = 500 * time.Millisecond
	sl.Animation = scrollanim

	// updateTween := func(){
	// 	sl.Animation.GetCurrentFrame().Motions[0].UpdateTween(
	// 		float64(sl.indicatorPos),
	// 		float64(sl.indicatorTgtPos),
	// 	)
	// }
	scrollanim.Start()
	scrollanim.Control.Play()
	return sl
}

// tif := tview.NewInputField().SetPlaceholder("ballzdeep")
// tif.SetLabel(" SHIT ")
// tif.SetLabelWidth(10)
// tif.SetFormAttributes(10, tcell.ColorGreen, ct.GetTheme().SidebarBackground, tcell.ColorYellow, ct.GetTheme().SidebarBackground)
// tif.SetPlaceholderStyle(*ct.GetTheme().Get("input_placeholder"))
// tif.SetFieldStyle(*ct.GetTheme().Get("input_field"))
// tf := tview.NewForm()
// tf.AddDropDown(" OPTIONS ", []string{"one", "two"}, 0, func(option string, optionIndex int) {})
// tdd := tview.NewDropDown()
// tdd.SetBorder(true).SetBorderPadding(2, 2, 2, 2)
// tdd.AddOption("SHIT", func() {})
// tdd.AddOption(" YES ", func() {})
// tdd.AddOption(" NO ", func() {})
// tf.AddFormItem(tdd)
// tf.AddFormItem(tif)
// tf.AddButton(" NO ", func() {})
// tf.AddButton(" YES ", func() {})
// tf.AddButton(" TITS ", func() {})
// tf.SetFocus(0)
// mc.floater = NewFloater(tf)
// mc.pages.AddPage("floater", mc.floater.GetRoot(), true, true)
// mc.pages.ShowPage("floater")
// mc.app.SetFocus(mc.floater.GetRoot())
//
// vim: ts=2 sw=2 et ft=go
