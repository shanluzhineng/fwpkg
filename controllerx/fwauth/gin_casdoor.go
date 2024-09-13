package fwauth

import (
	"strings"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/shanluzhineng/fwpkg/controllerx/responsex"
)

func GinCasdoorHandler() gin.HandlerFunc {
	// jwt 认证 - 8020逻辑，可去掉
	// jwtFunc := func(c *gin.Context) {
	// 	// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
	// 	token := c.Request.Header.Get("x-token")
	// 	if token == "" {
	// 		web.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问2", c)
	// 		c.Abort()
	// 		return
	// 	}
	// 	if service.GetServiceGroup().JwtService.IsBlacklist(token) {
	// 		web.FailWithDetailed(gin.H{"reload": true}, "您的帐户异地登陆或令牌失效", c)
	// 		c.Abort()
	// 		return
	// 	}
	// 	j := web.NewJWT()
	// 	// parseToken 解析token包含的信息
	// 	claims, err := j.ParseToken(token)
	// 	if err != nil {
	// 		if err == web.TokenExpired {
	// 			web.FailWithDetailed(gin.H{"reload": true}, "授权已过期", c)
	// 			c.Abort()
	// 			return
	// 		}
	// 		web.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
	// 		c.Abort()
	// 		return
	// 	}
	// 	// 用户被删除的逻辑 需要优化 此处比较消耗性能 如果需要 请自行打开
	// 	//if err, _ = userService.FindUserByUuid(claims.UUID.String()); err != nil {
	// 	//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
	// 	//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
	// 	//	c.Abort()
	// 	//}
	// 	if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
	// 		claims.ExpiresAt = time.Now().Unix() + config.AppConfig.JWT.ExpiresTime
	// 		newToken, _ := j.CreateTokenByOldToken(token, *claims)
	// 		newClaims, _ := j.ParseToken(newToken)
	// 		c.Header("new-token", newToken)
	// 		c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
	// 		if config.AppConfig.System.UseMultipoint {
	// 			RedisJwtToken, err := service.GetServiceGroup().JwtService.GetRedisJWT(newClaims.Username)
	// 			if err != nil {
	// 				log.Logger.Error("get redis jwt failed", zap.Error(err))
	// 			} else { // 当之前的取成功时才进行拉黑操作
	// 				_ = service.GetServiceGroup().JwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: RedisJwtToken})
	// 			}
	// 			// 无论如何都要记录当前的活跃状态
	// 			_ = service.GetServiceGroup().JwtService.SetRedisJWT(newToken, newClaims.Username)
	// 		}
	// 	}
	// 	c.Set("claims", claims)
	// 	c.Next()
	// }
	// casdoor 认证
	casdoorFunc := func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// jwtFunc(c)
			return
		}

		token := strings.Split(authHeader, "Bearer ")
		if len(token) != 2 {
			// jwtFunc(c)
			return
		}

		claims, err := casdoorsdk.ParseJwtToken(token[1])
		if err != nil {
			responsex.FailWithDetailed(gin.H{}, "ParseJwtToken() error: "+err.Error(), c)
			return
		}
		// record login info into current ctx, transfer to next handler
		c.Set(AuthKey, claims)
		// Passthrough to next handler if needed
		c.Next()
	}
	return casdoorFunc
}
