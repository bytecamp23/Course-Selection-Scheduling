package student

import (
	"Course-Selection-Scheduling/api/models/student"
	"Course-Selection-Scheduling/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func BookCourse(c *gin.Context) {
	var (
		requestData student.BookCourseRequest
		respondData student.BookCourseResponse
	)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
		return
	}
	respondData.Code = requestData.CheckValid()
	if respondData.Code != types.OK {
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
		return
	}

	//生成userCourse ,防止同一个用户多次抢同一门课
	success := fmt.Sprintf("success_%s_%s", requestData.StudentID, requestData.CourseID)
	frequency := fmt.Sprintf("frequency_%s", requestData.StudentID)

	respondData.Code = requestData.CheckRestriction(success, frequency)
	if respondData.Code != types.OK {
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
		return
	}

	respondData.Code = requestData.LockCourse(success)
	c.JSON(200, respondData)
	log.Println(requestData)
	log.Println(respondData)
}

// 根据学生信息查询课程列表
func QueryCourse(c *gin.Context) {
	var (
		requestData student.GetStudentCourseRequest
		respondData student.GetStudentCourseResponse
	)
	if err := c.ShouldBindQuery(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
		return
	}

	//学生不存在
	respondData.Code = requestData.CheckStudent()
	if respondData.Code != types.OK {
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
		return
	}
	//限制频率
	frequency := fmt.Sprintf("frequency_%s", requestData.StudentID)
	respondData.Code = requestData.CheckRestriction(frequency)
	if respondData.Code != types.OK {
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
		return
	}

	respondData.Data.CourseList, respondData.Code = requestData.GetCourses()
	c.JSON(200, &respondData)
	log.Println(requestData)
	log.Println(respondData)
}
