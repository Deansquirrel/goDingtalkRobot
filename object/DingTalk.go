package object

import (
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goMonitorV2/object"
	"github.com/kataras/iris/core/errors"
	"github.com/robfig/cron"
	"sync"
	"time"
)
import log "github.com/Deansquirrel/goToolLog"

//https://oapi.dingtalk.com/robot/send?access_token=7a84d09b83f9633ad37866505d2c0c26e39f4fa916b3af8f6a702180d3b9906b

const (
	DingTalkWebHookFormat = "https://oapi.dingtalk.com/robot/send?access_token=%s"
	//SendTooFastError      = "send too fast"
)

type DingTalk struct {
	lock           sync.Mutex
	mapSendRecords map[string]time.Time
}

func NewDingTalk() *DingTalk {
	dt := DingTalk{}
	dt.mapSendRecords = make(map[string]time.Time, 0)
	c := cron.New()
	err := c.AddFunc("0 0/1 * * * ?", dt.clearTimes)
	if err != nil {
		log.Error(err.Error())
	} else {
		c.Start()
	}
	return &dt
}

type aliResponseDat struct {
	//{"errmsg":"ok","errcode":0}
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (dt *DingTalk) GetWebHookUrl(key string) string {
	return fmt.Sprintf(DingTalkWebHookFormat, key)
}

func (dt *DingTalk) SendTextMsg(msg *DingTalkTextMsg) (*object.SimpleResponse, error) {
	if dTime, ok := dt.getDelayTime(msg.WebHookKey); ok {
		log.Debug("delay")
		time.Sleep(dTime)
	}
	rm, err := msg.GetAliMsgStr()
	if err != nil {
		errMsg := "获取请求文本时发生错误:" + err.Error()
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	log.Debug(rm)
	comm := Common{}
	re, err := comm.HttpPostJsonData([]byte(rm), dt.GetWebHookUrl(msg.WebHookKey))
	if err != nil {
		errMsg := "发送Http数据时发生错误:" + err.Error()
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return dt.CheckAliResponse(re), nil
}

//验证阿里返回的结果消息
func (dt *DingTalk) CheckAliResponse(resp string) *object.SimpleResponse {
	log.Debug(resp)
	var ar aliResponseDat
	err := json.Unmarshal([]byte(resp), &ar)
	if err != nil {
		return &object.SimpleResponse{
			ErrCode: -1,
			ErrMsg:  "验证返回结果时发生错误：" + err.Error(),
		}
		//c.GetErrReturn("验证返回结果时发生错误：" + err.Error())
	}
	return &object.SimpleResponse{
		ErrCode: ar.ErrCode,
		ErrMsg:  ar.ErrMsg,
	}
}

func (dt *DingTalk) getDelayTime(key string) (time.Duration, bool) {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	c, ok := dt.mapSendRecords[key]
	if ok {
		if c.Before(time.Now().Add(-time.Second * 3)) {
			dt.mapSendRecords[key] = time.Now()
			return time.Second, false
		} else {
			dt.mapSendRecords[key] = c.Add(time.Second * 3)
			return time.Until(c.Add(time.Second * 3)), true
		}
	} else {
		//无记录，新增
		log.Debug("add key " + key)
		dt.mapSendRecords[key] = time.Now()
		return time.Second, false
	}
}

//定时清除无效记录
func (dt *DingTalk) clearTimes() {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	outTime := time.Now().Add(-time.Minute)
	list := make([]string, 0)
	for key, val := range dt.mapSendRecords {
		if val.Before(outTime) {
			list = append(list, key)
		}
	}
	for _, key := range list {
		delete(dt.mapSendRecords, key)
	}
}
