package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Deansquirrel/goToolLog"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	TranErrStr = "{\"code\":-1,\"msg\":\"构造返回结果时发生错误 [%s]\"}"
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

func (c *common) GetMsgReturn(msg string) string {
	return c.getReturn(0, msg)
}

func (c *common) GetErrReturn(msg string) string {
	return c.getReturn(-1, msg)
}

func (c *common) getReturn(code int, msg string) string {
	rd := responseDao{
		ErrCode: code,
		ErrMsg:  msg,
	}
	rb, err := json.Marshal(rd)
	if err != nil {
		return fmt.Sprintf(TranErrStr, "err:"+err.Error()+",code:"+strconv.Itoa(code)+",msg:"+msg)
	} else {
		return string(rb)
	}
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

func (c *common) checkResponse(resp string) string {
	var ar aliResponseDat
	err := json.Unmarshal([]byte(resp), &ar)
	if err != nil {
		return c.GetErrReturn("验证返回结果时发生错误：" + err.Error())
	}
	return c.getReturn(ar.ErrCode, ar.ErrMsg)
}

type aliResponseDat struct {
	//{"errmsg":"ok","errcode":0}
	ErrMsg  string `json:"errmsg"`
	ErrCode int    `json:"errcode"`
}

type responseDao struct {
	ErrCode int    `json:"code"`
	ErrMsg  string `json:"msg"`
}
