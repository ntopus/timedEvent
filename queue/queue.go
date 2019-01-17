package queue

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type fnPublish func(msg interface{})
type fnConsume func(queueName string, msg []byte) bool

type Queue struct {
	activeChannel       *amqp.Channel
	cancelQueue         context.CancelFunc
	conn                *amqp.Connection
	confirm             chan amqp.Confirmation
	consumerTag         string
	context             context.Context
	mutex               sync.Mutex
	name                string
	params              QueueParams
	queueChErr          chan error
	threadLimit         int
	threadCount         int
	wg                  sync.WaitGroup
	OnPublishedEvent    fnPublish
	OnNotPublishedEvent fnPublish
}

func NewQueue(params QueueParams, connection *amqp.Connection) (*Queue, error) {
	nq := Queue{
		name:   params.Name(),
		conn:   connection,
		params: params,
	}
	var err error
	nq.activeChannel, err = nq.conn.Channel()
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	nq.context, nq.cancelQueue = context.WithCancel(ctx)
	if err != nil {
		return nil, err
	}
	if err := nq.activeChannel.Qos(params.ThreadLimit(), 0, true); err != nil {
		return nil, err
	}
	_, err = nq.activeChannel.QueueDeclare(
		params.Name(),
		params.Durable(),
		params.AutoDelete(),
		params.Exclusive(),
		params.NoWait(),
		params.Args(),
	)
	if err != nil {
		return nil, err
	}
	nq.confirm = nq.activeChannel.NotifyPublish(make(chan amqp.Confirmation, 1))
	if err := nq.activeChannel.Confirm(false); err != nil {
		return nil, err
	}
	nq.threadLimit = params.ThreadLimit()
	return &nq, nil
}

func (q *Queue) Close() {
	q.cancelQueue()
	q.activeChannel.Close()
	q.conn.Close()
}

func (q *Queue) Publish(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	q.mutex.Lock()
	defer q.mutex.Unlock()
	err = q.activeChannel.Publish(
		"",
		q.name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		return err
	}
	if confirmed := <-q.confirm; !confirmed.Ack {
		if q.OnNotPublishedEvent != nil {
			q.OnNotPublishedEvent(message)
		}
		return errors.New("message could not be published")
	}

	if q.OnPublishedEvent != nil {
		go q.OnPublishedEvent(message)
	}

	return nil
}

func (q *Queue) Cancel() error {
	return q.activeChannel.Cancel(q.consumerTag, false)
}

func (q *Queue) ErrorConsume() <-chan error {
	return q.queueChErr
}

func (q *Queue) StartConsume(consumeHandler fnConsume) error {

	q.queueChErr = make(chan error)
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	q.consumerTag = fmt.Sprintf("%x", b)
	consumerChannel, err := q.activeChannel.Consume(q.name, q.consumerTag, q.params.AutoAck(), q.params.Exclusive(), false, q.params.NoWait(), nil)
	if err != nil {
		return err
	}

	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for {
			select {
			case <-q.context.Done():
				fmt.Println("Closing queue")
				return
			case message, ok := <-consumerChannel:
				if !ok {
					err = q.activeChannel.Cancel(message.ConsumerTag, q.params.NoWait())
					if err != nil {
						q.queueChErr <- err
					}
					continue
				}
				q.waitThreadLimit()
				q.incThreadCount()
				go q.treatMessage(message, consumeHandler)
			}
		}
	}()
	return nil
}

func (q *Queue) GetThreadCount() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return q.threadCount
}

func (q *Queue) SetThreadLimit(limit int) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.threadLimit = limit
	return
}

func (q *Queue) WaitQueue() {
	for q.GetThreadCount() != 0 {
		time.Sleep(10 * time.Millisecond)
	}
	q.wg.Wait()
}

func (q *Queue) decThreadCount() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if q.threadCount > 0 {
		q.threadCount--
	}
	return
}

func (q *Queue) incThreadCount() {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.threadCount++
	return
}

func (q *Queue) treatMessage(message amqp.Delivery, consumeHandler fnConsume) {
	defer q.decThreadCount()
	ok := consumeHandler(q.name, message.Body)
	if ok {
		err := message.Ack(false)
		if err != nil {
			q.queueChErr <- err
		}
	} else {
		err := message.Nack(false, false)
		if err != nil {
			q.queueChErr <- err
		}
	}
	return
}

func (q *Queue) waitThreadLimit() {
	for q.GetThreadCount() >= q.threadLimit {
		time.Sleep(100 * time.Millisecond)
	}
	return
}
