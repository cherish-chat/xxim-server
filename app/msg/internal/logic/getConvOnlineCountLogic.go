package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConvOnlineCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConvOnlineCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConvOnlineCountLogic {
	return &GetConvOnlineCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetConvOnlineCount 获取一个会话里所有的在线用户
func (l *GetConvOnlineCountLogic) GetConvOnlineCount(in *pb.GetConvOnlineCountReq) (*pb.GetConvOnlineCountResp, error) {
	// 判断是否是单聊
	var userIds []string
	if pb.IsSingleConv(in.ConvId) {
		userIds = pb.ParseSingleConv(in.ConvId)
	} else if pb.IsGroupConv(in.ConvId) {
		groupId := pb.ParseGroupConv(in.ConvId)
		// 获取群成员
		memberList, err := l.svcCtx.GroupService().GetGroupMemberList(l.ctx, &pb.GetGroupMemberListReq{
			CommonReq: in.CommonReq,
			GroupId:   groupId,
			Page: &pb.Page{
				Page: 1,
				Size: 100000,
			},
			Filter: nil,
			Opt: &pb.GetGroupMemberListReq_GetGroupMemberListOpt{
				OnlyId: utils.AnyPtr(true),
			},
		})
		if err != nil {
			l.Errorf("get group member list error: %v", err)
			return &pb.GetConvOnlineCountResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		for _, member := range memberList.GroupMemberList {
			userIds = append(userIds, member.UserId)
		}
	}
	userIds = utils.Set(userIds)
	// 获取在线用户
	userConnResp, err := l.svcCtx.ImService().GetUserConn(l.ctx, &pb.GetUserConnReq{
		UserIds:   userIds,
		Platforms: nil,
		Devices:   nil,
	})
	if err != nil {
		l.Errorf("get user conn error: %v", err)
		return &pb.GetConvOnlineCountResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var userIdMap = make(map[string]bool)
	var deviceIdMap = make(map[string]bool)
	for _, conn := range userConnResp.ConnParams {
		userIdMap[conn.UserId] = true
		deviceIdMap[conn.DeviceId] = true
	}
	return &pb.GetConvOnlineCountResp{
		CommonResp: pb.NewSuccessResp(),
		User:       int32(len(userIdMap)),
		Device:     int32(len(deviceIdMap)),
	}, nil
}
