package teacher

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/mydb"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BindCourse(c *gin.Context) {
	var bindCourseRequest global.BindCourseRequest
	if err := c.ShouldBindJSON(&bindCourseRequest); err != nil {
		// TODO: ParamInvalid
		bindCourseResponse := global.BindCourseResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, bindCourseResponse)
		return
	}
	bindCourse := mydb.BindCourse{
		TeacherId: bindCourseRequest.TeacherID,
		CourseId:  bindCourseRequest.CourseID,
	}
	// TODO: CourseNotExisted
	if err := global.MysqlClient.
		Model(&mydb.Course{}).
		Where("course_id = ?", bindCourseRequest.CourseID).
		First(&mydb.Course{}); err.Error == gorm.ErrRecordNotFound {
		bindCourseResponse := global.BindCourseResponse{
			Code: global.CourseNotExisted,
		}
		c.JSON(200, bindCourseResponse)
	}
	// TODO: CourseHasBound
	if err := global.MysqlClient.Model(&mydb.BindCourse{}).Create(&bindCourse); err != nil {
		bindCourseResponse := global.BindCourseResponse{
			Code: global.CourseHasBound,
		}
		c.JSON(200, bindCourseResponse)
		return
	}
	bindCourseResponse := global.BindCourseResponse{
		Code: global.OK,
	}
	c.JSON(200, bindCourseResponse)
	return
}

func UnBindCourse(c *gin.Context) {
	var unbindCourseRequest global.UnbindCourseRequest
	if err := c.ShouldBindJSON(&unbindCourseRequest); err != nil {
		// TODO: ParamInvalid
		unbindCourseResponse := global.UnbindCourseResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, unbindCourseResponse)
		return
	}
	unbindCourse := mydb.BindCourse{
		TeacherId: unbindCourseRequest.TeacherID,
		CourseId:  unbindCourseRequest.CourseID,
	}
	// TODO: CourseNotExisted
	if err := global.MysqlClient.
		Model(&mydb.Course{}).
		Where("course_id = ?", unbindCourseRequest.CourseID).
		First(&mydb.Course{}); err.Error == gorm.ErrRecordNotFound {
		bindCourseResponse := global.BindCourseResponse{
			Code: global.CourseNotExisted,
		}
		c.JSON(200, bindCourseResponse)
	}
	// TODO: CourseHasBound
	if err := global.MysqlClient.Model(&mydb.BindCourse{}).
		Where("course_id = ? AND teacher_id = ?", unbindCourse.CourseId, unbindCourse.TeacherId).
		Delete(&unbindCourse); err != nil {
		bindCourseResponse := global.BindCourseResponse{
			Code: global.CourseNotBind,
		}
		c.JSON(200, bindCourseResponse)
		return
	}
	bindCourseResponse := global.BindCourseResponse{
		Code: global.OK,
	}
	c.JSON(200, bindCourseResponse)
	return
}

func GetTeacherCourse(c *gin.Context) {
	var getTeacherCourseRequest global.GetTeacherCourseRequest
	if err := c.ShouldBindJSON(&getTeacherCourseRequest); err != nil {
		// TODO: ParamInvalid
		getTeacherCourseResponse := global.GetTeacherCourseResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, getTeacherCourseResponse)
		return
	}
	var courseIDs []string
	global.MysqlClient.Model(&mydb.BindCourse{}).Where("teacher_id = ?", getTeacherCourseRequest.TeacherID).Select("course_id").Find(&courseIDs)
	var courses []*global.TCourse
	global.MysqlClient.Model(&mydb.Course{}).Where("course_id IN ?", courseIDs).Find(&courses)
	getTeacherCourseResponse := global.GetTeacherCourseResponse{
		Code: global.OK,
		Data: struct{ CourseList []*global.TCourse }{CourseList: courses},
	}
	c.JSON(200, getTeacherCourseResponse)
	return
}
