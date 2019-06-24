package timer_control

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"sync"
	"time"
)

type TimerControl struct {
	controlTime time.Duration
	list        *sync.Map
}

func NewTimerControl(controlTime int, list *sync.Map) *TimerControl {
	return &TimerControl{list: list, controlTime: time.Duration(controlTime)}
}

func (tc *TimerControl) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				tc.processList()
			}
		}
	}()
}

func (tc *TimerControl) processList() {
	time.Sleep(tc.controlTime * time.Second)
	horaAtual := time.Now()
	tc.list.Range(func(key interface{}, value interface{}) bool {
		if event, ok := value.(data_types.EventMapper); ok {
			if horaAtual.Sub(event.PublishDate) > 0 {
				return true
			}
		}
		return false
	})
	fmt.Println("TC")
}
