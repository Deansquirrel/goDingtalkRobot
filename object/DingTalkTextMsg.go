package object

import "github.com/Deansquirrel/goToolCommon"

type DingTalkTextMsg struct {
	WebHookKey string   `json:"webhookkey"`
	Content    string   `json:"content"`
	AtMobiles  []string `json:"atmobiles"`
	IsAtAll    bool     `json:"isatall"`
}

type aliDingTalkTextMsg struct {
	MsgType string                 `json:"msgtype"`
	Text    aliDingTalkTextMsgTest `json:"text"`
	At      aliDingTalkTextMsgAt   `json:"at"`
}

type aliDingTalkTextMsgTest struct {
	Content string `json:"content"`
}

type aliDingTalkTextMsgAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

func (dt *DingTalkTextMsg) GetAliMsgStr() (string, error) {
	ro := aliDingTalkTextMsg{}
	ro.MsgType = "text"
	ro.Text = aliDingTalkTextMsgTest{
		Content: dt.Content,
	}
	ro.At = aliDingTalkTextMsgAt{
		AtMobiles: make([]string, 0),
		IsAtAll:   dt.IsAtAll,
	}
	for _, s := range dt.AtMobiles {
		ro.At.AtMobiles = append(ro.At.AtMobiles, s)
	}

	return goToolCommon.GetJsonStr(ro)
}
