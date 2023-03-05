package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppMgmtLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppMgmtLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppMgmtLinkLogic {
	return &AddAppMgmtLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAppMgmtLinkLogic) AddAppMgmtLink(in *pb.AddAppMgmtLinkReq) (*pb.AddAppMgmtLinkResp, error) {
	model := &appmgmtmodel.Link{
		Id:         appmgmtmodel.GetId(l.svcCtx.Mysql(), &appmgmtmodel.Link{}, 10000),
		Sort:       in.AppMgmtLink.Sort,
		Name:       in.AppMgmtLink.Name,
		Url:        in.AppMgmtLink.Url,
		Icon:       in.AppMgmtLink.Icon,
		IsEnable:   in.AppMgmtLink.IsEnable,
		CreateTime: time.Now().UnixMilli(),
	}
	err := model.Insert(l.svcCtx.Mysql())
	if err != nil {
		l.Errorf("insert err: %v", err)
		return &pb.AddAppMgmtLinkResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.AddAppMgmtLinkResp{}, nil
}
