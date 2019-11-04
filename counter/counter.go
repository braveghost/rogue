package counter

type ICounter interface {
	// 计数加1
	Add() error

	// 计数减1
	Minus() error

	// 设置计数值
	Set(n int64) error

	// 重置计数器
	Reset() error

	// 计数比较
	Compare(int64) (bool, error)

	// 获取当前计数长度
	Get() (int64, error)

	Clear()error
}
