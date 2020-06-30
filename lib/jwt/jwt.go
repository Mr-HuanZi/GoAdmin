package jwt

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type LoginToken struct {
	Exp      int64
	Iat      int64
	Username string
	User     int64
}

//生成用户令牌
func GenerateUserToken(username string, uid int64) (string, error) {
	expireSecondsConf := beego.AppConfig.String("jwt::expireSeconds")
	var expireSeconds int
	if s, err := strconv.Atoi(expireSecondsConf); err == nil {
		expireSeconds = s
	} else {
		expireSeconds = 600 //默认10分钟
	}
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Second * time.Duration(expireSeconds)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["username"] = username
	claims["User"] = uid
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(beego.AppConfig.String("jwt::jwtKey")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//验证用户令牌
func ValidateUserToken(tokenString string) (bool, *jwt.Token) {
	if tokenString == "" {
		return false, nil
	}
	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(beego.AppConfig.String("jwt::jwtKey")), nil
	})
	logs.Info("%+v", t)
	if err != nil {
		logs.Error("转换为jwt claims失败.", err)
		return false, nil
	}
	return true, t
}

//刷新用户令牌
func RefreshUserToken(t *jwt.Token) (string, error) {
	LoginToken := GetTokenClaims(t)
	if LoginToken.Exp-LoginToken.Iat <= 300 {
		// 如果令牌有效期小于5分钟
		return GenerateUserToken(LoginToken.Username, LoginToken.User)
	}
	return "", nil
}

func GetTokenClaims(t *jwt.Token) LoginToken {
	claims := t.Claims.(jwt.MapClaims)
	var LoginToken LoginToken
	LoginToken.Exp = int64(claims["exp"].(float64))
	LoginToken.Username = claims["username"].(string)
	LoginToken.User = int64(claims["User"].(float64))
	LoginToken.Iat = time.Now().Unix()
	return LoginToken
}
