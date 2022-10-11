package coolor

import (
	// "fmt"
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/charmbracelet/harmonica"
	"github.com/digitallyserviced/tview"
	// "github.com/digitallyserviced/coolors/coolor/util"
)

const (
	fps              = 30
	defaultFrequency = 7.0
	defaultDamping   = 0.6
)

var (
	wFrequency = 2.0
	wDamping   = 1.0
	bFrequency = 5.5
	bDamping   = 1.0
)

type (
	FrameMotionType     int
	AnimationStatusType int
)

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
	idx          int
	CurrentFrame *FrameMotions
	Frames       KeyFrames
	Control      *AnimControl
	*Animated
	*eventNotifier
}

type Animated struct {
	Box  *tview.Box
	Item interface{}
}

type AnimControl struct {
	start   chan struct{}
	cancel  chan struct{}
	done    chan struct{}
	updates map[MotionType]chan MotionValues
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
	name       string
	motionType FrameMotionType
}

type MotionValues struct {
	X, Xvelocity float64
}

type Motion struct {
	MotionType
	MotionMutator
}

type FrameMotion struct {
	idx         int
	started     time.Time
	startOffset float64
	Motion
	*Tween
	*Animated
	CurrentMotion MotionValues
	Motions       []*MotionValues
}

type FrameMotions struct {
	// Updates chan Motion
	Motions []*FrameMotion
}

type (
	KeyFrames          []*FrameMotions
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
	controller = &AnimControl{}
	controller.o = sync.Once{}
	controller.cancel = make(chan struct{})
	controller.start = make(chan struct{})
	controller.done = make(chan struct{})
	// controller.updates = make(map[MotionType]chan MotionValues)
	return
}

func NewFrameMotion(
	mt Motion,
	a *Animated,
) (fm *FrameMotion) {
	fm = &FrameMotion{
		idx:     0,
		started: time.Time{},
		// startOffset:   startOffset,
		Motion: mt,
		// Tween:         NewTween(startX, targetX),
		Animated: a,
		// CurrentMotion: MotionValues{X: startX, Xvelocity: 0},
		Motions: make([]*MotionValues, 0),
	}
	return fm
}
func MakeFrameMotion(
	startX, targetX float64,
	startOffset float64,
	mt Motion,
	a *Animated,
) (fm *FrameMotion) {
	fm = NewFrameMotion(mt, a)
	fm.Tween = NewTween(startX, targetX)
	fm.CurrentMotion = MotionValues{X: startX, Xvelocity: 0}
	fm.startOffset = startOffset

	// fm = &FrameMotion{
	// 	idx:           0,
	// 	started:       time.Time{},
	// 	Motion:        mt,
	// 	Tween:         ,
	// 	Animated:      a,
	// 	CurrentMotion: MotionValues{X: startX, Xvelocity: 0},
	// 	Motions:       []*MotionValues{},
	// }
	return fm
}

func NewTween(startX, targetX float64, args ...float64) *Tween {
	freq, damp := defaultFrequency, defaultDamping
	if len(args) == 2 {
		freq = args[0]
		damp = args[1]
	}
	tw := &Tween{
		Spring:  harmonica.Spring{},
		startX:  startX,
		targetX: targetX,
	}
	tw.Spring = harmonica.NewSpring(harmonica.FPS(fps), freq, damp)

	return tw
}

func NewAnimation() *Animation {
	a := &Animation{
		idx:           0,
		CurrentFrame:  nil,
		Frames:        make(KeyFrames, 0),
		Control:       &AnimControl{},
		Animated:      &Animated{},
		eventNotifier: NewEventNotifier("animation"),
	}
	return a
}

func NewMotion(
	name string,
	motionType FrameMotionType,
	mm MotionMutator,
) (m Motion) {
	m.MotionMutator = mm
	m.name = name
	m.motionType = motionType
	return
}

func NewMotionValues(x, xVelocity float64) (mv MotionValues) {
	mv.X = x
	mv.Xvelocity = xVelocity
	return
}

func (m *Motion) AddMutatee() {
}

