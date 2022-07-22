package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/digitallyserviced/coolors/coolor"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
)

func main() {
	// setBordersChars()
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
      pages.SwitchToPage("palette")
      return nil
		case ch == 'e':
      pages.SwitchToPage("editor")
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

  pages.AddPage("editor", editor, true, false)
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

	for i := 0; i < colors.GetItemCount(); i++ {
		pcol := colors.GetItem(i)
		fmt.Printf("%s ", pcol)

	}
}

func NewCoolorPaletteFromCssStrings(values []string) {
	panic("unimplemented")
}

// func (s *StatusBar) AddStatus(id, label string) *StatusBar {
// 	c := tview.NewTableCell("").
// 		SetExpansion(2)
// 	s.statuses[id] = c
//
// 	t := s
// 	n := t.GetColumnCount()
// 	t.SetCell(0, n, tview.NewTableCell(label))
// 	t.SetCell(0, n+1, c)
//
// 	return s
// }
//
// func (s *StatusBar) UpdateStatusItem(id, msg string, ok bool) *StatusBar {
// 	// TODO use app.QueueUpdate() to make thread-safe
// 	c := s.statuses[id]
// 	if ok {
// 		c.SetText("✓ " + msg).
// 			SetTextColor(tcell.ColorGreen)
// 	} else {
// 		c.SetText("x " + msg).
// 			SetTextColor(tcell.ColorRed)
// 	}
// 	s.app.Draw()
// 	return s
// }
func startClockStatus() chan *Status {
	updates := make(chan *Status)
	go func() {
		for {
			time.Sleep(10 * time.Second)
			update := &Status{
				Severity: Healthy,
				Message:  time.Now().String(),
			}
			updates <- update
		}
	}()
	return updates
}

// vim: ts=2 sw=2 et ft=go
