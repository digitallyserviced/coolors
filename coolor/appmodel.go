package coolor

import (
	"fmt"
	"os"

	"github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	// "github.com/gdamore/tcell/v2/terminfo"
)

type Model struct {
	app      *tview.Application
	helpbar  *HelpBar
	rootView *tview.Flex
	status   *status.StatusBar
	scr      tcell.Screen
	pages    *MainContainer
}
// func init() {
// 	if err := test(); err != nil {
// 		log.Fatal(err)
// 	}
// }
  // scr.SetContent(x int, y int, primary rune, combining []rune, style tcell.Style)
  // ti, _ := terminfo.LookupTerminfo("xterm-256color")
  // ti.TPuts
  
  // if err != nil {
  //   // panic(err)
  // }
	// cmd := exec.Command("zsh", "-c", "echo", "-n", str)
	// output := &bytes.Buffer{}
	// errs := &bytes.Buffer{}
	// cmd.Stdout = output
 //  cmd.Stderr = errs
 //
	// if err := cmd.Run(); err != nil {
 //    panic(err)
	// }
 //  dump.P(errs.String(),output.String(), len(output.Bytes()))
var AppModel Model
func StartApp() {
  tty, ok := os.LookupEnv("GOTTY")
  var scr tcell.Screen
  if !ok {
    scr, _ = tcell.NewTerminfoScreen()
  } else {
    tty, _ := tcell.NewDevTtyFromDev(tty)
    scr, _ = tcell.NewTerminfoScreenFromTty(tty)
  }
	// if err != nil {
	// 	panic(err)
	// }
  err := scr.Init()
	if err != nil {
		panic(err)
	}
 if scr.HasMouse() {
  scr.EnableMouse()
  }
  scr.SetCursorStyle(tcell.CursorStyleBlinkingBar)
  scr.EnablePaste()
	AppModel.scr = scr
	AppModel.app = tview.NewApplication()
	
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
		switch ch {
		case 'Q':
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
  fmt.Println(colors)
	// for i := 0; i < colors.GetItemCount(); i++ {
	// 	pcol := colors.GetItem(i)
	// 	fmt.Printf("%s\n", pcol)
	// }
}
// vim: ts=2 sw=2 et ft=go
