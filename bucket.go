package rogue

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Counter struct {
	num int64
}

func (c *Counter) Add() {

	atomic.AddInt64(&c.num, 1)
}
func (c *Counter) Set(n int64) {

	atomic.StoreInt64(&c.num, n)
}

func (c *Counter) Reset() {
	atomic.AddInt64(&c.num, 0)
}

func (c *Counter) Compare(n int64) bool {
	if n > atomic.LoadInt64(&c.num) {
		return false
	}
	return true
}

func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.num)
}

func (c *Counter) Minus() {
	atomic.AddInt64(&c.num, -1)
}

func (c *Counter) String() string {
	return strconv.Itoa(int(c.num))
}

type BucketCounter struct {
	Counters  map[int64]*Counter
	Mutex     *sync.RWMutex
	Threshold int64
	Duration  int64
	delCh     chan *struct{}
	closeCh   chan *struct{}
}

// 创建计数器桶, ts为最大阈值, dt为统计的持续时间
func NewBucketCounter(ts, dt int64) *BucketCounter {
	r := &BucketCounter{
		Counters:  make(map[int64]*Counter, dt*2),
		Mutex:     &sync.RWMutex{},
		Threshold: ts,
		Duration:  dt,
		delCh:     make(chan *struct{}, 100),
		closeCh:   make(chan *struct{}, 1),
	}
	go r.daemon()
	return r
}

// 状态守护
func (bc *BucketCounter) daemon() {

	for {
		select {
		case <-bc.delCh:
			bc.delDeadline()
		case <-bc.closeCh:
			return
		}
	}
}

// 删除
func (bc *BucketCounter) delDeadline() {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()
	bc.deadlineCheck()
}

// 新建计数器
func (bc *BucketCounter) getCounter() *Counter {
	now := time.Now().Unix()
	var counter *Counter
	var ok bool

	if counter, ok = bc.Counters[now]; !ok {
		counter = &Counter{}
		bc.Counters[now] = counter
	}

	return counter
}

// 计数器ttl判断
func (bc *BucketCounter) deadlineCheck() {
	dtt := bc.diffTimestamp()
	for tt := range bc.Counters {
		if tt <= dtt {
			delete(bc.Counters, tt)
		}
	}
}

// 增量更新
func (bc *BucketCounter) Increment() {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()

	b := bc.getCounter()
	b.Add()
	bc.deadlineCheck()
}

// 计数器总计数
func (bc *BucketCounter) Sum() int64 {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	return bc.size()
}

// 总数
func (bc *BucketCounter) size() int64 {
	var (
		dtt = bc.diffTimestamp()
		sum int64
	)
	for tt, ct := range bc.Counters {
		if tt >= dtt {
			sum += ct.Get()
		}
	}
	bc.delCh <- &struct{}{}
	return sum
}

// 当前时间差, 起点时间
func (bc *BucketCounter) diffTimestamp() int64 {
	return time.Now().Unix() - bc.Duration
}

// 锁定状态
func (bc *BucketCounter) Overflow() bool {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()

	bc.deadlineCheck()
	if bc.size() > bc.Threshold {
		return true
	}
	return false
}

// 秒内最小值
func (bc *BucketCounter) Min() int64 {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	var (
		dtt  = bc.diffTimestamp()
		min  int64
		flag bool
	)

	for tt, ct := range bc.Counters {
		if tt >= dtt {

			if ! flag {
				min = ct.Get()
				flag = true

				continue
			}
			if ! ct.Compare(min) {
				min = ct.Get()
			}
		}
	}
	bc.delCh <- &struct{}{}

	return min
}

// 秒内最大值
func (bc *BucketCounter) Max() int64 {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	var (
		dtt  = bc.diffTimestamp()
		max  int64
		flag bool
	)
	for tt, ct := range bc.Counters {
		if tt >= dtt {
			if ! flag {
				max = ct.Get()
				flag = true
				continue
			}
			if ct.Compare(max) {
				max = ct.Get()

			}
		}
	}

	bc.delCh <- &struct{}{}
	return max
}

// 均值
func (bc *BucketCounter) Avg() float64 {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	return float64(bc.size()) / float64(len(bc.Counters))
}

// 重置计数器桶
func (bc *BucketCounter) String() string {
	return strconv.Itoa(int(bc.Sum()))

}

// 重置计数器桶
func (bc *BucketCounter) Clear() {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	bc.Counters = make(map[int64]*Counter, bc.Duration*2)

}

// 关闭计数器桶
func (bc *BucketCounter) Close() {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()
	bc.closeCh <- &struct{}{}
}
