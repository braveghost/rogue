package rogue

import (
	"fmt"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	b := NewBucket(1, 5, ECounterTypeRedis, UnitMinute,"test")
	_ = b.Increment()
	_ = b.Increment()
	_ = b.Increment()
	//time.Sleep(time.Second * 2)
	_ = b.Increment()
	_ = b.Increment()
	time.Sleep(time.Second * 4)
	fmt.Println(b.Counters)
	for k, v := range b.Counters {
		fmt.Println("====", k)
		fmt.Println(v.Get())
	}
	fmt.Println("sum")
	fmt.Println(b.Sum())

	fmt.Println("avg")

	fmt.Println(b.Avg())

	fmt.Println("max")
	fmt.Println(b.Max())

	fmt.Println("min")
	fmt.Println(b.Min())
	fmt.Println(b.Overflow())
}
