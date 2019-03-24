package main

import (
	"io/ioutil"
	"log"

	"github.com/nytimes/gizmo/server"
	"github.com/nytimes/video-transcoding-api/config"
	_ "github.com/nytimes/video-transcoding-api/provider/bitmovin"
	_ "github.com/nytimes/video-transcoding-api/provider/elastictranscoder"
	_ "github.com/nytimes/video-transcoding-api/provider/elementalconductor"
	_ "github.com/nytimes/video-transcoding-api/provider/encodingcom"
	_ "github.com/nytimes/video-transcoding-api/provider/hybrik"
	_ "github.com/nytimes/video-transcoding-api/provider/zencoder"
	"github.com/nytimes/video-transcoding-api/service"
	"github.com/google/gops/agent"
)

func main() {
	agent.Listen(agent.Options{})
	defer agent.Close()
	cfg := config.LoadConfig()
	server.Init("video-transcoding-api", cfg.Server)
	server.Log.Out = ioutil.Discard

	logger, err := cfg.Log.Logger()
	if err != nil {
		log.Fatal(err)
	}

	service, err := service.NewTranscodingService(cfg, logger)
	if err != nil {
		logger.Fatal("unable to initialize service: ", err)
	}
	err = server.Register(service)
	if err != nil {
		logger.Fatal("unable to register service: ", err)
	}
	err = server.Run()
	if err != nil {
		logger.Fatal("server encountered a fatal error: ", err)
	}
}
