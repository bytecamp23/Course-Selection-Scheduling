package myredis

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/config"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

func NewRedisClient(cfg *config.Redis) *redis.Pool {
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

func ZAddToRedis(key string, score int64, member interface{}) error {
	rc := global.RedisClient.Get()
	defer rc.Close()
	_, err := rc.Do("ZADD", key, score, member)
	return err
}

func ZGetAllFromRedis(key string) (interface{}, error) {
	rc := global.RedisClient.Get()
	defer rc.Close()
	return rc.Do("ZRANGE", key, 0, -1)
}

func SAddToRedisSet(key string, member interface{}) error {

	rc := global.RedisClient.Get()
	defer rc.Close()
	_, err := rc.Do("SADD", key, member)
	return err
}

func SIsNumberOfRedisSet(key string, member interface{}) (bool, error) {
	rc := global.RedisClient.Get()
	defer rc.Close()
	value, err := redis.Bool(rc.Do("SISMEMBER", key, member))
	return value, err
}

func GetFromRedis(key string) (interface{}, error) {

	rc := global.RedisClient.Get()
	defer rc.Close()
	value, err := rc.Do("GET", key)
	return value, err
}

func PutToRedis(key string, value interface{}, timeout int) error {
	rc := global.RedisClient.Get()
	defer rc.Close()
	_, err := rc.Do("SET", key, value, "EX", timeout)
	return err
}

func DeleteFromRedis(key string) error {
	rc := global.RedisClient.Get()
	defer rc.Close()
	_, err := rc.Do("DEL", key)
	return err
}
