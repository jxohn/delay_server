package model

import (
	"time"

	"github.com/Shopify/sarama"
)

type DelayMessage struct {
	DelayTime time.Time
	Msg       sarama.ProducerMessage // kafka消息

	Brokers []string
}

type MessageOffset struct {
	FileName string // 文件名
	Index    uint64 // 秒钟位置
}
