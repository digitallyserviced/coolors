package coolor

import (
	"fmt"
	"math"
	"strings"
	"time"

	// "github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/harmonica"
	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/samber/lo"

	// "github.com/gookit/goutil/dump"

	// "github.com/gookit/goutil/dump"

	. "github.com/digitallyserviced/coolors/coolor/anim"
	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/util"
	// "github.com/digitallyserviced/coolors/coolor/util"
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
    fm.SetStartMotion(MotionValues{start, 0})
		kf.AddMotion(fm)
		anim.AddKeyFrame(kf)

		rkf := anim.NewKeyFrame()
		rfm := rkf.NewFrameMotion(
			NewMotion("RevMotionAxisX", MotionAxisX, *NewMotionMutator(mut)),
			NewBoxAnimated(b),
		)
		rfm.Tween = fm.Tween.GetReverse()
    rfm.SetStartMotion(MotionValues{start, 0})
		rkf.AddMotion(rfm)
		anim.AddKeyFrame(rkf)

		return
	})
}

var notiIdx = 0

func Notid(t string) string {
	notiIdx += 1
	return fmt.Sprintf("%s_%d", t, notiIdx)
}

type NotificationStatus struct {
  icon string
  color tcell.Color
}

var (
  //  
  ErrorNotify = NewNotificationStatus(" ","#9F1b35")
  WarnNotify = NewNotificationStatus(" ","#FDC45C")
  InfoNotify = NewNotificationStatus(" ","#276979")
  // WarnNotify = &NotificationStatus{
  // 	icon:  " ",
  // 	color: tcell.GetColor("#FDC45C"),
  // }
  // InfoNotify = &NotificationStatus{
  // 	icon:  " ",
  // 	color: tcell.GetColor("#276979"),
  // }
)

func (ns *NotificationStatus) GetIconBox(args ...string) *tview.Box {
  icon := " "
  color := tcell.ColorBlack
  if ns == nil {
    if len(args) == 2 {
      icon = args[0]
      color = tcell.GetColor(args[1])
    }
  } else {
    icon = ns.icon
    color = ns.color
  }
    licon := MakeBoxItem("", "")
  licon.SetBackgroundColor(color)
    licon.SetBorderPadding(0, 0, 0, 0)
    licon.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
      tview.Print(screen, icon, x, y + (height / 2), width, tview.AlignCenter, NewIntCoolorColor(color.Hex()).GetFgColor())
      width=6
      height=3
      return x,y,width,height
    })
  return licon
}

func NewNotificationStatus(icon, color string) *NotificationStatus {
  ns := &NotificationStatus{
  	icon:  icon,
  	color: tcell.GetColor(color),
  }
  return ns
}

func NewNotification(id, text string, s *NotificationStatus) *Animation {
	return GetAnimator().NewAnimation(id, func(name string) *Animation {
    b := tview.NewFlex()
    tv := tview.NewTextView()
    tv.SetDynamicColors(true)
    // var licon *tview.Box
    licon := s.GetIconBox()
    tv.SetMaxLines(1)
    tv.SetWordWrap(false)
    tv.SetWrap(false)

    b.AddItem(licon, 6, 0, false)
    b.AddItem(tv, 34, 0, false)
		tv.SetText(text).SetTextAlign(tview.AlignCenter)
    tv.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetBorderSides(true, false, true, true).SetBorderColor(s.color)
		w, _ := AppModel.scr.Size()
    GetAnimator().LayerStack.Push(id, b, false)
		// b.SetRect(w-1, 3, 40, 3)
		b.SetAnimating(true)

		// tv.SetWordWrap(true).SetWrap(true)
		b.SetBorderPadding(0, 0, 0, 0)
		b.SetDontClear(true)
		tv.SetBackgroundColor(tcell.ColorBlack)
		anim := NewAnimation()
		b.SetAnimating(true)

		start := w+40
    yPos := 3
    currents := lo.CountBy[string](Animator.LayerStack.Names, func(s string) bool {
     return strings.Contains(s, "notif_")
    })
    if currents > 1 {
      yPos = (currents-1) * 4 + yPos
    }
		target := w - 46
		cm := NewCallbackMutator(func(m *MotionValues, i interface{}) bool {
			x, y, w, h := b.GetRect()
			x = int(m.X)
			b.SetRect(x, y, w, h)
			return true
		})
		kf := anim.NewKeyFrame()
		fm := kf.NewFrameMotion(
			NewMotion(
				"MotionAxisX",
				MotionAxisX,
				*NewMotionMutator(NewCallbackMutator(func(m *MotionValues, i interface{}) bool {
					x, y, w, h := b.GetRect()
          x = int(m.X)
					Animator.App.QueueUpdateDraw(func() {

						b.SetRect(x, y, w, h)
					})
					return true
				})),
			),
			NewBoxAnimated(b.Box),
		)
		fm.Tween = NewTween(float64(start), float64(target))
    fm.SetStartMotion(MotionValues{float64(start), 0})
		kf.AddMotion(fm)
		anim.AddKeyFrame(kf)

		rkf := anim.NewKeyFrame()
		cm.AddCallback(
			FinishedMutating,
			NotifyFinishedMutate(func(m *MotionValues, i interface{}) bool {
				anim.Control.Close()
				GetAnimator().Animations[id] = nil
				return true
			}),
		)
		finmot := *NewMotionMutator(cm)
		rfm := rkf.NewFrameMotion(
			NewMotion("RevMotionAxisX", MotionAxisX, finmot),
			NewBoxAnimated(b.Box),
		)
		rfm.Tween = NewTween(float64(target), float64(start))
    rfm.SetStartMotion(MotionValues{float64(target), 0})
		rkf.AddMotion(rfm)
		anim.AddKeyFrame(rkf)
		anim.Register(
			AnimationFinished,
			NewAnonymousHandlerFunc(func(o ObservableEvent) bool {
				if ObservableEventType(
					AnimationFinished | AnimationDone | AnimationCanceled,
				).Is(o.Type) {
          GetAnimator().LayerStack.Pop(id)
					return true

				}
				return true
			}),
		)

		MainC.app.SetFocus(MainC.pages)
		b.SetRect(start, yPos, 40, 3)
		MainC.app.Draw(b)
		anim.Start()

		return anim
	}, func(a *Animation) bool {
		if a != nil {
			a.Control.Close()
			a = nil
		}
		AppModel.anims.RemovePage(id)
		return true
	})
}

