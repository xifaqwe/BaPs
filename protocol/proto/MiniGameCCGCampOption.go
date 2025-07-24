package proto

type MiniGameCCGCampOption int32

const (
	MiniGameCCGCampOption_Invalid    MiniGameCCGCampOption = 0
	MiniGameCCGCampOption_Heal       MiniGameCCGCampOption = 1
	MiniGameCCGCampOption_Revive     MiniGameCCGCampOption = 2
	MiniGameCCGCampOption_RemoveCard MiniGameCCGCampOption = 3
	MiniGameCCGCampOption_Skip       MiniGameCCGCampOption = 4
)

var (
	MiniGameCCGCampOption_name = map[int32]string{
		0: "Invalid",
		1: "Heal",
		2: "Revive",
		3: "RemoveCard",
		4: "Skip",
	}
	MiniGameCCGCampOption_value = map[string]int32{
		"Invalid":    0,
		"Heal":       1,
		"Revive":     2,
		"RemoveCard": 3,
		"Skip":       4,
	}
)

func (x MiniGameCCGCampOption) String() string {
	return MiniGameCCGCampOption_name[int32(x)]
}

func (x MiniGameCCGCampOption) Value(sr string) MiniGameCCGCampOption {
	return MiniGameCCGCampOption(MiniGameCCGCampOption_value[sr])
}
