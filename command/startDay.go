package command

import (
	"fmt"
	"strconv"

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
		Name:     "startDay",
		Aliases:  []string{"s"},
		Usage:    "input `STARTDAY` format (YYYY-MM-DD)",
		Value:    "2023-01-01",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "endDay",
		Aliases:  []string{"e"},
		Usage:    "input `ENDDAY` format (YYYY-MM-DD)",
		Value:    "2023-01-01",
		Required: true,
	},
}

func startDayAction(ctx *cli.Context) error {

	log.Info().Msg(strconv.Itoa(ctx.NArg()))
	log.Info().Msg(ctx.Args().Get(0))
	log.Info().Msg(strconv.Itoa(ctx.Args().Len()))

	log.Info().Msg(ctx.String("startDay"))

	// if ctx.NArg() > 0 {
	// 	log.Info().Msg(ctx.Args().Get(0))
	// } else {
	// 	log.Info().Msg("asd")
	// }
	// log.Info().Msg(ctx.Args().Get(1))
	fmt.Println(ctx.Args().Tail())

	return nil
}
