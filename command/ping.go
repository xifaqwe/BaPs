package command

import (
	"fmt"
	"github.com/bytedance/sonic"
	"runtime"
	"sync/atomic"

	"github.com/gucooing/BaPs/common/check"
	"github.com/gucooing/BaPs/pkg"
	"github.com/gucooing/cdq"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type Ping struct {
	PlayerNum     int64       `json:"playerNum"`     // 在线玩家数量
	Tps           int64       `json:"tps"`           // 上一分钟请求量
	Rt            FloatString `json:"rt"`            // 上一分钟每一个请求平均处理时间 单位ms
	ClientVersion string      `json:"clientVersion"` // 客户端版本
	ServerVersion string      `json:"serverVersion"` // 服务端版本
	CpuOc         float64     `json:"cpuOc"`         // cpu占用
	MemoryOc      string      `json:"memoryOc"`      // 内存占用
	BaPsMemoryOc  string      `json:"baPsMemoryOc"`  // BaPs内存占用
}

func (c *Command) ApplicationCommandPing() {
	ping := &cdq.Command{
		Name:        "ping",
		AliasList:   []string{"ping"},
		Description: "检查服务是否存活",
		Permissions: cdq.Guest,
		Options:     nil,
		CommandFunc: c.ping,
	}
	c.c.ApplicationCommand(ping)
}

func (c *Command) ping(options map[string]*cdq.CommandOption) (string, error) {
	response := Ping{
		PlayerNum:     atomic.LoadInt64(&check.SessionNum),
		Tps:           check.OLDTPS,
		Rt:            FloatString(check.OLDRT),
		ClientVersion: pkg.ClientVersion,
		ServerVersion: pkg.ServerVersion,
		CpuOc:         GetCpuOc(),
		MemoryOc:      MemoryOc(),
		BaPsMemoryOc:  BaPsMemoryOc(),
	}
	bin, err := sonic.MarshalString(response)
	return bin, err
}

func GetCpuOc() float64 {
	percents, err := cpu.Percent(0, false)
	if err != nil {
		return 0
	}
	return percents[0]
}

func MemoryOc() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "0/0"
	}
	used := float64(memInfo.Used)
	total := float64(memInfo.Total)

	if used/1024/1024 > 1024 {
		return fmt.Sprintf("%.2f/%.2fGB", used/1024/1024/1024, total/1024/1024/1024)
	}
	return fmt.Sprintf("%.2f/%.2fMB", used/1024/1024, total/1024/1024)
}

func BaPsMemoryOc() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	used := float64(m.Alloc)
	if used/1024/1024 > 1024 {
		return fmt.Sprintf("%.2fGB", used/1024/1024/1024)
	}
	return fmt.Sprintf("%.2fMB", used/1024/1024)
}

type FloatString float64

func (f FloatString) MarshalJSON() ([]byte, error) {
	// 将浮点数格式化为字符串，保留必要精度
	return []byte(fmt.Sprintf(`"%v"`, float64(f))), nil
}
