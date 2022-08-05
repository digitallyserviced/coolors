package status

import (
	"fmt"
	// "time"

	"github.com/digitallyserviced/tview"
	// "github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
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

var updateCh chan *StatusUpdate

func init() {
	updateCh = make(chan *StatusUpdate)
}

func NewStatusUpdate(elem, content string) {
	if updateCh != nil {
		updateCh <- &StatusUpdate{
			elem: elem, content: content,
		}

	} else {
		updateCh = make(chan *StatusUpdate)
		NewStatusUpdate("main", "READY")

	}
}

type Status struct {
	Severity Severity
	Message  string
}
type StatusItem struct {
	*tview.TableCell
	row, col int
}
type StatusBar struct {
	*tview.Table
	statuses map[string]*StatusItem
	app      *tview.Application
	elems    int8
}

func NewStatusBar(app *tview.Application) *StatusBar {
	sb := &StatusBar{
		statuses: make(map[string]*StatusItem),
		app:      app,
		Table:    tview.NewTable().SetBorders(false).InsertRow(0),
		elems:    0,
	}

	sb.Init()
	return sb
}

func statusUpdater(s *StatusBar) {
	for {
		select {
		case status := <-updateCh:
			dump.P(status)
			go func() {
				s.app.QueueUpdateDraw(func() {
					s.statuses[status.elem].TableCell.SetText(status.content)
					s.Table.SetCell(s.statuses[status.elem].row, s.statuses[status.elem].col, s.statuses[status.elem].TableCell)
				})
			}()
		}
	}
}

// 
func (s *StatusBar) Init() {
	s.Table.SetBorder(false)
	s.Table.SetBorderPadding(0, 0, 0, 0)
	c := tview.NewTableCell(fmt.Sprintf("[black:red:b] %s ", "STATUS")).SetAlign(tview.AlignCenter).SetExpansion(0)
	s.AddStatusItem("main", c)
	c = tview.NewTableCell(fmt.Sprintf("[black:yellow:b] 🗲  %s ", "action")).SetExpansion(0)
	s.AddStatusItem("action", c)
	c = tview.NewTableCell(fmt.Sprintf(" %s ", "")).SetExpansion(2)
	s.AddStatusItem("action_str", c)
	c = tview.NewTableCell(fmt.Sprintf(" %s ", "")).SetExpansion(2)
	s.AddStatusItem("fill", c)
	c = tview.NewTableCell(fmt.Sprintf(" %s ", "")).SetExpansion(0)
	s.AddStatusItem("color", c)
	c = tview.NewTableCell(fmt.Sprintf("[black:blue:b] (%s) ", "untitled")).SetExpansion(0)
	s.AddStatusItem("name", c)
	go statusUpdater(s)
}

func (s *StatusBar) AddStatusItem(label string, c *tview.TableCell) {
	si := &StatusItem{
		TableCell: c,
		row:       0,
		col:       int(s.elems),
	}
	s.statuses[label] = si
	s.Table.SetCell(0, int(s.elems), s.statuses[label].TableCell)
	s.elems = s.elems + 1
}

// func (s *StatusBar) add(label string, updates <-chan *Status) *StatusBar {
// 	n := s.Table.GetColumnCount()
//
// 	c := tview.NewTableCell(label)
// 	s.Table.SetCell(0, n, c)
//
// 	c = tview.NewTableCell("SHIT").SetExpansion(2)
// 	s.Table.SetCell(0, n+1, c)
//
// 	return s
// }

// func updateStatusCell(c *tview.TableCell, status *Status) {
// 	msg := status.Message
// 	switch status.Severity {
// 	case Refresh:
// 		c.SetText(msg).SetTextColor(tcell.ColorBlue)
// 	case Healthy:
// 		c.SetText(msg).SetTextColor(tcell.ColorGreen)
// 	case Warning:
// 		c.SetText(msg).SetTextColor(tcell.ColorYellow)
// 	case Alert:
// 		c.SetText(msg).SetTextColor(tcell.ColorRed)
// 	default:
// 		c.SetText(msg).SetTextColor(tcell.ColorDefault)
// 	}
// }
// func startClockStatus() chan *Status {
// 	updates := make(chan *Status)
// 	go func() {
// 		for {
// 			time.Sleep(10 * time.Second)
// 			update := &Status{
// 				Severity: Healthy,
// 				Message:  time.Now().String(),
// 			}
// 			updates <- update
// 		}
// 	}()
// 	return updates
// }
//
// vim: ts=2 sw=2 et ft=go
