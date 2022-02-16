package server

import (
	"Course-Selection-Scheduling/api/controllers/auth"
	"Course-Selection-Scheduling/api/controllers/course"
	"Course-Selection-Scheduling/api/controllers/member"
	"Course-Selection-Scheduling/api/controllers/student"
	"Course-Selection-Scheduling/api/controllers/teacher"
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
		memberRouter.GET("/", member.GetMember)
		memberRouter.GET("/list", member.ListMember)
		memberRouter.POST("/update", member.UpdateMember)
		memberRouter.POST("/delete", member.DeleteMember)
	}

	// 登录
	authRouter := g.Group("auth")
	{
		authRouter.POST("/login", auth.Login)
		authRouter.POST("/logout", auth.Logout)
		authRouter.GET("/whoami", auth.Whoami)
	}

	// 排课
	courseRouter := g.Group("course")
	{
		courseRouter.POST("/create", course.CreateCourse)
		courseRouter.GET("/get", course.GetCourse)
		courseRouter.POST("/schedule", course.ScheduleCourse)
	}

	teacherRouter := g.Group("teacher")
	{
		teacherRouter.POST("/bind_course", teacher.BindCourse)
		teacherRouter.POST("/unbind_course", teacher.UnBindCourse)
		teacherRouter.GET("/get_course", teacher.GetTeacherCourse)
	}

	// 抢课
	g.POST("/student/book_course", student.BookCourse)
	g.GET("/student/course", student.QueryCourse)
}
