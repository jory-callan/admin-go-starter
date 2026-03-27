package idutil

import (
	"strings"

	"github.com/gofrs/uuid/v5"
)

// 全局复用同一个生成器，避免每次 new。
var defaultGen = uuid.NewGen()

// UUID 返回标准 UUIDv7 字符串。
func UUIDv7() string {
	return uuid.Must(defaultGen.NewV7()).String()
}

// ShortUUIDv7 返回去掉连字符的 UUIDv7。
func ShortUUIDv7() string {
	return strings.ReplaceAll(UUIDv7(), "-", "")
}
