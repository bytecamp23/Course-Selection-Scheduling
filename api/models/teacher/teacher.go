package teacher

import (
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/types"
	"gorm.io/gorm"
)

// 老师绑定课程
// Method： Post
// 注：这里的 teacherID 不需要做已落库校验
// 一个老师可以绑定多个课程 , 不过，一个课程只能绑定在一个老师下面

type BindCourseRequest struct {
	CourseID  string `binding:"required,IsDigitValidator"`
	TeacherID string `binding:"required,IsDigitValidator"`
}

type BindCourseResponse struct {
	Code types.ErrNo
}

// 老师解绑课程
// Method： Post
type UnbindCourseRequest struct {
	CourseID  string `binding:"required,IsDigitValidator"`
	TeacherID string `binding:"required,IsDigitValidator"`
}

type UnbindCourseResponse struct {
	Code types.ErrNo
}

// 获取老师下所有课程
// Method：Get
type GetTeacherCourseRequest struct {
	TeacherID string `binding:"required,IsDigitValidator"`
}

type GetTeacherCourseResponse struct {
	Code types.ErrNo
	Data struct {
		CourseList []*types.TCourse
	}
}

// -----------------------------------
//检验绑定课程合法性
func (bindCourseInfo BindCourseRequest) CheckBind() (errno types.ErrNo) {
	var course mydb.Course
	err := mydb.MysqlClient.
		Model(&course).
		Where("course_id = ?", bindCourseInfo.CourseID).
		First(&course)
	if err.Error == gorm.ErrRecordNotFound {
		return types.CourseNotExisted
	}
	if course.TeacherId != nil {
		return types.CourseHasBound
	}
	return types.OK
}

//检验解绑课程合法性
func (unbindCourseInfo UnbindCourseRequest) CheckUnBind() (errno types.ErrNo) {
	unbindCourse := mydb.BindCourse{
		TeacherId: unbindCourseInfo.TeacherID,
		CourseId:  unbindCourseInfo.CourseID,
	}
	//课程不存在
	var course mydb.Course
	err := mydb.MysqlClient.
		Model(&course).
		Where("course_id = ?", unbindCourseInfo.CourseID).
		First(&course)
	if err.Error == gorm.ErrRecordNotFound {
		return types.CourseNotExisted
	}
	//课程未绑定
	err = mydb.MysqlClient.Model(&mydb.BindCourse{}).
		Where("course_id = ? AND teacher_id = ?", unbindCourse.CourseId, unbindCourse.TeacherId).
		First(&unbindCourse)
	if err.Error != nil {
		return types.CourseNotBind
	}
	return types.OK
}

//绑定课程
func (bindCourseInfo BindCourseRequest) Bind() (errno types.ErrNo) {
	db := mydb.MysqlClient
	bindCourse := mydb.BindCourse{
		TeacherId: bindCourseInfo.TeacherID,
		CourseId:  bindCourseInfo.CourseID,
	}
	//恢复软删除
	db.Model(bindCourse).Unscoped().Update("deleted_at", nil)
	//绑定课程
	mydb.MysqlClient.Model(&mydb.BindCourse{}).Create(&bindCourse)
	//添加课程绑定的教师号
	var course mydb.Course
	db.Where("course_id = ?", bindCourseInfo.CourseID).First(&course)
	course.TeacherId = &bindCourseInfo.TeacherID
	if err := db.Save(&course); err.Error != nil {
		return types.UnknownError
	}
	myredis.PutToRedis(types.TeacherIDPre+bindCourseInfo.CourseID, course.TeacherId, -1)
	return types.OK
}

//解绑课程
func (unbindCourseInfo UnbindCourseRequest) UnBind() (errno types.ErrNo) {
	db := mydb.MysqlClient
	unbindCourse := mydb.BindCourse{
		TeacherId: unbindCourseInfo.TeacherID,
		CourseId:  unbindCourseInfo.CourseID,
	}
	err := db.Model(&mydb.BindCourse{}).
		Where("course_id = ? AND teacher_id = ?", unbindCourse.CourseId, unbindCourse.TeacherId).
		Delete(&unbindCourse)
	if err.Error != nil {
		return types.UnknownError
	}

	//删除课程绑定的教师号
	var course mydb.Course
	db.Where("course_id = ?", unbindCourseInfo.CourseID).First(&course)
	course.TeacherId = nil
	if err := db.Save(&course); err.Error != nil {
		return types.UnknownError
	}
	myredis.DeleteFromRedis(types.TeacherIDPre + unbindCourseInfo.CourseID)
	return types.OK
}

//获得教师绑定的课程
func (teacherCourseInfo GetTeacherCourseRequest) GetTeacherCourses() (courses []*types.TCourse, errno types.ErrNo) {
	var courseIDs []string
	db := mydb.MysqlClient
	db.Model(&mydb.BindCourse{}).
		Where("teacher_id = ?", teacherCourseInfo.TeacherID).
		Select("course_id").
		Find(&courseIDs)
	db.Model(&mydb.Course{}).
		Where("course_id IN ?", courseIDs).
		Find(&courses)
	return courses, types.OK
}
