package command

import "github.com/urfave/cli/v2"

type CliFlag struct {
	Flag cli.DocGenerationFlag
	Tags []string
}

type CliFlags []CliFlag

// Cli returns the currently active command-line parameters.
func (f CliFlags) Cli() (result []cli.Flag) {
	result = make([]cli.Flag, 0, len(f))

	for _, flag := range f {
		result = append(result, flag.Flag)
	}

	return result
}

var Flags = CliFlags{
	{
		Flag: &cli.StringFlag{
			Name:    "help2, ",
			Aliases: []string{"h2"},
			Usage:   "help",
		},
	},
	{
		Flag: &cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "config file path",
			Value:   "config/config.yaml",
		},
	},
}
