package object

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type Common struct {
}

//POST发送数据
func (c *Common) HttpPostJsonData(data []byte, url string) (string, error) {
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
