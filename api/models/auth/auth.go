package auth

import (
	"Course-Selection-Scheduling/api/models/member"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/pkg/myredis"
	"Course-Selection-Scheduling/types"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// ----------------------------------------
// 登录
type LoginRequest struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

// 登录成功后需要 Set-Cookie("camp-session", ${value})
// 密码错误范围密码错误状态码
type LoginResponse struct {
	Code types.ErrNo
	Data struct {
		UserID string
	}
}

// 登出
type LogoutRequest struct{}

// 登出成功需要删除 Cookie
type LogoutResponse struct {
	Code types.ErrNo
}

// WhoAmI 接口，用来测试是否登录成功，只有此接口需要带上 Cookie
type WhoAmIRequest struct {
}

// 用户未登录请返回用户未登录状态码
type WhoAmIResponse struct {
	Code types.ErrNo
	Data member.TMember
}

// ----------------------------------------

//验证用户名与密码
func (loginInfo LoginRequest) CheckLogin() (userID string, errno types.ErrNo) {
	var user mydb.User
	mydb.MysqlClient.Unscoped().Where("username = ?", loginInfo.Username).First(&user)
	if user.DeletedAt.Valid {
		return "", types.UserHasDeleted
	} else if user.Username != loginInfo.Username {
		return "", types.UserNotExisted
	} else if user.Password != loginInfo.Password {
		return "", types.WrongPassword
	} else {
		return user.UserId, types.OK
	}
}

//设置session
func (loginInfo LoginRequest) SetSession(c *gin.Context) (errno types.ErrNo) {
	data, _ := c.Cookie(types.CampSession)
	myredis.DeleteFromRedis(data) //删除旧session
	u1, _ := uuid.NewUUID()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:  types.CampSession,
		Value: u1.String(),
		Path:  "/api/v1",
	})
	myredis.PutToRedis(u1.String(), loginInfo.Username, 60*60*24)
	return types.OK
}

//清除session
func (logoutInfo LogoutRequest) ClearSession(c *gin.Context) (errno types.ErrNo) {
	data, _ := c.Cookie(types.CampSession)
	if val, _ := myredis.GetFromRedis(data); val == nil {
		return types.LoginRequired
	}
	myredis.DeleteFromRedis(data)
	c.SetCookie(types.CampSession, data, -1, "/", "localhost", false, true)
	return types.OK
}

//检验session
func (whoamiInfo WhoAmIRequest) CheckSession(c *gin.Context) (username string, errno types.ErrNo) {
	data, err := c.Cookie(types.CampSession)
	if err != nil {
		return "", types.LoginRequired
	}
	username, err = redis.String(myredis.GetFromRedis(data))
	if err != nil {
		return "", types.LoginRequired
	}
	return username, types.OK
}

//查询个人信息
func (whoamiInfo WhoAmIRequest) GetPersonInfo(username string) (user mydb.User, errno types.ErrNo) {
	mydb.MysqlClient.Where("username = ?", username).First(&user)
	if username == user.Username {
		return user, types.OK
	} else {
		return user, types.LoginRequired
	}
}
