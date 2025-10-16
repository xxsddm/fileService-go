package utils

import (
	"time"
)

// IDGenerator ID生成器（简单的实现）
type IDGenerator struct {
	lastTimestamp int64
	sequence      int64
}

var generator = &IDGenerator{}

// GenerateID 生成唯一ID
func GenerateID() uint64 {
	now := time.Now().UnixNano() / 1000000 // 毫秒时间戳
	
	// 简单的ID生成逻辑
	if now == generator.lastTimestamp {
		generator.sequence++
	} else {
		generator.sequence = 0
		generator.lastTimestamp = now
	}
	
	// 组合时间戳和序列号生成ID
	id := (uint64(now) << 22) | (uint64(generator.sequence) & 0x3FFFFF)
	
	return id
}