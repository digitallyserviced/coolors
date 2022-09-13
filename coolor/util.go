package coolor

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/errorx"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

type (
	Severity int
	Status   struct {
		Message  string
		Severity Severity
	}
)

const (
	Unknown Severity = iota
	Refresh
	Healthy
	Warning
	Alert
)

const (
	cssInteger       = "[-\\+]?\\d+%?"
	cssNumber        = "[-\\+]?\\d*\\.\\d+%?"
	cssUnit          = "(?:" + cssNumber + ")|(?:" + cssInteger + ")"
	permissiveMatch3 = "[\\s|\\(]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")\\s*\\)?"
	permissiveMatch4 = "[\\s|\\(]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")[,|\\s]+(" + cssUnit + ")\\s*\\)?"
	rgb              = "rgb" + permissiveMatch3
	rgba             = "RGBA" + permissiveMatch4
	hsl              = "hsl" + permissiveMatch3
	hsla             = "hsla" + permissiveMatch4
	hsv              = "hsv" + permissiveMatch3
	hsva             = "hsva" + permissiveMatch4
	// hex3             = `#?([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})`
	// hex6             = `#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})`
	// hex4             = `#?([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})([0-9a-fA-F]{1})`
	// hex8             = `#?([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})`
	hex3 = `(#[0-9a-fA-F]{3})\b`
	hex6 = `(#[0-9a-fA-F]{6})\b`
	// 0;#090300;1;#db2d20;2;#01a252;3;#fded02;4;#01a0e4;5;#a16a94;6;#b5e4f4;7;#a5a2a2;8;#5c5855;9;#e8bbd0;10;#3a3432;11;#4a4543;12;#807d7c;13;#d6d5d4;14;#cdab53;15;#f7f7f7
	// printf "\033]10;#4a4543;#f7f7f7;#4a4543\007"
	// printf "\033]17;#a5a2a2\007"
	// printf "\033]19;#4a4543\007"
	// printf "\033]5;0;#4a4543\007"
	set4BitDynamicColors string = "\033]4;%s\007"
	dynamicColorIndex    string = "%d;%s"
	setTextFgBgCursor    string = "\033]10;%s\007"
	setBgColor           string = "\033]17;%s\007"
	setSelectionFgColor  string = "\033]19;%s\007"
	setDynamicColorBold  string = "\033]5;%d;%s\007"
)

type ColorRepParserFunc func(cr *ColorRep) *Color

// Color representations hsl(n,n,n), #7bafcd, rgba(0.6,0.6,0.6,255)
// assume only need at most 4 value to make it in the color model/space
type ColorRep struct {
	parseRegex     *regexp.Regexp
	scanFormat     string
	v1, v2, v3, v4 float64
	parseFunc      ColorRepParserFunc
}

// var CssHex6 ColorRep
var (
	CssHex6 *ColorRep
)

func init() {
	CssHex6 = NewColorRep(hex6, "%s", func(cr *ColorRep) *Color {
		return &Color{cr.v1 / 255.0, cr.v2 / 255.0, cr.v3 / 255.0}
	})
}

func (cr *ColorRep) ParseSingle(s string) *Color {
	n, err := fmt.Sscanf(s, cr.scanFormat, &cr.v1, cr.v2, cr.v3, cr.v4)
	if err != nil {
		panic(err)
	}
	if n > 0 {
		return cr.parseFunc(cr)
	}
	return nil
}

func (cr *ColorRep) Matches(s string) bool {
	found := cr.parseRegex.Match([]byte(s))
	dump.P(s, found)
	return found
}

func NewColorRep(regex, scanFormat string, p ColorRepParserFunc) *ColorRep {
	rx := regexp.MustCompile(regex)
	cr := &ColorRep{
		parseRegex: rx,
		scanFormat: scanFormat,
		v1:         0,
		v2:         0,
		v3:         0,
		v4:         0,
		parseFunc:  p,
	}
	return cr
}



func StringColorizer(s string) (string, *CoolorColorsPalette) {
	for _, v := range []*ColorRep{CssHex6} {
		if v.Matches(s) {
			content, cols := v.FindAndColorize(s)
			if cols.Len() > 0 {
				return content, cols
			}
		}
	}
	return "", nil
}

func (cr *ColorRep) FindAndColorize(sc string) (string, *CoolorColorsPalette) {
	cols := NewCoolorColorsPalette()
	newsc := cr.parseRegex.ReplaceAllStringFunc(sc, func(s string) string {
		col := NewCoolorColor(s)
		cols.AddCoolorColor(col)
		return col.TVPreview()
	})
	return newsc, cols
	// if matchIdxs := cr.parseRegex.FindAllStringSubmatchIndex(sc, -1); len(matchIdxs) > 0 {
	// 	for _, c := range matchIdxs {
	//      var cuint int32 = 0
	//      str := sc[c[0]:c[1]]
	//      // fmt.Printf("%q", str)
	//      // n, err := fmt.Sscanf(strings.TrimSpace(str), "%s", &cuint)
	//      // if err != nil || n == 0{
	//      //   panic(fmt.Errorf("%v %v %v %v", c, n, err, str))
	//      // }
	//      // tcol := tcell.NewHexColor(cuint)
	//      // str = NewIntCoolorColor(cuint).TVPreview()
	//      newsc =
	// 		if len(c) == 2 {
	// 			colors = append(colors, str)
	// 		}
	// 	}
	//    dump.P(sc)
	//    return colors
	// }
	// if match := cr.parseRegex.FindAllStringSubmatch(sc, -1); match != nil {
	// 	colors := make([]string, 0)
	// 	for _, c := range match {
	// 		if len(c) == 2 {
	// 			colors = append(colors, c[1])
	// 		}
	//      dump.P(c)
	// 	}
	//    return colors
	// }
	// regexp.MustCompile(reg).FindAllSubmatch()
}

