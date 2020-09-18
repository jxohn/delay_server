package store

import "delay_server/model"

// DiskLog 磁盘存储器
type DiskLog struct {

}

func (d *DiskLog) Append(message model.DelayMessage) error {
	panic("implement me")
}

