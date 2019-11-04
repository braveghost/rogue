package rogue

import (
	"github.com/braveghost/joker"
	"github.com/braveghost/meteor/itime"
	"github.com/braveghost/rogue/counter"
	"sync"
	"time"
)

type unitTimeType string

const (
	UnitZero   unitTimeType = "zero"
	UnitSecond unitTimeType = "second"
	UnitMinute unitTimeType = "minute"
	UnitHour   unitTimeType = "hour"
	UnitDay    unitTimeType = "day"
	UnitWeek   unitTimeType = "week"
	UnitMonth  unitTimeType = "month"
	UnitYear   unitTimeType = "year"
)

type counterType string

const (
	ECounterTypeMemory counterType = "memory"
	ECounterTypeRedis  counterType = "redis"
)

var (
	defaultThreshold = int64(5)
	defaultDuration  = int64(10)
)

// 创建计数器桶, ts为最大阈值, dt为统计的持续时间
func NewBucket(ts, dt int64, ct counterType, unit unitTimeType) *Bucket {
	if ts <= 0 {
		ts = defaultThreshold
	}
	if dt <= 0 {
		dt = defaultDuration
	}

	r := &Bucket{
		Type:      ct,
		Unit:      unit,
		Counters:  make(map[int64]counter.ICounter, dt*2),
		Mutex:     &sync.RWMutex{},
		Threshold: ts,
		Duration:  dt,
		delCh:     make(chan *struct{}, 100),
		closeCh:   make(chan *struct{}, 1),
	}
	return r
}

type Bucket struct {
	Type      counterType
	Unit      unitTimeType
	Counters  map[int64]counter.ICounter
	Mutex     *sync.RWMutex
	Threshold int64
	Duration  int64 // 持续时长, duration * unit, 最低单位秒
	delCh     chan *struct{}
	closeCh   chan *struct{}
}

// 删除
func (bc *Bucket) delDeadline() {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()
	bc.deadlineCheck()
}

// 新建计数器
func (bc *Bucket) getCounter() counter.ICounter {
	now := bc.nowTimestamp()
	var ct counter.ICounter
	var ok bool

	if ct, ok = bc.Counters[now]; !ok {
		ct = bc.newCounter(now)
		bc.Counters[now] = ct
	}

	return ct
}

func (bc *Bucket) newCounter(now int64) counter.ICounter {
	switch bc.Type {
	case ECounterTypeRedis:
		return counter.NewRedisCounter(now)
	case ECounterTypeMemory:
		return counter.NewMemoryCounter(0)
	}
	return nil
}

// 计数器 ttl 判断
func (bc *Bucket) deadlineCheck() {
	dtt := bc.diffTimestamp()
	for tt, vv := range bc.Counters {
		if tt <= dtt {
			_ = vv.Clear()
			delete(bc.Counters, tt)
		}
	}
}

// 增量更新
func (bc *Bucket) Increment() error {
	bc.Mutex.Lock()
	defer bc.Mutex.Unlock()

	b := bc.getCounter()
	err := b.Add()
	if err != nil {
		return err
	}
	bc.deadlineCheck()

	return nil
}

// 计数器总计数
func (bc *Bucket) Sum() (int64, error) {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()

	return bc.size()
}

// 总数
func (bc *Bucket) size() (int64, error) {
	bc.deadlineCheck()
	var (
		dtt = bc.diffTimestamp()
		sum int64
	)
	for tt, ct := range bc.Counters {
		if tt >= dtt {
			val, err := ct.Get()
			if err != nil {
				return 0, err
			}
			sum += val
		}
	}
	return sum, nil
}

// 当前时间差, 起点时间
func (bc *Bucket) diffTimestamp() int64 {

	switch bc.Unit {
	case UnitSecond:
		return time.Now().Unix() - bc.Duration
	case UnitMinute:
		return itime.LastMinutesStart(bc.Duration)
	case UnitHour:
		return itime.LastHoursStart(bc.Duration)
	case UnitDay:
		return itime.LastDaysStart(bc.Duration)
	case UnitMonth:
		return itime.LastMonthsStart(bc.Duration)
	case UnitYear:
		return itime.LastYearsStart(bc.Duration)
	default:
		joker.Warnf("not support unit type '%s'", bc.Unit)
		return 0
	}
}


// 当前时间差, 起点时间
func (bc *Bucket) nowTimestamp() int64 {

	switch bc.Unit {
	case UnitSecond:
		return time.Now().Unix()
	case UnitMinute:
		return itime.NowMinuteStart()
	case UnitHour:
		return itime.NowHourStart()
	case UnitDay:
		return itime.NowDayStart()
	case UnitMonth:
		return itime.NowMonthStart()
	case UnitYear:
		return itime.NowYearStart()
	default:
		joker.Warnf("not support unit type '%s'", bc.Unit)
		return 0
	}
}


// 锁定状态
func (bc *Bucket) Overflow() (bool, error) {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	size, err := bc.size()
	if err != nil {
		return false, err
	}
	if size > bc.Threshold {
		return true, nil
	}
	return false, nil
}

// 秒内最小值
func (bc *Bucket) Min() (int64, error) {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	var (
		dtt  = bc.diffTimestamp()
		min  int64
		flag bool
		ok   bool
		err  error
	)
	bc.deadlineCheck()
	for tt, ct := range bc.Counters {
		if tt >= dtt {

			if ! flag {
				min, err = ct.Get()
				if err != nil {
					return 0, err
				}

				flag = true
				continue
			}
			ok, err = ct.Compare(min)
			if ok {
				min, err = ct.Get()
				if err != nil {
					return 0, err
				}
			}
		}
	}

	return min, nil
}

// 秒内最大值
func (bc *Bucket) Max() (int64, error) {

	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	var (
		dtt  = bc.diffTimestamp()
		max  int64
		flag bool
		ok   bool
		err  error
	)
	bc.deadlineCheck()

	for tt, ct := range bc.Counters {
		if tt >= dtt {
			if ! flag {
				max, err = ct.Get()
				if err != nil {
					return 0, err
				}
				flag = true
				continue
			}
			ok, err = ct.Compare(max)
			if ok {
				max, err = ct.Get()
				if err != nil {
					return 0, err
				}
			}
		}
	}
	return max, nil
}

// 均值
func (bc *Bucket) Avg() (float64, error) {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	val, err := bc.size()
	if err != nil {
		return 0, err
	}
	return float64(val) / float64(len(bc.Counters)), nil
}

// 重置计数器桶
func (bc *Bucket) Clear() {
	bc.Mutex.RLock()
	defer bc.Mutex.RUnlock()
	for _, vv := range bc.Counters {
		_ = vv.Clear()
	}
	bc.Counters = make(map[int64]counter.ICounter, bc.Duration*2)

}
