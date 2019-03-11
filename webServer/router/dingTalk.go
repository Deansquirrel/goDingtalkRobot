package router

import (
	"github.com/Deansquirrel/goDingtalkRobot/object"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/kataras/iris"
)

type dingTalk struct {
	app    *iris.Application
	c      common
	dtTool *object.DingTalk
}

func NewRouterDingTalk(app *iris.Application) *dingTalk {
	return &dingTalk{
		app:    app,
		c:      common{},
		dtTool: object.NewDingTalk(),
	}
}

func (dt *dingTalk) AddDingTalk() {
	//clientInfo := dt.app.Party("/DingTalk", dt.dingTalk)
	//clientInfo.Post("/Text", dt.sendTextMsg)
	dt.app.Post("/text", dt.sendTextMsg)
}

func (dt *dingTalk) dingTalk(ctx iris.Context) {
	ctx.Next()
}

func (dt *dingTalk) sendTextMsg(ctx iris.Context) {
	var tm object.DingTalkTextMsg
	err := ctx.ReadJSON(&tm)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		dt.c.WriteResponse(ctx, dt.c.GetErrReturn(err))
		//_, _ = ctx.WriteString(dt.c.GetErrReturn(err.Error()))
		log.Warn("转换请求文本时发生错误:" + err.Error() + ",requestBody:" + dt.c.GetRequestBody(ctx))
		return
	}
	re, err := dt.dtTool.SendTextMsg(&tm)
	if err != nil {
		dt.c.WriteResponse(ctx, dt.c.GetErrReturn(err))
		log.Warn(err.Error())
		return
	}
	dt.c.WriteResponse(ctx, re)
	return

}
