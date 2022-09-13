package coolor

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/digitallyserviced/coolors/status"
	ct "github.com/digitallyserviced/coolors/theme"
	"github.com/digitallyserviced/coolors/tree"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/structs"
	// "github.com/rivo/tview"
	// "github.com/josa42/term-finder/tree"
)

type MainContainer struct {
	*tview.Flex
	screen     *tcell.Screen
	editor     *CoolorColorEditor
	menu       *CoolorToolMenu
	palette    *CoolorMainPalette
	mixer      *CoolorBlendPalette
	shades     *CoolorShadePalette
	scratch    *PaletteFloater
	info       *CoolorColorFloater
	floater    Floater
	preview    *Square
	fileviewer *CoolorFileView
	pages      *tview.Pages
	sidebar    *FixedFloater
	main       *tview.Flex
	app        *tview.Application
	options    *structs.MapDataStore
	conf       *HistoryDataConfig
	*eventNotifier
}

var MainC *MainContainer

type CoolorFileView struct {
	*tview.Grid
	Detail      *tview.Grid
	treeView    *tree.FileTree
	contentView *PaletteFileView
}

func (mc *MainContainer) NewFileViewer() *CoolorFileView {
	setupLogging()
	pwd, _ := os.Getwd()
	log.Printf("open: %s", pwd)

	theme := tree.GetTheme()
	tt := ct.GetTheme()

	topgrid := tview.NewGrid().
		SetBordersColor(theme.Border).
		SetBorders(theme.Border != 0).
		SetColumns(25, 0)

	rightgrid := tview.NewGrid().
		SetBordersColor(theme.Border).
		SetBorders(theme.Border != 0).
		SetRows(0)

	cfv := &CoolorFileView{
		Grid:        topgrid,
		Detail:      rightgrid,
		treeView:    tree.NewFileTree(theme),
		contentView: NewPaletteFileView(tt),
	}

	cfv.treeView.OnChanged(func(fsnode *tree.FSNode) {
		MainC.app.QueueUpdateDraw(func() {
			cfv.contentView.SetPreview(fsnode)
		})
	})

	cfv.treeView.OnSelect(func(node *tree.FSNode) {
		// cfv.contentView.SetPreview(fsnode)
	})

	cfv.treeView.OnOpen(func(node *tree.FSNode) {
		go func() {
			exec.Command("open", node.Path)
		}()
	})

	configPath, _, _, _ := GetDataDirs()
	cfv.treeView.Load(configPath)

	MainC.app.SetAfterDrawFunc(func(screen tcell.Screen) {
		var x func()
		for len(cfv.treeView.AfterDraw) > 0 {
			x, cfv.treeView.AfterDraw = cfv.treeView.AfterDraw[0], cfv.treeView.AfterDraw[1:]
			x()
		}
	})

	cfv.Grid.
		AddItem(cfv.treeView, 0, 0, 1, 20, 1, 1, true).
		AddItem(cfv.Detail, 0, 20, 1, 60, 1, 1, false)

	cfv.Detail.
		AddItem(cfv.contentView, 0, 0, 1, 1, 1, 1, false)

	return cfv
}

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
		colsize = 8
	}

	var colors *CoolorMainPalette

	if len(values) > 0 {
		colors = NewCoolorPaletteFromCssStrings(values)
	} else {
		colors = NewCoolorPaletteWithColors(GenerateRandomColors(colsize))
	}

  // colors.UpdateHash()
  // // p := colors.Hash
  // p := 1728741516
  //
  // log.Printf("%s", rg.WithSeed(int64(p)).GenerateName(2))
	return colors
}

