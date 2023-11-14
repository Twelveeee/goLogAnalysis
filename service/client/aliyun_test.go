package client

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/yaml.v2"
)

type testAliyunConfig struct {
	Aliyun struct {
		Endpoint        string `yaml:"Endpoint"`
		AccessKeyId     string `yaml:"AccessKeyId"`
		AccessKeySecret string `yaml:"AccessKeySecret"`
	} `yaml:"Aliyun"`
}

func newTestAliyunClient() *oss.Client {
	config := testAliyunConfig{}
	yamlConfig, err := os.ReadFile("../../config/config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlConfig, &config)
	if err != nil {
		panic(err)
	}

	client, err := oss.New(config.Aliyun.Endpoint, config.Aliyun.AccessKeyId, config.Aliyun.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	return client

}

func TestAliyunClient(t *testing.T) {
	client := newTestAliyunClient()
	_, err := client.ListBuckets()
	if err != nil {
		panic(err)
	}
}

func TestAliyunOssRegex(t *testing.T) {
	logs := []string{
		` 1.1.1.1 - - [11/Nov/2023:23:31:40 +0800] "GET /Image/img/bgimg_31.webp HTTP/1.1" 200 166114 148 "https://www.example.com/" "Mozilla/5.0 (iPhone; CPU iPhone OS 16_5_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/8.0.38(0x1800262b) NetType/WIFI Language/zh_CN qcloudcdn-xinan Request-Source=4 Request-Channel=99" "example.oss-cn-beijing.aliyuncs.com" "654F9E5C35EB26393996A9A1" "true" "-" "GetObject" "exampleBucket" "bgimg_31.webp" 166114 31 "-" 578 "1165947387648789" - "-" "standard" "-" "-" "-"`,
		` 2.2.2.2 - - [11/Nov/2023:00:15:17 +0800] "GET /?bucketInfo HTTP/1.1" 200 1129 126 "-" "aliyun-sdk-java/1.1.8.3(Linux/4.19.91-007.ali4000.alios7.x86_64/amd64;1.8.0_332)" "example.oss-cn-shanghai-internal.aliyuncs.com" "654E5715DE97D439332E01EA" "true" "300024016803281053" "GetBucketInfo" "exampleBucket" "-" - - "-" 949 "1165947387648789" - "-" "standard" "-" "-" "STS.NTj3wVJJGW9NHzrue798inyDt"`,
	}

	pattern := `^ (\S+) - - \[(\S+ \+\d{4})\] "([^"]+)" (\d{3}) (\d+) (\d+) "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)" "([^"]*)"`
	re := regexp.MustCompile(pattern)

	for _, log := range logs {
		matches := re.FindStringSubmatch(log)
		if matches != nil {
			fmt.Printf("IP: %s\n", matches[1])
			fmt.Printf("Time: %s\n", matches[2])
			fmt.Printf("Operate: %s\n", matches[3])
			fmt.Printf("HTTPS Status: %s\n", matches[4])
			fmt.Printf("Content Length: %s\n", matches[5])
			fmt.Printf("Response Length: %s\n", matches[6])
			fmt.Printf("Host URL: %s\n", matches[7])
			fmt.Printf("User Agent: %s\n", matches[8])
			fmt.Printf("URL: %s\n", matches[9])
			fmt.Printf("Unknown ID: %s\n", matches[10])
			fmt.Printf("Unknown 2: %s\n", matches[11])
			fmt.Printf("Unknown 3: %s\n", matches[12])
			fmt.Printf("Action: %s\n", matches[13])
			fmt.Printf("Bucket: %s\n", matches[14])
			fmt.Println("------")
		} else {
			fmt.Println("No match found in the log line.")
		}
	}
}
