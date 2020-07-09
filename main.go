package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := cview.NewApplication()

	colors := make([]*PaletteColor, 5)
	currSelected := -1

	container := cview.NewFlex()

	for i := range colors {
		newBox := cview.NewBox()
		container.AddItem(newBox, 0, 1, false)
		colors[i] = NewPaletteColor(newBox, randomColor())
	}

	currSelected = 0
	colors[0].SetSelected(true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ch := event.Rune()
		kp := event.Key()
		switch {
		case ch == ' ' || ch == 'R':
			for _, pcol := range colors {
				if !pcol.locked {
					pcol.SetColor(randomColor())
				}
			}
			return nil
		case ch == 'h' || kp == tcell.KeyLeft:
			if currSelected != 0 {
				colors[currSelected].SetSelected(false)
				currSelected--
				colors[currSelected].SetSelected(true)
			}
			return nil
		case ch == 'l' || kp == tcell.KeyRight:
			if currSelected != len(colors)-1 {
				colors[currSelected].SetSelected(false)
				currSelected++
				colors[currSelected].SetSelected(true)
			}
			return nil
		case kp == tcell.KeyCtrlH:
			return nil
		case ch == 'w':
			currColor := colors[currSelected]
			currColor.SetLocked(!currColor.locked)
			return nil
		case ch == 'q' || kp == tcell.KeyEscape:
			app.Stop()
			return nil
		case ch == 'r':
			colors[currSelected].SetColor(randomColor())
			return nil
		}

		return event
	})

	app.SetRoot(container, true)

	if err := app.Run(); err != nil {
		panic(err)
	}

	for _, pcol := range colors {
		fmt.Printf("#%06x\t%t\t%t\n", pcol.Hex(), pcol.locked, pcol.box.HasFocus())
	}
}
