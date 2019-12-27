package counter

import (
	"github.com/go-redis/redis"
)

const Nil = "redis: nil"

var (
	defaultRedisDao *redisDao
	redisPotions    = &redis.Options{
		Addr:       "localhost:6379",
		MaxRetries: 3,
		Password:   "", // no password set
		DB:         0,  // use default DB
		PoolSize:   100,
	}
)

type redisDao struct {
	*redis.Client
}

func SetRedisServerAddr(addr string) {
	redisPotions.Addr = addr
}
func SetRedisMaxRetries(max int) {
	redisPotions.MaxRetries = max
}

func SetRedisPoolSize(size int) {
	redisPotions.PoolSize = size
}
func SetRedisPassword(pwd string) {
	redisPotions.Password = pwd
}

func RedisDao() *redisDao {
	if defaultRedisDao == nil {
		defaultRedisDao = &redisDao{redis.NewClient(redisPotions)}
	}
	return defaultRedisDao

}


func (rd *redisDao) SetValue(key string, v interface{}) error {
	res := rd.Client.Set(key, v, 0)
	return res.Err()
}

func (rd *redisDao) GetStr(key string) (string, error) {
	res := rd.Client.Get(key)
	if res.Err() != nil {
		if res.Err().Error() == Nil {
			return "", nil
		}
		return "", res.Err()
	}
	return res.Val(), nil
}
