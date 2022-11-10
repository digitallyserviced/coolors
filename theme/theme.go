package theme

import (
	// "fmt"

	"regexp"
	"strings"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
)

type Theme struct {
	HeaderBackground  tcell.Color
	GrayerBackground  tcell.Color
	SidebarBackground tcell.Color
	SidebarLines      tcell.Color
	ContentBackground tcell.Color
	Border            tcell.Color
	Primary           tcell.Color
	Secondary         tcell.Color
	TopbarBorder      tcell.Color
	InfoLabel         tcell.Color
	Styles            map[string]tcell.Style
}
// Theme{
// 	// PrimitiveBackgroundColor:    tcell.GetColor("#101010").TrueColor(),
//   PrimitiveBackgroundColor:    tcell.ColorBlack,
// 	ContrastBackgroundColor:     tcell.ColorBlue,
// 	MoreContrastBackgroundColor: tcell.ColorGreen,
// 	BorderColor:                 tcell.ColorWhite,
// 	BorderFocusColor:            tcell.ColorBlue,
// 	TitleColor:                  tcell.ColorWhite,
// 	GraphicsColor:               tcell.ColorWhite,
// 	PrimaryTextColor:            tcell.ColorWhite,
// 	SecondaryTextColor:          tcell.ColorYellow,
// 	TertiaryTextColor:           tcell.ColorGreen,
// 	InverseTextColor:            tcell.ColorBlue,
// 	ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
// }
var tvtheme *tview.Theme = &tview.Theme{
	PrimitiveBackgroundColor:    tcell.GetColor("#21252B"),
	ContrastBackgroundColor:     tcell.ColorBlue,
	MoreContrastBackgroundColor: tcell.ColorGreen,
	BorderColor:                 tcell.ColorWhite,
	BorderFocusColor:            tcell.ColorBlue,
	TitleColor:                  tcell.ColorWhite,
	// TitleColor:                  tcell.GetColor("#5c6370"),
	GraphicsColor:               tcell.ColorWhite,
	PrimaryTextColor:            tcell.ColorWhite,
	SecondaryTextColor:          tcell.ColorYellow,
	TertiaryTextColor:           tcell.ColorGreen,
	InverseTextColor:            tcell.ColorBlue,
	ContrastSecondaryTextColor:  tcell.ColorDarkCyan,
	// ContrastBackgroundColor:     0,
	// MoreContrastBackgroundColor: 0,
	// BorderColor:                 0,
	// BorderFocusColor:            0,
	// GraphicsColor:               0,
	// PrimaryTextColor:            0,
	// SecondaryTextColor:          0,
	// TertiaryTextColor:           0,
	// InverseTextColor:            0,
	// ContrastSecondaryTextColor:  0,
}
var theme *Theme

// func NewHexColor() tcell.Color {
//
// }
func init() {
  tview.Styles = *tvtheme
	theme = &Theme{ //0x303030
		HeaderBackground:  tcell.GetColor("#1C1C1C"),
		GrayerBackground:  tcell.GetColor("#282c34"),
		SidebarBackground: tcell.GetColor("#21252B"),
		ContentBackground: tcell.GetColor("#303030"),
		SidebarLines:      tcell.GetColor("#5c6370"),
		Border:            tcell.GetColor("#1C1C1C"),
		TopbarBorder:      tcell.GetColor("#5c6370"),
		InfoLabel:         tcell.GetColor("#5c6370"),
		Primary:           tcell.GetColor("#4ed6aa"),
		Secondary:         tcell.GetColor("#b5d1f6"),
		Styles:            make(map[string]tcell.Style),
	}
	theme.SetStyleFgBgAttr(
		"palette_name",
		tcell.ColorWhite,
		tcell.ColorRed,
		// tcell.GetColor("#890a37"),
		// tcell.ColorGreen,
		tcell.AttrBold,
	)
	// theme.SetStyleFgBg("action", tcell.ColorBlack, tcell.ColorYellow)
	theme.SetStyleFgBg("action", tcell.ColorBlack, tcell.ColorYellow)
	theme.SetStyleFg("list_main", tcell.ColorGreen)
	theme.SetStyleFg("list_second", tcell.ColorBlue)
	theme.SetStyleFgBg(
		"input_placeholder",
		tcell.ColorYellow,
		theme.SidebarBackground,
	)
	theme.SetStyleFgBg("input_field", tcell.ColorBlue, theme.SidebarBackground)
	theme.SetStyleFgBg(
		"input_autocomplete",
		tcell.ColorRed,
		theme.SidebarBackground,
	)
}

