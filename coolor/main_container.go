package coolor

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/gookit/goutil/structs"

	"github.com/digitallyserviced/coolors/coolor/anim"
	. "github.com/digitallyserviced/coolors/coolor/anim"
	"github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/shortcuts"
	"github.com/digitallyserviced/coolors/coolor/stack"
	"github.com/digitallyserviced/coolors/coolor/util"

	"github.com/digitallyserviced/coolors/theme"

	"github.com/digitallyserviced/coolors/status"
	ct "github.com/digitallyserviced/coolors/theme"
)

type MainContainer struct {
	floater        Floater
	current        tview.Primitive
	filetree       *PaletteFileTree
	palette        *CoolorMainPalette
	menu           *CoolorToolMenu
	mixer          *CoolorBlendPalette
	shades         *CoolorShadePalette
	scratch        *PaletteFloater
	info           *CoolorColorFloater
	editor         *CoolorColorEditor
	ScrollLine     *ScrollLine
	fileviewer     *CoolorFileView
	pages          *tview.Pages
	sidebar        *FixedFloater
	main           *tview.Flex
	app            *tview.Application
	options        *structs.Data
	conf           *HistoryDataConfig
	screen         *tcell.Screen
	preview        *Square
	inputSwallowed bool
	*stack.PageStack
	*tview.Flex
	*events.EventNotifier
	*shortcuts.Scope
}

// GetRef implements events.Referenced
func (mc *MainContainer) GetRef() interface{} {
	return mc
}

var (
	MainC  *MainContainer
	cpName = flag.String("palette", "", "Palette name to open")
)

func LoadPaletteFromArgs() *CoolorMainPalette {
	flag.Parse()
	if len(*cpName) > 0 {
		pal := GetStore().FindNamedPalette("stupefied_elastic_hypatia")
		if pal != nil {
			colors := NewCoolorColorsPaletteFromMeta(pal)
			events.Global.Notify(
				*events.Global.NewObservableEvent(events.PaletteLoadedEvent, "loaded_palette", colors, pal),
			)
			return colors
		} else {
			panic(fmt.Errorf("no palette named %s", *cpName))
		}
	}
	values := flag.Args()
	if len(values) > 0 {
		colors := NewCoolorPaletteFromCSSStrings(values)
		events.Global.Notify(
			*events.Global.NewObservableEvent(events.PaletteCreatedEvent, "startup_palette", colors, nil),
		)
		return colors
	}
	return nil
}

func DefaultStartupPalette() *CoolorMainPalette {
	colsize := 16
	colors := NewCoolorPaletteWithColors(GenerateRandomColors(colsize))

	events.Global.Notify(
		*events.Global.NewObservableEvent(events.PaletteCreatedEvent, "startup_palette", colors, nil),
	)
	return colors
}

func GetStartupPalette() *CoolorMainPalette {

	if colors := LoadPaletteFromArgs(); colors != nil {
		return colors
	}

	return DefaultStartupPalette()
}

func NewMainContainer(app *tview.Application) *MainContainer {
	tview.Borders = SimpleBorderStyle
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
		PageStack:     stack.NewPageStack(pgs),
	}

	addPluginOnLoad("register palette file view", func(pm *PluginsManager) {
		pm.Register(
			events.PluginEvents,
			events.NewAnonymousHandlerFunc(func(e events.ObservableEvent) bool {
				if e.Type.Is(events.PluginEvents) {
					if pe, ok := e.Ref.(*PluginEvent); ok {
						txt := fmt.Sprintf("Loaded Plugin %s", pe.name)
						status.NewStatusUpdate("action_str", txt)
						id := Notid("notif_")
						noti := anim.GetAnimator().GetAnimation(id)
						if noti == nil {
				noti = NewNotification(Notid("cpnotif_"), txt, InfoNotify)
						}
						noti.Control.Play()
					}
				}
				return true
			}),
		)
	})

	// events.Global.Register(
	// 	events.PaletteCreatedEvent,
	// 	events.NewAnonymousHandlerFunc(func(e events.ObservableEvent) bool {
	// 		if events.ObservableEventType(events.PaletteCreatedEvent | events.PaletteSavedEvent).
	// 			Is(e.Type) {
 //        fmt.Println(e)
	// 			ccpm, ok := e.Ref.(*CoolorColorsPaletteMeta)
	// 			if !ok {
	// 				return true
	// 			}
	// 			status.NewStatusUpdate("title", ccpm.Name)
	// 		}
	// 		return true
	// 	}),
	// )
	MainC.menu = NewCoolorColorMainMenu(app)

	MainC.palette = GetStartupPalette()
	MainC.conf.AddPalette("random", MainC.palette)
	MainC.menu.Init()
	// events.Global.Register(events.PaletteColorSelectedEvent, MainC.ScrollLine)
	MainC.palette.SetMenu(MainC.menu)
	MainC.palette.Register(events.PaletteColorSelectedEvent, MainC.menu)
	// MainC.editor = NewCoolorEditor(app, MainC.palette)
	MainC.preview = NewRecursiveSquare(MainC.palette.GetPalette(), 5)
	// MainC.fileviewer = NewFileViewer()
	// MainC.filetree = NewPaletteFileTree()
	MainC.Init()
	InitPlugins()
	return MainC
}

