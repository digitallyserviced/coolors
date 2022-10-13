package coolor

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/gookit/goutil/dump"
	"github.com/samber/lo"

	// "github.com/gookit/goutil/dump"

	// "github.com/digitallyserviced/coolors/theme"
	"github.com/digitallyserviced/coolors/theme"
	ct "github.com/digitallyserviced/coolors/theme"
	"github.com/digitallyserviced/coolors/tree"
)

type CoolorFileView struct {
	*tview.Grid
	Detail      *tview.Grid
	treeView    *tree.FileTree
	contentView *PaletteFileView
}

type PaletteFilePreviews struct {
	*tview.TextView
  currentPalette *CoolorColorsPalette
  fsnode *tree.FSNode
	header, footer         *tview.Frame
  headerText, footerText *tview.TextView
}

type PaletteFileTree struct {
	*tview.Flex
	treeView               *tree.FileTree
  preview *PaletteFilePreviews
}

// Hidden implements tview.Paged
func (*PaletteFileTree) Hidden(*tview.Pages) {
}

// Moved implements tview.Paged
func (*PaletteFileTree) Moved(*tview.Pages, tview.PageSentDirection) {
}

//ﮣ Shown implements tview.Paged
func (pft *PaletteFileTree) Shown(p *tview.Pages) {
	// MainC.app.QueueUpdateDraw(func() {
	configPath, _, _, _ := GetDataDirs()
  MainC.app.QueueUpdateDraw(func() {
    pft.treeView.LoadFiltered(configPath, PluginManager.SupportedFilenames())
  })
	MainC.app.SetFocus(pft.treeView)
	// pft.SetRect(p.GetInnerRect())
	// })
	//   x,y,w,h := p.GetInnerRect()
	//   dump.P(x,y,w,h)
	//   _,_,_,_ = x,y,w,h
	//   pft.SetRect(x,y,w/4,h)
	//   MainC.app.Draw(pft)
	// })
}

// // Focus implements tview.Primitive
func (pft *PaletteFileTree) Focus(delegate func(p tview.Primitive)) {
	// configPath, _, _, _ := GetDataDirs()
	// pft.treeView.LoadFiltered(configPath, PluginManager.SupportedFilenames())
	delegate(pft.treeView)
}

// InputHandler implements tview.Primitive
func (pft *PaletteFileTree) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return pft.WrapInputHandler(
		func(ek *tcell.EventKey, f func(p tview.Primitive)) {
			pft.treeView.InputHandler()(ek, f)
		},
	)
}

//
// // IsVisible implements tview.Primitive
// func (*PaletteFileTree) IsVisible() bool {
// 	panic("unimplemented")
// }
//
// // MouseHandler implements tview.Primitive
// func (*PaletteFileTree) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
// 	panic("unimplemented")
// }

// // Draw implements tview.Primitive
// func (pft *PaletteFileTree) DrawForSubclass(screen tcell.Screen, p tview.Primitive) {
//   pft.treeView.Box.DrawForSubclass(screen, p)
// }
func (pft *PaletteFileTree) Draw(screen tcell.Screen) {
	w, h := screen.Size()
	pft.SetRect(0, 1, w, h-2)
	pft.Box.DrawForSubclass(screen, pft)
	// pft.treeView.Box.DrawForSubclass(screen, pft.treeView)
	pft.Flex.Draw(screen)
}

// GetRect implements tview.Primitive
func (pft *PaletteFileTree) GetRect() (int, int, int, int) {
	return pft.treeView.GetRect()
}

// SetRect implements tview.Primitive
func (pft *PaletteFileTree) SetRect(x int, y int, width int, height int) {
	pft.Flex.SetRect(x, y, width/4, height)
	// pft.treeView.SetRect(x, y, width/4, height)
}

