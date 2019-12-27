package counter

import (
	"fmt"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	x := Redis{float64(time.Now().Unix()),"test"}

	fmt.Println(x.Get())
	for i := 0; i <= 5; i++ {
		_ = x.Add()
	}
	fmt.Println(x.Get())
	_ = x.Set(8)

	fmt.Println(x.Get())

	_ = x.Set(3)

	fmt.Println(x.Get())
	fmt.Println(x.Compare(4))
	fmt.Println(x.Compare(3))
	_ = x.Minus()
	fmt.Println(x.Compare(2))
	x.Clear()
	fmt.Println(x.Get())
	_= x.Reset()
	fmt.Println(x.Get())
}
