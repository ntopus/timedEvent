package main

import (
	"devgit.kf.com.br/core/lib-queue/queue"
	"devgit.kf.com.br/core/lib-queue/queue_repository"
	"encoding/json"
	"fmt"
	"github.com/onsi/gomega"
	"sync"
	"testing"
	"time"
)

type Message struct {
	MessageID       string
	ParentMessageID string
	Source          string
	MsgType         int
	Version         string
	Severity        uint8
	Payload         string
	DateTime        time.Time
	Err             error
}

type singleton struct {
	Date time.Time
}

var instance *singleton
var once sync.Once

func getDate() time.Time {
	once.Do(func() {
		instance = &singleton{Date: time.Now()}
	})
	return instance.Date
}

func compareMsg(msg1 Message, msg2 Message) {
	gomega.Expect(msg1.MessageID).To(gomega.Equal(msg2.MessageID))
	gomega.Expect(msg1.ParentMessageID).To(gomega.Equal(msg2.ParentMessageID))
	gomega.Expect(msg1.Source).To(gomega.Equal(msg2.Source))
	gomega.Expect(msg1.MsgType).To(gomega.Equal(msg2.MsgType))
	gomega.Expect(msg1.Version).To(gomega.Equal(msg2.Version))
	gomega.Expect(msg1.Severity).To(gomega.Equal(msg2.Severity))
	gomega.Expect(msg1.Payload).To(gomega.Equal(msg2.Payload))
	gomega.Expect(msg1.Err).To(gomega.BeNil())
	gomega.Expect(msg2.Err).To(gomega.BeNil())
}

func getMockMsg() Message {
	msg := Message{}
	msg.MessageID = "1234234"
	msg.ParentMessageID = "4565"
	msg.Source = "Teste source"
	msg.Severity = 10
	msg.Err = nil
	msg.DateTime = getDate()
	msg.MsgType = 100
	msg.Version = "1"
	msg.Payload = "\"data\":{ \"eventDate\": { \"date\": \"2018-07-09 19:17:36.000000\",\"timezone_type\": 3,\"timezone\": \"UTC\"}"
	return msg
}

func getQueue(threadLimit int) *queue.Queue {
	qr, err := queue_repository.NewQueueRepository(queue_repository.NewQueueRepositoryParams("miseravi", "trAfr@guR36a", "srvqueue.module.ntopus.com.br", 5672))
	params := queue.NewQueueParams("newQueue")
	params.SetThreadLimit(threadLimit)
	q, err := qr.QueueDeclare(params, true)
	gomega.Expect(err).To(gomega.BeNil())
	return q
}

func publishMockMessageOnQueue(q *queue.Queue) {
	msg := getMockMsg()
	err := q.Publish(msg)
	gomega.Expect(err).To(gomega.BeNil())
}

func TestSaveMsg(test *testing.T) {
	gomega.RegisterTestingT(test)

	fmt.Println("Trying to save a message on queue")

	mu := sync.Mutex{}
	q := getQueue(5)

	mu.Lock()
	count := 0
	mu.Unlock()
	q.OnPublishedEvent = func(message interface{}) {
		mu.Lock()
		defer mu.Unlock()
		count++
		fmt.Println("Publish OK")
	}
	mu.Lock()
	countErr := 0
	mu.Unlock()
	q.OnNotPublishedEvent = func(message interface{}) {
		mu.Lock()
		defer mu.Unlock()
		countErr++
		fmt.Println("Publish Err")
	}

	err := q.StartConsume(func(queueName string, msg []byte) bool {
		fmt.Println("New message")
		return true
	})
	gomega.Expect(err).To(gomega.BeNil())

	publishMockMessageOnQueue(q)

	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return count
	}).Should(gomega.BeEquivalentTo(1))
	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return countErr
	}).Should(gomega.BeEquivalentTo(0))

	time.Sleep(1 * time.Second)
	q.Close()
	q.WaitQueue()
}

func TestGetMsg(test *testing.T) {
	gomega.RegisterTestingT(test)

	const QTDE_MSGS = 1000

	fmt.Println("Trying to consume messages on queue")

	mu := sync.Mutex{}

	q := getQueue(500)
	defer q.Close()

	mu.Lock()
	newMessage := 0
	mu.Unlock()

	mu.Lock()
	receivedMsg := Message{}
	mu.Unlock()

	err := q.StartConsume(func(queueName string, msg []byte) bool {
		mu.Lock()
		defer mu.Unlock()
		newMessage++
		err := json.Unmarshal(msg, &receivedMsg)
		if err != nil {
			return false
		}
		mockMsg := getMockMsg()
		compareMsg(receivedMsg, mockMsg)
		return true
	})
	gomega.Expect(err).To(gomega.BeNil())

	mu.Lock()
	newMessage = 0
	mu.Unlock()

	for i := 0; i < QTDE_MSGS; i++ {
		publishMockMessageOnQueue(q)
	}

	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return newMessage
	}, 2, 1).Should(gomega.BeEquivalentTo(QTDE_MSGS))

	q.Close()
	q.WaitQueue()
}

func TestDelayedGetMsgAtOnce(test *testing.T) {
	gomega.RegisterTestingT(test)

	const QTDE_MSGS = 1000

	fmt.Println("Trying to consume messages on queue (delayed process)")

	mu := sync.Mutex{}

	q := getQueue(QTDE_MSGS)
	defer q.Close()

	mu.Lock()
	newMessage := 0
	mu.Unlock()

	err := q.StartConsume(func(queueName string, msg []byte) bool {
		defer func() {
			mu.Lock()
			newMessage++
			mu.Unlock()
		}()
		receivedMsg := Message{}
		err := json.Unmarshal(msg, &receivedMsg)
		if err != nil {
			return false
		}
		mockMsg := getMockMsg()
		compareMsg(receivedMsg, mockMsg)
		time.Sleep(750 * time.Millisecond)
		return true
	})
	gomega.Expect(err).To(gomega.BeNil())

	mu.Lock()
	newMessage = 0
	mu.Unlock()

	for i := 0; i < QTDE_MSGS; i++ {
		publishMockMessageOnQueue(q)
	}
	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return newMessage
	}, 1).Should(gomega.BeEquivalentTo(QTDE_MSGS))
	time.Sleep(time.Second)
	q.Close()
	q.WaitQueue()
}

func TestDelayedGetMsg(test *testing.T) {
	gomega.RegisterTestingT(test)

	const QTDE_MSGS = 1000

	fmt.Println("Trying to consume messages on queue (delayed process)")

	mu := sync.Mutex{}

	q := getQueue(QTDE_MSGS / 2)
	defer q.Close()

	mu.Lock()
	newMessage := 0
	mu.Unlock()

	err := q.StartConsume(func(queueName string, msg []byte) bool {
		defer func() {
			mu.Lock()
			newMessage++
			mu.Unlock()
		}()
		receivedMsg := Message{}
		err := json.Unmarshal(msg, &receivedMsg)
		if err != nil {
			return false
		}
		mockMsg := getMockMsg()
		compareMsg(receivedMsg, mockMsg)
		time.Sleep(750 * time.Millisecond)
		return true
	})
	gomega.Expect(err).To(gomega.BeNil())

	mu.Lock()
	newMessage = 0
	mu.Unlock()

	for i := 0; i < QTDE_MSGS; i++ {
		publishMockMessageOnQueue(q)
	}

	gomega.Eventually(func() int {
		mu.Lock()
		defer mu.Unlock()
		return newMessage
	}, 2).Should(gomega.BeEquivalentTo(QTDE_MSGS))
	time.Sleep(time.Second)
	q.Close()
	q.WaitQueue()
}
