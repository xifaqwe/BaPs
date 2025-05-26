package command

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gucooing/BaPs/game"

	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandGetPlayer() {
	getPlayer := &cdq.Command{
		Name:        "getPlayer",
		AliasList:   []string{"getPlayer", "gp"},
		Description: "获取指定玩家的db数据",
		Permissions: cdq.Admin,
		Options: []*cdq.CommandOption{
			{
				Name:        "uid",
				Description: "玩家游戏id",
				Required:    true,
			},
			{
				Name:        "json",
				Description: "输出为json或base64",
				Required:    false,
			},
			{
				Name:        "basis",
				Description: "获取玩家基础数据",
				Required:    false,
			},
		},
		CommandFunc: c.getPlayer,
	}

	c.C.ApplicationCommand(getPlayer)
}

func (c *Command) getPlayer(options map[string]string) (string, error) {
	uid := alg.S2I64(options["uid"])
	session := enter.GetSessionByUid(uid)
	if session == nil {
		return "", errors.New(fmt.Sprintf("玩家未注册 UID:%v", uid))
	}
	var info any
	if alg.S2I64(options["basis"]) != 0 {
		info = game.GetPlayerBin(session).GetBaseBin()
	} else {
		info = game.GetPlayerBin(session)
	}

	jsonInfo, err := sonic.Marshal(info)
	if err != nil {
		return "", err
	}

	if alg.S2I64(options["json"]) != 0 {
		return string(jsonInfo), nil
	} else {
		return base64.RawStdEncoding.EncodeToString(jsonInfo), nil
	}
}
