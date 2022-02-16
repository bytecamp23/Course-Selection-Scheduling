package utils

import (
	"Course-Selection-Scheduling/pkg/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var (
	MysqlCfg   Mysql
	LogCfg     Log
	ServerCfg  Server
	RedisCfg   Redis
	SessionCfg Session
	RMQCfg     RabbitMQ
)

type Log struct {
	Path string `yaml:"path"`
}
type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"db_name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
type Redis struct {
	Host      string        `yaml:"host"`
	Port      string        `yaml:"port"`
	Auth      string        `yaml:"auth"`
	ConType   string        `yaml:"con_type"`
	DbNum     int           `yaml:"db_num"`
	MaxIdle   int           `yaml:"max_idle"`
	MaxActive int           `yaml:"max_active"`
	Timeout   time.Duration `yaml:"timeout"`
}
type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Mode string `yaml:"mode"`
}
type Session struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

type RabbitMQ struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// LoadCfg 加载配置文件
func LoadCfg(env string) {
	basePath := path.Join("./config", env)

	logPath := path.Join(basePath, "log.yml")
	LogCfg = NewLogConfig(logPath)
	SetLogPath()

	mysqlPath := path.Join(basePath, "mysql.yml")
	MysqlCfg = NewMysqlConfig(mysqlPath)

	serverPath := path.Join(basePath, "server.yml")
	ServerCfg = NewServerConfig(serverPath)

	redisPath := path.Join(basePath, "redis.yml")
	RedisCfg = NewRedisConfig(redisPath)

	sessionPath := path.Join(basePath, "session.yml")
	SessionCfg = NewSessionConfig(sessionPath)

	rabbitMQPath := path.Join(basePath, "rabbitMQ.yml")
	RMQCfg = NewRabbitMQConfig(rabbitMQPath)
}

func SetLogPath() {
	if LogCfg.Path == "stdout" {
		return
	}
	LogCfg.Path = path.Join(LogCfg.Path, time.Now().String())
	logFile, err := os.Create(LogCfg.Path)
	logger.FilePanic(LogCfg.Path, err)
	//设置logger日志
	log.SetOutput(logFile)
	//设置gin日志
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
}

func NewLogConfig(path string) Log {
	cfg, err := ioutil.ReadFile(path)
	logger.FilePanic(path, err)
	var log Log
	err = yaml.Unmarshal(cfg, &log)
	logger.ParsePanic(path, err)
	return log
}
func NewMysqlConfig(path string) Mysql {
	cfg, err := ioutil.ReadFile(path)
	logger.FilePanic(path, err)
	var mysql Mysql
	err = yaml.Unmarshal(cfg, &mysql)
	logger.ParsePanic(path, err)
	return mysql
}
func NewRedisConfig(path string) Redis {
	cfg, err := ioutil.ReadFile(path)
	logger.FilePanic(path, err)
	var redis Redis
	err = yaml.Unmarshal(cfg, &redis)
	logger.ParsePanic(path, err)
	return redis
}
func NewServerConfig(path string) Server {
	cfg, err := ioutil.ReadFile(path)
	logger.FilePanic(path, err)
	var server Server
	err = yaml.Unmarshal(cfg, &server)
	logger.ParsePanic(path, err)
	return server
}
func NewSessionConfig(path string) Session {
	cfg, err := ioutil.ReadFile(path)
	logger.FilePanic(path, err)
	var session Session
	err = yaml.Unmarshal(cfg, &session)
	logger.ParsePanic(path, err)
	return session
}

func NewRabbitMQConfig(path string) RabbitMQ {
	cfg, err := ioutil.ReadFile(path)
	logger.FilePanic(path, err)
	var rabbitmq RabbitMQ
	err = yaml.Unmarshal(cfg, &rabbitmq)
	logger.ParsePanic(path, err)
	return rabbitmq
}
