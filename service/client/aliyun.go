package client

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/twelveeee/log_analysis/config"
	"github.com/twelveeee/log_analysis/service/entity"
	"github.com/twelveeee/log_analysis/service/util"
)

const aliyunOssRegex = `^ (\S+) - - \[(\S+ \+\d{4})\] "([^"]+)" (\d{3}) (\d+) (\d+) "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)"`

var aliyunOssRe = regexp.MustCompile(aliyunOssRegex)

type AliyunOssConfig struct {
	downloadPath string
	client       *oss.Client
	pathList     []config.AliyunBucketPath
}

type AliyunOssList struct {
	bucket       *oss.Bucket
	downloadPath string
	objectList   []oss.ObjectProperties
}

func StartAliyunOss(conf *config.Config) {
	aliConf := newAliyunOssConf(conf)
	aliyunOssList := newAliyunOssList(&aliConf)

	var wg sync.WaitGroup
	ossFilePathChannel := make(chan string, 10)

	// download aliyun oss file
	for _, aliyunOss := range aliyunOssList {
		wg.Add(len(aliyunOss.objectList))
		for _, object := range aliyunOss.objectList {
			go func(bucket *oss.Bucket, object oss.ObjectProperties, downloadPath string) {
				defer wg.Done()
				filePath := getFilePath(object.Key, downloadPath)
				if err := downLoadAliyunOss(bucket, object, filePath); err != nil {
					log.Err(err).Send()
					return
				}
				ossFilePathChannel <- filePath
			}(aliyunOss.bucket, object, aliyunOss.downloadPath)
		}
	}

	// download success close channel
	go func() {
		wg.Wait()
		close(ossFilePathChannel)
	}()

	// read oss file to db
	var fileWg sync.WaitGroup
	for filePath := range ossFilePathChannel {
		fileWg.Add(1)
		go func(filePath string) {
			defer fileWg.Done()
			scanAliyunOss(filePath)
		}(filePath)
	}

	fileWg.Wait()
}

func newAliyunOssConf(conf *config.Config) AliyunOssConfig {
	pathList := conf.GetAliyunOssBucketPathList()
	client := conf.GetAliyunClient()
	downLoadPath := conf.DownloadPath()

	return AliyunOssConfig{
		downloadPath: downLoadPath,
		client:       client,
		pathList:     pathList,
	}
}

func newAliyunOssList(conf *AliyunOssConfig) []AliyunOssList {
	ret := make([]AliyunOssList, 0)
	client := conf.client
	for _, bucketPath := range conf.pathList {
		bucket, err := client.Bucket(bucketPath.Bucket)
		if err != nil {
			log.Err(err).Send()
			continue
		}

		lsRes, err := bucket.ListObjects(oss.Prefix(bucketPath.Prefix))
		if err != nil {
			log.Err(err).Send()
			continue
		}

		ret = append(ret, AliyunOssList{
			bucket:       bucket,
			objectList:   lsRes.Objects,
			downloadPath: conf.downloadPath,
		})
	}

	return ret
}

func downLoadAliyunOss(bucket *oss.Bucket, object oss.ObjectProperties, filePath string) error {
	log.Debug().Msgf("download file to %s", filePath)

	if filePath == "" {
		return fmt.Errorf("download file to %s ,key %s", filePath, object.Key)
	}

	err := bucket.GetObjectToFile(object.Key, filePath)
	if err != nil {
		return err
	}

	return nil
}

func getFilePath(key, downloadPath string) string {
	lastSlashIndex := strings.LastIndex(key, "/")
	if lastSlashIndex == -1 {
		return ""
	}

	if len(downloadPath) > 1 && downloadPath[len(downloadPath)-1] != '/' {
		return downloadPath + "/" + key[lastSlashIndex+1:]
	}
	return downloadPath + key[lastSlashIndex+1:]
}

func scanAliyunOss(filePath string) {
	log.Debug().Msgf("read file %s", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Err(err).Send()
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ossList := make([]entity.OssLog, 0)
	for scanner.Scan() {
		line := scanner.Text()
		ossLog := aliyuLogToOssLog(line)
		if ossLog.IsFilter() {
			log.Info().Msgf("filter:oss filter line:%s", line)
			continue
		}
		ossList = append(ossList, ossLog)
	}

	ossLog := entity.OssLog{}
	rowsAffected, err := ossLog.CreateInBatches(ossList, 100)
	if err != nil {
		log.Err(err)
	}
	log.Debug().Msgf("read file save to db rowsAffected:%d , %s", rowsAffected, filePath)
}

func aliyuLogToOssLog(line string) entity.OssLog {
	matches := aliyunOssRe.FindStringSubmatch(line)
	if matches == nil {
		return entity.OssLog{}
	}

	return entity.OssLog{
		IP:          util.IpToInt(matches[1]),
		RequestTime: util.ConverToTimestamp(matches[2]),
		UA:          matches[8],
		Path:        getPath(matches[3]),
		Referer:     matches[9],
		Host:        getHost(matches[7]),
		Bucket:      matches[14],
	}
}

func getPath(path string) string {
	fields := strings.Split(path, " ")
	return fields[1]
}

func getHost(host string) string {
	if len(host) <= 1 {
		return host
	}

	u, err := url.Parse(host)
	if err != nil {
		log.Err(err).Msg("getHost Parse url fail")
		return host
	}

	return u.Scheme + "://" + u.Host
}
