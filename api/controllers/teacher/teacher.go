package teacher

import (
	"Course-Selection-Scheduling/api/models/teacher"
	"Course-Selection-Scheduling/types"
	"github.com/gin-gonic/gin"
)

func BindCourse(c *gin.Context) {
	var (
		requestData teacher.BindCourseRequest
		respondData teacher.BindCourseResponse
	)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		return
	}

	//检验绑定合法性
	respondData.Code = requestData.CheckBind()
	if respondData.Code != types.OK {
		c.JSON(200, respondData)
		return
	}
	respondData.Code = requestData.Bind()
	c.JSON(200, respondData)
	return
}
func UnBindCourse(c *gin.Context) {
	var (
		requestData teacher.UnbindCourseRequest
		respondData teacher.UnbindCourseResponse
	)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		return
	}
	respondData.Code = requestData.CheckUnBind()
	if respondData.Code != types.OK {
		c.JSON(200, respondData)
		return
	}
	respondData.Code = requestData.UnBind()
	c.JSON(200, respondData)
	return
}

func GetTeacherCourse(c *gin.Context) {
	var (
		requestData teacher.GetTeacherCourseRequest
		respondData teacher.GetTeacherCourseResponse
	)
	if err := c.ShouldBindQuery(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		return
	}
	respondData.Data.CourseList, respondData.Code = requestData.GetTeacherCourses()
	c.JSON(200, respondData)
	return
}
