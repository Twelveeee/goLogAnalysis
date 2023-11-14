package config

import (
	"fmt"
	"os"

	"github.com/twelveeee/log_analysis/service/util"
	"gopkg.in/yaml.v2"
)

type Options struct {
	File struct {
		DownLoadPath string `yaml:"DownLoadPath"`
	} `yaml:"File"`

	Database struct {
		Driver    string `yaml:"Driver"`
		Dsn       string `yaml:"Dsn"`
		DbName    string `yaml:"DbName"`
		Server    string `yaml:"Server"`
		User      string `yaml:"User"`
		Password  string `yaml:"Password"`
		Port      int    `yaml:"Port"`
		ConnsIdle int    `yaml:"ConnsIdle"`
		Conns     int    `yaml:"Conns"`
	} `yaml:"Database"`

	Aliyun struct {
		Endpoint          string             `yaml:"Endpoint"`
		AccessKeyId       string             `yaml:"AccessKeyId"`
		AccessKeySecret   string             `yaml:"AccessKeySecret"`
		OssBucketPathList []AliyunBucketPath `yaml:"OssBucketPath"`
	} `yaml:"Aliyun"`
}

type AliyunBucketPath struct {
	Bucket string `yaml:"Bucket"`
	Prefix string `yaml:"Prefix"`
}

// Load uses a yaml config file to initiate the configuration entity.
func (c *Options) Load(fileName string) error {
	if fileName == "" {
		return nil
	}

	if !util.FileExists(fileName) {
		return fmt.Errorf("%s not found", fileName)
	}

	yamlConfig, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(yamlConfig, c)
}
