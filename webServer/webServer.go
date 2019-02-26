package webServer

import (
	"context"
	"github.com/Deansquirrel/goDingtalkRobot/global"
	wsRouter "github.com/Deansquirrel/goDingtalkRobot/webServer/router"
	"github.com/Deansquirrel/goToolCommon"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type webServer struct {
	app  *iris.Application
	port int
}

func NewWebServer(port int) *webServer {
	return &webServer{
		app:  iris.New(),
		port: port,
	}
}

//启动Web服务
func (ws *webServer) StartWebService() {
	app := iris.New()

	app.Use(recover.New())
	app.Use(logger.New())

	//iris.RegisterOnInterrupt(func() {
	//	log.Info("StopWebService")
	//	timeout := 5 * time.Second
	//	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	//	defer cancel()
	//	_ = app.Shutdown(ctx)
	//})

	ws.irisInit(app)
	ws.irisRouter(app)
	ws.irisStart(app)

	select {
	case <-global.Ctx.Done():
	}
}

//iris初始化
func (ws *webServer) irisInit(app *iris.Application) {
	ws.setIrisLogWrap(app)
	ws.setIrisLogLevel(app)
	ws.setIrisLogTimeFormat(app)
	ws.setIrisLogFile(app)
}

//iris添加路由
func (ws *webServer) irisRouter(app *iris.Application) {
	routerBase := wsRouter.NewRouterBase(app)
	routerBase.AddBase()
	routerDingTalk := wsRouter.NewRouterDingTalk(app)
	routerDingTalk.AddDingTalk()
}

//iris启动
func (ws *webServer) irisStart(app *iris.Application) {
	log.Warn("StartWebService")
	defer log.Warn("StartWebService Complete")
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			os.Interrupt,
			syscall.SIGINT,
			os.Kill,
			syscall.SIGKILL,
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			ws.irisStop(app)
			defer global.Cancel()
		case <-global.Ctx.Done():
			ws.irisStop(app)
		}
	}()
	go func() {
		_ = app.Run(
			iris.Addr(":"+strconv.Itoa(ws.port)),
			iris.WithoutInterruptHandler,
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
		)
	}()
}

//iris停止
func (ws *webServer) irisStop(app *iris.Application) {
	log.Warn("StopWebService")
	defer log.Warn("StopWebService complete")
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_ = app.Shutdown(ctx)
}

//校验SysConfig中iris日志级别设置
func (ws *webServer) checkIrisLogLevel(level string) string {
	switch level {
	case "debug", "info", "warn", "error":
		return level
	default:
		return "warn"
	}
}

//设置iris日志级别
func (ws *webServer) setIrisLogLevel(app *iris.Application) {
	app.Logger().SetLevel(ws.checkIrisLogLevel(global.SysConfig.Iris.LogLevel))
}

//设置iris日志时间格式
func (ws *webServer) setIrisLogTimeFormat(app *iris.Application) {
	app.Logger().SetTimeFormat("2006-01-02 15:04:05")
}

//设置iris日志换行格式
func (ws *webServer) setIrisLogWrap(app *iris.Application) {
	app.Logger().DisableNewLine()
	app.Logger().SetPrefix(goToolCommon.GetWrapStr())
}

//设置日志输出文件对象,需按日分割
func (ws *webServer) setIrisLogFile(app *iris.Application) {
	ws.reSetLogFile(app)
	time.AfterFunc(ws.getRemainingTime(), func() {
		ws.setIrisLogFile(app)
	})
}

//获取当日所剩时间
func (ws *webServer) getRemainingTime() time.Duration {
	todayLast := goToolCommon.GetDateStr(time.Now()) + " 23:59:59"
	todayLastTime, err := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)
	if err != nil {
		log.Warn("获取当日所剩时间时发生错误:" + err.Error())
		return time.Minute
	}
	return time.Duration(todayLastTime.Unix()-time.Now().Local().Unix()+1) * time.Second
}

//设置日志输出对象
func (ws *webServer) reSetLogFile(app *iris.Application) {
	path, err := ws.getIrisLogPath()
	if err != nil {
		log.Warn("刷新iris日志对象,获取当前路径时遇到错误:" + err.Error())
		return
	}
	fileName := "iris_" + goToolCommon.GetDateStr(time.Now()) + ".log"
	w, err := ws.getIrisLogWriter(path, fileName)
	if err != nil {
		log.Warn("刷新iris日志对象,获取io.Writer遇到错误:" + err.Error())
		return
	}
	if w != nil {
		app.Logger().SetOutput(w)
		if global.SysConfig.Total.StdOut {
			app.Logger().AddOutput(os.Stdout)
		}
		log.Debug("SetLogFile")
	}
}

//获取iris写日志对象
func (ws *webServer) getIrisLogWriter(path string, fileName string) (io.Writer, error) {
	f, err := os.OpenFile(path+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}

//获取iris日志路径
func (ws *webServer) getIrisLogPath() (string, error) {
	path := log.Path
	if strings.Trim(path, " ") == "" {
		path, err := goToolCommon.GetCurrPath()
		if err != nil {
			return "", err
		}
		return path + "\\" + "log" + "\\", nil
	}
	if strings.HasSuffix(path, "\\") {
		path = path + "\\"
	}
	return path, nil
}