var colorRegexes = []string{rgb, rgba, hsl, hsla, hsv, hsva, hex3, hex6} //
func genSimilarHslColor(
	tcol Color,
	f func(r, h, s, l float64) Color,
) interface{} {
	rand.Seed(time.Now().UnixNano())
	h, s, l := tcol.Hsl()
	adjust := rand.Float64() * 360
	return f(adjust, h, s, l)
}

func checkHslColorDistance(tcol, tcol2 Color, distance float64) bool {
	return tcol2.DistanceRgb(tcol) <= distance
}

func GetColorName(col tcell.Color) string {
	for n, v := range tcell.ColorNames {
		if col == v {
			return n
		}
	}
	return ""
}

func errAss[R any](v R, e error) R {
	iserr := func(v interface{}) {
		if v == nil {
			return
		}
		e, ok := v.(error)
		if ok {
			doLog(errorx.WithPrevf(e, "%V", v))
			panic(e)
		}
		return
	}
	iserr(e)
	// for _, v := range vars {
	//   iserr(v)
	//   r, ok := v.(R)
	//   if ok {
	//     return r
	//   } else {
	//     continue
	//   }
	// }
	return v
}

func checkErrX(err error, vars ...interface{}) bool {
	if err != nil {
		doLog(errorx.WithPrevf(err, "%V", vars))
		return false
	}
	return true
}

func checkErr(err error) {
	if err != nil {
		doLog(err)
		panic(err)
	}
}


func handleSignals() {
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signal_chan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
			case syscall.SIGINT:
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
			GetStore().Close()
		}
	}()

	code := <-exit_chan
	os.Exit(code)
}

func startBoltStats() {
	go func() {
		// Grab the initial stats.
		// prev := Store.Bolt().Stats()
		tick := time.NewTicker(10 * time.Second)
    tick.Reset(1000 * time.Millisecond)

		for {
			// Wait for 10s.
			select {
			case <-tick.C:
			TrimSeentCoolors(RecentCoolorsMax)
				// stats := Store.Bolt().Stats()
				// diff := stats.Sub(&prev)

				// Encode stats to JSON and print to STDERR.
				// json.NewEncoder(os.Stderr).Encode(diff)

				// Save stats for the next loop.
				// prev = stats
			}
		}
	}()
}

func SeentColor(from string, cc *CoolorColor, src Referenced) {
	// fmt.Println(from, cc)
	if MainC == nil || MainC.eventNotifier == nil {
		return
	}
	MainC.eventNotifier.Notify(*MainC.eventNotifier.NewObservableEvent(ColorSeentEvent, from, cc, src))
}

func IfElse[T any](b bool, a,c T) T {
  if b {
    return a
  }
  return c
}

var IfElseStr = IfElse[string]
var IfElseUint64 = IfElse[uint64]
var IfElseTCol = IfElse[*tcell.Color]
var IfElseCCol = IfElse[*CoolorColor]

// custom msgpack decoding function for bolthold (faster than gobs)
func dec(data []byte, value interface{}) error {
	return msgpack.Unmarshal(data, value)
}

// custom msgpack encoding function for bolthold (faster than gobs)
func enc(value interface{}) ([]byte, error) {
	return msgpack.Marshal(value)
}

// var _ msgpack.CustomEncoder = (*CoolorColorPalette)(nil)
// var _ msgpack.CustomDecoder = (*CoolorColorPalette)(nil)
//
var _ msgpack.CustomEncoder = (*CoolorColor)(nil)
var _ msgpack.CustomDecoder = (*CoolorColor)(nil)

// Debug message categories.
// const (
// 	LogConn  = LogMask(1 << iota)   // Connection events
// 	LogState                        // State changes
// 	LogCmd                          // Command execution
// 	LogRaw                          // Raw data stream excluding literals
// 	LogGo                           // Goroutine execution
// 	LogAll   = LogMask(1<<iota - 1) // All messages
// 	LogNone  = LogMask(0)           // No messages
// )

// var logMasks = []enumName{
// 	{uint32(LogAll), "LogAll"},
// 	{uint32(LogConn), "LogConn"},
// 	{uint32(LogState), "LogState"},
// 	{uint32(LogCmd), "LogCmd"},
// 	{uint32(LogRaw), "LogRaw"},
// 	{uint32(LogNone), "LogNone"},
// }
//
// func (v LogMask) String() string   { return enumString(uint32(v), logMasks, false) }
// func (v LogMask) GoString() string { return enumString(uint32(v), logMasks, true) }

// enumName associates an enum value with its name for printing.
type enumName struct {
	v uint32
	s string
}

// enumString converts a flag-based enum value into its string representation.
func enumString(v uint32, names []enumName, goSyntax bool) string {
	s := ""
	for _, n := range names {
		if v&n.v == n.v && (n.v != 0 || v == 0) {
			if len(s) > 0 {
				s += "+"
			}
			if goSyntax {
				s += "imap."
			}
			s += n.s
			if v &= ^n.v; v == 0 {
				return s
			}
		}
	}
	if len(s) > 0 {
		s += "+"
	}
	return s + "0x" + strconv.FormatUint(uint64(v), 16)
}
// vim: ts=2 sw=2 et ft=go
