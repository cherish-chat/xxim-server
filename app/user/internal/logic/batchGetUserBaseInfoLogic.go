package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUserBaseInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUserBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUserBaseInfoLogic {
	return &BatchGetUserBaseInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetUserBaseInfoLogic) BatchGetUserBaseInfo(in *pb.BatchGetUserBaseInfoReq) (*pb.BatchGetUserBaseInfoResp, error) {
	usersByIds, err := usermodel.GetUsersByIds(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.Ids)
	if err != nil {
		l.Errorf("BatchGetUserBaseInfoLogic BatchGetUserBaseInfo err: %v", err)
		return &pb.BatchGetUserBaseInfoResp{CommonResp: pb.NewInternalErrorResp()}, err
	}
	var resp = make([]*pb.UserBaseInfo, 0)
	var userConnMap = make(map[string]*pb.GetUserLatestConnResp)
	userLatestConn, err := l.svcCtx.ImService().BatchGetUserLatestConn(l.ctx, &pb.BatchGetUserLatestConnReq{UserIds: in.Ids})
	if err != nil {
		l.Errorf("BatchGetUserBaseInfoLogic BatchGetUserLatestConn err: %v", err)
		return &pb.BatchGetUserBaseInfoResp{CommonResp: pb.NewInternalErrorResp()}, err
	}
	for _, conn := range userLatestConn.UserLatestConns {
		userConnMap[conn.UserId] = conn
	}
	for _, user := range usersByIds {
		if user.Id == "" {
			continue
		}
		userConn, ok := userConnMap[user.Id]
		if !ok {
			userConn = &pb.GetUserLatestConnResp{}
		}
		resp = append(resp, &pb.UserBaseInfo{
			Id:       user.Id,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Xb:       user.Xb,
			Birthday: user.Birthday,
			IpRegion: userConn.IpRegion, // latest connect ip region
		})
	}
	return &pb.BatchGetUserBaseInfoResp{UserBaseInfos: resp}, nil
}
