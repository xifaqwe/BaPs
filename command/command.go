package command

import (
	"github.com/gin-gonic/gin"
	"github.com/gucooing/BaPs/common/check"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/protocol/mx"
	"github.com/gucooing/cdq"
	cdqlog "github.com/gucooing/cdq/logger"
)

type Command struct {
	c *cdq.CDQ
}

func NewCommand(router *gin.Engine) {
	command := new(Command)
	command.c = cdq.New(&cdq.CDQ{Log: cdqlog.NewLog(cdqlog.LevelInfo, nil)})
	ginApi := cdq.NewGinApi(command.c, check.GateWaySync)
	ginApi.SetRouter(router)
	ginApi.SetApiKey(config.GetGucooingApiKey(), mx.Key)
	command.c.AddCommandRun(ginApi)

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
