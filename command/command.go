package command

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/check"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/protocol/mx"
	"github.com/gucooing/cdq"
	cdqlog "github.com/gucooing/cdq/logger"
)

var C *Command

type Command struct {
	C *cdq.CDQ
}

func NewCommand(router *gin.Engine) {
	command := new(Command)
	C = command
	command.C = cdq.New(&cdq.CDQ{Log: cdqlog.NewLog(cdqlog.LevelInfo, nil)})
	ginApi := cdq.NewGinApi(command.C, check.GateWaySync)
	ginApi.SetRouter(router)
	ginApi.SetApiKey(config.GetGucooingApiKey(), mx.Key)
	command.C.AddCommandRun(ginApi)

	// 注册指令
	command.ApplicationCommandGiveAll()
	command.ApplicationCommandGetPlayer()
	command.ApplicationCommandEmailCode()
	command.ApplicationCommandGameMail()
	command.ApplicationCommandMail()
	command.ApplicationCommandSet()
	command.ApplicationCommandPing()
	command.ApplicationCommandCharacter()
	command.ApplicationCommandAccount()
}
