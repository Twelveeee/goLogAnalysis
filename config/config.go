package config

import (
	"fmt"
	"os"

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
	// env     string
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

	// Test DownloadPath
	if err := c.CheckDownloadPath(); err != nil {
		return err
	}

	// Test Aliyun
	if err := c.initAliyunClient(); err != nil {
		return err
	}

	// Connect to database.
	if err := c.initDb(); err != nil {
		return err
	}

	// Register Db
	c.RegisterDb()

	return nil
}

func (c *Config) CheckDownloadPath() error {
	// 检查文件夹是否存在
	path := c.DownloadPath()
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 文件夹不存在，创建文件夹
		err := os.Mkdir(path, 0755)
		if err != nil {
			return fmt.Errorf("mkdir fail：%s", err)
		}
	}
	return nil
}

func (c *Config) DownloadPath() string {
	return c.options.File.DownLoadPath
}
