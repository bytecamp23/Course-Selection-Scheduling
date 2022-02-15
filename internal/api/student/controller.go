package student

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/config"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/pkg/rabbitmq"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.UnknownError},
		)
		return
	}
	if value != nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.RepeatRequest},
		)
		return
	} else {
		myredis.PutToRedis(freqentSC, "true", 3) //3秒内只能抢一次
	}

	//限制重复抢课
	value, err = myredis.GetFromRedis(successSC)
	if err != nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.UnknownError},
		)
		return
	}
	if value != nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.RepeatRequest},
		)
		return
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
		myredis.IncrForRedis(requestData.CourseID) //加回来
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.CourseNotAvailable},
		)
		return
	}
	//减库存成功后 标记为抢课成功
	myredis.PutToRedis(successSC, "true", -1)

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

func QueryCourse(c *gin.Context){
	var json global.GetStudentCourseRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		// TODO: ParamInvalid
		getStudentCourseResponse := global.GetStudentCourseResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, getStudentCourseResponse)
		return
	}

	var coursesInfo []mydb.Course
	var res global.GetStudentCourseResponse
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	db.Model(&mydb.Course{}).Where("student_id = ?", json.StudentID).Find(&coursesInfo)

	res.Data.CourseList = make([]global.TCourse, len(coursesInfo))
	for i := 0; i < len(coursesInfo); i++ {
		res.Data.CourseList[i].Name = coursesInfo[i].Name
		res.Data.CourseList[i].CourseID = coursesInfo[i].CourseId
		res.Data.CourseList[i].TeacherID = coursesInfo[i].TeacherId
	}
}