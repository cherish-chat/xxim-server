package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetGroupInfoLogic {
	return &ResetGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResetGroupInfoReq 重设群信息
func (l *ResetGroupInfoLogic) ResetGroupInfo(in *pb.ResetGroupInfoReq) (*pb.EditGroupInfoResp, error) {
	err := groupmodel.CleanGroupCache(l.ctx, l.svcCtx.Redis(), in.GroupId)
	if err != nil {
		l.Errorf("flush group cache failed, err: %v", err)
		return &pb.EditGroupInfoResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 更新用户信息
	updateMap := map[string]interface{}{}
	if in.Name != "" {
		updateMap["name"] = in.Name
	}
	if in.Avatar != "" {
		updateMap["avatar"] = in.Avatar
	}
	if in.Introduction != "" {
		updateMap["description"] = in.Introduction
	}
	updateMap["allMute"] = in.AllMute
	updateMap["allMuterType"] = pb.AllMuterType_NORMAL
	updateMap["memberCanAddFriend"] = in.MemberCanAddFriend
	updateMap["canAddMember"] = in.CanAddMember
	if len(updateMap) == 0 {
		return &pb.EditGroupInfoResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err = xorm.Update(tx, &groupmodel.Group{}, updateMap, xorm.Where("id = ?", in.GroupId))
		if err != nil {
			l.Errorf("update group failed, err: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvIdGroup(in.GroupId),
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvNotice: false,
			},
			ContentType: pb.NoticeContentType_UpdateGroupInfo,
			Content: utils.AnyToBytes(pb.NoticeContent_UpdateGroupInfo{
				GroupId:   in.GroupId,
				UpdateMap: updateMap,
			}),
			UniqueId: "updateGroupInfo",
			Title:    "",
			Ext:      nil,
		}
		err = notice.Insert(l.ctx, tx, l.svcCtx.Redis())
		if err != nil {
			l.Errorf("insert notice failed, err: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		err := groupmodel.CleanGroupCache(l.ctx, l.svcCtx.Redis(), in.GroupId)
		if err != nil {
			l.Errorf("flush group cache failed, err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.EditGroupInfoResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	utils.RetryProxy(l.ctx, 5, 1*time.Second, func() error {
		_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
			CommonReq: in.GetCommonReq(),
			UserId:    "",
			ConvId:    pb.HiddenConvIdGroup(in.GroupId),
			DeviceId:  nil,
		})
		if err != nil {
			l.Errorf("SendNoticeData failed, err: %v", err)
		}
		return err
	})
	return &pb.EditGroupInfoResp{}, nil
}
