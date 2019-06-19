package scheduler

import (
	"context"
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestCloudEventEntry(test *testing.T) {
	RegisterTestingT(test)
	fmt.Println("Trying to generate a event entry")
	scheduler := NewScheduler(2)
	scheduler.Run(context.Background())
	time.Sleep(20 * time.Second)
	//Expect(err).ShouldNot(HaveOccurred())
}
