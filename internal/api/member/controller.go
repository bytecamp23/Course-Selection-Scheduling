package member

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/config"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"strconv"
)

//创建成员接口
func CreateMember(c *gin.Context) {
	val, err := c.Cookie("camp-session")
	var res global.CreateMemberResponse
	if err != nil {
		res.Code = global.PermDenied
		c.JSON(200, &res)
		return
	}

	//检验是否有权限创建成员
	name, err := redis.String(myredis.GetFromRedis(val))
	if err != nil {
		fmt.Println(err)
		res.Code = global.PermDenied
		c.JSON(200, &res)
		return
	}
	var mem mydb.User // 当前操作者的身份
	db := global.MysqlClient
	db.Where("username = ?", name).First(&mem)
	if mem.UserType != global.Admin {
		res.Code = global.PermDenied
		c.JSON(200, &res)
		return
	}

	var json global.CreateMemberRequest
	var user mydb.User
	if err := c.ShouldBindJSON(&json); err != nil {
		fmt.Println(err.Error())
		// TODO: ParamInvalid
		createMemberResponse := global.CreateMemberResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, createMemberResponse)
		return
	}
	/*errs := validate.Struct(json)
	if errs != nil {
		fmt.Println(errs.Errors)
		createMemberResponse := global.CreateMemberResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, createMemberResponse)
		return
	}*/
	db.Unscoped().Where("username = ?", json.Username).First(&user)
	//检验用户名是否已经存在
	if user.Username == json.Username {
		res.Code = global.UserHasExisted
		res.Data.UserID = user.UserId
		c.JSON(200, &res)
	} else {
		res.Code = global.OK
		res.Data.UserID = saveMember(json.Nickname, json.Username, json.Password, json.UserType)
		if json.UserType == global.Student {
			myredis.PutToRedis("student_"+res.Data.UserID, "true", -1)
		}
		c.JSON(200, &res)
	}
}

//创建成员信息
func saveMember(nickname string, username string, password string, usertype global.UserType) string {
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	var user mydb.User
	user.Username = username
	user.Nickname = nickname
	user.Password = password
	user.UserType = usertype
	fmt.Println(user)
	db.Create(&user)
	fmt.Println(user.UserId)
	return user.UserId
}

//更新成员
func UpdateMember(c *gin.Context) {
	var json global.UpdateMemberRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		// TODO: ParamInvalid
		fmt.Println(err)
		updateMemberResponse := global.UpdateMemberResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, updateMemberResponse)
		return
	}
	/*errs := validate.Struct(json)
	if errs != nil {
		createMemberResponse := global.CreateMemberResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, createMemberResponse)
		return
	}*/
	var user mydb.User
	var res global.UpdateMemberResponse
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	db.Unscoped().Where("user_id = ?", json.UserID).First(&user)
	if user.UserId == json.UserID {
		if user.DeletedAt.Valid {
			res.Code = global.UserHasDeleted
		} else {
			user.Nickname = json.Nickname
			db.Save(&user)
			res.Code = global.OK
		}
	} else {
		res.Code = global.UserNotExisted
	}
	c.JSON(200, &res)
}

//删除成员
func DeleteMember(c *gin.Context) {
	var json global.DeleteMemberRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		// TODO: ParamInvalid
		deleteMemberResponse := global.DeleteMemberResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, deleteMemberResponse)
		return
	}
	var user mydb.User
	var res global.DeleteMemberResponse
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	db.Where("user_id = ?", json.UserID).First(&user)
	if user.UserId == json.UserID {
		if user.UserType == global.Student {
			myredis.DeleteFromRedis("student_" + user.UserId)
		}
		db.Delete(&user)
		res.Code = global.OK
	} else {
		res.Code = global.UserNotExisted
	}
	c.JSON(200, &res)
}

//批量获取成员
func ListMember(c *gin.Context) {
	var json global.GetMemberListRequest
	if err := c.ShouldBindQuery(&json); err != nil {
		// TODO: ParamInvalid
		getMemberResponse := global.GetMemberListResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, getMemberResponse)
		return
	}
	fmt.Println(json)
	var users []mydb.User
	var res global.GetMemberListResponse
	offset, _ := strconv.Atoi(json.Offset)
	limit, _ := strconv.Atoi(json.Limit)
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	db.Model(&mydb.User{}).Offset(offset).Limit(limit).Find(&users)
	res.Data.MemberList = make([]global.TMember, len(users))
	for i := 0; i < len(users); i++ {
		res.Data.MemberList[i].UserType = users[i].UserType
		res.Data.MemberList[i].Nickname = users[i].Nickname
		res.Data.MemberList[i].UserID = users[i].UserId
		res.Data.MemberList[i].Username = users[i].Username
	}
	res.Code = global.OK
	c.JSON(200, &res)
}

// 获取单个成员
func GetMember(c *gin.Context) {
	var json global.GetMemberRequest
	if err := c.ShouldBindQuery(&json); err != nil {
		// TODO: ParamInvalid
		getMemberResponse := global.GetMemberResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(200, getMemberResponse)
		return
	}
	var user mydb.User
	var res global.GetMemberResponse
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	db.Model(&mydb.User{}).Where("user_id = ?", json.UserID).First(&user)
	res.Data.UserType = user.UserType
	res.Data.Nickname = user.Nickname
	res.Data.UserID = user.UserId
	res.Data.Username = user.Username
	res.Code = global.OK
	c.JSON(200, &res)
}
