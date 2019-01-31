package hconcurrent

import "sync"

type concurrentItem struct {
	inputChan      chan interface{}
	outputChan     chan interface{}
	stopChan       chan *sync.WaitGroup
	goroutineCount int
	handle         func(interface{}) interface{}
}

func newConcurrentItem(option Option) *concurrentItem {
	return &concurrentItem{
		inputChan:      make(chan interface{}, option.channelSize),
		stopChan:       make(chan *sync.WaitGroup, option.goroutineCount),
		goroutineCount: option.goroutineCount,
		handle:         option.handle,
	}
}

func (ci *concurrentItem) setOutputChan(c chan interface{}) {
	ci.outputChan = c
}

func (ci *concurrentItem) start() {
	for i := 0; i < ci.goroutineCount; i++ {
		go ci.f()
	}
}

func (ci *concurrentItem) f() {
	for {
		select {
		case v := <-ci.inputChan:
			r := ci.handle(v)
			if ci.outputChan != nil {
				ci.outputChan <- r
			}
		case w := <-ci.stopChan:
			w.Done()
			return
		}
	}
}

func (ci *concurrentItem) stop(w *sync.WaitGroup) {
	for i := 0; i < ci.goroutineCount; i++ {
		ci.stopChan <- w
	}
}
