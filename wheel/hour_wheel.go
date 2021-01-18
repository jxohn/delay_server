package wheel

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"delay_server/model"
	"delay_server/tool"
)

const (
	WheelsLength = 0b1111111111111
	MinuteBit    = 0b111111000000
	SecondBit    = 0b000000111111
	NextHourBit  = 0b0
)

const (
	Locked = iota + 1
	UnLocked
)

type HourWheel struct {
	sync.RWMutex

	cursor       *uint32 // buffer游标, 用来标记当前小时和下一个小时的数据存储位
	CurrentIndex *uint64 // 当前位置

	CurrentHour time.Time // 当前小时

	tickDuration time.Duration // 秒级计时器

	mWheels []*SecondWheel

	notifyChan chan model.MessageOffset // 通知加载数据, 小时级别数据, 处理时需要加锁, 防止同时从文件和内存写入
	offsetChan chan model.MessageOffset // 更新消费offset

	loaderLock sync.RWMutex // loader锁

	processor Processor
}

type LockEvent int32

// NewHourWheel 新建一个时间轮
// @param notifyChan: 小时级别通知, 加在数据; offsetChan : 秒级别的通知, 更新处理索引位置
// @WARN : 传入的notifyChan和offsetChan必须由传入方去消费
func NewHourWheel(nc chan model.MessageOffset, oc chan model.MessageOffset, p Processor) *HourWheel {
	return &HourWheel{
		RWMutex:      sync.RWMutex{},
		cursor:       new(uint32),
		CurrentIndex: nil,
		CurrentHour:  time.Time{},
		tickDuration: time.Second,
		mWheels:      make([]*SecondWheel, 1<<13),
		notifyChan:   nc,
		offsetChan:   oc,
		processor:    p,
	}
}

// Start 开始处理
func (h *HourWheel) Start() {
	go h.listen()

	h.BeginLoop()
}

func (h *HourWheel) BeginLoop() {
	tick := time.NewTicker(h.tickDuration)
	defer tick.Stop()

	currentIndex := new(uint64)
	now := time.Now()
	*currentIndex = h.GetIndexFromTime(now.Add(time.Second))
	h.CurrentIndex = currentIndex
	h.CurrentHour = now
	for {
		select {
		case <-tick.C:
			now := time.Now()
			nowIndex := h.GetIndexFromTime(now)
			nextIndex := h.GetIndexFromTime(now.Add(time.Second))

			// @WARN : ntp时间同步可能会导致时间跳跃
			if !atomic.CompareAndSwapUint64(h.CurrentIndex, nowIndex, nextIndex) {
				// 由于ntp时间同步问题, 此处不能panic, 对于回拨的情况, 不做处理, 对于跳跃的情况, 需要把之前的数据都处理掉
				for i := *h.CurrentIndex; i < nowIndex; i++ {
					h.mWheels[(*h.cursor<<12+uint32(i))&WheelsLength].DoSendAndClean()
				}
				atomic.SwapUint64(h.CurrentIndex, nextIndex)
				// panic("index check failed, should be : " + strconv.FormatUint(nowIndex, 10) + " but actual " + strconv.FormatUint(*h.CurrentIndex, 10))
			}
			log.Printf("%d index", (*h.cursor<<12+uint32(nowIndex))&WheelsLength)
			h.mWheels[(*h.cursor<<12+uint32(nowIndex))&WheelsLength].DoSendAndClean()

			// 下一个小时
			if (nextIndex ^ NextHourBit) == NextHourBit {
				h.CurrentHour = now.Add(time.Second)
				log.Printf("next..., %s", h.CurrentHour.String())
				// 1. 换游标
				atomic.SwapUint32(h.cursor, 1^(*h.cursor))
				// 2. 通知加载数据
				h.notifyChan <- model.MessageOffset{
					FileName: tool.BuildFileNameByTime(h.CurrentHour),
					Index:    0,
				}
			}
		}
	}
}

// Put 放入一个延迟消息, 无需外围计算位置, 隔离处理和加载数据线程
func (h *HourWheel) Put(message model.DelayMessage) error {
	// 1. 计算延迟时间
	// 1.1 如果超出时间, 则抛弃
	now := time.Now()
	hour := message.DelayTime.Hour() - now.Hour()
	if hour > 1 {
		log.Printf("msg delay too long, not now , time is %s", message.DelayTime.String())
		return nil
	}

	// 1.2 如果放入时间小于或等于当前时间, 直接执行
	if now.Add(time.Second).After(message.DelayTime) {
		go h.processor(message)
		return nil
	}

	// 1.3 根据延迟时间计算位置
	cursorIndex := hour ^ int(*h.cursor)
	waterIndex := h.GetIndexFromTime(message.DelayTime)

	// todo 如果cursor在next, 需要判断next是否被file loader锁住, 如果锁住, 则放入等待队列...
	index := (cursorIndex<<12 + int(waterIndex)) & WheelsLength
	// 下一小时
	if (hour & int(*h.cursor)) == 1 {
		h.loaderLock.RLock()
		defer h.loaderLock.RUnlock()
	}

	// 2. 放入计算位置bucket
	h.mWheels[index].Put(message)
	return nil
}

// GetIndexFromTime 根据出传入的时间, 计算该时间的index
func (h *HourWheel) GetIndexFromTime(now time.Time) uint64 {
	minute := now.Minute()
	second := now.Second()
	return uint64((minute<<6)&MinuteBit+second&SecondBit) & WheelsLength
}

func (h *HourWheel) listen() {
	for {
		mo := <-h.notifyChan
		// 1. 加锁
		h.loaderLock.Lock()
		// 2. 加载数据
		// todo loader
		mo.Index
		// 3. 释放锁
		h.loaderLock.Unlock()
	}
}