func NewMainContainer(app *tview.Application) *MainContainer {
	MainC = &MainContainer{
		Flex:          tview.NewFlex(),
		floater:       nil,
		pages:         tview.NewPages(),
		main:          tview.NewFlex(),
		app:           app,
		options:       ActionOptions,
		conf:          NewPaletteHistoryFile(),
		eventNotifier: NewEventNotifier("main"),
	}

	// if GetStore().Store != nil {
	//   MainC.eventNotifier.Register(ColorSeentEvent, GetStore().MetaService)
	// }

	MainC.eventNotifier.Register(
		PaletteColorModifiedEvent|PaletteColorRemovedEvent|PaletteCreatedEvent|PaletteMetaUpdatedEvent|PaletteSavedEvent|SelectedEvent|ColorEvent|ColorSeentEvent,
		GetStore().MetaService,
	)
	// MainC.eventNotifier.Register(SelectedEvent | ColorEvent | ColorSeentEvent, GetStore().MetaService)

	MainC.menu = NewCoolorColorMainMenu(app)
	MainC.palette = CreateStartupPalette()
	MainC.conf.AddPalette("random", MainC.palette)
	MainC.menu.Init()
	MainC.palette.SetMenu(MainC.menu)
	MainC.palette.Register(PaletteColorSelectedEvent, MainC.menu)
	MainC.editor = NewCoolorEditor(app, MainC.palette)
	MainC.preview = NewRecursiveSquare(MainC.palette.GetPalette(), 5)
	MainC.fileviewer = MainC.NewFileViewer()
	MainC.Init()
	return MainC
}

func (mc *MainContainer) CloseConfig() {
	mc.conf.ref.Close()
}

func (mc *MainContainer) Init() {
	mc.SetBackgroundColor(ct.GetTheme().SidebarBackground)
	mc.SetDirection(tview.FlexColumn)
	mc.AddItem(mc.pages, 0, 80, false)
	mc.pages.AddPage("editor", mc.editor, true, false)
	mc.pages.AddPage("fileviewer", mc.fileviewer, true, false)
	mc.pages.AddPage("preview", mc.preview, true, true)
	mc.pages.AddAndSwitchToPage("palette", mc.palette, true)
	mc.pages.SetChangedFunc(func() {
		name, _ := mc.pages.GetFrontPage()
		status.NewStatusUpdate("action", name)
	})
}

type AnonymousObserver struct {
	*eventObserver
	Callbacks []Observer
}

type AnonymousHandler struct {
	Callback func(e ObservableEvent) bool
}

func (ah *AnonymousHandler) HandleEvent(e ObservableEvent) bool {
	fmt.Println(e)
	ah.Callback(e)
	return true
}

func (ah *AnonymousObserver) ObserverFunc(
	n Notifier,
	t ObservableEventType,
	f func(e ObservableEvent) bool,
) *AnonymousObserver {
	ahh := &AnonymousHandler{
		Callback: func(e ObservableEvent) bool {
			return f(e)
		},
	}
	n.Register(t, ahh)
	ah.Callbacks = append(ah.Callbacks, ahh)
	return ah
}

