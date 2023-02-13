package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAppMgmtShieldWordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAppMgmtShieldWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppMgmtShieldWordLogic {
	return &DeleteAppMgmtShieldWordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAppMgmtShieldWordLogic) DeleteAppMgmtShieldWord(in *pb.DeleteAppMgmtShieldWordReq) (*pb.DeleteAppMgmtShieldWordResp, error) {
	model := &appmgmtmodel.ShieldWord{}
	err := l.svcCtx.Mysql().Model(model).Where("id in (?)", in.Ids).Delete(model).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
		return &pb.DeleteAppMgmtShieldWordResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	for _, pod := range l.svcCtx.MsgPodsMgr.AllMsgServices() {
		_, err := pod.FlushShieldWordTireTree(context.Background(), &pb.FlushShieldWordTireTreeReq{})
		if err != nil {
			l.Errorf("flush shield word tire tree err: %v", err)
			return &pb.DeleteAppMgmtShieldWordResp{CommonResp: pb.NewRetryErrorResp()}, nil
		}
	}
	return &pb.DeleteAppMgmtShieldWordResp{}, nil
}
