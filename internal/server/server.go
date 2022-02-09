package server

import (
	"Course-Selection-Scheduling/pkg/config"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Run() {
	serverCfg := config.ServerCfg
	sessionCfg := config.SessionCfg
	//设置模式
	gin.SetMode(serverCfg.Mode)
	httpServer := gin.Default()
	//创建session存储引擎
	sessionStore := cookie.NewStore([]byte(sessionCfg.Key))
	sessionStore.Options(sessions.Options{
		MaxAge: sessionCfg.Age,
	})
	//使用session中间件
	httpServer.Use(sessions.Sessions(sessionCfg.Name, sessionStore))
	//恢复恐慌
	httpServer.Use(gin.Recovery())
	//注册路由
	registerRouter(httpServer)

	err := httpServer.Run(serverCfg.Host + ":" + serverCfg.Port)
	if err != nil {
		panic(fmt.Sprintf("server error,address: %v. err: %v", serverCfg.Host+":"+serverCfg.Port, err.Error()))
	}
}
