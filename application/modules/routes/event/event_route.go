package event

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
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

func bindQueryFilterParams(context *gin.Context) interface{} {
	filter := make(map[string]interface{})
	for i, value := range context.Request.URL.Query() {
		switch i {
		case "contract":
			filter["contract"] = map[string]interface{}{"$in": value}
		case "customer_code":
			filter["customercode"] = value[0]
		case "personaltag":
			filter["personaltag.number"] = value[0]
		case "category":
			filter["category.name"] = value[0]
		default:
			filter[i] = value[0]
		}
	}
	return filter
}

func HTTPCreateEvent(context *gin.Context) {
	response := routes.JsendMessage{}
	driver, err := bindEventInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	driver.Id = bson.NewObjectId()
	err = fleetDB.CreateDriver(*driver)
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(driver)
	context.JSON(int(response.Status()), &response)
}

func HTTPDeleteEvent(context *gin.Context) {
	id := context.Param("driver_id")
	err := fleetDB.DeleteDriver(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPUpdateEvent(context *gin.Context) {
	id := context.Param("driver_id")
	response := routes.JsendMessage{}
	driver, err := bindEventInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	err = fleetDB.UpdateDriver(id, *driver)
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPGetEvent(context *gin.Context) {
	id := context.Param("driver_id")
	driver := fleetDB.GetDriverById(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(driver)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetAllEvent(context *gin.Context) {
	data := fleetDB.GetDriverCollection(bindQueryFilterParams(context))
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}
