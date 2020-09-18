package biz

import "delay_server/model"

type DefaultProcessor struct {

}

func (d *DefaultProcessor) Process(message model.DelayMessage) error {
	return nil
}

