package event

import (
	"encoding/json"
	"errors"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/gin-gonic/gin"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/routes"
	"io/ioutil"
	"net/http"
)

const (
	F_ID            = 1
	F_SPEC_VERSION  = 2
	F_SOURCE        = 4
	F_PUBLISH_DATE  = 8
	F_PUBLISH_QUEUE = 16
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
	event.Data = data
	headers := context.Request.Header

	if value, ok := headers["specversion"]; ok {
		err := event.SetSpecVersion(value[0])
		if err != nil {
			return nil, errors.New("could not get spec version: " + err.Error())
		}
	} else {
		event.SetSpecVersion(cloudevents.CloudEventsVersionV02)
	}
	if value, ok := headers["type"]; ok {
		err = event.SetType(value[0])
		if err != nil {
			return nil, errors.New("could not get type: " + err.Error())
		}
	}
	if value, ok := headers["source"]; ok {
		err = event.SetSource(value[0])
		if err != nil {
			return nil, errors.New("could not get source: " + err.Error())
		}
	} else {
		event.SetSource(context.ClientIP())
	}
	if value, ok := headers["id"]; ok {
		err = event.SetID(value[0])
		if err != nil {
			return nil, errors.New("could not get id: " + err.Error())
		}
	}
	if value, ok := headers["Content-Type"]; ok {
		err = event.SetDataContentType(value[0])
		if err != nil {
			return nil, errors.New("could not get content type: " + err.Error())
		}
	}
	if value, ok := headers["publishDate"]; ok {
		time, err := data_types.GetTime(value[0])
		if err != nil {
			return nil, errors.New("could not get time: " + err.Error())
		}
		event.PublishDate = time.String()
	}
	if value, ok := headers["publishQueue"]; ok {
		event.PublishQueue = value[0]
	}

	for name, value := range headers {
		switch name {
		case "specversion":
		case "type":
		case "source":
		case "id":
		case "Expires":
		case "Content-Type":
		case "publishDate":
		case "publishQueue":
		default:
			err := event.SetExtension(name, value)
			if err != nil {
				return nil, errors.New("could not get extension " + name + ": " + err.Error())
			}
		}
	}
	validationError := event.Validate()
	if validationError != nil {
		return nil, errors.New("could not read validate data: " + validationError.Error())
	}
	//eventValidation := validateEvent(&event)
	//if eventValidation != nil {
	//	return nil, errors.New(`could not validate event: ` + eventValidation.Error())
	//}
	return &event, nil
}

func validateEvent(event *data_types.CloudEvent) error {
	if event.PublishQueue == "" {
		return errors.New(`could not validate publish queue`)
	}
	return nil
}

func HTTPCreateEvent(context *gin.Context) {
	response := routes.JsendMessage{}
	switch context.Request.Header.Get("Content-Type") {
	case "application/cloudevents":
		event, err := ceHttpCreate(context)
		if err != nil {
			response = collection_managment.DefaultErrorHandler(err)

		}
		insertedItem, err := collection_managment.NewEventCollection().Insert(event)
		if err != nil {
			response = collection_managment.DefaultErrorHandler(err)
			context.JSON(int(response.Status()), &response)
			return
		}
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
