package coolor

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/dump"
)

type Square struct {
	color *CoolorColor
	*tview.Box
	cp         *CoolorColorsPalette
	items      []*rSquare
	num, count int
}

var nextIdx int

type rSquare struct {
	*tview.Box
	color      *CoolorColor
	idx        int
	p          int
	x, y, size int
	num, count int
}

func (sq *Square) NewSquare(cp *CoolorColorsPalette, x, y, size, count int, rsq int) *rSquare {
	nextIdx += 1
	count--
	nrsq := &rSquare{
		Box:   tview.NewBox(),
		num:   count,
		idx:   nextIdx,
		x:     x,
		y:     y,
		size:  size,
		count: count,
		color: cp.RandomColor(),
		p:     rsq,
	}
	nrsq.SetRect(x, y, size, size/2)
	if count >= 0 && pRand(count, 0.5) {
		hs := size / 2
		// if pRand(count + size){
		//   hs = size / 4
		// }
		if pRand(count, 0.2) {
			sq.items = append(sq.items, sq.NewSquare(cp, x, y, hs, count, nrsq.idx))
		}
		if pRand(count, 0.4) {
			sq.items = append(sq.items, sq.NewSquare(cp, x+hs, y, hs, count, nrsq.idx))
		}
		if pRand(count, 0.6) {
			sq.items = append(sq.items, sq.NewSquare(cp, x+hs, y+hs, hs, count, nrsq.idx))
		}
		if pRand(count, 0.8) {
			sq.items = append(sq.items, sq.NewSquare(cp, x, y+hs, hs, count, nrsq.idx))
		}
	}

	return nrsq
}

func pRand(count int, pc float64) bool {
	r := rand.Float64()
	p := MapVal(float64(count), 0, float64(count-1), pc, 0)
	return r > p
}

func NewRecursiveSquare(cp *CoolorColorsPalette, count int) *Square {
	nextIdx = 0
	sq := &Square{
		Box:   tview.NewBox(),
		cp:    cp.GetPalette(),
		items: make([]*rSquare, 0),
		count: count,
		color: cp.RandomColor(),
	}
	return sq
}

func Lerp(a, b, t float64) float64 {
	return (1.0-t)*a + b*t
}
func InvLerp(a, b, t float64) float64 {
	return (t - a) / (b - a)
}
func MapVal(val, imin, imax, omin, omax float64) float64 {
	return Lerp(omin, omax, InvLerp(imin, imax, val))
}

func (sq *Square) TopInit(count int) {
	x, y, w, _ := sq.GetRect()
	// cc,_  := sq.color.pallette.GetSelected()
	// x,y,w,_ = cc.GetInnerRect()
	sq.items = append(sq.items, sq.NewSquare(sq.cp, x, y, w, count, -1))
	// sq.SetBackgroundColor(tcell.ColorBlack)
	// sq.SetBorder(true).SetTitle(fmt.Sprintf("%d", count))
	// sq.ColorPass(sq.cp)
}

func (sq *Square) ColorPass(cp *CoolorPaletteMainView) {
	sq.color = cp.RandomColor()
	for _, v := range sq.items {
		if v == nil {
			continue
		}
		if v.num != 0 {
			v.color = cp.RandomColor()
		}
	}

}

// ██▀⬛⬛⬠⬡⬢⬣⬤⬟⬧⬯⭓⭔⭕⭖⭗⭘⭙⮰⮱⮲⮳⮴⮵⮶⮷⮸⮹⮺⮻⮼⮽⮾⮿⯀⯁⯂⯃⯄⯅⯆⯇⯈⯉⯊⯋⯌⯍⯎⯏
// var fullBlock = '⬤'

// var fullBlock = '▀'
// var fullBlock = '█'

func (sq *Square) DrawBlock(fg string, screen tcell.Screen, bx int, by int, bwidth int, bheight int, idx int) (int, int, int, int) {

	fgcol := tcell.GetColor(fg)
	if idx == -1 {
		dump.P(fmt.Sprintf("shiz %d %d %d %d %d %s", bx, by, bwidth, bheight, idx, NewIntCoolorColor(fgcol.Hex()).TerminalPreview()))
	}
	for y := by; y < by+bheight; y++ {
		for x := bx; x < bx+bwidth; x++ {
			c := '█'
			if idx == -1 && x == bx {
				dump.P(fmt.Sprintf("fux %d %d %d %d %d %s", x, y, bwidth, bheight, idx, NewIntCoolorColor(fgcol.Hex()).TerminalPreview()))
				screen.SetContent(x, y, c, nil, tcell.Style.Background(tcell.StyleDefault, fgcol).Foreground(fgcol))
				continue
			} // c := fullBlock
			if idx%2 == 0 {
				c = '⬤'
				screen.SetContent(x, y, c, nil, tcell.Style.Background(tcell.StyleDefault, fgcol).Foreground(tcell.Color(NewIntCoolorColor(fgcol.Hex()).GetFgColor().Hex())))
			} else {
				screen.SetContent(x, y, c, nil, tcell.Style.Foreground(tcell.StyleDefault, fgcol))
			}
		}
	}

	return bx, by, bwidth, bheight
}

