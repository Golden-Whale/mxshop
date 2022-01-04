package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/forms"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// RemoveTopStruct 剔除errs前缀
func RemoveTopStruct(errs map[string]string) map[string]string {
	newErrors := make(map[string]string, 0)
	for key, val := range errs {
		newErrors[key[strings.Index(key, ".")+1:]] = val
	}
	return newErrors
}

func HandelGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc的code转换成Http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"msg": e.Message()})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		}
	}
}

func HandelValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
}

func GetUserList(ctx *gin.Context) {
	zap.S().Debug("获取用户列表页")

	// 获取参数
	pn := ctx.DefaultQuery("pn", "1")
	pnInt, _ := strconv.Atoi(pn)

	pSize := ctx.DefaultQuery("psize", "1")
	pSizeInt, _ := strconv.Atoi(pSize)

	// 拨号连接用户grpc服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	// 调用接口
	userSrvClient := proto.NewUserClient(userConn)

	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandelGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Gender: value.NickName,
			Mobile: value.Mobile,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
	return
}

func PasswordLogin(ctx *gin.Context) {
	// 表单验证
	passwordForm := forms.PasswordLoginForm{}
	if err := ctx.ShouldBind(&passwordForm); err != nil {
		HandelValidatorError(ctx, err)
		return
	}

	// 拨号连接用户grpc服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	// 调用接口
	userSrvClient := proto.NewUserClient(userConn)

	// 通过手机号码获取用户
	rsp, err := userSrvClient.GetUserByMobile(context.Background(), &proto.MobileReqeust{Mobile: passwordForm.Mobile})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户手机号】失败")
		HandelGrpcErrorToHttp(err, ctx)
		return
	}

	passwordSame, err := userSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Password:          passwordForm.Password,
		EncryptedPassword: rsp.Password,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【验证密码】失败")
		HandelGrpcErrorToHttp(err, ctx)
		return
	}
	if passwordSame.Success {
		ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": nil, "msg": "登录成功"})
	} else {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "data": nil, "msg": "用户名或密码错误"})
	}
}
