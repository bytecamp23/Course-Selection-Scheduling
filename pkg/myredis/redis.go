package myredis

import (
	"Course-Selection-Scheduling/utils"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var RedisClient *redis.Pool

func NewRedisClient(cfg *utils.Redis) *redis.Pool {
	client := &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: time.Second * cfg.Timeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(cfg.ConType, cfg.Host+":"+cfg.Port, redis.DialDatabase(cfg.DbNum))
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			/*if _, err := c.Do("AUTH", redisConf["auth"].(string)); err != nil {
				_ = c.Close()
				fmt.Println(err.Error())
				return nil, err
			}*/
			return c, nil
		},
	}
	return client
}

func GetFromRedis(key string) (interface{}, error) {

	rc := RedisClient.Get()
	defer rc.Close()
	value, err := rc.Do("GET", key)
	return value, err
}

func PutToRedis(key string, value interface{}, timeout int) error {
	rc := RedisClient.Get()
	defer rc.Close()
	var err error
	if timeout == -1 {
		_, err = rc.Do("SET", key, value)
	} else {
		_, err = rc.Do("SET", key, value, "EX", timeout)
	}
	return err
}

func DeleteFromRedis(key string) error {
	rc := RedisClient.Get()
	defer rc.Close()
	_, err := rc.Do("DEL", key)
	return err
}

func DecrForRedis(key string) (interface{}, error) {
	rc := RedisClient.Get()
	defer rc.Close()
	value, err := rc.Do("DECR", key)
	return value, err
}

func IncrForRedis(key string) (interface{}, error) {
	rc := RedisClient.Get()
	defer rc.Close()
	value, err := rc.Do("INCR", key)
	return value, err
}

func SAddToRedisSet(key string, member interface{}) error {
	rc := RedisClient.Get()
	defer rc.Close()
	_, err := rc.Do("SADD", key, member)
	return err
}

func Exsits(key string) (bool, error) {
	rc := RedisClient.Get()
	defer rc.Close()
	return redis.Bool(rc.Do("EXISTS", key))
}

func SGetAllFromRedis(key string) ([]string, error) {
	rc := RedisClient.Get()
	defer rc.Close()
	return redis.Strings(rc.Do("SMEMBERS", key))
}

func Flushdb() error {
	rc := RedisClient.Get()
	defer rc.Close()
	_, err := rc.Do("flushdb")
	return err
}
