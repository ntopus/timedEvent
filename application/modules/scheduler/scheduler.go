package scheduler

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"sync"
	//"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"time"
)

type EventMapper struct {
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
				es.DBPoll()
			}
		}
	}()
}

func (es *EventScheduler) DBPoll() {
	fmt.Println("pool")
	data, err := collection_managment.NewEventCollection().Read(nil)
	//if err != nil {
	//}
	//fmt.Println(data)
}
