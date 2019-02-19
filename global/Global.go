package global

import (
	"context"
	"github.com/Deansquirrel/goDingtalkRobot/config"
)

var SysConfig *config.SysConfig
var Ctx context.Context
var Cancel func()

const (
	Version = "0.0.0 Build20190218"
)
