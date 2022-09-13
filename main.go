package main

import (
	"fmt"
	"syscall"

	// "fmt"
	"log"
	_ "net/http/pprof"
	"os"

	// _ "github.com/divan/expvarmon"
	_ "expvar"
	"net/http"

	"github.com/digitallyserviced/coolors/coolor"
	"github.com/gookit/goutil/dump"
	"golang.org/x/sys/unix"
)

// func init() {
//   initVT100()
//   fmt.Println(getPosition())
//   byt := make([]byte,1024)
//   n, err := os.Stdin.Read(byt)
//   fmt.Print(n, err, string(byt))
//   os.Setenv("TCELL_TRUECOLOR", "1")
// 	gr := expvar.NewInt("Goroutines")
// 	go func() {
// 		for range time.Tick(100 * time.Millisecond) {
// 			gr.Set(int64(runtime.NumGoroutine()))
// 		}
// 	}()
//   // fmt.Println(os.Environ())
// }
const (
	ioctlReadTermios  = syscall.TCGETS
	ioctlWriteTermios = syscall.TCSETS
)

var (
	fd          int
	termios     *unix.Termios
	inputBuffer = make([]byte, 3)
)
func initVT100() {
	fd = int(os.Stdout.Fd())
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		panic(err)
	}

	newState := *termios
	newState.Lflag &^= unix.ECHO   // Disable echo
	newState.Lflag &^= unix.ICANON // Disable buffering
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		panic(err)
	}

	// hideCursor()
	// clear()
}
func getPosition() (int, int) {
	fmt.Printf("\033]4;0;?\033")

	x, y := 0, 0
 //  _,_ = x,y
 //  var dat string
	// _, _ = fmt.Scanf("\033%st", &dat)
	return x, y
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
