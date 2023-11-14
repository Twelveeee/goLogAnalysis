package config

import (
	"errors"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var lockAliyun = sync.Mutex{}

func (c *Config) initAliyunClient() error {
	lockAliyun.Lock()
	defer lockAliyun.Unlock()

	endpoint := c.options.Aliyun.Endpoint
	accessKeyId := c.options.Aliyun.AccessKeyId
	accessKeySecret := c.options.Aliyun.AccessKeySecret

	if endpoint == "" {
		return errors.New("config: aliyun endpoint not specified")
	}

	if accessKeyId == "" {
		return errors.New("config: aliyun accessKeyId not specified")
	}

	if accessKeySecret == "" {
		return errors.New("config: aliyun accessKeySecret not specified")
	}

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	if err = c.checkAliyunClient(client); err != nil {
		return err
	}

	c.client.AliyunClient = client

	return nil
}

func (c *Config) checkAliyunClient(client *oss.Client) error {
	_, err := client.ListBuckets()
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) GetAliyunClient() *oss.Client {
	return c.client.AliyunClient
}

func (c *Config) GetAliyunOssBucketPathList() []AliyunBucketPath {
	return c.options.Aliyun.OssBucketPathList
}
