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
	// 清理缓存
	m.Flush(context.Background(), "")
	return m
}

func (m *ConfigMgr) Flush(ctx context.Context, userId string) error {
	_, err := m.rc.DelCtx(ctx, m.rediskey+":"+userId)
	return err
}

//func (m *ConfigMgr) GetAll(ctx context.Context) ([]*appmgmtmodel.Config, error) {
//	logger := logx.WithContext(ctx)
//	val, err := m.rc.GetCtx(ctx, m.rediskey)
//	if err != nil || val == "" {
//		return m.GetAllFromMysql(ctx, "")
//	}
//	if val == xredis.NotFound {
//		logger.Errorf("get config from redis not found")
//		return make([]*appmgmtmodel.Config, 0), nil
//	}
//	var configs []*appmgmtmodel.Config
//	err = json.Unmarshal([]byte(val), &configs)
//	if err != nil {
//		logger.Errorf("get config from redis unmarshal error: %v. val: %s", err, val)
//		return m.GetAllFromMysql(ctx, "")
//	}
//	return configs, nil
//}

func (m *ConfigMgr) GetAll(ctx context.Context, userId string) ([]*appmgmtmodel.Config, error) {
	logger := logx.WithContext(ctx)
	val, err := m.rc.GetCtx(ctx, m.rediskey+":"+userId)
	if err != nil || val == "" {
		return m.GetAllFromMysql(ctx, "")
	}
	if val == xredis.NotFound {
		logger.Errorf("get config from redis not found")
		return make([]*appmgmtmodel.Config, 0), nil
	}
	var configs []*appmgmtmodel.Config
	err = json.Unmarshal([]byte(val), &configs)
	if err != nil {
		logger.Errorf("get config from redis unmarshal error: %v. val: %s", err, val)
		return m.GetAllFromMysql(ctx, userId)
	}
	return configs, nil
}

func (m *ConfigMgr) GetCtx(ctx context.Context, k string, userId string) string {
	configs, err := m.GetAll(ctx, userId)
	if err != nil {
		return ""
	}
	for _, config := range configs {
		if (utils.InSlice(config.GetScopePlatforms(), m.platform) || config.ScopePlatforms == "") && config.K == k {
			return config.V
		}
	}
	configs, err = m.GetAll(ctx, "")
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

func (m *ConfigMgr) GetByPlatformCtx(ctx context.Context, k string, platform string, userId string) string {
	configs, err := m.GetAll(ctx, userId)
	if err != nil {
		return ""
	}
	for _, config := range configs {
		if (utils.InSlice(config.GetScopePlatforms(), platform) || config.ScopePlatforms == "") && config.K == k {
			return config.V
		}
	}
	configs, err = m.GetAll(ctx, "")
	if err != nil {
		return ""
	}
	for _, config := range configs {
		if (utils.InSlice(config.GetScopePlatforms(), platform) || config.ScopePlatforms == "") && config.K == k {
			return config.V
		}
	}
	return ""
}

func (m *ConfigMgr) GetOrDefaultCtx(ctx context.Context, k string, defaultValue string, userId string) string {
	configs, err := m.GetAll(ctx, userId)
	if err != nil {
		return defaultValue
	}
	for _, config := range configs {
		if (utils.InSlice(config.GetScopePlatforms(), m.platform) || config.ScopePlatforms == "") && config.K == k {
			return config.V
		}
	}
	configs, err = m.GetAll(ctx, "")
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

func (m *ConfigMgr) GetSliceCtx(ctx context.Context, k string, userId string) []string {
	configs, err := m.GetAll(ctx, userId)
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
	configs, err = m.GetAll(ctx, "")
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

func (m *ConfigMgr) MGetOrDefaultCtx(ctx context.Context, kv map[string]string, userId string) map[string]string {
	configs, err := m.GetAll(ctx, userId)
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
	if len(mp) == 0 {
		configs, err = m.GetAll(ctx, "")
		if err != nil {
			return kv
		}
		for _, config := range configs {
			mp[config.K] = config.V
		}
		for k, v := range kv {
			if _, ok := mp[k]; !ok {
				mp[k] = v
			}
		}
	}
	return mp
}

func (m *ConfigMgr) GetAllFromMysql(ctx context.Context, userId string) ([]*appmgmtmodel.Config, error) {
	logger := logx.WithContext(ctx)
	// 删除缓存
	err := m.Flush(ctx, userId)
	if err != nil {
		logger.Errorf("flush config error: %v", err)
		return nil, err
	}
	var configs []*appmgmtmodel.Config
	err = m.mysql.Where("userId = ?", userId).Find(&configs).Error
	if err != nil {
		logger.Errorf("get config from mysql error: %v", err)
		return nil, err
	}
	// 缓存
	{
		err := m.rc.SetCtx(ctx, m.rediskey+":"+userId, utils.AnyToString(configs))
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

func (m *ConfigMgr) upsert(k string, config *appmgmtmodel.Config) {
	// 查询是否存在
	var c appmgmtmodel.Config
	err := m.mysql.Where("k = ?", k).First(&c).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			m.mysql.Create(config)
		}
	} else {
		m.mysql.Model(&c).Where("k = ? AND userId = ''", k).Updates(map[string]any{
			"group":          config.Group,
			"scopePlatforms": config.ScopePlatforms,
			"type":           config.Type,
			"name":           config.Name,
			"options":        config.Options,
		})
	}
}

func (m *ConfigMgr) delete(id string) {
	m.mysql.Where("id = ?", id).Delete(&appmgmtmodel.Config{})
}
