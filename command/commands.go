package command

import (
	"github.com/twelveeee/log_analysis/service/eventLog"
	"github.com/urfave/cli/v2"
)

var log = eventLog.NewLog()

var Commands = []*cli.Command{
	&StartCommand,
	&StartDayCommand,
}
