package member

import (
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/types"
	"Course-Selection-Scheduling/utils"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"strconv"
)

// -----------------------------------
type TMember struct {
	UserID   string
	Nickname string
	Username string
	UserType types.UserType
}

// 系统内置管理员账号
// 账号名：JudgeAdmin 密码：JudgePassword2022

// 创建成员
// 参数不合法返回 ParamInvalid

// 只有管理员才能添加

type CreateMemberRequest struct {
	Nickname string         `binding:"required,min=4,max=20"`      // required，不小于 4 位 不超过 20 位
	Username string         `binding:"required,UserNameValidator"` // required，只支持大小写，长度不小于 8 位 不超过 20 位
	Password string         `binding:"required,PasswordValidator"` // required，同时包括大小写、数字，长度不少于 8 位 不超过 20 位
	UserType types.UserType `binding:"UserTypeValidator"`          // required, 枚举值
}

type CreateMemberResponse struct {
	Code types.ErrNo
	Data struct {
		UserID string // int64 范围
	}
}

// 获取成员信息

type GetMemberRequest struct {
	UserID string `binding:"required,IsDigitValidator"`
}

// 如果用户已删除请返回已删除状态码，不存在请返回不存在状态码

type GetMemberResponse struct {
	Code types.ErrNo
	Data TMember
}

// 批量获取成员信息

type GetMemberListRequest struct {
	Offset string `binding:"required"`
	Limit  string `binding:"required"`
}

type GetMemberListResponse struct {
	Code types.ErrNo
	Data struct {
		MemberList []TMember
	}
}

// 更新成员信息

type UpdateMemberRequest struct {
	UserID   string `binding:"required,IsDigitValidator"`
	Nickname string `binding:"required,min=4,max=20"`
}

type UpdateMemberResponse struct {
	Code types.ErrNo
}

// 删除成员信息
// 成员删除后，该成员不能够被登录且不应该不可见，ID 不可复用

type DeleteMemberRequest struct {
	UserID string `binding:"required,IsDigitValidator"`
}

type DeleteMemberResponse struct {
	Code types.ErrNo
}

// -----------------------------------

//检验是否为admin
func (createMemberInfo CreateMemberRequest) CheckAdmin(c *gin.Context) (errno types.ErrNo) {
	val, err := c.Cookie(types.CampSession)
	if err != nil {
		return types.PermDenied
	}
	//检验是否有权限创建成员
	name, err := redis.String(myredis.GetFromRedis(val))
	if err != nil {
		return types.PermDenied
	}
	var mem mydb.User // 当前操作者的身份
	mydb.MysqlClient.Where("username = ?", name).First(&mem)
	if mem.UserType != types.Admin {
		return types.PermDenied
	}
	return types.OK
}

//创建用户
func (createMemberInfo CreateMemberRequest) CreateUser() (UserID string, errno types.ErrNo) {
	var user mydb.User
	db := mydb.MysqlClient
	db.Unscoped().Where("username = ?", createMemberInfo.Username).First(&user)
	//检验用户名是否已经存在
	if user.Username == createMemberInfo.Username {
		if user.DeletedAt.Valid {
			return "", types.UserHasDeleted
		} else {
			return "", types.UserHasExisted
		}
	} else {
		user = mydb.User{
			Username: createMemberInfo.Username,
			Nickname: createMemberInfo.Nickname,
			Password: createMemberInfo.Password,
			UserType: createMemberInfo.UserType,
		}
		db.Create(&user)
		if createMemberInfo.UserType == types.Student {
			myredis.PutToRedis(types.StudentPre+user.UserId, "true", -1)
		}
		return user.UserId, types.OK
	}
}

//查询个人信息
func (getMemberInfo GetMemberRequest) GetPersonInfo() (user mydb.User, errno types.ErrNo) {
	mydb.MysqlClient.Unscoped().Where("user_id = ?", getMemberInfo.UserID).First(&user)
	if user.UserId == getMemberInfo.UserID {
		if user.DeletedAt.Valid {
			return mydb.User{}, types.UserHasDeleted
		} else {
			return user, types.OK
		}
	} else {
		return mydb.User{}, types.UserNotExisted
	}
}

//更新个人信息
func (updateInfo UpdateMemberRequest) Update() (errno types.ErrNo) {
	var user mydb.User
	db := mydb.NewMysqlConn(&utils.MysqlCfg)
	db.Unscoped().Where("user_id = ?", updateInfo.UserID).First(&user)
	if user.UserId == updateInfo.UserID {
		if user.DeletedAt.Valid {
			return types.UserHasDeleted
		} else {
			user.Nickname = updateInfo.Nickname
			db.Save(&user)
			return types.OK
		}
	} else {
		return types.UserNotExisted
	}
}

//删除账户
func (deleteInfo DeleteMemberRequest) Delete() (errno types.ErrNo) {
	var user mydb.User
	db := mydb.MysqlClient
	db.Where("user_id = ?", deleteInfo.UserID).First(&user)
	if user.UserId == deleteInfo.UserID {
		if user.UserType == types.Student {
			myredis.DeleteFromRedis("student_" + user.UserId)
		}
		db.Delete(&user)
		return types.OK
	} else {
		return types.UserNotExisted
	}
}

//查询列表信息
func (memberListInfo GetMemberListRequest) GetListInfo() (users []mydb.User, errno types.ErrNo) {
	offset, _ := strconv.Atoi(memberListInfo.Offset)
	limit, _ := strconv.Atoi(memberListInfo.Limit)
	mydb.MysqlClient.Model(&mydb.User{}).Offset(offset).Limit(limit).Find(&users)
	return users, types.OK
}
