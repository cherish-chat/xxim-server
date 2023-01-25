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
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvIdFriend(in.CommonReq.UserId),
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvMsg:    false,
			},
			ContentType: pb.NoticeContentType_UpdateUserInfo,
			Content: utils.AnyToBytes(pb.NoticeContent_UpdateUserInfo{
				UserId:    in.CommonReq.UserId,
				UpdateMap: updateMap,
			}),
			Title: "",
			Ext:   nil,
		}
		err = notice.Insert(l.ctx, tx)
		if err != nil {
			l.Errorf("insert notice failed, err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.UpdateUserInfoResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	utils.RetryProxy(l.ctx, 5, 1*time.Second, func() error {
		_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
			CommonReq: in.GetCommonReq(),
			UserId:    "",
			ConvId:    pb.HiddenConvIdFriend(in.CommonReq.UserId),
			DeviceId:  nil,
		})
		if err != nil {
			l.Errorf("SendNoticeData failed, err: %v", err)
		}
		return err
	})
	return &pb.UpdateUserInfoResp{}, nil
}
