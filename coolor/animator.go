package coolor

import (
	// "fmt"
	// "fmt"
	"fmt"
	"math"

	// "time"

	"github.com/charmbracelet/harmonica"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/digitallyserviced/coolors/coolor/util"
	// "github.com/samber/lo"
	// "github.com/samber/lo"
)

type animator struct {
	Animations map[string]*Animation
	ViewLayer  *tview.Pages
	LayerName  string
	*eventObserver
}

var Animator *animator

func getAnimator() *animator {
	if Animator == nil {
		Animator = &animator{
			Animations:    make(map[string]*Animation, 0),
			ViewLayer:     AppModel.anims,
			LayerName:     "",
			eventObserver: NewEventObserver("animator"),
		}
	}
	return Animator
}

func (anim *animator) GetAnimation(name string) (a *Animation) {
	aa, ok := Animator.Animations[name]
	if ok {
		a = aa
	}
	return
}

func (anim *animator) HandleEvent(ev ObservableEvent) bool {
	// log.Printf("%s", ev.String())
	return true
}

type Creanimation func(name string) *Animation

func (anim *animator) NewAnimation(
	name string,
	f Creanimation,
) (aa *Animation) {
	aa = anim.GetAnimation(name)
	if aa != nil {
    aa.Reset()
		return aa
	}
	aa = f(name)
	anim.Animations[name] = aa
	aa.Register(AllEvents, anim)
	return
}

func NewAnimatedBox(
	name string,
	b *tview.Box,
	start, target float64,
	mut PropertyMutator,
) *Animation {
	return Animator.NewAnimation(name, func(name string) (anim *Animation) {
		anim = NewAnimation()
		fm := MakeFrameMotion(
			start,
			target,
			0.0,
			NewMotion("MotionAxisX", MotionAxisX, *NewMotionMutator(mut)),
			NewBoxAnimated(b),
		)
		revfm := fm.GetReverse()
		anim.MakeKeyFrame(fm)
		anim.MakeKeyFrame(revfm)
		return

	})
}

func NewNotification(name, text string) *Animation {
	// b := tview.NewBox()
	b := tview.NewTextView()
	b.SetTitle(name)
	b.SetText(text).SetTextAlign(tview.AlignCenter)
	b.SetWordWrap(true).SetWrap(true)
	b.SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	b.SetDontClear(false)
	b.SetBackgroundColor(tcell.ColorBlack)
	w, _ := AppModel.scr.Size()
	b.SetVisible(true)
	ab := NewAnimatedBox(
		"notif",
		b.Box,
		float64(w+36),
		float64(w-40),
		RectXMutator,
	)
	AppModel.anims.AddPage("boxtest", b, false, true)
	AppModel.anims.ShowPage("boxtest")
	MainC.app.SetFocus(MainC.pages)
	MainC.app.QueueUpdateDraw(func() {
		b.SetRect(w+30, 3, 40, 4)
	})
	_ = ab
	go ab.Start()

	return ab
}

func NewBoxAnimated(b *tview.Box) *Animated {
	return &Animated{
		Box:  b,
		Item: b,
	}
}

func NewFrameAnimator() *Animation {
	return getAnimator().NewAnimation("waves", func(name string) *Animation {
		// math.▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁
		// frames := "▔🮂🮃🮄🮅🮆▉"
		// frames := "▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁"
		// frames := "▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁"
    frames := "▁▂▃▄▅▆▇█"
		// frames := "▔🮂🮃🮄🮅🮆▉█" // █▇▆▅▃▂▁
    // rframes := strings.Clone
    fms := []rune(frames)
    // fmss := fms[:0]
    fmss := append(fms[0:], []rune(util.Reverse(frames))...)
    fmt.Printf("%s\n", string(fmss))
		start := 2.0
    target := 6.0
		// target := (float64(len(fms)) / 2) + 1
		// mut := NewMotionMutator(NewBoxRectPropertyMutator(func(r Rect, m MotionValues) Rect {
		//
		// }))
		// 🮇🮈🮉🮊🮋▉  ▉▊▋▌▍▎▏🮇🮈🮉🮊🮋▉
		// ▁▂▃▅▆▇█▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁
		makeMut := func(x int) *CallbackMutator {
			return NewCallbackMutator(func(m MotionValues, i interface{}) bool {
				// fmt.Println(m.X)
				fmn := int(
					MapVal(
						math.Floor(m.X),
						0,
						float64(len(fms)-1),
						0,
						float64(len(fms)-1),
					),
				)
				frame := fmss[fmn]
				AppModel.app.QueueUpdate(func() {
					AppModel.scr.SetContent(
						x,
						AppModel.h-1,
						frame,
						nil,
						tcell.StyleDefault.Foreground(tcell.ColorBlue),
						// Background(tcell.ColorBlack),
					)
					AppModel.scr.Show()
				})
				return true
			})
		}
		anim := NewAnimation()
		// muts := make([]*CallbackMutator, 10)
		// width := AppModel.w
    width := 2
		motions := &FrameMotions{
			Motions: make([]*FrameMotion, width),
		}
		for i := 0; i < width; i++ {
			motions.Motions[i] = NewFrameMotion(
				NewMotion(
					fmt.Sprintf("FrameMotion%d", i),
					MotionNil,
					*NewMotionMutator(makeMut(i)),
				),
				&Animated{
					Item: frames,
				},
			)
			motions.Motions[i].SetTween(start, target, 10.0, 0.3)
			off := 10 * float64(i+1) * harmonica.FPS(fps)
			// fmt.Println(off)
			motions.Motions[i].SetStartOffset(off)
		}
		anim.AddKeyFrame(motions)
		return anim
	})
}

