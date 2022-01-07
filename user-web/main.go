package main

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/initialize"
	"mxshop-api/user-web/utils"
	myValidator "mxshop-api/user-web/validator"
)

func main() {
	// 1. 初始化logger
	initialize.InitLogger()
	// 2. 初始化配置文件
	initialize.InitConfig()
	// 3. 初始化routers
	router := initialize.Routers()
	// 4. 初始化验证器翻译
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Error("初始化验证器错误", err.Error())
	}
	// 5. 初始化srv的连接
	initialize.InitSrvConn()

	debug := viper.GetInt("MXSHOP")
	// 如果是本地开发环境端口号固定， 线上环境获取随机端口
	if debug == 0 {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	// 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myValidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "非法的手机号码", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}

	/*
		1. S()可以获取一个全局的sugar， 可以让我们自己设置一个全局的logger
		2. 日志是分级别的，debug，info, warn, error, fetal
		3. S函数和L函数很有用
	*/
	zap.S().Debugf("启动服务器, 端口:%d", global.ServerConfig.Port)

	if err := router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动失败", err.Error())
	}
}
