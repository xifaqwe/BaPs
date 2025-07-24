package proto

import (
	"github.com/bytedance/sonic"
)

type MiniGameCCGSaveDB struct {
	EventContentId       int64
	CCGId                int64
	GameOver             bool
	Clear                bool
	Strikers             []*MiniGameCCGCharacterDB
	Specials             []*MiniGameCCGCharacterDB
	OverflowedStrikerIds []int64
	OverflowedSpecialIds []int64
	Deck                 []*MiniGameCCGCardDB
	LevelId              int64
	CurrentNodeId        int64
	RewardPoint          int32
	CurrentStageDB       *MiniGameCCGStagePlayDB
	ClearedHistoryDBs    []*MiniGameCCGPlayHistory
	Perks                []int64
	TotalUsedCost        int32
	TotalDamageAmount    int32
	TotalKillCount       int32
	TotalSkillCount      map[int64]int32
	CurrentStageIndex    int32
	LastStageIndex       int32
}

func (x *MiniGameCCGSaveDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
