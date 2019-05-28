package event

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/routes"
	"io/ioutil"
	"net/http"
)

func bindEventInformation(context *gin.Context) (*data_types.CloudEvent, error) {
	var event data_types.CloudEvent
	err := context.ShouldBind(&event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

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
	event, err := bindEventInformation(context)
	if err != nil {
		return nil, errors.New("could not read request data: " + err.Error())
	}
	validationError := event.Validate()
	if validationError != nil {
		return nil, errors.New("could not read validate data: " + err.Error())
	}
	return event, nil
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
		case "schemaurl":
			err = event.SetSchemaURL(value[0])
			if err != nil {
				return nil, errors.New("could not get schemaURL: " + err.Error())
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
	return &event, nil
}

func HTTPCreateEvent(context *gin.Context) {
	response := routes.JsendMessage{}

	//switch context.Request.Header.Get("Content-Type") {
	//	case "application/cloudevents":
	//		event,err := ceHttpCreate(context)
	//	case "application/json":
	//		event,err := jsonHttpCreate(context)
	//	default:
	//
	//}

	//
	//["eventId"]=>
	//string(44) "timeout:location-expires:Tetra:724-1121:5000"
	//["eventSource"]=>
	//string(22) "Native.Location.Expire"
	//["eventDate"]=>
	//string(27) "2019-05-15T14:38:48.000000Z"
	//["eventPublishDate"]=>
	//string(27) "2019-05-15T14:38:52.000000Z"
	//["eventPublishQueue"]=>
	//string(22) "Timer.Resource.ThrowAt"

	insertedItem, err := collection_managment.NewEventCollection().Insert(event)
	if err != nil {
		response = collection_managment.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(insertedItem.CloudEvent)
	context.JSON(int(response.Status()), &response)
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
