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
	SC := fmt.Sprintf("%s_%s", requestData.StudentID, requestData.CourseID)
	freqentSC := "frequent_" + SC
	successSC := "success_" + SC

	//限制抢课频率
	value, err := myredis.GetFromRedis(freqentSC)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if value != nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.RepeatRequest},
		)
		return
	} else {
		myredis.PutToRedis(freqentSC, "true", 5) //5秒内只能抢一次
	}

	//限制重复抢课
	value, err = myredis.GetFromRedis(successSC)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if value != nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.RepeatRequest},
		)
		return
	} else {
		myredis.PutToRedis(successSC, "true", -1)
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
