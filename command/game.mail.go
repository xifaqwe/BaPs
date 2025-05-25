package command

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	dbstruct "github.com/gucooing/BaPs/db/struct"
	"strconv"
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
				Name:        "uid",
				Description: "玩家游戏id",
				Required:    false,
			},
			{
				Name:        "player",
				Description: "是否为玩家邮件,0为全局邮件",
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
				Name:        "uid",
				Description: "如果是玩家邮件,则此选项必写",
				Required:    false,
			},
			{
				Name:        "parcelInfoList",
				Description: "附件 json格式",
				Required:    false,
			},
		},
		CommandFunc: c.gameMail,
	}

	c.c.ApplicationCommand(mail)
}

func (c *Command) gameMail(options map[string]string) (string, error) {
	parcelInfoListOption := options["parcelInfoList"]

	if player, _ := strconv.ParseBool(options["player"]); !player {
		parcelInfoList, err := genParcelInfo[dbstruct.ParcelInfo](parcelInfoListOption)
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
	} else {
		// 玩家验证
		uid := alg.S2I64(options["uid"])
		s := enter.GetSessionByAccountServerId(uid)
		if s == nil {
			return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
		}
		parcelInfoList, err := genParcelInfo[sro.ParcelInfo](parcelInfoListOption)
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
