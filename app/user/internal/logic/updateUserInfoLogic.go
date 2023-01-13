package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *pb.UpdateUserInfoReq) (*pb.UpdateUserInfoResp, error) {
	err := usermodel.FlushUserCache(l.ctx, l.svcCtx.Redis(), []string{in.CommonReq.UserId})
	if err != nil {
		l.Errorf("flush user cache failed, err: %v", err)
		return &pb.UpdateUserInfoResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 更新用户信息
	updateMap := map[string]interface{}{}
	if in.Nickname != nil {
		updateMap["nickname"] = *in.Nickname
	}
	if in.Avatar != nil {
		updateMap["avatar"] = *in.Avatar
	}
	if in.Signature != nil {
		updateMap["signature"] = *in.Signature
	}
	if len(updateMap) == 0 {
		return &pb.UpdateUserInfoResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err = xorm.Update(tx, &usermodel.User{}, updateMap, xorm.Where("id = ?", in.CommonReq.UserId))
		if err != nil {
			l.Errorf("update user failed, err: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		// 发送一条订阅号消息 订阅号的convId = notice:selfId  noticeId = UpdateUserInfo
		data := &pb.NoticeData{
			ConvId:         noticemodel.ConvIdUser(in.CommonReq.UserId),
			UnreadCount:    0,
			UnreadAbsolute: false,
			NoticeId:       "UpdateUserInfo",
			ContentType:    0,
			Content:        []byte{},
			Options: &pb.NoticeData_Options{
				StorageForClient: false,
				UpdateConvMsg:    false,
				OnlinePushOnce:   false,
			},
			Ext: nil,
		}
		m := noticemodel.NoticeFromPB(data, true, "")
		err := m.Upsert(tx)
		if err != nil {
			l.Errorf("Upsert failed, err: %v", err)
		}
		return err
	})
	if err != nil {
		return &pb.UpdateUserInfoResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 刷新订阅
	utils.RetryProxy(l.ctx, 5, 1*time.Second, func() error {
		_, err = l.svcCtx.NoticeService().SendNoticeData(l.ctx, &pb.SendNoticeDataReq{
			CommonReq:   in.CommonReq,
			NoticeData:  &pb.NoticeData{NoticeId: "UpdateUserInfo"},
			UserId:      nil,
			IsBroadcast: utils.AnyPtr(true),
			Inserted:    utils.AnyPtr(true),
		})
		if err != nil {
			l.Errorf("SendNoticeData failed, err: %v", err)
		}
		return err
	})
	return &pb.UpdateUserInfoResp{}, nil
}
