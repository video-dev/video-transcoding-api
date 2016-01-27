package main

import (
	gizmoConfig "github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/service"
)

func main() {
	var cfg config.Config
	gizmoConfig.LoadJSONFile("./config.json", &cfg)
	gizmoConfig.LoadEnvConfig(&cfg)

	server.Init("video-transcoding-api", cfg.Server)
	service, err := service.NewTranscodingService(&cfg)
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
