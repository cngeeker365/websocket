package ucenter

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var(
	//定义HASH密码时用的盐，该部分目前为固定，后续应保证其生成和保存的安全性，如用MD5生成
	privateKey = []byte("`s#ax_1-!")
)

type CustomClaims struct{
	User *User
	//使用标准的 payload
	jwt.StandardClaims
}

type Auth interface{
	Decode(token string) (*CustomClaims, error)
	Encode(user *User) (string, error)
}

type TokenService struct {}

// 将 User 加密为 JWT 字符串
func (this *TokenService) Encode(user *User) (string, error){

	var (
		//设定过期时间
		expireTime int64
		claims CustomClaims
	)

	//1024天后过期
	expireTime = time.Now().Add(time.Hour * 24 * 1024).Unix()
	claims = CustomClaims{
		User:user,
		StandardClaims: jwt.StandardClaims{
			Issuer: "sminedata.com",	//签发者
			ExpiresAt: expireTime,
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(privateKey)
}

// 将 JWT 字符串解密为 CustomClaims 对象
func (this *TokenService) Decode(token string) (*CustomClaims, error)  {
	t, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return privateKey, nil
	})
	// 解密转换类型并返回
	if claims, ok := t.Claims.(*CustomClaims); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
