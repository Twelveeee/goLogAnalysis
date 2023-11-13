package config

import (
	"github.com/twelveeee/log_analysis/service/eventLog"
	"github.com/twelveeee/log_analysis/service/util"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

var log = eventLog.NewLog()

type Config struct {
	// once   sync.Once
	cliCtx  *cli.Context
	options *Options
	db      *gorm.DB
	client  *Client
	env     string
}

// NewConfig initialises a new configuration
func NewConfig(ctx *cli.Context) *Config {
	configFilePath := ctx.String("config")

	c := &Config{
		cliCtx:  ctx,
		options: &Options{},
		client:  &Client{},
	}

	// Overwrite values with options.yml from config path.
	if util.FileExists(configFilePath) {
		if err := c.options.Load(configFilePath); err != nil {
			log.Err(err).Msgf("config: failed loading values from %s (%s)", configFilePath, err)
		} else {
			log.Debug().Msgf("config: overriding config with values from %s", configFilePath)
		}
	}

	return c
}

func (c *Config) Init() error {

	// Test Aliyun
	if err := c.initAliyunClient(); err != nil {
		return err
	}

	// Connect to database.
	if err := c.initDb(); err != nil {
		return err
	}
	return nil
}
