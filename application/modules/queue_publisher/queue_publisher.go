package queue_publisher

import "github.com/kofre/lib-queue/rabbitmq_repository"

type queue_publisher struct {
	queueRepo rabbitmq_repository.QueueRepository
	queueMap  map[string]rabbitmq_repository.Queue
}

func InitQueuePublisher() {

}
