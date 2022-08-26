package coolor

import (
	// "fmt"
	"regexp"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/gookit/color"
)

type Theme struct {
	SidebarBackground tcell.Color
	SidebarLines      tcell.Color
	ContentBackground tcell.Color
	Border            tcell.Color
	TopbarBorder      tcell.Color
	InfoLabel         tcell.Color
}

var theme *Theme
func GetTheme() *Theme {
  tags := color.GetColorTags()
  tags["infolabel"] = RgbHex256toCode("5c6370", false)
  tags["sckey"] = RgbHex256toCode("fda47f", false)
  tags["scicon"] = RgbHex256toCode("7aa4a1", false)
  tags["scname"] = RgbHex256toCode("7aa4a1", false)
  tags["scdesc"] = RgbHex256toCode("5a93aa", false)
  tags["colorinfolabel"] = RgbHex256toCode("7aa4a1", false)
  tags["colorinfovalue"] = RgbHex256toCode("fda47f", false)
  //"#cb7985" 
  //"#ff8349" 
  //"#2f3239", "#e85c51", "#7aa4a1", "#fda47f", "#5a93aa", "#ad5c7c", "#a1cdd8", "#ebebeb" 
  //"#4e5157", "#eb746b", "#8eb2af", "#fdb292", "#73a3b7", "#b97490", "#afd4de", "#eeeeee" 
	if theme == nil {
		theme = &Theme{
			SidebarLines:      tcell.NewHexColor(0x5c6370),
			SidebarBackground: tcell.NewHexColor(0x21252B),
			ContentBackground: tcell.NewHexColor(0x282c34),
			Border:            tcell.NewHexColor(0x5c6370),
			TopbarBorder:      tcell.NewHexColor(0x5c6370),
			InfoLabel:         tcell.NewHexColor(0x5c6370),
		}
	}
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
