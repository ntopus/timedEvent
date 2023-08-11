package scheduler

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/application/modules/config"
	. "github.com/onsi/gomega"
	"time"
	"testing"
)

// unity test disabled
func TestSchedulerPoll(test *testing.T) {
	return
	RegisterTestingT(test)
	c := config.GetConfig()
	c.DataBase.ServerHost = "http://localhost"
	c.DataBase.ServerPort = "8529"
	c.DataBase.ServerUser = "root"
	c.DataBase.ServerPassword = "rootpass"
	c.DataBase.DbName = "testDb"
	fmt.Println("Trying to poll database")
	scheduler := NewScheduler(2, 1, 1800)
	scheduler.Run(context.Background())
	time.Sleep(20000 * time.Second)
}
