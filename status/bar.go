package status

import (
	"fmt"
	"sync"
	"time"

	"github.com/digitallyserviced/coolors/theme"
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

type Status struct {
	Severity Severity
	Message  string
}

type StatusItem struct {
	name, format string
	*tview.TableCell
	w, row, col int
	updates     chan *StatusUpdate
	l           sync.Mutex
	bar         *StatusBar
	positioner  func(s string, n int) string
}

type StatusBars struct {
	*tview.Flex
	Main         *StatusBar
	Transient    *StatusBar
	Title        *StatusBar
	app          *tview.Application
	elems        int8
	lastStatuses map[string]StatusUpdate
}

type StatusUpdate struct {
	elem    string
	content string
	timeOut time.Duration
}

var updateCh chan *StatusUpdate

func init() {
	updateCh = make(chan *StatusUpdate)
}

type StatusBar struct {
	*tview.Table
	bars     *StatusBars
	statuses map[string]*StatusItem
}

func NewStatusUpdateWithTimeout(elem, content string, to time.Duration) {
	if to == 0*time.Second {
		to = 3000 * time.Millisecond
	}
	if updateCh != nil {
		updateCh <- &StatusUpdate{
			elem: elem, content: content, timeOut: to,
		}
	} else {
		updateCh = make(chan *StatusUpdate)
	}
}
func NewStatusUpdate(elem, content string) {
	// if elem == "action_str" {
	//   NewStatusUpdateWithTimeout(elem, content, time.Second * )
	// }
	if updateCh != nil {
		updateCh <- &StatusUpdate{
			elem: elem, content: content, timeOut: time.Second * 0,
		}
	} else {
		updateCh = make(chan *StatusUpdate)
	}
}

func NewStatusBar(app *tview.Application) *StatusBars {
	sb := &StatusBars{
		Flex: tview.NewFlex(),
		// Main:         tview.NewTable(),
		// Transient:    tview.NewTable(),
		// Title:        tview.NewTable(),
		app:          app,
		elems:        0,
		lastStatuses: make(map[string]StatusUpdate),
	}
	sb.Title = sb.SetupStatusBar()
	sb.Title.SetBackgroundColor(theme.GetTheme().ContentBackground)
	sb.Main = sb.SetupStatusBar()
	sb.Main.SetBackgroundColor(theme.GetTheme().ContentBackground)
	sb.Transient = sb.SetupStatusBar()
	sb.Init()
	sb.Flex.SetDirection(tview.FlexRow)
	sb.Flex.AddItem(sb.Main, 1, 0, false)
	sb.Flex.AddItem(sb.Transient, 1, 0, false)
	return sb
}
func (s *StatusBars) Init() {
	// tv := tview.NewTextView()
	done := make(chan struct{})
	// pn := NewStatusItem("name", "[red:gray:-][-:-:-][-:gray:-]  [-:-:-][gray:red:-] [-:-:-][black:red:b] %s [-:-:-][red:gray:-][-:gray:-] [gray:-:-][-:-:-]", "untitled", s, done) // [-:gray:-]  [-:-:-]
	pn := s.Main.NewStatusItem(
		"name",
		"%s",
		"untitled",
		s,
		done,
	) // [-:gray:-]  [-:-:-]
	// sy := theme.GetTheme().Get("palette_name")
	// fmt.Printf("%v", sy)
	// pn.SetStyle(*sy)
	pn.SetAlign(tview.AlignCenter)
	s.Main.AddStatusItem(pn)
	sidot := s.Main.NewStatusItem("dots", "%s", "", s, done)
	sidot.SetExpansion(2)
	sidot.SetAlign(tview.AlignCenter)
	sidot.UpdateItem("")
	s.Main.AddStatusItem(sidot)
	// s.AddStatusItem(NewStatusItem("fill1", "%s", "", s, done))
	// tag := NewStatusItem("tag", "%s", "  ", s, done)
	// tag.SetAlign(tview.AlignCenter)
	// tag.SetExpansion(1)
	// s.AddStatusItem(tag)
	// action := NewStatusItem("action", "   %s ", "action", s, done)
	// action.SetStyle(*theme.GetTheme().Get("action"))
	// s.AddStatusItem(action)
	// pn := NewStatusItem("name", "[red:gray:-][-:-:-][-:gray:-]  [-:-:-][gray:red:-] [-:-:-][black:red:b] %s [-:-:-][red:gray:-][-:gray:-] [gray:-:-][-:-:-]", "untitled", s, done) // [-:gray:-]  [-:-:-]
	// 
	// █
  cpname := 
		s.Title.NewStatusItem(
			"title",
			"[-:gray:-] 識 [-:-:-][gray:green:-] [-:-:-][black:green:b] %s [-:-:-][green:gray:-][-:gray:-] [gray:-:-][-:-:-]",
			"title",
			s,
			done,
		)
  cpname.SetExpansion(2)
	s.Title.AddStatusItem(cpname)
  titact := s.Title.NewStatusItem(
			"action",
			"%s",
			"help",
			s,
			done,
		)
  // [blue:-:-][gray:blue:-][-:gray:-] ﬤ [blue:gray:-][-:-:-][black:blue:b] [black::d]?[::b] %s [gray:blue:-][-:gray:-] [blue:gray:-][-:-:-]
  titact.SetAlign(tview.AlignRight)
  titact.SetExpansion(2)
	s.Title.AddStatusItem(titact)
	s.Main.AddStatusItem(
		s.Main.NewStatusItem(
			"help",
			"[blue:-:-][gray:blue:-][-:gray:-] ﬤ [blue:gray:-][-:-:-][black:blue:b] [black::d]?[::b] %s [gray:blue:-][-:gray:-] [blue:gray:-][-:-:-]",
			"help",
			s,
			done,
		),
	)
	actMsg := s.Transient.NewStatusItem("action_str", "%s", "", s, done)
	actMsg.SetExpansion(2)
	s.Transient.AddStatusItem(actMsg)
  go statusUpdater(s)
}

func (sbs *StatusBars) SetupStatusBar() (sb *StatusBar) {
	sb = &StatusBar{
		Table: tview.NewTable(),
		bars:  sbs,
		// statuses: map[string]*StatusItem{},
	}
	sb.Table.SetBorders(false).InsertRow(0)
	sb.Table.SetSeparator(0)
	sb.statuses = make(map[string]*StatusItem)
	return
}

func debounce[T any](
	min time.Duration,
	max time.Duration,
	input chan T,
) chan T {
	output := make(chan T)

	go func() {
		var (
			buffer   T
			ok       bool
			minTimer <-chan time.Time
			maxTimer <-chan time.Time
		)

		// Start debouncing
		for {
			select {
			case buffer, ok = <-input:
				if !ok {
					return
				}
				minTimer = time.After(min)
				if maxTimer == nil {
					maxTimer = time.After(max)
				}
			case <-minTimer:
				minTimer, maxTimer = nil, nil
				output <- buffer
			case <-maxTimer:
				minTimer, maxTimer = nil, nil
				output <- buffer
			}
		}
	}()

	return output
}

func (sbs *StatusBars) FindBarForItem(name string) *StatusItem {
  if elem, ok := sbs.Title.statuses[name]; ok {
    return elem
  }
  if elem, ok := sbs.Transient.statuses[name]; ok {
    return elem
  }
  if elem, ok := sbs.Main.statuses[name]; ok {
    return elem
  }
  return nil
}

func statusUpdater(s *StatusBars) {
	for status := range updateCh {
    ss := s.FindBarForItem(status.elem)
    if ss != nil {
			ss.updates <- status
    }
	}
}

func (si *StatusItem) SetStyle(sty tcell.Style) {
	fg, bg, attr := sty.Decompose()
	si.SetBackgroundColor(bg)
	si.SetTextColor(fg)
	si.SetAttributes(attr)
}

func (si *StatusItem) SetPositioner(w int, pos func(s string, n int) string) {
	si.w = w
	si.positioner = pos
}

func (si *StatusItem) Position(txt string) string {
	return si.positioner(txt, (si.w - len(txt)))
}

func (si *StatusItem) UpdateItem(txt string) {
	txt = fmt.Sprintf(si.format, txt)
	si.bar.bars.app.QueueUpdateDraw(func() {
		si.SetText(txt)
		si.SetMaxWidth(si.w + len(txt))
	})
}

// 
// █
func (s *StatusBar) UpdateCell(status StatusUpdate) {
	s.statuses[status.elem].SetText(status.content)
	// s.statuses[status.elem].TableCell.SetExpansion(s.statuses[status.elem].TableCell.Expansion)
	s.Table.SetCell(
		0,
		s.statuses[status.elem].col,
		s.statuses[status.elem].TableCell,
	)
	// } else {
	// 	s.Transient.SetCell(0, s.statuses[status.elem].col, s.statuses[status.elem].TableCell)
	// }
}

func (sb *StatusBar) NewStatusItem(
	name, format, content string,
	bar *StatusBars,
	done <-chan struct{},
	// args ...int,
) *StatusItem {
	// row := 0
	// if len(args) > 0 {
	// 	row = args[0]
	// }
	tc := tview.NewTableCell(fmt.Sprintf(format, content))
	si := &StatusItem{
		name:      name,
		format:    format,
		TableCell: tc,
		w:         8,
		// row:       row,
		// col:       0,
		updates: make(chan *StatusUpdate),
		l:       sync.Mutex{},
		bar:     sb,
	}
	// si.SetBackgroundColor(color tcell.Color)
	si.positioner = func(s string, n int) string {
		return theme.Jcenter(s, n)
	}

	go func() {
		done := make(chan struct{})
		defer close(done)

		debouncedChan := debounce[*StatusUpdate](
			50*time.Millisecond,
			200*time.Millisecond,
			si.updates,
		)
	Done:
		for {
			select {
			case <-done:
				break Done
			case event := <-debouncedChan:
				if event == nil {
					break
				}
				si.UpdateItem(event.content)
				if event.timeOut.Milliseconds() != 0 {
					go func() {
						time.AfterFunc(event.timeOut, func() {
							updateCh <- &StatusUpdate{
								elem:    event.elem,
								content: "",
							}
						})
					}()
				}
			}
		}
	}()

	return si
}

func (s *StatusBar) AddStatusItem(si *StatusItem) {
	s.statuses[si.name] = si
	si.col = len(s.statuses) - 1
	s.Table.SetCell(0, si.col, s.statuses[si.name].TableCell)
	// if si.row == 0 {
	// } else if si.row == 1 {
	// 	s.Transient.SetCell(0, si.col, s.statuses[si.name].TableCell)
	// } else {
	// 	s.Title.SetCell(0, si.col, s.statuses[si.name].TableCell)
	// }
	// s.elems = s.elems + 1
}

// sifill := NewStatusItem("fill", " %s ", "", s, done)
// sifill.SetExpansion(2).SetAlign(tview.AlignCenter)
// s.AddStatusItem(sifill)
// sifill2 := NewStatusItem("fill2", " %s ", "", s, done)
// sifill2.SetExpansion(2).SetAlign(tview.AlignCenter)
// s.AddStatusItem(sifill2)
//  csi := NewStatusItem("color", "  %s", "", s, done)
//  csi.SetAlign(tview.AlignRight)
// s.AddStatusItem(csi)
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
