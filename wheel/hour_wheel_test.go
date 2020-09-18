package wheel

import (
	"log"
	"sync"
	"testing"
	"time"

	"delay_server/model"

	"github.com/Shopify/sarama"
)

func TestHourWheel_BeginLoop(t *testing.T) {
	nChan := make(chan model.MessageOffset, 1)
	wheels := make([]*SecondWheel, 1<<13)
	for i := range wheels {
		wheels[i] = &SecondWheel{
			RWMutex: sync.RWMutex{},
			data:    nil,
			processor: func(message model.DelayMessage) error {
				log.Println("dead")
				return nil
			},
		}
	}

	c := new(uint32)
	*c = 1
	wheel := &HourWheel{
		RWMutex:      sync.RWMutex{},
		cursor:       c,
		CurrentIndex: nil,
		CurrentHour:  time.Time{},
		tickDuration: time.Second,
		mWheels:      wheels,
		notifyChan:   nChan,
		processor: func(message model.DelayMessage) error {
			log.Printf("deal, %s", message.DelayTime.String())
			return nil
		},
	}
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				msg := model.DelayMessage{
					DelayTime: time.Now().Add(10 * time.Second),
					Msg:       sarama.ProducerMessage{},
				}
				wheel.Put(msg)
			}
		}
	}()

	go func() {
		for {
			select {
			case a := <-nChan:
				log.Println(a.FileName)
			}
		}
	}()

	wheel.BeginLoop()
}
