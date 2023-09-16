package middleware

import (
	"errors"
	"orderingsystem/app/common/response"
	"orderingsystem/app/services"
	"orderingsystem/global"
	"orderingsystem/utils"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 调用方式 authRouter := router.Group("").Use(middleware.JWTAuth(services.AppGuardName))

// 首先封装验证 JWT 基础功能
func ValidJWT(tokenStr string) (token *jwt.Token, claims *services.CustomClaims, err error) {
	token, err = jwt.ParseWithClaims(tokenStr, &services.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.App.Config.Jwt.Secret), nil
	})

	if err != nil || !token.Valid {
		_ = token
		err = errors.New("认证失败")
		return
	}
	claims = token.Claims.(*services.CustomClaims)
	return
}

// 根据JWT 获取的用户信息获取用户权限，细化中间件
func JWTAuth(auth string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			response.TokenFail(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}
		tokenStr = tokenStr[len(services.TokenType)+1:]
		token, claims, err := ValidJWT(tokenStr)

		if err != nil || services.JwtService.IsInBlacklist(tokenStr) {
			response.TokenFail(c, "登录已过期，请重新登录")
			c.Abort()
			return
		}

		userRoles := claims.Roles
		userStatus := claims.UserStatus
		if !userStatus {
			response.TokenFail(c, "对不起,您已被封禁，请联系管理员")
			c.Abort()
			return
		}
		if auth == "" {
			// 普通接口
			c.Set("id", claims.Id)
			c.Set("token", token)
			c.Set("userTel", claims.UserTel)
			return
		}
		// 判断角色
		is_right := false
		for _, v := range userRoles {
			if v == auth {
				is_right = true
				break
			}
		}
		if !is_right {
			response.TokenFail(c, "对不起,您没有权限")
			c.Abort()
			return
		}
		// token 续签 utils.MD5([]byte(tokenStr))
		if claims.ExpiresAt-time.Now().Unix() < global.App.Config.Jwt.RefreshGracePeriod {
			lock := global.Lock(utils.MD5([]byte(claims.UserName+claims.UserTel)), global.App.Config.Jwt.JwtBlacklistGracePeriod)
			if lock.Get() {
				// 需要创建新token，需要user对象
				userId, _ := strconv.ParseInt(claims.Id, 10, 64)
				err, user := services.UserServices.GetUserInfoById(userId)
				if err != nil {
					global.App.Log.Error(err.Error())
					lock.Release()
				} else {
					tokenData, _, _ := services.JwtService.CreateToken(services.AppFuardName, user)
					c.Header("new-token", tokenData.AccessToken)
					c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
					_ = services.JwtService.JoinBlackList(token)
				}
			}
		}

		c.Set("id", claims.Id)
		c.Set("token", token)
		c.Set("userTel", claims.UserTel)
	}
}
