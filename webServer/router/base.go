package router

import (
	"github.com/Deansquirrel/goDingtalkRobot/global"
	"github.com/Deansquirrel/goToolCommon"
	"github.com/kataras/iris"
)

type base struct {
	app *iris.Application
	c   common
}

type versionInfo struct {
	Version string `json:"version"`
}

func NewRouterBase(app *iris.Application) *base {
	return &base{
		app: app,
		c:   common{},
	}
}

func (base *base) AddBase() {
	base.app.Get("/version", base.version)
}

func (base *base) version(ctx iris.Context) {
	v := versionInfo{
		Version: global.Version,
	}
	str, err := goToolCommon.GetJsonStr(v)
	if err != nil {
		_, _ = ctx.WriteString(base.c.GetErrReturn(err.Error()))
		return
	}
	_, _ = ctx.WriteString(str)
	return
}
