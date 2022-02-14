package server

import (
	"Course-Selection-Scheduling/internal/api/course"
	"Course-Selection-Scheduling/internal/api/login"
	"Course-Selection-Scheduling/internal/api/member"
	"Course-Selection-Scheduling/internal/api/student"
	"github.com/gin-gonic/gin"
)

func registerRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", member.CreateMember)
	g.GET("/member", login.GetMember)
	g.GET("/member/list", member.ListMember)
	g.POST("/member/update", member.UpdateMember)
	g.POST("/member/delete", member.DeleteMember)

	// 登录
	g.POST("/auth/login", login.Login)
	g.POST("/auth/logout", login.Logout)
	g.GET("/auth/whoami", login.Whoami)

	// 排课
	g.POST("/course/create")
	g.GET("/course/get")

	g.POST("/teacher/bind_course")
	g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course")
	g.POST("/course/schedule", course.ScheduleCourse)

	// 抢课
	g.POST("/student/book_course", student.BookCourse)
	g.GET("/student/course")
}
