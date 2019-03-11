package global

import (
	"context"
	"github.com/Deansquirrel/goDingtalkRobot/config"
)

const (
	//HisVersion = "0.0.0 Build20000101"
	Version = "1.0.1 Build201903111529"
)

var SysConfig *config.SysConfig
var Ctx context.Context
var Cancel func()
