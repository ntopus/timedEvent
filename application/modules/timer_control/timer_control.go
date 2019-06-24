package timer_control

import "sync"

type TimerControl struct {
	list *sync.Map
}

func (tc *TimerControl) Check() {
	tc.list.Range(func(key interface{}, value interface{}) bool {
		return true
	})
}
