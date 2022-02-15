package course

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/mydb"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	request
	{"TeacherCourseRelationShip":{"1":["7","8"],"2":["6","9","10"],"3":["7","8"],"4":["7","8"],"5":["10"]}}
	respond
	{
		"Code": 0,
		"Data": {
			"2": "9",
			"3": "7",
			"4": "8",
			"5": "10"
		}
	}
*/
func ScheduleCourse(c *gin.Context) {
	requestData := global.ScheduleCourseRequest{TeacherCourseRelationShip: map[string][]string{}}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(
			http.StatusOK,
			global.ScheduleCourseResponse{Code: global.ParamInvalid},
		)
	}
	discretize(requestData.TeacherCourseRelationShip)
	var respondData map[string]string
	if pointCnt*len(requestData.TeacherCourseRelationShip) > 10000000 {
		respondData = dinic(requestData.TeacherCourseRelationShip)
	} else {
		respondData = hungarian(requestData.TeacherCourseRelationShip)
	}
	c.JSON(
		http.StatusOK,
		global.ScheduleCourseResponse{Code: global.OK, Data: respondData},
	)
}

func CreateCourse(c *gin.Context) {
	var createCourseRequest global.CreateCourseRequest
	if err := c.ShouldBindJSON(&createCourseRequest); err != nil {
		// TODO: ParamInvalid
		fmt.Println(err.Error())
		createCourseResponse := global.CreateCourseResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, createCourseResponse)
		return
	}

	course := mydb.Course{
		Name: createCourseRequest.Name,
		Cap:  createCourseRequest.Cap,
	}
	_ = global.MysqlClient.Create(&course)
	fmt.Println(course)
	createCourseResponse := global.CreateCourseResponse{
		Code: global.OK,
		Data: struct{ CourseID string }{CourseID: course.CourseId},
	}
	c.JSON(200, createCourseResponse)
	return
}

func GetCourse(c *gin.Context) {
	var getCourseRequest global.GetCourseRequest
	if err := c.ShouldBindQuery(&getCourseRequest); err != nil {
		// TODO: ParamInvalid
		fmt.Println(err.Error())
		getCourseResponse := global.GetCourseResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, getCourseResponse)
		return
	}
	var course mydb.Course
	course.TeacherId = new(string)
	err := global.MysqlClient.Model(&course).Where("course_id = ?", getCourseRequest.CourseID).First(&course)
	if err != nil {
		c.JSON(200, global.GetCourseResponse{
			Code: global.CourseNotExisted,
		})
		return
	}
	getCourseResponse := global.GetCourseResponse{
		Code: global.OK,
		Data: struct {
			CourseID  string
			Name      string
			TeacherID string
		}{CourseID: course.CourseId, Name: course.Name},
	}
	if course.TeacherId != nil {
		getCourseResponse.Data.TeacherID = *course.TeacherId
	}
	c.JSON(200, getCourseResponse)
	return
}
