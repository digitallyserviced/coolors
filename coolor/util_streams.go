package coolor

import(
  "github.com/digitallyserviced/coolors/coolor/util"
)

// TakeNColors function ﳑ
func TakeNColors(
	done <-chan struct{},
	valueStream <-chan interface{},
	num int,
) []Color {
	colors := make([]Color, 0)
	for cv := range util.TakeN(done, valueStream, num) {
		if cv == nil {
			continue
		}
		col := cv.(Color)
		colors = append(colors, col)
	}
	return colors
}
