package stack

import (
	// "fmt"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/gookit/goutil/strutil"
	"github.com/samber/lo"

	"github.com/digitallyserviced/coolors/coolor/shortcuts"
	// "github.com/digitallyserviced/coolors/coolor/util"
)

type PageStack struct {
	Stack
	Pages *tview.Pages
	Names []string
  InputScope *shortcuts.Scope
  InputScopes map[string]*shortcuts.Scope
}

func NewPageStack(p *tview.Pages) *PageStack {
	ps := &PageStack{
		Stack:       *NewStack(),
		Pages:       p,
		Names:       make([]string, 0),
		InputScope:  shortcuts.GlobalScope,
		InputScopes: make(map[string]*shortcuts.Scope),
	}
	return ps
}

func (ps *PageStack) FrontPageStacked() (stacked bool) {
  name, _ := ps.Pages.GetFrontPage()
  stacked = lo.Contains[string](ps.Names, name)
  return
}

func (ps *PageStack) RouteInput(event *tcell.EventKey, setFocus func(p tview.Primitive), hdlrs ...shortcuts.DefaultInputHandler) {
  if ps.InputScope != nil {
    ps.InputScope.HandleShortcuts(event, setFocus)
  }
	if len(hdlrs) > 0 {
		for _, v := range hdlrs {
			v(event, setFocus)
		}
	}
}

func (ps *PageStack) Pop(pages ...string) (good bool) {
	if len(pages) > 0 {
    ps.Names = lo.Without[string](ps.Names, pages...)
		for _, ipg := range ps.Stack {
			if pg, ok := ipg.(*tview.Page); ok {
				if strutil.HasOneSub(pg.Name, pages) {
					ps.Pages.RemovePage(pg.Name)
				}
			}
		}
		good= true
	}
	ipg := ps.Stack.Pop()
	if pg, ok := ipg.(*tview.Page); ok {
    ps.Names = lo.Without[string](ps.Names, pg.Name)
		ps.Pages.RemovePage(pg.Name)
		good= true
	}

  if ps.InputScope != nil {
    ps.InputScope = ps.InputScope.Parent
  }

	return
}

// type UIWidget interface {
//   tview.AnimatedPrimitive | shortcuts.ShortcutsHandler 
// }

func (ps *PageStack) Push(name string, p tview.Primitive, rz bool) *tview.Page {
  // ps.newMethod()
  // p := util.GetI[tview.Primitive](i)
  // if p == nil {
  //
  // }
	pg := ps.Pages.NewPage(name, p, rz, true)
	ps.Stack.Push(pg)
	ps.Names = append(ps.Names, name)
	ps.Pages.Addpage(pg)
	ps.Pages.ShowPage(pg.Name)
	return pg
}

func (ps *PageStack) UpdateScope(scoped shortcuts.ShortcutsHandler) {
	if scoped != nil {
		ps.InputScope = scoped.GetScope()
		ps.InputScopes[scoped.GetScope().Name] = scoped.GetScope()
	} else {
		ps.InputScope = shortcuts.GlobalScope
		ps.InputScopes[shortcuts.GlobalScope.Name] = shortcuts.GlobalScope
	}
}

type Stack []interface{}

// Create a new stack
func NewStack() *Stack {
	return &Stack{}
}

// Get size of stack
func (this *Stack) Len() int {
	return len(*this)
}

// View the top item on the stack
func (this *Stack) Peek() interface{} {
	if len(*this) == 0 {
		return nil
	}
	return (*this)[0]
}

// Pop the top item of the stack and return it
func (this *Stack) Pop() interface{} {
	if len(*this) == 0 {
		return nil
	}
	elem := this.Peek()
	*this = (*this)[:(len(*this) - 1)]
	return elem
}

// Push a value onto the top of the stack
func (this *Stack) Push(elem interface{}) {
	*this = append(*this, elem)
}
