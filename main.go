package main

import (
	"Course-Selection-Scheduling/api/models/student"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/pkg/rabbitmq"
	"Course-Selection-Scheduling/pkg/server"
	"Course-Selection-Scheduling/types"
	"Course-Selection-Scheduling/utils"
	"flag"
)

func main() {
	defer func() {
		_ = myredis.RedisClient.Close() //关闭redis
		rabbitmq.RMQClient.Destory()
		//mysql由gorm管理
	}()
	var env = flag.String("env", "dev", "配置环境")
	flag.Parse() //获取命令行参数 根据参数选择配置文件 默认dev
	utils.LoadCfg(*env)
	mydb.MysqlClient = mydb.NewMysqlConn(&utils.MysqlCfg)
	mydb.CreateTables()
	myredis.RedisClient = myredis.NewRedisClient(&utils.RedisCfg)
	rabbitmq.RMQClient = rabbitmq.NewRabbitMQSimple(types.RMQName) //打开rmq消费者
	go rabbitmq.RMQClient.ConsumeSimple(student.Consume)
	server.Run()
}
