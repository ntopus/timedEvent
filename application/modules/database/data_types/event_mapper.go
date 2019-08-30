package data_types

import (
	"sync"
	"time"
)

type EventMapperEntry struct {
	PublishDate   time.Time
	EventRevision string
	EventID       string
	Event         ArangoCloudEvent
	ControlTimer  *time.Timer
}

type EventMapper struct {
	list map[string]EventMapperEntry
	sync.Mutex
}

func (em *EventMapper) Delete(eventKey string) {
	em.Lock()
	defer em.Unlock()
	delete(em.list, eventKey)
}

func (em *EventMapper) Load(eventKey string) (EventMapperEntry, bool) {
	em.Lock()
	defer em.Unlock()
	if entry, ok := em.list[eventKey]; ok {
		return entry, true
	}
	return EventMapperEntry{}, false
}

func (em *EventMapper) LoadOrStore(eventKey string, eventEntry EventMapperEntry) (EventMapperEntry, bool) {
	em.Lock()
	defer em.Unlock()
	if entry, ok := em.list[eventKey]; ok {
		return entry, false
	}
	em.list[eventKey] = eventEntry
	return eventEntry, true
}

func (em *EventMapper) Store(eventKey string, eventEntry EventMapperEntry) {
	em.Lock()
	defer em.Unlock()
	em.list[eventKey] = eventEntry
	return
}