func (mc *MainContainer) SetupKeys() {
	mc.Scope.NewShortcut(
		"ESC",
		"cancel",
		tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone),
		func(i ...interface{}) bool {
			events.Global.Notify(
				*events.Global.NewObservableEvent(events.CancelledEvent, "esc", mc, nil),
			)
			if mc.PageStack.FrontPageStacked() {
				mc.Pop()
				name, pg := mc.PageStack.Pages.GetFrontPage()
				mc.app.SetFocus(pg)
				events.Global.Notify(
					*events.Global.NewObservableEvent(events.StatusEvent, name, mc, nil),
				)
			}
			return true
		},
	)
}
func (mc *MainContainer) Init() {
	mc.SetBackgroundColor(ct.GetTheme().SidebarBackground)
	mc.SetDirection(tview.FlexRow)
	MainC.ScrollLine = NewScrollLine()
	mc.AddItem(mc.pages, 0, 80, false)
	mc.AddItem(MainC.ScrollLine, 1, 0, false)
	// mc.pages.AddPage("editor", mc.editor, true, false)
	mc.pages.AddPage("preview", mc.preview, true, true)
	// mc.pages.AddPage("filetree", mc.filetree, false, false)
	mc.pages.AddAndSwitchToPage("palette", mc.palette, true)
	mc.pages.SetChangedFunc(func() {
		name, page := mc.pages.GetFrontPage()
		mc.current = page
		status.NewStatusUpdate("action", name)
	})
	// AppModel.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// if event.Key() == tcell.KeyF40 {
		// 	mc.inputSwallowed = true
		// } else if event.Key() == tcell.KeyF41 {
		// 	mc.inputSwallowed = false
		// }
		//
		// if mc.inputSwallowed {
		// 	// fmt.Printf("swallow:%+v %+v", mc.inputSwallowed, event)
		// 	p := MainC.app.GetFocus()
		// 	p.InputHandler()(event, func(tp tview.Primitive) {
		// 		MainC.app.SetFocus(tp)
		// 	})
		// 	return nil
		// }
		//
		// ch := event.Rune()
		// if event.Modifiers() == tcell.ModShift {
		// }
		// switch ch {
		// // case 'Q':
		// // 	AppModel.app.Stop()
		// // 	return nil
		// }
		// return event
	// })
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
			// if kp == tcell.KeyEscape {
			// 	events.Global.Notify(
			// 		*events.Global.NewObservableEvent(events.CancelledEvent, "esc", mc, nil),
			// 	)
			// 	// if mc.menu.Activated() != nil {
			// 	// 	mc.menu.Activated().Cancel()
			// 	// }
			// 	if mc.PageStack.FrontPageStacked() {
			// 		mc.Pop()
			// 		mc.app.SetFocus(page)
			// 	}
			//
			// 	// strutil.HasOneSub(name, []string{"floater", "sidebar", "colorpick", "info"})
			// 	// if !strutil.HasOneSub(name, []string{"palette", "mixer", "shades"}) {
			// 	// 	mc.Pop()
			// 	// }
			// 	// if name == "floater" {
			// 	// 	mc.pages.HidePage("floater")
			// 	// 	page.Blur()
			// 	// 	mc.pages.RemovePage("floater")
			// 	// 	mc.floater = nil
			// 	// }
			// 	return
			// }
			switch ch {
			case 'P':
				mc.Pop()
				shouldReturn := OpenColorPicker(mc)
				if shouldReturn {
					return
				}
			case 'F':
				mc.Pop()
				sb := OpenFavoritesView()
				_ = sb
				// mc.app.QueueUpdateDraw(func() {
				// 	if mc.sidebar == nil {
				// 		mc.pages.AddPage("sidebar", mc.sidebar.GetRoot(), true, true)
				// 		mc.pages.ShowPage("sidebar")
				// 		mc.app.SetFocus(mc.sidebar.GetRoot().Item)
				// 	} else {
				// 		// name, page := mc.pages.GetFrontPage()
				// 		if name == "sidebar" {
				// 			mc.pages.HidePage("sidebar")
				// 			page.Blur()
				// 			mc.pages.RemovePage("sidebar")
				// 			mc.sidebar = nil
				// 		} else {
				// 			mc.pages.ShowPage("sidebar")
				// 			mc.app.SetFocus(mc.sidebar.GetRoot().Item)
				// 		}
				// 		AppModel.helpbar.SetTable("sidebar")
				// 	}
				// })
			// case 't':
			//  
			// case 'i':
			// 	cc, _ := mc.palette.GetSelected()
			// 	if mc.info == nil {
			// 		mc.info = NewCoolorColorModal(cc)
			// 		mc.pages.AddPage("info", mc.info, true, false)
			// 		mc.pages.ShowPage("info")
			// 	} else {
			// 		// name, page := mc.pages.GetFrontPage()
			// 		if name == "info" {
			// 			mc.pages.HidePage("info")
			// 			page.Blur()
			// 		} else {
			// 			mc.pages.ShowPage("info")
			// 			mc.app.SetFocus(mc.info.Flex)
			// 			mc.info.Color.UpdateColor(cc)
			// 		}
			// 		AppModel.helpbar.SetTable("info")
			// 	}
			case '$':

				AppModel.app.GetComponentAt(20, 20)
				// mc.OpenTagView(mc.palette.CoolorColorsPalette)
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
				if mc.filetree == nil {
					mc.filetree = NewPaletteFileTree()
				}
				mc.Pop()
				mc.Push("filetree", mc.filetree, false)
				mc.app.SetFocus(mc.filetree)
			case 'S':
				mc.pages.SwitchToPage("shades") // .HidePage("editor")
				AppModel.helpbar.SetTable("shades")
			case 'N':
      OpenPalettesHistory()
			case 'M':
      OpenMainMenu("MENU")
				// mc.pages.SwitchToPage("mixer") // .HidePage("editor")
				// AppModel.helpbar.SetTable("mixer")
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
			case 'V':
				ani := NewFrameAnimator()
				ani.Start()
				ani.Control.Play()
			case 'u':
				noti := NewNotification(Notid("cpnotif_"), fmt.Sprintf("Favorited %s some really long wholly shit", NewCoolorColor("#fb3292").TVPreview()), NewNotificationStatus(" ", "#f7e6b5"))
				noti.Control.Play()

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
			// case "info":
			// 	HandleVimNavigableHorizontal(mc.info, ch, kp)
			// 	HandleVimNavigableVertical(mc.info, ch, kp)
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
				// dump.P(fmt.Sprintf("%s input handled", name))
				handler(event, setFocus)
			}
		},
	)
}

