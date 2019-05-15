package tracking_device

import (
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB/data_types"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/routes"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func bindTrackingDeviceInformation(context *gin.Context) (*data_types.GsmTrackingDevice, error) {
	var trackingDevice data_types.GsmTrackingDevice
	err := context.ShouldBind(&trackingDevice)
	if err != nil {
		return nil, err
	}
	return &trackingDevice, nil
}

func bindQueryFilterParams(context *gin.Context) interface{} {
	filter := make(map[string]interface{})
	for i, value := range context.Request.URL.Query() {
		switch i {
		case "contract":
			filter["contract"] = map[string]interface{}{"$in": value}
		default:
			filter[i] = value[0]
		}
	}
	return filter
}

func HTTPCreateTrackingDevice(context *gin.Context) {
	response := routes.JsendMessage{}
	trackingDevice, err := bindTrackingDeviceInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	trackingDevice.Id = bson.NewObjectId()
	err = fleetDB.CreateTrackingDevice(*trackingDevice)
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(trackingDevice)
	context.JSON(int(response.Status()), &response)
}

func HTTPDeleteTrackingDevice(context *gin.Context) {
	id := context.Param("device_id")
	err := fleetDB.DeleteTrackingDevice(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPUpdateTrackingDevice(context *gin.Context) {
	id := context.Param("device_id")
	response := routes.JsendMessage{}
	trackingDevice, err := bindTrackingDeviceInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	err = fleetDB.UpdateTrackingDevice(id, *trackingDevice)
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPGetTrackingDevice(context *gin.Context) {
	id := context.Param("device_id")
	trackingDevice := fleetDB.GetTrackingDeviceById(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(trackingDevice)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetAllTrackingDevice(context *gin.Context) {
	data := fleetDB.GetTrackingDeviceCollection(bindQueryFilterParams(context))
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}
