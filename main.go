package main

import (
	"os"

	"github.com/twelveeee/log_analysis/command"
	"github.com/twelveeee/log_analysis/service/eventLog"
	"github.com/urfave/cli/v2"
)

var version = "v0.0.1 development"

var log = eventLog.NewLog()

const appName = "AliyunLogAnalysis"
const appAbout = "Twelveeee"

const appDescription = "日志分析"
const appCopyright = "(c) 2023 Twelveeee @ Twelveeee"

// Metadata contains build specific information.
var Metadata = map[string]interface{}{
	"Name":        appName,
	"About":       appAbout,
	"Description": appDescription,
	"Version":     version,
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal().Msgf("run recover %v", err)
			os.Exit(1)
		}
	}()

	eventLog.InitLog()

	app := cli.NewApp()
	app.Usage = appAbout
	app.Description = appDescription
	app.Version = version
	app.Copyright = appCopyright
	app.EnableBashCompletion = true
	app.Flags = command.Flags.Cli()
	app.Commands = command.Commands
	app.Metadata = Metadata
	// os.Args = append(os.Args, "start")

	log.Info().Msgf("start app os.Args: %v", os.Args[1:])
	if err := app.Run(os.Args); err != nil {
		log.Err(err).Msg("run error")
	}

	log.Info().Msg("task done")
}