func NewPaletteFilePreview() *PaletteFilePreviews {
	headerText := tview.NewTextView()
	headerText.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	footerText := tview.NewTextView()
	footerText.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
  footerText.SetBackgroundColor(theme.GetTheme().SidebarBackground)
	header := tview.NewFrame(headerText)
	footer := tview.NewFrame(footerText)
	header.SetBorders(0, 0, 1, 0, 0, 0)
	header.SetBorder(false).
		SetBorderColor(theme.GetTheme().HeaderBackground).
		SetBorderPadding(1, 0, 0, 0)
  headerText.SetText(" ")

	footer.SetBorderPadding(0, 0, 1, 1)
  footer.SetBackgroundColor(theme.GetTheme().SidebarBackground)
  footer.SetBorder(true).SetBorderVisible(false)
	footer.SetBorders(1, 1, 1, 1, 0, 0)

		// footer.SetBorders(0, 0, 1, 1, 0, 0)

	pfp := &PaletteFilePreviews{
		// TextView:       tview.NewTextView(),
		currentPalette: nil,
		header:     header,
		footer:     footer,
		headerText: headerText,
    footerText: footerText,
	}

	// pfp.SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	// pfp.SetBorderVisible(false)
	// pfp.SetBorder(true).
	// 	SetBorderPadding(0, 0, 0, 0).
	// 	SetBorderColor(theme.GetTheme().SidebarLines)
	return pfp
}
func (pfp *PaletteFilePreviews) UpdatePreview(fsnode *tree.FSNode) {
  pfp.SetFile(fsnode)
		if fsnode.IsDir {
			// pfp.Clear()
			return
		}
		var cols *CoolorColorsPalette
		i := []string{}
		b := LoadFile(fsnode.Path)
		mapd := b.All()
		colors := getColorsFromMap(mapd)
		cols = CoolorStrings(colors).GetPalette()
		i = append(
			i,
			fmt.Sprintf("  %s ", formatSize(fsnode.Size)),
		)
		i = append(
			i,
			fmt.Sprintf("  %d ", len(cols.Colors)),
		)
		if cols.Len() != 0 {
			cols.Sort()
			barChars := cols.MakeSquarePalette(false)
    x, y, w, h := pfp.footer.GetInnerRect()
			_, _, _, _ = x, y, w, h
			lbc := len(barChars)
			rows := math.Ceil(float64(lbc) / (float64(w) - 2))
			lines := strings.Join(
				lo.Map[[]string, string](
					lo.Chunk[string](barChars, lbc/int(rows)),
					func(s []string, i int) string {
						return strings.Join(s, "")
					},
				),
				"\n",
			)
    pfp.footerText.SetText(lines)
			// i = append(i, lines)
		}
		// dump.P(x, y, w, h, rows, lbc, lines)
    pfp.header.AddText(strings.Join(i, " | "), false, tview.AlignCenter, theme.GetTheme().InfoLabel)
		// 🬫 🬛       🬴     🬻 🬺     🬨🬕     🬷🬲   🬪   🬜
    // pfp.header.AddText(fmt.Sprintf("[blue:black:-]🬸[-:-:-][black:blue:-] %s [-:-:-][blue:black:-]🬴[-:-:-]", fsnode.Name), true, AlignCenter, theme.GetTheme().InfoLabel)
}
	// icon := "  "
	// if n.IsDir {
	// 	if n.IsExpanded() {
	// 		icon = " ﱮ"
	// 	} else {
	// 		icon = " "
	// 	}
	// }
func (pfp *PaletteFilePreviews) SetFile(fsnode *tree.FSNode) {
  pfp.fsnode = fsnode
  pfp.header.Clear()
  if pfp.fsnode.IsDir {
    pfp.header.AddText(fmt.Sprintf("[blue:black:-]🬸[-:-:-][black:blue:-] ﱮ  %s [-:-:-][blue:black:-]🬴[-:-:-]", fsnode.Name), true, tview.AlignCenter, theme.GetTheme().InfoLabel)
  } else {
    pfp.header.AddText(fmt.Sprintf("[blue:black:-]🬸[-:-:-][black:blue:-]   %s [-:-:-][blue:black:-]🬴[-:-:-]", fsnode.Name), true, tview.AlignCenter, theme.GetTheme().InfoLabel)
  }

}

