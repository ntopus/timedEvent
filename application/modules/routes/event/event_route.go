package event

import (
	"encoding/json"
	"errors"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/gin-gonic/gin"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/queue_publisher"
	"github.com/ivanmeca/timedEvent/application/modules/routes"
	"github.com/ivanmeca/timedEvent/application/modules/scheduler"
	"github.com/pborman/uuid"
	"io/ioutil"
	"net/http"
	"time"
)

func bindQueryFilterParams(context *gin.Context) []database.AQLComparator {
	var filter []database.AQLComparator
	for i, value := range context.Request.URL.Query() {
		switch i {
		case "id":
			filter = append(filter, database.AQLComparator{Field: "Context.id", Comparator: "==", Value: value[0]})
		case "initialDate":
			filter = append(filter, database.AQLComparator{Field: "Context.time", Comparator: ">=", Value: value[0]})
		case "finalDate":
			filter = append(filter, database.AQLComparator{Field: "Context.time", Comparator: "<=", Value: value[0]})
		default:
			filter = append(filter, database.AQLComparator{Field: i, Comparator: "==", Value: value[0]})
		}
	}
	return filter
}

func ceHttpCreate(context *gin.Context) (*data_types.CloudEvent, error) {
	data, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		return nil, errors.New("could not read request data: " + err.Error())
	}
	event := data_types.CloudEvent{}
	err = json.Unmarshal(data, &event)
	if err != nil {
		return nil, errors.New("could not parse request data: " + err.Error())
	}
	validationError := event.Validate()
	if validationError != nil {
		return nil, errors.New("could not read validate data: " + validationError.Error())
	}
	eventValidation := validateEvent(&event)
	if eventValidation != nil {
		return nil, errors.New("could not validate event: " + eventValidation.Error())
	}
	return &event, nil
}

func jsonHttpCreate(context *gin.Context) (*data_types.CloudEvent, error) {
	data, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		return nil, errors.New("could not read request data: " + err.Error())
	}
	event := data_types.CloudEvent{}
	err = json.Unmarshal(data, &event.Data)
	if err != nil {
		return nil, errors.New("could not parse request data: " + err.Error())
	}
	headers := context.Request.Header

	if value, ok := headers["Specversion"]; ok {
		err := event.SetSpecVersion(value[0])
		if err != nil {
			return nil, errors.New("could not get spec version: " + err.Error())
		}
	} else {
		err := event.SetSpecVersion(cloudevents.CloudEventsVersionV02)
		if err != nil {
			return nil, errors.New("could not get spec version: " + err.Error())
		}
	}
	if value, ok := headers["Type"]; ok {
		err = event.SetType(value[0])
		if err != nil {
			return nil, errors.New("could not get type: " + err.Error())
		}
	} else {
		err = event.SetType("Received.Schedule.Request")
		if err != nil {
			return nil, errors.New("could not get type: " + err.Error())
		}
	}
	if value, ok := headers["Source"]; ok {
		err = event.SetSource(value[0])
		if err != nil {
			return nil, errors.New("could not get source: " + err.Error())
		}
	} else {
		err = event.SetSource(context.ClientIP())
		if err != nil {
			return nil, errors.New("could not generate source: " + err.Error())
		}
	}
	if value, ok := headers["Id"]; ok {
		err = event.SetID(value[0])
		if err != nil {
			return nil, errors.New("could not get id: " + err.Error())
		}
	} else {
		err = event.SetID(uuid.NewUUID().String())
		if err != nil {
			return nil, errors.New("could not generate id: " + err.Error())
		}
	}
	if value, ok := headers["Content-Type"]; ok {
		err = event.SetDataContentType(value[0])
		if err != nil {
			return nil, errors.New("could not get content type: " + err.Error())
		}
	}
	if value, ok := headers["Publishdate"]; ok {
		event.PublishDate = value[0]
	} else {
		return nil, errors.New("could not get publishDate")
	}
	if value, ok := headers["Publishqueue"]; ok {
		event.PublishQueue = value[0]
	} else {
		return nil, errors.New("could not get publishQueue")
	}
	if value, ok := headers["Publishtype"]; ok {
		event.PublishType = value[0]
	}

	for name, value := range headers {
		switch name {
		case "specversion":
		case "type":
		case "source":
		case "id":
		case "Expires":
		case "Content-Type":
		case "publishdate":
		case "publishqueue":
		case "Publishtype":
		default:
			err := event.SetExtension(name, value)
			if err != nil {
				return nil, errors.New("could not get extension " + name + ": " + err.Error())
			}
		}
	}
	validationError := event.Validate()
	if validationError != nil {
		return nil, errors.New("could not validate data: " + validationError.Error())
	}
	eventValidation := validateEvent(&event)
	if eventValidation != nil {
		return nil, errors.New(`could not validate event: ` + eventValidation.Error())
	}
	return &event, nil
}

func validateEvent(event *data_types.CloudEvent) error {
	if queue_publisher.QueuePublisher().ValidateQueue(event.PublishQueue) != true {
		return errors.New(`could not validate publish queue`)
	}
	_, err := time.Parse("2006-01-02 15:04:05Z", event.PublishDate)
	if err != nil {
		return errors.New(`could not validate publish date`)
	}
	if event.PublishType != data_types.DataOnly {
		event.PublishType = data_types.EntireCloudEvent
	}
	return nil
}

func checkScheduler(event *data_types.ArangoCloudEvent) {
	s := scheduler.GetScheduler()
	s.CheckEvent(event)
}

func HTTPCreateEvent(context *gin.Context) {
	response := routes.JsendMessage{}
	switch context.Request.Header.Get("Content-Type") {
	case "application/cloudevents":
		event, err := ceHttpCreate(context)
		if err != nil {
			response = collection_managment.DefaultErrorHandler(err)
			context.JSON(int(response.Status()), &response)
			return
		}
		insertedItem, err := collection_managment.NewEventCollection().Insert(event)
		if err != nil {
			response = collection_managment.DefaultErrorHandler(err)
			context.JSON(int(response.Status()), &response)
			return
		}
		checkScheduler(insertedItem)
		response.SetStatus(http.StatusCreated)
		response.SetData(insertedItem.CloudEvent)
		context.JSON(int(response.Status()), &response)
		return
	case "application/json":
		event, err := jsonHttpCreate(context)
		if err != nil {
			response = collection_managment.DefaultErrorHandler(err)
			context.JSON(int(response.Status()), &response)
			return
		}
		insertedItem, err := collection_managment.NewEventCollection().Insert(event)
		if err != nil {
			response = collection_managment.DefaultErrorHandler(err)
			context.JSON(int(response.Status()), &response)
			return
		}
		checkScheduler(insertedItem)
		response.SetStatus(http.StatusCreated)
		response.SetData(insertedItem.CloudEvent)
		context.JSON(int(response.Status()), &response)
		return
	default:
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage("unknown Content-Type")
		return
	}
}

func HTTPDeleteEvent(context *gin.Context) {
	id := context.Param("event_id")
	response := routes.JsendMessage{}
	data, err := collection_managment.NewEventCollection().DeleteItem([]string{id})
	if err != nil {
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetEvent(context *gin.Context) {
	id := context.Param("event_id")
	response := routes.JsendMessage{}
	data, err := collection_managment.NewEventCollection().ReadItem(id)
	if err != nil {
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	if data == nil {
		response.SetStatus(http.StatusNotFound)
		response.SetMessage("event not found")
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetAllEvent(context *gin.Context) {
	response := routes.JsendMessage{}
	data, err := collection_managment.NewEventCollection().Read(bindQueryFilterParams(context))
	if err != nil {
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}