// 
//        ﱣ                             ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//     ﱣ ﱣ ﱣ                          ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//    ﱣ ﱣ                        ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                        ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                       ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                          ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                          ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                           ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                             ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                               ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ  ﯟ ⬢ ⬡
//                              ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                                         ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ⬢ ⬣  ⬡
//                                             ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                                 ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ⬢ ⬣  ⬡
//                     ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                      ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                 ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ⬢ ⬣  ⬡
//                    ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ
//                      ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ⬢ ⬣  ⬡
//                    ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ  
//                   ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ ﯟ  ⬢ ⬣   ⬡

func (mc *MainContainer) Draw(s tcell.Screen) {
	mc.Box.DrawForSubclass(s, mc)
	mc.Flex.Draw(s)
	mc.ScrollLine.Box.DrawForSubclass(s, mc.ScrollLine)
	mc.ScrollLine.Draw(s)
	// mc.menu.Draw(s)
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
	// fmt.Println(sl.indicatorTgtPos, sl.indicatorPos)
	sl.Animation.GetCurrentFrame().Motions[0].UpdateTween(
		float64(sl.indicatorPos),
		float64(sl.indicatorTgtPos),
	)
	// fmt.Println(
	// 	"delta",
	// 	math.Abs(float64(sl.indicatorTgtPos)-float64(sl.indicatorPos)),
	// 	sl.Animation.State.String(),
	// )
	if math.Abs(float64(sl.indicatorTgtPos)-float64(sl.indicatorPos)) > 1 {
		if sl.Animation.State.Is(events.AnimationNext) ||
			sl.Animation.State.Is(events.AnimationPaused) ||
			sl.Animation.State.Is(events.AnimationFinished) {
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
			//
			// ◢◣◤◥◧◸◹◺◿🭼🭽🭾🭿           ─━╸╺ ╼ ╾ ╍  ═ ┅    ╌╍═▪▫▬▭▮▯▰▱▶▷▸▹►▻▼▽▾▿◀◁◂◃◄◅◆◇◈◉◊○◌◍◍◎●◻◼◽◾◦🭶🭷🭸🭹🭺🭻🮂🮃🬂🬋🬭🬏🬇🬋🬃
			// ╶ ◂─◆◇◈─▸ ─╴ ╺▬━▰╸▬▭▮▯▰▱   ╺ ╼ ╾ ╍   ╍
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
				glyph := '─'
				col := tcell.Color238
				//  ━━━━━━━━━━━━━━━━━╸╺━━━━━━━━━━━━━━━━━━━━━━
				if int(i) == mid-padl-1 {
					glyph = '╼'
				}
				if int(i) == mid-padl {
					col = *sl.selectedColor
					glyph = rightHalfBar
				}
				if int(i) > mid-padl && int(i) < mid+padr {
					glyph = fullBar
					col = *sl.selectedColor
				}
				if int(i) == mid+padr {
					glyph = leftHalfBar
					col = *sl.selectedColor
				}
				if int(i) == mid+padr+1 {
					glyph = '╾'
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
		NewCallbackMutator(func(m *MotionValues, i interface{}) bool {
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
