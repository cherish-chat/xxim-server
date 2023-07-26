package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserRemarkLogic {
	return &UpdateUserRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserRemarkLogic) UpdateUserRemark(in *pb.UpdateUserRemarkReq) (*pb.UpdateUserRemarkResp, error) {
	relationmodel.FlushUserRemarkCache(l.svcCtx.Redis(), in.GetCommonReq().GetUserId())
	defer relationmodel.FlushUserRemarkCache(l.svcCtx.Redis(), in.GetCommonReq().GetUserId())
	model := &relationmodel.UserRemark{
		UserId:   in.GetCommonReq().GetUserId(),
		TargetId: in.TargetId,
		Remark:   in.Remark,
	}
	// upsert
	_ = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		return tx.Save(model).Error
	}, func(tx *gorm.DB) error {
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvIdCommand(),
			UserId: in.CommonReq.UserId,
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvNotice: false,
			},
			ContentType: pb.NoticeContentType_SyncFriendList,
			Content: utils.AnyToBytes(pb.NoticeContent_SyncFriendList{
				Comment: "updateUserRemark",
			}),
			UniqueId: "syncFriendList",
			Title:    "",
			Ext:      nil,
		}
		err := notice.Insert(l.ctx, tx, l.svcCtx.Redis())
		if err != nil {
			l.Errorf("insert notice failed, err: %v", err)
			return err
		}
		return nil
	})
	userId := in.GetCommonReq().GetUserId()
	go utils.RetryProxy(context.Background(), 12, time.Second*3, func() error {
		_, err := l.svcCtx.NoticeService().GetUserNoticeData(context.Background(), &pb.GetUserNoticeDataReq{
			UserId: userId,
			ConvId: pb.HiddenConvIdCommand(),
		})
		if err != nil {
			l.Errorf("SendNoticeData failed, err: %v", err)
			return err
		}
		return nil
	})
	return &pb.UpdateUserRemarkResp{}, nil
}
