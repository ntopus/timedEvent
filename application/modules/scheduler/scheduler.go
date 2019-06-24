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
	elapsedTime   int64
	totalTime     int64
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
}

func (es *EventScheduler) pooler() {
	es.DBPoll()
}

func (es *EventScheduler) DBPoll() {

	horaAtual := time.Now().Format("2006-01-02 15:04:05Z")
	data, err := collection_managment.NewEventCollection().Read([]database.AQLComparator{{Field: "publishdate", Comparator: ">=", Value: horaAtual}})
	if err != nil {
	}
	fmt.Println(data)
}
