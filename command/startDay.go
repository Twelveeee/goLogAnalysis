package command

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/twelveeee/log_analysis/config"
	"github.com/twelveeee/log_analysis/service/client"
	"github.com/urfave/cli/v2"
)

var StartDayCommand = cli.Command{
	Name:    "startDay",
	Aliases: []string{"upd"},
	Usage:   "Starts the server with start day and end day",
	Flags:   startDayFlags,
	Action:  startDayAction,
}

var startDayFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "startDay",
		Aliases: []string{"s"},
		Usage:   "input `STARTDAY` format (YYYYMMDD)",
		Value:   "",
	},
	&cli.IntFlag{
		Name:    "days",
		Aliases: []string{"d"},
		Usage:   "input `days` int ,from startDay to startDay+days",
		Value:   1,
	},
}

func startDayAction(ctx *cli.Context) error {

	log.Info().Msg(strconv.Itoa(ctx.NArg()))
	log.Info().Msg(ctx.Args().Get(0))
	log.Info().Msg(strconv.Itoa(ctx.Args().Len()))

	startDay := ctx.String("startDay")
	days := ctx.Int("days")

	var startTime time.Time
	if len(startDay) == 0 {
		return fmt.Errorf("startDay is empty")
		// startTime = time.Now()
	}

	startTime, _ = time.Parse("20060102", startDay)

	log.Info().Msg("startDay: " + startTime.Format("20060102"))
	log.Info().Msg("days: " + strconv.Itoa(days))

	conf, err := config.InitConfig(ctx)
	if err != nil {
		return err
	}

	// Pass this context down the chain.
	cctx, cancel := context.WithCancel(context.Background())

	go StartDay(cctx, conf, cancel, startTime, days)

	<-cctx.Done()

	log.Info().Msg("down...")
	cancel()

	return nil
}

func StartDay(ctx context.Context, conf *config.Config, cancel context.CancelFunc, startTime time.Time, days int) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal().Msgf("StartDay recover %v", err)
		}
	}()

	client.StartAliyunOssDay(conf, startTime, days)

	cancel()
}
