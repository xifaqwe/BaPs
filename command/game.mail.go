package command

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	dbstruct "github.com/gucooing/BaPs/db/struct"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandGameMail() {
	mail := &cdq.Command{
		Name:        "gameMail",
		AliasList:   []string{"gameMail", "gm"},
		Description: "发送一封游戏邮件,对象可以是全局或单个玩家",
		Permissions: cdq.Admin,
		Options: []*cdq.CommandOption{
			{
				Name:        "recipient",
				Description: "玩家游戏id或all全局对象",
				Required:    true,
			},
			{
				Name:        "sender",
				Description: "发件人",
				Required:    true,
			},
			{
				Name:        "comment",
				Description: "邮件文本内容",
				Required:    true,
			},
			{
				Name:        "sendDate",
				Description: "领取开始时间",
				Required:    true,
			},
			{
				Name:        "expireDate",
				Description: "邮件有效截至时间",
				Required:    true,
			},
			{
				Name:        "parcelInfoList",
				Description: "附件 json格式",
				Required:    false,
			},
			{
				Name:        "del",
				Description: "删除邮件 可传入邮件id或all",
				Required:    false,
			},
		},
		CommandFunc: c.gameMail,
	}

	c.C.ApplicationCommand(mail)
}

func (c *Command) gameMail(options map[string]string) (string, error) {
	switch options["recipient"] {
	case "all", "All": // 全局邮件操作
		return gmYostarMail(options)
	default:
		return gmPlayerMail(options)
	}
}

func gmYostarMail(options map[string]string) (string, error) {
	if strDel := options["del"]; strDel == "all" || strDel == "All" {
		if enter.DelAllYostarMail() != nil {
			return "", errors.New("删除所有全局邮件失败")
		} else {
			return "删除所有全局邮件成功", nil
		}
	} else if delId := alg.S2I64(strDel); delId != 0 {
		if enter.DelYostarMail(delId) != nil {
			return "", errors.New("删除全局邮件失败")
		} else {
			return "删除全局邮件成功", nil
		}
	} else {
		parcelInfoList, err := genParcelInfo[dbstruct.ParcelInfo](options["parcelInfoList"])
		if err != nil {
			return "", errors.New(fmt.Sprintf("解析邮件附件失败:%s", err.Error()))
		}
		sendMail := &dbstruct.YostarMail{
			Sender:         options["sender"],
			Comment:        options["comment"],
			SendDate:       sql.NullTime{Time: time.Unix(alg.S2I64(options["sendDate"]), 0), Valid: true},
			ExpireDate:     sql.NullTime{Time: time.Unix(alg.S2I64(options["expireDate"]), 0), Valid: true},
			ParcelInfoList: parcelInfoList,
		}
		if enter.AddYostarMail(sendMail) {
			return "全局邮件发送成功", nil
		}
		return "", errors.New("全局邮件发送失败")
	}
}

func gmPlayerMail(options map[string]string) (string, error) {
	uid := alg.S2I64(options["recipient"])
	s := enter.GetSessionByAccountServerId(uid)
	if s == nil {
		return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
	}
	if strDel := options["del"]; strDel == "all" || strDel == "All" {
		if game.DelAllMail(s) {
			return "删除玩家全部邮件成功", nil
		} else {
			return "", errors.New("删除玩家全部邮件失败,这是一个神奇的bug")
		}
	} else if delId := alg.S2I64(strDel); delId != 0 {
		if game.DelMail(s, delId) {
			return "删除玩家邮件成功", nil
		} else {
			return "", errors.New("删除玩家邮件失败,邮件可能不存在")
		}
	} else {
		parcelInfoList, err := genParcelInfo[sro.ParcelInfo](options["parcelInfoList"])
		if err != nil {
			return "", errors.New(fmt.Sprintf("解析邮件附件失败:%s", err.Error()))
		}
		sendMail := &sro.MailInfo{
			Sender:         options["sender"],
			Comment:        options["comment"],
			SendDate:       alg.S2I64(options["sendDate"]),
			ExpireDate:     alg.S2I64(options["expireDate"]),
			ParcelInfoList: parcelInfoList,
		}
		if game.AddMail(s, sendMail) {
			return "请查询游戏内邮箱获取结果", nil
		}
		return "", errors.New("私人邮件发送失败")
	}
}

func genParcelInfo[T any](str string) ([]*T, error) {
	list := make([]*T, 0)
	if str == "" {
		return list, nil
	}
	err := json.Unmarshal([]byte(str), &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}
