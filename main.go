package main

import (
	"expvar"
	// "fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	// _ "github.com/divan/expvarmon"
	_ "expvar"
	"net/http"

	"github.com/digitallyserviced/coolors/coolor"
	"github.com/gookit/goutil/dump"
)

func init() {
  os.Setenv("TCELL_TRUECOLOR", "1")
	gr := expvar.NewInt("Goroutines")
	go func() {
		for range time.Tick(100 * time.Millisecond) {
			gr.Set(int64(runtime.NumGoroutine()))
		}
	}()
  // fmt.Println(os.Environ())
}

func main() {
	go http.ListenAndServe(":1234", nil)
	f, err := os.Create("dump")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	dump.Config(func(opts *dump.Options) {
		opts.Output = f
		opts.ShowFlag = dump.Ffunc | dump.Fline | dump.Ffname
	})

	coolor.StartApp()
}

// vim: ts=2 sw=2 et ft=go
