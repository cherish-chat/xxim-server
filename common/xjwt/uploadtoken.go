package xjwt

import (
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type uploadTokenX struct {
}

var UploadToken = uploadTokenX{}

func (uploadTokenX) GenerateToken(
	objectId string,
	expireSeconds int32,
	secret string,
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        objectId,
		ExpiresAt: time.Now().Add(time.Second * time.Duration(expireSeconds)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (uploadTokenX) VerifyToken(
	tokenString string,
	secret string,
) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		err := claims.Valid()
		if err != nil {
			return "", err
		} else {
			return claims.Id, nil
		}
	} else {
		return "", err
	}
}
