package coolor_test

import (
	"fmt"
	"testing"

	"github.com/digitallyserviced/coolors/coolor"
)

func TestHueModIncrDecr(t *testing.T) {
   cc := coolor.NewRandomCoolorColor()
   // fmt.Println(cc)
   coolor.HueMod.SetColor(cc)
   fmt.Println(coolor.HueMod.String())
   coolor.HueMod.Decr(0.0)
   fmt.Println(coolor.HueMod.String())
   for _, v := range coolor.HueMod.Above() {
      fmt.Println(coolor.NewCoolorColor(v.Hex()).TerminalPreview())
   }
   for _, v := range coolor.HueMod.Below() {
      fmt.Println(coolor.NewCoolorColor(v.Hex()).TerminalPreview())
   }
   fmt.Println()
   // fmt.Println(coolor.HueMod.Below())
   coolor.HueMod.Incr(0.0)
   for _, v := range coolor.HueMod.Above() {
      fmt.Println(coolor.NewCoolorColor(v.Hex()).TerminalPreview())
   }
   for _, v := range coolor.HueMod.Below() {
      fmt.Println(coolor.NewCoolorColor(v.Hex()).TerminalPreview())
   }
   fmt.Println()
   fmt.Println(coolor.HueMod.String())
}
