package coolor

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	// "github.com/gookit/goutil/dump"

	. "github.com/digitallyserviced/coolors/coolor/anim"
	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/util"
	// "github.com/digitallyserviced/coolors/coolor/zzlog"
)

func NewAnimatedBox(
	name string,
	b *tview.Box,
	start, target float64,
	mut PropertyMutator,
) *Animation {
	b.SetAnimating(true)
	return Animator.NewAnimation(name, func(name string) (anim *Animation) {
		anim = NewAnimation()

		kf := anim.NewKeyFrame()
		fm := kf.NewFrameMotion(
			NewMotion("MotionAxisX", MotionAxisX, *NewMotionMutator(mut)),
			NewBoxAnimated(b),
		)
		fm.Tween = NewTween(start, target)
		kf.AddMotion(fm)
		anim.AddKeyFrame(kf)

		rkf := anim.NewKeyFrame()
		rfm := rkf.NewFrameMotion(
			NewMotion("RevMotionAxisX", MotionAxisX, *NewMotionMutator(mut)),
			NewBoxAnimated(b),
		)
		rfm.Tween = fm.Tween.GetReverse()
		rkf.AddMotion(rfm)
		anim.AddKeyFrame(rkf)

		return
	})
}

func NewNotification(name, text string) *Animation {
	return GetAnimator().NewAnimation("notif", func(name string) *Animation {
		// anim := NewAnimatedBox(
		// 	"notif",
		// 	b.Box,
		// 	float64(w+40),
		// 	float64(w-46),
		// 	RectXMutator,
		// )
		// b.Box.SetAnimating(true)
		// anim.Reanimate = func(anim *Animation) bool {
		// 	if AppModel.anims.HasPage("notif") {
		// 		AppModel.anims.ShowPage("notif")
		// 	}
		// 	// MainC.app.QueueUpdateDraw(func() {
		// 	b.SetRect(w+30, 3, 40, 4)
		// 	MainC.app.Draw(b)
		// 	// })
		//
		// 	return true
		// }
	b := tview.NewTextView()
	b.SetTitle(name)
	b.SetText(text).SetTextAlign(tview.AlignCenter)
	w, _ := AppModel.scr.Size()
    AppModel.app.QueueUpdateDraw(func() {

	AppModel.anims.AddPage("notif", b, false, true)
	AppModel.anims.ShowPage("notif")
	b.SetRect(30, 3, 40, 4)
    })
	// MainC.app.Draw(AppModel.anims)
	// b.SetRect(w+30, 3, 40, 4)
		b.SetAnimating(true)

	b.SetWordWrap(true).SetWrap(true)
	b.SetBorder(false).SetBorderPadding(0, 0, 1, 1)
	b.SetDontClear(false)
	b.SetBackgroundColor(tcell.ColorBlack)
	b.SetVisible(true)
	// b.SetRect(w+30, 3, 40, 4)
    anim := NewAnimation()
    b.SetAnimating(true)

    start := w+30
    target := w - 46
    cm := NewCallbackMutator(func(m MotionValues, i interface{}) bool {
        x,y,w,h := b.GetRect()
        x = int(m.X)
        b.SetRect(x, y, w, h)
        return true
      })
		kf := anim.NewKeyFrame()
		fm := kf.NewFrameMotion(
			NewMotion("MotionAxisX", MotionAxisX, *NewMotionMutator(NewCallbackMutator(func(m MotionValues, i interface{}) bool {
        x,y,w,h := b.GetRect()
        x = int(m.X)
        b.SetRect(x, y, w, h)
        return true
      }))),
			NewBoxAnimated(b.Box),
		)
		fm.Tween = NewTween(float64(start), float64(target))
		kf.AddMotion(fm)
		anim.AddKeyFrame(kf)

		rkf := anim.NewKeyFrame()
    cm.AddCallback(FinishedMutating, NotifyFinishedMutate(func(m MotionValues, i interface{}) bool {
      anim.Control.Close()
      GetAnimator().Animations["notif"]=nil
      return true
    }))
    finmot := *NewMotionMutator(cm)
    // finmot.Mutator.Finished(idx int, startX float64, targetX float64, prevMv .MotionValues, newMv .MotionValues, i interface{})
    // finmot.
		rfm := rkf.NewFrameMotion(
			NewMotion("RevMotionAxisX", MotionAxisX, finmot),
			NewBoxAnimated(b.Box),
		)
		rfm.Tween = NewTween(float64(target), float64(start))
		rkf.AddMotion(rfm)
		anim.AddKeyFrame(rkf)
    anim.Register(AnimationFinished, NewAnonymousHandlerFunc(func(o ObservableEvent) bool {
      if !o.Type.Is(AnimationFinished){
        return true
      }
        return true
    }))
    

		// b.Box.SetRect(w+30, 3, 40, 4)
		MainC.app.SetFocus(MainC.pages)
		MainC.app.QueueUpdateDraw(func() {
			b.SetRect(w+20, 3, 40, 4)
			MainC.app.Draw(b)
    // b.SetAnimating(true)
			// b.SetRect(123, 3, 40, 4)
		})
		// b.SetRect(w+30, 3, 40, 4)
		anim.Start()

		return anim
	},func(a *Animation) bool{
      if a != nil {
        a.Control.Close()
        a = nil
      }
      AppModel.anims.RemovePage("notif")
      return false
    })
}

