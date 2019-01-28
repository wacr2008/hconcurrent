package hconcurrent

import (
	"sync"
)

type concurrentItem struct {
	lock        *sync.Mutex
	inputChan   chan interface{}
	outputChan  chan interface{}
	doFuncCount int
	doFunc      func(interface{}) interface{}
	started     bool
}

func newConcurrentItem(
	inputChan chan interface{},
	doFuncCount int,
	doFunc func(interface{}) interface{},
	outputChan chan interface{},
) *concurrentItem {
	return &concurrentItem{
		lock:        new(sync.Mutex),
		inputChan:   inputChan,
		doFuncCount: doFuncCount,
		doFunc:      doFunc,
		outputChan:  outputChan,
	}
}

func (ci *concurrentItem) start() {
	ci.lock.Lock()
	if !ci.started {
		for i := 0; i < ci.doFuncCount; i++ {
			go ci.f()
		}
		ci.started = true
	}
	ci.lock.Unlock()
}

func (ci *concurrentItem) f() {
	for {
		v := <-ci.inputChan
		if v == nil {
			return
		}
		i := ci.doFunc(v)
		if i != nil && ci.outputChan != nil {
			ci.outputChan <- i
		}
	}
}

func (ci *concurrentItem) stop() {
	ci.lock.Lock()
	ci.stopNoLock()
	ci.lock.Unlock()
}

func (ci *concurrentItem) stopNoLock() {
	if !ci.started {
		return
	}
	for i := 0; i < ci.doFuncCount; i++ {
		ci.inputChan <- nil
	}
	ci.started = false
}
