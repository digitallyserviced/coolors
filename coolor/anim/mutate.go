package anim

import (
	"fmt"
	"math"

	"github.com/digitallyserviced/tview"

	// "github.com/digitallyserviced/coolors/coolor/zzlog"
	. "github.com/digitallyserviced/coolors/coolor/events"
)



type PropertyFinishedMutating interface {
	Finished(idx int, startX, targetX float64, prevMv, newMv MotionValues, i interface{}) ObservableEventType
}
type PropertyMutator interface {
	PropertyFinishedMutating
	Mutate(m MotionValues, i interface{}) bool
}
type MutatorNotify interface {
	Notify(m MotionValues, i interface{}) bool
}

type RectAnimator func(r Rect, m Motion, rm RectMutator)
type PropertyMutatorFunction func(m MotionValues, i interface{}) bool
type NotifyBeforeMutate func(m MotionValues, i interface{}) bool
type NotifyAfterMutate func(m MotionValues, i interface{}) bool
type NotifyFinishedMutate func(m MotionValues, i interface{}) bool
type RectMutator func(r Rect, m MotionValues) Rect
type AnyMutator func(a interface{}, m MotionValues) Rect
type MutatorFinishedCallback func(idx int, startX float64, targetX float64, prevMv MotionValues, newMv MotionValues, i interface{}) ObservableEventType

type MotionMutator struct {
	data    map[string]interface{}
	Mutator PropertyMutator
}
type Rect struct {
	x, y, width, height int
}
type CallbackMutator struct {
	Mutator   PropertyMutatorFunction
	Callbacks map[MutateCallbackType]MutatorNotify
  FinishedCallback PropertyFinishedMutating
}
type BoxMutator struct {
	Mutator PropertyMutatorFunction
}
type MutateCallbackType int
const (
	BeforeMutate MutateCallbackType = 1 << iota
	AfterMutate
  FinishedMutating
)

var (
	RectYMutator = NewBoxRectPropertyMutator(func(r Rect, m MotionValues) Rect {
		r.y = int(m.X)
		return r
	})
	RectXMutator = NewBoxRectPropertyMutator(func(r Rect, m MotionValues) Rect {
		r.x = int(m.X)
		return r
	})
	RectHMutator = NewBoxRectPropertyMutator(func(r Rect, m MotionValues) Rect {
		r.height = int(m.X)
		return r
	})
	RectWMutator = NewBoxRectPropertyMutator(func(r Rect, m MotionValues) Rect {
		r.width = int(m.X)
		return r
	})
)

func NewPropertyMutator(rMut AnyMutator) PropertyMutator {
	return NewBoxMutator(func(m MotionValues, i interface{}) bool {
		rMut(i, m)
		return true
	})
}
func NewBoxRectPropertyMutator(rMut RectMutator) PropertyMutator {
	return NewBoxMutator(func(m MotionValues, i interface{}) bool {
		var b *tview.Box
		var ok bool
		b, ok = i.(*tview.Box)
		if !ok {
			return false
		}
		r := NewRect(b.GetRect())
		r = rMut(r, m)
		b.SetRect(r.x, r.y, r.width, r.height)
		return true
	})
}

func NewMotionMutator(mut PropertyMutator) (mm *MotionMutator) {
	mm = &MotionMutator{}
	mm.data = make(map[string]interface{})
	mm.Mutator = mut
	return
}
func NewCallbackMutator(pmf PropertyMutatorFunction) *CallbackMutator {
	bm := &CallbackMutator{
		Mutator: pmf,
		Callbacks: make(map[MutateCallbackType]MutatorNotify),
	}

	return bm
}
func NewBoxMutator(pmf PropertyMutatorFunction) *CallbackMutator {
	bm := &CallbackMutator{
		Mutator:   pmf,
		Callbacks: make(map[MutateCallbackType]MutatorNotify),
	}

	bm.AddCallback(AfterMutate, NotifyAfterMutate(DrawAfterMutate))

	return bm
}
func NewRect(x, y, width, height int) Rect {
	r := Rect{x, y, width, height}
	return r
}
func DrawAfterMutate(m MotionValues, i interface{}) bool {
	var b *tview.Box
	var ok bool
	b, ok = i.(*tview.Box)
	if !ok {
		return false
	}
  fmt.Println(b)
	return true
}
func (bm *BoxMutator) Mutate(m MotionValues, i interface{}) bool {
	return bm.Mutator(m, i)
}
func (nbm NotifyFinishedMutate) Notify(m MotionValues, i interface{}) bool {
	return nbm(m, i)
}
func (nbm NotifyBeforeMutate) Notify(m MotionValues, i interface{}) bool {
	return nbm(m, i)
}

func (nbm NotifyAfterMutate) Notify(m MotionValues, i interface{}) bool {
	return nbm(m, i)
}


func (mfc MutatorFinishedCallback) Finished(idx int, startX float64, targetX float64, prevMv MotionValues, newMv MotionValues, i interface{}) ObservableEventType {
  return mfc(idx, startX, targetX, prevMv, newMv, i)
}

// Finished implements PropertyMutator
func (cm *CallbackMutator) Notify(t MutateCallbackType, m MotionValues, i interface{}) bool {
	if fn, ok := cm.Callbacks[t]; ok {
		return fn.Notify(m, i)
	}
	return true
}
func (cm *CallbackMutator) Finished(idx int, startX float64, targetX float64, prevMv MotionValues, newMv MotionValues, i interface{}) ObservableEventType {
  fini := AnimationPlaying
  if cm.FinishedCallback != nil {
    fini = cm.FinishedCallback.Finished(idx, startX, targetX, prevMv, newMv, i)
  } else {
    if idx > 1 && math.Abs(prevMv.Xvelocity - newMv.Xvelocity) < 0.001 && math.Abs(targetX - newMv.X) < 1 {
      fini = AnimationIdle
    }
  }
  if fini.Is(AnimationFinished) || fini.Is(AnimationIdle) || fini.Is(AnimationPaused) {
    cm.Notify(FinishedMutating, newMv, i)
  }
  return fini
}

func (cm *CallbackMutator) AddCallback(t MutateCallbackType, n MutatorNotify) {
	cm.Callbacks[t] = n
}
func (cm *CallbackMutator) Mutate(m MotionValues, i interface{}) (res bool) {
	if res = cm.Notify(BeforeMutate, m, i); !res {
		return
	}
	if res = cm.Mutator(m, i); !res {
		return
	}
	if res = cm.Notify(AfterMutate, m, i); !res {
		return
	}
	return true
}
