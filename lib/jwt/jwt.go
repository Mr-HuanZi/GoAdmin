package jwt

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

var (
	// Jwt加密Key，请勿泄露
	jwtKey = "ClNNHdJLif0XrhxD"
)

type MyCustomClaims struct {
	User int64
	jwt.StandardClaims
}

//生成用户令牌
func GenerateUserToken(uid int64) (string, error) {
	expireSecondsConf := beego.AppConfig.String("jwt::expireSeconds")
	var expireSeconds int
	if s, err := strconv.Atoi(expireSecondsConf); err == nil {
		expireSeconds = s
	} else {
		expireSeconds = 600 //默认10分钟
	}
	newClaims := MyCustomClaims{
		User: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(expireSeconds)).Unix(), // 过期时间
			IssuedAt:  time.Now().Unix(),                                                 // 签发时间
			Issuer:    "GoAdmin",                                                         // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//验证用户令牌
func ValidateUserToken(tokenString string) (bool, *MyCustomClaims) {
	if tokenString == "" {
		return false, nil
	}
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(jwtKey), nil
	})
	logs.Debug("%+v", token)
	if err != nil {
		logs.Error("转换为jwt claims失败.", err)
		return false, nil
	}
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		// 解析成功
		logs.Debug("%v %v", claims.User, claims.StandardClaims.ExpiresAt)
		logs.Debug("token will be expired at ", time.Unix(claims.StandardClaims.ExpiresAt, 0))
		return true, claims
	} else {
		logs.Warning("validate tokenString failed !!!", err)
		return false, nil
	}
}

//刷新用户令牌
func RefreshUserToken(claims *MyCustomClaims) (string, error) {
	if claims.StandardClaims.ExpiresAt-claims.StandardClaims.IssuedAt <= 300 {
		// 如果令牌有效期小于5分钟
		return GenerateUserToken(claims.User)
	}
	return "", nil
}
