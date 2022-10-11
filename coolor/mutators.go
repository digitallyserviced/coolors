package coolor

import (
	// "fmt"

	"github.com/digitallyserviced/tview"
)



type MotionMutator struct {
	data    map[string]interface{}
	Mutator PropertyMutator
}


type Rect struct {
	x, y, width, height int
}

type RectAnimator func(r Rect, m Motion, rm RectMutator)
type PropertyMutator interface {
	Mutate(m MotionValues, i interface{}) bool
}

type PropertyMutatorFunction func(m MotionValues, i interface{}) bool

var (
	RectXMutator = NewBoxRectPropertyMutator(func(r Rect, m MotionValues) Rect {
		r.x = int(m.X)
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

func NewRect(x, y, width, height int) Rect {
	r := Rect{x, y, width, height}
	return r
}
type MutateCallbackType int
const (
  BeforeMutate MutateCallbackType = 1 << iota 
  AfterMutate
)
type NotifyBeforeMutate func(m MotionValues, i interface{}) bool
type NotifyAfterMutate func(m MotionValues, i interface{}) bool
type MutatorNotify interface {
  Notify(m MotionValues, i interface{}) bool
}

func (nbm NotifyBeforeMutate) Notify(m MotionValues, i interface{}) bool {
  return nbm(m, i)
}

func (nbm NotifyAfterMutate) Notify(m MotionValues, i interface{}) bool {
  return nbm(m, i)
}

func DrawAfterMutate(m MotionValues, i interface{}) bool {
		var b *tview.Box
		var ok bool
		b, ok = i.(*tview.Box)
		if !ok {
			return false
		}
  AppModel.app.QueueUpdateDraw(func() {
    AppModel.app.Draw(b)
  })
    // AppModel.scr.Show()
  return true
}

func (cm *CallbackMutator) Notify(t MutateCallbackType, m MotionValues, i interface{}) bool {
  if fn, ok := cm.Callbacks[t]; ok {
    return fn.Notify(m, i)
  }
  return true
}

type CallbackMutator struct {
	Mutator PropertyMutatorFunction
  Callbacks map[MutateCallbackType]MutatorNotify
}

type BoxMutator struct {
	Mutator PropertyMutatorFunction
}

type RectMutator func(r Rect, m MotionValues) Rect
type AnyMutator func(a interface{}, m MotionValues) Rect

func NewMotionMutator(mut PropertyMutator) (mm *MotionMutator) {
	mm = &MotionMutator{}
	mm.data = make(map[string]interface{})
	mm.Mutator = mut
	return
}
func NewCallbackMutator(pmf PropertyMutatorFunction) *CallbackMutator {
	bm := &CallbackMutator{
		Mutator: pmf,
	}

	return bm
}
func NewBoxMutator(pmf PropertyMutatorFunction) *CallbackMutator {
	bm := &CallbackMutator{
		Mutator: pmf,
    Callbacks: make(map[MutateCallbackType]MutatorNotify),
	}

  bm.AddCallback(AfterMutate, NotifyAfterMutate(DrawAfterMutate))

	return bm
}

func (bm *BoxMutator) Mutate(m MotionValues, i interface{}) bool {
  return bm.Mutator(m, i)
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
