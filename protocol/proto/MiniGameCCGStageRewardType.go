package proto

type MiniGameCCGStageRewardType int32

const (
	MiniGameCCGStageRewardType_Invalid MiniGameCCGStageRewardType = 0
	MiniGameCCGStageRewardType_All     MiniGameCCGStageRewardType = 1
	MiniGameCCGStageRewardType_Select  MiniGameCCGStageRewardType = 2
)

var (
	MiniGameCCGStageRewardType_name = map[int32]string{
		0: "Invalid",
		1: "All",
		2: "Select",
	}
	MiniGameCCGStageRewardType_value = map[string]int32{
		"Invalid": 0,
		"All":     1,
		"Select":  2,
	}
)

func (x MiniGameCCGStageRewardType) String() string {
	return MiniGameCCGStageRewardType_name[int32(x)]
}

func (x MiniGameCCGStageRewardType) Value(sr string) MiniGameCCGStageRewardType {
	return MiniGameCCGStageRewardType(MiniGameCCGStageRewardType_value[sr])
}
