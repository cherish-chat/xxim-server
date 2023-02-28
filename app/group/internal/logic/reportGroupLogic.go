package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReportGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReportGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReportGroupLogic {
	return &ReportGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ReportGroup 举报群
func (l *ReportGroupLogic) ReportGroup(in *pb.ReportGroupReq) (*pb.ReportGroupResp, error) {
	model := &groupmodel.ReportRecord{
		Id:            utils.GenId(),
		ReporterId:    in.CommonReq.UserId,
		ReportedId:    in.GroupId,
		ReportType:    "",
		ReportContent: in.Reason,
		ReportImages:  make([]string, 0),
		ReportTime:    time.Now().UnixMilli(),
		ReportStatus:  "",
		HandleTime:    0,
		HandlerId:     "",
	}
	l.svcCtx.Mysql().Create(model)
	return &pb.ReportGroupResp{}, nil
}
