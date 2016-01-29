package main

import (
	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/service"
)

func main() {
	cfg := config.LoadConfig("./config.json")

	server.Init("video-transcoding-api", cfg.Server)
	service, err := service.NewTranscodingService(cfg)
	if err != nil {
		server.Log.Fatal("unable to initialize service: ", err)
	}
	err = server.Register(service)
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}
	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
