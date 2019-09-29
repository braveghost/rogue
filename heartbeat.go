package rogue

import (
	"sync"
)

type iSignal interface {
	Status() bool
}

func NewHeartBeat(ts, dt int64) *HeartBeat {
	hb := &HeartBeat{
		lock:    sync.Mutex{},
		counter: NewBucketCounter(ts, dt),
	}
	return hb
}

type HeartBeat struct {
	counter *BucketCounter
	lock    sync.Mutex
}

func (hb *HeartBeat) Status() bool {
	hb.lock.Lock()
	defer hb.lock.Unlock()
	return hb.counter.Overflow()
}

func (hb *HeartBeat) AddSignal(cf iSignal) {
	hb.disposeSignal(cf)
}

func (hb *HeartBeat) disposeSignal(s iSignal) {
	hb.lock.Lock()
	defer hb.lock.Unlock()
	if !s.Status() {
		hb.counter.Increment()
	}
}
