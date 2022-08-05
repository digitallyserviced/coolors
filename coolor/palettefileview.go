package coolor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/digitallyserviced/coolors/tree"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
)

var _ tview.Primitive = &PaletteFileView{}

type PaletteFileView struct {
	theme       *Theme
	view        *tview.Grid
	topbarView  *tview.TextView
	infoView    *tview.TextView
	colorsView  *ColorsView
	contentView *tview.TextView
}

type ColorsView struct {
	*tview.Grid
	palettes []Palette
}

func (cv *ColorsView) Draw(screen tcell.Screen) {
	// cv.Grid.Clear().SetSize(len(cv.palettes), 1, -2, -2).SetOffset(0, 0)
	x, y, w, h := cv.Grid.GetRect()
	dump.P(x, y, w, h)
	x, y, w, h = cv.Grid.GetInnerRect()
	dump.P(x, y, w, h)
	cv.Grid.SetGap(0, 0)
	rows := make([]int, 0)
	for i, v := range cv.palettes {
		p := v.GetPalette()
		p.Plainify(true)
		f := tview.NewFrame(p)
		f.SetTitle("")
		f.SetBorder(true).SetBorderPadding(0, 0, 2, 2).SetBorderColor(tree.GetTheme().TopbarBorder)
		f.SetBorders(0, 0, 0, 0, 1, 1)
		cv.Grid.AddItem(f, i, 0, 1, 1, 0, 0, false)
		rows = append(rows, 4)
	}
	cv.Grid.SetRows(rows...)
	cv.Box.DrawForSubclass(screen, cv)
	cv.Grid.Draw(screen)
}

func (cv *ColorsView) Clear() {
	cv.palettes = nil
	cv.Grid.Clear()
}

func (cv *ColorsView) SetPalettes(pdc *PaletteDataConfig) {
	cv.Clear()
	pals := pdc.GetPalettes()
	for _, v := range pals {
		pal := pdc.LoadPalette(v)
		cv.palettes = append(cv.palettes, pal.GetPalette())
	}
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

func NewPaletteFileView(theme *Theme) *PaletteFileView {
	view := tview.NewGrid()
	view.SetRows(3, 5, 0)

	topbar := tview.NewTextView()
	topbar.SetBorderPadding(0, 0, 0, 0)
	topbar.ScrollTo(0, 0)
	topbar.SetBorder(true)
	topbar.SetBorderColor(theme.TopbarBorder)
	topbar.SetDynamicColors(true)
	topbar.SetRegions(true)

	info := tview.NewTextView()
	info.SetDynamicColors(true)
	info.SetRegions(true)

	colorsView := NewColorsView()

	content := tview.NewTextView()
	content.SetBorderPadding(0, 0, 2, 2)
	content.SetBackgroundColor(theme.ContentBackground)

	view.AddItem(topbar, 0, 0, 1, 1, 3, 0, false)
	view.AddItem(info, 1, 0, 1, 1, 5, 0, false)
	view.AddItem(colorsView, 2, 0, 1, 1, 0, 0, false)
	// view.AddItem(content, 3, 0, 1, 1, 0, 0, false)

	ft := &PaletteFileView{
		theme:       theme,
		view:        view,
		topbarView:  topbar,
		infoView:    info,
		colorsView:  colorsView,
		contentView: content,
	}

	return ft
}

// Primitive interface

func (ft *PaletteFileView) Draw(screen tcell.Screen) {
	// x, y, w, h := ft.topbarView.GetInnerRect()
	// dump.P(x, y, w, h)
	// x, y, w, h = ft.topbarView.GetRect()
	// dump.P(x, y, w, h)
	ft.view.Draw(screen)
}
func (ft *PaletteFileView) GetRect() (int, int, int, int) {
	return ft.view.GetRect()
}
func (ft *PaletteFileView) SetRect(x, y, width, height int) {
	ft.view.SetRect(x, y, width, height)
}
func (ft *PaletteFileView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ft.view.InputHandler()
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
				log.Printf("Error loading config file: %s %v", fsnode.Path, errorx.WithStack(errorx.New(err.Error())))
			}
		}
	}()
	d := LoadConfigFrom(fsnode.Path)
	v.colorsView.SetPalettes(d)
	v.topbarView.SetText(formatPath(fsnode.Path))

	if !fsnode.IsDir {
		v.contentView.SetText(fsnode.Path)
		// v.contentView.SetTitle(fsnode.Path)
		content, _ := ioutil.ReadFile(fsnode.Path)
		v.contentView.SetText(string(content))
		v.contentView.ScrollTo(0, 0)
	} else {
		v.contentView.SetText("")
	}

	i := []string{}

	i = append(i, fmt.Sprintf(" [#5c6370]│      Mode:[normal] %v", fsnode.Mode))
	i = append(i, fmt.Sprintf(" [#5c6370]│  Modified:[normal] %v", fsnode.ModTime))
	if fsnode.Size != -1 {
		i = append(i, fmt.Sprintf(" [#5c6370]│      Size:[normal] %v", formatSize(fsnode.Size)))
	}

	if !fsnode.IsDir {
		i = append(i, fmt.Sprintf(" [#5c6370]│ Mime Type:[normal] %v", fsnode.MimeType))
	}

	v.infoView.SetText(strings.Join(i, "\n"))
	v.view.SetRows(3, len(i), len(d.GetPalettes())*4, 0)
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
