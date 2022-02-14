package main

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/internal/server"
	"Course-Selection-Scheduling/pkg/config"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/pkg/rabbitmq"
	"Course-Selection-Scheduling/utils"
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
	utils.SetLogPath()
	global.MysqlClient = mydb.NewMysqlConn(&config.MysqlCfg)
	global.RedisClient = myredis.NewRedisClient(&config.RedisCfg)

	rmq := rabbitmq.NewRabbitMQSimple("bookcourse") //打开rmq消费者
	rmq.ConsumeSimple()
	server.Run()
}
