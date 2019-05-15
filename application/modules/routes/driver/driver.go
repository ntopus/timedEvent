package driver

import (
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB/data_types"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/routes"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func bindDriverInformation(context *gin.Context) (*data_types.Driver, error) {
	var driver data_types.Driver
	err := context.ShouldBind(&driver)
	if err != nil {
		return nil, err
	}
	return &driver, nil
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

func HTTPCreateDriver(context *gin.Context) {
	response := routes.JsendMessage{}
	driver, err := bindDriverInformation(context)
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

func HTTPDeleteDriver(context *gin.Context) {
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

func HTTPUpdateDriver(context *gin.Context) {
	id := context.Param("driver_id")
	response := routes.JsendMessage{}
	driver, err := bindDriverInformation(context)
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

func HTTPGetDriver(context *gin.Context) {
	id := context.Param("driver_id")
	driver := fleetDB.GetDriverById(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(driver)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetAllDrivers(context *gin.Context) {
	data := fleetDB.GetDriverCollection(bindQueryFilterParams(context))
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}
