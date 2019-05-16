package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpServer struct {
	engine *gin.Engine
	server *http.Server
	auth   authenticate.IAuthenticate
}

type contractParser struct {
	Contract string `json:"contract"`
}

func NewHttpServer(port string, auth authenticate.IAuthenticate) *HttpServer {
	engine := gin.Default()
	engine.HandleMethodNotAllowed = true
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}
	httpServer := &HttpServer{
		server: srv,
		engine: engine,
		auth:   auth,
	}
	return httpServer
}

func (httpServer *HttpServer) GETApiHandler(c *gin.Context) {
	fmt.Println("GET MW!")
	token := c.Request.Header.Get("token")
	if token == "" {
		fmt.Println("Without token!")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
	requestContract := c.Request.URL.Query().Get("contract")
	if requestContract == "" {
		fmt.Println("Without contract!")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Abort()
		return
	}
	ok, err := httpServer.auth.CheckTokenPermission(token, requestContract)
	if err != nil {
		fmt.Println("Token wrong!")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
	if !ok {
		fmt.Println("Token unauthorized!")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
}

func (httpServer *HttpServer) POSTApiHandler(c *gin.Context) {
	fmt.Println("POST MW!")
	token := c.Request.Header.Get("token")
	if token == "" {
		fmt.Println("Without token!")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
	buf, _ := ioutil.ReadAll(c.Request.Body)
	pretempBody := ioutil.NopCloser(bytes.NewBuffer(buf))
	postempBody := ioutil.NopCloser(bytes.NewBuffer(buf))
	c.Request.Body = pretempBody
	var Contract contractParser
	err := c.ShouldBind(&Contract)
	if err != nil {
		fmt.Println("Without contract!")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Abort()
		return
	}
	c.Request.Body = postempBody
	ok, err := httpServer.auth.CheckTokenPermission(token, Contract.Contract)
	if err != nil {
		fmt.Println("Token wrong!")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
	if !ok {
		fmt.Println("Token unauthorized!")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Abort()
		return
	}
}

func (httpServer *HttpServer) MiddlewareApiHandler(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		httpServer.GETApiHandler(c)
	case "DELETE":
		httpServer.GETApiHandler(c)
	default:
		httpServer.POSTApiHandler(c)
	}
}

func (httpServer *HttpServer) RunServer(ctx context.Context) error {
	v1 := httpServer.engine.Group("/v1")
	{
		v1.Use(httpServer.MiddlewareApiHandler)
		driverGroup := v1.Group("/driver")
		{
			driverGroup.GET("", driver.HTTPGetAllDrivers)
			driverGroup.POST("", driver.HTTPCreateDriver)
			driverCRUD := driverGroup.Group("/:driver_id")
			{
				driverCRUD.GET("", driver.HTTPGetDriver)
				driverCRUD.PUT("", driver.HTTPUpdateDriver)
				driverCRUD.DELETE("", driver.HTTPDeleteDriver)
			}
		}
		categoryGroup := v1.Group("/driver_category")
		{
			categoryGroup.GET("", driver.HTTPGetAllDriverCategories)
			categoryGroup.POST("", driver.HTTPCreateDriverCategory)
			categoryCRUD := categoryGroup.Group("/:category_id")
			{
				categoryCRUD.GET("", driver.HTTPGetDriverCategory)
				categoryCRUD.PUT("", driver.HTTPUpdateDriverCategory)
				categoryCRUD.DELETE("", driver.HTTPDeleteDriverCategory)
			}
		}
		vehicleGroup := v1.Group("/vehicle")
		{
			vehicleGroup.GET("", vehicle.HTTPGetAllVehicles)
			vehicleGroup.POST("", vehicle.HTTPCreateVehicle)
			vehicleCRUD := vehicleGroup.Group("/:vehicle_id")
			{
				vehicleCRUD.GET("", vehicle.HTTPGetVehicle)
				vehicleCRUD.PUT("", vehicle.HTTPUpdateVehicle)
				vehicleCRUD.DELETE("", vehicle.HTTPDeleteVehicle)
			}
		}
		gsmSimGroup := v1.Group("/gsm_sim")
		{
			gsmSimGroup.GET("", gsm_sim.HTTPGetAllGsmSIM)
			gsmSimGroup.POST("", gsm_sim.HTTPCreateGsmSIM)
			gsmSimCRUD := gsmSimGroup.Group("/:gsmsim_id")
			{
				gsmSimCRUD.GET("", gsm_sim.HTTPGetGsmSIM)
				gsmSimCRUD.PUT("", gsm_sim.HTTPUpdateGsmSIM)
				gsmSimCRUD.DELETE("", gsm_sim.HTTPDeleteGsmSIM)
			}
		}
		trackingDeviceGroup := v1.Group("/tracking_device")
		{
			trackingDeviceGroup.GET("", tracking_device.HTTPGetAllTrackingDevice)
			trackingDeviceGroup.POST("", tracking_device.HTTPCreateTrackingDevice)
			trackingDeviceCRUD := trackingDeviceGroup.Group("/:device_id")
			{
				trackingDeviceCRUD.GET("", tracking_device.HTTPGetTrackingDevice)
				trackingDeviceCRUD.PUT("", tracking_device.HTTPUpdateTrackingDevice)
				trackingDeviceCRUD.DELETE("", tracking_device.HTTPDeleteTrackingDevice)
			}
		}
		refuelGroup := v1.Group("/refuel")
		{
			refuelGroup.GET("", refuel.HTTPGetAllRefuel)
			refuelGroup.POST("", refuel.HTTPCreateRefuel)
			refuelCRUD := refuelGroup.Group("/:refuel_id")
			{
				refuelCRUD.GET("", refuel.HTTPGetRefuel)
				refuelCRUD.PUT("", refuel.HTTPUpdateRefuel)
				refuelCRUD.DELETE("", refuel.HTTPDeleteRefuel)
			}
		}
	}
	go func() {
		if err := httpServer.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()
	go func() {
		<-ctx.Done()
		err := httpServer.server.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()
	return nil
}
