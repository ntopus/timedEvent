package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ivanmeca/timedEvent/application/modules/routes/event"
	"log"
	"net/http"
)

type HttpServer struct {
	engine *gin.Engine
	server *http.Server
}

func NewHttpServer(port string, debugMode bool) *HttpServer {
	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	engine.HandleMethodNotAllowed = true
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: engine,
	}
	httpServer := &HttpServer{
		server: srv,
		engine: engine,
	}
	return httpServer
}

func (httpServer *HttpServer) RunServer(ctx context.Context) error {
	v1 := httpServer.engine.Group("/v1")
	{
		eventGroup := v1.Group("/event")
		{
			eventGroup.GET("", event.HTTPGetAllEvent)
			eventGroup.POST("", event.HTTPCreateEvent)
			eventCRUD := eventGroup.Group("/:event_id")
			{
				eventCRUD.GET("", event.HTTPGetEvent)
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
