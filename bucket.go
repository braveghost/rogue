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

func (bc *BucketCounter) daemon() {
	dlc := bc.deadlineCheck
	for {
		select {
		case <-bc.delCh:
			dlc()
		case <-bc.closeCh:
			return
		}
	}

}
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

func (bc *BucketCounter) deadlineCheck() {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()

	dtt := bc.diffTimestamp()
	for tt := range bc.Counters {
		if tt <= dtt {
			delete(bc.Counters, tt)
		}
	}
}

func (bc *BucketCounter) Increment() {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()

	b := bc.getCounter()
	b.Add()
	bc.deadlineCheck()
}

func (bc *BucketCounter) Size() int64 {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()

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
func (bc *BucketCounter) diffTimestamp() int64 {
	return time.Now().Unix() - bc.Duration

}
func (bc *BucketCounter) Overflow() bool {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	bc.deadlineCheck()
	if bc.Size() > bc.Threshold {
		return false
	}
	return true

}
func (bc *BucketCounter) Min() int64 {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	var (
		dtt = bc.diffTimestamp()
		min int64
	)

	for tt, ct := range bc.Counters {
		if tt >= dtt {
			if ! ct.Compare(min) {
				min = ct.Get()
			}
		}
	}
	bc.delCh <- &struct{}{}

	return min
}
func (bc *BucketCounter) Max() int64 {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	var (
		dtt = bc.diffTimestamp()
		max int64
	)

	for tt, ct := range bc.Counters {
		if tt >= dtt {
			if ct.Compare(max) {
				max = ct.Get()

			}
		}
	}

	bc.delCh <- &struct{}{}
	return max
}

func (bc *BucketCounter) Avg() int64 {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	return bc.Size() / bc.Duration
}

func (bc *BucketCounter) Clear() {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	bc.Counters = make(map[int64]*Counter, bc.Duration*2)

}

func (bc *BucketCounter) Close() {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()
	bc.closeCh <- &struct{}{}
}
