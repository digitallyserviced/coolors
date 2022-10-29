package anim

import (
	// "fmt"
	"container/ring"
	"fmt"

	// "log"
	"math"
	"sync"
	"time"

	"github.com/charmbracelet/harmonica"
	// "github.com/digitallyserviced/tview"

	. "github.com/digitallyserviced/coolors/coolor/events"
	"github.com/digitallyserviced/coolors/coolor/util"
	"github.com/digitallyserviced/coolors/coolor/xp"
)

const (
	DefaultFPS              = 30
	DefaultFrequency = 7.0
	DefaultDamping   = 0.6
)

var (
	wFrequency = 2.0
	wDamping   = 1.0
	bFrequency = 5.5
	bDamping   = 1.0
)

type (
	FrameMotionType int
	// AnimationStatus int
)

// const (
//   AnimationPlaying AnimationStatus = 1 << iota
//   AnimationPaused
//   AnimationTransition
// )

const (
	MotionNil FrameMotionType = 1 << iota

	MotionAxisX
	MotionAxisY

	MotionSizeW
	MotionSizeH

	MotionColor

	MotionAxis = MotionAxisX | MotionAxisY
	MotionSize = MotionSizeH | MotionSizeW
)

type Animation struct {
	Frames  *KeyFrames
	Control *AnimControl
	*Animated
	*EventNotifier
	idx            int
	AutoNext, Loop bool
  Reanimate Reanimation
  State ObservableEventType
}

type Animated struct {
	Item interface{}
}

type AnimControl struct {
	start   chan struct{}
	cancel  chan struct{}
	paused  chan struct{}
	playing chan struct{}
	idle    chan struct{}
	done    chan struct{}
	ch      chan struct{}
	clch    chan struct{}
	nilch   chan struct{}
	updates map[MotionType]chan MotionValues
	master  time.Ticker
	l       sync.RWMutex
	wg      sync.WaitGroup
	o       sync.Once
}

type Tween struct {
	Spring  harmonica.Spring
	startX  float64
	targetX float64
}

type MotionType struct {
	Name       string
	MotionType FrameMotionType
}

type MotionValues struct {
	X, Xvelocity float64
}

type Motion struct {
	MotionType
	MotionMutator
}

type FrameMotion struct {
	started time.Time
	Motion
	*Animated
	*Tween
	Motions       []*MotionValues
	CurrentMotion MotionValues
	idx           int
	startOffset   float64
}

type FrameMotions struct {
	anim     *Animation
	Motions  []*FrameMotion
	IdleTime time.Duration
}

type KeyFrame struct {
	anim *Animation
	idx  int
	*ring.Ring
	Motions []*FrameMotion
}

type KeyFrames struct {
	*KeyFrame
	tail     *KeyFrame
	IdleTime time.Duration
}

type (
	KeyFrameMotions    []*FrameMotions
	AnimationStartable interface {
		Start()
	}
)

type AnimationNavigable interface {
	Next()
	Prev()
}

type AnimationControl interface {
	AnimationStartable
	AnimationNavigable
}

type AnimFrameHandler interface {
	AnimFrameStart(func())
	AnimFrameUpdate(frame *FrameMotion)
	AnimFrameEnd()
}
type MotionHandler interface {
	HandleMotionUpdate(fm *FrameMotion)
}

func NewAnimControl() (controller *AnimControl) {
	controller = &AnimControl{
		l:       sync.RWMutex{},
		wg:      sync.WaitGroup{},
		o:       sync.Once{},
	}
	controller.o = sync.Once{}
	controller.cancel = make(chan struct{})
	controller.ch = make(chan struct{})
	controller.clch = make(chan struct{})
	close(controller.clch)
	controller.nilch = nil
	controller.paused = controller.nilch
	controller.playing = controller.nilch
	controller.idle = controller.ch
	controller.start = make(chan struct{})
	controller.done = make(chan struct{})
	// controller.updates = make(map[MotionType]chan MotionValues)
	return
}

