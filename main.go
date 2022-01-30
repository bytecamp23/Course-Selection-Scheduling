package main

import (
	"bytecamp/internal/global"
	"bytecamp/internal/server"
	"bytecamp/pkg/config"
	"bytecamp/pkg/mydb"
	"bytecamp/pkg/myredis"
	"bytecamp/utils"
	"flag"
)

func main() {
	defer func() {
		_ = global.RedisClient.Close()   //关闭redis
		db, _ := global.MysqlClient.DB() //获取已有sql连接
		_ = db.Close()                   //关闭sql连接
	}()

	var env = flag.String("env", "dev", "配置环境")
	flag.Parse() //获取命令行参数 根据参数选择配置文件 默认dev
	utils.LoadCfg(*env)
	global.MysqlClient = mydb.NewMysqlConn(&config.MysqlCfg)
	global.RedisClient = myredis.NewRedisClient(&config.RedisCfg)
	server.Run()
}
