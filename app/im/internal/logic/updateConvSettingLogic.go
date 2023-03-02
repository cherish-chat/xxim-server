package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConvSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateConvSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConvSettingLogic {
	return &UpdateConvSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateConvSettingLogic) UpdateConvSetting(in *pb.UpdateConvSettingReq) (*pb.UpdateConvSettingResp, error) {
	isSingleConv := pb.IsSingleConv(in.ConvId)
	isGroupConv := pb.IsGroupConv(in.ConvId)
	if isSingleConv {
		return l.updateSingleConvSetting(in)
	} else if isGroupConv {
		return l.updateGroupConvSetting(in)
	}
	return &pb.UpdateConvSettingResp{}, nil
}

func (l *UpdateConvSettingLogic) updateSingleConvSetting(in *pb.UpdateConvSettingReq) (*pb.UpdateConvSettingResp, error) {
	req := &pb.SetSingleConvSettingReq{
		CommonReq: in.CommonReq,
		Setting: &pb.SingleConvSetting{
			ConvId:            in.ConvId,
			UserId:            in.CommonReq.UserId,
			IsTop:             in.IsTop,
			IsDisturb:         nil,
			NotifyPreview:     nil,
			NotifySound:       nil,
			NotifyCustomSound: nil,
			NotifyVibrate:     nil,
			IsShield:          nil,
			ChatBg:            nil,
		},
	}
	resp, err := l.svcCtx.RelationService().SetSingleConvSetting(l.ctx, req)
	if err != nil {
		l.Errorf("SetSingleConvSetting err: %v", err)
		return &pb.UpdateConvSettingResp{}, err
	}
	return &pb.UpdateConvSettingResp{CommonResp: resp.GetCommonResp()}, nil
}

func (l *UpdateConvSettingLogic) updateGroupConvSetting(in *pb.UpdateConvSettingReq) (*pb.UpdateConvSettingResp, error) {
	groupId := pb.ParseGroupConv(in.ConvId)
	req := &pb.SetGroupMemberInfoReq{
		CommonReq:   in.CommonReq,
		GroupId:     groupId,
		MemberId:    in.CommonReq.UserId,
		Notice:      "",
		Remark:      nil,
		Role:        nil,
		UnbanTime:   nil,
		GroupRemark: nil,
		IsTop:       in.IsTop,
	}
	resp, err := l.svcCtx.GroupService().SetGroupMemberInfo(l.ctx, req)
	if err != nil {
		l.Errorf("SetGroupMemberInfo err: %v", err)
		return &pb.UpdateConvSettingResp{}, err
	}
	return &pb.UpdateConvSettingResp{CommonResp: resp.GetCommonResp()}, nil
}
