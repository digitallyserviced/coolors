package tree

import (
	// "os"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	// "github.com/gookit/goutil/dump"
	// "github.com/gookit/goutil/fsutil/finder"
	// "github.com/gookit/goutil/dump"
)

var _ tview.Primitive = &FileTree{}

type FileTree struct {
	theme *Theme
	*tview.Box
	collapsed   bool
	view        *tview.TreeView
	root        *tview.TreeNode
	onSelect    func(node *FSNode)
	onChanged   func(node *FSNode)
	onOpen      func(node *FSNode)
	AfterDraw   []func(*FileTree)
	filters     []string
	localRoot   *tview.TreeNode
	virtualRoot *tview.TreeNode
}

func get(node *tview.TreeNode) *FSNode {
	ref := node.GetReference()
	if ref == nil {
		return nil
	}
	return ref.(*FSNode)
}

func NewFileTree(theme *Theme) *FileTree {
	view := tview.NewTreeView().
		SetTopLevel(1)

	ft := &FileTree{
		Box:   tview.NewBox(),
		theme: theme,
		view:  view,
	}
	ft.Box.SetBackgroundColor(0)
	ft.Box.SetDontClear(false)
	view.SetIndicateOverflow(true)

	view.SetBorder(true)
	view.SetBorderPadding(0, 1, 1, 1)
	view.SetBorderSides(true, false, true, true)
	view.SetGraphicsColor(theme.SidebarLines)
	view.SetBackgroundColor(theme.SidebarBackground)
	// view.SetBorderVisible(false)
	view.SetDontClear(true)

	view.SetSelectedFunc(func(node *tview.TreeNode) {
		ft.selected(node)
	})

	view.SetChangedFunc(func(node *tview.TreeNode) {
		ft.changed(node)
	})

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return ft.inputCapture(event)
	})

	// Disable mouse scroll
	view.SetMouseCapture(
		func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			return ft.mouseCapture(action, event)
		},
	)

	// ft.Box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
	// 	arrow := ">"
	//    dump.P(x,y,width,height)
	// 	if ft.collapsed {
	// 		x, y, w, h := ft.Box.GetInnerRect()
	// 		_, _, _, _ = x, y, w, h
	// 		centerY := h / 2
	// 		centerTop, centerBottom := centerY-(centerY/2), centerY+(centerY/2)
	// 		ft.Box.SetRect(x, y, 2, h)
	// 		ft.Box.SetBorder(true).SetBorderColor(theme.SidebarLines).SetBorderSides(false, false, false, true)
	// 		// ft.Box.DrawForSubclass(screen, ft)
	// 		tview.Print(screen, arrow, x, centerTop, 1, tview.AlignCenter, theme.ContentBackground)
	// 		tview.Print(screen, arrow, x, centerBottom, 1, tview.AlignCenter, theme.ContentBackground)
	// 		return x, y, 2, h
	// 	} else {
	// 		x, y, w, h := ft.Box.GetRect()
	//    dump.P(x,y,w,h)
	// 		_, _, _, _ = x, y, w, h
	// 		ft.Box.SetBorder(false)
	// 		ft.Box.SetRect(x, y, w, h)
	// 		ft.view.SetRect(x, y, w, h)
	// 		ft.view.Draw(screen)
	// 		return x, y, w, h
	// 	}
	// })

	return ft
}

// Primitive interface
func (ft *FileTree) ToggleCollapsed() {
	ft.SetCollapsed(!ft.GetCollapsed())
}

func (ft *FileTree) GetCollapsed() bool {
	return ft.collapsed
}

func (ft *FileTree) SetCollapsed(c bool) {
	ft.collapsed = c
}

func (ft *FileTree) Draw(screen tcell.Screen) {
	ft.Box.DrawForSubclass(screen, ft)
	ft.view.Box.DrawForSubclass(screen, ft.view)
	ft.view.Draw(screen)
	// ft.tr.DrawForSubclass(screen, ft)
}

func (ft *FileTree) GetRect() (int, int, int, int) {
	// if ft.collapsed {
	// 	return ft.Box.GetRect()
	// }
	return ft.view.GetRect()
}

func (ft *FileTree) SetRect(x, y, width, height int) {
	ft.Box.SetRect(x, y, width, height)
	ft.view.SetRect(x, y, width, height)
}

func (ft *FileTree) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return ft.view.InputHandler()
}

func (ft *FileTree) Focus(delegate func(p tview.Primitive)) {
	ft.view.Focus(delegate)
}

func (ft *FileTree) HasFocus() bool {
	return ft.view.HasFocus()
}

func (ft *FileTree) Blur() {
	ft.view.Blur()
}

func (ft *FileTree) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return ft.view.MouseHandler()
}

