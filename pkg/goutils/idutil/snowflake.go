package idutil

import (
	"sync/atomic"
	"time"

	"github.com/sony/sonyflake"
)

// 原子指针直接存 *sonyflake.Sonyflake，避免接口转换。
var currentSF atomic.Pointer[sonyflake.Sonyflake]

func init() {
	ResetSnowflake(sonyflake.Settings{})
}

// ResetSnowflake 运行时热替换雪花生成器。
func ResetSnowflake(settings sonyflake.Settings) {
	sf := sonyflake.NewSonyflake(settings)
	if sf == nil {
		// 通常因为 StartTime 在未来，直接 panic 比默默 nil 更安全。
		panic("sonyflake: invalid settings")
	}
	currentSF.Store(sf)
}

// SnowflakeID 获取全局唯一 64 bit ID。
// 时钟回拨时内部重试 3 次，再失败即 panic（防止重复 ID）。
func SnowflakeID() uint64 {
	sf := currentSF.Load()
	// 重试 3 次
	for i := 0; i < 3; i++ {
		id, err := sf.NextID()
		if err == nil {
			return id
		}
		// 生产环境可替换为日志 + 告警
		time.Sleep(time.Millisecond)
	}
	panic("sonyflake: clock backwards too long")
}
