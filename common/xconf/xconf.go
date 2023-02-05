package xconf

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ConfigMgr struct {
	mysql    *gorm.DB
	rc       *redis.Redis
	rediskey string
	platform string
}

func NewConfigMgr(tx *gorm.DB, rc *redis.Redis, platform string) *ConfigMgr {
	m := &ConfigMgr{
		mysql:    tx,
		rc:       rc,
		rediskey: rediskey.AllAppMgmtConfig(),
		platform: platform,
	}
	m.initData()
	return m
}

func (m *ConfigMgr) Flush(ctx context.Context) error {
	_, err := m.rc.DelCtx(ctx, m.rediskey)
	return err
}

func (m *ConfigMgr) GetAll(ctx context.Context) ([]*appmgmtmodel.Config, error) {
	logger := logx.WithContext(ctx)
	val, err := m.rc.GetCtx(ctx, m.rediskey)
	if err != nil || val == "" {
		return m.getAllFromMysql(ctx)
	}
	if val == xredis.NotFound {
		logger.Errorf("get config from redis not found")
		return make([]*appmgmtmodel.Config, 0), nil
	}
	var configs []*appmgmtmodel.Config
	err = json.Unmarshal([]byte(val), &configs)
	if err != nil {
		logger.Errorf("get config from redis unmarshal error: %v. val: %s", err, val)
		return m.getAllFromMysql(ctx)
	}
	return configs, nil
}

func (m *ConfigMgr) GetCtx(ctx context.Context, k string) string {
	configs, err := m.GetAll(ctx)
	if err != nil {
		return ""
	}
	for _, config := range configs {
		if (utils.InSlice(config.GetScopePlatforms(), m.platform) || config.ScopePlatforms == "") && config.K == k {
			return config.V
		}
	}
	return ""
}

func (m *ConfigMgr) GetOrDefaultCtx(ctx context.Context, k string, defaultValue string) string {
	configs, err := m.GetAll(ctx)
	if err != nil {
		return defaultValue
	}
	for _, config := range configs {
		if (utils.InSlice(config.GetScopePlatforms(), m.platform) || config.ScopePlatforms == "") && config.K == k {
			return config.V
		}
	}
	// 默认值
	return defaultValue
}

func (m *ConfigMgr) GetSliceCtx(ctx context.Context, k string) []string {
	configs, err := m.GetAll(ctx)
	if err != nil {
		return []string{}
	}
	for _, config := range configs {
		if (utils.InSlice(config.GetScopePlatforms(), m.platform) || config.ScopePlatforms == "") && config.K == k {
			var res []string
			err := json.Unmarshal([]byte(config.V), &res)
			if err != nil {
				return []string{}
			}
			return res
		}
	}
	// 默认值
	return []string{}
}

func (m *ConfigMgr) MGetOrDefaultCtx(ctx context.Context, kv map[string]string) map[string]string {
	configs, err := m.GetAll(ctx)
	if err != nil {
		return kv
	}
	mp := make(map[string]string)
	for _, config := range configs {
		mp[config.K] = config.V
	}
	for k, v := range kv {
		if _, ok := mp[k]; !ok {
			mp[k] = v
		}
	}
	return mp
}

func (m *ConfigMgr) getAllFromMysql(ctx context.Context) ([]*appmgmtmodel.Config, error) {
	logger := logx.WithContext(ctx)
	// 删除缓存
	err := m.Flush(ctx)
	if err != nil {
		logger.Errorf("flush config error: %v", err)
		return nil, err
	}
	var configs []*appmgmtmodel.Config
	err = m.mysql.Find(&configs).Error
	if err != nil {
		logger.Errorf("get config from mysql error: %v", err)
		return nil, err
	}
	// 缓存
	{
		err := m.rc.SetCtx(ctx, m.rediskey, utils.AnyToString(configs))
		if err != nil {
			logger.Errorf("set config to redis error: %v", err)
		}
	}
	return configs, nil
}

func (m *ConfigMgr) insertIfNotFound(k string, config *appmgmtmodel.Config) {
	var c appmgmtmodel.Config
	err := m.mysql.Where("k = ?", k).First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			m.mysql.Create(config)
		}
	}
}
