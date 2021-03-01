package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/SotirisAlfonsos/chaos-bot/config"
	"github.com/SotirisAlfonsos/chaos-bot/web"
	"github.com/SotirisAlfonsos/chaos-master/chaoslogger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {
	debugLevel := flag.String(
		"debug.level",
		"info",
		"the debug level for the chaos bot. Can be one of debug, info, warn, error.")

	port := flag.String(
		"port",
		"8081",
		"the port used by the grpc server.")

	configFile := flag.String(
		"config.file",
		"",
		"the file that contains the configuration for the chaos master")

	flag.Parse()

	logger := createLogger(*debugLevel)

	conf, err := config.GetConfig(*configFile)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	grpcHandler, err := web.NewGRPCHandler(*port, logger, conf)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}

	if err := grpcHandler.Run(); err != nil {
		_ = level.Error(logger).Log("msg", "Failed to start Grpc server on port "+*port, "err", err)
	}
}

func createLogger(debugLevel string) log.Logger {
	allowLevel := &chaoslogger.AllowedLevel{}
	if err := allowLevel.Set(debugLevel); err != nil {
		fmt.Printf("%v", err)
	}

	return chaoslogger.New(allowLevel)
}
