package main

import (
	"flag"

	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	_ "github.com/nytm/video-transcoding-api/provider/elastictranscoder"
	_ "github.com/nytm/video-transcoding-api/provider/elementalconductor"
	_ "github.com/nytm/video-transcoding-api/provider/encodingcom"
	"github.com/nytm/video-transcoding-api/service"
)

func main() {
	flag.Parse()
	cfg := config.LoadConfig()
	if cfg.Server.RouterType == "" {
		cfg.Server.RouterType = "fast"
	}

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
