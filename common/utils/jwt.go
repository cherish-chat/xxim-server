package utils

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"math"
	"time"
)

type JwtConfig struct {
	PrivateKey string
	ExpireHour int    `json:",default=24"`
	MaxToken   int    `json:",default=5"`    // 每一个key 最大token数量
	Scene      string `json:",default=user"` // user,admin,other
}

type Jwt struct {
	rc     *redis.Redis
	Config JwtConfig
}

func NewJwt(
	config JwtConfig,
	rc *redis.Redis,
) *Jwt {
	return &Jwt{
		rc:     rc,
		Config: config,
	}
}

type TokenObject struct {
	UserId    string   `json:"userId"`    // 用户id
	UniqueKey string   `json:"uniqueKey"` // 唯一标识 比如设备id/平台id 就可以实现单点登录/单平台登录/任意登录
	Status    int      `json:"status"`    // 状态 用户自行定义
	Token     string   `json:"token"`     // token
	Extra     string   `json:"extra"`     // 额外信息 用户自行定义
	AliveTime int64    `json:"aliveTime"` // 存活时间
	ExpiredAt int64    `json:"expiredAt"` // 过期时间
	Scope     []string `json:"scope"`     // 权限
	Scene     string   `json:"scene"`     // 场景
}

func (o *TokenObject) Marshal() string {
	return Json.MarshalToString(o)
}

func (x *Jwt) GenerateToken(
	userId string,
	uniqueKey string,
	status int,
	extra string,
	scope []string,
) *TokenObject {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "",
		Subject:   "",
		Audience:  nil,
		ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(100, 0, 0)),
		NotBefore: nil,
		IssuedAt:  nil,
		ID:        userId,
	})
	// password to sign
	tokenString, err := token.SignedString([]byte(Md5(x.Config.PrivateKey)))
	if err != nil {
		return nil
	}
	t := &TokenObject{
		UserId:    userId,
		UniqueKey: uniqueKey,
		Status:    status,
		Token:     tokenString,
		Extra:     extra,
		AliveTime: time.Now().UnixMilli(),
		ExpiredAt: time.Now().Add(time.Hour * time.Duration(x.Config.ExpireHour)).UnixMilli(),
		Scope:     scope,
		Scene:     x.Config.Scene,
	}
	return t
}

var SetTokenLockError = fmt.Errorf("set token lock error")
var TokenInvalidError = fmt.Errorf("token invalid error")
var TokenReplaceError = fmt.Errorf("token replace error")
var TokenExpiredError = fmt.Errorf("token expired error")

func (x *Jwt) SetToken(
	ctx context.Context,
	tokenObject *TokenObject,
) error {
	key := fmt.Sprintf("token:%s:%s", x.Config.Scene, tokenObject.UserId) // redis key
	// key 上锁
	ok, err := x.rc.SetnxExCtx(ctx, "lock:"+key, "1", 5)
	if err != nil {
		return err
	}
	if !ok {
		// 同一时间只能有一个协程操作
		return SetTokenLockError
	}
	// 解锁
	defer x.rc.DelCtx(ctx, "lock:"+key)
	hkey := tokenObject.UniqueKey // redis hkey
	// 获取所有的token hgetall
	hgetall, err := x.rc.HgetallCtx(ctx, key)
	if err != nil {
		if err != redis.Nil {
			return err
		}
	}
	// 是否到达token上限
	if len(hgetall) >= x.Config.MaxToken {
		// 删 len(hgetall)+1-x.Config.MaxToken 个
		for i := 0; i < len(hgetall)+1-x.Config.MaxToken; i++ {
			// 删除最早的token
			// hgetall 排序
			minAliveTime := int64(math.MaxInt64)
			minAliveTimeHkey := ""
			for k, v := range hgetall {
				to := &TokenObject{}
				err := Json.Unmarshal([]byte(v), to)
				if err != nil {
					// 删除
					x.rc.HdelCtx(ctx, key, k)
					continue
				}
				if to.AliveTime < minAliveTime {
					minAliveTime = to.AliveTime
					minAliveTimeHkey = k
				}
			}
			// 删除
			x.rc.HdelCtx(ctx, key, minAliveTimeHkey)
			delete(hgetall, minAliveTimeHkey)
		}
	}
	hvalue := tokenObject.Marshal() // tokenObject序列化
	err = x.rc.HsetCtx(ctx, key, hkey, hvalue)
	if err != nil {
		return err
	}
	return nil
}

func (x *Jwt) VerifyToken(
	ctx context.Context,
	tokenString string,
	uniqueKey string,
) (*TokenObject, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(Md5(x.Config.PrivateKey)), nil
	})
	if err != nil {
		return nil, TokenInvalidError
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		err := claims.Valid()
		if err != nil {
			return nil, TokenInvalidError
		} else {
			userId := claims.ID
			key := fmt.Sprintf("token:%s:%s", x.Config.Scene, userId) // redis key
			hkey := uniqueKey                                         // redis hkey
			hget, err := x.rc.HgetCtx(ctx, key, hkey)
			if err != nil {
				return nil, TokenInvalidError
			}
			to := &TokenObject{}
			err = Json.Unmarshal([]byte(hget), to)
			if err != nil {
				return to, TokenInvalidError
			}
			if to.Token != tokenString {
				return to, TokenReplaceError
			}
			if to.ExpiredAt < time.Now().UnixMilli() {
				return to, TokenExpiredError
			}
			return to, nil
		}
	} else {
		return nil, TokenInvalidError
	}
}

func (x *Jwt) RevokeToken(
	ctx context.Context,
	userId string,
	uniqueKey string,
) error {
	key := fmt.Sprintf("token:%s:%s", x.Config.Scene, userId) // redis key
	hkey := uniqueKey                                         // redis hkey
	_, err := x.rc.HdelCtx(ctx, key, hkey)
	if err != nil {
		return err
	}
	return nil
}
