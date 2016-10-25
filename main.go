package main

import (
	"io/ioutil"

	"github.com/NYTimes/gizmo/server"
	"github.com/knq/sdhook"
	"github.com/nytm/video-transcoding-api/config"
	_ "github.com/nytm/video-transcoding-api/provider/elastictranscoder"
	_ "github.com/nytm/video-transcoding-api/provider/elementalconductor"
	_ "github.com/nytm/video-transcoding-api/provider/encodingcom"
	"github.com/nytm/video-transcoding-api/service"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.Server.RouterType == "" {
		cfg.Server.RouterType = "fast"
	}
	server.Init("video-transcoding-api", cfg.Server)
	if cfg.GCPCredentials.String() != "" {
		gcpLoggingHook, err := sdhook.New(
			sdhook.GoogleServiceAccountCredentialsJSON([]byte(cfg.GCPCredentials.String())),
		)
		if err != nil {
			server.Log.Fatal("unable to initialize GCP logging hook: ", err)
		}
		server.Log.Hooks.Add(gcpLoggingHook)
		server.Log.Out = ioutil.Discard
	} else {
		server.Log.Debug("GCP credentials were not set")
	}
	service, err := service.NewTranscodingService(cfg, server.Log)
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
