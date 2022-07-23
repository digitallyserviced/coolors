package main

import (
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
)

type Severity int

const (
	Unknown Severity = iota
	Refresh
	Healthy
	Warning
	Alert
)

type StatusUpdate struct {
	elem    string
	content string
}

type Status struct {
	Severity Severity
	Message  string
}

type StatusBar struct {
	*tview.Table
	statuses map[string]*tview.TableCell
	app      *tview.Application
	elems    int8
}

func NewStatusBar(app *tview.Application) *StatusBar {
	sb := &StatusBar{
		statuses: make(map[string]*tview.TableCell),
		app:      app,
		Table:    tview.NewTable().SetBorders(false).InsertRow(0),
		elems:    0,
	}
	sb.Init()
	return sb
}

func (s *StatusBar) Init() <-chan *StatusUpdate {
	s.Table.SetBorder(false)
	s.Table.SetBorderPadding(0, 0, 0, 0)
	updatech := make(chan *StatusUpdate)
  return updatech
	// c := tview.NewTableCell(fmt.Sprintf("[black:red:b] %s ", "STATUS")).SetAlign(tview.AlignCenter).SetExpansion(0)
 //  s.AddStatusItem("main", c, updatech)
	// c = tview.NewTableCell(fmt.Sprintf("%s", "")).SetExpansion(10)
 //  s.AddStatusItem("fill", c, updatech)
	// c = tview.NewTableCell(fmt.Sprintf("[black:blue:b] (%s) ", "untitled")).SetExpansion(0)
 //  s.AddStatusItem("name", c, updatech)
 //  return updatech
}

func (s *StatusBar) AddStatusItem(label string, c *tview.TableCell, updatech <-chan *StatusUpdate)  {
	s.statuses[label] = c
	s.Table.SetCell(0, int(s.elems), c)
  s.elems = s.elems + 1
	go func() {
		for {
			select {
			case status := <-updatech:
				s.app.QueueUpdateDraw(func() {
          s.statuses[status.elem].SetText(status.content)
				})
			case <-time.After(time.Second * 1):
        continue
      default:
        continue
			}
		}
	}()
}

func (s *StatusBar) add(label string, updates <-chan *Status) *StatusBar {
	n := s.Table.GetColumnCount()

	c := tview.NewTableCell(label)
	s.Table.SetCell(0, n, c)

	c = tview.NewTableCell("SHIT").SetExpansion(2)
	s.Table.SetCell(0, n+1, c)

	return s
}

func updateStatusCell(c *tview.TableCell, status *Status) {
	msg := status.Message
	switch status.Severity {
	case Refresh:
		c.SetText(msg).SetTextColor(tcell.ColorBlue)
	case Healthy:
		c.SetText(msg).SetTextColor(tcell.ColorGreen)
	case Warning:
		c.SetText(msg).SetTextColor(tcell.ColorYellow)
	case Alert:
		c.SetText(msg).SetTextColor(tcell.ColorRed)
	default:
		c.SetText(msg).SetTextColor(tcell.ColorDefault)
	}
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