func (sq *Square) DrawItems(screen tcell.Screen) {
	x, y, width, height := sq.GetInnerRect()
	_, _, _, _ = x, y, width, height

	mh := width / 2 / 2
	leftover := height - mh
	extras := math.Ceil(float64(leftover))
	dump.P(mh, leftover, extras)
	sq.DrawBlock(NewIntCoolorColor(tcell.Color236.Hex()).Html(), screen, x, y+int(extras)-1, width, mh-(mh/2), -1)

	// x, y, width, height := sq.cp.GetRect()
	// DrawBlock(sq.color.Html(), screen, x, y, width, height, -1)
	col := sq.cp.RandomColor().Html()
	count := sq.count
	if len(sq.items) > 1 {
		for j := 0; j < count; j++ {
			for i := len(sq.items) - 1; i > 0; i-- {
				v := sq.items[i]
				if v != nil {
					if v.p != v.idx-1 {
						col = sq.cp.RandomColor().Html()
					}
					// if v.y > y-v.size+(mh-(mh/2)) {
					// 	continue
					// }
					dump.P(fmt.Sprintf("%d %d %d %d %d %d %s", v.x, v.y+int(extras), v.size, v.p, v.idx, v.num, sq.color.TerminalPreview()))
					if v.num == j {
						sq.DrawBlock(col, screen, v.x, v.y+int(extras), v.size, v.size/2, v.idx)
					}
				}
			}

		}
		// sq.DrawBlock(NewIntCoolorColor(tcell.Color236.Hex()).Html(), screen, x, y+mh, width, 1, -1)
		// for i := len(sq.items) - 1; i > 0; i-- {
		// 	v := sq.items[i]
		// 	if v != nil {
		// 		if v.p != v.idx-1 {
		// 			col = sq.cp.RandomColor().Html()
		// 		}
		// 		dump.P(fmt.Sprintf("%d %d %d %d %d %d %s", v.x, v.y, v.size, v.p, v.idx, v.num, sq.color.TerminalPreview()))
		// 		DrawBlock(col, screen, v.x, v.y, v.size, v.size/2)
		// 	}
		// }
	}
}

func (sq *Square) Draw(screen tcell.Screen) {
	sq.DrawItems(screen)
}

//   func (sq *Square) Draw(screen tcell.Screen)  {
// DrawBlock(NewIntCoolorColor(tcell.ColorBlack.Hex()).Html(), screen, x, y, width, height)
// v.DrawItems(screen)
// col = sq.cp.RandomColor().Html()
//   if len(sq.items) > 1 {
//     // sq.SetBorder(true)
//     // tview.Borders = MyBorderStyle
//   }
//        // sq.Box.SetDontClear(false)
//    sq.Box.DrawForSubclass(screen, sq)
//    // tview.Borders = OrigBorders
//   // newsq:= NewSquare(x,y,width-4,5)
//   // newsq.Draw(screen)
//    // sq.items = append(sq.items, NewSquare(x,y,width-4,count))
//
//    // tview.Print()
//    // if sq.num != 0 {
//    //   return
//    // }
//        // sq.Box.SetDontClear(true)
//    sq.Box.DrawForSubclass(screen, sq)
//    for _, v := range sq.items {
//     x,y,width,height := sq.GetInnerRect()
//     dump.P(fmt.Sprintf("%d %d %d %d %d %d",x,y,width,height,sq.num, sq.count))
//     if v.num > 0 {
//       v.SetBackgroundColor(*v.color.color)
//     } else {
//       v.SetBackgroundColor(tcell.ColorBlack)
//     }
//      if v != nil && len(v.items) == 1 {
//        // v.SetDontClear(true)
//        v.Box.DrawForSubclass(screen, v)
//        // v.SetDontClear(false)
//      }
//    }
//    // for i := sq.count; i > 0; i-- {
//    //
//    // }
// }

// func (sq *Square) DrawForSubclass(screen tcell.Screen, p tview.Primitive) {
// }
// vim: ts=2 sw=2 et ft=go
