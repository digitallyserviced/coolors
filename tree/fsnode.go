package tree

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	// "github.com/gookit/goutil/dump"
	// "github.com/gookit/goutil/fsutil"
)

type FSNode struct {
	Name     string
  Virtual bool
  Children func(*FSNode) []*FSNode
	Path     string
  Icon string
	IsDir    bool
	Size     int64
	Node     *tview.TreeNode
	Mode     fs.FileMode
	ModTime  time.Time
	MimeType string
}

func newRootFsnode(path string) *FSNode {
	stat, _ := os.Stat(path)
	return newFsnode(filepath.Dir(path), stat)
}

func NewRootNode(path string) *tview.TreeNode {
	fsnode := newRootFsnode(path)

	if !fsnode.Node.IsExpanded() {
		fsnode.Node.Expand()
		fsnode.ReadChildren()
	}

	return fsnode.Node
}
// 
// '▀', '▁', '▂', '▃', '▅', '▆' , '▇', '█', '▉', '▊', '▋', '▌', '▍', '▎', '▏', 
//'▐', '░', '▒', '▓', '▔', '▕', '▖', '▗', '▘', '▙', '▚', '▛', '▜', '▝', '▞', '▟', '🭰', '🭱', '🭲', '🭳', '🭴', '🭶', '🭷', '🭸', '🭹', '🭺', '🭻', '🭼', '🭽', '🭾', '🭿', '🮀', '🮁', '🮂', '🮃', '🮄', '🮅', '🮆', '🮇', '🮈', '🮉', '🮊', '🮋', '🮌', '🮍', '🮎', '🮏', '🮐', '🮑', '🮒', '■', '□', '▢', '▣', '▥', '▦', '▧', '▨', '▩', '▪', '▫', '▬', '▭', '▮', '▯', '▰', '▱', '▲', '△', '▴', '▵', '▶', '▷', '▸', '▹', '►', '▻', '▼', '▽', '▾', '▿', '◀', '◁', '◂', '◃', '◄', '◅', '◆', '◇', '◈', '◉', '◊', '○', '◌', '◍', '◎', '●', '◐', '◑', '◒', '◓', '◔', '◕', '◖', '◗', '◘', '◙', '◚', '◛', '◜', '◝', '◞', '◟', '◠', '◡', '◢', '◣', '◤', '◥', '◦', '◧', '◨', '◩', '◪', '◫', '◬', '◭', '◮', '◯', '◰', '◱', '◲', '◳', '◴', '◵', '◶', '◷', '◸', '◹', '◺', '◻', '◼', '◽', '◾', '◿', '⑀', '⑁', '⑂', '⑃', '⑄', '⑅', '⑆', '⑇', '⑈', '⑉', '⑊', 

func NewVirtualNode(name, icon, path string) *FSNode {
  return newVirtualNode(name, icon, path, nil)
}
func newVirtualNode(name, icon, path string, children []*FSNode) *FSNode {
	fsnode := &FSNode{
      Virtual: true,
		Name:     name,
    Icon: icon,
    IsDir: true,
		Path:     path,
		Size:     -1,
	}
	node := tview.NewTreeNode("").
		SetSelectable(true)

	node.SetExpanded(true)


  fsnode.Node = node
  node.SetReference(fsnode)
  node.SetText(fsnode.Title())
    return fsnode
}

func newFsnode(parentPath string, stat fs.FileInfo) *FSNode {
	name := stat.Name()
	fpath := filepath.Join(parentPath, name)

	fsnode := &FSNode{
		Name:     name,
		Path:     fpath,
		IsDir:    stat.IsDir(),
		Size:     -1,
		Mode:     stat.Mode(),
		ModTime:  stat.ModTime(),
		// MimeType: mime,
	}

	if !stat.IsDir() {
		fsnode.Size = stat.Size()
	} else {
		go func() {
			size, _ := dirSize(fpath)
			fsnode.Size = size
			log.Printf("dir size: %v | %s", size, fpath)
		}()
	}

	fsnode.Node = createNode(fsnode)

	return fsnode
}


func NewNode(parentPath string, file fs.FileInfo) *tview.TreeNode {
	fsnode := newFsnode(parentPath, file)
	return fsnode.Node
}

func (n *FSNode) Expand() {
	n.ReadChildren()
	n.Node.Expand()
	n.Node.SetText(n.Title())
}

func (n *FSNode) Collapse() {
	n.Node.ClearChildren()
	n.Node.Collapse()
	n.Node.SetText(n.Title())
}

func (n *FSNode) IsExpanded() bool {
	return n.Node != nil && n.Node.IsExpanded()
}

func (n *FSNode) readChildren(node *FSNode) {
  if n.Virtual && n.Path == "" {
    return
  }
	if n.IsDir {
		n.Node.ClearChildren()

		files, err := ioutil.ReadDir(n.Path)
		if err != nil {
			panic(err)
		}

		nodes := []*tview.TreeNode{}

		for _, file := range files {
  dump.P(file.Name())
			fpath := filepath.Join(n.Path, file.Name())

			if node != nil && node.Path == fpath {
				nodes = append(nodes, node.Node)
			} else {
				nodes = append(nodes, NewNode(n.Path, file))
			}
		}

		sort.Slice(nodes, func(i, j int) bool {
			a := nodes[i].GetReference().(*FSNode)
			b := nodes[j].GetReference().(*FSNode)

			if a.IsDir == b.IsDir {
				return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name)) < 0
			}

			return a.IsDir
		})

		for _, node := range nodes {
			n.Node.AddChild(node)
		}
    n.Node.SetExpanded(true)
	}
}

func (n *FSNode) ReadChildren() {
	n.readChildren(nil)
}

func (n *FSNode) CreateParent() *FSNode {
	dir := filepath.Dir(n.Path)
	log.Printf("Create parent for: %s => %s", n.Path, dir)

	if n.Path == dir {
		return n
	}

	rnode := newRootFsnode(dir)

	rnode.readChildren(n)
	rnode.Node.SetExpanded(true)

	return rnode
}

func (n *FSNode) Title(args... string) string {
	icon := "  "
	if n.IsDir {
		if n.IsExpanded() {
			icon = " ﱮ"
		} else {
			icon = " "
		}
	}
  if n.Virtual && n.Icon != "" {
    eicon := ""
		if n.IsExpanded() {
			eicon = ""
		}
    icon = fmt.Sprintf("%s %s", eicon, n.Icon)
  } else if n.Virtual {
    icon = "  "
  }
  str := n.Name
  if len(args) > 0 {
    str = strings.Join(args, " ")
  }
	return fmt.Sprintf("%s %s%s", icon, str, strings.Repeat(" ", 50))
}

func createNode(n *FSNode) *tview.TreeNode {
	node := tview.NewTreeNode(n.Title()).
		SetReference(n).
		SetSelectable(true)

	if n.IsDir {
		node.SetColor(tcell.ColorBlue)
	}

	node.SetExpanded(true)

	return node
}

func getFileContentType(file *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}
