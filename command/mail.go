package command

import (
	"fmt"
	"github.com/gucooing/BaPs/common/mail"
	"github.com/gucooing/cdq"
	"strings"
)

func (c *Command) ApplicationCommandMail() {
	apiMail := &cdq.Command{
		Name:        "mail",
		AliasList:   []string{"mail", "m"},
		Description: "发送一封邮件,对象可以是多个用户",
		Permissions: cdq.Admin,
		Options: []*cdq.CommandOption{
			{
				Name:        "header",
				Description: "邮件标题",
				Required:    true,
			},
			{
				Name:        "body",
				Description: "邮件内容,必须是string格式",
				Required:    true,
			},
			{
				Name:        "usernames",
				Description: "送达邮箱,可以是多个邮箱,用;隔开",
				Required:    true,
			},
		},
		CommandFunc: c.mail,
	}

	c.C.ApplicationCommand(apiMail)
}

func (c *Command) mail(options map[string]string) (string, error) {
	usernames := strings.Split(options["usernames"], ";")
	err := mail.SendTextMail(options["header"], options["body"], usernames...)
	if err != nil {
		return "", fmt.Errorf("邮件发送失败err:%s", err)
	}
	return "邮件发送成功", nil
}
