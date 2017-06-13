package main

import (
	"github.com/Gurpartap/logrus-stack"
	"github.com/NYTimes/gizmo/server"
	"github.com/NYTimes/video-transcoding-api/config"
	_ "github.com/NYTimes/video-transcoding-api/provider/bitmovin"
	_ "github.com/NYTimes/video-transcoding-api/provider/elastictranscoder"
	_ "github.com/NYTimes/video-transcoding-api/provider/elementalconductor"
	_ "github.com/NYTimes/video-transcoding-api/provider/encodingcom"
	_ "github.com/NYTimes/video-transcoding-api/provider/hybrik"
	_ "github.com/NYTimes/video-transcoding-api/provider/zencoder"
	"github.com/NYTimes/video-transcoding-api/service"
	"github.com/google/gops/agent"
	"github.com/knq/sdhook"
	"github.com/marzagao/logrus-env"
)

func main() {
	agent.Listen(&agent.Options{NoShutdownCleanup: true})
	defer agent.Close()
	cfg := config.LoadConfig()
	if cfg.Server.RouterType == "" {
		cfg.Server.RouterType = "fast"
	}
	server.Init("video-transcoding-api", cfg.Server)
	server.Log.Hooks.Add(logrus_stack.StandardHook())
	server.Log.Hooks.Add(logrus_env.NewHook([]string{"ENVIRONMENT"}))

	gcpLoggingHook, err := sdhook.New(
		sdhook.GoogleLoggingAgent(),
		sdhook.ErrorReportingService("video-transcoding-api"),
		sdhook.ErrorReportingLogName("error_log"),
	)
	if err != nil {
		server.Log.Debug("unable to initialize GCP logging hook: ", err)
	} else {
		server.Log.Hooks.Add(gcpLoggingHook)
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
