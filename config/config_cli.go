package config

import (
	"github.com/urfave/cli/v2"
)

var InitConfig = func(ctx *cli.Context) (*Config, error) {
	c := NewConfig(ctx)
	// get.SetConfig(c)
	return c, c.Init()
}
