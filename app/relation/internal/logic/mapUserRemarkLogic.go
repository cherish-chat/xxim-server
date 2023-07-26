package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MapUserRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMapUserRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MapUserRemarkLogic {
	return &MapUserRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MapUserRemarkLogic) MapUserRemark(in *pb.MapUserRemarkReq) (*pb.MapUserRemarkResp, error) {
	if len(in.TargetIds) == 0 {
		return &pb.MapUserRemarkResp{
			RemarkMap: make(map[string]string),
		}, nil
	}
	remarkMap, err := relationmodel.GetUserRemarkMap(l.svcCtx.Redis(), l.svcCtx.Mysql(), in.GetCommonReq().GetUserId())
	if err != nil {
		l.Errorf("MapUserRemark - error: %v", err)
		return nil, err
	}
	m := make(map[string]string)
	for _, targetId := range in.TargetIds {
		v, ok := remarkMap[targetId]
		if ok {
			m[targetId] = v
		} else {
			m[targetId] = ""
		}
	}
	return &pb.MapUserRemarkResp{
		CommonResp: pb.NewSuccessResp(),
		RemarkMap:  m,
	}, nil
}
