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

type EditGroupInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditGroupInfoLogic {
	return &EditGroupInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EditGroupInfo 编辑群信息
func (l *EditGroupInfoLogic) EditGroupInfo(in *pb.EditGroupInfoReq) (*pb.EditGroupInfoResp, error) {
	err := groupmodel.CleanGroupCache(l.ctx, l.svcCtx.Redis(), in.GroupId)
	if err != nil {
		l.Errorf("flush group cache failed, err: %v", err)
		return &pb.EditGroupInfoResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 更新用户信息
	updateMap := map[string]interface{}{}
	if in.Name != nil {
		updateMap["name"] = *in.Name
	}
	if in.Avatar != nil {
		updateMap["avatar"] = *in.Avatar
	}
	if in.Introduction != nil {
		updateMap["description"] = *in.Introduction
	}
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
