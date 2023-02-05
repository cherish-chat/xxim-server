package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAppMgmtNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAppMgmtNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAppMgmtNoticeLogic {
	return &AddAppMgmtNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddAppMgmtNoticeLogic) AddAppMgmtNotice(in *pb.AddAppMgmtNoticeReq) (*pb.AddAppMgmtNoticeResp, error) {
	model := &appmgmtmodel.Notice{
		Id:         appmgmtmodel.GetId(l.svcCtx.Mysql(), &appmgmtmodel.Notice{}, 10000),
		Position:   int8(in.AppMgmtNotice.Position),
		Platform:   in.AppMgmtNotice.Platform,
		Title:      in.AppMgmtNotice.Title,
		Image:      in.AppMgmtNotice.Image,
		Content:    in.AppMgmtNotice.Content,
		Sort:       in.AppMgmtNotice.Sort,
		IsEnable:   in.AppMgmtNotice.IsEnable,
		StartTime:  in.AppMgmtNotice.StartTime,
		EndTime:    in.AppMgmtNotice.EndTime,
		CreateTime: time.Now().UnixMilli(),
	}
	err := model.Insert(l.svcCtx.Mysql())
	if err != nil {
		l.Errorf("insert err: %v", err)
		return &pb.AddAppMgmtNoticeResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.AddAppMgmtNoticeResp{}, nil
}
