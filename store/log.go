package store

import "delay_server/model"

// LogAppender 为消息存储器, 分为内存和磁盘两种
type LogAppender interface {
	Append(message model.DelayMessage) error
}
