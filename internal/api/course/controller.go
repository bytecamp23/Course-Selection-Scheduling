package course

import (
	"Course-Selection-Scheduling/internal/global"
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
