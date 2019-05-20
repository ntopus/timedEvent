package event

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/ivanmeca/timedEvent/application/modules/database"
	"github.com/ivanmeca/timedEvent/application/modules/database/collection_managment"
	"github.com/ivanmeca/timedEvent/application/modules/database/data_types"
	"github.com/ivanmeca/timedEvent/application/modules/routes"
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
	event, err := bindEventInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	err = event.Context.SetID(bson.NewObjectId().String())
	if err != nil {
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	_, err = collection_managment.NewEventCollection().Insert(event)
	if err != nil {
		response = collection_managment.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(event)
	context.JSON(int(response.Status()), &response)
}

func HTTPDeleteEvent(context *gin.Context) {
	//id := context.Param("driver_id")
	//err := fleetDB.DeleteDriver(id)
	//response := routes.JsendMessage{}
	//response.SetStatus(http.StatusOK)
	//response.SetData("OK")
	//if err != nil {
	//	response = fleetDB.DefaultErrorHandler(err)
	//}
	//context.JSON(int(response.Status()), &response)
}

func HTTPUpdateEvent(context *gin.Context) {
	//id := context.Param("driver_id")
	//response := routes.JsendMessage{}
	//event, err := bindEventInformation(context)
	//if err != nil {
	//	response.SetStatus(http.StatusBadRequest)
	//	response.SetMessage(err.Error())
	//	context.JSON(int(response.Status()), &response)
	//	return
	//}
	//err = fleetDB.UpdateDriver(id, *event)
	//response.SetStatus(http.StatusOK)
	//response.SetData("OK")
	//if err != nil {
	//	response = fleetDB.DefaultErrorHandler(err)
	//	context.JSON(int(response.Status()), &response)
	//	return
	//}
	//context.JSON(int(response.Status()), &response)
}

func HTTPGetEvent(context *gin.Context) {
	//id := context.Param("driver_id")
	//event := fleetDB.GetDriverById(id)
	//response := routes.JsendMessage{}
	//response.SetStatus(http.StatusOK)
	//response.SetData(event)
	//context.JSON(int(response.Status()), &response)
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
