package wheel

import (
	"log"
	"sync"

	"delay_server/model"
)

type Processor func(message model.DelayMessage) error

// SecondWheel 秒级时间轮
type SecondWheel struct {
	sync.RWMutex

	data []model.DelayMessage // 该秒需要发送的消息

	processor Processor // 消息处理器
}

// 放入一个消息
func (s *SecondWheel) Put(message model.DelayMessage) {
	s.Lock()
	defer s.Unlock()
	log.Println("put in")
	s.data = append(s.data, message)
}

func (s *SecondWheel) GetAllMsg() (messages []model.DelayMessage) {
	s.RLock()
	defer s.RUnlock()
	return s.data
}

// DoSendAndClean 将该秒的所有消息发送出去
func (s *SecondWheel) DoSendAndClean() {
	s.Lock()
	defer s.Unlock()

	log.Println("now do")
	for _, v := range s.data {
		go func(msg model.DelayMessage) {
			err := s.processor(msg)
			log.Println(err)
		}(v)
	}

	// clean
	s.data = nil
}
