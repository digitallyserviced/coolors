package coolor

import (
	"fmt"
	"os"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/digitallyserviced/coolors/status"
	// "github.com/gdamore/tcell/v2/terminfo"
)

type Model struct {
	app      *tview.Application
	helpbar  *HelpBar
	rootView *tview.Flex
	status   *status.StatusBar
	scr      tcell.Screen
  pages *tview.Pages
	main    *MainContainer
  anims *tview.Pages
  w,h int
}

type Animations struct {
  *tview.Pages
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
	setupLogger()
  tty, ok := os.LookupEnv("GOTTY")
  var scr tcell.Screen
  if !ok {
    scr, _ = tcell.NewTerminfoScreen()
  } else {
    ok := false
    for !ok {
      l, err := os.Readlink(tty)
      if checkErrX(err, l){
        ok= true
        err = nil
      }
    }
    tty, _ := tcell.NewDevTtyFromDev(tty)
    scr, _ = tcell.NewTerminfoScreenFromTty(tty)
  }
  // simscr := tcell.NewSimulationScreen("")
  // AppModel.simscr = simscr

	// if err != nil {
	// 	panic(err)
	// }
  zlog.Info("SHIT")
  setupExpVars()
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
  AppModel.w,AppModel.h = AppModel.scr.Size()
	AppModel.status = status.NewStatusBar(AppModel.app)
	AppModel.helpbar = &HelpBar{}
	AppModel.main = NewMainContainer(AppModel.app)
	AppModel.rootView = tview.NewFlex().SetDirection(tview.FlexRow)
	AppModel.helpbar = NewHelpBar(AppModel.app)
	AppModel.rootView.
		AddItem(AppModel.status, 1, 1, false).
		AddItem(AppModel.main, 0, 10, true).
		AddItem(AppModel.helpbar, 1, 1, false)
  AppModel.pages = tview.NewPages()
  AppModel.pages.AddAndSwitchToPage("main", AppModel.rootView, true)
  AppModel.anims = tview.NewPages()
  AppModel.anims.SetBackgroundColor(0)
  AppModel.anims.Box.SetDontClear(true)
  AppModel.pages.AddPage("animation", AppModel.anims, true, true)
  spaceB := tview.NewBox()
  spaceB.SetDontClear(true)
  spaceB.SetRect(0, 0, AppModel.w, AppModel.h)
  AppModel.anims.AddPage("spacer", spaceB, true, true)
  AppModel.anims.ShowPage("spacer")
  spaceB.SetRect(0, 0, AppModel.w, AppModel.h)
  AppModel.pages.ShowPage("animation")
  AppModel.anims.SetRect(0, 0, AppModel.w, AppModel.h)
  AppModel.anims.SetVisible(true)
// AppModel.anims.Box.
//   spaceB := tview.NewBox()
//   spaceB.SetDontClear(false)
//   spaceB.SetBorder(true)
//   AppModel.anims.AddPage("spacer", spaceB, true, true)
//   AppModel.pages.ShowPage("spacer")
//   spaceB.SetRect(0, 0, w, h)
//   AppModel.anims.HidePage("spacer")
	if err := AppModel.app.SetRoot(AppModel.pages, true).Run(); err != nil {
		panic(err)
	}

	AppModel.main.CloseConfig()
	colors := AppModel.main.GetPalette()
  fmt.Println(colors)
	// for i := 0; i < colors.GetItemCount(); i++ {
	// 	pcol := colors.GetItem(i)
	// 	fmt.Printf("%s\n", pcol)
	// }
}
// vim: ts=2 sw=2 et ft=go
