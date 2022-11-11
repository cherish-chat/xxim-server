package xjwt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type (
	VerifyTokenCode int
	exTyp           int
	tokenInfo       struct {
		jwt.StandardClaims
		UniqueSuffix string `json:"uniqueSuffix"`
	}
	TokenObj struct {
		UserId       string `json:"userId"`
		UniqueSuffix string `json:"uniqueSuffix"`
		Typ          exTyp  `json:"typ"`
		Token        string `json:"token"`
		Base64PubKey string `json:"base64PubKey"`
		Data         string `json:"data"`
		Platform     string `json:"platform"`
		DeviceId     string `json:"serviceId"`
		DeviceModel  string `json:"deviceModel"`
	}
)

const (
	VerifyTokenCodeOK            VerifyTokenCode = iota
	VerifyTokenCodeInternalError                 // 服务器内部错误
	VerifyTokenCodeError                         // 验证失败
	VerifyTokenCodeExpire                        // 过期
	VerifyTokenCodeBaned                         // 被封禁
	VerifyTokenCodeReplace                       // 被顶替

	exTypNormal exTyp = 0 // 普通
	exTypBaned  exTyp = 1 // 用户被禁用
)

func GenerateToken(
	userId string,
	uniqueSuffix string,
	opts ...OptionFunc,
) *TokenObj {
	option := &Option{}
	for _, o := range opts {
		o(option)
	}
	rsaKey := rsaKeyGenToKey(2048)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tokenInfo{
		StandardClaims: jwt.StandardClaims{
			Id: userId,
		},
		UniqueSuffix: uniqueSuffix,
	})
	tokenString, _ := token.SignedString(rsaKey)
	t := &TokenObj{
		UserId:       userId,
		UniqueSuffix: uniqueSuffix,
		Typ:          exTypNormal,
		Token:        tokenString,
		Base64PubKey: pubKeyToBytes(&rsaKey.PublicKey),
	}
	if option.Platform != nil {
		t.Platform = *option.Platform
	}
	if option.DeviceId != nil {
		t.DeviceId = *option.DeviceId
	}
	if option.DeviceModel != nil {
		t.DeviceModel = *option.DeviceModel
	}
	if option.Data != nil {
		t.Data = *option.Data
	}
	return t
}

func SaveToken(
	ctx context.Context,
	rc *redis.Redis,
	token *TokenObj,
) error {
	key := rediskey.UserToken(token.UserId)
	hkey := tokenClaimsString(token.Token)
	err := rc.HsetCtx(ctx, key, hkey, utils.AnyToString(token))
	if err != nil {
		return err
	}
	return nil
}

func VerifyToken(
	ctx context.Context,
	rc *redis.Redis,
	userId string,
	tokenStr string,
	opts ...OptionFunc,
) (VerifyTokenCode, string) {
	option := &Option{}
	for _, o := range opts {
		o(option)
	}
	logger := logx.WithContext(ctx)
	key := rediskey.UserToken(userId)
	hkey := tokenClaimsString(tokenStr)
	tokenObj := &TokenObj{}
	val, err := rc.HgetCtx(ctx, key, hkey)
	if err != nil {
		logger.Errorf("redis Hget error: %v, key: %s, hkey: %s", err, key, hkey)
		return VerifyTokenCodeInternalError, "服务器内部错误"
	}
	if val == "" {
		// token不存在 过期
		return VerifyTokenCodeExpire, "登录已过期"
	}
	err = json.Unmarshal([]byte(val), tokenObj)
	if err != nil {
		logger.Errorf("json Unmarshal error: %v", err)
		return VerifyTokenCodeInternalError, "服务器内部错误"
	}
	if tokenObj.Typ == exTypBaned {
		return VerifyTokenCodeBaned, tokenObj.Data
	}
	if option.Platform != nil && tokenObj.Platform != *option.Platform {
		return VerifyTokenCodeError, "登录平台不匹配"
	}
	if tokenObj.Token != tokenStr {
		return VerifyTokenCodeReplace, fmt.Sprintf("您的账号在%s上登录", tokenObj.DeviceModel)
	}
	if option.DeviceId != nil && tokenObj.DeviceId != *option.DeviceId {
		return VerifyTokenCodeReplace, fmt.Sprintf("您的账号在%s上登录", tokenObj.DeviceModel)
	}
	// jwt验证
	token, err := jwt.ParseWithClaims(tokenStr, &tokenInfo{}, func(token *jwt.Token) (interface{}, error) {
		return bytesToPubKey(tokenObj.Base64PubKey)
	})
	if err != nil {
		logger.Errorf("jwt ParseWithClaims error: %v", err)
		return VerifyTokenCodeError, "登录已过期"
	}
	if claims, ok := token.Claims.(*tokenInfo); ok && token.Valid {
		if claims.Id != userId {
			return VerifyTokenCodeError, "登录已过期"
		}
		if claims.UniqueSuffix != tokenObj.UniqueSuffix {
			return VerifyTokenCodeReplace, fmt.Sprintf("您的账号在%s上登录", tokenObj.DeviceModel)
		}
		return VerifyTokenCodeOK, ""
	}
	return VerifyTokenCodeOK, ""
}

// BanToken 封禁userId的所有token
func BanToken(
	ctx context.Context,
	rc *redis.Redis,
	userId string,
	data string,
) error {
	key := rediskey.UserToken(userId)
	vals, err := rc.HgetallCtx(ctx, key)
	if err != nil {
		return err
	}
	for _, val := range vals {
		tokenObj := &TokenObj{}
		err = json.Unmarshal([]byte(val), tokenObj)
		if err != nil {
			return err
		}
		tokenObj.Typ = exTypBaned
		tokenObj.Data = data
		err = SaveToken(ctx, rc, tokenObj)
		if err != nil {
			return err
		}
	}
	return nil
}
