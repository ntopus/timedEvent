package main

import (
	"fmt"
	"github.com/ivanmeca/timedEvent/tests"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"testing"
	"time"
)

func TestApplication(t *testing.T) {
	gomega.RegisterTestingT(t)
	gomega.RegisterFailHandler(ginkgo.Fail)
	fmt.Println("Starting application")
	ginkgo.RunSpecs(t, "main_test_suite")
}

var App *gexec.Session

var _ = ginkgo.Describe("main_test_suite", func() {
	ginkgo.BeforeSuite(func() {
		tests.BuildApplication()
		tests.SaveConfigFile()
		App = tests.RunApp()
		time.Sleep(time.Second)
	})
	ginkgo.AfterSuite(func() {
		fmt.Println("Killing application")
		App.Kill()
	})
	ginkgo.Context("Test DB generator", tests.CreateEventRequest)
})
