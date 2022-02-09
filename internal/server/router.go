package server

import (
	"Course-Selection-Scheduling/internal/api/course"
	"github.com/gin-gonic/gin"
)

func registerRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create")
	g.GET("/member")
	g.GET("/member/list")
	g.POST("/member/update")
	g.POST("/member/delete")

	// 登录

	g.POST("/auth/login")
	g.POST("/auth/logout")
	g.GET("/auth/whoami")

	// 排课
	g.POST("/course/create")
	g.GET("/course/get")

	g.POST("/teacher/bind_course")
	g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course")
	g.POST("/course/schedule", course.Schedule)

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")
}
