package course

//从数据库获得绑定课程信息 仅在服务开启后执行一次，之后通过缓存增量维护
/*func LoadBindCourses() (ret global.ScheduleCourseRequest) {
	var bindCourses []mydb.BindCourse
	global.MysqlClient.Find(&bindCourses)
	for _, bindCourse := range bindCourses {
		ret[bindCourse.CourseId]
	}

}*/
