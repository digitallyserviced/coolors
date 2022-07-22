package coolor

import (
	"fmt"
	"sync"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
)
type CoolorColors []*CoolorColor

type CoolorPalette struct {
	*tview.Flex
	colors      CoolorColors
	selectedIdx int
	name        string
	l           *sync.RWMutex
  handlers      map[string]EventHandlers
}

func BlankCoolorPalette() *CoolorPalette {
	cp := &CoolorPalette{
		Flex:        tview.NewFlex(),
		colors:      CoolorColors{},
		selectedIdx: 0,
		name:        "",
		l:           &sync.RWMutex{},
    handlers: make(map[string]EventHandlers),
	}
	cp.Flex.SetDirection(tview.FlexColumn)
	return cp
}

func DefaultCoolorPalette() *CoolorPalette {
	tcols := GenerateRandomColors(5)
	cp := NewCoolorPaletteWithColors(tcols)
	return cp
}

func NewCoolorPaletteFromCssStrings(cols []string) *CoolorPalette {
	cp := BlankCoolorPalette()
	for _, v := range cols {
		cp.AddCssCoolorColor(v)
	}
  cp.SetSelected(0)
	return cp
}

func NewCoolorPaletteWithColors(tcols []tcell.Color) *CoolorPalette {
	cp := BlankCoolorPalette()
	for _, v := range tcols {
  fmt.Printf("%06x", v.Hex())
		cp.AddCssCoolorColor(fmt.Sprintf("#%06x", v.Hex()))
	}
    cp.SetSelected(0)
	return cp
}

func (cp *CoolorPalette) AddCoolorColor(color *CoolorColor) *CoolorColor {
  fmt.Printf("%s", color)
	color.pallette = cp
	cp.l.Lock()
	defer cp.l.Unlock()
	cp.colors = append(cp.colors, color)
	cp.Flex.AddItem(cp.colors[len(cp.colors)-1], 0, 1, true)
	return cp.colors[len(cp.colors)-1]
}

func (cp *CoolorPalette) AddCssCoolorColor(c string) *CoolorColor {
	color := cp.AddCoolorColor(NewCoolorColor(c))
	return color
}

func (cp *CoolorPalette) AddRandomCoolorColor() *CoolorColor {
	newc := NewRandomCoolorColor()
	newc.pallette = cp
	cp.colors = append(cp.colors, newc)
	cp.Flex.AddItem(newc, 0, 1, false)
	return newc
}

func (cp *CoolorPalette) GetSelected() (*CoolorColor, int) {
	cp.l.RLock()
	defer cp.l.RUnlock()
  
  if cp.colors[cp.selectedIdx] != nil {
    return cp.colors[cp.selectedIdx], cp.selectedIdx
  }
  return nil, -1
}

func (cp *CoolorPalette) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return cp.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		ch := event.Rune()
		kp := event.Key()
		switch {
		case ch == ' ' || ch == 'R':
			cp.Randomize()

		case ch == 'h' || kp == tcell.KeyLeft:
			cp.NavSelection(-1)

		case ch == 'l' || kp == tcell.KeyRight:
			cp.NavSelection(1)

		case ch == 'w':
			cp.ToggleLockSelected()

		case ch == 'r':
			color, _ := cp.GetSelected()
			color.Random()

		case ch == '+': // Add a color
			cp.AddRandomCoolorColor()

		case ch == '-': // Remove a color
			remcolor, idx := cp.GetSelected()
			cp.SetSelected(idx - 1)
			remcolor.Remove()
		}
	})
}

func (cp *CoolorPalette) NavSelection(idx int) error {
	cp.l.RLock()
	newidx := cp.selectedIdx + idx
	if newidx >= len(cp.colors) {
		newidx = len(cp.colors) - 1
	}
	if newidx < 0 {
		newidx = 0
	}
	cp.l.RUnlock()
	cp.SetSelected(newidx)
	return nil
}

func (cp *CoolorPalette) Randomize() int {
	changed := 0
	for _, v := range cp.colors {
		if v.Random() {
			changed += 1
		}
	}
	cp.ResetViews()
	return changed
}

func (cpfv *CoolorPalette) ResetViews() {
}

func (cc *CoolorPalette) SpawnSelectionEvent(c *CoolorColor, idx int) bool {
  if len(cc.handlers["selected"]) > 0 {
    ev := &SelectionEvent{
      color: c,
      idx: int8(idx),
    }
    for _, v := range cc.handlers["selected"] {
      if v != nil {
        eh, ok := v.(tcell.EventHandler)
        if !ok {
          panic(ok)
        }
        eh.HandleEvent(ev)
      }
    }
  }
  return true
}

func (cc *CoolorPalette) AddEventHandler(t string, h tcell.EventHandler) {
	cc.l.Lock()
	defer cc.l.Unlock()

  if cc.handlers[t] == nil {
    cc.handlers[t] = make(EventHandlers, 0)
  }
  cc.handlers[t] = append(cc.handlers[t], h)
}

func (cp *CoolorPalette) SetSelected(idx int) error {
	sel, _ := cp.GetSelected()
	sel.SetSelected(false)
	if idx < len(cp.colors) {
		if idx < 0 {
			idx = 0
		}
    cp.l.Lock()
		cp.selectedIdx = idx
		cp.colors[cp.selectedIdx].SetSelected(true)
		cp.l.Unlock()
    cp.SpawnSelectionEvent(cp.colors[cp.selectedIdx],cp.selectedIdx)
		return nil
	}
	return fmt.Errorf("No valid color at idx: %d", idx)
}

func (cp *CoolorPalette) ToggleLockSelected() (*CoolorColor, int) {
	cc, _ := cp.GetSelected()
	cc.ToggleLocked()
	return cp.colors[cp.selectedIdx], cp.selectedIdx
}

// vim: ts=2 sw=2 et ft=go
