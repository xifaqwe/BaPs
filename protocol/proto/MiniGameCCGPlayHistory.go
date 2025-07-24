package proto

import (
	"github.com/bytedance/sonic"
)

type MiniGameCCGPlayHistory struct {
	*ValueType
	LevelId     int64
	NodeId      int64
	StagePlayDB *MiniGameCCGStagePlayDB
}

func (x *MiniGameCCGPlayHistory) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
