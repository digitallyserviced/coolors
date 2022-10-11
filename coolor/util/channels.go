package util

import (
	"sync"
	"time"
)



func TakeN(
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

func LossySend[T any](
  done <-chan struct{},
  valueChan chan<- T,
  value T,
  t time.Duration,
){
  select {
  case <-done:
    return
  case <-time.After(t):
    return
  case valueChan <- value:
}

}

func TakeFn(
	done <-chan struct{},
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
				if v == nil {
					continue
				}
				if fn(v) {
					takeStream <- v
				}
			}
		}
	}()
	return takeStream
}

func AsStream(
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

func FanIn(chans ...<-chan interface{}) <-chan interface{} {
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
