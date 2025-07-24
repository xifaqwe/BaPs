package proto

import (
	"github.com/bytedance/sonic"
)

type MiniGameCCGStagePlayDB struct {
	StageId      int64
	EnemyGroupId int64
	IsClear      bool
	RewardDBs    []*MiniGameCCGStageRewardDB
	RerollPoint  int32
}

func (x *MiniGameCCGStagePlayDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
