package main

import (
	"flag"

	"github.com/NYTimes/gizmo/server"
	"github.com/nytm/video-transcoding-api/config"
	"github.com/nytm/video-transcoding-api/service"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "c", "", "path to the configuration file")
	flag.Parse()
}

func main() {
	cfg := config.LoadConfig(configFile)

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
