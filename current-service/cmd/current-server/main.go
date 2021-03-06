// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: d5b3153b9f
// Version Date: Thu Jul 27 18:20:46 UTC 2017

package main

import (
	"flag"
	"os"

	// Go Kit
	"github.com/go-kit/kit/log"

	// This Service
	"github.com/zaquestion/current/current-service/svc/server"
	"github.com/zaquestion/current/current-service/svc/server/cli"
)

func main() {
	// Update addresses if they have been overwritten by flags
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}
	server.Run(cli.Config, logger)
}
