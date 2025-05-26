package command

import (
	"errors"
	"fmt"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	sro "github.com/gucooing/BaPs/common/server_only"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/BaPs/protocol/proto"
	"github.com/gucooing/cdq"
)

type ApiGiveAll struct {
	Uid  int64  `json:"uid"`
	Type string `json:"type"`
	Num  int64  `json:"num"`
}

type giveAllFunc func(req *ApiGiveAll) []*sro.ParcelInfo

var GiveAllType = map[string]giveAllFunc{
	"Material":    GiveAllMaterial,
	"Character":   GiveAllCharacter,
	"Equipment":   GiveAllEquipment,
	"Furniture":   GiveAllFurniture,
	"Favor":       GiveAllFavor,
	"Emblem":      GiveAllEmblem,
	"Sticker":     GiveAllSticker,
	"MemoryLobby": GiveAllMemoryLobby,
}

func (c *Command) ApplicationCommandGiveAll() {
	giveAll := &cdq.Command{
		Name:        "giveAll",
		AliasList:   []string{"giveAll", "ga"},
		Description: "获取某个类型的全部物品",
		Permissions: cdq.User,
		Options: []*cdq.CommandOption{
			{
				Name:        "uid",
				Description: "玩家游戏id",
				Required:    true,
			},
			{
				Name:        "t",
				Description: "需要给予物品的类型",
				Required:    true,
			},
			{
				Name:        "num",
				Description: "需要给予物品的数量 默认值:1",
				Required:    false,
			},
		},
		CommandFunc: c.giveALL,
	}

	c.C.ApplicationCommand(giveAll)
}

func (c *Command) giveALL(options map[string]string) (string, error) {
	// 玩家验证
	uid := alg.S2I64(options["uid"])
	s := enter.GetSessionByAccountServerId(uid)
	if s == nil {
		return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
	}

	// 执行
	parcelInfoList := make([]*sro.ParcelInfo, 0)
	parcelInfoList = GiveAllJsonToProtobuf(&ApiGiveAll{
		Uid:  uid,
		Type: options["t"],
		Num:  alg.MaxInt64(alg.S2I64(options["num"]), 1),
	})

	if len(parcelInfoList) == 0 {
		return "", errors.New(fmt.Sprintf("不存在此物品类型 Type:%s", options["t"]))
	}

	mail := &sro.MailInfo{
		Sender:         "gucooing",
		Comment:        "请查收您的意外奖励",
		SendDate:       time.Now().Unix(),
		ExpireDate:     time.Now().Add(10 * time.Minute).Unix(),
		ParcelInfoList: parcelInfoList,
	}
	if game.AddMail(s, mail) {
		return "请查询游戏内邮箱获取结果", nil
	}
	return "", errors.New("游戏邮箱错误")
}

func GiveAllJsonToProtobuf(req *ApiGiveAll) []*sro.ParcelInfo {
	if req.Type == "All" {
		return GiveAllTypes(req)
	}
	v, ok := GiveAllType[req.Type]
	if ok {
		return v(req)
	}
	return nil
}

func GiveAllTypes(req *ApiGiveAll) []*sro.ParcelInfo {
	parcelInfoList := make([]*sro.ParcelInfo, 0)
	for _, v := range GiveAllType {
		parcelInfoList = append(parcelInfoList, v(req)...)
	}

	return parcelInfoList
}

func GiveAllMaterial(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	num := alg.MinInt64(req.Num, 9999)
	for _, conf := range gdconf.GetItemExcelCategoryMapByCategory("Material") {
		list = append(list, &sro.ParcelInfo{
			Type: int32(proto.ParcelType_Item),
			Id:   conf.Id,
			Num:  num,
		})
	}

	return list
}

func GiveAllCharacter(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetCharacterMap() {
		list = append(list, &sro.ParcelInfo{
			Type: int32(proto.ParcelType_Character),
			Id:   conf.Id,
			Num:  1,
		})
	}
	return list
}

func GiveAllEquipment(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetEquipmentExcelMap() {
		num := req.Num
		if conf.MaxLevel < 10 {
			num = alg.MinInt64(req.Num, 9999)
		} else {
			num = alg.MinInt64(req.Num, 20)
		}

		list = append(list, &sro.ParcelInfo{
			Type: int32(proto.ParcelType_Equipment),
			Id:   conf.Id,
			Num:  num,
		})
	}

	return list
}

func GiveAllFurniture(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetFurnitureExcelTableMap() {
		num := req.Num
		num = alg.MinInt64(req.Num, 3)
		list = append(list, &sro.ParcelInfo{
			Type: int32(proto.ParcelType_Furniture),
			Id:   conf.Id,
			Num:  num,
		})
	}

	return list
}

func GiveAllFavor(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	num := req.Num
	num = alg.MinInt64(req.Num, 9999)
	for _, conf := range gdconf.GetItemExcelCategoryMapByCategory("Favor") {
		list = append(list, &sro.ParcelInfo{
			Type: int32(proto.ParcelType_Item),
			Id:   conf.Id,
			Num:  num,
		})
	}

	return list
}

func GiveAllEmblem(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetEmblemExcelList() {
		list = append(list, &sro.ParcelInfo{
			Type: int32(proto.ParcelType_Emblem),
			Id:   conf.Id,
		})
	}
	return list
}

func GiveAllSticker(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetStickerPageContentExcelList() {
		list = append(list, &sro.ParcelInfo{
			Type: int32(proto.ParcelType_Sticker),
			Id:   conf.Id,
			Num:  1,
		})
	}
	return list
}

func GiveAllMemoryLobby(req *ApiGiveAll) []*sro.ParcelInfo {
	list := make([]*sro.ParcelInfo, 0)
	for _, conf := range gdconf.GetMemoryLobbyExcelList() {
		list = append(list, &sro.ParcelInfo{
			Type: int32(proto.ParcelType_MemoryLobby),
			Id:   conf.Id,
			Num:  1,
		})
	}
	return list
}
