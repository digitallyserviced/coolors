package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"gitlab.com/tslocum/cview"
)

func PaletteColorFlex(pcols []*PaletteColor) *cview.Flex {
	newFlex := cview.NewFlex()
	for _, pcol := range pcols {
		newFlex.AddItem(pcol.box, 0, 1, false)
	}
	return newFlex
}

func randomiseColors(pcols []*PaletteColor) {
	for _, pcol := range pcols {
		pcol.SetColor(randomColor())
	}
}

const (
	FLEX_MIN = 2
	FLEX_MAX = 10
)

func main() {
	rand.Seed(time.Now().UnixNano())

	app := cview.NewApplication()

	colors := make([]*PaletteColor, 5)

	for i := range colors {
		colors[i] = NewPaletteColor(cview.NewBox(), randomColor())
	}

	container := PaletteColorFlex(colors)

	currSelected := 0
	colors[0].SetSelected(true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ch := event.Rune()
		kp := event.Key()
		switch {
		case ch == ' ' || ch == 'R':
			randomiseColors(colors)
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

		case ch == '+' && len(colors) < FLEX_MAX: // Add a color
			colors = append(colors, NewPaletteColor(cview.NewBox(), randomColor()))
			container = PaletteColorFlex(colors)
			app.SetRoot(container, true)

		case ch == '-' && len(colors) > 2: // Remove a color
			colors = append(colors[:currSelected], colors[currSelected+1:]...)
			container = PaletteColorFlex(colors)
			app.SetRoot(container, true)
			if currSelected != 0 {
				currSelected--
			}
			colors[currSelected].SetSelected(true)
		}

		return event
	})

	app.SetRoot(container, true)

	if err := app.Run(); err != nil {
		panic(err)
	}

	isTerminal := terminal.IsTerminal(int(os.Stdout.Fd()))

	for _, pcol := range colors {
		r, g, b := pcol.RGB()
		br, bg, bb := getFGColor(pcol.col).RGB()
		if isTerminal {
			fmt.Printf(
				"\033[48;2;%d;%d;%d;38;2;%d;%d;%dm #%06x \033[0m\n",
				r, g, b, br, bg, bb, pcol.Hex(),
			)
		} else {
			fmt.Printf(" %06x \n", pcol.Hex())
		}
	}
}
