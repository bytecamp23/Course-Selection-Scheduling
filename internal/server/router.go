package server

import (
	"Course-Selection-Scheduling/internal/api/course"
	"Course-Selection-Scheduling/internal/api/login"
	"Course-Selection-Scheduling/internal/api/member"
	"Course-Selection-Scheduling/internal/api/student"
	"Course-Selection-Scheduling/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func registerRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("UserNameValidator", utils.UserNameValidator)
		_ = v.RegisterValidation("PasswordValidator", utils.PasswordValidator)
		_ = v.RegisterValidation("UserTypeValidator", utils.UserTypeValidator)
		_ = v.RegisterValidation("UserIDValidator", utils.UserIDValidator)
	}

	// 成员管理
	memberRouter := g.Group("member")
	{
		memberRouter.POST("/create", member.CreateMember)
		memberRouter.GET("/", login.GetMember)
		memberRouter.GET("/list", member.ListMember)
		memberRouter.POST("/update", member.UpdateMember)
		memberRouter.POST("/delete", member.DeleteMember)
	}

	// 登录
	authRouter := g.Group("auth")
	{
		authRouter.POST("/login", login.Login)
		authRouter.POST("/logout", login.Logout)
		authRouter.GET("/whoami", login.Whoami)
	}

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
