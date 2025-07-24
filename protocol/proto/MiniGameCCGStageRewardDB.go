package proto

import (
	"github.com/bytedance/sonic"
)

type MiniGameCCGStageRewardDB struct {
	Type        MiniGameCCGStageRewardType
	RewardIndex int32
	RewardIds   *[]int64
}

func (x *MiniGameCCGStageRewardDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