func NewAnonymousHandler(callbacks []Observer) *AnonymousObserver {
	ah := &AnonymousObserver{
		eventObserver: NewEventObserver("anon"),
	}
	return ah
}

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
			case 'L':
				mc.app.QueueUpdateDraw(func() {
					if mc.sidebar == nil {
						// ccsTwo := NewCoolorColorSwatch(
						// 	func(cs *CoolorColorSwatch) *CoolorColorsPalette {
						// 		cp := NewCoolorPaletteWithColors(
						// 			GenerateRandomColors(20),
						// 		).GetPalette()
						// 		return cp
						// 	},
						// )
						// tv.AddTab(NewTabView(" Favorites", tv.TakeNext(), ccsTwo))
						// tv.AddTab(NewTabView("ﱬ Stream", tv.TakeNext(), ccsTwo))
						// tv.AddTab(NewTabView("Alter Fuck",tv.TakeNext(), ccs111))
						tv := NewTabbedView()
						// tv.tabsView.SetBackgroundColor(tcell.Color101)
						mc.sidebar = NewFixedFloater("QuickColors", tv)
						ccs := NewCoolorColorSwatch(
							func(cs *CoolorColorSwatch) *CoolorColorsPalette {
								cp := NewCoolorColorsPalette()
								// cp.Colors = *GetStore().MetaService.RecentColors
								// fmt.Println(len(cp.Colors))
								return cp
							},
						)
						ccstv := NewTabView(" Recents", tv.TakeNext(), ccs)
						ccstv.SetBackgroundColor(ct.GetTheme().GrayerBackground)
						changed := NewAnonymousHandler(
							nil,
						).ObserverFunc(ccs, ColorSelectedEvent, func(e ObservableEvent) bool {
							col, ok := e.Ref.(*CoolorColor)
							if ok {
								ccstv.Frame.Clear().
									AddText(col.GetMeta().String(), false, AlignCenter, tcell.ColorRed)
								return true
							}
							return true
						})
						selected := NewAnonymousHandler(
							nil,
						).ObserverFunc(ccs, ColorSelectionEvent, func(e ObservableEvent) bool {
							col, ok := e.Ref.(*CoolorColor)
							if ok {
								MainC.palette.AddCoolorColor(col)
								return true
							}
							return true
						})
						// ccs.UpdateItems()
						// ccs.Register(ColorSelectionEvent, selected)
						// ccs.Register(ColorSelectedEvent, changed)
						_, _ = changed, selected
						tv.AddTab(ccstv)
						tv.UpdateView()
						mc.pages.AddPage("sidebar", mc.sidebar.GetRoot(), true, true)
						mc.pages.ShowPage("sidebar")
						mc.app.SetFocus(mc.sidebar.GetRoot().Item)
					} else {
						name, page := mc.pages.GetFrontPage()
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
					mc.info = NewCoolorColorInfoFloater(cc)
					mc.pages.AddPage("info", mc.info, true, false)
					mc.pages.ShowPage("info")
				} else {
					name, page := mc.pages.GetFrontPage()
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
			case 'C':
				if mc.scratch == nil {
					p := NewCoolorPaletteWithColors(GenerateRandomColors(8))
					cpc := NewScratchPaletteFloater(p.GetPalette())
					mc.scratch = cpc
					mc.pages.AddPage("scratch", mc.scratch, true, false)
					mc.pages.ShowPage("scratch")
				} else {
					name, page := mc.pages.GetFrontPage()
					if name == "scratch" {
						mc.pages.HidePage("scratch")
						page.Blur()
					} else {
						mc.pages.ShowPage("scratch")
						mc.app.SetFocus(mc.scratch.Palette)
					}
					AppModel.helpbar.SetTable("scratch")
				}
			case 'F':
				mc.pages.SwitchToPage("fileviewer") // .HidePage("editor")
				AppModel.helpbar.SetTable("fileviewer")
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
				HandleVimNavigableHorizontal(mc.scratch.Palette.Palette.Palette, ch, kp)
				HandleCoolorSelectable(mc.scratch.Palette.Palette.Palette, ch, kp)
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

func GetColorTagItems() []TagListItem {
	tis := make([]TagListItem, 0)
	tis = append(
		tis,
		*NewTagListItem("foreground", "default foreground color", 'f'),
	)
	return tis
}

func (cfv *CoolorFileView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cfv.WrapInputHandler(cfv.treeView.InputHandler())
}

func (cfv *CoolorFileView) Focus(delegate func(p tview.Primitive)) {
	// cfv.treeView.Focus(delegate)
}

func (mc *MainContainer) NewShades(base *CoolorColor) {
	mc.app.QueueUpdateDraw(func() {
		if mc.pages.HasPage("shades") && mc.shades != nil {
			mc.shades.UpdateColors(base)
		} else {
			mc.shades = BlankCoolorShadePalette(base, 8)
			mc.shades.SetMenu(MainC.menu)
			mc.pages.AddPage("shades", mc.shades, true, true)
		}
		mc.pages.SwitchToPage("shades") // .HidePage("editor")
	})
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
