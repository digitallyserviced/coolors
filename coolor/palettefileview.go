package coolor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
	"github.com/gookit/goutil/errorx"

	"github.com/digitallyserviced/coolors/theme"
	"github.com/digitallyserviced/coolors/tree"
)

var _ tview.Primitive = &PaletteFileView{}

type PaletteFileView struct {
	*tview.Box
	theme       *theme.Theme
	view        *tview.Grid
	infoGrid    *tview.Grid
	contentGrid *tview.Grid
	topbarView  *tview.TextView
	infoView    *tview.TextView
	colorsView  *ColorsView
	contentView *tview.TextView
}

type ColorsView struct {
	*tview.Grid
	palettes []Palette
}


func (cv *ColorsView) Clear() {
	cv.palettes = nil
	cv.Grid.Clear()
}

func (cv *ColorsView) SetPalettes(pdc *HistoryDataConfig) {
	MainC.app.QueueUpdateDraw(func() {
		cv.Clear()
		pals := pdc.GetPalettes()

		// dump.P(pals)
		for _, v := range pals {
			pal := pdc.LoadPalette(v)
			cv.palettes = append(cv.palettes, pal.GetPalette())
		}
	})
}

func NewColorsView() *ColorsView {
	colorsView := tview.NewGrid()
	// colorsView.SetBorderPadding(0, 0, 1, 1)

	return &ColorsView{
		Grid:     colorsView,
		palettes: make([]Palette, 0),
	}
}

