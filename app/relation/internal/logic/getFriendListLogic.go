package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(in *pb.GetFriendListReq) (*pb.GetFriendListResp, error) {
	if in.Opt == pb.GetFriendListReq_WithBaseInfo {
		return l.getFriendListWithBaseInfo(in)
	} else if in.Opt == pb.GetFriendListReq_OnlyId {
		return l.getFriendListOnlyId(in)
	}
	return &pb.GetFriendListResp{}, nil
}

func (l *GetFriendListLogic) getFriendListWithBaseInfo(in *pb.GetFriendListReq) (*pb.GetFriendListResp, error) {
	myFriendList, err := l.getFriendListOnlyId(in)
	if err != nil {
		return myFriendList, err
	}
	userMap := make(map[string]*pb.UserBaseInfo)
	userBaseInfos, err := l.svcCtx.UserService().BatchGetUserBaseInfo(l.ctx, &pb.BatchGetUserBaseInfoReq{Ids: myFriendList.Ids})
	if err != nil {
		l.Errorf("get friend list error: %v", err)
		return &pb.GetFriendListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	for _, userBaseInfo := range userBaseInfos.UserBaseInfos {
		userMap[userBaseInfo.Id] = userBaseInfo
	}
	myFriendList.UserMap = userMap
	return myFriendList, nil
}

func (l *GetFriendListLogic) getFriendListOnlyId(in *pb.GetFriendListReq) (*pb.GetFriendListResp, error) {
	myFriendList, err := relationmodel.GetMyFriendList(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.CommonReq.Id)
	if err != nil {
		l.Errorf("get friend list error: %v", err)
		return &pb.GetFriendListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetFriendListResp{CommonResp: pb.NewSuccessResp(), Ids: myFriendList}, nil
}
