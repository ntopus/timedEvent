package scheduler

import (
	"context"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/logger"
	"github.com/ivanmeca/timedEvent/application/modules/queue_publisher"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type Scheduler interface {
	Run(ctx context.Context)
	CheckEvent(event *data_types.ArangoCloudEvent)
}

const TimerControlUnit = time.Millisecond

var instance *EventScheduler

type FnTimer func()

func GetScheduler() Scheduler {
	return instance
}

func NewScheduler(pollTime int, controlTime int, expirationTime int) Scheduler {
	instance = &EventScheduler{
		poolTime:       time.Duration(pollTime),
		eventTimerList: data_types.NewEventMapper(),
	}
	return instance
}

type EventScheduler struct {
	poolTime       time.Duration
	eventTimerList data_types.EventMapper
	logger         *logger.StdLogger
	mu             sync.Mutex
}

func (es *EventScheduler) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(es.poolTime * TimerControlUnit)
				es.pooler()
			}
		}
	}()
	es.logger = logger.GetLogger()
}

func (es *EventScheduler) CheckEvent(event *data_types.ArangoCloudEvent) {
	horaAtual := time.Now().UTC()
	publishDate, err := time.Parse("2006-01-02 15:04:05Z", event.PublishDate)
	if err != nil {
		es.logger.ErrorPrintln("error on date parsing (value.PublishDate,event id: " + event.ArangoKey + ") : " + err.Error())
		return
	}
	timeDiffInSecond := horaAtual.Sub(publishDate)
	timeDiffInSecond /= TimerControlUnit

	if timeDiffInSecond >= (es.poolTime * -1) {
		ev := data_types.EventMapperEntry{}
		ev.PublishDate = publishDate
		ev.Event = *event
		ev.EventRevision = event.ArangoRev
		ev.EventID = event.ArangoKey
		es.scheduleEvent(&ev)
	}
}

func (es *EventScheduler) pooler() {
	horaAtual := time.Now().UTC()
	horaLimite := horaAtual.Add(es.poolTime * TimerControlUnit).Format("2006-01-02 15:04:05Z")
	es.logger.DebugPrintln("Scheduler:" + horaAtual.Format("2006-01-02 15:04:05Z") + " timeLimit:" + horaLimite)
	data, err := collection_managment.NewEventCollection().Read([]database.AQLComparator{{Field: "publishdate", Comparator: "<=", Value: horaLimite}})
	if err != nil {
		return
	}
	for _, value := range data {
		ev := data_types.EventMapperEntry{}
		publishDate, err := time.Parse("2006-01-02 15:04:05Z", value.PublishDate)
		if err != nil {
			es.logger.ErrorPrintln("error on date parsing (value.PublishDate,event id: " + value.ArangoKey + ") : " + err.Error())
			continue
		}
		ev.PublishDate = publishDate
		ev.Event = value
		ev.EventRevision = value.ArangoRev
		ev.EventID = value.ArangoKey
		es.scheduleEvent(&ev)
	}
	return
}

func (es *EventScheduler) scheduleEvent(event *data_types.EventMapperEntry) {
	listEntry, ok := es.eventTimerList.Load(event.EventID)
	if ok {
		if listEntry.EventRevision == event.EventRevision {
			return
		}
		listEntry.ControlTimer.Stop()
	}
	horaAtual := time.Now().UTC()
	timeDiff := event.PublishDate.Sub(horaAtual)
	event.ControlTimer = time.AfterFunc(timeDiff, es.buildPublishFunc(event))
	es.eventTimerList.Store(event.EventID, *event)
}

func (es *EventScheduler) buildPublishFunc(event *data_types.EventMapperEntry) FnTimer {
	return func() {
		es.mu.Lock()
		defer es.mu.Unlock()
		defer es.eventTimerList.Delete(event.EventID)
		data, err := collection_managment.NewEventCollection().ReadItem(event.EventID)
		if err != nil || data == nil {
			es.logger.NoticePrintln(errors.Wrap(err, "event check fail").Error())
			return
		}
		if data.ArangoRev != event.EventRevision {
			es.logger.DebugPrintln("event rev check fail")
			return
		}
		es.logger.DebugPrintln("Publicar ID " + event.EventID)
		var dataToPublish interface{}
		if event.Event.PublishType == data_types.DataOnly {
			dataToPublish = event.Event.CloudEvent.Data
		} else {
			dataToPublish = event.Event.CloudEvent
		}
		if queue_publisher.QueuePublisher().PublishInQueue(data.PublishQueue, dataToPublish) {
			_, err := collection_managment.NewEventCollection().DeleteItem([]string{event.EventID})
			if err != nil {
				es.logger.ErrorPrintln(errors.Wrap(err, "falha ao excluir ID (PublishOK): "+event.EventID).Error())
			}
			es.logger.DebugPrintln("ID excluido: " + event.EventID)
		}
	}
}
