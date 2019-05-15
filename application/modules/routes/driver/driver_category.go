package driver

import (
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/fleetDB/data_types"
	"devgit.kf.com.br/comercial/fleet-management-api/application/modules/routes"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
)

func bindDriverCategoryInformation(context *gin.Context) (*data_types.DriverCategory, error) {
	var category data_types.DriverCategory
	err := context.ShouldBind(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func bindDriverCategoryQueryFilterParams(context *gin.Context) interface{} {
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

func HTTPCreateDriverCategory(context *gin.Context) {
	response := routes.JsendMessage{}
	category, err := bindDriverCategoryInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	category.Id = bson.NewObjectId()
	err = fleetDB.CreateDriverCategory(*category)
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	response.SetStatus(http.StatusCreated)
	response.SetData(category)
	context.JSON(int(response.Status()), &response)
}

func HTTPDeleteDriverCategory(context *gin.Context) {
	id := context.Param("category_id")
	err := fleetDB.DeleteDriverCategory(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPUpdateDriverCategory(context *gin.Context) {
	id := context.Param("category_id")
	response := routes.JsendMessage{}
	category, err := bindDriverCategoryInformation(context)
	if err != nil {
		response.SetStatus(http.StatusBadRequest)
		response.SetMessage(err.Error())
		context.JSON(int(response.Status()), &response)
		return
	}
	err = fleetDB.UpdateDriverCategory(id, *category)
	response.SetStatus(http.StatusOK)
	response.SetData("OK")
	if err != nil {
		response = fleetDB.DefaultErrorHandler(err)
		context.JSON(int(response.Status()), &response)
		return
	}
	context.JSON(int(response.Status()), &response)
}

func HTTPGetDriverCategory(context *gin.Context) {
	id := context.Param("category_id")
	driver := fleetDB.GetDriverCategoryById(id)
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(driver)
	context.JSON(int(response.Status()), &response)
}

func HTTPGetAllDriverCategories(context *gin.Context) {
	data := fleetDB.GetDriverCategoryCollection(bindDriverCategoryQueryFilterParams(context))
	response := routes.JsendMessage{}
	response.SetStatus(http.StatusOK)
	response.SetData(data)
	context.JSON(int(response.Status()), &response)
}
