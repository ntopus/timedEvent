package queue_publisher

import (
	"github.com/ivanmeca/timedEvent/application/modules/config"
	"github.com/ivanmeca/timedEvent/application/modules/logger"
	"github.com/kofre/lib-queue/rabbitmq_repository"
	"os"
	"strconv"
	"sync"
)

var once sync.Once
var instance *queue_publisher

type queue_publisher struct {
	queueRepo *rabbitmq_repository.QueueRepository
	queueMap  map[string]*rabbitmq_repository.Queue
}

func QueuePublisher() *queue_publisher {
	once.Do(func() {
		instance = &queue_publisher{}
		instance.init()
	})
	return instance
}

func (qp *queue_publisher) init() {
	AppLogger := logger.GetLogger()
	QueueConf := config.GetConfig().PublishQueue
	for _, qConf := range QueueConf {
		port, err := strconv.Atoi(qConf.ServerPort)
		if err != nil {
			AppLogger.ErrorPrintln("could not get queue port on queue " + qConf.QueueName)
			os.Exit(1)
		}
		qr, err := rabbitmq_repository.NewQueueRepository(rabbitmq_repository.NewQueueRepositoryParams(qConf.ServerUser, qConf.ServerPassword, qConf.ServerHost, port))
		if err != nil {
			AppLogger.ErrorPrintln("could not init queue repository on queue " + qConf.QueueName)
			os.Exit(1)
		}
		queueName := qConf.QueueName
		qParam := rabbitmq_repository.NewQueueParams(queueName)
		qParam.SetThreadLimit(200)
		q, err := qr.QueueDeclare(qParam, false)
		if err != nil {
			AppLogger.ErrorPrintln("could not declare queue " + queueName)
			os.Exit(1)
		}
		if queue, ok := q.(*rabbitmq_repository.Queue); ok {
			qp.queueMap[queueName] = queue
		}
	}
}

func (qp *queue_publisher) ValidateQueue(queueName string) bool {
	if _, ok := qp.queueMap[queueName]; ok {
		return true
	}
	return false
}

func (qp *queue_publisher) PublishInQueue(queueName string, data interface{}) bool {
	AppLogger := logger.GetLogger()
	if val, ok := qp.queueMap[queueName]; ok {
		err := val.Publish(data)
		if err != nil {
			AppLogger.ErrorPrintln("could not publish on queue: " + err.Error())
			return false
		}
		return true
	}
	return false
}
