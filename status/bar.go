package status

import (
	"fmt"
	"sync"
	// "time"

	// "time"

	"github.com/digitallyserviced/tview"
	// "github.com/gdamore/tcell/v2"
	// "github.com/gookit/goutil/dump"
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
	}
}

type Status struct {
	Severity Severity
	Message  string
}
type StatusItem struct {
	name, format string
	*tview.TableCell
	row, col int
	l        sync.Mutex
}
type StatusBar struct {
	*tview.Table
	statuses     map[string]*StatusItem
	app          *tview.Application
	elems        int8
	lastStatuses map[string]StatusUpdate
}

func NewStatusBar(app *tview.Application) *StatusBar {
	sb := &StatusBar{
		Table:        tview.NewTable(),
		statuses:     make(map[string]*StatusItem),
		app:          app,
		elems:        0,
		lastStatuses: make(map[string]StatusUpdate),
	}

	sb.Table.SetBorders(false).InsertRow(0)
	sb.Init()
	return sb
}

func statusUpdater(s *StatusBar) {
	// ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case status := <-updateCh:
				s.app.QueueUpdateDraw(func() {
					s.UpdateCell(*status)
				})
		}
	}
}

func (s *StatusBar) UpdateCell(status StatusUpdate) {
	s.statuses[status.elem].TableCell.SetText(status.content)
	s.statuses[status.elem].TableCell.SetExpansion(s.statuses[status.elem].TableCell.Expansion)
	s.Table.SetCell(s.statuses[status.elem].row, s.statuses[status.elem].col, s.statuses[status.elem].TableCell)
}

func (si *StatusItem) UpdateItem(txt string) {
	si.SetText(fmt.Sprintf(si.format, txt))
}

func NewStatusItem(name, format, content string) *StatusItem {
	tc := tview.NewTableCell(fmt.Sprintf(format, content))
	tc.SetExpansion(1)
	si := &StatusItem{
		name:      name,
		format:    format,
		TableCell: tc,
		row:       0,
		col:       0,
		l:         sync.Mutex{},
	}
	return si
}

func (s *StatusBar) Init() {
	s.Table.SetBorder(false)
	s.Table.SetBorderPadding(0, 0, 0, 0)
	s.AddStatusItem(NewStatusItem("name", "[black:blue:b] (%s) [-:-:-]", "untitled"))
	s.AddStatusItem(NewStatusItem("action_str", " %s ", ""))
	sidot := NewStatusItem("dots", " %s ", "")
	sidot.SetExpansion(4).SetAlign(tview.AlignCenter)
	sidot.UpdateItem("")
	s.AddStatusItem(sidot)
	sifill := NewStatusItem("fill", " %s ", "")
	sifill.SetExpansion(3).SetAlign(tview.AlignCenter)
	s.AddStatusItem(sifill)
	s.AddStatusItem(NewStatusItem("color", " %s ", ""))
	s.AddStatusItem(NewStatusItem("action", "[black:yellow:b]   %s [-:-:-]", "action"))
	go statusUpdater(s)
}

func (s *StatusBar) AddStatusItem(si *StatusItem) {
	s.statuses[si.name] = si
	si.row = 0
	si.col = len(s.statuses) - 1
	s.Table.SetCell(0, si.col, s.statuses[si.name].TableCell)
	// s.elems = s.elems + 1
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
