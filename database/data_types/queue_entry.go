package data_types

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type QueueEntry struct {
	Id              bson.ObjectId         `json:"id" bson:"_id"`
	QueueName       string                `json:"destination_queue_name"`
	QueueRepository QueueRepositoryParams `json:"queue_repository"`
	Payload         interface{}           `json:"payload"`
	PublishDate     time.Time             `json:"publish_date"`
}
