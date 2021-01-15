package delay_server

import (
	"sync"
	"time"
)

type (
	ring struct {
		sync.Mutex

		started bool
		// 当前时间轮的位置
		current int
		// 时间轮的大小, 在start方法输入的slotSize
		slotSize int
		// 当前执行的时间轮
		now *TimeWheel
		// 下一次轮换执行的时间轮
		next *TimeWheel
		// 通知预加载
		reloadChan chan struct{}

		HandlerFunc DataHandler
	}
	TimeWheel struct {
		// 存放数据, map的大小=slotSize
		data map[int]*InnerData
	}

	InnerData struct {
		sync.Mutex
		SlotData []DelayMessage
	}
)

// NewRing
func NewRing(slotSize int) *ring {
	r := &ring{
		Mutex:    sync.Mutex{},
		current:  0,
		slotSize: slotSize,
	}
	r.now = &TimeWheel{data: make(map[int]*InnerData)}
	r.next = &TimeWheel{data: make(map[int]*InnerData)}
	for i := 0; i < slotSize; i++ {
		r.now.data[i] = new(InnerData)
		r.next.data[i] = new(InnerData)
	}
	return r
}

// Start
func (r *ring) Start() {
	if r.started {
		return
	}
	r.Lock()
	if r.started {
		return
	}
	r.started = true
	defer r.Unlock()

	// 开启一些监听
	// 1. 监听reload
	go r.listen()

	// 开始定时处理
	go r.loop()
}

func (r *ring) loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		r.current += 1
		if r.current == r.slotSize {
			r.Lock()
			*r.now, *r.next = *r.next, *r.now
			r.reloadChan <- struct{}{}
			r.Unlock()
		}

		data := r.now.data[r.current]
		data.Lock()
		go r.HandlerFunc(data.SlotData)
		data.SlotData = nil
		data.Unlock()
	}
}

func (r *ring) listen() {
	for range r.reloadChan {
		// @1. reload, 加锁

		// @2. reload完毕, 加载reload过程中put阻塞的数据, 加锁

		// @3. 释放 put阻塞锁, 释放reload锁
	}
}

func (r *ring) PutOne(msg DelayMessage) {

	// @1. 计算放置位置
	delaySecond := msg.DelayTime - time.Now().Unix()
	slot := delaySecond / (60 * 60)
	index := delaySecond % (60 * 60)
	// @2. 执行|放置|丢弃
	if delaySecond <= 1 {
		// 立即执行
		r.HandlerFunc([]DelayMessage{msg})
	} else if slot <= 1 {
		var data *InnerData
		switch slot {
		case 0:
			// now
			data = r.now.data[int(index)]
		case 1:
			// next
			// todo 判断是否可写, 防止load和put同时进行, 导致放入了重复的消息
			// todo 如果不可写, 加锁写待处理msg arr
			data = r.next.data[int(index)]
		default:
			return
		}
		data.Lock()
		data.SlotData = append(data.SlotData, msg)
		data.Unlock()
	}

}
