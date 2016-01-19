package main

import (
	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/service"
)

func main() {
	var cfg *service.Config
	config.LoadJSONFile("./config.json", &cfg)
	config.LoadEnvConfig(&cfg)

	server.Init("video-transcoding-api", cfg.Server)
	err := server.Register(service.NewJSONService(cfg))
	if err != nil {
		server.Log.Fatal("unable to register service: ", err)
	}

	err = server.Run()
	if err != nil {
		server.Log.Fatal("server encountered a fatal error: ", err)
	}
}
