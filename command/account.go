package command

import (
	"errors"
	"fmt"
	"github.com/gucooing/BaPs/sdk"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandAccount() {
	account := &cdq.Command{
		Name:        "account",
		AliasList:   []string{"account", "ac"},
		Description: "注册账户",
		Permissions: cdq.User,
		Options: append(playerOptions, []*cdq.CommandOption{
			{
				Name:        "name",
				Description: "账户昵称",
				Required:    true,
			},
		}...),
		CommandFunc: c.account,
	}

	c.c.ApplicationCommand(account)
}

func (c *Command) account(options map[string]*cdq.CommandOption) (string, error) {
	nameOption, ok := options["name"]
	if !ok {
		return "", errors.New("缺少参数 name")
	}

	ya, err := sdk.AddYostarAccount(nameOption.Option, true)
	if err != nil || ya.YostarAccount != nameOption.Option {
		return "", errors.New(fmt.Sprintf("账户注册失败 Account:%s", nameOption.Option))
	}

	return fmt.Sprintf("账户注册成功 Account:%s", nameOption.Option), nil
}