func (ac *AnimControl) Close() {
	ac.o.Do(func() {
		close(ac.cancel)
		close(ac.done)
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
			if elapsed >= float64(offset * delta) {
				break
			}
		} else {
			if fm.Finished(mv) {
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

func (fm *FrameMotion) SetTween(
	startX, targetX float64,
	vel, damp float64,
) *FrameMotion {
	fm.Tween = NewTween(startX, targetX, vel, damp)
	fm.CurrentMotion = MotionValues{X: startX, Xvelocity: 0}
	return fm
}

func (fm *FrameMotion) GetReverse() *FrameMotion {
	revfm := MakeFrameMotion(fm.targetX, fm.startX, 1.0, fm.Motion, fm.Animated)
	return revfm
}

func (fm *FrameMotion) Start(controller *AnimControl) {
	tick := time.NewTicker(time.Second / fps)
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
			mv, done := fm.Update()
			xp.motionX.Set(math.Abs(mv.X))
			xp.motionVel.Set(math.Abs(mv.Xvelocity))
			if !fm.MotionMutator.Mutator.Mutate(mv, fm.Animated.Item) {
				return
			}
			if done {
				return
			}
			// xp.toC.Add(1)
			// select {
			// case <-controller.cancel:
			// 	return
			// case <-controller.done:
			// 	return
			// case updater <- m:
			// default:
			// }
		}
	}
}

func (fm *FrameMotion) Elapsed() (d time.Duration) {
	d = time.Since(fm.started)
	xp.elapsed.Set(int64(d))
	return
}

func (fm *FrameMotion) Updater(ac *AnimControl) {
	ac.AddFrameMotionUpdater(fm)
	defer fmt.Println("updater stopped")
	ac.l.Lock()
	updater := ac.updates[fm.MotionType]
	ac.l.Unlock()
	go func() {
		for {
			<-ac.start
			select {
			case <-ac.cancel:
				return
			case <-ac.done:
				return
			case mv := <-updater:
				xp.fromC.Add(1)
				xp.motionX.Set(math.Abs(mv.X))
				xp.motionVel.Set(math.Abs(mv.Xvelocity))
				if !fm.MotionMutator.Mutator.Mutate(mv, fm.Animated.Item) {
					return
				}
			}
		}
	}()
	go fm.Start(ac)
}

func (fm *FrameMotion) Reset() {
	fm.idx = 0
	fm.CurrentMotion.X = fm.startX
	fm.CurrentMotion.Xvelocity = 0
	fm.SetStartOffset(fm.startOffset)
}

func (fm *FrameMotion) Finished(mv MotionValues) bool {
	if fm.idx > 1 && math.Abs(fm.CurrentMotion.Xvelocity-mv.Xvelocity) < 0.001 &&
		math.Abs(fm.targetX-mv.X) < 1 {
		return true
	}
	return false
}
func (fm *FrameMotion) Update() (mv MotionValues, done bool) {
	if fm.idx == 0 {
		fm.started = time.Now()
	}
	mv = fm.CurrentMotion
	oldVel := mv.Xvelocity
	mv.X, mv.Xvelocity = fm.Tween.Spring.Update(mv.X, oldVel, fm.targetX)
	fm.Motions = append(fm.Motions, &mv)
	// if math.Abs(mv.xVelocity) < 0.1 {
	//   fmt.Println(mv.xVelocity, fm.CurrentMotion)
	// }
	if fm.Finished(mv) {
		log.Printf("%d %s %v", fm.idx, fm.Elapsed().String(), fm)
		done = true
	}

	fm.CurrentMotion = mv
	fm.idx++
	xp.motionIdx.Add(1)
	return
	//  if fm.idx > 1 && math.Abs(oldVel-mv.Xvelocity) < 0.001 &&
	// 	math.Abs(fm.targetX-mv.X) < 1 {
	// } else {
	// 	fm.idx++
	// 	done = false
	// }

}

func (fms *FrameMotions) Start() (controller *AnimControl) {
	defer fmt.Println("motions stopped stopped")
	controller = NewAnimControl()
	for _, f := range fms.Motions {
		go f.Start(controller)
	}
	return
}
func (fms *FrameMotions) GetRef() interface{} {
	return fms
}

func (fms *FrameMotions) AddMotion(frame *FrameMotion) *FrameMotions {
	fms.Motions = append(fms.Motions, frame)
	return fms
}

func (fms KeyFrames) AddMotions(motions *FrameMotions) KeyFrames {
	fms = append(fms, motions)
	// fms = append(fms, elems ...Type)
	return fms
}

func (a *Animation) GetCurrentFrame() *FrameMotions {
	a.CurrentFrame = a.Frames[a.idx]
	return a.CurrentFrame
}

func (a *Animation) GetFrame(idx int) *FrameMotions {
	idx = iclamp(idx, 0, len(a.Frames)-1)
	return a.Frames[idx]
}

func (a *Animation) GetFrameIdx() int {
	return a.idx
}

func (a *Animation) SetFrameIdx(idx int) *Animation {
	a.idx = iclamp(idx, 0, len(a.Frames)-1)
	a.CurrentFrame = a.Frames[a.idx]
	return a
}

func (a *Animation) GetRef() interface{} {
	return a
}

func (a *Animation) Prev() *Animation {
	a.SetFrameIdx(a.idx - 1)
	a.Notify(
		*a.NewObservableEvent(AnimationPrevious, "prev", a.GetCurrentFrame(), a),
	)
	return a
}

func (a *Animation) Next() *Animation {
	a.SetFrameIdx(a.idx + 1)
	a.Notify(*a.NewObservableEvent(AnimationNext, "next", a.GetCurrentFrame(), a))
	return a
}

func (a *Animation) Reset() *Animation {
	select {
	case <-a.Control.start:
		select {
		case <-a.Control.done:
			a.Control.Close()
			for _, v := range a.CurrentFrame.Motions {
				v.Reset()
			}
		}
	default:
	}
	return a
}

func (a *Animation) Start() *Animation {
	a.Control = a.GetCurrentFrame().Start()
	close(a.Control.start)
	a.Notify(
		*a.NewObservableEvent(AnimationStarted, "start", a.GetCurrentFrame(), a),
	)
	// control.wg.Wait()
	// control.Close()
	<-a.Control.done
	a.Notify(
		*a.NewObservableEvent(AnimationFinished, "done", a.GetCurrentFrame(), a),
	)
	// a.Next()
	// control = a.GetCurrentFrame().Start()
	// close(control.start)
	// <-control.done
	return a
}

func (a *Animation) AddKeyFrame(frame *FrameMotions) *Animation {
	a.Frames = append(a.Frames, frame)
	return a
}

func (a *Animation) MakeKeyFrame(frames ...*FrameMotion) *Animation {
	fms := &FrameMotions{
		Motions: make([]*FrameMotion, 0),
	}
	fms.Motions = append(fms.Motions, frames...)
	a.Frames = append(a.Frames, fms)
	return a
}
