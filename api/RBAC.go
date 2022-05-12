package api

import (
	"course-selection/service"
	"course-selection/tool"
	"github.com/gin-gonic/gin"
	"github.com/storyicon/grbac"
	"log"
	"strings"
	"time"
)

func Authorization(ctx *gin.Context) {
	rbac, err := grbac.New(grbac.WithJSON("config/ruleConfig.json", time.Minute*10))
	if err != nil {
		tool.Failure(ctx, 500, "服务器错误")
		log.Fatal("解析权限配置文件失败", err)
	}
	//获取并解析token
	tokenString := ctx.Request.Header.Get("token")
	if tokenString == "" {
		tool.Failure(ctx, 401, "token不能为空")
		return
	}
	tokenClaims, err := service.ParseToken(tokenString)
	if err != nil {
		tool.Failure(ctx, 500, "服务器错误")
		ctx.Abort()
		log.Fatal("token解析失败", err)
	}
	//获取角色权限等级
	roleLevel, err := service.HashGet(tokenClaims.Identify, "roleLevel")
	if err != nil {
		tool.Failure(ctx, 500, "服务器错误")
		log.Fatal("查询角色权限等级失败", err)
	}
	roles := strings.Split(roleLevel, ",")
	//鉴权
	state, err := rbac.IsRequestGranted(ctx.Request, roles)
	if err != nil {
		tool.Failure(ctx, 500, "服务器错误")
		log.Fatal("鉴权失败", err)
	}

	if state.IsGranted() {
		tool.Failure(ctx, 400, "不相互打扰是你的温柔")
		ctx.Abort()
	}
}
