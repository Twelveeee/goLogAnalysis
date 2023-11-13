package command

import (
	"github.com/twelveeee/log_analysis/config"
	"github.com/urfave/cli/v2"
)

var StartCommand = cli.Command{
	Name:    "start",
	Aliases: []string{"up"},
	Usage:   "Starts the server",
	Flags:   startFlags,
	Action:  startAction,
}

var startFlags = []cli.Flag{
	// &cli.BoolFlag{
	// 	Name:    "defaultDay",
	// 	Aliases: []string{"d"},
	// 	Value:   true,
	// 	Usage:   "is default day?",
	// },
}

func startAction(ctx *cli.Context) error {
	_, err := config.InitConfig(ctx)
	if err != nil {
		return err
	}

	return nil
}
