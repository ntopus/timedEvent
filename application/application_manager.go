package application

//
//import (
//	"context"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/app_config"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/app_log"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/converter/flight_converter"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/converter/flight_converter/flight_mapper"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/converter/flight_converter/flight_mapper/arrival_flight_mapper"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/converter/flight_converter/flight_mapper/departure_flight_mapper"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/http_request"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/import_routine"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/queue_producer"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/queue_producer/event_creator"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/queue_service/queue_connection"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/queue_service/queue_event_repository"
//	"devgit.kf.com.br/customers/galeao-flight-integration/importer/application/modules/task_ticker"
//	"github.com/micro/go-config"
//	"github.com/micro/go-config/source/file"
//	"github.com/onsi/gomega"
//	"time"
//)
//
//const timeoutQueueMessageAck = 30
//
//type ApplicationManager struct {
//	configPath string
//	config     app_config.AppConfig
//	services   ApplicationServices
//	cancelChan chan interface{}
//}
//
//type ApplicationServices struct {
//	taskTicker      task_ticker.ITaskTicker
//	queueRepository queue_event_repository.IQueueEventRepository
//	cancelTasks     []func()
//}
//
//func NewApplicationManager(configPath string) *ApplicationManager {
//	return &ApplicationManager{
//		configPath: configPath,
//		cancelChan: make(chan interface{}),
//	}
//}
//
//func (app *ApplicationManager) RunApplication(ctx context.Context) error {
//	err := app.loadConfig()
//	if err != nil {
//		return err
//	}
//	err = app.initializeServices()
//	if err != nil {
//		return err
//	}
//	app.initializeArrivalImport()
//	app.initializeDepartureImport()
//	go app.waitForGracefulStop(ctx)
//	return nil
//}
//
//func (app *ApplicationManager) initializeServices() error {
//	app.services.taskTicker = task_ticker.NewTaskTicker()
//	stanConn, err := queue_connection.GetStanConn(app.config.TargetQueueConfig.ClusterID, app.config.TargetQueueConfig.ClientID, app.config.TargetQueueConfig.Address)
//	if err != nil {
//		return err
//	}
//	queueRepository := queue_repository.NewQueueRepository(stanConn, timeoutQueueMessageAck, int(app.config.ThreadLimit))
//	app.services.queueRepository = queue_event_repository.NewEventQueueRepository(queueRepository)
//	app.services.cancelTasks = append(app.services.cancelTasks, func() { _ = app.services.queueRepository.CloseConn() })
//	return nil
//}
//
//func (app *ApplicationManager) initializeDepartureImport() {
//	app.initializeImportRoutine(app.config.ImportSourceDepartureConfig, departure_flight_mapper.DepartureFlightMapFromRawFlight, service_models.EventImportedFlightDeparture, service_models.SubjectImportedFlightDeparture)
//}
//
//func (app *ApplicationManager) initializeArrivalImport() {
//	app.initializeImportRoutine(app.config.ImportSourceArrivalConfig, arrival_flight_mapper.ArrivalFlightMapFromRawFlight, service_models.EventImportedFlightArrival, service_models.SubjectImportedFlightArrival)
//}
//
//func (app *ApplicationManager) initializeImportRoutine(importConfig app_config.ImportSourceConfig, mapperFunc flight_mapper.FlightMapperFunc, eventContentType string, subjectQueue string) {
//	httpRequestManager := http_request.NewHttpRequestManager(importConfig.Address, importConfig.RetryAmount)
//	flightConverter := flight_converter.NewFlightConverter(mapperFunc, app.config.ThreadLimit)
//	queueProducer := queue_producer.NewQueuePr
//
//		mu.Lock()
//	count := 0
//	mu.Unlock()
//	gomega.Expect(err).To(gomega.BeNil())
//
//
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return count
//	}).Should(gomega.BeEquivalentTo(1))
//	gomega.Eventually(func() int {
//		mu.Lock()
//		defer mu.Unlock()
//		return countErr
//	}).Should(gomega.BeEquivalentTo(0))
//
//	time.Sleep(1 * time.Second)oducer(app.services.queueRepository, event_creator.NewEventCreator(eventContentType), subjectQueue)
//	routine := import_routine.NewImportRoutine(httpRequestManager, flightConverter, queueProducer, app.config.ThreadLimit)
//	cancelTask := app.services.taskTicker.NewPeriodicallyTask(routine.Routine, importConfig.ImportingRatePerSecond, func() {})
//	app.services.cancelTasks = append(app.services.cancelTasks, cancelTask)
//}
//
//func (app *ApplicationManager) loadConfig() error {
//	err := config.Load(file.NewSource(
//		file.WithPath(app.configPath),
//	))
//	if err != nil {
//		return err
//	}
//	err = config.Scan(&app.config)
//	if err != nil {
//		return err
//	}
//	app_config.SetConfig(&app.config)
//	app_log.ReloadDefaultLoggerConfig(app.config.LogConfig)
//	app_log.Debugln("load config")
//	return nil
//}
//
//func (app *ApplicationManager) waitForGracefulStop(ctx context.Context) {
//	func() {
//		<-ctx.Done()
//		app.stopApplication()
//	}()
//}
//
//func (app *ApplicationManager) stopApplication() {
//	for _, ct := range app.services.cancelTasks {
//		ct()
//	}
//	app.cancelChan <- nil
//}
//
//func (app *ApplicationManager) Done() chan interface{} {
//	return app.cancelChan
//}
