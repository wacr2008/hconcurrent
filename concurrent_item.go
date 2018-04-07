package hconcurrent

import "sync"

type concurrentItem struct {
	lock        *sync.Mutex
	wait        *sync.WaitGroup
	inputChan   chan interface{}
	outputChan  chan interface{}
	doFuncCount int
	doFunc      func(interface{}) (interface{}, error)
	started     bool
}

func newConcurrentItem(
	inputChan chan interface{},
	doFuncCount int,
	doFunc func(interface{}) (interface{}, error),
	outputChan chan interface{},
) *concurrentItem {
	return &concurrentItem{
		lock:        new(sync.Mutex),
		wait:        new(sync.WaitGroup),
		inputChan:   inputChan,
		doFuncCount: doFuncCount,
		doFunc:      doFunc,
		outputChan:  outputChan,
	}
}

func (ci *concurrentItem) do() {
	ci.lock.Lock()
	defer ci.lock.Unlock()
	if !ci.started {
		for i := 0; i < ci.doFuncCount; i++ {
			go ci.f()
		}
		ci.started = true
	}
}

func (ci *concurrentItem) f() {
	ci.wait.Add(1)
	for {
		v := <-ci.inputChan
		if v == nil {
			ci.wait.Done()
			return
		}
		i, e := ci.doFunc(v)
		if e == nil && ci.outputChan != nil {
			ci.outputChan <- i
		}
	}
}

func (ci *concurrentItem) stop() {
	ci.lock.Lock()
	defer ci.lock.Unlock()
	for i := 0; i < ci.doFuncCount; i++ {
		ci.inputChan <- nil
	}
	ci.wait.Wait()
	ci.started = false
}
