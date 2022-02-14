package login

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/config"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

//验证用户名与密码
func signUsername(username string, password string) global.ErrNo {
	db := mydb.NewMysqlConn(&config.MysqlCfg)
	usernameLen := len(username)
	passwordLen := len(password)
	if usernameLen < 9 || usernameLen > 20 || passwordLen < 9 || passwordLen > 20 {
		return global.ParamInvalid
	}
	var user mydb.User
	db.Unscoped().Where("username = ?", username).First(&user)
	if user.DeletedAt.Valid {
		return global.UserHasDeleted
	} else if user.Username != username {
		return global.UserNotExisted
	} else if  user.Password != password{
		return global.WrongPassword
	} else {
		return global.OK
	}
}


//登录
func Login(c *gin.Context) {
	var json global.LoginRequest
	var res global.LoginResponse
	err := c.ShouldBindJSON(&json)
	if err != nil {
		fmt.Println("json fail ", err)
		return
	}
	//print(json.Username, json.Password)
	res.Code = signUsername(json.Username, json.Password)
	if res.Code != global.OK {
		c.JSON(401, &res)
		return
	}
	u1, _ := uuid.NewUUID()
	http.SetCookie(c.Writer, &http.Cookie{
		Name : "camp-session",
		Value: u1.String(),
		Path: "/api/v1",
	})
	//c.SetCookie("camp-session", u1.String(), -1, "/", "localhost" ,false, false)
	global.RedisClient = myredis.NewRedisClient(&config.RedisCfg)
	rdb := global.RedisClient.Get()
	rdb.Do("SET", u1.String(), json.Username)
	rdb.Do("EXPIRE", u1.String(), 60*60)
	c.JSON(401, &res)
}