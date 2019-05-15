package gsm_sim

import (
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB/data_types"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/routes"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func bindGsmSimInformation(context *gin.Context) (*data_types.GsmSIM, error) {
	var gsmSim data_types.GsmSIM
	err := context.ShouldBind(&gsmSim)
	if err != nil {
		return nil, err
	}
	return &gsmSim, nil
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

func HTTPCreateGsmSIM(context *gin.Context) {
	response := routes.JsendMessage{}
	gsmSim, err := bindGsmSimInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	gsmSim.Id = bson.NewObjectId()
	err = fleetDB.CreateGsmSim(*gsmSim)
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(gsmSim)
	context.JSON(int(response.Status()), &response)
}

func HTTPDeleteGsmSIM(context *gin.Context) {
	id := context.Param("gsmsim_id")
	err := fleetDB.DeleteGsmSim(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPUpdateGsmSIM(context *gin.Context) {
	id := context.Param("gsmsim_id")
	response := routes.JsendMessage{}
	gsmSim, err := bindGsmSimInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	err = fleetDB.UpdateGsmSim(id, *gsmSim)
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPGetGsmSIM(context *gin.Context) {
	id := context.Param("gsmsim_id")
	gsmSim := fleetDB.GetGsmSimById(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(gsmSim)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetAllGsmSIM(context *gin.Context) {
	data := fleetDB.GetGsmSimCollection(bindQueryFilterParams(context))
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}
