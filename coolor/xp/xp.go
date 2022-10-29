package xp

import (
	"expvar"
	"runtime"
	"time"
)


type expvars struct {
	MotionX   *expvar.Float
	MotionVel *expvar.Float
	MotionIdx *expvar.Int
	Motions   *expvar.Var
	FromC     *expvar.Int
	Elapsed   *expvar.Int
	ToC       *expvar.Int
	Gr        *expvar.Int
}

var Xp *expvars

func SetupExpVars() {
	Xp = &expvars{
		MotionX:   expvar.NewFloat("MotionX"),
		MotionVel: expvar.NewFloat("MotionVel"),
		MotionIdx: expvar.NewInt("MotionIdx"),
		FromC:     expvar.NewInt("FromC"),
		ToC:       expvar.NewInt("ToC"),
		Gr:        expvar.NewInt("Goroutines"),
		Elapsed:   expvar.NewInt("Elapsed"),
		// gr := expvar.NewInt("Goroutines")
	}
	go func() {
		for range time.Tick(100 * time.Millisecond) {
			Xp.Gr.Set(int64(runtime.NumGoroutine()))
		}
	}()
}