func NewPaletteFileTree() *PaletteFileTree {
  pfp := NewPaletteFilePreview()
	tv := tree.NewFileTree(tree.GetTheme())
	pft := &PaletteFileTree{
		Flex:       tview.NewFlex(),
		treeView:   tv,
		preview: pfp,
	}

  tv.Box.SetDontClear(false)
	tv.Box.SetBackgroundColor(ct.GetTheme().SidebarBackground)
	pft.Box.SetDontClear(true)
	pft.Flex.SetDirection(tview.FlexRow)
	pft.Flex.AddItem(pfp.header, 5, 0, false)
	pft.Flex.AddItem(tv, 0, 60, true)
	pft.Flex.AddItem(pfp.footer, 0, 25, false)

	// configPath, _, _, _ := GetDataDirs()
	// tv.Load(configPath)
	cwd, _ := os.Getwd()
	rt := tree.NewVirtualNode("/", "", "")
	lr := tree.NewVirtualNode("Local", "", cwd)
	vr := tree.NewVirtualNode("Plugins", " ", "")
  vr.Children = PluginManager.GetTreeEntries()
	tv.SetRoot(rt)
  tv.SetVirtualRoot(vr)
  tv.SetLocalRoot(lr)

	tv.OnChanged(func(node *tree.FSNode) {
    pft.preview.UpdatePreview(node)
	})

	tv.OnSelect(func(node *tree.FSNode) {
	})

	tv.OnOpen(func(node *tree.FSNode) {
	})

	MainC.app.SetAfterDrawFunc(func(screen tcell.Screen) {
		var x func(*tree.FileTree)
		for len(pft.treeView.AfterDraw) > 0 {
			x, pft.treeView.AfterDraw = pft.treeView.AfterDraw[0], pft.treeView.AfterDraw[1:]
			x(pft.treeView)
		}
	})
	return pft
}

