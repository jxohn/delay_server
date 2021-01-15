package delay_server

import (
	"fmt"
)

type (
	// delayServer 延迟server端
	delayServer struct {
		storagePath string

		ring *ring
	}
	// SOption 延迟server端配置
	SOption     func(*delayServer)
	DataHandler func(messages []DelayMessage)
)

func NewDelayServer(opt []SOption) (*delayServer, error) {
	server := new(delayServer)

	for _, f := range opt {
		f(server)
	}

	err := server.validate()
	if err != nil {
		return nil, fmt.Errorf("validate server error, error is %+v", err)
	}

	ring := NewRing(60 * 60)
	server.ring = ring

	return server, nil
}

// WithStoragePath 存储本地文件地址
func (s *delayServer) WithStoragePath(path string) SOption {
	return func(s *delayServer) {
		s.storagePath = path
	}
}

// PutMsg 收到一条延迟消息
func (s *delayServer) PutMsg(msg DelayMessage) error {
	s.ring.PutOne(msg)
	// todo persist
	return nil
}

// validate server各参数校验
func (s *delayServer) validate() error {
	return nil
}