func NewTween(startX, targetX float64, args ...float64) *Tween {
	freq, damp := DefaultFrequency, DefaultDamping
	if len(args) == 2 {
		freq = args[0]
		damp = args[1]
	}
	tw := &Tween{
		Spring:  harmonica.Spring{},
		startX:  startX,
		targetX: targetX,
	}
	tw.Spring = harmonica.NewSpring(harmonica.FPS(DefaultFPS), freq, damp)
  // zlog.Debug("new tween", zzlog.Reflect("tween", tw))

	return tw
}

func NewAnimation() *Animation {
	a := &Animation{
		idx:      0,
		AutoNext: true,
		State:   AnimationInit,
		EventNotifier: NewEventNotifier("animation"),
	}
	a.Frames = a.NewKeyFrames()
	return a
}

func NewMotion(
	name string,
	motionType FrameMotionType,
	mm MotionMutator,
) (m Motion) {
	m.MotionMutator = mm
	m.Name = name
	m.MotionType.MotionType = motionType
	return
}

func NewMotionValues(x, xVelocity float64) (mv MotionValues) {
	mv.X = x
	mv.Xvelocity = xVelocity
	return
}

// func (m *Tween) SetTarget(targetX float64) {
//
// }
// func (m *Motion) AddMutatee() {
// }

func (tw *Tween) GetReverse() *Tween {
	ntw := &Tween{
		Spring: harmonica.NewSpring(
			harmonica.FPS(DefaultFPS),
			DefaultFrequency,
			DefaultDamping,
		),
		startX:  tw.targetX,
		targetX: tw.startX,
	}
	return ntw
}

func (fm *FrameMotion) Evolve(delta, offset float64, full bool) MotionValues {
	// delta := harmonica.FPS(fps)
	elapsed := 0.0
	mv := MotionValues{fm.CurrentMotion.X, fm.CurrentMotion.Xvelocity}
	motions := make([]MotionValues, 0)
	for {
		mv.X, mv.Xvelocity = fm.Spring.Update(
			mv.X,
			mv.Xvelocity,
			fm.targetX,
		)
		// fm.Motions = append(fm.Motions, &mv)
		motions = append(motions, mv)
		if !full {
			elapsed += delta
			if elapsed >= float64(offset*delta) {
				break
			}
		} else {
			switch fm.Transition(mv) {
			case AnimationPlaying:
			case AnimationPaused:
				fallthrough
			case AnimationIdle:
				break
			}
		}
	}
	return mv
}

// bar = "-" if ascii else "━"
//        half_bar_right = " " if ascii else "╸"
//        half_bar_left = " " if ascii else "╺"

func (fm *FrameMotion) SetStartOffset(
	offset float64,
) *FrameMotion {
	fm.startOffset = offset
	if fm.Tween != nil {
		fm.CurrentMotion = fm.Evolve(harmonica.FPS(30), fm.startOffset, false)
	}
	return fm
}

func (fm *FrameMotion) UpdateTween(
	startX, targetX float64,
) bool {
	cont := false
	if fm.Tween.targetX != targetX || fm.Tween.startX != startX {
		cont = true
	}
	fm.Tween.targetX = targetX
	fm.Tween.startX = startX
	fm.CurrentMotion.X = startX
	return cont
}

func (fm *FrameMotion) SetTween(
	startX, targetX float64,
	vel, damp float64,
) *FrameMotion {
	fm.Tween = NewTween(startX, targetX, vel, damp)
	fm.CurrentMotion = MotionValues{X: startX, Xvelocity: 0}
	return fm
}


func (fm *FrameMotion) Elapsed() (d time.Duration) {
	d = time.Since(fm.started)
	xp.Xp.Elapsed.Set(int64(d))
	return
}

func (fm *FrameMotion) Reset() {
	fm.idx = 0
	fm.CurrentMotion.X = fm.startX
	fm.CurrentMotion.Xvelocity = 0
	fm.SetStartOffset(fm.startOffset)
}

func (fm *FrameMotion) Update() (mv MotionValues, done bool) {
	if fm.idx == 0 {
		fm.started = time.Now()
	}
	mv = fm.CurrentMotion
	oldVel := mv.Xvelocity
	mv.X, mv.Xvelocity = fm.Tween.Spring.Update(mv.X, oldVel, fm.targetX)
	// fm.Motions = append(fm.Motions, &mv)
	switch fm.Transition(mv) {
    case AnimationPlaying:
    case AnimationPaused:
      fallthrough
    case AnimationIdle:
      done = true
	}

	fm.CurrentMotion = mv
	fm.idx++
	xp.Xp.MotionIdx.Add(1)
	return
}

