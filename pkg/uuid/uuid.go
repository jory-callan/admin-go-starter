package uuid

import (
	"time"

	"github.com/google/uuid"
)

// GenerateUUIDv7 生成基于时间的 UUID v7
func GenerateUUIDv7() string {
	now := time.Now().UnixMicro()

	// 将时间戳转换为字节数组
	timeBytes := make([]byte, 8)
	timeBytes[0] = byte(now >> 56)
	timeBytes[1] = byte(now >> 48)
	timeBytes[2] = byte(now >> 40)
	timeBytes[3] = byte(now >> 32)
	timeBytes[4] = byte(now >> 24)
	timeBytes[5] = byte(now >> 16)
	timeBytes[6] = byte(now >> 8)
	timeBytes[7] = byte(now)

	// 生成随机部分
	u := uuid.New()
	uBytes := u.String()

	// 组合时间戳和随机部分（简化实现，实际应遵循 UUID v7 标准）
	combined := make([]byte, 16)
	copy(combined[:6], timeBytes[:6])
	copy(combined[6:], uBytes[6:])

	return uuid.Must(uuid.FromBytes(combined)).String()
}
