package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"gorm.io/gorm"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberInfoLogic {
	return &GetGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupMemberInfo 获取群成员信息
func (l *GetGroupMemberInfoLogic) GetGroupMemberInfo(in *pb.GetGroupMemberInfoReq) (*pb.GetGroupMemberInfoResp, error) {
	members, err := groupmodel.ListGroupMemberFromRedis(l.ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.GroupId, []string{in.MemberId})
	if err != nil {
		l.Errorf("getGroupMemberInfoLogic err: %v", err)
		return &pb.GetGroupMemberInfoResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(members) == 0 {
		return &pb.GetGroupMemberInfoResp{CommonResp: pb.NewAlertErrorResp(
			l.svcCtx.T(in.CommonReq.Language, "操作失败"),
			l.svcCtx.T(in.CommonReq.Language, "群成员不存在"),
		)}, gorm.ErrRecordNotFound
	}
	return &pb.GetGroupMemberInfoResp{
		GroupMemberInfo: members[0].Pb(),
	}, nil
}
