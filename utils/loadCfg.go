package utils

import (
	"Course-Selection-Scheduling/pkg/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

//加载配置文件
func LoadCfg(env string) {
	basePath := path.Join("./config", env)

	logPath := path.Join(basePath, "log.yml")
	config.LogCfg = NewLogConfig(logPath)
	config.LogCfg.Path = path.Join(config.LogCfg.Path, time.Now().String())
	logFile, err := os.Create(config.LogCfg.Path)
	if err != nil {
		panic(fmt.Sprintf("create logs failed, path: %v. err: %v", config.LogCfg.Path, err.Error()))
	}
	//设置logger日志
	log.SetOutput(logFile)
	//设置gin日志
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)

	mysqlPath := path.Join(basePath, "mysql.yml")
	config.MysqlCfg = NewMysqlConfig(mysqlPath)

	serverPath := path.Join(basePath, "server.yml")
	config.ServerCfg = NewServerConfig(serverPath)

	redisPath := path.Join(basePath, "redis.yml")
	config.RedisCfg = NewRedisConfig(redisPath)

	sessionPath := path.Join(basePath, "session.yml")
	config.SessionCfg = NewSessionConfig(sessionPath)
}

func NewLogConfig(path string) config.Log {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("read config file failed, path: %v. err: %v", path, err.Error()))
	}
	var log config.Log
	err = yaml.Unmarshal(cfg, &log)
	if err != nil {
		panic(fmt.Sprintf("parse config file failed, path: %v. err: %v", path, err.Error()))
	}
	return log
}
func NewMysqlConfig(path string) config.Mysql {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("read config file failed, path: %v. err: %v", path, err.Error()))
	}
	var mysql config.Mysql
	err = yaml.Unmarshal(cfg, &mysql)
	if err != nil {
		panic(fmt.Sprintf("parse config file failed, path: %v. err: %v", path, err.Error()))
	}
	return mysql
}
func NewRedisConfig(path string) config.Redis {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("read config file failed, path: %v. err: %v", path, err.Error()))
	}
	var redis config.Redis
	err = yaml.Unmarshal(cfg, &redis)
	if err != nil {
		panic(fmt.Sprintf("parse config file failed, path: %v. err: %v", path, err.Error()))
	}
	return redis
}
func NewServerConfig(path string) config.Server {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("read config file failed, path: %v. err: %v", path, err.Error()))
	}
	var server config.Server
	err = yaml.Unmarshal(cfg, &server)
	if err != nil {
		panic(fmt.Sprintf("parse config file failed, path: %v. err: %v", path, err.Error()))
	}
	return server
}
func NewSessionConfig(path string) config.Session {
	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("read config file failed, path: %v. err: %v", path, err.Error()))
	}
	var session config.Session
	err = yaml.Unmarshal(cfg, &session)
	if err != nil {
		panic(fmt.Sprintf("parse config file failed, path: %v. err: %v", path, err.Error()))
	}
	return session
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
