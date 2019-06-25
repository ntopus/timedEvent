package timer_control

import (
	"context"
	"devgit.kf.com.br/comercial/gateway-maxtrack/module_maxtrack/logger"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"sync"
	"time"
)

type TimerControl struct {
	expirationTime time.Duration
	controlTime    time.Duration
	list           *sync.Map
	logger         *logger.StdLogger
}

func NewTimerControl(controlTime int, expirationTime int, list *sync.Map) *TimerControl {
	return &TimerControl{list: list, controlTime: time.Duration(controlTime), expirationTime: time.Duration(expirationTime)}
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
	tc.logger = logger.GetLogger()
}

func (tc *TimerControl) processList() {
	time.Sleep(tc.controlTime * time.Second)
	horaAtual := time.Now().UTC()
	tc.list.Range(func(key interface{}, value interface{}) bool {
		if event, ok := value.(data_types.EventMapper); ok {
			timeDiffInSecond := horaAtual.Sub(event.PublishDate)
			timeDiffInSecond /= time.Second
			tc.logger.DebugPrintln("Hora atual: %s, hora do evento: %s\n", horaAtual.Format("2006-01-02 15:04:05Z"), event.PublishDate.Format("2006-01-02 15:04:05Z"))
			if timeDiffInSecond > tc.expirationTime {
				_, err := collection_managment.NewEventCollection().DeleteItem([]string{event.EventID})
				if err != nil {
					tc.logger.NoticePrintln("falha ao excluir ID: " + event.EventID)
				}
				tc.list.Delete(key)
				tc.logger.DebugPrintln("ID excluido: " + event.EventID)
			} else {
				if timeDiffInSecond > 0 {
					tc.logger.DebugPrintln("Publicar ID" + event.EventID)
					data, err := collection_managment.NewEventCollection().ReadItem(event.EventID)
					if err != nil {
						tc.logger.ErrorPrintln("event check fail: " + err.Error())
						return true
					}
					if data.ArangoRev == event.EventRevision {
						//todo: publicar evento
					} else {
						tc.list.Delete(key)
					}
				}
			}
		}
		return true
	})
}
