package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"strings"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListByUserIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListByUserIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListByUserIdLogic {
	return &GetFriendListByUserIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListByUserIdLogic) GetFriendListByUserId(in *pb.GetFriendListByUserIdReq) (*pb.GetFriendListByUserIdResp, error) {
	if len(in.Filter) == 0 {
		in.Filter = make(map[string]string)
	}
	var friendIds []string
	var total int64
	model := &relationmodel.Friend{}
	tx := l.svcCtx.Mysql().Model(model).Where("userId = ?", in.UserId)
	var nicknameLike *string
	var fids []string
	var role *int32
	for k, v := range in.Filter {
		if v == "" {
			continue
		}
		switch k {
		case "nickname":
			v := v
			nicknameLike = utils.AnyPtr(v)
		case "role":
			if v == "normal" {
				role = utils.AnyPtr(int32(0))
			} else if v == "service" {
				role = utils.AnyPtr(int32(1))
			} else if v == "guest" {
				role = utils.AnyPtr(int32(2))
			}
		}
	}
	if uid, ok := in.Filter["id"]; ok && uid != "" {
		tx = tx.Where("friendId = ?", uid)
	} else {
		if nicknameLike != nil {
			// 先查询所有用户的好友
			resp, err := NewGetFriendListLogic(l.ctx, l.svcCtx).GetFriendList(&pb.GetFriendListReq{
				CommonReq: &pb.CommonReq{UserId: in.UserId},
				Page:      &pb.Page{Page: 1, Size: 10000},
				Opt:       pb.GetFriendListReq_WithBaseInfo,
			})
			if err != nil {
				l.Errorf("GetFriendListByUserId err: %v", err)
				return &pb.GetFriendListByUserIdResp{CommonResp: pb.NewRetryErrorResp()}, err
			}
			// 再查询用户信息
			for _, user := range resp.UserMap {
				if nicknameLike != nil {
					if strings.Contains(user.Nickname, *nicknameLike) {
						fids = append(fids, user.Id)
					}
				}
				if role != nil {
					if user.Role == *role {
						fids = append(fids, user.Id)
					}
				}
			}
			fids = utils.Set(fids)
			if len(fids) == 0 {
				return &pb.GetFriendListByUserIdResp{CommonResp: pb.NewSuccessResp(), Total: total}, nil
			}
			tx = tx.Where("friendId in (?)", fids)
		}
	}
	tx.Count(&total)
	err := xorm.Paging(tx.Order("friendId ASC"), in.Page.Page, in.Page.Size).Pluck("friendId", &friendIds).Error
	if err != nil {
		l.Errorf("GetFriendListByUserId err: %v", err)
		return &pb.GetFriendListByUserIdResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(friendIds) == 0 {
		return &pb.GetFriendListByUserIdResp{CommonResp: pb.NewSuccessResp(), Total: total}, nil
	}
	// 获取用户信息去
	mapUserByIds, err := l.svcCtx.UserService().MapUserByIds(l.ctx, &pb.MapUserByIdsReq{Ids: friendIds})
	if err != nil {
		l.Errorf("GetFriendListByUserId err: %v", err)
		return &pb.GetFriendListByUserIdResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.GetFriendListByUserIdItem
	for _, v := range friendIds {
		user := &usermodel.User{}
		bytes, ok := mapUserByIds.Users[v]
		if ok {
			user = usermodel.UserFromBytes(bytes)
		}
		resp = append(resp, &pb.GetFriendListByUserIdItem{
			UserId:   v,
			Avatar:   user.Avatar,
			Nickname: user.Nickname,
			Role:     user.Role.String(),
		})
	}
	return &pb.GetFriendListByUserIdResp{
		FriendList: resp,
		Total:      total,
	}, nil
}
