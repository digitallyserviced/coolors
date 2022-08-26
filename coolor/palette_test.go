package coolor_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/digitallyserviced/coolors/status"
	"github.com/digitallyserviced/coolors/coolor"
	// "github.com/gookit/goutil/dump"
)

//   ﴞ ﵆ ﵃ ﵄ 🅴 🅰 🅱 🅲 🅳 🅴 🅵 🅶 🅷 🅸 🅹 🅺 🅻 🅼 🅽 🅾 🅿 🆀 🆁 🆂 🆃 🆄 🆅 🆆 🆇 🆈 🆉 🄰 🄱 🄲 🄳 🄴 🄵 🄶 🄷 🄸 🄹 🄺 🄻 🄼 🄽 🄾 🄿 🅀 🅁 🅂 🅃 🅄 🅅 🅆 🅇 🅈 🅉 🅐 🅑 🅒 🅓 🅔 🅕 🅖 🅗 🅘 🅙 🅚 🅛 🅜 🅝 🅞 🅟 🅠 🅡 🅢 🅣 🅤 🅥 🅦 🅧 🅨 🅩 ① ② ③ ④ ⑤ ⑥ ⑦ ⑧ ⑨ ⑩ ⑪ ⑫ ⑬ ⑭ ⑮ ⑯ ⑰ ⑱ ⑲ ⑳ ⓪ ⓫ ⓬ ⓭ ⓮ ⓯ ⓰ ⓱ ⓲ ⓳ ⓴ ⓵ ⓶ ⓷ ⓸ ⓹ ⓺ ⓻ ⓼ ⓽ ⓾ ⓿ ㉑ ㉒ ㉓ ㉔ ㉕ ㉖ ㉗ ㉘ ㉙ ㉚ ㉛ ㉜ ㉝ ㉞ ㉟ ㉠ ㊱ ㊲ ㊳ ㊴ ㊵ ㊶ ㊷ ㊸ ㊹ ㊺ ㊻ ㊼ ㊽ ㊾ ㊿ 

func TestColorStream(t *testing.T) {
  // cbp := coolor.BlankCoolorShadePalette()
	// tcol, _ := coolor.Hex("#be1685")
	// tcol, _ := coolor.Hex("#00be00")
	tcol, _ := coolor.Hex("#ffffff")
	// tcol, _ := coolor.Hex("#000000")
  done := make(chan struct{})
  defer close(done)
	colors := coolor.RandomShadesStream(tcol, 0.05)
  colors.Status.SetProgressHandler(coolor.NewProgressHandler(func(u uint32) {
    status := fmt.Sprintf("Found Shades (%d / %d)", u, colors.Status.GetItr())
    fmt.Println(status)
		status.NewStatusUpdate("action", fmt.Sprintf("Found Shades (%d / %d)", u, colors.Status.GetItr()))
    // dump.P(u,colors.Status.Itr)
  },func(i uint32) {
    status := fmt.Sprintf("Iterating Shades (%d)", i)
    fmt.Println(status)
		// status.NewStatusUpdate("action_str", fmt.Sprintf("Iterating Shades (%d)", i))
  }))
  cbp := make([]string, 0)
  colors.Run(done)
  for _, v := range coolor.TakeNColors(done, colors.OutColors, 20) {
    cbp = append(cbp, v.GetCC().TerminalPreview())
  }
  fmt.Println(strings.Join(cbp, " "))
  // close(done)
}
