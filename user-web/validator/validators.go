package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	// 使用正则表达式判断是否合法
	if ok, _ := regexp.MatchString(`^1([38][0-9]|14[57]|5[^4])\d{8}$`, mobile); !ok {
		return false
	}
	return true
}
