package proto

import (
	"github.com/bytedance/sonic"
)

type MiniGameCCGCharacterDB struct {
	SlotIndex   int32
	CharacterId int64
	Health      int32
}

func (x *MiniGameCCGCharacterDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
