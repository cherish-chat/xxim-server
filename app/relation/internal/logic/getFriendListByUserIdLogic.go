package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

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
	var friendIds []string
	var total int64
	model := &relationmodel.Friend{}
	tx := l.svcCtx.Mysql().Model(model).Where("userId = ?", in.UserId)
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
