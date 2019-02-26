package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Deansquirrel/goDingtalkRobot/object"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
)

const (
	TranErrStr = "{\"errcode\":-1,\"errmsg\":\"构造返回结果时发生错误 [%s]\"}"
)

type common struct {
}

func (c *common) GetRequestBody(ctx iris.Context) string {
	body := ctx.Request().Body
	defer func() {
		_ = body.Close()
	}()
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Error("获取Http请求文本时发生错误：" + err.Error())
		return ""
	}
	return string(b)
}

func (c *common) GetOKReturn(msg string) *object.SimpleResponse {
	return c.GetMsgReturn("OK")
}

func (c *common) GetMsgReturn(msg string) *object.SimpleResponse {
	return &object.SimpleResponse{
		ErrCode: 0,
		ErrMsg:  msg,
	}
}

func (c *common) GetErrReturn(err error) *object.SimpleResponse {
	return &object.SimpleResponse{
		ErrCode: -1,
		ErrMsg:  err.Error(),
	}
}

//func (c *common) getReturn(code int, msg string) string {
//	rd := object.SimpleResponse{
//		ErrCode: code,
//		ErrMsg:  msg,
//	}
//	rb, err := json.Marshal(rd)
//	if err != nil {
//		return fmt.Sprintf(TranErrStr, "err:"+err.Error()+",code:"+strconv.Itoa(code)+",msg:"+msg)
//	} else {
//		return string(rb)
//	}
//}

//向ctx中添加返回内容
func (c *common) WriteResponse(ctx iris.Context, v interface{}) {
	str, err := json.Marshal(v)
	if err != nil {
		_, _ = ctx.WriteString(fmt.Sprintf(TranErrStr, "err:"+err.Error()))
		return
	}
	_, _ = ctx.WriteString(string(str))
	return
}

//POST发送数据
func (c *common) httpPostJsonData(data []byte, url string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	rData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(rData), nil
}

//验证阿里返回的结果消息
func (c *common) CheckAliResponse(resp string) *object.SimpleResponse {
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

type aliResponseDat struct {
	//{"errmsg":"ok","errcode":0}
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
