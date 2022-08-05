package coolor_test

import (
	"fmt"
	"testing"

	"github.com/digitallyserviced/coolors/coolor"
	// "github.com/gookit/goutil/dump"
)

func TestHueModIncrDecr(t *testing.T) {
   cc := coolor.NewRandomCoolorColor()
   cc.Random()
   coolor.HueMod.SetColor(cc)
   // fmt.Println(cc.TerminalPreview())
   coolor.HueMod.ColorModifier.ChannelMod.Range(cc, 12, 10)
   // coolor.HueMod.ColorModifier.Ran
      // fmt.Printf("%f\n", coolor.HueMod.GetChannelValue(cc))

   // for i := 0; i < 5; i++ {
   //    c:=coolor.HueMod.Next()
   //    fmt.Printf("%f\n", coolor.HueMod.GetChannelValue(c))
   //    c.GetCC().TerminalPreview()
   //    // ccc := coolor.NewCoolorColor(c.GetCC().Html())
   //    // _ = ccc
   // // dump.P(coolor.HueMod)
   //    fmt.Println(c.GetCC().TerminalPreview())
   //    fmt.Print("")
   // }
   // coolor.NewEditorStrip("Hue", coolor.NewCoolorEditor())
   // for _, v := range coolor.HueMod.Above() {
   // }
   for _, v := range coolor.HueMod.ColorModifier.ChannelMod.Range(cc, 24, 10){
      fmt.Println(coolor.NewCoolorColor(v.Html()).TerminalPreview())
   }
   // coolor.HueMod.Incr(0.0)
   // for _, v := range coolor.HueMod.Above() {
   //    // fmt.Println(v.GetCC().TerminalPreview())
   // }
   // for _, v := range coolor.HueMod.Below() {
   //    fmt.Println(coolor.NewCoolorColor(v.Html()).TerminalPreview())
   // }
   // fmt.Println()
   // fmt.Println(coolor.HueMod.String())
}
