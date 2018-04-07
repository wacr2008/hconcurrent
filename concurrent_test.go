package hconcurrent

import (
	"sync"
	"testing"
)

func TestConcurrent(t *testing.T) {
	inputChanSize := 4
	doFuncCount := 4
	listCount := 3
	inputs := []int{1, 2, 3, 4, 5}
	list := []interface{}{}
	for i := 0; i < listCount; i++ {
		list = append(list, inputChanSize, doFuncCount,
			func(i interface{}) (interface{}, error) {
				n := i.(int)
				return n + 1, nil
			})
	}
	m := map[int]int{}
	l := new(sync.Mutex)
	list = append(list, inputChanSize, doFuncCount,
		func(i interface{}) (interface{}, error) {
			l.Lock()
			defer l.Unlock()
			n := i.(int)
			m[n] = n
			return n, nil
		})
	c := NewConcurrent2(list...)
	c.Do()
	for i := 0; i < len(inputs); i++ {
		c.Input(inputs[i])
	}
	c.Stop()
	for i := 0; i < len(inputs); i++ {
		n := inputs[i] + listCount
		if m[n] != n {
			t.Errorf("concurrent do error")
		}
	}
}

func testDo(i interface{}) (interface{}, error) {
	n := i.(int)
	return n + 1, nil
}