func (t *Theme) SetStyleFgBgAttr(
	name string,
	fg, bg tcell.Color,
	attr tcell.AttrMask,
) *tcell.Style {
	sty := t.SetStyle(name)
	t.Styles[name] = sty.Foreground(fg).Background(bg).Attributes(attr)
	return sty
}

func (t *Theme) SetStyleFg(name string, fg tcell.Color) *tcell.Style {
	sty := t.SetStyle(name)
	t.Styles[name] = sty.Foreground(fg)
	return sty
}

func (t *Theme) SetStyleFgBg(name string, fg, bg tcell.Color) *tcell.Style {
	sty := t.SetStyle(name)
	t.Styles[name] = sty.Foreground(fg).Background(bg)
	return sty
}

func (t *Theme) SetStyle(name string) *tcell.Style {
	sty := &tcell.Style{}
	t.Styles[name] = *sty
	return sty
}

func (t *Theme) Get(name string) *tcell.Style {
	sty, ok := t.Styles[name]
	if ok {
		return &sty
	}
	return &tcell.Style{}
}

func Jright(s string, n int) string {
	if n < 0 {
		n = 0
	}
	return strings.Repeat(" ", n) + s
}

func Jleft(s string, n int) string {
	if n < 0 {
		n = 0
	}
	return s + strings.Repeat(" ", n)
}

func Jcenter(s string, n int) string {
	if n < 0 {
		n = 0
	}
	// div := ((2 * (n-len(s))) / 2) + 1
  rem := n - len(s)
  rem = ((rem * 2) / 2) / 2
  rem = rem
	return strings.Repeat(" ", rem) + s + strings.Repeat(" ", rem)
}

func (t *Theme) FixedSize(w int) string {
	return strings.Repeat(" ", w)
}

func GetTheme() *Theme {
	// tags := color.GetColorTags()
	// tags["infolabel"] = RgbHex256toCode("5c6370", false)
	// tags["sckey"] = RgbHex256toCode("fda47f", false)
	// tags["scicon"] = RgbHex256toCode("7aa4a1", false)
	// tags["scname"] = RgbHex256toCode("7aa4a1", false)
	// tags["scdesc"] = RgbHex256toCode("5a93aa", false)
	// tags["colorinfolabel"] = RgbHex256toCode("7aa4a1", false)
	// tags["colorinfovalue"] = RgbHex256toCode("fda47f", false)
	//"#cb7985"
	//"#ff8349"
	//"#2f3239", "#e85c51", "#7aa4a1", "#fda47f", "#5a93aa", "#ad5c7c", "#a1cdd8", "#ebebeb"
	//"#4e5157", "#eb746b", "#8eb2af", "#fdb292", "#73a3b7", "#b97490", "#afd4de", "#eeeeee"

	return theme
}

const (
	TplFgRGB = "38;2;%d;%d;%d"
	TplBgRGB = "48;2;%d;%d;%d"
	FgRGBPfx = "38;2;"
	BgRGBPfx = "48;2;"
)

const (
	TplFg256 = "38;5;%d"
	TplBg256 = "48;5;%d"
	Fg256Pfx = "38;5;"
	Bg256Pfx = "48;5;"
)

var (
	rxNumStr  = regexp.MustCompile("^[0-9]{1,3}$")
	rxHexCode = regexp.MustCompile("^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$")
)

func RgbHex256toCode(val string, isBg bool) (code string) {
	if len(val) == 6 && rxHexCode.MatchString(val) { // hex: "fc1cac"
		code = color.HEX(val, isBg).String()
	} else if strings.ContainsRune(val, ',') { // rgb: "231,178,161"
		code = strings.Replace(val, ",", ";", -1)
		if isBg {
			code = BgRGBPfx + code
		} else {
			code = FgRGBPfx + code
		}
	} else if len(val) < 4 && rxNumStr.MatchString(val) { // 256 code
		if isBg {
			code = Bg256Pfx + val
		} else {
			code = Fg256Pfx + val
		}
	}
	return
}
