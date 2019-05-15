package vehicle

import (
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB/data_types"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/routes"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func bindVehicleInformation(context *gin.Context) (*data_types.Vehicle, error) {
	var vehicle data_types.Vehicle
	err := context.ShouldBind(&vehicle)
	if err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func bindQueryFilterParams(context *gin.Context) interface{} {
	filter := make(map[string]interface{})
	for i, value := range context.Request.URL.Query() {
		switch i {
		case "contract":
			filter["contract"] = map[string]interface{}{"$in": value}
		case "networkdeviceid":
			filter["trackingdevice.networkdeviceid"] = value[0]
		case "gsmnumber":
			filter["trackingdevice.gsmsim.gsmnumber"] = value[0]
		default:
			filter[i] = value[0]
		}
	}
	return filter
}

func HTTPCreateVehicle(context *gin.Context) {
	response := routes.JsendMessage{}
	vehicle, err := bindVehicleInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	vehicle.Id = bson.NewObjectId()
	err = fleetDB.CreateVehicle(*vehicle)
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(vehicle)
	context.JSON(int(response.Status()), &response)
}

func HTTPDeleteVehicle(context *gin.Context) {
	id := context.Param("vehicle_id")
	err := fleetDB.DeleteVehicle(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPUpdateVehicle(context *gin.Context) {
	id := context.Param("vehicle_id")
	response := routes.JsendMessage{}
	vehicle, err := bindVehicleInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	err = fleetDB.UpdateVehicle(id, *vehicle)
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPGetVehicle(context *gin.Context) {
	id := context.Param("vehicle_id")
	vehicle := fleetDB.GetVehicleById(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(vehicle)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetAllVehicles(context *gin.Context) {
	data := fleetDB.GetVehicleCollection(bindQueryFilterParams(context))
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}
