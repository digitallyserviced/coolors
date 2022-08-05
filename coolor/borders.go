package coolor

import (
   "github.com/digitallyserviced/tview"
)
// Borders defines various borders used when primitives are drawn.
// These may be changed to accommodate a different look and feel.
var OrigBorders = struct {
	Horizontal  rune
	Vertical    rune
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune

	LeftT   rune
	RightT  rune
	TopT    rune
	BottomT rune
	Cross   rune

	HorizontalFocus  rune
	VerticalFocus    rune
	TopLeftFocus     rune
	TopRightFocus    rune
	BottomLeftFocus  rune
	BottomRightFocus rune
}{
	Horizontal:  tview.BoxDrawingsLightHorizontal,
	Vertical:    tview.BoxDrawingsLightVertical,
	TopLeft:     tview.BoxDrawingsLightDownAndRight,
	TopRight:    tview.BoxDrawingsLightDownAndLeft,
	BottomLeft:  tview.BoxDrawingsLightUpAndRight,
	BottomRight: tview.BoxDrawingsLightUpAndLeft,

	LeftT:   tview.BoxDrawingsLightVerticalAndRight,
	RightT:  tview.BoxDrawingsLightVerticalAndLeft,
	TopT:    tview.BoxDrawingsLightDownAndHorizontal,
	BottomT: tview.BoxDrawingsLightUpAndHorizontal,
	Cross:   tview.BoxDrawingsLightVerticalAndHorizontal,

	HorizontalFocus:  tview.BoxDrawingsDoubleHorizontal,
	VerticalFocus:    tview.BoxDrawingsDoubleVertical,
	TopLeftFocus:     tview.BoxDrawingsDoubleDownAndRight,
	TopRightFocus:    tview.BoxDrawingsDoubleDownAndLeft,
	BottomLeftFocus:  tview.BoxDrawingsDoubleUpAndRight,
	BottomRightFocus: tview.BoxDrawingsDoubleUpAndLeft,
}

var MyBorderStyle = struct {
	Horizontal  rune
	Vertical    rune
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune

	LeftT   rune
	RightT  rune
	TopT    rune
	BottomT rune
	Cross   rune

	HorizontalFocus  rune
	VerticalFocus    rune
	TopLeftFocus     rune
	TopRightFocus    rune
	BottomLeftFocus  rune
	BottomRightFocus rune
}{
	Horizontal:  tview.BoxDrawingsHeavyHorizontal,
	Vertical:    tview.BoxDrawingsHeavyVertical,
	TopLeft:     tview.BoxDrawingsHeavyDownAndRight,
	TopRight:    tview.BoxDrawingsHeavyDownAndLeft,
	BottomLeft:  tview.BoxDrawingsHeavyUpAndRight,
	BottomRight: tview.BoxDrawingsHeavyUpAndLeft,

	LeftT:   tview.BoxDrawingsHeavyVerticalAndRight,
	RightT:  tview.BoxDrawingsHeavyVerticalAndLeft,
	TopT:    tview.BoxDrawingsHeavyDownAndHorizontal,
	BottomT: tview.BoxDrawingsHeavyUpAndHorizontal,
	Cross:   tview.BoxDrawingsHeavyVerticalAndHorizontal,

	HorizontalFocus:  tview.BoxDrawingsDoubleHorizontal,
	VerticalFocus:    tview.BoxDrawingsDoubleVertical,
	TopLeftFocus:     tview.BoxDrawingsDoubleDownAndRight,
	TopRightFocus:    tview.BoxDrawingsDoubleDownAndLeft,
	BottomLeftFocus:  tview.BoxDrawingsDoubleUpAndRight,
	BottomRightFocus: tview.BoxDrawingsDoubleUpAndLeft,
}
