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
	SendTooFastError      = "send too fast"
)

var lock sync.Mutex
var mapSendRecords map[string]*sendRecords

func init() {
	mapSendRecords = make(map[string]*sendRecords, 0)
	c := cron.New()
	err := c.AddFunc("0 0/1 * * * ?", clearTimes)
	if err != nil {
		log.Error(err.Error())
	} else {
		c.Start()
	}
}

//定时清除无效记录
func clearTimes() {
	lock.Lock()
	defer lock.Unlock()
	timeOut := time.Now().Add(-time.Minute)
	list := make([]string, 0)
	for key, val := range mapSendRecords {
		if val.LastTime.Before(timeOut) {
			list = append(list, key)
		}
	}
	for _, delKey := range list {
		log.Debug("del key " + delKey)
		delete(mapSendRecords, delKey)
	}
}

type DingTalk struct {
}

type sendRecords struct {
	Key       string
	Times     int
	FirstTime time.Time
	LastTime  time.Time
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
	if err := dt.addTimes(msg.WebHookKey); err != nil {
		return nil, err
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

func (dt *DingTalk) addTimes(key string) error {
	lock.Lock()
	defer lock.Unlock()
	now := time.Now()
	c, ok := mapSendRecords[key]
	if ok {
		if now.After(c.FirstTime.Add(time.Minute)) {
			//首条消息已超过1分钟,充值状态
			c.FirstTime = now
			c.Times = 1
			return nil
		}
		if c.Times < 20 {
			//未超过规定的20次，增加次数
			c.Times = c.Times + 1
			c.LastTime = now
			return nil
		}
		return errors.New(SendTooFastError)
	} else {
		//无记录，新增
		log.Debug("add key " + key)
		mapSendRecords[key] = &sendRecords{
			Key:       key,
			Times:     1,
			FirstTime: now,
			LastTime:  now,
		}
		return nil
	}
}
