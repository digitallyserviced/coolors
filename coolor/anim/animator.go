package anim

import (
	// "fmt"

	// "fmt"

	"fmt"
	"time"

	"github.com/digitallyserviced/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/gookit/goutil/errorx"

	. "github.com/digitallyserviced/coolors/coolor/events"
)

type animator struct {
	Animations map[string]*Animation
	handlers   map[string]AnimatorHandler
	App        *tview.Application
	RootLayer  *tview.Pages
	AnimLayer  *tview.Pages
	Screen     tcell.Screen
	LayerName  string
	*EventObserver
}
type Creanimation func(name string) *Animation
type Reanimation func(anim *Animation) bool

var Animator *animator
var NotifMgr *NotificationManager

func Init(
	app *tview.Application,
	scr tcell.Screen,
	viewRoot, viewAnim *tview.Pages,
) *animator {
	if Animator == nil {
		Animator = &animator{
			Animations:    make(map[string]*Animation, 0),
			handlers:      make(map[string]AnimatorHandler),
			App:           app,
			RootLayer:     viewRoot,
			AnimLayer:     viewAnim,
			Screen:        scr,
			LayerName:     "anims",
			EventObserver: NewEventObserver("animator"),
		}
		go Animator.watcher()
	}
	return Animator
}

func (a *animator) watcher() {
	tick := time.NewTicker(time.Millisecond * 200)
	for range tick.C {
		for _, anim := range a.Animations {
			if anim == nil {
				continue
			}
			if ObservableEventType(
				AnimationPlaying | AnimationIdle | AnimationLooped | AnimationNext | AnimationPrevious | AnimationUpdate,
			).Is(anim.State) {
        a.App.Draw()
				// a.App.QueueUpdateDraw(func() {
				// 	if anim == nil {
				// 		return
				// 	}
				// 	// fmt.Printf("animating %s : %s \n", n, anim.State.String())
				// })
			}
		}
	}
}

func GetAnimator() *animator {
	if Animator == nil {
		panic(errorx.New("Not initialized with root and anim page layers"))
	}
	return Animator
}

func (anim *animator) GetAnimation(name string) (a *Animation) {
	if aa, ok := Animator.Animations[name]; ok && aa != nil {
		return aa
	}
	return nil
}

func (anim *animator) HandleEvent(ev ObservableEvent) bool {
	fmt.Println(fmt.Sprintf("%+v", ev))
	if ev.Type.Is(AnimationFinished) || ev.Type.Is(AnimationDone) {
		if a, ok := ev.Ref.(*Animation); ok {
			if h, ok := anim.handlers[a.NotifierName]; ok {
				if !h.Handle(a) {
					anim.Animations[a.NotifierName] = nil
					delete(anim.Animations, a.NotifierName)
				}
			}
		}
	}
	return true
}

type AnimatorHandlerFunc func(a *Animation) bool

func (animFunc AnimatorHandlerFunc) Handle(a *Animation) bool {
	return animFunc(a)
}

func NewAnimatorHandler(f AnimatorHandlerFunc) AnimatorHandler {
	return f
}

type AnimatorHandler interface {
	Handle(*Animation) bool
}

func (anim *animator) NewAnimation(
	name string,
	f Creanimation,
	handlers ...AnimatorHandlerFunc,
) (aa *Animation) {
	aa = anim.GetAnimation(name)
	if aa != nil {
		aa.NotifierName = name
		return aa
	}
	aa = f(name)
	aa.NotifierName = name
	anim.Animations[name] = aa
	aa.Register(AllEvents, anim)
	return
}

func NewBoxAnimated(b tview.Primitive) *Animated {
	return &Animated{
		Item: b,
	}
}
