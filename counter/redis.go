package counter

import (
	"fmt"
	"github.com/braveghost/meteor/random"
	"github.com/go-redis/redis"
	"math"
	"time"
)

var ttl time.Duration = -1

const (
	defaultCounterRedisKey = "braveghost.rogue.redis.counter"
)

func SetRedisCounterTtl(t int64) {
	ttl = time.Duration(t) * time.Second
}

type Redis struct {
	score float64
	name  string
}

func (c *Redis) getKey() string {
	return defaultCounterRedisKey + "." + c.name
}

func NewRedisCounter(t int64, name string) ICounter {
	return &Redis{float64(t), name}
}

func (c *Redis) Add() error {
	err := c.add(redis.Z{Score: c.score, Member: random.SnowFlake().String()})
	return err
}
func (c *Redis) add(val ...redis.Z) error {
	key := c.getKey()
	res := RedisDao().ZAdd(key, val...)
	if ttl > 0 {
		RedisDao().Expire(key, ttl)
	}
	return res.Err()
}

func (c *Redis) del(val ...interface{}) error {
	res := RedisDao().ZRem(c.getKey(), val...)
	return res.Err()
}

func (c *Redis) rangeVal() ([]string, error) {
	x := RedisDao().ZRangeByScore(c.getKey(),
		redis.ZRangeBy{
			Min: fmt.Sprint(c.score),
			Max: fmt.Sprint(c.score),
		})

	if x.Err() != nil {
		return nil, x.Err()
	}
	return x.Val(), nil
}
func (c *Redis) Set(n int64) error {
	val, err := c.rangeVal()

	if err != nil {
		return err
	}

	vLen := int64(len(val))
	m := n - vLen
	if m > 0 {
		var val []redis.Z
		var i int64

		for i = 0; i < m; i++ {
			val = append(val, redis.Z{Score: c.score, Member: random.SnowFlake().String()})
		}
		return c.add(val...)
	}
	if m < 0 {
		m = int64(math.Abs(float64(m)))
		return c.del(val[0:m])

	}
	return nil

}

func (c *Redis) Reset() error {
	res := RedisDao().ZRemRangeByScore(c.getKey(), fmt.Sprint(c.score), fmt.Sprint(c.score))
	return res.Err()
}

func (c *Redis) Compare(n int64) (bool, error) {
	count, err := c.Get()
	if err != nil {
		return false, err
	}
	return count == n, nil
}

func (c *Redis) Get() (int64, error) {
	res := RedisDao().ZCount(c.getKey(), fmt.Sprint(c.score), fmt.Sprint(c.score))
	if res.Err() != nil {
		return 0, res.Err()
	}
	return res.Val(), nil
}

func (c *Redis) Minus() error {
	val, err := c.rangeVal()

	if err != nil {
		return err
	}
	if len(val) == 0 {
		return nil
	}
	res := RedisDao().ZRem(c.getKey(), val[1])
	return res.Err()
}

func (c *Redis) Clear() error {
	res := RedisDao().ZRemRangeByScore(c.getKey(), "-inf", fmt.Sprint(c.score))
	return res.Err()

}
