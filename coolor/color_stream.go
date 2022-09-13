package coolor

import (
	"context"
	"sync/atomic"
	"time"
)

type ColorCount uint32

type ColorStreamProgress struct {
	ProgressHandler ColorStreamProgressHandler
	Valid           uint32
	Itr             uint32
}
type ColorStream struct {
	OutColors <-chan interface{}
	Start     chan struct{}
	Status    *ColorStreamProgress
	Cancel    context.CancelFunc
	Generator func() interface{}
	Validator func(interface{}) bool
	Context   context.Context
}

type FunctionalProgressHandler struct {
	v func(uint32)
	i func(uint32)
}

type ColorStreamIterationProgressHandler interface {
	OnItr(uint32)
}
type ColorStreamValidProgressHandler interface {
	OnValid(uint32)
}
type ColorStreamProgressHandler interface {
	ColorStreamIterationProgressHandler
	ColorStreamValidProgressHandler
}
type NilProgressHandler struct{}

func (NilProgressHandler) OnItr(i uint32)   {}
func (NilProgressHandler) OnValid(v uint32) {}

// func init() {
// 	rand.Seed(time.Now().UnixMilli())
// }
func checkColorDistance(tcol, tcol2 Color, distance float64) bool {
	return tcol2.DistanceCIEDE2000(tcol) <= distance
}

func NewNilProgressHandler() ColorStreamProgressHandler {
	nph := &NilProgressHandler{}
	return nph
}

func (fph FunctionalProgressHandler) OnItr(i uint32) {
	fph.i(i)
}

func (fph FunctionalProgressHandler) OnValid(v uint32) {
	fph.v(v)
}

func NewProgressHandler(
	v func(uint32),
	i func(uint32),
) *FunctionalProgressHandler {
	nph := &FunctionalProgressHandler{
		v: v,
		i: i,
	}
	return nph
}

func NewColorStreamProgress() *ColorStreamProgress {
	csp := &ColorStreamProgress{
		Valid:           0,
		Itr:             0,
		ProgressHandler: NewNilProgressHandler(),
	}
	return csp
}

func (csp *ColorStreamProgress) SetProgressHandler(
	csph ColorStreamProgressHandler,
) {
	csp.ProgressHandler = csph
}

func (csp *ColorStreamProgress) GetValid() uint32 {
	return atomic.LoadUint32(&csp.Valid)
}

func (csp *ColorStreamProgress) GetItr() uint32 {
	return atomic.LoadUint32(&csp.Itr)
}

func (csp *ColorStreamProgress) Itrd() {
	res := atomic.AddUint32(&csp.Itr, 1)
	csp.ProgressHandler.OnItr(res)
}

func (csp *ColorStreamProgress) Validd() {
	res := atomic.AddUint32(&csp.Itr, 1)
	csp.ProgressHandler.OnValid(res)
}
func (cs *ColorStream) Run(done <-chan struct{}) {
	numRoutines := 4
	generators := make([]<-chan interface{}, 0)
	for i := 0; i < numRoutines; i++ {
		generators = append(
			generators,
			takeFn(
				done,
				cs.Status,
				asStream(done, genRandomSeededColor, time.Millisecond*10),
				cs.Validator,
			),
		)
	}
	cs.OutColors = fanIn(generators...)
}

func StartColorStream(
	g func() interface{},
	v func(interface{}) bool,
) *ColorStream {
	ctx, cancel := context.WithCancel(context.Background())
	// ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	cs := &ColorStream{
		OutColors: make(<-chan interface{}),
		Start:     make(chan struct{}),
		Status:    NewColorStreamProgress(),
		Cancel:    cancel,
		Generator: func() interface{} {
			return &Color{0, 0, 0}
		},
		Validator: func(interface{}) bool {
			return true
		},
		Context: ctx,
	}

	if g != nil {
		cs.Generator = g
	}

	if v != nil {
		cs.Validator = v
	}

	return cs
}

func RandomHuesStream(
	tcol Color,
	maxDistance float64,
) *ColorStream { // , cs *ColorStream
	cs := StartColorStream(genRandomSeededColor, func(c interface{}) bool {
		return checkColorDistance(tcol, c.(Color), maxDistance)
	})
	return cs
}

func RandomShadesStream(
	tcol Color,
	maxDistance float64,
) *ColorStream { // , cs *ColorStream
	cs := StartColorStream(genRandomSeededColor, func(c interface{}) bool {
		return checkColorDistance(tcol, c.(Color), maxDistance)
	})
	return cs
}

func ShutdownStream(cs *ColorStream) {
	cs.Cancel()
}
