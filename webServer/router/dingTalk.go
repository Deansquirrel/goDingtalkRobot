package router

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goDingtalkRobot/object"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/kataras/iris"
)

//https://oapi.dingtalk.com/robot/send?access_token=7a84d09b83f9633ad37866505d2c0c26e39f4fa916b3af8f6a702180d3b9906b

const (
	DingTalkWebHookFormat = "https://oapi.dingtalk.com/robot/send?access_token=%s"
)

type dingTalk struct {
	app *iris.Application
	c   common
}

func NewRouterDingTalk(app *iris.Application) *dingTalk {
	return &dingTalk{
		app: app,
		c:   common{},
	}
}

func (dt *dingTalk) AddDingTalk() {
	clientInfo := dt.app.Party("/DingTalk", dt.dingTalk)
	clientInfo.Post("/Text", dt.sendTextMsg)
}

func (dt *dingTalk) dingTalk(ctx iris.Context) {
	ctx.Next()
}

func (dt *dingTalk) sendTextMsg(ctx iris.Context) {
	var tm object.DingTalkTextMsg
	err := ctx.ReadJSON(&tm)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.WriteString(dt.c.GetErrReturn(err.Error()))
		log.Warn("转换请求文本时发生错误，" + err.Error())
		log.Warn(dt.c.GetRequestBody(ctx))
		return
	}
	ts, err := json.Marshal(tm)
	fmt.Println(string(ts))
	rm, err := tm.GetAliMsgStr()
	if err != nil {
		_, _ = ctx.WriteString(dt.c.GetErrReturn(err.Error()))
		log.Warn("获取阿里请求文本时发生错误，" + err.Error())
		return
	}
	log.Debug(rm)
	re, err := dt.c.httpPostJsonData([]byte(rm), dt.getWebHookUrl(tm.WebHookKey))
	if err != nil {
		_, _ = ctx.WriteString(dt.c.GetErrReturn(err.Error()))
		log.Warn("发送Http数据时发生错误，" + err.Error())
		return
	}
	_, _ = ctx.WriteString(dt.c.checkResponse(re))
	return
}

func (dt *dingTalk) getWebHookUrl(key string) string {
	return fmt.Sprintf(DingTalkWebHookFormat, key)
}
