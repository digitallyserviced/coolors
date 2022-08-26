package coolor

import (
	"fmt"

	"github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
)

type Model struct {
	app      *tview.Application
	helpbar  *HelpBar
	rootView *tview.Flex
	status   *status.StatusBar
	scr      tcell.Screen
	pages    *MainContainer
}

var AppModel Model

func StartApp() {
	AppModel.app = tview.NewApplication()
	scr, err := tcell.NewTerminfoScreen()
	if err != nil {
		panic(err)
	}
	AppModel.scr = scr
	AppModel.scr.Init()
	AppModel.app.SetScreen(AppModel.scr)
	AppModel.status = status.NewStatusBar(AppModel.app)
	AppModel.helpbar = &HelpBar{}
	AppModel.pages = NewMainContainer(AppModel.app)
	AppModel.rootView = tview.NewFlex().SetDirection(tview.FlexRow)
	AppModel.helpbar = NewHelpBar(AppModel.app)
	AppModel.rootView.
		AddItem(AppModel.helpbar, 1, 1, false).
		AddItem(AppModel.pages, 0, 10, true).
		AddItem(AppModel.status, 1, 1, false)
	AppModel.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ch := event.Rune()
		// kp := event.Key()
		switch {
		case ch == 'Q':
			AppModel.app.Stop()
			return nil
		}
		return event
	})
	if err := AppModel.app.SetRoot(AppModel.rootView, true).Run(); err != nil {
		panic(err)
	}

	AppModel.pages.CloseConfig()
	colors := AppModel.pages.GetPalette()
	for i := 0; i < colors.GetItemCount(); i++ {
		pcol := colors.GetItem(i)
		fmt.Printf("%s ", pcol)
	}
}
// vim: ts=2 sw=2 et ft=go
