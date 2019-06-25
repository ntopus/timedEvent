package scheduler

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/timer_control"
	"sync"
	"time"
)

type Scheduler interface {
	Run(ctx context.Context)
}

func NewScheduler(pollTime int, controlTime int, exclusionTime int) Scheduler {
	sc := &EventScheduler{poolTime: time.Duration(pollTime)}
	sc.tc = timer_control.NewTimerControl(controlTime, exclusionTime, &sc.eventList)
	return sc
}

type EventScheduler struct {
	controlTime time.Duration
	poolTime    time.Duration
	eventList   sync.Map
	tc          *timer_control.TimerControl
}

func (es *EventScheduler) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(es.poolTime * time.Second)
				es.pooler()
			}
		}
	}()
	es.tc.Run(ctx)
}

func (es *EventScheduler) pooler() {
	horaAtual := time.Now().UTC()
	data, err := collection_managment.NewEventCollection().Read([]database.AQLComparator{{Field: "publishdate", Comparator: "<=", Value: horaAtual.Add(es.poolTime).Format("2006-01-02 15:04:05Z")}})
	if err != nil {
		return
	}
	for _, value := range data {
		ev := data_types.EventMapper{}
		publishDate, err := time.Parse("2006-01-02 15:04:05Z", value.PublishDate)
		if err != nil {
			fmt.Println("Erro no parse da data")
			continue
		}
		ev.PublishDate = publishDate
		ev.Event = value
		ev.EventRevision = value.ArangoRev
		ev.EventID = value.ArangoKey
		es.eventList.Store(value.ID, ev)
	}
	return
}
