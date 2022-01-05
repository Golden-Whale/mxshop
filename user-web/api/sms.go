package api

import (
	"context"
	"fmt"
	"github.com/cloopen/go-sms-sdk/cloopen"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"math/rand"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"net/http"
	"strings"
	"time"
)

func GenerateSmsCode(width int) string {
	// 生成width长度的短信验证码
	numberic := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numberic)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numberic[rand.Intn(r)])
	}
	return sb.String()
}

func SendSms(ctx *gin.Context) {
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HandelValidatorError(ctx, err)
		return
	}
	smsCode := GenerateSmsCode(6)

	cfg := cloopen.DefaultConfig().
		// 开发者主账号,登陆云通讯网站后,可在控制台首页看到开发者主账号ACCOUNT SID和主账号令牌AUTH TOKEN
		WithAPIAccount(global.ServerConfig.RLYSmsInfo.APIAccount).
		// 主账号令牌 TOKEN,登陆云通讯网站后,可在控制台首页看到开发者主账号ACCOUNT SID和主账号令牌AUTH TOKEN
		WithAPIToken(global.ServerConfig.RLYSmsInfo.APIToken)
	sms := cloopen.NewJsonClient(cfg).SMS()
	// 下发包体参数
	input := &cloopen.SendRequest{
		// 应用的APPID
		AppId: global.ServerConfig.RLYSmsInfo.APPID,
		// 手机号码
		To: sendSmsForm.Mobile,
		// 模版ID
		TemplateId: "1",
		// 模版变量内容 验证码 过期时间(分钟)
		Datas: []string{GenerateSmsCode(6), "10"},
	}
	// 下发
	_, _ = sms, input
	//resp, err := sms.Send(input)
	//if err != nil {
	//	zap.S().Error("短信发送错误", err)
	//	return
	//}
	//if resp.StatusCode != "000000" {
	//	zap.S().Debugw("发送短信失败, error msg=", resp.StatusMsg, "error code=", resp.StatusCode)
	//	ctx.JSON(http.StatusInternalServerError, gin.H{
	//		"msg": "短信发送失败",
	//	})
	//	return
	//}
	// 保存 验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	rdb.Set(context.Background(), sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
	zap.S().Infof("给手机号%v发送短信 验证码为:%v", sendSmsForm.Mobile, smsCode)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}
