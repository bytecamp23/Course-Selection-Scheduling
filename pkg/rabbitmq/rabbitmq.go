package rabbitmq

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/mydb"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
)

const MQURL = "amqp://guest:guest@localhost:5672/"

type RabbitMQ struct {
	coon      *amqp.Connection
	channel   *amqp.Channel
	QueueName string // 队列名称
	Exchange  string // 交换机
	key       string
	Mqurl     string // 连接信息
}

func NewRabbitMQ(quequName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: quequName,
		Exchange:  exchange,
		key:       key,
		Mqurl:     MQURL,
	}
	return rabbitmq
}

func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.coon.Close()
}

func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// NewRabbitMQSimple 创建简单模式下RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	rabbitmq := NewRabbitMQ(queueName, "", "") //use default exchange, key->nil

	var err error
	rabbitmq.coon, err = amqp.Dial(rabbitmq.Mqurl) //连接RMQ服务器
	rabbitmq.failOnErr(err, "创建连接错误")
	rabbitmq.channel, err = rabbitmq.coon.Channel() //创建通道,大多数API通过该通道操作
	rabbitmq.failOnErr(err, "获取channel失败")

	return rabbitmq
}

func (r *RabbitMQ) PublishSimple(message []byte) {
	// 申请队列，固定用法，如果队列不存在会自动创建，如果存在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名称
		false,       // 是否持久化
		false,       // 是否自动删除
		false,       // 是否具有排他性, 如果为true则仅自己可见
		false,       // 是否阻塞
		nil,         // 额外属性
	)
	if err != nil {
		fmt.Println(err)
	}

	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		false, // 如果为true: 找不到符合条件的队列会返回给消费者
		false, // 如果为true: 当exchange发送消息到队列后，发现队列没有绑定消费者，则把消息返还给发送者
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
}

func (r *RabbitMQ) ConsumeSimple() {
	// 申请队列，固定用法，如果队列不存在会自动创建，如果存在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName, // 队列名称
		false,       // 是否持久化
		false,       // 是否自动删除
		false,       // 是否具有排他性, 如果为true则仅自己可见
		false,       // 是否阻塞
		nil,         // 额外属性
	)
	if err != nil {
		fmt.Println(err)
	}

	msgs, err := r.channel.Consume(
		q.Name,
		"",   // 用来区分多个消费者
		true, // 是否自动应答，是否通知RabbitMQ删掉消息
		false,
		false, // 如果为true: 不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool) //一个布尔类型的chan
	// 启用协程处理消息
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			Consume(d.Body)
		}
	}()
	log.Printf("[*] Waiting for messages, To exit press CTRL+C")
	<-forever
}
func Consume(msgByte []byte) {
	//解析message
	var msg global.BookCourseRequest
	err := json.Unmarshal(msgByte, &msg)
	if err != nil {
		log.Fatalln(err)
	}
	//扣减课程余量
	global.MysqlClient.AutoMigrate(&mydb.Course{}) //迁移表到Course
	global.MysqlClient.Model(&mydb.Course{}).
		Where("course_id = ?", msg.CourseID).
		Update("cap", gorm.Expr("cap- ?", 1))
	//插入课表
	global.MysqlClient.AutoMigrate(&mydb.SelectCourse{}) //迁移表到SelectCourse
	global.MysqlClient.Create(&mydb.SelectCourse{StudentId: msg.StudentID, CourseId: msg.CourseID})
}
