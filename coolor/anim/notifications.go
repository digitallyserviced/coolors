package anim

import (
	"fmt"
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"

	"github.com/digitallyserviced/coolors/coolor/events"
)

func getNotifier() *NotificationManager {
	if NotifMgr == nil {
    NotifMgr = NewNotificationManager()
  }
  return NotifMgr
}


type Notification struct {
	idx int
	*tview.Box
	title, msg string
	*Animation
	*NotificationMods
	m *NotificationManager
}

type NotificationState int

const (
	NotificationNop NotificationState = 1 << iota
	NotificationIdle
	NotificationInit
	NotificationEmit
	NotificationOnMove
	NotificationTarget
	NotificationBail
	NotificationDie
)

type NotificationMod func(s NotificationState, n *Notification, nm *NotificationManager) (bool, string)

type NotificationMods struct {
	mods []NotificationMod
}

type NotificationManager struct {
	Notifications      []*Notification
	max, width, height int
	targetX, targetY   int
	offsetX, offsetY   int
	startX, startY     int
}

var (
	NotificationInitialPos = func(x, y int) NotificationMod {
		return func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "initial_pos"
			if s == NotificationInit {
				_, _, w, h := n.Box.GetRect()
				n.Box.SetRect(x, y, w, h)
        ok = true
			}
      return
		}
	}
	NotificationTimeoutCancel = func(t *time.Timer) NotificationMod {
		return func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "to_cancel"
			if s == NotificationDie || s == NotificationOnMove {
				t.Stop()
        ok = true
			}
      return
		}
	}
	NotificationTimeout = func(t time.Duration) (*time.Timer, NotificationMod) {
		var timer *time.Timer
		return timer, func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "to"

			if s == NotificationEmit {
				timer = time.AfterFunc(t, func() {
					n.Bail()
				})
        ok = true
			}
      return
		}
	}
	NotificationAnyCallback = func(ns NotificationState, f NotificationMod) NotificationMod {
		return func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "custom"

			if ns == s {
				f(s, n, nm)
        ok = true
			}
      return
		}
	}

	NotificationKilled = func(id string) NotificationMod {
		return func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "die"

			if s == NotificationDie {
				if Animator.AnimLayer.HasPage(id) {
					Animator.AnimLayer.RemovePage(id)
				}
        ok = true
			}
      return
		}
	}
	NotificationPlaced = func(id string) NotificationMod {
		return func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "placed"

			if s == NotificationInit {
				if !Animator.AnimLayer.HasPage(id) {
					Animator.AnimLayer.AddPage(id, n, false, false)
				}
				Animator.AnimLayer.ShowPage(id)
				Animator.AnimLayer.SendToFront(id)
        ok=true
			}
      return
		}
	}
	NotificationFire = func(kf *KeyFrame) NotificationMod {
		return func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "FIRE"

			if s == NotificationEmit {
				kf.Jump()
        ok=true
			}
      return
		}
	}
	NotificationBailout = func(kf *KeyFrame) NotificationMod {
		return func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "bailout"

			if s == NotificationBail {
				kf.Jump()
        ok=true
			}
      return
		}
	}

	NotificationSize = func(w, h int) NotificationMod {
		return func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok =true 
      desc = "size"

			if s == NotificationInit {
        x, y, _, _ := n.Box.GetRect()
        n.Box.SetRect(x, y, w, h)
      }
      return
		}
	}
)

func NewNotificationBox(
	name, text string,
	mods ...NotificationMod,
) *Notification {
	b := tview.NewTextView()
	b.SetTitle(name)
	b.SetText(text).SetTextAlign(tview.AlignCenter)
	b.SetWordWrap(true).SetWrap(true)
	b.SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	b.SetDontClear(false)
	b.SetBackgroundColor(tcell.ColorBlack)
	w, _ := Animator.Screen.Size()
	b.SetVisible(true)
	b.SetRect(w+30, 3, 40, 4)
	notif := &Notification{
		idx:       0,
		Box:       b.Box,
		title:     "shit",
		msg:       "shit",
		Animation: &Animation{},
		NotificationMods: &NotificationMods{
			mods: mods,
		},
	}
	return notif
}

