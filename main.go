package main

import (
	"bytecamp/global"
	"log"
)

func sqlExample() {
	student := global.User{Nickname: "student", Username: "student123", Password: "123456", UserType: 2}
	result := global.MysqlClient.Create(&student)
	if result.Error != nil {
		log.Println("insert error")
	}

	teacher := global.User{Nickname: "teacher", Username: "teacher123", Password: "123456", UserType: 3}
	result = global.MysqlClient.Create(&teacher)
	if result.Error != nil {
		log.Println("insert error")
	}

	course := global.Course{Name: "data structure", Cap: 100}
	result = global.MysqlClient.Create(&course)
	if result.Error != nil {
		log.Println("insert error")
	}

	bindCourse := global.BindCourse{TeacherId: teacher.UserId, CourseId: course.CourseId}
	result = global.MysqlClient.Create(&bindCourse)
	if result.Error != nil {
		log.Println("insert error")
	}

	selectCourse := global.SelectCourse{StudentId: student.UserId, CourseId: course.CourseId}
	result = global.MysqlClient.Create(&selectCourse)
	if result.Error != nil {
		log.Println("insert error")
	}

	result = global.MysqlClient.Where("username = ?", "teacher123").Delete(&global.User{})
	if result.Error != nil {
		log.Println("delete error")
	}
}
func main() {
	defer func() {
		db, _ := global.MysqlClient.DB() //获取已有sql连接
		_ = db.Close()                   //关闭sql连接
	}()
	global.NewMysqlConn("root", "1234567890", "bytecamp")
	sqlExample()
}
