package router

import (
	"github.com/gin-gonic/gin"
	"mxshop-api/user-web/api"
	"mxshop-api/user-web/middlewares"
)

func InitUserRouter(router *gin.RouterGroup) {
	userRouter := router.Group("user")

	{
		userRouter.GET("/list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		userRouter.POST("/pwd_login", api.PasswordLogin)
		userRouter.POST("/register", api.Register)
	}
}
