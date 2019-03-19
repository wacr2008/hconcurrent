package hconcurrent

import (
	"sync"
	"time"
)

type Option struct {
	channelSize    int
	goroutineCount int
	handle         func(interface{}) interface{}
}

func NewOption(channelSize, goroutineCount int, handle func(interface{}) interface{}) Option {
	return Option{channelSize: channelSize, goroutineCount: goroutineCount, handle: handle}
}

type Concurrent struct {
	lock            *sync.Mutex
	concurrentItems []*concurrentItem
	inputChan       chan interface{}
	started         bool
}

func NewConcurrent(channelSize int, goroutineCount int, handle func(interface{}) interface{}) *Concurrent {
	return NewConcurrentWithOptions(NewOption(channelSize, goroutineCount, handle))
}

//
//                     |-      -|                    |-      -|                    |-      -|
//                     | handle |                    | handle |                    | handle |
//                     | handle |                    | handle |                    | handle |
//                     | handle |                    | handle |                    | handle |
// >>> input chan >>> -| handle |- >>> out chan >>> -| handle |- >>> out chan >>> -| handle |
//    (channelSize)    | handle |    (channelSize)   | handle |    (channelSize)   | handle |
//                     | handle |                    | handle |                    | handle |
//                     |   .    |                    |   .    |                    |   .    |
//                     |   .    |                    |   .    |                    |   .    |
//                     |-  .   -|                    |-  .   -|                    |-  .   -|
//               goroutineCount x handle       goroutineCount x handle       goroutineCount x handle
//
func NewConcurrentWithOptions(options ...Option) *Concurrent {
	c := &Concurrent{
		lock: new(sync.Mutex),
	}
	c.initConcurrentItems(options)
	return c
}

func (c *Concurrent) initConcurrentItems(options []Option) {
	c.concurrentItems = []*concurrentItem{}
	var preItem *concurrentItem
	for _, op := range options {
		if op.channelSize <= 0 || op.goroutineCount <= 0 || op.handle == nil {
			panic("init concurrent fail,option should be:channelSize>0,goroutineCount>0,handle!=nil")
		}
		item := newConcurrentItem(op)
		c.concurrentItems = append(c.concurrentItems, item)

		if preItem != nil {
			preItem.setOutputChan(item.inputChan)
		}
		if c.inputChan == nil {
			c.inputChan = item.inputChan
		}

		preItem = item
	}
}

func (c *Concurrent) Start() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.started {
		c.run()
		c.started = true
	}
}

func (c *Concurrent) run() {
	for _, item := range c.concurrentItems {
		item.start()
	}
}

func (c *Concurrent) Input(i interface{}) bool {
	select {
	case c.inputChan <- i:
		return true
	default:
		return false
	}
}

func (c *Concurrent) MustInput(i interface{}) bool {
	c.inputChan <- i
	return true
}

func (c *Concurrent) InputWithTimeout(i interface{}, timeout time.Duration) bool {
	return c.InputWithTimer(i, time.NewTimer(timeout))
}

func (c *Concurrent) InputWithTimer(i interface{}, t *time.Timer) bool {
	select {
	case c.inputChan <- i:
		return true
	case <-t.C:
		return false
	}
}

func (c *Concurrent) Stop() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if !c.started {
		return
	}
	w := new(sync.WaitGroup)
	for _, item := range c.concurrentItems {
		w.Add(item.goroutineCount)
		item.stop(w)
		w.Wait()
	}
	c.started = false
}
