package command

import (
	"errors"
	"fmt"
	"github.com/gucooing/BaPs/pkg/alg"

	"github.com/gucooing/BaPs/common/code"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandEmailCode() {
	getEmailCode := &cdq.Command{
		Name:        "emailCode",
		AliasList:   []string{"emailCode", "ec"},
		Description: "操作指定邮箱的登录验证码",
		Permissions: cdq.Admin,
		Options: []*cdq.CommandOption{
			{
				Name:        "account",
				Description: "目标邮箱",
				Required:    true,
			},
			{
				Name:        "type",
				Description: "操作类型",
				Required:    true,
			},
			{
				Name:        "code",
				Description: "需要设置的验证码",
				Required:    false,
			},
		},
		CommandFunc: c.emailCode,
	}

	c.C.ApplicationCommand(getEmailCode)
}

// 通过邮箱拉取验证码
func (c *Command) emailCode(options map[string]string) (string, error) {
	switch options["type"] {
	case "set": // 设置验证码
		if err := code.SetCode(options["account"], alg.S2I32(options["code"])); err != nil {
			return "", err
		} else {
			return fmt.Sprintf("set account:%s coed:%v up!", options["account"], options["code"]), nil
		}
	case "get": // 获取验证码
		if codeInfo := code.GetCodeInfo(options["account"]); codeInfo != nil {
			return fmt.Sprintf("coed:%v", codeInfo.Code), nil
		} else {
			return "", errors.New("验证码已过期或失效")
		}
	case "del": // 删除验证码
		code.DelCode(options["account"])
		return fmt.Sprintf("del account:%v up!", options["account"]), nil
	default:
		return "", errors.New("error type")
	}

}
