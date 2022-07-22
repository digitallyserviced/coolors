package coolor

import (
	"fmt"
	"math/rand"
	"strings"
	"text/template"

	"github.com/gdamore/tcell/v2"
)


func GenerateRandomColors(count int) []tcell.Color {
	tcols := make([]tcell.Color, count)
	for i := range tcols {
		tcols[i] = *MakeRandomColor()
	}
	return tcols
}


func MakeRandomColor() *tcell.Color {
	col := tcell.NewRGBColor(int32(randRange(0, 255)), int32(randRange(0, 255)), int32(randRange(0, 255)))
	return &col
}

func randRange(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func randomColor() tcell.Color {
	r := int32(randRange(0, 255))
	g := int32(randRange(0, 255))
	b := int32(randRange(0, 255))
	return tcell.NewRGBColor(r, g, b)
}

func getFGColor(col tcell.Color) tcell.Color {
	r, g, b := col.RGB()
	if (float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114) > 150 {
		return tcell.ColorBlack
	}
	return tcell.ColorWhite
}

func inverseColor(col tcell.Color) tcell.Color {
	r, g, b := col.RGB()
	return tcell.NewRGBColor(255-r, 255-g, 255-b)
}

func MakeTemplate(name, tpl string, funcMap template.FuncMap) func(s string, data interface{})(string) {
    status_tpl := template.New(name)
    status_tpl.Funcs(funcMap)

    status_tpl.Parse(tpl)
  
    return func(s string, data interface{}) string {
      out := &strings.Builder{}
      ntpl, ok := template.Must(status_tpl.Clone()).Parse(s)
      if ok != nil {
        fmt.Println(fmt.Errorf("%s", ok))
      }
      ntpl.Execute(out, data)
      return out.String()
    }
}

// vim: ts=2 sw=2 et ft=go
