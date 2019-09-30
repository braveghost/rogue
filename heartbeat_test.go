package rogue

import (
	"fmt"
	"testing"
	"time"
)

func TestNewHeartBeat(t *testing.T) {
	hb := NewHeartBeat(5, 10)
	hb.AddSignal(&SrvSignal{false})
	hb.AddSignal(&SrvSignal{true})
	hb.AddSignal(&SrvSignal{false})
	hb.AddSignal(&SrvSignal{false})
	hb.AddSignal(&SrvSignal{false})
	hb.AddSignal(&SrvSignal{false})
	hb.AddSignal(&SrvSignal{false})
	time.Sleep(time.Second * 11)
	//hb.AddBeat(&SrvSignal{false})
	//hb.AddBeat(&SrvSignal{false})
	fmt.Println(hb.Status())
}

type SrvSignal struct {
	b bool
}

// 心跳状态, 每一次就计算一次
func (hc *SrvSignal) Status() bool {
	return hc.b
}