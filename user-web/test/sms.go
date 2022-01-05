package main

import (
	"github.com/cloopen/go-sms-sdk/cloopen"
	"log"
)

func main() {
	cfg := cloopen.DefaultConfig().
		// 开发者主账号,登陆云通讯网站后,可在控制台首页看到开发者主账号ACCOUNT SID和主账号令牌AUTH TOKEN
		WithAPIAccount("8a216da87a332d53017a5688823a0f64").
		// 主账号令牌 TOKEN,登陆云通讯网站后,可在控制台首页看到开发者主账号ACCOUNT SID和主账号令牌AUTH TOKEN
		WithAPIToken("59f63148e6094f609a4cbcbd76138a56")
	sms := cloopen.NewJsonClient(cfg).SMS()
	// 下发包体参数
	input := &cloopen.SendRequest{
		// 应用的APPID
		AppId: "8a216da87a332d53017a5688833a0f6a",
		// 手机号码
		To: "130026725",
		// 模版ID
		TemplateId: "1",
		// 模版变量内容 非必填
		Datas: []string{"12345", "10"},
	}
	// 下发
	resp, err := sms.Send(input)
	if err != nil {
		log.Fatal("短信发送错误", err)
		return
	}
	log.Printf("Response MsgId: %s \n", resp.TemplateSMS.SmsMessageSid)
}
