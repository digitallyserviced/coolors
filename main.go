package main

import (
	"fmt"
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"
	"runtime/pprof"

	"github.com/digitallyserviced/coolors/coolor"
	"github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
)

func main() {
	// setBordersChars()
	go func() { log.Println(http.ListenAndServe("localhost:6060", nil)) }()
	f, err := os.Create("dump")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	dump.Config(func(opts *dump.Options) {
		opts.Output = f
		opts.ShowFlag = dump.Ffunc | dump.Fline | dump.Ffname
	})
	fmem, err := os.Create("mem.mprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(fmem)
	defer fmem.Close()
	app := tview.NewApplication()
	statusbar := status.NewStatusBar(app)

	helpbar := &coolor.HelpBar{}
	helpbar.Init(app)
	// statusbar.AddStatusItem("Clock:", startClockStatus())
	// pages := tview.NewPages()
	pages := coolor.NewMainContainer(app)
	flexview := tview.NewFlex().SetDirection(tview.FlexRow)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ch := event.Rune()
		kp := event.Key()
		switch {
		case ch == 'q' || kp == tcell.KeyEscape:
			app.Stop()
			return nil
		}
		return event
	})

	flexview.
		AddItem(helpbar, 1, 1, false).
		AddItem(pages, 0, 10, true).
		AddItem(statusbar, 1, 1, false)
		// _ = statusbar
	if err := app.SetRoot(flexview, true).Run(); err != nil {
		panic(err)
	}
	fmem.Close()

	coolor.MainC.CloseConfig()
	colors := pages.GetPalette()
	for i := 0; i < colors.GetItemCount(); i++ {
		pcol := colors.GetItem(i)
		fmt.Printf("%s ", pcol)
	}
}

// vim: ts=2 sw=2 et ft=go
