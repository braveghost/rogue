package rogue

import (
	"sync"
)

type iSignal interface {
	Status() error
}

func NewHeartBeat(ts, dt int64) *HeartBeat {
	hb := &HeartBeat{
		lock:    sync.Mutex{},
		Counter: NewBucketCounter(ts, dt),
	}
	return hb
}

type HeartBeat struct {
	Counter *BucketCounter
	lock    sync.Mutex
}

func (hb *HeartBeat) Status() bool {
	hb.lock.Lock()
	defer hb.lock.Unlock()
	return hb.Counter.Overflow()
}

func (hb *HeartBeat) AddSignal(cf iSignal) error {
	return hb.disposeSignal(cf)
}

func (hb *HeartBeat) disposeSignal(s iSignal) error {
	hb.lock.Lock()
	defer hb.lock.Unlock()
	if err := s.Status(); err != nil {
		hb.Counter.Increment()
		return err
	}
	return nil
}
