package login

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/mydb"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

//获取成员信息
func GetMember(c *gin.Context) {
	var json global.GetMemberRequest
	var res global.GetMemberResponse
	var user mydb.User
	_ = c.ShouldBindJSON(&json)
	global.MysqlClient.Unscoped().Where("user_id = ?", json.UserID).First(&user)
	if user.UserId == json.UserID {
		if user.DeletedAt.Valid {
			res.Code = global.UserHasDeleted
		} else {
			res.Code = global.OK
			res.Data.UserID = user.UserId
			res.Data.UserType = user.UserType
			res.Data.Nickname = user.Nickname
			res.Data.Username = user.Username
		}
	} else {
		res.Code = global.UserNotExisted
	}
	c.JSON(200, &res)
}

func Whoami(c *gin.Context) {
	data, err := c.Cookie("camp-session")
	var res global.WhoAmIResponse
	if err != nil {
		res.Code = global.LoginRequired
		c.JSON(401, &res)
		return
	}
	username, err := redis.String(global.RedisClient.Get().Do("GET", data))
	if err != nil {
		res.Code = global.LoginRequired
		c.JSON(401, &res)
		return
	}
	var user mydb.User
	global.MysqlClient.Where("username = ?", username).First(&user)
	if username == user.Username {
		res.Code = global.OK
		res.Data.Username = user.Username
		res.Data.UserType = user.UserType
		res.Data.Nickname = user.Nickname
		res.Data.UserID = user.UserId
		c.JSON(200, &res)
	}  else {
		res.Code = global.LoginRequired
		c.JSON(401, &res)
	}
}