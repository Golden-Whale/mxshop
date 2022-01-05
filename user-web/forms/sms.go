package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"`
	Type   string `from:"type" json:"type" binding:"required,oneof=1 2"`
	// 1. 注册发送短信验证码和动态登录发送短信验证码
}
