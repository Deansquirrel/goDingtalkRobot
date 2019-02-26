package config

import (
	"github.com/Deansquirrel/goToolCommon"
	"strings"
	"time"
)

type SysConfig struct {
	Total         Total         `toml:"total"`
	Iris          iris          `toml:"iris"`
	ServiceConfig serviceConfig `toml:"serviceConfig"`
}

type iris struct {
	Port     int    `toml:"port"`
	LogLevel string `toml:"logLevel"`
}

type serviceConfig struct {
	Name        string `toml:"name"`
	DisplayName string `toml:"displayName"`
	Description string `toml:"description"`
}

//返回配置字符串
func (sc *SysConfig) GetConfigStr() (string, error) {
	return goToolCommon.GetJsonStr(sc)
}

//配置检查并格式化
func (sc *SysConfig) FormatConfig() {
	sc.Total.FormatConfig()
	sc.Iris.FormatConfig()
	sc.ServiceConfig.FormatConfig()
}

//格式化
func (i *iris) FormatConfig() {
	//设置默认端口 8000
	if i.Port == 0 {
		i.Port = 8000
	}
	//去除首尾空格
	i.LogLevel = strings.Trim(i.LogLevel, " ")
	//设置Iris默认日志级别
	if i.LogLevel == "" {
		i.LogLevel = "warn"
	}
	//设置字符串转换为小写
	i.LogLevel = strings.ToLower(i.LogLevel)
}

//格式化
func (sc *serviceConfig) FormatConfig() {
	sc.Name = strings.Trim(sc.Name, " ")
	sc.DisplayName = strings.Trim(sc.DisplayName, " ")
	sc.Description = strings.Trim(sc.Description, " ")
	if sc.Name == "" {
		sc.Name = "DingTalkMsgService" + goToolCommon.GetDateTimeStr(time.Now())
	}
	if sc.DisplayName == "" {
		sc.DisplayName = "DTMS" + goToolCommon.GetDateTimeStr(time.Now())
	}
	if sc.Description == "" {
		sc.Description = sc.Name
	}
}
