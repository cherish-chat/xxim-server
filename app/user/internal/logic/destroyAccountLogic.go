package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xjwt"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DestroyAccountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDestroyAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DestroyAccountLogic {
	return &DestroyAccountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DestroyAccountLogic) DestroyAccount(in *pb.DestroyAccountReq) (*pb.DestroyAccountResp, error) {
	err := usermodel.FlushUserCache(l.ctx, l.svcCtx.Redis(), []string{in.CommonReq.UserId})
	if err != nil {
		l.Errorf("flush user cache failed, err: %v", err)
		return &pb.DestroyAccountResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 更新用户信息
	updateMap := map[string]interface{}{
		"destroyTime": time.Now().UnixMilli(),
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
				UpdateConvNotice: false,
			},
			ContentType: pb.NoticeContentType_UpdateUserInfo,
			Content: utils.AnyToBytes(pb.NoticeContent_UpdateUserInfo{
				UserId:    in.CommonReq.UserId,
				UpdateMap: updateMap,
			}),
			UniqueId: "updateUserInfo",
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
		err := xjwt.BanToken(l.ctx, l.svcCtx.Redis(), in.CommonReq.UserId, "您的账号已注销")
		if err != nil {
			l.Errorf("flush user cache failed, err: %v", err)
			return err
		}
		err = usermodel.FlushUserCache(l.ctx, l.svcCtx.Redis(), []string{in.CommonReq.UserId})
		if err != nil {
			l.Errorf("flush user cache failed, err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.DestroyAccountResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	_, err = l.svcCtx.ImService().KickUserConn(l.ctx, &pb.KickUserConnReq{GetUserConnReq: &pb.GetUserConnReq{
		UserIds: []string{in.CommonReq.UserId},
	}})
	if err != nil {
		l.Errorf("kick user conn failed, err: %v", err)
		return &pb.DestroyAccountResp{CommonResp: pb.NewRetryErrorResp()}, err
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
	return &pb.DestroyAccountResp{}, nil
}
