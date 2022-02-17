package student

import (
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/pkg/rabbitmq"
	"Course-Selection-Scheduling/types"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"gorm.io/gorm"
	"log"
)

type BookCourseRequest struct {
	StudentID string
	CourseID  string
}

// 课程已满返回 CourseNotAvailable

type BookCourseResponse struct {
	Code types.ErrNo
}

type GetStudentCourseRequest struct {
	StudentID string
}

type GetStudentCourseResponse struct {
	Code types.ErrNo
	Data struct {
		CourseList []types.TCourse
	}
}

//检查课程学生合法性
func (bookCourseInfo BookCourseRequest) CheckValid() (errno types.ErrNo) {
	//课程不存在
	log.Println(types.CoursePre + bookCourseInfo.CourseID)
	value, _ := myredis.Exsits(types.CoursePre + bookCourseInfo.CourseID)
	if value == false {
		return types.CourseNotExisted
	}
	//学生不存在
	value, _ = myredis.Exsits(types.StudentPre + bookCourseInfo.StudentID)
	if value == false {
		return types.StudentNotExisted
	}
	return types.OK
}

//限制重复抢课和抢课频度
func (bookCourseInfo BookCourseRequest) CheckRestriction(success, frequency string) (errno types.ErrNo) {
	//限制抢课频率
	value, _ := myredis.Exsits(frequency)
	if value == true {
		return types.RepeatRequest
	} else {
		myredis.PutToRedis(frequency, "true", 3) //3秒内只能抢一次
	}
	//限制重复抢课
	cnt, _ := redis.Int(myredis.DecrForRedis(success))
	//0-1=-1为初次抢课
	if cnt < (-1) {
		myredis.IncrForRedis(success)
		return types.StudentHasCourse
	}
	return types.OK
}

//锁定课程
func (bookCourseInfo BookCourseRequest) LockCourse(success string) (errno types.ErrNo) {
	//查询课程余量并减库存 , 数据库操作送入消息队列中
	value, err := myredis.DecrForRedis(types.CoursePre + bookCourseInfo.CourseID)
	if err != nil {
		myredis.IncrForRedis(success) //锁定失败
		return types.UnknownError
	}
	if value.(int64) < 0 {
		myredis.IncrForRedis(types.CoursePre + bookCourseInfo.CourseID) //加回来
		myredis.IncrForRedis(success)                                   //锁定失败
		return types.CourseNotAvailable
	}
	myredis.SAddToRedisSet(types.SelectPre+bookCourseInfo.StudentID, bookCourseInfo.CourseID)
	//放到消息队列中,进行数据库操作
	msgByte, err := json.Marshal(bookCourseInfo)
	if err != nil {
		log.Fatalln(err)
	}
	rabbitmq.RMQClient.PublishSimple(msgByte)
	return types.OK
}

//检查课程学生合法性
func (studentCourseInfo GetStudentCourseRequest) CheckStudent() (errno types.ErrNo) {
	//学生不存在
	value, _ := myredis.Exsits(types.StudentPre + studentCourseInfo.StudentID)
	if value == false {
		return types.StudentNotExisted
	}
	return types.OK
}

//限制频度
func (bookCourseInfo GetStudentCourseRequest) CheckRestriction(frequency string) (errno types.ErrNo) {
	//限制频率
	value, _ := myredis.Exsits(frequency)
	if value == true {
		return types.RepeatRequest
	} else {
		myredis.PutToRedis(frequency, "true", 3) //3秒内只能抢查一次
	}
	return types.OK
}

//得到课程表
func (studentCourseInfo GetStudentCourseRequest) GetCourses() (CourseList []types.TCourse, errno types.ErrNo) {
	courseIDs, _ := myredis.SGetAllFromRedis(types.SelectPre + studentCourseInfo.StudentID)
	log.Println(courseIDs)
	CourseList = make([]types.TCourse, len(courseIDs))
	if len(CourseList) == 0 {
		return CourseList, types.StudentHasNoCourse
	}
	for i, courseID := range courseIDs {
		CourseList[i].CourseID = courseID
		CourseList[i].Name, _ = redis.String(myredis.GetFromRedis(types.CourseNamePre + courseID))
		CourseList[i].TeacherID, _ = redis.String(myredis.GetFromRedis(types.TeacherIDPre + courseID))
	}
	return CourseList, types.OK
}

func Consume(msgByte []byte) {
	//解析message
	var msg BookCourseRequest
	err := json.Unmarshal(msgByte, &msg)
	if err != nil {
		log.Fatalln(err)
	}
	//扣减课程余量
	mydb.MysqlClient.Model(&mydb.Course{}).
		Where("course_id = ?", msg.CourseID).
		Update("cap", gorm.Expr("cap- ?", 1))
	//插入课表
	mydb.MysqlClient.Create(&mydb.SelectCourse{StudentId: msg.StudentID, CourseId: msg.CourseID})
}