func formatSize(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func NewPaletteFileView(theme *theme.Theme) *PaletteFileView {
	view := tview.NewGrid()
	box := tview.NewBox()
	// view.SetRows(3, 5, 0)

	topbar := tview.NewTextView()
	topbar.SetBorderPadding(1,1,0,0)
	topbar.ScrollTo(0, 0)
	topbar.SetBorder(true)
	topbar.SetMaxLines(1)
  topbar.SetDontClear(true)
	topbar.SetTextAlign(AlignCenter)
	topbar.SetBorderColor(theme.Border)
  // topbar.SetBorderAttributes(tcell.AttrReverse)
  topbar.SetBorderVisible(false)
  topbar.SetBackgroundColor(theme.SidebarBackground)
	topbar.SetDynamicColors(true)
	topbar.SetRegions(true)

	info := tview.NewTextView()
	info.SetDynamicColors(true)
	info.SetRegions(true)
  info.SetBackgroundColor(theme.SidebarBackground)

	colorsView := NewColorsView()
  // colorsView.SetBorderFocusColor(theme.TopbarBorder)
  // colorsView.SetBorder(true)
  // colorsView.SetBorderColor(theme.Border)
  // colorsView.SetBackgroundColor(theme.SidebarBackground)

	content := tview.NewTextView()
	content.SetBorderPadding(1, 1, 2, 2)
  content.SetBorder(true)
  content.SetBorderFocusColor(theme.TopbarBorder)
  content.SetBorderColor(theme.Border)
  info.SetDynamicColors(true)
	content.SetBackgroundColor(theme.SidebarBackground)

	infoGrid := tview.NewGrid()
	infoGrid.AddItem(topbar, 0, 0, 1, 1, 5, 0, false)
	infoGrid.AddItem(info, 1, 0, 1, 1, 5, 0, false)

	contentGrid := tview.NewGrid()

	contentGrid.AddItem(colorsView, 0, 0, 1, 1, -5, 1, false)
	contentGrid.AddItem(content, 1, 0, 1, 1, -15, 1, false)

	view.AddItem(infoGrid, 0, 0, 1, 1, -8, 1, false)
	view.AddItem(contentGrid, 1, 0, 1, 1, -32, 1, false)

	// .NextFocusableComponent(tview.Down, )

	ft := &PaletteFileView{
		Box:         box,
		theme:       theme,
		view:        view,
		infoGrid:    infoGrid,
		contentGrid: contentGrid,
		topbarView:  topbar,
		infoView:    info,
		colorsView:  colorsView,
		contentView: content,
	}

	return ft
}

// Primitive interface

func (ft *PaletteFileView) UpdateSizes() {
	//  x,y,w,h := ft.view.GetRect()
	//  x, y, w, h := ft.view.GetInnerRect()
	// dump.P(x, y, w, h)
	//  ft.colorsView.SetRows(3, 5, h - 8)
}
func (ft *PaletteFileView) Draw(screen tcell.Screen) {
	// x, y, w, h := ft.topbarView.GetInnerRect()
	ft.view.Draw(screen)
}

func (ft *PaletteFileView) GetRect() (int, int, int, int) {
	// ft.UpdateSizes()
	return ft.view.GetRect()
}

func (ft *PaletteFileView) SetRect(x, y, width, height int) {
	ft.view.SetRect(x, y, width, height)
}

func (ft *PaletteFileView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ft.WrapInputHandler(
		func(ek *tcell.EventKey, f func(p tview.Primitive)) {
			if ek.Modifiers() == tcell.ModShift {
				ek = DirectionalFocusHandling(ek, AppModel.app)
				return
			}
			ft.view.InputHandler()(ek, f)
		},
	)
}

func (ft *PaletteFileView) Focus(delegate func(p tview.Primitive)) {
	ft.view.Focus(delegate)
}

func (ft *PaletteFileView) HasFocus() bool {
	return ft.view.HasFocus()
}

func (ft *PaletteFileView) Blur() {
	ft.view.Blur()
}

func (ft *PaletteFileView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return ft.view.MouseHandler()
}

func (v *PaletteFileView) SetPreview(fsnode *tree.FSNode) {
	defer func() {
		if err := recover(); err != nil {
			err, ok := err.(error)
			if ok {
				log.Printf(
					"Error loading config file: %s %v",
					fsnode.Path,
					errorx.WithStack(errorx.New(err.Error())),
				)
			}
		}
	}()
	i := []string{}
	// v.contentView.SetDynamicColors(true)
	d := LoadConfigFrom(fsnode.Path)
	v.colorsView.SetPalettes(d)
	var content []byte
  var cs string
	var cols *CoolorColorsPalette
  v.contentView.SetDynamicColors(true)

	if !fsnode.IsDir {
		// v.contentView.SetText(fsnode.Path)
		// v.contentView.SetTitle(fsnode.Path)
		content, _ = ioutil.ReadFile(fsnode.Path)
		cs, cols = StringColorizer(string(content))
		v.topbarView.SetText(strings.Join(cols.MakeSquarePalette(false), " "))
    v.contentView.Clear()
		v.contentView.SetText(cs)
		v.contentGrid.SetOffset(0, 0)
    v.infoGrid.SetOffset(0, 0)
	} else {
		v.contentView.SetText("")
	}

	// infoHex := fmt.Sprintf("#%6x", theme.InfoLabel.Hex())
	// i = append(i, )
	i = append(
		i,
		color.Render(
			fmt.Sprintf(" <infolabel>│      Path:</> %v", formatPath(fsnode.Path)),
		),
	)
	// i = append(i, color.Render(fmt.Sprintf(" <infolabel>│      Mode:</> %v", fsnode.Mode)))
	i = append(
		i,
		color.Render(fmt.Sprintf(" <infolabel>│  Modified:</> %v", fsnode.ModTime)),
	)
	if fsnode.Size != -1 {
		i = append(
			i,
			color.Render(
				fmt.Sprintf(" <infolabel>│      Size:</> %v", formatSize(fsnode.Size)),
			),
		)
	}

	if !fsnode.IsDir {
		i = append(
			i,
			color.Render(
				fmt.Sprintf(" <infolabel>│ Mime Type:</> %v", fsnode.MimeType),
			),
		)
	}
	if cols != nil {

		i = append(
			i,
			color.Render(
				fmt.Sprintf(" <infolabel>│ Colors:</> %v", len(cols.Colors)),
			),
		)
	}

	v.infoView.SetDynamicColors(true)
	v.infoView.SetText(tview.TranslateANSI(strings.Join(i, "\n")))
	// v.view.SetRows(3, len(i), len(d.GetPalettes())*5, 0)
  // v.view.SetRows(3 + len(i))
// x,y,width,height := v.colorsView.Grid.GetInnerRect()
//   _,_,_,_ = x,y,width,height
//   dump.P(x,y,width,height,v.contentView.GetOriginalLineCount())
  v.infoGrid.SetRows(5, len(i))
  v.contentGrid.SetRows(len(d.GetPalettes())*-5, -10)
  v.view.SetRows(5 + len(i), -10)
	v.view.SetOffset(0, 0)
}

func formatPath(p string) string {
	dir := filepath.Dir(p)
	base := filepath.Base(p)

	home := os.Getenv("HOME")

	if strings.HasPrefix(dir, home) {
		dir = strings.Replace(dir, home, "~", 1)
	}

	if dir == "/" {
		return fmt.Sprintf("[blue]/[normal]%s", base)
	}

	return fmt.Sprintf("[blue]%s/[normal]%s", dir, base)
}
func (cv *ColorsView) Draw(screen tcell.Screen) {
	cv.Box.DrawForSubclass(screen, cv)
	// cv.Grid.DrawForSubclass(screen, cv)
	x, y, w, h := cv.GetRect()
	_, _, _, _ = x, w, y, h
	x, y, w, h = cv.GetInnerRect()
	rows := make([]int, 0)
	for i, v := range cv.palettes {
		p := v.GetPalette()
		// p.Plainify(true)
		pt := NewPaletteTable(p)
    pt.SetBackgroundColor(theme.GetTheme().SidebarBackground)
		// pt.SetBackgroundColor(theme.GetTheme().ContentBackground)
		// pt.SetBorders(true).SetBorder(false)
		// pt.SetFixed(2, p.Len())
		f := tview.NewFrame(pt)
		f.SetTitleAlign(AlignCenter)
    f.SetBackgroundColor(theme.GetTheme().SidebarBackground)

		// f.SetTitle("")
		f.SetBorderPadding(0, 0, 0, 0)
		f.SetBorders(0, 0, 0, 0, 0, 0)
		cv.AddItem(f, i, 0, 1, 1, -5, 1, true)
		rows = append(rows, pt.rows*-5)
	}
	cv.SetRows(rows...)
	cv.Grid.Draw(screen)
}
