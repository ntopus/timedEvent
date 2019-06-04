package scheduler

import "github.com/ivanmeca/timedEvent/application/modules/database/data_types"

type EventMapper struct {
	eventRevision string
	eventID       string
	event         data_types.ArangoCloudEvent
}

type EventScheduler struct {
	eventList map[string]EventMapper
}

func (es *EventScheduler) DBPoll() {

}
