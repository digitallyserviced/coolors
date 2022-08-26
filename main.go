package main

import (
	"log"
	"os"

	_ "net/http/pprof"

	"github.com/digitallyserviced/coolors/coolor"
	"github.com/gookit/goutil/dump"
)



func main() {
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
