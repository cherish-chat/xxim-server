package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
)

// getUsersByIds
func getUsersByIds(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, ids []string) ([]*usermodel.User, error) {
	users, err := getUsersByIdsFromRedis(ctx, rc, ids)
	if err != nil {
		return userFromMongo(ctx, rc, c, ids)
	}
	// 判断是否有缺失
	userMap := make(map[string]*usermodel.User)
	for _, user := range users {
		userMap[user.Id] = user
	}
	var missIds []string
	for _, id := range ids {
		if _, ok := userMap[id]; !ok {
			missIds = append(missIds, id)
		}
	}
	if len(missIds) > 0 {
		missUsers, err := userFromMongo(ctx, rc, c, missIds)
		if err != nil {
			return nil, err
		}
		users = append(users, missUsers...)
	}
	return users, nil
}

func userFromMongo(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, ids []string) ([]*usermodel.User, error) {
	users := make([]*usermodel.User, 0)
	err := error(nil)
	xtrace.StartFuncSpan(ctx, "FindUserByIds", func(ctx context.Context) {
		err = c.Find(ctx, bson.M{
			"_id": bson.M{
				"$in": ids,
			},
		}).All(&users)
	})
	if err != nil {
		logx.WithContext(ctx).Errorf("find users by ids %v error: %s", ids, err.Error())
		return users, err
	}
	// 存入 redis
	userMap := make(map[string]*usermodel.User)
	for _, user := range users {
		userMap[user.Id] = user
		key := rediskey.UserKey(user.Id)
		err = rc.SetexCtx(ctx, key, utils.AnyToString(user), user.ExpireSeconds())
		if err != nil {
			logx.WithContext(ctx).Errorf("set user %s to redis error: %s", user.Id, err.Error())
			continue
		}
	}
	var notFoundIds []string
	for _, id := range ids {
		if _, ok := userMap[id]; !ok {
			notFoundIds = append(notFoundIds, id)
		}
	}
	if len(notFoundIds) > 0 {
		for _, id := range notFoundIds {
			key := rediskey.UserKey(id)
			err = rc.SetexCtx(ctx, key, xredis.NotFound, xredis.ExpireMinutes(5))
			if err != nil {
				logx.WithContext(ctx).Errorf("set user %s to redis error: %s", id, err.Error())
				continue
			}
		}
	}
	return users, nil
}

func getUsersByIdsFromRedis(ctx context.Context, rc *redis.Redis, ids []string) ([]*usermodel.User, error) {
	users := make([]*usermodel.User, 0)
	vals, err := rc.MgetCtx(ctx, utils.UpdateSlice(ids, func(id string) string {
		return rediskey.UserKey(id)
	})...)
	if err != nil {
		logx.WithContext(ctx).Errorf("get users by ids %v from redis error: %s", ids, err.Error())
		return users, err
	}
	for i, val := range vals {
		user := &usermodel.User{}
		if val == xredis.NotFound {
			id := ids[i]
			user.NotFound(id)
		} else {
			err = json.Unmarshal([]byte(val), user)
			if err != nil {
				logx.WithContext(ctx).Errorf("convert user error: %s", err.Error())
				continue
			}
		}
		users = append(users, user)
	}
	return users, nil
}

func flushUserCache(ctx context.Context, rc *redis.Redis, ids []string) error {
	var err error
	if len(ids) > 0 {
		xtrace.StartFuncSpan(ctx, "DeleteCache", func(ctx context.Context) {
			redisKeys := utils.UpdateSlice(ids, func(v string) string {
				return rediskey.UserKey(v)
			})
			_, err = rc.DelCtx(ctx, redisKeys...)
		})
	}
	return err
}
