package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"
  "runtime/pprof"

	"github.com/digitallyserviced/coolors/coolor"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
)

func main() {
	// setBordersChars()
 	go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()
	flag.Parse()
	values := flag.Args()
	app := tview.NewApplication()
	colsize := 0
	if len(values) > 0 {
		colsize = len(values)
	} else {
		colsize = 5
	}
	var colors *coolor.CoolorPalette

	if len(values) > 0 {
		colors = coolor.NewCoolorPaletteFromCssStrings(values)
	} else {
		colors = coolor.NewCoolorPaletteWithColors(coolor.GenerateRandomColors(colsize))
	}
  colors.SetSelected(0)


  scr, _ := tcell.NewScreen()

  scr.Size()
	statusbar := NewStatusBar(app)
	helpbar := &HelpBar{}
	helpbar.Init(app)
	// statusbar.AddStatusItem("Clock:", startClockStatus())
  pages := tview.NewPages()
	flexview := tview.NewFlex().SetDirection(tview.FlexRow)
  editor := coolor.NewCoolorEditor(app, colors)

	flexview.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ch := event.Rune()
		// kp := event.Key()
		switch {
		case ch == 'p':
      pages.SwitchToPage("palette").HidePage("editor")
      return nil
		case ch == 'e':
      pages.SwitchToPage("editor").HidePage("palette")
      return nil
    }
    return nil
    // return event
  })

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ch := event.Rune()
		kp := event.Key()
		switch {
		case ch == ' ' || ch == 'R':
			colors.Randomize()
			return nil

		case ch == 'h' || kp == tcell.KeyLeft:
			colors.NavSelection(-1)
			return nil

		case ch == 'l' || kp == tcell.KeyRight:
			colors.NavSelection(1)
			return nil

		case ch == 'w':
			colors.ToggleLockSelected()
			return nil

		case ch == 'q' || kp == tcell.KeyEscape:
			app.Stop()
			return nil

		case ch == 'r':
			color, _ := colors.GetSelected()
			color.Random()
			return nil

		case ch == '+': // Add a color
			colors.AddRandomCoolorColor()
			return nil

		case ch == '-': // Remove a color
			remcolor, idx := colors.GetSelected()
			colors.SetSelected(idx - 1)
			remcolor.Remove()
			return nil
		}
		
    // return nil
    // return nil
		return event
	})

  pages.AddPage("editor", editor, true, true)
  pages.AddAndSwitchToPage("palette", colors, true)
  pages.SetChangedFunc(func() {
    name, _ := pages.GetFrontPage()
    helpbar.UpdateRegion(name)
  })
  pages.SetFocusFunc(nil)
      pages.SwitchToPage("palette")
  flexview.
		AddItem(helpbar, 1, 1, false).
		AddItem(pages, 0, 10, false).
		AddItem(statusbar, 1, 1, false)
// _ = statusbar
		if err := app.SetRoot(flexview, true).Run(); err != nil {
			panic(err)
		}
  f, err := os.Create("mem.mprof")
        if err != nil {
            log.Fatal(err)
        }
        pprof.WriteHeapProfile(f)
        f.Close()

	for i := 0; i < colors.GetItemCount(); i++ {
		pcol := colors.GetItem(i)
		fmt.Printf("%s ", pcol)

	}
}


// vim: ts=2 sw=2 et ft=go