func (ft *FileTree) inputCapture(event *tcell.EventKey) *tcell.EventKey {
	n := ft.view.GetCurrentNode()
	fsnode := get(n)
	if event.Modifiers() == tcell.ModShift {
		return event
	}
	switch event.Key() {
	case tcell.KeyLeft:
		if fsnode.Virtual {
			return event
		}
		parent := ft.GetParentNode(fsnode)
		if fsnode.IsDir && fsnode.IsExpanded() {
			fsnode.Collapse()
		} else if !ft.IsRoot(parent) {
			ft.SetCurrent(parent)
		}
		return nil

	case tcell.KeyRight:
		// if fsnode.Virtual && !fsnode.IsExpanded() {
		// fsnode.Node.Expand()
		// return nil
		// }
		// if !fsnode.IsExpanded() {
			fsnode.Expand()
		// }
		return nil

	case tcell.KeyEnter:
		if ft.onOpen != nil {
			ft.onOpen(fsnode)
		}
		return nil

		// case tcell.KeyRune:
		// 	switch event.Rune() {
		// 	case 'K':
		// 		return nil // noop
		//
		//    case '^':
		// 		ft.SetRoot(get(ft.root).CreateParent())
		// 		return nil
		//
		// 	case 'o':
		// 		if ft.onOpen != nil {
		// 			ft.onOpen(fsnode)
		// 		}
		// 		return nil
		// 	}
	}
	return event
}

func (ft *FileTree) mouseCapture(
	action tview.MouseAction,
	event *tcell.EventMouse,
) (tview.MouseAction, *tcell.EventMouse) {
	switch action {
	case tview.MouseScrollUp:
		return action, nil
	case tview.MouseScrollDown:
		return action, nil
	default:
		return action, event
	}
}

func (ft *FileTree) selected(node *tview.TreeNode) {
	fsnode := get(node)
	if fsnode.IsExpanded() {
		fsnode.Collapse()
	} else if fsnode.IsDir {
		fsnode.Expand()
	} else {
		if ft.onSelect != nil {
			ft.onSelect(fsnode)
		}
	}
}

func (ft *FileTree) changed(node *tview.TreeNode) {
	if ft.onChanged != nil {
		ft.onChanged(get(node))
	}
}

func (ft *FileTree) GetParentNode(fsnode *FSNode) *FSNode {
	var currParent *tview.TreeNode
	ft.root.Walk(func(node, parent *tview.TreeNode) bool {
		if node == fsnode.Node {
			currParent = parent
			return false
		}
		return true
	})

	if currParent == nil {
		parent := get(ft.localRoot).CreateParent()
		ft.UpdateLocalRoot(parent)
		return parent
	}
	return get(currParent)
}

func (ft *FileTree) SetRoot(fsnode *FSNode) {
	ft.root = fsnode.Node
	ft.view.SetRoot(ft.root)
	// ft.root.ClearChildren()
	// if ft.virtualRoot != nil {
	// 	// ft.root.AddChild(ft.virtualRoot)
	// 	// defer ft.view.SetCurrentNode(ft.virtualRoot)
	// }
	// if ft.localRoot != nil {
	// 	ft.root.AddChild(ft.localRoot)
	// 	defer ft.view.SetCurrentNode(ft.localRoot)
	// }
	ft.root.Expand()
}
func (ft *FileTree) SetVirtualRoot(fsnode *FSNode) {
	ft.virtualRoot = fsnode.Node
	ft.virtualRoot.ClearChildren()
	ft.root.AddChild(fsnode.Node)
  ft.view.SetCurrentNode(ft.virtualRoot)
}
func (ft *FileTree) SetLocalRoot(fsnode *FSNode) {
	ft.localRoot = fsnode.Node
	ft.localRoot.ClearChildren()
	ft.root.AddChild(fsnode.Node)
  ft.view.SetCurrentNode(ft.localRoot)
}
func (ft *FileTree) UpdateLocalRoot(fsnode *FSNode) {
	if fsnode != nil {
		ft.localRoot.ClearChildren()
    ft.localRoot.AddChild(fsnode.Node)
		fsn := get(ft.localRoot)
		fsn.Path = fsnode.Path
		// ft.localRoot.Expand()
		fsnode.Node.SetSelectable(true)
		ft.view.SetCurrentNode(fsnode.Node)

		if ft.onChanged != nil {
			ft.AfterDraw = append(ft.AfterDraw, func(ft *FileTree) {
				ft.onChanged(get(ft.view.GetCurrentNode()))
			})
		}
	}
}

func (ft *FileTree) SetCurrent(fsnode *FSNode) {
	if fsnode != nil {
		ft.view.SetCurrentNode(fsnode.Node)
	}
}

func (ft *FileTree) IsRoot(fsnode *FSNode) bool {
	return ft.root == fsnode.Node
}

func (ft *FileTree) LoadFiltered(dir string, patterns []string) {
	ft.filters = patterns
	ft.UpdateLocalRoot(newRootFsnode(dir))
}

func (ft *FileTree) Load(dir string) {
	ft.SetLocalRoot(newRootFsnode(dir))
}

func (ft *FileTree) OnSelect(fn func(node *FSNode)) {
	ft.onSelect = fn
}

func (ft *FileTree) OnOpen(fn func(node *FSNode)) {
	ft.onOpen = fn
}

func (ft *FileTree) OnChanged(fn func(node *FSNode)) {
	ft.onChanged = fn
}
