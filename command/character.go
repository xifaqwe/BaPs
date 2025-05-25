package command

import (
	"errors"
	"fmt"
	"github.com/gucooing/BaPs/gdconf"
	"github.com/gucooing/BaPs/protocol/mx"
	"time"

	"github.com/gucooing/BaPs/common/enter"
	"github.com/gucooing/BaPs/game"
	"github.com/gucooing/BaPs/pkg/alg"
	"github.com/gucooing/cdq"
)

func (c *Command) ApplicationCommandCharacter() {
	character := &cdq.Command{
		Name:        "character",
		AliasList:   []string{"character", "c"},
		Description: "直接设置角色,角色不存在时将添加角色",
		Permissions: cdq.User,
		Options: []*cdq.CommandOption{
			{
				Name:        "uid",
				Description: "玩家游戏id",
				Required:    true,
			},
			{
				Name:        "id",
				Description: "角色id",
				Required:    true,
			},
			{
				Name:        "level",
				Description: "要设置的角色类型",
				Required:    false,
			},
			{
				Name:        "starGrade",
				Description: "要设置的角色星级",
				Required:    false,
			},
			{
				Name:        "favorRank",
				Description: "要设置的角色好感度等级",
				Required:    false,
			},
			{
				Name:        "max",
				Description: "角色属性设置满级",
				Required:    false,
			},
		},
		CommandFunc: c.character,
	}

	c.c.ApplicationCommand(character)
}

func (c *Command) character(options map[string]string) (string, error) {
	// 玩家验证
	uid := alg.S2I64(options["uid"])
	s := enter.GetSessionByAccountServerId(uid)
	if s == nil {
		return "", errors.New(fmt.Sprintf("玩家不在线或未注册 UID:%v", uid))
	}
	characterId := alg.S2I64(options["id"])
	if gdconf.GetCharacterExcel(characterId) == nil {
		return "", errors.New(fmt.Sprintf("角色不存在 CharacterId:%v", characterId))
	}
	respMsg := ""

	characterInfo := game.GetCharacterInfo(s, characterId)
	if characterInfo == nil {
		if game.AddCharacter(s, characterId) {
			characterInfo = game.GetCharacterInfo(s, characterId)
		} else {
			return "", errors.New("角色id错误")
		}
	}
	if alg.S2I64(options["max"]) != 0 {
		if game.SetMaxCharacter(characterInfo) {
			return "已设置角色满级", nil
		}
	}
	if level := alg.S2I32(options["level"]); level != 0 {
		if game.SetCharacterLevel(characterInfo, level) {
			respMsg += "角色等级设置成功|"
		}
	}
	if starGrade := alg.S2I32(options["starGrade"]); starGrade != 0 {
		if game.SetCharacterStarGrade(characterInfo, starGrade) {
			respMsg += "角色星级设置成功|"
		}
	}
	if favorRank := alg.S2I32(options["favorRank"]); favorRank != 0 {
		if game.SetCharacterFavorRank(characterInfo, favorRank) {
			respMsg += "角色好感度等级设置成功|"
		}
	}

	game.AddToast(s, &enter.Toast{
		Text:      respMsg,
		BeginDate: mx.Now(),
		EndDate:   mx.Now().Add(30 * time.Second),
	})

	return respMsg, nil
}
