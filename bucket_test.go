package rogue

import (
	"fmt"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	c := &Counter{0}
	fmt.Println(c)

	c.Add()
	fmt.Println(c)
	c.Minus()
	fmt.Println(c)
	c.Set(10)
	fmt.Println(c)

	fmt.Println(c.Compare(11))
	xx := c.Get()
	fmt.Println(c, xx)
	c.Reset()
	fmt.Println(c)
}
func TestNewBucketCounter(t *testing.T) {

	x := NewBucketCounter(5, 5)

	x.Increment()
	x.Increment()
	time.Sleep(time.Second * 3)
	x.Increment()
	x.Increment()
	time.Sleep(time.Second * 2)
	x.Increment()
	fmt.Println(x.Overflow())
	fmt.Println(x.Avg())
	fmt.Println(x.Max())
	fmt.Println(x.Min())
	fmt.Println(x)
}
