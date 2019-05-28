package event

import (
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

func HTTPCreateEvent(context *gin.Context) {
	response := routes.JsendMessage{}
	data, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		response.SetMessage("could not read request data: " + err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	event := data_types.CloudEvent{}
	event.Data = data
	headers := context.Request.Header
	for name, value := range headers {
		switch name {
		case "specversion":
			err := event.Context.SetSpecVersion(value[0])
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
		case "type":
			err = event.Context.SetType(value[0])
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
		case "Source":
			err = event.Context.SetSource(value[0])
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
		case "id":
			err = event.Context.SetID(value[0])
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
		case "Expires":
			time, err := data_types.GetTime(value[0])
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
			err = event.Context.SetTime(*time)
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
		case "schemaurl":
			err = event.Context.SetSchemaURL(value[0])
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
		case "Content-Type":
			err = event.Context.SetDataContentType(value[0])
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
		default:
			err := event.Context.SetExtension(name, value)
			if err != nil {
				response.SetMessage(err.Error())
				context.JSON(int(response.Status()), &response)
				return
			}
		}
	}
	if event.Context.GetID() == "" {
		err = event.Context.SetID(bson.NewObjectId().String())
		if err != nil {
			response.SetMessage(err.Error())
			context.JSON(int(response.Status()), &response)
			return
		}
	}
	validationError := event.Validate()
	if validationError != nil {
		response = collection_managment.DefaultErrorHandler(validationError)
		context.JSON(int(response.Status()), &response)
		return
	}
	insertedItem, err := collection_managment.NewEventCollection().Insert(&event)
	if err != nil {
		response = collection_managment.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(insertedItem)
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
