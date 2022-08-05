package coolor

import (
  "github.com/digitallyserviced/tview"
  "github.com/gdamore/tcell/v2"
)

type MainMenuStyle struct {
  SelectedBg tcell.Color
  SelectedFg tcell.Color
  DefaultBg tcell.Color
  DefaultFg tcell.Color
  HoverBg tcell.Color 
  HoverFg tcell.Color 
}

var DefaultMenuStyle = MainMenuStyle{
  SelectedBg: tview.Styles.MoreContrastBackgroundColor,
  SelectedFg: tview.Styles.SecondaryTextColor,
  DefaultBg: tview.Styles.PrimitiveBackgroundColor,
  DefaultFg: tview.Styles.TertiaryTextColor,
  HoverBg: tview.Styles.ContrastBackgroundColor,
  HoverFg: tview.Styles.SecondaryTextColor,
}
// vim: ts=2 sw=2 et ft=go