func (fm *FrameMotion) Transition(mv MotionValues) ObservableEventType {
	return fm.MotionMutator.Mutator.Finished(
		fm.idx,
		fm.startX,
		fm.targetX,
		fm.CurrentMotion,
		mv,
		fm.Animated.Item,
	)
}

func (fm *FrameMotion) Tick(ac *AnimControl) {
	mv, done := fm.Update()
	xp.Xp.MotionX.Set(math.Abs(mv.X))
	xp.Xp.MotionVel.Set(math.Abs(mv.Xvelocity))
	if !fm.MotionMutator.Mutator.Mutate(mv, fm.Animated.Item) {
		ac.Pause()
	}
	if done {
		ac.Idle()
	}
}

func (fm *FrameMotion) Start(controller *AnimControl) {
	tick := time.NewTicker(time.Second / DefaultFPS)
	defer tick.Stop()
	defer controller.Close()
	defer fmt.Println("frame motion stopped stopped")
	<-controller.start
	for {
		select {
		case <-controller.cancel:
			return
		case <-controller.done:
			return
		case <-tick.C:
			// if !fm.Tick(controller) {
			// 	return
			// }
		}
	}
}

// func (fms *FrameMotions) GetRef() interface{} {
// 	return fms
// }
//
//
// func (fms *FrameMotions) Start() (controller *AnimControl) {
// 	defer fmt.Println("motions stopped stopped")
// 	controller = NewAnimControl()
// 	for _, f := range fms.Motions {
// 		go f.Start(controller)
// 	}
// 	return
// }

func (fms *KeyFrame) Jump(){
  fms.anim.Frames.KeyFrame = fms
  fms.anim.Control.Play()
}

func (fms *KeyFrame) NewFrameMotion(
	mt Motion,
	a *Animated,
) (fm *FrameMotion) {
	fm = &FrameMotion{
		idx:      len(fms.Motions),
		started:  time.Time{},
		Motion:   mt,
		Animated: a,
		Motions:  make([]*MotionValues, 0),
	}
	return fm
}

func (fms *KeyFrame) MakeFrameMotion(
	startX, targetX float64,
	startOffset float64,
	mt Motion,
	a *Animated,
) (fm *FrameMotion) {
	fm = fms.NewFrameMotion(mt, a)
	fm.Tween = NewTween(startX, targetX)
	fm.CurrentMotion = MotionValues{X: startX, Xvelocity: 0}
	fm.startOffset = startOffset
	return fm
}

func (fms *KeyFrame) AddMotion(frame *FrameMotion) *KeyFrame {
	fms.Motions = append(fms.Motions, frame)
	return fms
}
func (a *KeyFrame) GetRef() interface{} {
  return a
}
func (a *KeyFrames) GetRef() interface{} {
	return a
}

func (kfs *KeyFrames) Last() *KeyFrame {
	n := kfs.tail.Prev()
	if n.Value == nil {
		return kfs.tail
	}
	if kf, ok := n.Value.(*KeyFrame); ok {
		return kf
	}
	return nil
}
func (kfs *KeyFrames) First() *KeyFrame {
	n := kfs.tail.Next()
	if n.Value == nil {
		return kfs.tail
	}
	if kf, ok := n.Value.(*KeyFrame); ok {
		return kf
	}
	return nil
}
func (kfs *KeyFrames) End() bool {
	n := kfs.KeyFrame.Next()
	// zlog.Debug(
	// 	"endnex",
	// 	zzlog.Reflect("nextframe", n),
	// 	zzlog.Reflect("thisframe", kfs.KeyFrame),
	// )
	if n.Value == nil {
		if kfs.anim.Loop {
			n = kfs.tail.Next()
      kfs.anim.UpdateState(AnimationLooped)
		} else {
      kfs.anim.UpdateState(AnimationFinished)
    }
	}
	return n.Value == nil
}

