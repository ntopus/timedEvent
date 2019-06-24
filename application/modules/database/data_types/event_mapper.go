package data_types

import "time"

type EventMapper struct {
	PublishDate   time.Time
	EventRevision string
	EventID       string
	Event         ArangoCloudEvent
}
