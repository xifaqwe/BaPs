package command

import (
	"errors"
	"fmt"
	"github.com/gucooing/BaPs/config"
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/cdq"
)

type ApiSet struct {
	s    *enter.Session
	Type string
	Sub1 string
}

func (c *Command) ApplicationCommandSet() {
	set := &cdq.Command{
		Name:        "set",
		AliasList:   []string{"set"},
		Description: "直接设置一个值",
		Permissions: cdq.User,
		Options: []*cdq.CommandOption{
			{
				Name:        "uid",
				Description: "玩家游戏id",
				Required:    true,
			},
			{
				Name:        "type",
				Description: "设置类型",
				Required:    true,
			},
			{
				Name:        "sub1",
				Description: "设置参数1",
				Required:    true,
			},
		},
		CommandFunc: c.set,
	}

	c.C.ApplicationCommand(set)
}

func (c *Command) set(options map[string]string) (string, error) {
	// 玩家验证
	uid := alg.S2I64(options["uid"])
	s := enter.GetSessionByAccountServerId(uid)
	if s == nil {
		return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
	}

	req := &ApiSet{
		s:    s,
		Type: options["type"],
		Sub1: options["sub1"],
	}

	return req.parseSet()
}

func (x *ApiSet) parseSet() (string, error) {
	switch x.Type {
	case "AccountLevel":
		return x.setAccountLevel()
	case "Toast":
		return x.setToast()
	case "GetConfig":
		return config.GetConfig().String(), nil
	default:
		return "", errors.New(fmt.Sprintf("Set Type 未实现:%s", x.Type))
	}
}

func (x *ApiSet) setAccountLevel() (string, error) {
	game.SetAccountLevel(x.s, alg.S2I32(x.Sub1))
	game.AddToast(x.s, &enter.Toast{
		Text:      "已设置账号等级,请重新登录以刷新",
		BeginDate: mx.Now(),
		EndDate:   mx.Now().Add(30 * time.Second),
	})
	return fmt.Sprintf("已设置玩家等级:%v级", game.GetAccountLevel(x.s)), nil
}

func (x *ApiSet) setToast() (string, error) {
	game.AddToast(x.s, &enter.Toast{
		Text:      x.Sub1,
		BeginDate: mx.Now(),
		EndDate:   mx.Now().Add(30 * time.Second),
	})
	return fmt.Sprintf("已设置玩家通知:%s", x.Sub1), nil
}
