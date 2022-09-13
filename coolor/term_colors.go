package coolor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/creack/pty"
	"golang.org/x/term"
)

const (
	tesc  = "\033Ptmux;\033\033]"
	tcesc = "\007\033\\"
	esc   = "\033]"
	cesc  = "\007"
)

func EscReqColor(idx int) (string, string) {
	oe := esc
	ce := cesc
	// if tmux variable is set then we need extra escape codes
	if os.Getenv("TMUX") != "" {
		oe = tesc
		ce = tcesc
		// fmt.Println("tmux")
	}
	// fmt.Printf("%q", fmt.Sprintf("%s4;%d;?%s\n", oe, idx, ce))
	return oe, ce
}

const (
	ioctlReadTermios  = syscall.TCGETS
	ioctlWriteTermios = syscall.TCSETS
)

var (
	fd          int
	termios     *unix.Termios
	inputBuffer = make([]byte, 32)
)

func write(c byte) {
	fmt.Printf("%c", c)
}

func get(ptmx *os.File) string {
	n, _ := ptmx.Read(inputBuffer)
	return string(inputBuffer[:n])
}

func size() (int, int) {
	ws, _ := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	return int(ws.Col), int(ws.Row)
}

func hideCursor() {
	fmt.Printf("\x1b[?25l")
}

func showCursor() {
	fmt.Printf("\x1b[?25h")
}

func clear() {
	fmt.Print("\x1b[2J")
}

func setCursor(x, y int) {
	fmt.Printf("\x1b[%d;%dH", y, x)
}

func reset() {
	showCursor()
	clear()
	setCursor(0, 0)
	_ = unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)
}

func initVT100(fd int) {
	// fd = int(os.Stderr.Fd())
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		panic(err)
	}

	newState := *termios
	// newState.Lflag &^= unix.ECHO   // Disable echo
	newState.Lflag &^= unix.ICANON // Disable buffering
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		panic(err)
	}
}

func queryOscColor(ptmx *os.File, idx int) {
	s, e := EscReqColor(0)
	fmt.Fprintf(ptmx, "%s%d;?%s\n", s, idx, e)
	fmt.Fprintf(os.Stdin, "\n")
}

func queryIndexColor(ptmx *os.File, idx int) {
	s, e := EscReqColor(0)
	fmt.Fprintf(ptmx, "%s4;%d;?%s\n", s, idx, e)
	fmt.Fprintf(os.Stdin, "\n")
}

func readOscColor(ptmx *os.File, idx int) (r, g, b uint) {
	var i, rr, gg, bb uint = 0, 0, 0, 0
	txt := get(ptmx)
	// fmt.Printf("%q", txt)
	n, err := fmt.Sscanf(txt, "\x1b]%d;rgb:%02x%02x/%02x%02x/%02x%02x\x1b", &i, &r, &rr, &g, &gg, &b, &bb)
	if err != nil || n != 7 {
		// log.Fatal(fmt.Errorf("input does not match format %d != 7", n))
	}
	// _ = fmt.Sprintf("%d: %d %d %d %d",n, i,r,g,b, txt)
	return
}

func readIndexColor(ptmx *os.File, idx int) (r, g, b uint) {
	var i, rr, gg, bb uint = 0, 0, 0, 0
	txt := get(ptmx)
	// fmt.Println(txt)
	n, err := fmt.Sscanf(txt, "\x1b]4;%d;rgb:%02x%02x/%02x%02x/%02x%02x\x1b", &i, &r, &rr, &g, &gg, &b, &bb)
	if err != nil || n != 7 {
		// log.Fatal(fmt.Errorf("input does not match format %d != 7", n))
	}
	// _ = fmt.Sprintf("%d: %d %d %d %d",n, i,r,g,b, txt)
	return
}

func QueryColorScheme(n int) []Color {
	initVT100(int(os.Stderr.Fd()))
	clear()
	cols := make([]Color, n+5)
	for i := 0; i < n; i++ {
		queryIndexColor(os.Stderr, i)
		r, g, b := readIndexColor(os.Stdin, i)
		cols[i] = Color{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}
	}

	ps := []int{10, 11, 12, 17, 19}
	for i, v := range ps {
		queryOscColor(os.Stderr, v)
		r, g, b := readOscColor(os.Stdin, v)
		cols[n+i] = Color{float64(r) / 255.0, float64(g) / 255.0, float64(b) / 255.0}
	}

	clear()
	reset()
	return cols
}

func test() error {
	// Create arbitrary command.
	c := exec.Command("zsh")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(ptmx.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(ptmx.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	// go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	// go func() {_, _ = io.Copy(os.Stderr, ptmx)}()
	// i := queryIndexColor(ptmx, 2)

	return nil
}

func init() {
	// initVT100(int(os.Stdout.Fd()))
	// cols := QueryColorScheme(16)
	// for i, v := range cols {
	// 	fmt.Printf("%d %s", i, v.GetCC().TerminalPreview())
	// }
	// fmt.Println(cols)
	// fmt.Println(queryIndexColor(2))
	// queryIndexColor(os.Stderr,2)
	// str := readIndexColor(os.Stdin, 2)
	// fmt.Println(len(str),str)
	// byt := make([]byte, 1024)
	// n, err := os.Stdin.Read(byt)
	// fmt.Print(n, err, byt)
}
