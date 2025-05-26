package command

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gucooing/BaPs/db"
	"github.com/gucooing/BaPs/sdk"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandAccount() {
	account := &cdq.Command{
		Name:        "account",
		AliasList:   []string{"account", "ac"},
		Description: "操作玩家账户",
		Permissions: cdq.User,
		Options: []*cdq.CommandOption{
			{
				Name:        "name",
				Description: "账户昵称",
				Required:    true,
			},
			{
				Name:        "type",
				Description: "操作类型",
				Required:    true,
			},
			{
				Name:        "banMsg",
				Description: "封禁原因",
				Required:    false,
			},
		},
		CommandFunc: c.account,
	}

	c.C.ApplicationCommand(account)
}

func (c *Command) account(options map[string]string) (string, error) {
	switch options["type"] {
	case "login": // 注册
		ya, err := sdk.GetORAddYostarAccount(options["name"], true)
		if err != nil || ya.YostarAccount != options["name"] {
			return "", errors.New(fmt.Sprintf("账户注册失败 Account:%s", options["name"]))
		}
		return fmt.Sprintf("账户注册成功 Account:%s", options["name"]), nil
	case "ban": // 封禁
		yul := sdk.GetYostarUserLoginByAccount(options["name"])
		if yul == nil {
			return "", errors.New(fmt.Sprintf("账户不存在 Account:%s", options["name"]))
		}
		yul.Ban = true
		yul.BanMsg = options["banMsg"]
		if db.GetDBGame().UpdateYoStarUserLogin(yul) != nil {
			return "", errors.New(fmt.Sprintf("数据库操作失败 Account:%s", options["name"]))
		}
		return fmt.Sprintf("ban Account:%s up!", options["name"]), nil
	case "unban": // 解除封禁
		yul := sdk.GetYostarUserLoginByAccount(options["name"])
		if yul == nil {
			return "", errors.New(fmt.Sprintf("账户不存在 Account:%s", options["name"]))
		}
		yul.Ban = false
		yul.BanMsg = ""
		if db.GetDBGame().UpdateYoStarUserLogin(yul) != nil {
			return "", errors.New(fmt.Sprintf("数据库操作失败 Account:%s", options["name"]))
		}
		return fmt.Sprintf("unBan Account:%s up!", options["name"]), nil
	case "get": // 获取账户详情
		yul := sdk.GetYostarUserLoginByAccount(options["name"])
		if yul == nil {
			return "", errors.New(fmt.Sprintf("账户不存在 Account:%s", options["name"]))
		}
		return sonic.MarshalString(&AccountInfo{
			Account:         options["name"],
			AccountServerId: yul.AccountServerId,
			YostarUid:       yul.YostarUid,
			Ban:             yul.Ban,
			BanMsg:          yul.BanMsg,
		})
	default:
		return "", errors.New("error type")
	}
}

type AccountInfo struct {
	Account         string `json:"account"`
	AccountServerId int64  `json:"account_server_id"`
	YostarUid       int64  `json:"yostar_uid"`
	Ban             bool   `json:"ban"`
	BanMsg          string `json:"banMsg"`
}
