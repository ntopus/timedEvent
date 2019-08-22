package timer_control

import (
	"fmt"
	"github.com/onsi/gomega"
	"sync"
	"testing"
	"time"
)

func TestReadDocumentsWithFilter(test *testing.T) {
	gomega.RegisterTestingT(test)
	//fmt.Println("Trying to a read collection with filters")
	//horaAtual := time.Now().AddDate(0, 0, 3)
	//publishedDate, err := time.Parse("2006-01-02 15:04:05Z", mock.PublishDate)
	//timeDiff := time.Now().UTC().Sub(publishedDate)
	//gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

	wg := sync.WaitGroup{}
	wg.Add(1)
	valor := 5
	timeDiff := 5000 * time.Millisecond
	time.AfterFunc(timeDiff, func() {
		func(delay int) {
			defer wg.Done()
			fmt.Println("estorou com ", timeDiff)
		}(valor)
	})
	time.Sleep(time.Second)
	//t.Stop()
	time.Sleep(5 * time.Second)
	//wg.Wait()
}
