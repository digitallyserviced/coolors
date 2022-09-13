package coolor

import (
	"sync"
	"time"
)

func takeN(
	done <-chan struct{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func takeFn(
	done <-chan struct{},
	csp *ColorStreamProgress,
	valueStream <-chan interface{},
	fn func(interface{}) bool,
) <-chan interface{} {
	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for {
			select {
			case <-done:
				return
			case v := <-valueStream:
				csp.Itrd()
				if v == nil {
					continue
				}
				if fn(v) {
					csp.Validd()
					takeStream <- v
				}
			}
		}
	}()
	return takeStream
}

func asStream(
	done <-chan struct{},
	fn func() interface{},
	throttle time.Duration,
) <-chan interface{} {
	s := make(chan interface{})
	tttl := time.NewTicker(throttle)
	go func() {
		defer close(s)

		for {
			select {
			case <-done:
				return
			case s <- fn():
				<-tttl.C
			}
		}
	}()
	return s
}

func fanIn(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		var wg sync.WaitGroup
		wg.Add(len(chans))

		for _, c := range chans {
			go func(c <-chan interface{}) {
				for v := range c {
					out <- v
				}
				wg.Done()
			}(c)
		}

		wg.Wait()
		close(out)
	}()
	return out
}

func TakeNColors(
	done <-chan struct{},
	valueStream <-chan interface{},
	num int,
) []Color {
	colors := make([]Color, 0)
	for cv := range takeN(done, valueStream, num) {
		if cv == nil {
			continue
		}
		col := cv.(Color)
		colors = append(colors, col)
	}
	return colors
}
