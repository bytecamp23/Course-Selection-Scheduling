package login

import (
	"Course-Selection-Scheduling/internal/global"
	"Course-Selection-Scheduling/pkg/myredis"
	"github.com/gin-gonic/gin"
	"net/http"
)

//登出
func Logout(c *gin.Context) {
	var res global.LogoutResponse
	data, err := c.Cookie("camp-session")
	if err != nil {
		res.Code = global.LoginRequired
		c.JSON(http.StatusUnauthorized, &res)
		return
	}
	myredis.DeleteFromRedis(data)
	c.SetCookie("camp-session", data, -1, "/", "localhost", false, true)
	c.JSON(200, &res)
}
