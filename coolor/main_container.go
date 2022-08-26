package coolor

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/digitallyserviced/coolors/tree"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
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
	shades      *CoolorShadePalette
	scratch    *PaletteFloater
	info       *CoolorColorFloater
	preview    *Square
	fileviewer *CoolorFileView
	pages      *tview.Pages
	main       *tview.Flex
	app        *tview.Application
	options    *structs.MapDataStore
	conf       *HistoryDataConfig
}

type CoolorFileView struct {
	*tview.Grid
	Detail      *tview.Grid
	treeView    *tree.FileTree
	contentView *PaletteFileView
}

func setupLogging() func() error {
	f, _ := os.OpenFile("dumps", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
	// f, _ := os.OpenFile(os.DevNull, os.O_RDWR|os.O_APPEND, 0666)

	log.SetOutput(f)

	return f.Close
}

func (mc *MainContainer) NewFileViewer() *CoolorFileView {
	setupLogging()
	pwd, _ := os.Getwd()
	log.Printf("open: %s", pwd)

	theme := tree.GetTheme()
	tt := GetTheme()

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
		cfv.contentView.SetPreview(fsnode)
	})

	cfv.treeView.OnSelect(func(node *tree.FSNode) {
		if !node.IsDir {
		}
	})

	cfv.treeView.OnOpen(func(node *tree.FSNode) {
		go func() {
			exec.Command("open", node.Path).Run()
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
		AddItem(cfv.treeView, 0, 0, 1, 1, 0, 10, true).
		AddItem(cfv.Detail, 0, 1, 1, 1, 0, 50, true)

	cfv.Detail.
		AddItem(cfv.contentView, 0, 0, 1, 1, 0, 0, false)

	return cfv
}

var MainC *MainContainer

func CreateStartupPalette() *CoolorMainPalette {
	flag.Parse()
	values := flag.Args()

	colsize := 0
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

	return colors
}

func NewMainContainer(app *tview.Application) *MainContainer {
	MainC = &MainContainer{
		Flex:       tview.NewFlex(),
		editor:     &CoolorColorEditor{},
		menu:       &CoolorToolMenu{},
		palette:    &CoolorMainPalette{},
		mixer:      &CoolorBlendPalette{},
		preview:    &Square{},
		fileviewer: &CoolorFileView{},
		pages:      tview.NewPages(),
		main:       tview.NewFlex(),
		app:        app,
		options:    ActionOptions,
		conf:       NewPaletteHistoryFile(),
	}

	MainC.menu = NewCoolorColorMainMenu(app)
	MainC.palette = CreateStartupPalette()
	MainC.conf.AddPalette("random", MainC.palette)
	MainC.menu.Init()
	MainC.palette.SetMenu(MainC.menu)
	MainC.editor = NewCoolorEditor(app, MainC.palette)
	MainC.preview = NewRecursiveSquare(MainC.palette.GetPalette(), 5)
	MainC.fileviewer = MainC.NewFileViewer()
	MainC.Init()
	return MainC
}

func (mc *MainContainer) CloseConfig() {
	mc.conf.PaletteFile.ref.Close()
}

func (mc *MainContainer) Init() {
	mc.Flex.SetDirection(tview.FlexColumn)
	mc.Flex.AddItem(mc.pages, 0, 80, false)
	mc.pages.AddPage("editor", mc.editor, true, false)
	mc.pages.AddPage("fileviewer", mc.fileviewer, true, false)
	mc.pages.AddPage("preview", mc.preview, true, true)
	mc.pages.AddAndSwitchToPage("palette", mc.palette, true)
	mc.pages.SetChangedFunc(func() {
		// mc.mixer.menu
	})
}

// InputHandler returns the handler for this primitive.
func (mc *MainContainer) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return mc.pages.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		ch := event.Rune()
		kp := event.Key()
		switch {
		case ch == 'i':
			cc, _ := mc.palette.GetSelected()
			if mc.info == nil {
				mc.info = NewCoolorColorFloater(cc)
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
		case ch == 'C':
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
		case ch == 'F':
			mc.pages.SwitchToPage("fileviewer") // .HidePage("editor")
      AppModel.helpbar.SetTable("fileviewer")
		case ch == 'S':
			mc.pages.SwitchToPage("shades") // .HidePage("editor")
      AppModel.helpbar.SetTable("shades")
		case ch == 'M':
			mc.pages.SwitchToPage("mixer") // .HidePage("editor")
      AppModel.helpbar.SetTable("mixer")
		case ch == 'p':
			mc.pages.SwitchToPage("palette") // .HidePage("editor")
      AppModel.helpbar.SetTable("palette")
		case ch == 'Y':
			mc.pages.ShowPage("preview").SendToFront("preview") //.HidePage("palette")
      AppModel.helpbar.SetTable("preview")
			x, y, w, h := mc.pages.GetInnerRect()
			_, _, _, _ = x, y, w, h
			mc.preview.SetRect(x, y, w/2, h/2)
			mc.preview.TopInit(8)
		case ch == 'e':
			mc.pages.SwitchToPage("editor") //.HidePage("palette")
      AppModel.helpbar.SetTable("editor")
		case kp == tcell.KeyEscape:
			if mc.menu.Activated() != nil {
				mc.menu.Activated().Cancel()
			}
		}

		name, page := mc.pages.GetFrontPage()
		if page == nil {
			return
		}
		switch {
		case name == "shades":
			HandleVimNavigableHorizontal(mc.shades, ch, kp)
			HandleCoolorSelectable(mc.shades, ch, kp)
			// dump.P(fmt.Sprintf("%s horiz input handled", name))
		case name == "mixer":
			HandleVimNavigableHorizontal(mc.mixer, ch, kp)
			HandleCoolorSelectable(mc.mixer, ch, kp)
			// dump.P(fmt.Sprintf("%s horiz input handled", name))
		case name == "palette":
			HandleVimNavigableHorizontal(mc.palette, ch, kp)
			HandleCoolorSelectable(mc.palette, ch, kp)
			// dump.P(fmt.Sprintf("%s horiz input handled", name))
		case name == "info":
			HandleVimNavigableHorizontal(mc.info, ch, kp)
			HandleVimNavigableVertical(mc.info, ch, kp)
			// HandleCoolorSelectable(mc.info, ch, kp)
		case name == "scratch":
			HandleVimNavigableHorizontal(mc.scratch.Palette.Palette.Palette, ch, kp)
			HandleCoolorSelectable(mc.scratch.Palette.Palette.Palette, ch, kp)
		}

		if name == "mixer" || name == "palette" || name == "shades" {
			// dump.P(fmt.Sprintf("%s vertical input handled", name))
			HandleVimNavigableVertical(MainC.menu, ch, kp)
		}

		if handler := page.InputHandler(); handler != nil {
			handler(event, setFocus)
		}
	})
}

func (cfv *CoolorFileView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cfv.treeView.InputHandler()
}

func (cfv *CoolorFileView) Focus(delegate func(p tview.Primitive)) {
	cfv.treeView.Focus(delegate)
}

func (mc *MainContainer) NewShades(base *CoolorColor) {
	if mc.pages.HasPage("shades") && mc.shades != nil {
		// mc.mixer.UpdateColors(start, end)
    mc.shades.UpdateColors(base)
	} else {
		mc.shades = BlankCoolorShadePalette(base, 8)
		mc.shades.SetMenu(MainC.menu)
		mc.pages.AddPage("shades", mc.shades, true, false)
	}
	// name, page := mc.pages.GetFrontPage()
	// if page != nil {
	// 	mc.pages.HidePage(name)
	// }
	mc.pages.SwitchToPage("shades") // .HidePage("editor")
	// mc.pages.HidePage("palette")
	mc.palette.Blur()
	mc.palette.ColorContainer.Blur()
	mc.app.SetFocus(mc.shades)
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

func (mc *MainContainer) GetPalette() *CoolorPalette {
	return mc.palette.GetPalette()
}

// vim: ts=2 sw=2 et ft=go
