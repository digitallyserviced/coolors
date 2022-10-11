package lister

import (
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/digitallyserviced/coolors/theme"
)

type ListStyle interface {
	GetMainTextStyle() tcell.Style
	GetSecondaryTextStyle() tcell.Style
	GetShortcutStyle() tcell.Style
	GetSelectedStyle() tcell.Style
}

func NewListStyle() ListStyles {
	lis := ListStyles{
		main:  "",
		sec:   "",
		short: "",
		sel:   "",
	}

	return lis
}

func (f ListStyles) GetSelectedStyle() tcell.Style {
	if f.sel != "" {
		return *theme.GetTheme().Get(f.sel)
	}
	if theme.GetTheme().Get("list_sel") != nil {
		return *theme.GetTheme().Get("list_sel")
	}
	return tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor)
}

func (f ListStyles) GetShortcutStyle() tcell.Style {
	if f.short != "" {
		return *theme.GetTheme().Get(f.short)
	}
	if theme.GetTheme().Get("list_short") != nil {
		return *theme.GetTheme().Get("list_short")
	}
	return tcell.StyleDefault.Foreground(tview.Styles.SecondaryTextColor)
}

func (f ListStyles) GetSecondaryTextStyle() tcell.Style {
	if f.sec != "" {
		return *theme.GetTheme().Get(f.sec)
	}
	if theme.GetTheme().Get("list_second") != nil {
		return *theme.GetTheme().Get("list_second")
	}
	return tcell.StyleDefault.Foreground(tcell.ColorBlue)
}

func (f ListStyles) GetMainTextStyle() tcell.Style {
	if f.main != "" {
		return *theme.GetTheme().Get(f.main)
	}
	if theme.GetTheme().Get("list_main") != nil {
		return *theme.GetTheme().Get("list_main")
	}
	return tcell.StyleDefault.Foreground(tcell.ColorGreen)
}
