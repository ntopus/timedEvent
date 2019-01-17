package queue_repository

import (
	"devgit.kf.com.br/core/lib-queue/queue"
	"github.com/streadway/amqp"
	"strconv"
)

type QueueRepository struct {
	params     QueueRepositoryParams
	connection *amqp.Connection
}

func NewQueueRepository(params QueueRepositoryParams) (*QueueRepository, error) {
	queueRp := QueueRepository{params: params}

	auth := amqp.PlainAuth{Username: queueRp.params.Login(), Password: queueRp.params.Password()}
	var arrAuth []amqp.Authentication
	arrAuth = append(arrAuth, &auth)
	config := amqp.Config{
		SASL: arrAuth,
	}
	conn, err := amqp.DialConfig("amqp://"+queueRp.params.Host()+":"+strconv.Itoa(queueRp.params.Port())+"/", config)
	if err != nil {
		return nil, err
	}
	queueRp.connection = conn
	return &queueRp, nil
}

func (q *QueueRepository) QueueBind(params QueueBindParams) error {
	activeChannel, err := q.connection.Channel()
	if err != nil {
		return err
	}
	if err := activeChannel.QueueBind(
		params.Name(),
		params.Key(),
		params.Exchange(),
		params.NoWait(),
		params.Args()); err != nil {
		return err
	}
	return nil
}

func (q *QueueRepository) QueueDeclare(params queue.QueueParams, withErrorQueue bool) (*queue.Queue, error) {
	if withErrorQueue {
		q.errorQueueDeclare(params)
	}
	return q.queueDeclare(params)
}

func (q *QueueRepository) ExchangeDeclare(params queue.QueueParams) error {
	activeChannel, err := q.connection.Channel()
	if err != nil {
		return err
	}
	if err := activeChannel.ExchangeDeclare(
		params.Name(),
		params.Kind(),
		params.Durable(),
		params.AutoDelete(),
		params.Internal(),
		params.NoWait(),
		params.Args()); err != nil {
		return err
	}
	return nil
}

func (q *QueueRepository) errorQueueDeclare(params queue.QueueParams) error {

	errorQueueName := params.Name() + "-error"

	exchangeParams := queue.NewQueueParams("Error")
	q.ExchangeDeclare(exchangeParams)

	args := map[string]interface{}{
		"x-dead-letter-exchange":    "Error",
		"x-dead-letter-routing-key": errorQueueName,
	}
	qParam := queue.NewQueueParams(errorQueueName)
	qParam.SetArgs(args)
	_, err := q.queueDeclare(qParam)

	if err != nil {
		return err
	}

	err = q.QueueBind(NewQueueBindParams(errorQueueName, errorQueueName, "Error"))
	if err != nil {
		return err
	}
	return nil
}

func (q *QueueRepository) queueDeclare(params queue.QueueParams) (*queue.Queue, error) {
	queue, err := queue.NewQueue(params, q.connection)
	if err != nil {
		return nil, err
	}
	return queue, nil
}
