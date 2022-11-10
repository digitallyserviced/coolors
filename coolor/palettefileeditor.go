package coolor

import (
	// "fmt"

	// "fmt"
	// "fmt"
	"io/ioutil"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/pgavlin/femto"

	"github.com/digitallyserviced/tview"

	"github.com/digitallyserviced/coolors/runtime"
)

type PaletteFileEditor struct {
	*femto.View
	colorscheme   femto.Colorscheme
	files         *femto.RuntimeFiles
	filePath      string
	currentBuffer *femto.Buffer
	swallowed     bool
	escapeTime    time.Time
}

func NewPaletteFileEditor(content, path string) *PaletteFileEditor {
	buffer := femto.NewBufferFromString(string(content), path)
	root := femto.NewView(buffer)
	pfe := &PaletteFileEditor{
		View:       root,
		files:      runtime.Files,
		swallowed:  false,
		escapeTime: time.Now(),
	}
	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "dukedark-tc"); monokai != nil {
		if data, err := monokai.Data(); err == nil {
			pfe.colorscheme = femto.ParseColorscheme(string(data))
		}
	}
	//  pfe.colorscheme["statusline"]=tcell.StyleDefault.Reverse(false).Background(tcell.GetColor("#232526")).Foreground(tcell.GetColor("#656866,"))
	// root.SetRuntimeFiles(runtime.Files)
	// root.SetColorscheme(pfe.colorscheme)
	// pfe.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	return event
	// })
	pfe.SetOnFocus(func() {
		pfe.Swallow()
		pfe.SetOnFocus(nil)
		pfe.SetOnBlur(pfe.Vomit)
	})
	// pfe.SetOnBlur(func(){
	//   pfe.Vomit()
	//   pfe.SetOnFocus(pfe.Swallow)
	// })
	return pfe
}

func (pfe *PaletteFileEditor) Vomit() {
	ev := tcell.NewEventKey(tcell.KeyF41, 0, tcell.ModNone)
	MainC.app.QueueEvent(ev)
	pfe.SetOnFocus(pfe.Swallow)
	pfe.SetOnBlur(nil)
}

func (pfe *PaletteFileEditor) Swallow() {
	ev := tcell.NewEventKey(tcell.KeyF40, 0, tcell.ModNone)
  pfe.swallowed = true
	MainC.app.QueueEvent(ev)
	pfe.SetOnFocus(nil)
	pfe.SetOnBlur(pfe.Vomit)
}

func (pfe *PaletteFileEditor) saveBuffer() error {
	return ioutil.WriteFile(pfe.filePath, []byte(pfe.currentBuffer.SaveString(false)), 0600)
}

func (pfe *PaletteFileEditor) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return pfe.WrapInputHandler(func(ek *tcell.EventKey, f func(p tview.Primitive)) {
		switch ek.Key() {
		case tcell.KeyEscape:
      // fmt.Println(time.Since(pfe.escapeTime).String())
			if time.Since(pfe.escapeTime) <= time.Duration(time.Millisecond*1000) {
				pfe.swallowed = false
				ev := tcell.NewEventKey(tcell.KeyF41, 0, tcell.ModNone)
				MainC.app.QueueEvent(ev)
			}
			pfe.escapeTime = time.Now()
		case tcell.KeyCtrlT:
			pfe.OpenFile("js/template.js")
		case tcell.KeyCtrlS:
			pfe.saveBuffer()
		default:
			if pfe.swallowed == true {
				// pfe.WrapInputHandler(func(ek *tcell.EventKey, f func(p tview.Primitive)) {
				//   handler := pfe.View.InputHandler()
				//   handler(ek, f)
				// })
			}
		}
		handler := pfe.View.InputHandler()
		handler(ek, f)
	})
}

func (pfe *PaletteFileEditor) OpenFile(path string) {
	content, err := ioutil.ReadFile(path)
	if checkErrX(err, pfe) {
		pfe.SetContent(string(content), path)
	}
}

func (pfe *PaletteFileEditor) SetContent(content, path string) {
	pfe.currentBuffer = femto.NewBufferFromString(string(content), path)
	pfe.filePath = path
	pfe.View.OpenBuffer(pfe.currentBuffer)
	pfe.currentBuffer.Settings["tabsize"] = float64(2)
	// pfe.currentBuffer.Settings["filetype"] = "toml"
	pfe.currentBuffer.Settings["matchbrace"] = true
	pfe.currentBuffer.Settings["basename"] = true
	pfe.currentBuffer.Settings["matchbraceleft"] = true
	pfe.currentBuffer.Settings["rmtrailingws"] = true
	pfe.currentBuffer.Settings["saveundo"] = true
	pfe.currentBuffer.Settings["scrollbar"] = true
	pfe.currentBuffer.Settings["statusline"] = true
	pfe.currentBuffer.Settings["tabstospaces"] = true
	pfe.SetRuntimeFiles(pfe.files)
	pfe.View.SetColorscheme(pfe.colorscheme)
}
