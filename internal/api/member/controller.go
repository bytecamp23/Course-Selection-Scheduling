package member

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/config"
	"Course-Selection-Scheduling/pkg/mydb"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
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
	var mem mydb.User // 当前操作者的身份
	name, err := redis.String(global.RedisClient.Get().Do("GET", val))
	if err != nil {
		res.Code = global.PermDenied
		c.JSON(200, &res)
		return
	}
	db := mydb.NewMysqlConn(&config.MysqlCfg)
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
		c.JSON(400, createMemberResponse)
		return
	}

	fmt.Println(json)
	db.Where("username = ?", json.Username).First(&user)
	//检验用户名是否已经存在
	if user.Username == json.Username {
		res.Code = global.UserHasExisted
		res.Data.UserID = user.UserId
		c.JSON(401, &res)
	} else {
		res.Code = global.OK
		res.Data.UserID = saveMember(json.Nickname, json.Username, json.Password, json.UserType)
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
	return user.UserId
}

//更新成员
func UpdateMember(c *gin.Context) {
	var json global.UpdateMemberRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		// TODO: ParamInvalid
		updateMemberResponse := global.UpdateMemberResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(400, updateMemberResponse)
		return
	}
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
		c.JSON(400, deleteMemberResponse)
		return
	}
	var user mydb.User
	var res global.DeleteMemberResponse
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	db.Where("user_id = ?", json.UserID).First(&user)
	if user.UserId == json.UserID {
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
	if err := c.ShouldBindJSON(&json); err != nil {
		// TODO: ParamInvalid
		getMemberResponse := global.GetMemberListResponse{
			Code: global.ParamInvalid,
		}
		c.JSON(400, getMemberResponse)
		return
	}
	var users []mydb.User
	var res global.GetMemberListResponse
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	db.Model(&mydb.User{}).Offset(json.Offset).Limit(json.Limit).Find(&users)
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
