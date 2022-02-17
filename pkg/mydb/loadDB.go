package mydb

import (
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/types"
	"fmt"
)

func LoadDB() {
	myredis.Flushdb()
	db := MysqlClient
	//课程余量 课程信息
	var courses []Course
	db.Find(&courses)
	for _, course := range courses {
		myredis.PutToRedis(types.CoursePre+course.CourseId, course.Cap, -1)
		myredis.PutToRedis(types.CourseNamePre+course.CourseId, course.Name, -1)
		if course.TeacherId != nil {
			myredis.PutToRedis(types.TeacherIDPre+course.CourseId, *course.TeacherId, -1)
		}
	}
	//学生身份
	var users []User
	db.Find(&users)
	for _, user := range users {
		if user.UserType == types.Student {
			myredis.PutToRedis(types.StudentPre+user.UserId, "true", -1)
		}
	}
	//选课关系
	var selectCourses []SelectCourse
	db.Find(&selectCourses)
	for _, selectCourse := range selectCourses {
		myredis.SAddToRedisSet(types.SelectPre+selectCourse.StudentId, selectCourse.CourseId)
		success := fmt.Sprintf("success_%s_%s", selectCourse.StudentId, selectCourse.CourseId)
		myredis.PutToRedis(success, 0, -1)
	}
}
