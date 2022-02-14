package login

import (
	"Course-Selection-Scheduling/internal/global"
	"github.com/gin-gonic/gin"
	"net/http"
)

//登出
func Logout(c *gin.Context) {
	data, err := c.Cookie("camp-session")
	var res global.LogoutResponse
	if err != nil {
		res.Code = global.LoginRequired
		c.JSON(http.StatusUnauthorized, &res)
		return
	}
	c.SetCookie("camp-session", data, -1, "/","localhost", false, true)
	c.JSON(200, &res)
}
