package timer_control

import (
	"github.com/ivanmeca/timedEvent/application/modules/scheduler"
	"sync"
	"time"
)

type TimerControl struct {
	controlTime time.Duration
	list        *sync.Map
}

func NewTimerControl(list *sync.Map) *TimerControl {
	return &TimerControl{list: list}
}

func (tc *TimerControl) Run() {
	go func() {
		for {
			time.Sleep(time.Second)
			horaAtual := time.Now()
			tc.list.Range(func(key interface{}, value interface{}) bool {
				if event, ok := value.(scheduler.EventMapper); ok {
					if horaAtual.Sub(event.PublishDate) > 0 {
						return true
					}
				}
				return false
			})
		}
	}()
}
