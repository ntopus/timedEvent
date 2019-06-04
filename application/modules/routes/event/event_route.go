package event

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/routes"
	"io/ioutil"
	"net/http"
)

const (
	EVENT_WHEN  = "publishData"
	EVENT_WHERE = "publishQueue"
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
	for name, value := range headers {
		switch name {
		case "specversion":
			err := event.SetSpecVersion(value[0])
			if err != nil {
				return nil, errors.New("could not get spec version: " + err.Error())
			}
		case "type":
			err = event.SetType(value[0])
			if err != nil {
				return nil, errors.New("could not get type: " + err.Error())
			}
		case "Source":
			err = event.SetSource(value[0])
			if err != nil {
				return nil, errors.New("could not get source: " + err.Error())
			}
		case "id":
			err = event.SetID(value[0])
			if err != nil {
				return nil, errors.New("could not get id: " + err.Error())
			}
		case "Expires":
			time, err := data_types.GetTime(value[0])
			if err != nil {
				return nil, errors.New("could not get time: " + err.Error())
			}
			err = event.SetTime(*time)
			if err != nil {
				return nil, errors.New("could not get time: " + err.Error())
			}
		case "Content-Type":
			err = event.SetDataContentType(value[0])
			if err != nil {
				return nil, errors.New("could not get content type: " + err.Error())
			}
		default:
			err := event.SetExtension(name, value)
			if err != nil {
				return nil, errors.New("could not get extension " + name + ": " + err.Error())
			}
		}
	}
	eventValidation := validateEvent(&event)
	if eventValidation != nil {
		return nil, errors.New(`could not validate event: ` + eventValidation.Error())
	}
	return &event, nil
}

func validateEvent(event *data_types.CloudEvent) error {
	var PublishWhen string
	err := event.ExtensionAs(EVENT_WHEN, PublishWhen)
	if err != nil {
		return errors.New(`could not get publish date: ` + err.Error())
	}
	var PublishWhere string
	err = event.ExtensionAs(EVENT_WHERE, PublishWhere)
	if err != nil {
		return errors.New(`could not get publish queue: ` + err.Error())
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