// func NewKeyFrame(motions... *FrameMotion) *KeyFrame {
//   r := ring.New(1)
//   kf := &KeyFrame{
//   	idx:     0,
//   	Ring:    r,
//   	Motions: motions,
//   }
//   r.Value = kf
//   return kf
// }

func (kfs *KeyFrames) Add(kf *KeyFrame) *KeyFrames {
	kf.idx = kfs.Len() - 1
	o := kfs.tail.Prev().Link(kf.Ring)
	if o.Value == nil {
		o = kf.Ring
	}
	kfs.tail.idx = kfs.Len()
	if kfs.KeyFrame.Value == nil {
		kfs.KeyFrame = kf
	}
	return kfs
}
func (ac *AnimControl) Close() {
  drainAndClose := func(ch... chan struct {}){
    for _, v := range ch {
      close(v)
      for range v {}
    }
  }
	ac.o.Do(func() {
    drainAndClose(ac.cancel, ac.done, ac.ch)
		// time.AfterFunc(10*time.Millisecond, func() {
		// 	for x, c := range ac.updates {
		// 		go func(x MotionType, c chan MotionValues) {
		// 			close(c)
		// 			// delete(ac.updates, x)
		// 			go func() {
		// 				for range c {
		// 				}
		// 			}()
		// 		}(x, c)
		// 	}
		// })
	})
}

func (ac *AnimControl) AddFrameMotionUpdater(fm *FrameMotion) {
	ac.l.Lock()
	ac.updates[fm.MotionType] = make(chan MotionValues)
	ac.l.Unlock()
}
func (ac *AnimControl) Idle() {
	ac.l.Lock()
	defer ac.l.Unlock()
	ac.idle = ac.clch
	ac.playing = ac.nilch
	ac.paused = ac.nilch
}
func (ac *AnimControl) Pause() {
	ac.l.Lock()
	defer ac.l.Unlock()
	ac.paused = ac.clch
	ac.playing = ac.nilch
	ac.idle = ac.nilch
}
func (ac *AnimControl) Play() {
	ac.l.Lock()
	defer ac.l.Unlock()
	ac.playing = ac.clch
	ac.paused = ac.nilch
	ac.idle = ac.nilch
}

func (a *Animation) UpdateState(s ObservableEventType) bool {
  if !a.State.Is(s) {
    a.State = s
    a.Notify(*a.NewObservableEvent(s, s.String(), a, a.Frames))
    return true
  }
  return false
}

func (a *Animation) Animate() {
	frames := a.GetCurrentFrame()
	if frames == nil || frames.Motions == nil || len(frames.Motions) == 0 {
		// zlog.Error("No frames!", zzlog.String("anim", a.notifierName))
		// fmt.Println("No motions")
		return
	}

  longtick := time.Millisecond * 1000

	tick := time.NewTicker(time.Second / DefaultFPS)
  ltick := time.NewTicker(longtick)
  ltick.Stop()
	for {
		select {
		case <-a.Control.paused:
    a.UpdateState(AnimationPaused)
			select {
      case <-ltick.C:
      // fmt.Println("ltick")
			case <-a.Control.playing:
    a.UpdateState(AnimationPlaying)
						frames = a.Next()
					// if a.Frames.End() {
					// 	// a.Control.Play()
					// 	continue
					// }
			case <-a.Control.idle:
frames = a.GetCurrentFrame()
    a.UpdateState(AnimationIdle)
			case <-a.Control.cancel:
    a.UpdateState(AnimationCanceled)
				return
			case <-a.Control.done:
    a.UpdateState(AnimationDone)
				return
			}
		case <-a.Control.playing:
			select {
			case <-tick.C:
				// zlog.Debug("playing", zzlog.String("anim", a.notifierName))
				for _, fm := range frames.Motions {
					go fm.Tick(a.Control)
				}
			}
		case <-a.Control.idle:
    a.UpdateState(AnimationIdle)
			if a.Frames.IdleTime == time.Duration(0) {
				a.Frames.IdleTime = time.Millisecond * 500
			}
			timer := time.NewTimer(a.Frames.IdleTime)
    ltick.Reset(longtick)
			select {
			case <-timer.C:
				// zlog.Debug("idle finish", zzlog.String("anim", a.notifierName))
				timer.Stop()
				if a.AutoNext {
					if !a.Frames.End() {
						frames = a.Next()
						a.Control.Play()
						continue
					}
				}
				a.Control.Pause()
			case <-a.Control.cancel:
    a.UpdateState(AnimationCanceled)
				return
			case <-a.Control.done:
    a.UpdateState(AnimationDone)
				return
			}
      case <-ltick.C:
      // fmt.Println("ltick")
        // default:
		}
	}
}

