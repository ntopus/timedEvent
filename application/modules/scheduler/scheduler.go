package scheduler

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"sync"
	"time"
)

type EventMapper struct {
	publishDate   time.Time
	eventRevision string
	eventID       string
	event         data_types.ArangoCloudEvent
}

type Scheduler interface {
	Run(ctx context.Context)
}

func NewScheduler(pollTime int) Scheduler {
	return &EventScheduler{poolTime: time.Duration(pollTime)}
}

type EventScheduler struct {
	poolTime  time.Duration
	eventList sync.Map
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
	go func() {
		for {
			time.Sleep(time.Second)
			es.timerControl()
		}
	}()
}

func (es *EventScheduler) pooler() {
	fmt.Println("Pooler" + time.Now().Format("2006-01-02 15:04:05Z"))
	es.DBPoll()
}

func (es *EventScheduler) timerControl() {
	fmt.Println("Control" + time.Now().Format("2006-01-02 15:04:05Z"))
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
		}
		ev.publishDate = publishDate
		ev.event = value
		ev.eventRevision = value.ArangoRev
		ev.eventID = value.ArangoId
		es.eventList.Store(value.ID, value)
	}
	fmt.Println("Pool ok")
	return
}
