package global

import (
	"context"
	"github.com/Deansquirrel/goDingtalkRobot/config"
)

const (
	//LastVersion = "1.0.3 Build20190312"
	//Version = "0.0.0 Build20100101"
	Version = "0.0.0 Build20100101"
)

var SysConfig *config.SysConfig
var Ctx context.Context
var Cancel func()
