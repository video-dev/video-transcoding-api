package main

import (
	"io/ioutil"
	"log"

	"github.com/NYTimes/gizmo/server"
	"github.com/google/gops/agent"
	"github.com/video-dev/video-transcoding-api/v2/config"
	_ "github.com/video-dev/video-transcoding-api/v2/provider/bitmovin"
	_ "github.com/video-dev/video-transcoding-api/v2/provider/elastictranscoder"
	_ "github.com/video-dev/video-transcoding-api/v2/provider/elementalconductor"
	_ "github.com/video-dev/video-transcoding-api/v2/provider/encodingcom"
	_ "github.com/video-dev/video-transcoding-api/v2/provider/hybrik"
	_ "github.com/video-dev/video-transcoding-api/v2/provider/mediaconvert"
	_ "github.com/video-dev/video-transcoding-api/v2/provider/zencoder"
	"github.com/video-dev/video-transcoding-api/v2/service"
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
