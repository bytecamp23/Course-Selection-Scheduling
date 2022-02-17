package course

import (
	"Course-Selection-Scheduling/api/models/course"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func ScheduleCourse(c *gin.Context) {
	var (
		requestData = course.ScheduleCourseRequest{TeacherCourseRelationShip: map[string][]string{}}
		respondData course.ScheduleCourseResponse
	)
	if err := c.BindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
	}
	pointCnt := requestData.Discretize()
	if pointCnt*len(requestData.TeacherCourseRelationShip) > types.ChooseFactor {
		respondData.Data = requestData.Dinic()
	} else {
		respondData.Data = requestData.Hungarian()
	}
	c.JSON(http.StatusOK, respondData)
	log.Println(requestData)
	log.Println(respondData)
}

func CreateCourse(c *gin.Context) {
	var (
		requestData course.CreateCourseRequest
		respondData course.CreateCourseResponse
	)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
		return
	}
	respondData.Data.CourseID, respondData.Code = requestData.CreateCourse()
	c.JSON(200, respondData)
	log.Println(requestData)
	log.Println(respondData)
	return
}

func GetCourse(c *gin.Context) {
	var (
		requestData course.GetCourseRequest
		respondData course.GetCourseResponse
	)
	if err := c.ShouldBindQuery(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(requestData)
		log.Println(respondData)
		return
	}
	var course mydb.Course
	course, respondData.Code = requestData.GetCourseInfo()
	respondData.Data = struct {
		CourseID  string
		Name      string
		TeacherID string
	}{CourseID: course.CourseId, Name: course.Name}
	if course.TeacherId != nil {
		respondData.Data.TeacherID = *course.TeacherId
	}
	c.JSON(200, respondData)
	log.Println(requestData)
	log.Println(respondData)
	return
}
