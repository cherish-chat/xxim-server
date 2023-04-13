package xjwt

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	{
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Id:        "aaa",
			ExpiresAt: time.Now().Add(time.Second * 3600).Unix(),
		})
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		t.Logf("tokenString: %v", tokenString)

		// 验证
		token, err = jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			t.Logf("claims: %+v", claims)
			err := claims.Valid()
			if err != nil {
				t.Fatalf("err: %v", err)
			} else {
				t.Logf("valid")
			}
		} else {
			t.Fatalf("token invalid")
		}
	}
	{
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Id:        "aaa",
			ExpiresAt: time.Now().Add(time.Second * 3600).UnixNano(),
		})
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		t.Logf("tokenString: %v", tokenString)
	}
}
