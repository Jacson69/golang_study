package service

import (
	jwt2 "github.com/golang-jwt/jwt/v4"
	"time"
)

// secret key
var secretKey = []byte("abcd1234!@#$")

// ArithmeticCustomClaims 自定义声明
type ArithmeticCustomClaims struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`

	jwt2.StandardClaims
}

// jwtKeyFunc 返回密钥
func JwtKeyFunc(token *jwt2.Token) (interface{}, error) {
	return secretKey, nil
}

// Sign 生成token
func Sign(name, uid string) (string, error) {

	//为了演示方便，设置两分钟后过期
	expAt := time.Now().Add(time.Duration(2) * time.Minute).Unix()

	// 创建声明
	claims := ArithmeticCustomClaims{
		UserId: uid,
		Name:   name,
		StandardClaims: jwt2.StandardClaims{
			ExpiresAt: expAt,
			Issuer:    "system",
		},
	}

	//创建token，指定加密算法为HS256
	token := jwt2.NewWithClaims(jwt2.SigningMethodHS256, claims)

	//生成token
	return token.SignedString(secretKey)
}
