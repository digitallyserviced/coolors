package coolor

import (
	"fmt"
	"os"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/digitallyserviced/coolors/coolor/anim"
	xxp "github.com/digitallyserviced/coolors/coolor/xp"
	"github.com/digitallyserviced/coolors/status"
	// "github.com/gdamore/tcell/v2/terminfo"
)

// import {WezTerm} from "./wezterm"

type Model struct {
	app      *tview.Application
	helpbar  *HelpBar
	rootView *tview.Flex
	status   *status.StatusBar
	scr      tcell.Screen
	pages    *tview.Pages
	main     *MainContainer
	anims    *tview.Pages
	w, h     int
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
	// setupOut()
	setupLogger()
	tty, ok := os.LookupEnv("GOTTY")
	var scr tcell.Screen
	if !ok {
		scr, _ = tcell.NewTerminfoScreen()
	} else {
		ok := false
		for !ok {
			l, err := os.Readlink(tty)
			if checkErrX(err, l) {
				ok = true
				err = nil
			} else {
				stat, err := os.Stat(tty)
				if err != nil || stat.IsDir() {
					fmt.Println(err)
				} else {
					fmt.Println(stat)
					ok = true
					err = nil
				}
				// fmt.Println(err)
			}
		}
		tty, _ := tcell.NewDevTtyFromDev(tty)
		scr, _ = tcell.NewTerminfoScreenFromTty(tty)
	}
	zlog.Info("SHIT")
	xxp.SetupExpVars()
	err := scr.Init()
	if err != nil {
		panic(err)
	}
	if scr.HasMouse() {
		// scr.EnableMouse()
	}
	scr.SetCursorStyle(tcell.CursorStyleBlinkingBar)
	// scr.EnablePaste()
	AppModel.scr = scr
	AppModel.app = tview.NewApplication()
	AppModel.app.SetScreen(AppModel.scr)
	AppModel.w, AppModel.h = AppModel.scr.Size()
	AppModel.status = status.NewStatusBar(AppModel.app)
	AppModel.helpbar = &HelpBar{}
	AppModel.pages = tview.NewPages()
	AppModel.anims = tview.NewPages()
  anim.Init(AppModel.app, AppModel.scr, AppModel.pages, AppModel.anims)
	AppModel.main = NewMainContainer(AppModel.app)
	AppModel.rootView = tview.NewFlex().SetDirection(tview.FlexRow)
	AppModel.helpbar = NewHelpBar(AppModel.app)
	AppModel.rootView.
		AddItem(MakeBoxItem("", ""), 1, 1, false).
		AddItem(AppModel.main, 0, 10, true).
		AddItem(AppModel.status, 2, 1, false)
		// AddItem(AppModel.helpbar, 1, 1, false)
	AppModel.pages.AddAndSwitchToPage("main", AppModel.rootView, true)
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
	if err := AppModel.app.SetRoot(AppModel.pages, true).Run(); err != nil {
		panic(err)
	}

	AppModel.main.CloseConfig()
	colors := AppModel.main.GetPalette()
	fmt.Println(colors)
}

// sp_block_l_begin='▌'         ; sp_block_l_middl=''       ; sp_block_l_close='▐'
// sp_pentagon_begin=''        ; sp_pentagon_middl=''      ; sp_pentagon_close=''
// sp_tiny_begin=' '            ; sp_tiny_middl=' '          ; sp_tiny_close=' '
// sp_blank_begin='  '           ; sp_blank_middl='  '         ; sp_blank_close='  '            ;
// sp_block_l_begin=''         ; sp_block_l_middl=''       ; sp_block_l_close=''
// sp_block_c_begin='█'         ; sp_block_c_middl=''       ; sp_block_c_close='█'
// sp_block_d_begin='█'         ; sp_block_d_middl='██'     ; sp_block_d_close='█'
// sp_block_e_begin='▆ '         ; sp_block_e_middl='▆ '       ; sp_block_e_close='▆'
// sp_block_r_begin=''         ; sp_block_r_middl=''       ; sp_block_r_close=''
// sp_line_top_begin='┏╸━'       ; sp_line_top_middl='━╸'      ; sp_line_top_close='━━╸━┓'
// sp_line_begin='╺╸'           ; sp_line_middl='··'         ; sp_line_close='·╺╸'
// sp_line_top_begin='┏╸'       ; sp_line_top_middl='·'      ; sp_line_top_close='·╺┓'       ; #   
// sp_cross_begin=' '           ; sp_cross_middl=' '         ; sp_cross_close=' '
// sp_dot_begin=' '             ; sp_dot_middl=' '           ; sp_dot_close=' '
// sp_lash_begin='●'            ; sp_lash_middl='●'          ; sp_lash_close='●●'
// sp_dotline_begin='╸⏽'         ; sp_dotline_middl='●⏽'       ; sp_dotline_close='╺'
// sp_lash_begin='╸⏽'           ; sp_lash_middl='●⏽'          ; sp_lash_close='●╺'
// sp_box_slant_begin='█┣╸●'    ; sp_box_slant_middl=' '     ; sp_box_slant_close=' ●╺┫█'
// sp_circle_slant_begin='█┣ ●' ; sp_circle_slant_middl=' '  ; sp_circle_slant_close=' ● ┫█'
// sp_dot_slant_begin='█🮈╸'    ; sp_dot_slant_middl='·'     ; sp_dot_slant_close='·╺▍█'
// sp_line_top_mini_begin='┏╸'  ; sp_line_top_mini_middl='·' ; sp_line_top_mini_close='·╺╺┓'
// sp_line_bot_begin='┗━╺╸╺╸'    ; sp_line_bot_middl='╺━'      ; sp_line_bot_close='┛'
// sp_line_bo2_begin='┗╸'        ; sp_line_bo2_middl='╺╸'      ; sp_line_bo2_close='┛'
// sp_box_begin='■'              ; sp_box_middl='■'            ; sp_box_close='■'
// sp_box2_begin='█🮈'            ; sp_box2_middl='▍'           ; sp_box2_close='▍█'
//
// vim: ts=2 sw=2 et ft=go