func NewFileViewer() *CoolorFileView {
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
	tv := tree.NewFileTree(theme)

	cfv := &CoolorFileView{
		Grid:        topgrid,
		Detail:      rightgrid,
		treeView:    tv,
		contentView: NewPaletteFileView(tt),
	}

	fm := tview.NewFocusManager(func(p tview.Primitive) {
		MainC.app.SetFocus(p)
	})
	fm.Add(cfv.treeView, cfv.contentView.colorsView, cfv.contentView.contentView)

	cfv.SetFocusManager(fm)

	cfv.treeView.SetNextFocusableComponents(
		tview.Right,
		cfv.contentView.colorsView,
	)
	cfv.treeView.SetNextFocusableComponents(
		tview.Down,
		cfv.contentView.contentView,
	)
	cfv.contentView.contentView.SetNextFocusableComponents(
		tview.Left,
		cfv.treeView,
	)
	cfv.contentView.colorsView.SetNextFocusableComponents(
		tview.Down,
		cfv.contentView.contentView,
	)
	cfv.contentView.colorsView.SetNextFocusableComponents(
		tview.Left,
		cfv.treeView,
	)
	cfv.contentView.contentView.SetNextFocusableComponents(
		tview.Up,
		cfv.contentView.colorsView,
	)

	cfv.treeView.OnChanged(func(fsnode *tree.FSNode) {
		if fsnode.IsDir {
			return
		}
		MainC.app.QueueUpdateDraw(func() {
			cfv.contentView.SetPreview(fsnode)
		})
	})

	cfv.treeView.OnSelect(func(node *tree.FSNode) {
		if node.IsDir {
			return
		}
		// cfv.contentView.SetPreview(fsnode)
	})

	cfv.treeView.OnOpen(func(node *tree.FSNode) {
		if node.IsDir {
			return
		}
		MainC.app.QueueUpdateDraw(func() {
			x, y, w, h := cfv.treeView.Box.GetRect()
			collapsed := cfv.treeView.GetCollapsed()
			if !collapsed {
				fmt.Println("collapsing", fmt.Sprintf("%d %d %d %d", x, y, w, h))
				cfv.Grid.SetColumns(2, -98)
				// MainC.app.SetFocus(cfv.contentView.contentView)
			} else {
				fmt.Println("expanding", fmt.Sprintf("%d %d %d %d", x, y, w, h))
				cfv.Grid.SetColumns(-15, -85)
			}
			cfv.treeView.SetCollapsed(!collapsed)
		})
		// go func() {
		// 	exec.Command("open", node.Path)
		// }()
	})

	configPath, _, _, _ := GetDataDirs()
	cfv.treeView.LoadFiltered(configPath, PluginManager.SupportedFilenames())

	MainC.app.SetAfterDrawFunc(func(screen tcell.Screen) {
		var x func(*tree.FileTree)
		for len(cfv.treeView.AfterDraw) > 0 {
			x, cfv.treeView.AfterDraw = cfv.treeView.AfterDraw[0], cfv.treeView.AfterDraw[1:]
			x(cfv.treeView)
		}
	})

	cfv.Grid.
		AddItem(cfv.treeView, 0, 0, 1, 1, 1, -15, true).
		AddItem(cfv.Detail, 0, 1, 1, 1, 1, -85, false)

	cfv.Detail.
		AddItem(cfv.contentView, 0, 0, 1, 1, 1, 1, false)

	return cfv
}

func (cfv *CoolorFileView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cfv.WrapInputHandler(
		func(event *tcell.EventKey, f func(p tview.Primitive)) {
			if event.Modifiers() == tcell.ModShift {
				if DirectionalFocusHandling(event, AppModel.app) == nil {
					return
				}
			}
			cfv.Grid.InputHandler()(event, f)
		},
	)
}

func (cfv *CoolorFileView) Focus(delegate func(p tview.Primitive)) {
	// cfv.treeView.Focus(delegate)
}

//
// // Focus implements tview.Primitive
// func (*PaletteFileTree) Focus(delegate func(p tview.Primitive)) {
// 	panic("unimplemented")
// }
//
// // GetFocusable implements tview.Primitive
// func (*PaletteFileTree) GetFocusable() tview.Focusable {
// 	panic("unimplemented")
// }
//
// // GetParent implements tview.Primitive
// func (*PaletteFileTree) GetParent() tview.Primitive {
// 	panic("unimplemented")
// }
//
//
// // HasFocus implements tview.Primitive
// func (*PaletteFileTree) HasFocus() bool {
// 	panic("unimplemented")
// }
//
//
// // NextFocusableComponent implements tview.Primitive
// func (*PaletteFileTree) NextFocusableComponent(tview.FocusDirection) tview.Primitive {
// 	panic("unimplemented")
// }
//
// // OnPaste implements tview.Primitive
// func (*PaletteFileTree) OnPaste([]rune) {
// 	panic("unimplemented")
// }
//
// // SetOnBlur implements tview.Primitive
// func (*PaletteFileTree) SetOnBlur(handler func()) {
// 	panic("unimplemented")
// }
//
// // SetOnFocus implements tview.Primitive
// func (*PaletteFileTree) SetOnFocus(handler func()) {
// 	panic("unimplemented")
// }
//
// // SetParent implements tview.Primitive
// func (*PaletteFileTree) SetParent(tview.Primitive) {
// 	panic("unimplemented")
// }
// // SetVisible implements tview.Primitive
// func (*PaletteFileTree) SetVisible(bool) {
// 	panic("unimplemented")
// }
