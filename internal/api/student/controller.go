package student

import (
	"Course-Selection-Scheduling/internal/global"
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
	//课程不存在
	value, err := myredis.GetFromRedis("course_" + requestData.CourseID)
	if value == nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.CourseNotExisted},
		)
		return
	}
	//学生不存在
	value, err = myredis.GetFromRedis("student_" + requestData.StudentID)
	if value == nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.StudentNotExisted},
		)
		return
	}

	//生成userCourse ,防止同一个用户多次抢同一门课
	SC := fmt.Sprintf("%s_%s", requestData.StudentID, requestData.CourseID)
	freqentSC := "frequent_" + SC
	successSC := "success_" + SC

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
			global.BookCourseResponse{Code: global.StudentHasCourse},
		)
		return
	}

	//限制抢课频率
	value, err = myredis.GetFromRedis(freqentSC)
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

	//查询课程余量并减库存 , 数据库操作送入消息队列中
	value, err = myredis.DecrForRedis("course_" + requestData.CourseID)
	if err != nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.UnknownError},
		)
		return
	}
	if value.(int64) < 0 {
		myredis.IncrForRedis("course_" + requestData.CourseID) //加回来
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

// 根据学生信息查询课程列表
func QueryCourse(c *gin.Context) {
	var json global.GetStudentCourseRequest
	if err := c.ShouldBindQuery(&json); err != nil {
		// TODO: ParamInvalid
		getStudentCourseResponse := global.GetStudentCourseResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, getStudentCourseResponse)
		return
	}
	//学生不存在
	value, _ := myredis.GetFromRedis("student_" + json.StudentID)
	if value == nil {
		c.JSON(
			http.StatusOK,
			global.BookCourseResponse{Code: global.StudentNotExisted},
		)
		return
	}
	db := global.MysqlClient
	var selectCourses []mydb.SelectCourse
	db.Model(&mydb.SelectCourse{}).Where("student_id = ?", json.StudentID).Find(&selectCourses)

	var res global.GetStudentCourseResponse
	res.Data.CourseList = make([]global.TCourse, len(selectCourses))
	for i, selectCourse := range selectCourses {
		var courseInfo mydb.Course
		db.Model(&mydb.Course{}).Where("course_id = ?", selectCourse.CourseId).Find(&courseInfo)
		fmt.Println(selectCourse.CourseId)
		res.Data.CourseList[i].CourseID = courseInfo.CourseId
		res.Data.CourseList[i].Name = courseInfo.Name
		if courseInfo.TeacherId != nil {
			res.Data.CourseList[i].TeacherID = *courseInfo.TeacherId
		}
	}
	if len(res.Data.CourseList) == 0 {
		res.Code = global.StudentHasNoCourse
	} else {
		res.Code = global.OK
	}
	c.JSON(200, &res)
}
