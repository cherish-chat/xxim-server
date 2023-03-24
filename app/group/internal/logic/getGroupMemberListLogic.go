package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupMemberList 获取群成员列表
func (l *GetGroupMemberListLogic) GetGroupMemberList(in *pb.GetGroupMemberListReq) (*pb.GetGroupMemberListResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{
			Page: 1,
			Size: 100,
		}
	}
	whereMap := map[string]interface{}{
		"groupId": in.GroupId,
	}
	whereOrMap := map[string][]interface{}{}
	offset := (in.Page.Page - 1) * in.Page.Size
	limit := in.Page.Size
	if in.Filter != nil {
		if in.Filter.NoDisturb != nil {
			// 是否免打扰
			whereMap["noDisturb"] = *in.Filter.NoDisturb
		}
		if in.Filter.OnlyOwner != nil {
			// 是否只获取群主
			if _, ok := whereOrMap["role"]; !ok {
				whereOrMap["role"] = []interface{}{}
			}
			if *in.Filter.OnlyOwner {
				whereOrMap["role"] = append(whereOrMap["role"], groupmodel.RoleType_OWNER)
			}
		}
		if in.Filter.OnlyAdmin != nil {
			// 是否只获取管理员
			if _, ok := whereOrMap["role"]; !ok {
				whereOrMap["role"] = []interface{}{}
			}
			if *in.Filter.OnlyAdmin {
				whereOrMap["role"] = append(whereOrMap["role"], groupmodel.RoleType_MANAGER)
			}
		}
		if in.Filter.OnlyMember != nil {
			// 是否只获取成员
			if _, ok := whereOrMap["role"]; !ok {
				whereOrMap["role"] = []interface{}{}
			}
			if *in.Filter.OnlyMember {
				whereOrMap["role"] = append(whereOrMap["role"], groupmodel.RoleType_MEMBER)
			}
		}
	}
	rdsKey := rediskey.GroupMemberSearchKey(in.GroupId, whereMap, whereOrMap, offset, limit)
	saveKeysKey := rediskey.GroupMemberSearchKeyList(in.GroupId)
	// 从缓存中获取
	// zrange all
	var userIds []string
	var err error
	userIds, err = l.svcCtx.Redis().ZrangeCtx(l.ctx, rdsKey, 0, -1)
	if err != nil || (len(userIds) == 0) {
		if err != nil {
			l.Errorf("redis zrange %s error: %v", saveKeysKey, err)
		}
		userIds, err = l.getGroupMemberListFromDb(in, whereMap, whereOrMap, offset, limit, rdsKey, saveKeysKey)
		if err != nil {
			l.Errorf("getGroupMemberListFromDb error: %v", err)
			return &pb.GetGroupMemberListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	if len(userIds) == 1 && userIds[0] == xredis.NotFound {
		return &pb.GetGroupMemberListResp{
			CommonResp: pb.NewSuccessResp(),
		}, nil
	}
	if in.Opt != nil && in.Opt.OnlyId != nil && *in.Opt.OnlyId {
		var groupMemberList []*pb.GroupMemberInfo
		for _, id := range userIds {
			groupMemberList = append(groupMemberList, &pb.GroupMemberInfo{
				MemberId: id,
			})
		}
		return &pb.GetGroupMemberListResp{
			CommonResp:      pb.NewSuccessResp(),
			GroupMemberList: groupMemberList,
		}, nil
	} else {
		// 获取成员信息
		var memberInfoByIds *pb.MapGroupMemberInfoByIdsResp
		xtrace.StartFuncSpan(l.ctx, "", func(ctx context.Context) {
			memberInfoByIds, err = NewMapGroupMemberInfoByIdsLogic(ctx, l.svcCtx).MapGroupMemberInfoByIds(&pb.MapGroupMemberInfoByIdsReq{
				CommonReq: in.CommonReq,
				GroupId:   in.GroupId,
				MemberIds: userIds,
				Opt:       &pb.MapGroupMemberInfoByIdsReq_Opt{UserBaseInfo: true},
			})
		})
		if err != nil {
			l.Errorf("MapGroupMemberInfoByIds error: %v", err)
			return &pb.GetGroupMemberListResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		var groupMemberList []*pb.GroupMemberInfo
		for _, id := range userIds {
			info, ok := memberInfoByIds.GroupMemberInfoMap[id]
			if !ok {
				continue
			}
			groupMemberList = append(groupMemberList, info)
		}
		return &pb.GetGroupMemberListResp{
			CommonResp:      pb.NewSuccessResp(),
			GroupMemberList: groupMemberList,
		}, nil
	}
}

func (l *GetGroupMemberListLogic) getGroupMemberListFromDb(in *pb.GetGroupMemberListReq, whereMap map[string]interface{}, orMap map[string][]interface{}, offset int32, limit int32, rdsKey string, saveKeysKey string) ([]string, error) {
	// 先把 key 保存起来
	// hset
	// 加锁
	lockKey := "lock:" + saveKeysKey
	// setnx lockKey 1
	acquire, err := l.svcCtx.Redis().SetnxEx(lockKey, lockKey, 10)
	if err != nil {
		l.Errorf("redis setnx %s error: %v", lockKey, err)
		return []string{}, err
	}
	if acquire {
		// 获取到锁
		err = l.svcCtx.Redis().HsetCtx(l.ctx, saveKeysKey, rdsKey, rdsKey)
		if err != nil {
			l.Errorf("redis hset %s error: %v", saveKeysKey, err)
			return make([]string, 0), err
		}
		// 立刻释放锁
		// del lockKey
		_, _ = l.svcCtx.Redis().Del(lockKey)
	}
	tx := l.svcCtx.Mysql().Model(&groupmodel.GroupMember{})
	if len(whereMap) > 0 {
		tx = tx.Where(whereMap)
	}
	if len(orMap) > 0 {
		where := ""
		args := make([]interface{}, 0)
		for k, v := range orMap {
			whereOr := ""
			for _, vv := range v {
				whereOr += fmt.Sprintf("%s = ? or ", k)
				args = append(args, vv)
			}
			whereOr = whereOr[:len(whereOr)-4]
			where += fmt.Sprintf("(%s) or ", whereOr)
		}
		where = where[:len(where)-4]
		where = fmt.Sprintf("(%s)", where)
		tx = tx.Where(where, args...)
	}
	var userIds []string
	err = tx.Order("role desc, remark asc, userId asc").Offset(int(offset)).Limit(int(limit)).Pluck("userId", &userIds).Error
	if err != nil {
		l.Errorf("mysql pluck error: %v", err)
		return make([]string, 0), err
	}
	// 保存到缓存
	// zadd
	if acquire {
		var pairs []redis.Pair
		for i, userId := range userIds {
			pairs = append(pairs, redis.Pair{
				Key:   userId,
				Score: int64(i),
			})
		}
		if len(pairs) > 0 {
			_, err = l.svcCtx.Redis().ZaddsCtx(l.ctx, rdsKey, pairs...)
			if err != nil {
				l.Errorf("redis zadd %s error: %v", rdsKey, err)
				return userIds, nil
			}
		} else {
			// 保存空值
			_, err = l.svcCtx.Redis().ZaddsCtx(l.ctx, rdsKey, redis.Pair{
				Key:   xredis.NotFound,
				Score: 0,
			})
			if err != nil {
				l.Errorf("redis zadd %s error: %v", rdsKey, err)
				return userIds, nil
			}
		}
		_ = l.svcCtx.Redis().ExpireCtx(l.ctx, rdsKey, rediskey.GroupMemberSearchKeyExpire())
	}
	return userIds, nil
}
