package student

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func BookCourse(c *gin.Context) {
	requestData := global.BookCourseRequest{}
	if err := c.BindJSON(&requestData); err != nil { //若绑定出错
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.ParamInvalid},
		)
		return
	}
	//生成userCourse ,防止同一个用户多次抢同一门课
	stuCourse := fmt.Sprintf("%s_%s", requestData.StudentID, requestData.CourseID)
	value, err := myredis.GetFromRedis(stuCourse)
	log.Println(value)
	if err != nil {
		log.Println(err.Error())
		return
	} else {
		if value != nil {
			c.JSON(
				http.StatusOK,
				global.BookCourseResponse{Code: global.RepeatRequest},
			)
			return
		} else {
			myredis.PutToRedis(stuCourse, "true", 5) //5秒内只能抢一次
		}
	}

	//查询课程余量并减库存 , 数据库操作送入消息队列中
	value, err = myredis.DecrForRedis(requestData.CourseID)
	if err != nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.UnknownError},
		)
		return
	}
	if value.(int64) < 0 {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.CourseNotAvailable},
		)
		return
	}

	//放到消息队列中,进行数据库操作
	msgByte, err := json.Marshal(requestData)
	if err != nil {
		log.Fatalln(err)
	}
	rmq := rabbitmq.NewRabbitMQSimple("bookcourse")
	rmq.PublishSimple(msgByte)

	//回传OK
	c.JSON(
		http.StatusOK,
		global.BookCourseResponse{Code: global.OK},
	)
}
