package rogue

import (
	"github.com/braveghost/joker"
	"sync"
)

type iSignal interface {
	Status() error
}

func NewHeartBeat(ts, dt int64, ct counterType, unit unitTimeType, name string) *HeartBeat {
	hb := &HeartBeat{
		lock:    sync.Mutex{},
		Counter: NewBucket(ts, dt, ct, unit, name),
	}
	return hb
}

type HeartBeat struct {
	Counter *Bucket
	lock    sync.Mutex
}

func (hb *HeartBeat) Status() (bool, error) {
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
		innerErr := hb.Counter.Increment()
		if innerErr != nil {
			logging.Errorw("Rogue.HeartBeat.DisposeSignal.Error", "inner_error", innerErr, "err", err)
			return innerErr
		}

		return err
	}
	return nil
}
