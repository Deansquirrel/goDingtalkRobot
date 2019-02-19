package config

import (
	"github.com/Deansquirrel/goToolCommon"
	"strings"
)

type SysConfig struct {
	Total Total `toml:"total"`
	Iris  iris  `toml:"iris"`
}

type iris struct {
	Port     int    `toml:"port"`
	LogLevel string `toml:"logLevel"`
}

//返回配置字符串
func (sc *SysConfig) GetConfigStr() (string, error) {
	return goToolCommon.GetJsonStr(sc)
}

//配置检查并格式化
func (sc *SysConfig) FormatConfig() {
	sc.Total.FormatConfig()
	sc.Iris.FormatConfig()
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
