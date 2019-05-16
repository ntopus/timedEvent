package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ivanmeca/timedEvent/application/modules/authenticate"
	"github.com/ivanmeca/timedEvent/application/modules/routes/event"
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
		eventGroup := v1.Group("/driver")
		{
			eventGroup.GET("", event.HTTPGetAllEvent)
			eventGroup.POST("", event.HTTPCreateEvent)
			eventCRUD := eventGroup.Group("/:event_id")
			{
				eventCRUD.GET("", event.HTTPGetEvent)
				eventCRUD.PUT("", event.HTTPUpdateEvent)
				eventCRUD.DELETE("", event.HTTPDeleteEvent)
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
