package server

import (
	"Course-Selection-Scheduling/pkg/config"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path"
	"time"
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

	//生成日志
	logPath := path.Join(config.LogCfg.Path, time.Now().String())
	logFile, err := os.Create(logPath)
	if err != nil {
		panic(fmt.Sprintf("create logs failed, path: %v. err: %v", logPath, err.Error()))
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	//恢复恐慌
	httpServer.Use(gin.Recovery())
	//注册路由
	registerRouter(httpServer)

	err = httpServer.Run(serverCfg.Host + ":" + serverCfg.Port)
	if err != nil {
		panic(fmt.Sprintf("server error,address: %v. err: %v", serverCfg.Host+":"+serverCfg.Port, err.Error()))
	}
}
