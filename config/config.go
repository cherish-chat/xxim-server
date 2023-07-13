package config

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	zeroredis "github.com/zeromicro/go-zero/core/stores/redis"
	"net"
	"strings"
	"time"
)

type Config struct {
	Etcd   EtcdConfig
	Log    LogConfig
	Mode   string `json:",default=pro,options=dev|test|rt|pre|pro"`
	Jaeger JaegerConfig

	Gateway struct {
		Mode            string   `json:",default=p2p,options=p2p"`
		SignalingServer string   `json:",optional"` // xxx.xxx.xxx:xx
		AppId           string   `json:",optional"`
		AppSecret       string   `json:",optional"`
		StunUrls        []string `json:",optional"`
	}

	Rsa struct {
		PrivateKey string
		PublicKey  string `json:",optional"`
	}

	Redis struct {
		Addrs    []string
		Username string `json:",optional"`
		Password string `json:",optional"`
	}

	MongoCollection struct {
		User        xmgo.MongoCollectionConf
		UserSetting xmgo.MongoCollectionConf

		Group               xmgo.MongoCollectionConf
		GroupMember         xmgo.MongoCollectionConf
		GroupSubscribeCache xmgo.MongoCollectionConf

		Friend            xmgo.MongoCollectionConf
		FriendApplyRecord xmgo.MongoCollectionConf

		Channel               xmgo.MongoCollectionConf
		ChannelMember         xmgo.MongoCollectionConf
		ChannelSubscribeCache xmgo.MongoCollectionConf

		Message xmgo.MongoCollectionConf
		Notice  xmgo.MongoCollectionConf
	}

	Account struct {
		Register struct {
			AllowPlatform []peerpb.Platform `json:",optional"`
			// 是否必填password
			RequirePassword bool `json:",optional"`
			// 是否必填nickname
			RequireNickname bool `json:",optional"`
			// 默认昵称规则
			DefaultNicknameRule string `json:",options=random|fixed"` // random:随机生成，fixed:固定昵称
			//FixedNickname 固定昵称
			FixedNickname string `json:",default=用户"` // 固定昵称 只有DefaultNicknameRule=fixed时有效
			//RandomNicknamePrefix 随机昵称前缀
			RandomNicknamePrefix string `json:",default=用户"` // 随机昵称前缀 只有DefaultNicknameRule=random时有效
			//UsernameRegex 用户名正则
			UsernameRegex string `json:",optional"`
			//UsernameUnique 用户名是否唯一
			UsernameUnique bool `json:",default=true"`
			// 是否必填avatar
			RequireAvatar bool `json:",optional"`
			//默认头像规则
			DefaultAvatarRule string `json:",options=byName|fixed"` // byName:根据昵称生成，fixed:固定头像
			//ByNameAvatarBgColors 根据昵称生成头像的背景颜色 ex: ["#ffffff","#000000"]
			ByNameAvatarBgColors []string
			//ByNameAvatarFgColors 根据昵称生成头像的字体颜色 ex: ["#ffffff","#000000"]
			ByNameAvatarFgColors []string
			//FixedAvatar 固定头像
			FixedAvatar string `json:",default=group_avatar.png"` // 固定头像 只有DefaultAvatarRule=fixed时有效
			// 是否必须绑定手机号
			RequireBindPhone bool `json:",optional"`
			//PhoneCountryCode
			PhoneRules []struct {
				CountryCode string
				PhoneRegex  string `json:",optional"`
			} `json:",optional"`
			//PhoneUnique 手机号是否唯一
			PhoneUnique bool `json:",default=true"`
			// 是否必须绑定邮箱
			RequireBindEmail bool `json:",optional"`
			//EmailRegex
			EmailRegex string `json:",optional"`
			//EmailUnique
			EmailUnique bool `json:",default=true"`
			// 是否验证图形验证码
			RequireCaptcha bool `json:",optional"`
		}
		Login struct {
			AllowPlatform []peerpb.Platform `json:",optional"`
			// RequireCaptcha 是否验证图形验证码
			RequireCaptcha bool `json:",optional"`
			//JwtConfig jwt配置
			JwtConfig utils.JwtConfig
		}
	}

	Friend struct {
		MaxFriendCount     int64    `json:",default=2000"`
		AllowRoleApply     []string `json:",optional"` // 允许角色申请好友 user|robot
		AllowRoleBeApplied []string `json:",optional"` // 允许角色被申请好友 user|robot

		DefaultSayHello string `json:",default=我们已经是好友了，来聊天吧！"` // 默认打招呼内容
	}

	Group struct {
		AllowRoleCreate []string `json:",optional"`     // 允许角色创建群组 user|robot
		JoinedMaxCount  int64    `json:",default=2000"` // 允许加入的最大群组数量

		RequiredName    bool   `json:",default=true"`           // 是否必填群组名称
		DefaultNameRule string `json:",options=byMember|fixed"` // byMember:根据成员生成，fixed:固定名称
		FixedName       string `json:",default=群组"`             // 固定名称 只有DefaultNameRule=fixed时有效

		RequiredAvatar       bool   `json:",default=true"`             // 是否必填群组头像
		DefaultAvatarRule    string `json:",options=byName|fixed"`     // byName:根据名称生成，fixed:固定头像
		FixedAvatar          string `json:",default=group_avatar.png"` // 固定头像 只有DefaultAvatarRule=fixed时有效
		ByNameAvatarBgColors []string
		ByNameAvatarFgColors []string

		DefaultWelcomeMessage string `json:",default=欢迎加入群组"` // 默认欢迎消息
	}
}

func (c Config) GetRedis(db int) redis.UniversalClient {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:       c.Redis.Addrs,
		DB:          db,
		Username:    c.Redis.Username,
		Password:    c.Redis.Password,
		DialTimeout: time.Second * 3,
	})
	pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := client.Ping(pingCtx).Err()
	if err != nil {
		logx.Errorf("redis ping error: %v", err)
		panic(err)
	}
	return client
}

func (c Config) GetPrivateKey() *rsa.PrivateKey {
	if c.Rsa.PrivateKey == "" {
		logx.Errorf("private key is empty")
		panic("private key is empty")
	}
	block, rest := pem.Decode([]byte(c.Rsa.PrivateKey))
	if len(rest) > 0 {
		logx.Errorf("private key is invalid")
		panic("private key is invalid")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logx.Errorf("parse private key failed: %v", err)
		panic(err)
	}
	return privateKey
}

func (c Config) GetPublicKey() *rsa.PublicKey {
	if c.Rsa.PublicKey == "" {
		return nil
	}
	block, rest := pem.Decode([]byte(c.Rsa.PublicKey))
	if len(rest) > 0 {
		logx.Errorf("public key is invalid")
		panic("public key is invalid")
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		logx.Errorf("parse public key failed: %v", err)
		panic(err)
	}
	return publicKey
}

func (c Config) GetZeroRedisConf() zeroredis.RedisConf {
	typ := "node"
	if len(c.Redis.Addrs) > 1 {
		typ = "cluster"
	}
	return zeroredis.RedisConf{
		Host:     strings.Join(c.Redis.Addrs, ","),
		Type:     typ,
		Pass:     c.Redis.Password,
		NonBlock: true,
	}
}

func RpcPort() string {
	// 获取空闲端口
	var a *net.TCPAddr
	var err error
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return fmt.Sprintf("0.0.0.0:%d", l.Addr().(*net.TCPAddr).Port)
		}
	}
	logx.Errorf("get rpc port failed: %v", err)
	panic(err)
}