func (a *Animation) Start() *Animation {
	if a.Control != nil {
		if a.Control.playing == nil {
			a.Control.Play()
		}
		return a
	}
	a.Control = NewAnimControl()
	a.Control.Pause()
	go a.Animate()
	return a
}

func (a *Animation) GetCurrentFrame() *KeyFrame {
	if fms, ok := a.Frames.Value.(*KeyFrame); ok {
		return fms
	} else {
		return nil
	}
}

func (a *Animation) GetFrame(idx int) *KeyFrame {
	idx = util.Clamp(idx, 0, a.Frames.Len()-1) + 1
	if fms, ok := a.Frames.tail.Move(idx).Value.(*KeyFrame); ok {
		return fms
	} else {
		return nil
	}
}

func (a *Animation) GetFrameIdx() int {
	return a.Frames.idx
}

// func (a *Animation) SetFrameIdx(idx int) *Animation {
// 	idx = iclamp(idx, 0, a.Frames.Len()-1) + 1
//   if fms, ok := a.Frames.tail.Move(idx).Value.(*KeyFrame); ok {
//     a.Frames = fms
//   } else {
//     return nil
//   }
// 	// a.idx = iclamp(idx, 0, len(a.Frames)-1)
// 	// a.CurrentFrame = a.Frames[a.idx]
// 	return a
// }

func (a *Animation) GetRef() interface{} {
	return a
}
func (a *Animation) Next() *KeyFrame {
	n := a.Frames.Next()
	if n.Value == nil {
		a.Control.Pause()
		n = n.Next()
	}
	if kf, ok := n.Value.(*KeyFrame); ok && kf.Value != nil {
		a.Frames.KeyFrame = kf
    a.UpdateState(AnimationNext)
		return kf
	}
	return nil
}

func (a *Animation) Prev() *KeyFrame {
	n := a.Frames.Prev()
	if n.Value == nil {
		n = n.Prev()
	}
	if kf, ok := n.Value.(*KeyFrame); ok && kf.Value != nil {
		a.Frames.KeyFrame = kf
    a.UpdateState(AnimationPrevious)
		return kf
	}
	return nil
}
func (a *Animation) NewKeyFrames(frames ...*KeyFrame) *KeyFrames {
	// if len(frames) < 1 {
	//   return nil
	// }
	r := ring.New(1)
	r.Value = nil

	tail := a.NewKeyFrame()
	tail.Ring.Value = nil
	kfs := &KeyFrames{
		// KeyFrame: &KeyFrame{},
		tail:     tail,
		IdleTime: time.Millisecond * 1000,
	}
	kfs.KeyFrame = kfs.tail

	for _, kf := range frames {
		kfs.Add(kf)
	}

	if kfs.KeyFrame.Value == nil && kfs.KeyFrame.Len() > 1 {
		if kf, ok := kfs.tail.Next().Value.(*KeyFrame); ok && kf.Value != nil {
			kfs.KeyFrame = kf
		}
	}

	return kfs
}

func (a *Animation) AddKeyFrame(frame *KeyFrame) *KeyFrame {
	a.Frames.Add(frame)
	return frame
}

func (a *Animation) NewKeyFrame(motions ...*FrameMotion) *KeyFrame {
	r := ring.New(1)
	kf := &KeyFrame{
		anim:    a,
		idx:     0,
		Ring:    r,
		Motions: motions,
	}
	r.Value = kf
	return kf
	//  kf := NewKeyFrame(motions...)
	//  kf.anim = a
	// return kf
}

func (a *Animation) MakeKeyFrame(motions ...*FrameMotion) *KeyFrame {
	kf := a.NewKeyFrame(motions...)
	a.AddKeyFrame(kf)
	return kf
}
