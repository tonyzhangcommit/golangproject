package services

import (
	"context"
	"orderingsystem/app/models"
	"orderingsystem/global"
	"orderingsystem/utils"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtService struct {
}

var JwtService = new(jwtService)

type JwtUser interface {
	GetUid() string
}

type CustomClaims struct {
	Roles      []string
	UserStatus bool
	UserName   string
	UserTel    string
	jwt.StandardClaims
}

const (
	TokenType    = "bearer"
	AppFuardName = "manage"
)

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expire_in"`
	TokenType   string `json:"token_type"`
}

func (jwtService *jwtService) CreateToken(GuardName string, user models.User) (tokenData TokenOutPut, err error, token *jwt.Token) {
	var userRole []string
	for _, value := range user.Roles {
		userRole = append(userRole, value.Name)
	}
	token = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		CustomClaims{
			Roles:      userRole,       // 用户角色
			UserStatus: user.Status,    // 用户状态
			UserName:   user.Name,      // 用户名
			UserTel:    user.Telnumber, // 用户密码
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Unix() + global.App.Config.Jwt.JwtTtl,
				Id:        user.GetUid(),
				Issuer:    GuardName,
				NotBefore: time.Now().Unix() - 1000,
			},
		},
	)
	tokenStr, err := token.SignedString([]byte(global.App.Config.Jwt.Secret))

	tokenData = TokenOutPut{
		tokenStr,
		int(global.App.Config.Jwt.JwtTtl),
		TokenType,
	}

	return
}

func (jwtService *jwtService) getBlackListKey(tokenStr string) string {
	return "jwt_black_list:" + utils.MD5([]byte(tokenStr))
}

// 加入黑名单
func (jwtService *jwtService) JoinBlackList(token *jwt.Token) (err error) {
	// 如果参数为string,还需要增加一步解析
	now := time.Now().Unix()
	timer := time.Duration(token.Claims.(*CustomClaims).ExpiresAt-now) * time.Second
	err = global.App.Redis.SetNX(context.Background(), JwtService.getBlackListKey(token.Raw), now, timer).Err()
	return
}

// 判断是否在黑名单里
func (jwtService *jwtService) IsInBlacklist(tokenStr string) bool {
	joinUnixStr, err := global.App.Redis.Get(context.Background(), JwtService.getBlackListKey(tokenStr)).Result()
	if joinUnixStr == "" || err != nil {
		return false
	}
	joinUnix, err := strconv.ParseInt(joinUnixStr, 10, 64)
	if err != nil {
		return false
	}

	if time.Now().Unix()-joinUnix < global.App.Config.Jwt.JwtBlacklistGracePeriod {
		return false
	}
	return true

}
