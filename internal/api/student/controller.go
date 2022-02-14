package student

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/mydb"
	"encoding/json"
	"log"
)

func Consume(msgByte []byte) {
	//解析message
	var msg global.BookCourseRequest
	err := json.Unmarshal(msgByte, &msg)
	if err != nil {
		log.Fatalln(err)
	}
	//扣减课程余量
	global.MysqlClient.AutoMigrate(&mydb.Course{}) //迁移表到Course
	var course mydb.Course
	global.MysqlClient.Model(&course).Update("CourseId", course.Cap-1)
	//插入课表
	global.MysqlClient.AutoMigrate(&mydb.SelectCourse{}) //迁移表到SelectCourse
	global.MysqlClient.Create(&mydb.SelectCourse{StudentId: msg.StudentID, CourseId: msg.CourseID})
}
