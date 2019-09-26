package rogue

import (
	"sync"
)

type iSignal interface {
	Status() bool
}

func NewHeartBeat(ts, dt int64) *HeartBeat {
	hb := &HeartBeat{
		signalCh: make(chan iSignal, 100),
		lock:     sync.Mutex{},
		counter:  NewBucketCounter(ts, dt),
		closeCh:  make(chan *struct{}),
	}
	return hb
}

type HeartBeat struct {
	counter  *BucketCounter
	status   bool
	lock     sync.Mutex
	signalCh chan iSignal
	closeCh  chan *struct{}
}

func (hb *HeartBeat) Status() bool {
	hb.lock.Lock()
	defer hb.lock.Unlock()
	return hb.status
}
func (hb *HeartBeat) changeStatus(bl bool) {
	hb.lock.Lock()
	defer hb.lock.Unlock()
	hb.status = bl
}

func (hb *HeartBeat) AddBeat(cf iSignal) {
	hb.signalCh <- cf
}

func (hb *HeartBeat) run() {
	counter := hb.counter
	signalCh := hb.signalCh
	for {
		select {
		case h := <-signalCh:

			if !h.Status() {
				counter.Increment()
			}

			if counter.Overflow() {
				hb.changeStatus(true)
			}

		}

	}
}

func (hb *HeartBeat) Run() {
	go hb.run()
}