func NewFrameAnimator() *Animation {
	return GetAnimator().NewAnimation("waves", func(name string) *Animation {
		// math.▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁
		// frames := "▔🮂🮃🮄🮅🮆▉"
		// frames := "▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁"
		// frames := "▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁"
    nf := "🬋🬋🬋🬋🬋🬋🬋🬋"
		// frames := "▁▂▃▄▅▆▇█"
		// frames := "▔🮂🮃🮄🮅🮆▉█" // █▇▆▅▃▂▁
		fms := []rune(nf)
    frames := lo.Chunk[rune]([]rune(fms), 2)
		// fmss := append(fms[0:], []rune(util.Reverse(frames))...)
		// fmt.Printf("%s\n", string(fmss))
		start := 1.0
		target := 40.0
		// 🮇🮈🮉🮊🮋▉  ▉▊▋▌▍▎▏🮇🮈🮉🮊🮋▉
		// ▁▂▃▅▆▇█▔🮂🮃🮄🮅🮆▉██▇▆▅▃▂▁
		makeMut := func(x int) *CallbackMutator {
			return NewCallbackMutator(func(m *MotionValues, i interface{}) bool {
				// fmt.Println(m.X)
				fmn := int(
					MapVal(
            m.X,
						// math.Floor(m.X) % float64(len(frames)),
						0,
            target + 20.0,
						// float64(len(frames)-1),
						0,
						float64(len(frames)-1),
					),
				)
        fmn = util.Clamp[int](fmn, 0, len(frames)-1)
        // fmn := math.Mod(math.Floor(m.X), float64(len(frames)))
				frame := frames[int(fmn)]
				// AppModel.app.QueueUpdateDraw(func() {
					AppModel.scr.SetContent(
						2 * x,
						AppModel.h-1,
						frame[0],
						nil,
						tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(0),
						// Background(tcell.ColorBlack),
					)
					AppModel.scr.SetContent(
						2 * x+1,
						AppModel.h-1,
						frame[1],
						nil,
						tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(0),
						// Background(tcell.ColorBlack),
					)
					AppModel.scr.Show()
				// })
				return true
			})
		}
		anim := NewAnimation()
		width := 10
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
			nfm.SetTween(start, target, 1.0, 0.3)
      nfm.Tween.Spring = harmonica.NewSpring(harmonica.FPS(5), 3.0, 1.0)
			off := 10 * float64(i+1) * harmonica.FPS(5)
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
		mut8r := NewCallbackMutator(func(m *MotionValues, i interface{}) bool {
			mut.Mutate(m, i)
			return true
		})
		mut8r.FinishedCallback = MutatorFinishedCallback(
			func(idx int, startX, targetX float64, prevMv, newMv *MotionValues, i interface{}) ObservableEventType {
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
