package refuel

import (
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB/data_types"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/routes"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func bindRefuelInformation(context *gin.Context) (*data_types.RefuelRegister, error) {
	var refuel data_types.RefuelRegister
	err := context.ShouldBind(&refuel)
	if err != nil {
		return nil, err
	}
	return &refuel, nil
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

func HTTPCreateRefuel(context *gin.Context) {
	response := routes.JsendMessage{}
	refuel, err := bindRefuelInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	refuel.Id = bson.NewObjectId()
	err = fleetDB.CreateRefuel(*refuel)
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(refuel)
	context.JSON(int(response.Status()), &response)
}

func HTTPDeleteRefuel(context *gin.Context) {
	id := context.Param("refuel_id")
	err := fleetDB.DeleteRefuel(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPUpdateRefuel(context *gin.Context) {
	id := context.Param("refuel_id")
	response := routes.JsendMessage{}
	refuel, err := bindRefuelInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	err = fleetDB.UpdateRefuel(id, *refuel)
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPGetRefuel(context *gin.Context) {
	id := context.Param("refuel_id")
	refuel := fleetDB.GetRefuelById(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(refuel)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetAllRefuel(context *gin.Context) {
	data := fleetDB.GetRefuelCollection(bindQueryFilterParams(context))
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}
