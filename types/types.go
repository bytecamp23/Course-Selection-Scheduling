package types

const (
	CampSession   string = "camp-session"
	ChooseFactor  int    = 10000000
	SelectPre     string = "sel_"
	StudentPre    string = "stu_"
	CoursePre     string = "cour_"
	CourseNamePre string = "cour_name_"
	TeacherIDPre  string = "tea_id_"
	RMQName       string = "bookcourse"
)

// 说明：
// 1. 所提到的「位数」均以字节长度为准
// 2. 所有的 ID 均为 int64（以 string 方式表现）

// 通用结构

type ErrNo int

const (
	OK                 ErrNo = 0
	ParamInvalid       ErrNo = 1   // 参数不合法
	UserHasExisted     ErrNo = 2   // 该 Username 已存在
	UserHasDeleted     ErrNo = 3   // 用户已删除
	UserNotExisted     ErrNo = 4   // 用户不存在
	WrongPassword      ErrNo = 5   // 密码错误
	LoginRequired      ErrNo = 6   // 用户未登录
	CourseNotAvailable ErrNo = 7   // 课程已满
	CourseHasBound     ErrNo = 8   // 课程已绑定过
	CourseNotBind      ErrNo = 9   // 课程未绑定过
	PermDenied         ErrNo = 10  // 没有操作权限
	StudentNotExisted  ErrNo = 11  // 学生不存在
	CourseNotExisted   ErrNo = 12  // 课程不存在
	StudentHasNoCourse ErrNo = 13  // 学生没有课程
	StudentHasCourse   ErrNo = 14  // 学生有课程
	RepeatRequest      ErrNo = 254 // 重复请求
	UnknownError       ErrNo = 255 // 未知错误
)

type TCourse struct {
	CourseID  string
	Name      string
	TeacherID string
}
type UserType int

const (
	Admin   UserType = 1
	Student UserType = 2
	Teacher UserType = 3
)
