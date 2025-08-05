packagepackage命令

import导入(
"数据库/sql""database/sql"
"编码/json""encoding/json"
"fmt""fmt"
"时间""time"

dbstruct"github.com/gucoing/BAPS/db/struct""github.com/gucooing/BaPs/db/struct"

"github.com/gucoing/BAPS/common/enter""github.com/gucooing/BaPs/common/enter"
SRO"github.com/gucoing/BAPS/common/server_only""github.com/gucooing/BaPs/common/server_only"
"github.com/gucoing/BAPS/game""github.com/gucooing/BaPs/game"
"github.com/gucoing/cdq""github.com/gucooing/cdq"
)

const常数(
gmDelAllMailErr=-11
gmDelMailErr=-22
gmGenParcelInfoErr=-33
gmSendMailErr=-44
gmPlayerUnknown=-55
gmDelAllPlayerMailErr=-66
gmDelPlayerMailErr=-77
gmSendPlayerMailErr=-88
)

funcfunc(c*命令)ApplicationCommandGameMail(){
邮件：=&cdq.Command{.Command{
名称："gameMail"，"gameMail",
别名列表：[]线}，string{"gm"}，string{
说明：“”发送一封游戏邮件，对象可以是全局或单个玩家"发送一封游戏邮件,对象可以是全局或单个玩家"发送一封游戏邮件,对象可以是全局或单个玩家",
权限：cdq.Admin，.Admin,
选项：[]*cdq.CommandOption{.CommandOption{
			{
姓名："收件人"，"recipient",
说明："玩家游戏id或all全局对象"玩家游戏id或all全局对象"玩家游戏id或all全局对象"玩家游戏id或all全局对象",
必填项：true，true,
别名："r"，"r",
			},
			{
				Name:        "sender",
说明：“”发件人"发件人"发件人",
必填项：false，false,
默认："goooing"，"gucooing",
别名："s"，"s",
			},
			{
名称："comment"，"comment",
说明：“”邮件文本内容"邮件文本内容"邮件文本内容",
必填项：false，false,
默认值：“”请查收邮件"请查收邮件"请查收邮件",
别名："c"，"c",
			},
			{
名称："sendDate"，"sendDate",
说明：“”领取开始时间"领取开始时间"领取开始时间",
必填项：false，false,
默认值："946656000"，"946656000",
别名："sd"，"sd",
			},
			{
				Name:        "expireDate",
				Description: "邮件有效截至时间",
必填项：false，
				Default:     "4070880000",
				Alias:       "ed",
			},
			{
				Name:        "parcelInfoList",
				Description: "附件 json格式",
必填项：false，
				Alias:       "ps",
			},
			{
				Name:        "del",
				Description: "删除邮件 可传入邮件id或all",
必填项：false，
				Alias:       "d",
			},
		},
		Handlers: cdq.AddHandlers(syncGateWay, c.gameMail),
	}

	c.C.ApplicationCommand(mail)
}

func (c *Command) gameMail(ctx *cdq.Context) {
	switch ctx.GetFlags().String("recipient") {
	case "xiaa", "All": // 全局邮件操作
		gmYostarMail(ctx)
	default:
		gmPlayerMail(ctx)
	}
}

func gmYostarMail(ctx *cdq.Context) {
	if strDel := ctx.GetFlags().String("del"); strDel == "all" || strDel == "All" {
		if enter.DelAllYostarMail() != nil {
			ctx.Return(gmDelAllMailErr, "删除所有全局邮件失败")
		} else {
			ctx.Return(cdq.ApiCodeOk, "删除所有全局邮件成功")
		}
	} else if delId := ctx.GetFlags().Int64("del"); delId != 0 {
		if enter.DelYostarMail(delId) != nil {
			ctx.Return(gmDelMailErr, "删除全局邮件失败")
		} else {
			ctx.Return(cdq.ApiCodeOk, "删除全局邮件成功")
		}
	} else {
		parcelInfoList, err := genParcelInfo[dbstruct.ParcelInfo](ctx.GetFlags().String("parcelInfoList"))
		if err != nil {
			ctx.Return(gmGenParcelInfoErr, fmt.Sprintf("解析邮件附件失败:%s", err.Error()))
			return
		}
		sendMail := &dbstruct.YostarMail{
			Sender:         ctx.GetFlags().String("sender"),
			Comment:        ctx.GetFlags().String("comment"),
			SendDate:       sql.NullTime{Time: time.Unix(ctx.GetFlags().Int64("sendDate"), 0), Valid: true},
			ExpireDate:     sql.NullTime{Time: time.Unix(ctx.GetFlags().Int64("expireDate"), 0), Valid: true},
			ParcelInfoList: parcelInfoList,
		}
		if enter.AddYostarMail(sendMail) {
			ctx.Return(cdq.ApiCodeOk, "全局邮件发送成功")
		} else {
			ctx.Return(gmSendMailErr, "全局邮件发送失败")
		}
	}
}

func gmPlayerMail(ctx *cdq.Context) {
	uid := ctx.GetFlags().Int64("recipient")
	s := enter.GetSessionByAccountServerId(uid)
	if s == nil {
		ctx.Return(gmPlayerUnknown, fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
		return
	}
	if strDel := ctx.GetFlags().String("del"); strDel == "all" || strDel == "All" {
		if game.DelAllMail(s) {
			ctx.Return(cdq.ApiCodeOk, "删除玩家全部邮件成功")
		} else {
			ctx.Return(gmDelAllPlayerMailErr, "删除玩家全部邮件失败,这是一个神奇的bug")
		}
	} else if delId := ctx.GetFlags().Int64("del"); delId != 0 {
		if game.DelMail(s, delId) {
			ctx.Return(cdq.ApiCodeOk, "删除玩家邮件成功")
		} else {
			ctx.Return(gmDelPlayerMailErr, "删除玩家邮件失败,邮件可能不存在")
		}
	} else {
		parcelInfoList, err := genParcelInfo[sro.ParcelInfo](ctx.GetFlags().String("parcelInfoList"))
		if err != nil {
			ctx.Return(gmGenParcelInfoErr, fmt.Sprintf("解析邮件附件失败:%s", err.Error()))
			return
		}
		sendMail := &sro.MailInfo{
			Sender:         ctx.GetFlags().String("sender"),
			Comment:        ctx.GetFlags().String("comment"),
			SendDate:       ctx.GetFlags().Int64("sendDate"),
			ExpireDate:     ctx.GetFlags().Int64("expireDate"),
			ParcelInfoList: parcelInfoList,
		}
		if game.AddMail(s, sendMail) {
			ctx.Return(cdq.ApiCodeOk, "玩家邮件发送成功")
		} else {
			ctx.Return(gmSendPlayerMailErr, "玩家邮件发送失败")
		}
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
