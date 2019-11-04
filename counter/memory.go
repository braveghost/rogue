package counter

import (
	"sync/atomic"
)

func NewMemoryCounter(num int64) ICounter {
	return &Memory{num}
}

type Memory struct {
	num  int64
}

// 计数加1
func (c *Memory) Add() error {

	atomic.AddInt64(&c.num, 1)
	return nil
}

// 设置计数值
func (c *Memory) Set(n int64) error {

	atomic.StoreInt64(&c.num, n)
	return nil
}

// 重置计数器
func (c *Memory) Reset() error {
	atomic.AddInt64(&c.num, 0)
	return nil
}

// 计数比较
func (c *Memory) Compare(n int64) (bool, error) {
	if n > atomic.LoadInt64(&c.num) {
		return false, nil
	}
	return true, nil
}

// 获取当前计数长度
func (c *Memory) Get() (int64, error) {
	return atomic.LoadInt64(&c.num), nil
}

// 计数减1
func (c *Memory) Minus() error {
	atomic.AddInt64(&c.num, -1)
	return nil
}

// 清理该计数器及之前的
func (c *Memory) Clear() error {
	return c.Reset()
}
