package scheduler

import (
	"context"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/logger"
	"github.com/ivanmeca/timedEvent/application/modules/timer_control"
	"sync"
	"time"
)

type Scheduler interface {
	Run(ctx context.Context)
	CheckEvent(event *data_types.ArangoCloudEvent)
}

var instance *EventScheduler

func GetScheduler() Scheduler {
	return instance
}

func NewScheduler(pollTime int, controlTime int, expirationTime int) Scheduler {
	instance = &EventScheduler{poolTime: time.Duration(pollTime)}
	instance.tc = timer_control.NewTimerControl(controlTime, expirationTime, &instance.eventList)
	return instance
}

type EventScheduler struct {
	controlTime time.Duration
	poolTime    time.Duration
	eventList   sync.Map
	tc          *timer_control.TimerControl
	logger      *logger.StdLogger
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
	es.logger = logger.GetLogger()
	es.tc.Run(ctx)
}

func (es *EventScheduler) CheckEvent(event *data_types.ArangoCloudEvent) {
	horaAtual := time.Now().UTC()
	publishDate, err := time.Parse("2006-01-02 15:04:05Z", event.PublishDate)
	if err != nil {
		es.logger.ErrorPrintln("error on date parsing (value.PublishDate,event id: " + event.ArangoKey + ") : " + err.Error())
		return
	}
	timeDiffInSecond := horaAtual.Sub(publishDate)
	timeDiffInSecond /= time.Second

	if timeDiffInSecond >= (es.poolTime * -1) {
		ev := data_types.EventMapper{}
		ev.PublishDate = publishDate
		ev.Event = *event
		ev.EventRevision = event.ArangoRev
		ev.EventID = event.ArangoKey
		es.eventList.Store(event.ID, ev)
	}
}

func (es *EventScheduler) pooler() {
	horaAtual := time.Now().UTC()
	horaLimite := horaAtual.Add(es.poolTime * time.Second).Format("2006-01-02 15:04:05Z")
	es.logger.DebugPrintln("Scheduler:" + horaAtual.Format("2006-01-02 15:04:05Z") + " timeLimit:" + horaLimite)
	data, err := collection_managment.NewEventCollection().Read([]database.AQLComparator{{Field: "publishdate", Comparator: "<=", Value: horaLimite}})
	if err != nil {
		return
	}
	for _, value := range data {
		ev := data_types.EventMapper{}
		publishDate, err := time.Parse("2006-01-02 15:04:05Z", value.PublishDate)
		if err != nil {
			es.logger.ErrorPrintln("error on date parsing (value.PublishDate,event id: " + value.ArangoKey + ") : " + err.Error())
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
