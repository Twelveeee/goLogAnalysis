package command

import (
	"context"

	"github.com/twelveeee/log_analysis/config"
	"github.com/twelveeee/log_analysis/service/client"
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
	conf, err := config.InitConfig(ctx)
	if err != nil {
		return err
	}

	// Pass this context down the chain.
	cctx, cancel := context.WithCancel(context.Background())

	go Start(cctx, conf, cancel)

	<-cctx.Done()

	log.Info().Msg("down...")
	cancel()

	return nil
}

func Start(ctx context.Context, conf *config.Config, cancel context.CancelFunc) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal().Msgf("Start recover %v", err)
		}
	}()

	client.StartAliyunOss(conf)

	cancel()
}
