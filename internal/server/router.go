package server

import (
	"Course-Selection-Scheduling/internal/api/course"
	"Course-Selection-Scheduling/internal/api/login"
	"Course-Selection-Scheduling/internal/api/member"
	"Course-Selection-Scheduling/internal/api/student"
	"Course-Selection-Scheduling/internal/api/teacher"
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
		_ = v.RegisterValidation("IsDigitValidator", utils.IsDigitValidator)
		_ = v.RegisterValidation("IsUpperOrLowerOrDigitValidator", utils.IsUpperOrLowerOrDigitValidator)
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
	courseRouter := g.Group("course")
	{
		courseRouter.POST("/create", course.CreateCourse)
		courseRouter.GET("/get", course.GetCourse)
		courseRouter.POST("/course/schedule", course.ScheduleCourse)
	}

	teacherRouter := g.Group("teacher")
	{
		teacherRouter.POST("/bind_course", teacher.BindCourse)
		teacherRouter.POST("/unbind_course", teacher.UnBindCourse)
		teacherRouter.GET("/get_course", teacher.GetTeacherCourse)
	}

	// 抢课
	g.POST("/student/book_course", student.BookCourse)
	g.GET("/student/course")
}
