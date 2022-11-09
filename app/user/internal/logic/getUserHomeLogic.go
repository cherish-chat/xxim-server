package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserHomeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserHomeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserHomeLogic {
	return &GetUserHomeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserHomeLogic) GetUserHome(in *pb.GetUserHomeReq) (*pb.GetUserHomeResp, error) {
	users, err := usermodel.GetUsersByIds(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mongo().Collection(&usermodel.User{}), []string{in.Id})
	if err != nil {
		l.Errorf("getUsersByIds failed, err: %v", err)
		return &pb.GetUserHomeResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(users) == 0 {
		l.Errorf("user not found, id: %s", in.Id)
		return &pb.GetUserHomeResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.Requester.Language, "用户已注销"))}, nil
	}
	user := users[0]
	resp := &pb.GetUserHomeResp{
		Id:        user.Id,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Xb:        user.Xb,
		Birthday:  user.Birthday,
		IpRegion:  nil,
		Signature: user.InfoMap.Get("signature", l.svcCtx.SystemConfigMgr.Get("signature.if_not_set")),
		LevelInfo: user.LevelInfo.Pb(),
	}
	latestConn, err := l.svcCtx.ImService().GetUserLatestConn(l.ctx, &pb.GetUserLatestConnReq{UserId: user.Id})
	if err != nil {
		l.Errorf("get user latest conn failed, err: %v", err)
	} else {
		resp.IpRegion = latestConn.IpRegion
	}
	return resp, nil
}
