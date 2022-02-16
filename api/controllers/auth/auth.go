package auth

import (
	"Course-Selection-Scheduling/api/models/auth"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/types"
	"github.com/gin-gonic/gin"
	"log"
)

// Login 登录
func Login(c *gin.Context) {
	var (
		requestData auth.LoginRequest
		respondData auth.LoginResponse
	)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, &respondData)
		log.Println(respondData)
		return
	}
	respondData.Data.UserID, respondData.Code = requestData.CheckLogin()
	if respondData.Code != types.OK {
		c.JSON(200, &respondData)
		log.Println(respondData)
		return
	}
	requestData.SetSession(c)
	c.JSON(200, &respondData)
	log.Println(respondData)
}

// Logout 登出
func Logout(c *gin.Context) {
	var (
		requestData auth.LogoutRequest
		respondData auth.LogoutResponse
	)
	respondData.Code = requestData.ClearSession(c)
	c.JSON(200, &respondData)
	log.Println(respondData)
}

// Whoami 获取个人信息
func Whoami(c *gin.Context) {
	var (
		requestData auth.WhoAmIRequest
		respondData auth.WhoAmIResponse
	)
	var (
		username string
		user     mydb.User
	)
	username, respondData.Code = requestData.CheckSession(c)
	if respondData.Code != types.OK {
		c.JSON(200, &respondData)
		log.Println(respondData)
		return
	}

	user, respondData.Code = requestData.GetPersonInfo(username)
	if respondData.Code != types.OK {
		c.JSON(200, &respondData)
		log.Println(respondData)
		return
	}
	respondData.Data.Username = user.Username
	respondData.Data.UserType = user.UserType
	respondData.Data.Nickname = user.Nickname
	respondData.Data.UserID = user.UserId
	c.JSON(200, &respondData)
	log.Println(respondData)
}
