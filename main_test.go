package main

import (
	"context"
	"fmt"
	"github.com/ivanmeca/timedEvent/tests"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
	"time"
)

func TestApplication(t *testing.T) {
	gomega.RegisterTestingT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	fmt.Println("Starting application")
	ginkgo.RunSpecs(t, "main_test_suite")
}

var App context.Context

var _ = ginkgo.Describe("main_test_suite", func() {
	ginkgo.BeforeSuite(func() {
		tests.BuildApplication()
		tests.SaveConfigFile()
		App = tests.RunApp()
		time.Sleep(time.Second)
	})
	ginkgo.AfterSuite(func() {
		fmt.Println("Killing application")
		App.Done()
	})
	ginkgo.BeforeEach(func() {
		tests.PurgeQueue(tests.TEST_PUBLISH_QUEUE)
	})
	ginkgo.Context("Test webserver", tests.CreateEventTester)
	ginkgo.Context("Test scheduler", tests.SchedulerTester)
})
