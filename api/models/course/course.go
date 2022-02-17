package course

import (
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/types"
	"log"
)

// -------------------------------------
// 排课

// 创建课程
// Method: Post
type CreateCourseRequest struct {
	Name string `binding:"required"`
	Cap  int    `binding:"required"`
}

type CreateCourseResponse struct {
	Code types.ErrNo
	Data struct {
		CourseID string
	}
}

// 获取课程
// Method: Get
type GetCourseRequest struct {
	CourseID string `binding:"required,IsDigitValidator"`
}

type GetCourseResponse struct {
	Code types.ErrNo
	Data types.TCourse
}

// 排课求解器，使老师绑定课程的最优解， 老师有且只能绑定一个课程
// Method： Post
type ScheduleCourseRequest struct {
	TeacherCourseRelationShip map[string][]string // key 为 teacherID , val 为老师期望绑定的课程 courseID 数组
}

type ScheduleCourseResponse struct {
	Code types.ErrNo
	Data map[string]string // key 为 teacherID , val 为老师最终绑定的课程 courseID
}

// -------------------------------------
//创建课程
func (createCourseInfo CreateCourseRequest) CreateCourse() (CourseID string, errno types.ErrNo) {
	var course mydb.Course
	db := mydb.MysqlClient
	db.Where("name = ?", createCourseInfo.Name).First(&course)
	//检验课程是否已经存在
	if course.Name == createCourseInfo.Name {
		return "", types.UnknownError
	} else {
		course = mydb.Course{
			Name: createCourseInfo.Name,
			Cap:  createCourseInfo.Cap,
		}
		_ = mydb.MysqlClient.Create(&course)
		log.Println(types.CoursePre + course.CourseId)
		myredis.PutToRedis(types.CoursePre+course.CourseId, course.Cap, -1)
		myredis.PutToRedis(types.CourseNamePre+course.CourseId, course.Name, -1)
		return course.CourseId, types.OK
	}
}

//查询课程信息
func (getCourseInfo GetCourseRequest) GetCourseInfo() (user mydb.Course, errno types.ErrNo) {
	log.Println(getCourseInfo)
	var course mydb.Course
	err := mydb.MysqlClient.Model(&course).Where("course_id = ?", getCourseInfo.CourseID).First(&course)
	if err.Error != nil {
		return mydb.Course{}, types.CourseNotExisted
	}
	return course, types.OK
}
