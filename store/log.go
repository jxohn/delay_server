package store

import "delay_server/model"

// LogAppender 为消息存储器, 分为内存和磁盘两种
// 内存写需要在磁盘写的前面, 因为内存写需要加RLock(读磁盘)
type LogAppender interface {
	Append(message model.DelayMessage) error
}
