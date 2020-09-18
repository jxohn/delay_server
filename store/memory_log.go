package store

import (
	"delay_server/model"
	"delay_server/wheel"
)

// MemoryLog 内存存储器
type MemoryLog struct {
	timeWheel wheel.HourWheel
}

func (m *MemoryLog) Append(message model.DelayMessage) error {
	return m.timeWheel.Put(message)
}

