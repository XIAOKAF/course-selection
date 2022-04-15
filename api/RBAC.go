package api

import (
	"course-selection/service"
	"course-selection/tool"
	"github.com/gin-gonic/gin"
	"github.com/storyicon/grbac"
	"log"
	"net/http"
	"strings"
	"time"
)

func Authorization(ctx *gin.Context) {
	//获取用户账号
	userNumber := ctx.PostForm("userNumber")
	//获取角色权限等级
	err, roleLevel := service.SelectRoleLevel(userNumber)
	roles := strings.Fields(roleLevel)
	if err != nil {
		tool.Failure(ctx, 500, "服务器错误")
		log.Fatal("查询用户等级错误", err)
		return
	}
	//通过权限等级获取权限（解析yaml文件中的权限配置
	//判断角色是否有权限
	//以一分钟的频率获取最新的身份
	rbac, err := grbac.New(grbac.WithLoader(service.ParseRule, time.Minute))
	if err != nil {
		tool.Failure(ctx, 400, "你还没有这个权限哦")
		log.Fatal("解析权限配置文件错误", err)
		return
	}
	state, err := rbac.IsRequestGranted(ctx.Request, roles)
	if err != nil {
		tool.Failure(ctx, 500, "服务器错误")
		log.Fatal("查询用户等级错误", err)
		return
	}
	if !state.IsGranted() {
		tool.Failure(ctx, http.StatusUnauthorized, "未满十八岁🈲止访问")
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

}
