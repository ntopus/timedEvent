package timer_control

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"sync"
	"time"
)

type TimerControl struct {
	exclusionTime time.Duration
	controlTime   time.Duration
	list          *sync.Map
}

func NewTimerControl(controlTime int, exclusionTime int, list *sync.Map) *TimerControl {
	return &TimerControl{list: list, controlTime: time.Duration(controlTime), exclusionTime: time.Duration(exclusionTime)}
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
	horaAtual := time.Now().UTC()
	tc.list.Range(func(key interface{}, value interface{}) bool {
		if event, ok := value.(data_types.EventMapper); ok {
			timeDiffInSecond := horaAtual.Sub(event.PublishDate)
			timeDiffInSecond /= time.Second
			fmt.Printf("Hora atual: %s, hora do evento: %s\n", horaAtual.Format("2006-01-02 15:04:05"), event.PublishDate.Format("2006-01-02 15:04:05"))
			if timeDiffInSecond > tc.exclusionTime {
				_, err := collection_managment.NewEventCollection().DeleteItem([]string{event.EventID})
				if err != nil {
					fmt.Printf("falha ao excluir ID: " + event.EventID)
				}
				tc.list.Delete(key)
				fmt.Println("ID excluido: " + event.EventID)
			} else {
				if timeDiffInSecond > 0 {
					fmt.Println("Publicar ID" + event.EventID)
					//todo: publicar evento
				}
			}
		}
		return true
	})
	fmt.Println("TC")
}
