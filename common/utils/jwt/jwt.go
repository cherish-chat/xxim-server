package jwt

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/cherish-chat/xxim-server/common/xenv"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type (
	VerifyTokenCode int
	exTyp           int
	tokenInfo       struct {
		jwt.StandardClaims
	}
	tokenEx struct {
		Typ          exTyp
		Token        string
		Base64PubKey string
		Data         string
		Platform     string `json:"platform"`
		UnbanTime    string `json:"unbanTime"`
	}
)

var (
	rsaKey             *rsa.PrivateKey
	generateRsaKeyTime time.Time
)

const (
	VerifyTokenCodeOK VerifyTokenCode = iota // 4
	VerifyTokenCodeError
	VerifyTokenCodeExpire // 2
	VerifyTokenCodeBaned  // 1
	VerifyTokenCodePlatformError
	VerifyTokenCodeReplace // 3

	exTypNormal exTyp = 0 // 普通
	exTypBaned  exTyp = 1 // 用户被禁用
)

func VerifyToken(ctx context.Context, rc redis.UniversalClient, token string, platform string) (VerifyTokenCode, string) {
	logger := logx.WithContext(ctx)
	if strings.HasPrefix(token, "token.uid.") {
		if !xenv.EnableDebugToken {
			return VerifyTokenCodeError, ""
		}
		return VerifyTokenCodeOK, strings.TrimPrefix(token, "token.uid.")
	}
	// 生产环境
	claimsString := tokenClaimsString(token)
	key := rediskey.AuthToken(claimsString)
	ret, err := rc.HGet(ctx, key, platform).Result()
	if err != nil {
		return VerifyTokenCodeExpire, ""
	}
	if ret == "" {
		logger.Errorf("token not found: %s", token)
		return VerifyTokenCodeError, ""
	}
	ex := tokenEx{}
	err = json.Unmarshal([]byte(ret), &ex)
	if err != nil {
		logger.Errorf("token unmarshal error: %s", err.Error())
		return VerifyTokenCodeError, ""
	}
	switch ex.Typ {
	case exTypBaned:
		unbanTime, _ := time.Parse(time.RFC3339, ex.UnbanTime)
		if time.Now().Before(unbanTime) {
			return VerifyTokenCodeBaned, ex.Data
		}
	}
	if ex.Token != token {
		logger.Errorf("token not match: %s", token)
		return VerifyTokenCodeReplace, ""
	}
	tf := &tokenInfo{}
	toPubKey, err := bytesToPubKey(ex.Base64PubKey)
	if err != nil {
		logger.Errorf("token pub key error: %s", err.Error())
		return VerifyTokenCodeError, ""
	}
	getToken, err := jwt.ParseWithClaims(ex.Token, tf, func(token *jwt.Token) (interface{}, error) {
		return toPubKey, nil
	})
	if err != nil {
		logger.Errorf("token parse error: %s", err.Error())
		return VerifyTokenCodeError, ""
	}
	if getToken == nil || !getToken.Valid {
		logger.Errorf("token parse error: %s", err.Error())
		return VerifyTokenCodeError, ""
	}
	if ti, ok := getToken.Claims.(*tokenInfo); !ok {
		logger.Errorf("token parse error: %s", err.Error())
		return VerifyTokenCodeError, ""
	} else {
		if ex.Platform != platform {
			logger.Errorf("token platform error: ti.Platform: %s, platform: %s", ex.Platform, platform)
			return VerifyTokenCodePlatformError, ti.Id
		}
		return VerifyTokenCodeOK, ti.Id
	}
}

func BanToken(ctx context.Context, rc redis.UniversalClient, uid string, msg string, unbanTime time.Time) VerifyTokenCode {
	token := GenerateTokenButNotSet(uid)
	logger := logx.WithContext(ctx)
	claimsString := tokenClaimsString(token)
	keyr := rediskey.AuthToken(claimsString)
	{
		kvs, err := rc.HGetAll(ctx, keyr).Result()
		if err != nil {
			logger.Errorf("redis scan error: %s", err.Error())
			return VerifyTokenCodeError
		}
		for key, ret := range kvs {
			if ret == "" {
				logger.Errorf("token not found: %s", token)
				continue
			}
			ex := tokenEx{}
			err = json.Unmarshal([]byte(ret), &ex)
			if err != nil {
				logger.Errorf("token unmarshal error: %s", err.Error())
				return VerifyTokenCodeError
			}
			newTokenEx := tokenEx{
				Typ:          exTypBaned,
				Token:        ex.Token,
				Base64PubKey: ex.Base64PubKey,
				Data:         msg,
				Platform:     ex.Platform,
				UnbanTime:    unbanTime.Format(time.RFC3339),
			}
			buf, _ := json.Marshal(newTokenEx)
			err = rc.HSet(ctx, keyr, key, string(buf)).Err()
			if err != nil {
				logger.Errorf("token ban error: %s", err.Error())
				return VerifyTokenCodeError
			}
		}
	}
	return VerifyTokenCodeBaned
}

func SetToken(ctx context.Context, rc redis.UniversalClient, token string, platform string) error {
	if code, tip := VerifyToken(ctx, rc, token, platform); code == VerifyTokenCodeBaned {
		return errors.New(tip)
	}
	claimsString := tokenClaimsString(token)
	key := rediskey.AuthToken(claimsString)
	buf, _ := json.Marshal(tokenEx{Typ: exTypNormal, Token: token, Base64PubKey: pubKeyToBytes(&rsaKey.PublicKey), Platform: platform})
	err := rc.HSet(ctx, key, platform, string(buf)).Err()
	if err != nil {
		return err
	}
	return nil
}

func GenerateTokenButNotSet(id string) string {
	now := time.Now()
	if rsaKey == nil || now.After(generateRsaKeyTime.Add(time.Minute)) {
		rsaKey = rsaKeyGenToKey(1024)
		generateRsaKeyTime = now
	}
	ti := tokenInfo{}
	ti.Id = id
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodRS256, ti)
	token, _ := tokenObj.SignedString(rsaKey)
	return token
}

func tokenClaimsString(token string) string {
	split := strings.Split(token, ".")
	if len(split) < 2 {
		return ""
	}
	return split[1]
}

func rsaKeyGenToKey(bits int) *rsa.PrivateKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, bits)
	return privateKey
}

func pubKeyToBytes(key *rsa.PublicKey) string {
	return base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(key))
}

func bytesToPubKey(base64Ket string) (interface{}, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Ket)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PublicKey(keyBytes)
}
