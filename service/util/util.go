package util

import (
	"net"
	"strings"
	"time"

	"github.com/twelveeee/log_analysis/service/eventLog"
)

var log = eventLog.NewLog()

func IpToInt(ipStr string) uint32 {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		log.Error().Msgf("ip format err ,value : %s", ipStr)
		return 0
	}

	ip = ip.To4()
	if ip == nil {
		log.Error().Msgf("ipv4 format  err ,value : %s", ipStr)
		return 0
	}

	ipInt := uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
	return ipInt
}

func ConverToTimestamp(timeStr string) time.Time {
	// 去除头尾中括号
	timeStr = strings.Trim(timeStr, "[]")
	timeStr = strings.Trim(timeStr, ":")

	// 解析时间字符串
	layout := "02/Jan/2006:15:04:05 -0700"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Error().Msgf("unkwon timeStr :%s", timeStr)
		return time.Time{}
	}

	return t
}
