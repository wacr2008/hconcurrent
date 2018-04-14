package hconcurrent

import (
	"sync"
)

type Concurrent struct {
	lock            *sync.RWMutex
	concurrentItems []*concurrentItem
	inputChan       chan interface{}
	started         bool
	destroyed       bool
}

func NewConcurrent(
	inputChanSize int,
	doFuncCount int,
	doFunc func(interface{}) interface{},
) *Concurrent {
	return NewConcurrent2(inputChanSize, doFuncCount, doFunc)
}

//NewConcurrent2 v:intputChan->doFuncCount->doFunc->midChan->doFuncCount->doFunc->midChan->....->doFuncCount->doFunc
func NewConcurrent2(v ...interface{}) *Concurrent {
	c := &Concurrent{
		lock: new(sync.RWMutex),
	}
	c.initConcurrentItems(v...)
	return c
}

func (c *Concurrent) initConcurrentItems(v ...interface{}) {
	if c.concurrentItems != nil {
		return
	}
	c.concurrentItems = []*concurrentItem{}
	var outputChan chan interface{}
	for i := len(v) - 1; i > -1; i -= 3 {
		f := v[i].(func(interface{}) interface{})
		count := v[i-1].(int)
		size := v[i-2].(int)
		ch := make(chan interface{}, size)
		c.concurrentItems = append([]*concurrentItem{
			newConcurrentItem(ch, count, f, outputChan),
		}, c.concurrentItems...)

		outputChan = ch
	}
}

func (c *Concurrent) Do() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.started {
		c.do()
		c.started = true
	}
}

func (c *Concurrent) do() {
	for _, item := range c.concurrentItems {
		item.start()
	}
}

func (c *Concurrent) Input(i interface{}) {
	if i == nil {
		return
	}
	c.lock.RLock()
	c.concurrentItems[0].inputChan <- i
	c.lock.RUnlock()
}

func (c *Concurrent) Stop() {
	c.lock.Lock()
	c.stop()
	c.lock.Unlock()
}

//Destroy stop&destroy channels
func (c *Concurrent) Destroy() {
	c.lock.Lock()
	if c.destroyed {
		return
	}
	c.destroy()
	c.lock.Unlock()
}

func (c *Concurrent) destroy() {
	c.stop()
	for _, item := range c.concurrentItems {
		item.destroy()
	}
	c.destroyed = true
}

func (c *Concurrent) stop() {
	if !c.started {
		return
	}
	for _, item := range c.concurrentItems {
		item.stop()
	}
	c.started = false
}
