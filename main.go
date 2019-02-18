package main

import (
	"context"
	"github.com/Deansquirrel/goDingtalkRobot/common"
	"github.com/Deansquirrel/goDingtalkRobot/global"
	log "github.com/Deansquirrel/goToolLog"
	"time"
)

func main() {
	//==================================================================================================================
	config, err := common.GetSysConfig("config.toml")
	if err != nil {
		log.Error("加载配置文件时遇到错误：" + err.Error())
		return
	}
	config.FormatConfig()
	global.SysConfig = config
	err = common.RefreshSysConfig(*global.SysConfig)
	if err != nil {
		log.Error("刷新配置时遇到错误：" + err.Error())
		return
	}
	global.Ctx, global.Cancel = context.WithCancel(context.Background())
	//==================================================================================================================
	log.Info("程序启动")
	defer log.Info("程序退出")
	//==================================================================================================================
	time.AfterFunc(time.Second*5, func() {
		global.Cancel()
	})
	//==================================================================================================================
	select {
	case <-global.Ctx.Done():
	}
}