func NewFrameAnimator() *Animation {
	return GetAnimator().NewAnimation("waves", func(name string) *Animation {
		// math.▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁
		// frames := "▔🮂🮃🮄🮅🮆▉"
		// frames := "▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁"
		// frames := "▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁"
		frames := "▁▂▃▄▅▆▇█"
		// frames := "▔🮂🮃🮄🮅🮆▉█" // █▇▆▅▃▂▁
		fms := []rune(frames)
		fmss := append(fms[0:], []rune(util.Reverse(frames))...)
		// fmt.Printf("%s\n", string(fmss))
		start := 2.0
		target := 6.0
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
		width := 2
		kf := anim.NewKeyFrame()
		for i := 0; i < width; i++ {
			nfm := kf.NewFrameMotion(
				NewMotion(
					fmt.Sprintf("FrameMotion%d", i),
					MotionNil,
					*NewMotionMutator(makeMut(i)),
				),
				&Animated{
					Item: frames,
				},
			)
			nfm.SetTween(start, target, 10.0, 0.3)
			off := 10 * float64(i+1) * harmonica.FPS(DefaultFPS)
			// fmt.Println(off)
			nfm.SetStartOffset(off)
			kf.AddMotion(nfm)
		}
		anim.AddKeyFrame(kf)
		return anim
	})
}

type DynamicTargetUpdaterCallback func(a *Animation) (int, int, bool)

func NewDynamicTargetAnimation(
	name string,
	b interface{},
	start, target float64,
	mut PropertyMutator,
	tgtUpd DynamicTargetUpdaterCallback,
) *Animation {
	return GetAnimator().NewAnimation(name, func(name string) *Animation {
		anim := NewAnimation()
		mut8r := NewCallbackMutator(func(m MotionValues, i interface{}) bool {
			mut.Mutate(m, i)
			return true
		})
		mut8r.FinishedCallback = MutatorFinishedCallback(
			func(idx int, startX, targetX float64, prevMv, newMv MotionValues, i interface{}) ObservableEventType {
				current, _, _ := tgtUpd(anim)
				if idx > 1 && math.Abs(prevMv.Xvelocity-newMv.Xvelocity) < 0.001 &&
					math.Abs(float64(current)-newMv.X) < 1 {
					return AnimationPaused
				}
				return AnimationPlaying
			},
		)

		i := 0
		kf := anim.NewKeyFrame()
		nfm := kf.NewFrameMotion(
			NewMotion(
				fmt.Sprintf("DynamicTargetAnimation%d", i),
				MotionNil,
				*NewMotionMutator(mut8r),
			),
			&Animated{
				Item: b,
			},
		)

		nfm.SetTween(start, target, 7.0, 1.0)
		kf.AddMotion(nfm)
		anim.Frames.IdleTime = 400 * time.Millisecond
		anim.Frames.Add(kf)
		return anim
	})
}
