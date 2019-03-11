package global

import (
	"context"
	"github.com/Deansquirrel/goDingtalkRobot/config"
)

const (
	//LastVersion = "1.0.1 Build201903111529"
	Version = "1.0.2 Build20190311"
)

var SysConfig *config.SysConfig
var Ctx context.Context
var Cancel func()
