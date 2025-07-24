package proto

import (
	"github.com/bytedance/sonic"
)

type MiniGameCCGCardDB struct {
	CardDBId int32
	CardId   int64
}

func (x *MiniGameCCGCardDB) String() string {
	j, _ := sonic.MarshalString(x)
	return j
}
