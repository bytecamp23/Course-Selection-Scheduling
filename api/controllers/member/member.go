package member

import (
	"Course-Selection-Scheduling/api/models/member"
	"Course-Selection-Scheduling/pkg/mydb"
	"Course-Selection-Scheduling/types"
	"github.com/gin-gonic/gin"
	"log"
)

//创建成员
func CreateMember(c *gin.Context) {
	var (
		requestData member.CreateMemberRequest
		respondData member.CreateMemberResponse
	)
	respondData.Code = requestData.CheckAdmin(c)
	if respondData.Code != types.OK {
		c.JSON(200, &respondData)
		log.Println(respondData)
		return
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(respondData)
		return
	}
	respondData.Data.UserID, respondData.Code = requestData.CreateUser()
	c.JSON(200, respondData)
	log.Println(respondData)
}

//获取成员信息
func GetMember(c *gin.Context) {
	var (
		requestData member.GetMemberRequest
		respondData member.GetMemberResponse
	)
	if err := c.ShouldBindQuery(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(respondData)
		return
	}
	var user mydb.User
	user, respondData.Code = requestData.GetPersonInfo()
	respondData.Data.Nickname = user.Nickname
	respondData.Data.UserID = user.UserId
	respondData.Data.UserType = user.UserType
	respondData.Data.Username = user.Username
	c.JSON(200, &respondData)
	log.Println(respondData)
}

//更新成员
func UpdateMember(c *gin.Context) {
	var (
		requestData member.UpdateMemberRequest
		respondData member.UpdateMemberResponse
	)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(respondData)
		return
	}
	respondData.Code = requestData.Update()
	c.JSON(200, &respondData)
	log.Println(respondData)
}

//删除成员
func DeleteMember(c *gin.Context) {
	var (
		requestData member.DeleteMemberRequest
		respondData member.DeleteMemberResponse
	)
	if err := c.ShouldBindJSON(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(respondData)
		return
	}
	respondData.Code = requestData.Delete()
	c.JSON(200, &respondData)
	log.Println(respondData)
}

//批量获取成员
func ListMember(c *gin.Context) {
	var (
		requestData member.GetMemberListRequest
		respondData member.GetMemberListResponse
	)
	if err := c.ShouldBindQuery(&requestData); err != nil {
		respondData.Code = types.ParamInvalid
		c.JSON(200, respondData)
		log.Println(respondData)
		return
	}

	var users []mydb.User
	users, respondData.Code = requestData.GetListInfo()
	respondData.Data.MemberList = make([]member.TMember, len(users))
	for i := 0; i < len(users); i++ {
		respondData.Data.MemberList[i].UserType = users[i].UserType
		respondData.Data.MemberList[i].Nickname = users[i].Nickname
		respondData.Data.MemberList[i].UserID = users[i].UserId
		respondData.Data.MemberList[i].Username = users[i].Username
	}
	c.JSON(200, &respondData)
	log.Println(respondData)
}
