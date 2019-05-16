package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ivanmeca/timedEvent/application/modules/authenticate"
	"github.com/ivanmeca/timedEvent/application/modules/routes/event"
	"log"
	"net/http"
)

type HttpServer struct {
	engine *gin.Engine
	server *http.Server
	auth   authenticate.IAuthenticate
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

func (httpServer *HttpServer) RunServer(ctx context.Context) error {
	v1 := httpServer.engine.Group("/v1")
	{
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
