package client

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

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
	aliyunOssList := newAliyunOssList(&aliConf, time.Now(), 1)
	saveAliyunOss(aliyunOssList)
}

func StartAliyunOssDay(conf *config.Config, startTime time.Time, days int) {
	aliConf := newAliyunOssConf(conf)
	log.Debug().Msgf("startTime:%s days:%d", startTime.Format("20060102"), days)
	aliyunOssList := newAliyunOssList(&aliConf, startTime, days)
	saveAliyunOss(aliyunOssList)
}

func saveAliyunOss(aliyunOssList []AliyunOssList) {
	ossFilePathChannel := productFilePath(aliyunOssList)
	// ossFilePathChannel := productFilePathByStorage()
	scanAliyunOssFile(ossFilePathChannel)
}

func productFilePath(aliyunOssList []AliyunOssList) chan string {
	ossFilePathChannel := make(chan string, 10)
	go func() {
		var wg sync.WaitGroup
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

		go func() {
			wg.Wait()
			close(ossFilePathChannel)
		}()
	}()

	return ossFilePathChannel
}

func productFilePathByStorage(storagePath string) chan string {
	ossFilePathChannel := make(chan string, 10)

	var wg sync.WaitGroup
	err := filepath.Walk(storagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		wg.Add(1)
		go func() {
			// 忽略文件夹，只处理文件
			ossFilePathChannel <- path
			wg.Done()
			fmt.Println(path) // 打印文件路径
		}()

		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}

	go func() {
		wg.Wait()
		close(ossFilePathChannel)
	}()

	return ossFilePathChannel
}

func scanAliyunOssFile(ossFilePathChannel chan string) {
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

func newAliyunOssList(conf *AliyunOssConfig, startDay time.Time, days int) []AliyunOssList {
	ret := make([]AliyunOssList, 0)
	client := conf.client
	for _, bucketPath := range conf.pathList {
		bucket, err := client.Bucket(bucketPath.Bucket)
		if err != nil {
			log.Err(err).Send()
			continue
		}

		for i := 0; i < days; i++ {
			prefix := fmt.Sprintf("%s%s", bucketPath.Prefix, startDay.AddDate(0, 0, i).Format("2006-01-02"))
			log.Debug().Msgf("bucket: %s prefix: %s", bucketPath.Bucket, prefix)
			lsRes, err := bucket.ListObjects(oss.Prefix(prefix))
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
		if reason, ok := ossLog.IsFilter(); ok {
			if reason == "aliyun-sdk" {
				continue
			}

			log.Info().Msgf("filter:oss filter reason:%s line :%s", reason, line)
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
