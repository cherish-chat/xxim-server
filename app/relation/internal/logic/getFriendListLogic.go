package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendListLogic {
	return &GetFriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendListLogic) GetFriendList(in *pb.GetFriendListReq) (*pb.GetFriendListResp, error) {
	if in.Opt == pb.GetFriendListReq_WithBaseInfo {
		resp, err := l.getFriendListWithBaseInfo(in)
		if err != nil {
			return resp, err
		}
		if in.WithConvSetting {
			return l.withConvSetting(in, resp)
		} else {
			return resp, err
		}
	} else if in.Opt == pb.GetFriendListReq_OnlyId {
		resp, err := l.getFriendListOnlyId(in)
		if err != nil {
			return resp, err
		}
		if in.WithConvSetting {
			return l.withConvSetting(in, resp)
		} else {
			return resp, err
		}
	} else if in.Opt == pb.GetFriendListReq_WithBaseInfoAndRemark {
		getFriendListResp, err := l.getFriendListWithBaseInfo(in)
		if err != nil {
			return getFriendListResp, err
		}
		resp, err := l.getFriendListWithRemark(in, getFriendListResp)
		if err != nil {
			return resp, err
		}
		if in.WithConvSetting {
			return l.withConvSetting(in, resp)
		} else {
			return resp, err
		}
	}
	logx.Errorf("get friend list error: opt is not support: %v", in.Opt)
	return &pb.GetFriendListResp{}, nil
}

func (l *GetFriendListLogic) getFriendListWithBaseInfo(in *pb.GetFriendListReq) (*pb.GetFriendListResp, error) {
	myFriendList, err := l.getFriendListOnlyId(in)
	if err != nil {
		return myFriendList, err
	}
	userMap := make(map[string]*pb.UserBaseInfo)
	userBaseInfos, err := l.svcCtx.UserService().BatchGetUserBaseInfo(l.ctx, &pb.BatchGetUserBaseInfoReq{Ids: myFriendList.Ids})
	if err != nil {
		l.Errorf("get friend list error: %v", err)
		return &pb.GetFriendListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	for _, userBaseInfo := range userBaseInfos.UserBaseInfos {
		if userBaseInfo.Id == "" {
			continue
		}
		userMap[userBaseInfo.Id] = userBaseInfo
	}
	myFriendList.UserMap = userMap
	return myFriendList, nil
}

func (l *GetFriendListLogic) getFriendListOnlyId(in *pb.GetFriendListReq) (*pb.GetFriendListResp, error) {
	myFriendList, err := relationmodel.GetMyFriendList(l.ctx, l.svcCtx.Redis(), l.svcCtx.Mysql(), in.CommonReq.UserId)
	if err != nil {
		l.Errorf("get friend list error: %v", err)
		return &pb.GetFriendListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	myFriendList = append(myFriendList, in.CommonReq.UserId)
	return &pb.GetFriendListResp{CommonResp: pb.NewSuccessResp(), Ids: myFriendList}, nil
}

func (l *GetFriendListLogic) getFriendListWithRemark(in *pb.GetFriendListReq, getFriendListResp *pb.GetFriendListResp) (*pb.GetFriendListResp, error) {
	var targetIds []string
	for _, user := range getFriendListResp.UserMap {
		targetIds = append(targetIds, user.Id)
	}
	var mapUserRemarkResp *pb.MapUserRemarkResp
	var err error
	xtrace.StartFuncSpan(l.ctx, "MapUserRemark", func(ctx context.Context) {
		mapUserRemarkResp, err = NewMapUserRemarkLogic(ctx, l.svcCtx).MapUserRemark(&pb.MapUserRemarkReq{
			CommonReq: in.CommonReq,
			TargetIds: targetIds,
		})
	})
	if err != nil {
		l.Errorf("get friend list error: %v", err)
		return &pb.GetFriendListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	getFriendListResp.RemarkMap = mapUserRemarkResp.RemarkMap
	return getFriendListResp, nil
}

func (l *GetFriendListLogic) withConvSetting(in *pb.GetFriendListReq, resp *pb.GetFriendListResp) (*pb.GetFriendListResp, error) {
	convIds := make([]string, 0)
	for _, userId := range resp.Ids {
		convIds = append(convIds, pb.SingleConvId(userId, in.GetCommonReq().GetUserId()))
	}
	resp.ConvSettingMap = make(map[string]*pb.ConvSetting)
	resp.ConvSetting2Map = make(map[string]*pb.ConvSettingProto2)
	if len(convIds) == 0 {
		return resp, nil
	}
	getConvSettingResp, err := l.svcCtx.ImService().GetConvSetting(l.ctx, &pb.GetConvSettingReq{
		CommonReq: in.GetCommonReq(),
		ConvIds:   convIds,
	})
	if err != nil {
		l.Errorf("GetConvSetting err: %v", err)
		return resp, err
	}
	for _, setting := range getConvSettingResp.ConvSettings {
		resp.ConvSettingMap[setting.ConvId] = setting
	}
	for _, setting := range getConvSettingResp.ConvSetting2S {
		resp.ConvSetting2Map[setting.ConvId] = setting
	}
	return resp, nil
}
