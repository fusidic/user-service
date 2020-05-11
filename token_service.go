package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	pb "github.com/fusidic/user-service/proto/user"
)

var (
	// 此处声明一个安全密钥作为哈希
	// 此处仅作参考，实际使用时应该用随机生成的md5值或者其他
	key = []byte("fusidicsSuperSecretKey")
)

// CustomClaims 作为元数据，在被哈希之后作为第二段数据被发送给JWT
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

// Authable ...
type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *pb.User) (string, error)
}

// TokenService ...
type TokenService struct {
	repo Repository
}

// Decode 将token字符串解码为token对象
func (srv *TokenService) Decode(tokenString string) (*CustomClaims, error) {

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// Validate the token and return the custom claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// Encode 将claim编码为JWT
func (srv *TokenService) Encode(user *pb.User) (string, error) {

	expireToken := time.Now().Add(time.Hour * 72).Unix()

	// 创建Claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "user",
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 注册token并返回
	return token.SignedString(key)
}
