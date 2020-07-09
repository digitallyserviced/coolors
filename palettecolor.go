package main

import (
	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

type PaletteColor struct {
	box      *cview.Box
	col      tcell.Color
	locked   bool
	selected bool
}

func NewPaletteColor(box *cview.Box, col tcell.Color) *PaletteColor {
	box.
		SetBackgroundColor(col)
	p := &PaletteColor{box, col, false, false}
	box.
		SetMouseCapture(func(action cview.MouseAction, event *tcell.EventMouse) (cview.MouseAction, *tcell.EventMouse) {
			p.SetLocked(!p.locked)
			return action, event
		})
	return p
}

func (p *PaletteColor) Hex() int32 {
	return p.col.Hex()
}

func (p *PaletteColor) RGB() (int32, int32, int32) {
	return p.col.RGB()
}

func (p *PaletteColor) SetColor(col tcell.Color) {
	p.box.SetBackgroundColor(col)
	p.col = col
}

func (p *PaletteColor) SetLocked(b bool) {
	p.locked = b
	p.updateStyle()
}

func (p *PaletteColor) SetSelected(b bool) {
	p.selected = b
	p.updateStyle()
}

func (p *PaletteColor) updateStyle() {
	if p.locked || p.selected {
		inverse := getBorderColor(p.col)
		p.box.
			SetBorder(true).
			SetBorderColor(inverse).
			SetTitleColor(inverse)
		title := ""
		if p.locked {
			title += ""
		}
		if p.selected {
			title += ""
		}
		p.box.SetTitle(title)
	} else {
		p.box.SetBorder(false)
	}
}

func getBorderColor(col tcell.Color) tcell.Color {
	r, g, b := col.RGB()
	if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
		return tcell.NewHexColor(0x000000)
	}
	return tcell.NewHexColor(0xFFFFFF)
}

func inverseColor(col tcell.Color) tcell.Color {
	r, g, b := col.RGB()
	return tcell.NewRGBColor(255-r, 255-g, 255-b)
}