func NewNotificationManager() *NotificationManager {
	width, height := 40, 5
	w, h := Animator.Screen.Size()
	nm := &NotificationManager{
		Notifications: make([]*Notification, 0),
		width:         width,
		height:        height,
		offsetX:       5,
		offsetY:       2,
	}
	nm.startX = w + nm.offsetX + nm.width
	nm.startY = h + nm.offsetY + nm.height
	nm.max = (h - (nm.offsetY * 2)) / (nm.height + nm.offsetY)
	nm.targetX = w - width + nm.offsetX
	nm.targetY = nm.offsetY
	nm.max = (h - (nm.offsetY * 2)) / (nm.height + nm.offsetY)
	// var o Observer = nm
	return nm
}

func (n *Notification) Bail(mods ...NotificationMod) {
	n.ApplyStyles(NotificationBail, mods...)
  n.Control.Play()
}

func (n *Notification) Fire(mods ...NotificationMod) {
	n.ApplyStyles(NotificationEmit, mods...)
  n.Control.Play()
}

func (n *Notification) Emit(mods ...NotificationMod) {
	n.ApplyStyles(NotificationEmit, mods...)
  Animator.App.QueueUpdateDraw(func() {
Animator.App.Draw(n.Box)
  })
  // n.Control.Play()
}

func (n *Notification) GetId() string {
	pageId := fmt.Sprintf("notif_%d", n.idx)
	return pageId
}

func (n *Notification) ApplyStyles(
	s NotificationState,
	mods ...NotificationMod,
) {
	for _, v := range append(n.mods, mods...) {
    if ok, name := v(s, n, n.m); ok {
      _ = name
      // zlog.Debug(name, zzlog.Bool("good", ok))
    }
	}
}

func (nm *NotificationManager) NewNotification(
	title, msg string,
	mods ...NotificationMod,
) *Notification {
	n := NewNotificationBox(title, msg, mods...)
  n.idx = len(nm.Notifications)
	n.m = nm
			Animator.App.QueueUpdateDraw(func() {
				n.SetRect(30, 3, 40, 4)
				Animator.App.Draw(n)
			})
	nm.Notifications = append(nm.Notifications, n)
	n.mods = append(n.mods, NotificationPlaced(n.GetId()))
	n.mods = append(n.mods, NotificationSize(nm.width, nm.height))
	n.mods = append(n.mods, NotificationInitialPos(nm.startX, nm.startY))
	cancel, notifTimeout := NotificationTimeout(time.Millisecond * 7000)
	n.mods = append(n.mods, notifTimeout)
	n.ApplyStyles(NotificationInit)
	anim := NewAnimation()
  anim.AutoNext = true
  anim.Loop = false
	kf := anim.NewKeyFrame()
	fm := kf.NewFrameMotion(
		NewMotion("MotionAxisY", MotionAxisY, *NewMotionMutator(RectYMutator)),
		NewBoxAnimated(n.Box),
	)
	fm.Tween = NewTween(float64(nm.startY), float64(nm.targetY))
	kf.AddMotion(fm)
	anim.AddKeyFrame(kf)
	n.mods = append(n.mods, NotificationFire(kf))
	bkf := anim.NewKeyFrame()
	bfm := bkf.NewFrameMotion(
		NewMotion("MotionAxisX", MotionAxisY, *NewMotionMutator(RectYMutator)),
		NewBoxAnimated(n.Box),
	)
	bfm.Tween = NewTween(float64(nm.targetY), float64(nm.startY))
	bkf.AddMotion(fm)
	anim.AddKeyFrame(bkf)
	n.mods = append(n.mods, NotificationBailout(bkf))
	n.mods = append(n.mods, NotificationKilled(n.GetId()))
	n.mods = append(
		n.mods,
		NotificationAnyCallback(
			NotificationDie,
			func(s NotificationState, n *Notification, nm *NotificationManager) (ok bool, desc string) {
      ok = false
      desc = "customcc"

				if !cancel.Stop() {
					<-cancel.C
          ok = true
				}
        return
			},
		),
	)

  n.Animation = anim
  anim.Start()
  n.Emit()

	return n
}

func (nm *NotificationManager) Name() string {
	return "notif"
}

func (nm *NotificationManager) HandleEvent(o events.ObservableEvent) bool {

	return true
}

