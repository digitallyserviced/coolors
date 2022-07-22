package main

import (
	"github.com/digitallyserviced/tview"
)

func setBordersChars() {
	tview.Borders.Horizontal = tview.BoxDrawingsLightDoubleDashHorizontal
	tview.Borders.Vertical = tview.BoxDrawingsLightQuadrupleDashVertical
	tview.Borders.TopLeft = tview.BoxDrawingsLightDownAndRight
	tview.Borders.TopRight = tview.BoxDrawingsLightDownAndLeft
	tview.Borders.BottomLeft = tview.BoxDrawingsLightUpAndRight
	tview.Borders.BottomRight = tview.BoxDrawingsLightUpAndLeft
	// tview.Borders.Horizontal = tview.BoxDrawingsLightHorizontal
	// tview.Borders.Vertical = tview.BoxDrawingsLightVertical
	// tview.Borders.TopLeft = tview.BoxDrawingsLightDownAndRight
	// tview.Borders.TopRight = tview.BoxDrawingsLightDownAndLeft
	// tview.Borders.BottomLeft = tview.BoxDrawingsLightUpAndRight
	// tview.Borders.BottomRight = tview.BoxDrawingsLightUpAndLeft

	tview.Borders.LeftT = tview.BoxDrawingsLightVerticalAndRight
	tview.Borders.RightT = tview.BoxDrawingsLightVerticalAndLeft
	tview.Borders.TopT = tview.BoxDrawingsLightDownAndHorizontal
	tview.Borders.BottomT = tview.BoxDrawingsLightUpAndHorizontal
	tview.Borders.Cross = tview.BoxDrawingsLightVerticalAndHorizontal

	tview.Borders.HorizontalFocus = tview.BoxDrawingsHeavyDoubleDashHorizontal
	tview.Borders.VerticalFocus = tview.BoxDrawingsHeavyQuadrupleDashVertical
	tview.Borders.TopLeftFocus = tview.BoxDrawingsHeavyDownAndRight
	tview.Borders.TopRightFocus = tview.BoxDrawingsHeavyDownAndLeft
	tview.Borders.BottomLeftFocus = tview.BoxDrawingsHeavyUpAndRight
	tview.Borders.BottomRightFocus = tview.BoxDrawingsHeavyUpAndLeft
}
