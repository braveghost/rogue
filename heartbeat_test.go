package rogue

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestNewHeartBeat(t *testing.T) {
	hb := NewHeartBeat(5, 10,ECounterTypeRedis, UnitSecond,"test")
	hb.AddSignal(&SrvSignal{nil})
	hb.AddSignal(&SrvSignal{errors.New("")})
	hb.AddSignal(&SrvSignal{errors.New("")})
	hb.AddSignal(&SrvSignal{errors.New("")})
	hb.AddSignal(&SrvSignal{errors.New("")})
	hb.AddSignal(&SrvSignal{errors.New("")})
	hb.AddSignal(&SrvSignal{errors.New("")})
	hb.AddSignal(&SrvSignal{errors.New("")})
	hb.AddSignal(&SrvSignal{errors.New("")})
	hb.AddSignal(&SrvSignal{errors.New("")})
	time.Sleep(time.Second * 1)
	//hb.AddBeat(&SrvSignal{false})
	//hb.AddBeat(&SrvSignal{false})
	fmt.Println(hb.Status())
}

type SrvSignal struct {
	b error
}

// 心跳状态, 每一次就计算一次
func (hc *SrvSignal) Status() error {
	return hc.b
}
