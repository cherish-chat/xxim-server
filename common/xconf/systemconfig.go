package xconf

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SystemConfig struct {
	Namespace string `bson:"namespace" json:"namespace" gorm:"column:namespace;type:varchar(255);not null;index:idx_namespace_serviceName_key,unique;comment:命名空间"`
	// 服务名称
	ServiceName string `bson:"serviceName" json:"serviceName" gorm:"column:serviceName;type:varchar(255);not null;index:idx_namespace_serviceName_key,unique;comment:服务名称"`
	// 键
	Key string `bson:"key" json:"key" gorm:"column:key;type:varchar(255);not null;index:idx_namespace_serviceName_key,unique;comment:键"`
	// 值
	Value string `bson:"value" json:"value" gorm:"column:value;type:varchar(255);not null;comment:值"`
	// 输入选项
	InputOptions xorm.M `bson:"inputOptions,omitempty" json:"inputOptions,omitempty" gorm:"column:inputOptions;type:json;comment:输入选项"`
}

func (m *SystemConfig) TableName() string {
	return "system_config"
}

type SystemConfigMgr struct {
	Namespace   string
	ServiceName string
	mysql       *gorm.DB
}

func NewSystemConfigMgr(namespace string, serviceName string, mysql *gorm.DB) *SystemConfigMgr {
	s := &SystemConfigMgr{Namespace: namespace, ServiceName: serviceName, mysql: mysql}
	s.initData()
	return s
}

func (s *SystemConfigMgr) Get(key string) (value string) {
	return s.GetCtx(context.Background(), key)
}

func (s *SystemConfigMgr) GetCtx(ctx context.Context, key string) (value string) {
	config := &SystemConfig{
		Namespace:   s.Namespace,
		ServiceName: s.ServiceName,
		Key:         key,
		Value:       "",
	}
	err := s.mysql.WithContext(ctx).Model(config).Where("namespace = ? and serviceName = ? and `key` = ?", s.Namespace, s.ServiceName, key).FirstOrCreate(&config).Error
	if err != nil {
		logx.WithContext(ctx).Errorf("获取配置失败: %v", err)
	}
	return config.Value
}

func (s *SystemConfigMgr) GetOrDefault(key string, defaultValue string) (value string) {
	ctx := context.Background()
	return s.GetOrDefaultCtx(ctx, key, defaultValue)
}

func (s *SystemConfigMgr) GetOrDefaultCtx(ctx context.Context, key string, defaultValue string) (value string) {
	config := &SystemConfig{
		Namespace:   s.Namespace,
		ServiceName: s.ServiceName,
		Key:         key,
		Value:       defaultValue,
	}
	err := s.mysql.WithContext(ctx).Model(config).Where("namespace = ? and serviceName = ? and `key` = ?", s.Namespace, s.ServiceName, key).FirstOrCreate(&config).Error
	if err != nil {
		logx.WithContext(ctx).Errorf("获取配置失败: %v", err)
	}
	return config.Value
}

func (s *SystemConfigMgr) GetSlice(key string) []string {
	return s.GetSliceCtx(context.Background(), key)
}

func (s *SystemConfigMgr) GetSliceCtx(ctx context.Context, key string) []string {
	config := &SystemConfig{
		Namespace:   s.Namespace,
		ServiceName: s.ServiceName,
		Key:         key,
		Value:       "[]",
	}
	err := s.mysql.WithContext(ctx).Model(config).Where("namespace = ? and serviceName = ? and `key` = ?", s.Namespace, s.ServiceName, key).FirstOrCreate(&config).Error
	if err != nil {
		logx.WithContext(ctx).Errorf("获取配置失败: %v", err)
	}
	var value []string
	err = json.Unmarshal([]byte(config.Value), &value)
	if err != nil {
		logx.WithContext(ctx).Errorf("获取配置失败: %v", err)
	}
	return value
}

func (s *SystemConfigMgr) initData() {
	configs := []*SystemConfig{
		{
			Namespace:    "system",
			ServiceName:  "user",
			Key:          "signature_if_not_set", // 未设置签名时的默认签名
			Value:        "这个人很懒，还没有设置签名哦～",
			InputOptions: nil,
		},
		{
			Namespace:    "system",
			ServiceName:  "user",
			Key:          "nickname_default", // 默认昵称
			Value:        "XX用户",
			InputOptions: nil,
		},
		{
			Namespace:    "system",
			ServiceName:  "user",
			Key:          "avatars_default", // 默认头像
			Value:        `["https://go-zero.dev/img/footer/go-zero.svg"]`,
			InputOptions: nil,
		},
		{
			Namespace:    "system",
			ServiceName:  "relation",
			Key:          "app.friend_max_count", // 好友最大数量
			Value:        `20000`,
			InputOptions: nil,
		},
		{
			Namespace:   "system",
			ServiceName: "user",
			Key:         "app.register_max_count_per_day_ip", // 每个IP每天最大注册数量
			Value:       `10`,
		},
	}
	for _, config := range configs {
		err := xorm.InsertOne(s.mysql, config)
		if err != nil {
			logx.Errorf("初始化配置失败: %v", err)
		}
	}
}
