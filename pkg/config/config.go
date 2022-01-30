package config

import "time"

var (
	MysqlCfg   Mysql
	LogCfg     Log
	ServerCfg  Server
	RedisCfg   Redis
	SessionCfg Session
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
