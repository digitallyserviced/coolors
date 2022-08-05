package coolor

import (
	"fmt"

	"github.com/digitallyserviced/tview"
)

var ()

type HelpBar struct {
	*tview.TextView
	app *tview.Application
}

func (s *HelpBar) Init(app *tview.Application) {
	s.app = app
	s.TextView = tview.NewTextView()
	// s.UpdateRegion()
	s.SetRegions(true).SetBorder(false).SetBorderPadding(0, 0, 0, 0)
	// s.SetToggleHighlights(true)
	s.SetDynamicColors(true)
	s.SetText(fmt.Sprintf(`[red:black:b] ["palette"] %s [""] ["editor"] %s [""] ["mixer"] %s [""] `, "editor", "palette", "mixer")).SetTextAlign(tview.AlignCenter)
	s.UpdateRegion("editor")
}
func (s *HelpBar) UpdateRegion(r string) {
	s.Highlight(r)
}

// vim: ts=2 sw=2 et ft=go
