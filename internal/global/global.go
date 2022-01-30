package global

import (
	"github.com/garyburd/redigo/redis"
	"gorm.io/gorm"
)

var (
	MysqlClient *gorm.DB
	RedisClient *redis.Pool
)
