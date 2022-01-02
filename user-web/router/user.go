package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop-api/user-web/api"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")
	zap.S().Info("配置用户相关的URL")

	{
		userRouter.GET("/list", api.GetUserList)
		userRouter.POST("/pwd_login", api.PasswordLogin)
	}
}
