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

type EventMapper struct {
	PublishDate   time.Time
	EventRevision string
	EventID       string
	Event         data_types.ArangoCloudEvent
}

type Scheduler interface {
	Run(ctx context.Context)
}

func NewScheduler(pollTime int) Scheduler {
	sc := &EventScheduler{poolTime: time.Duration(pollTime)}
	sc.tc = timer_control.NewTimerControl(&sc.eventList)
	return sc
}

type EventScheduler struct {
	poolTime  time.Duration
	eventList sync.Map
	tc        *timer_control.TimerControl
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
}

func (es *EventScheduler) pooler() {
	es.DBPoll()
}

func (es *EventScheduler) DBPoll() {
	horaAtual := time.Now()
	data, err := collection_managment.NewEventCollection().Read([]database.AQLComparator{{Field: "publishdate", Comparator: "<=", Value: horaAtual.Add(es.poolTime).Format("2006-01-02 15:04:05Z")}})
	if err != nil {
		return
	}
	for _, value := range data {
		ev := EventMapper{}

		publishDate, err := time.Parse("2006-01-02 15:04:05Z", value.PublishDate)
		if err != nil {
			fmt.Println("Erro no parse da data")
			continue
		}
		ev.PublishDate = publishDate
		ev.Event = value
		ev.EventRevision = value.ArangoRev
		ev.EventID = value.ArangoId
		es.eventList.Store(value.ID, value)
	}
	fmt.Println("Pool ok")
	return
}
