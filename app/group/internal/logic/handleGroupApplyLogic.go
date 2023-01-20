package logic

import (
	"context"
	"fmt"
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

type HandleGroupApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleGroupApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleGroupApplyLogic {
	return &HandleGroupApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// HandleGroupApply 处理群聊申请
func (l *HandleGroupApplyLogic) HandleGroupApply(in *pb.HandleGroupApplyReq) (*pb.HandleGroupApplyResp, error) {
	// 查询 apply
	apply := &groupmodel.GroupApply{}
	err := l.svcCtx.Mysql().Model(&groupmodel.GroupApply{}).Where("id = ?", in.ApplyId).Limit(1).Find(apply).Error
	if err != nil {
		l.Errorf("HandleGroupApply Find error: %v", err)
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if apply.Id == "" {
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	if apply.HandleTime != 0 {
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	// 获取群里所有的管理员
	groupManagers, err := getAllGroupManager(l.ctx, l.svcCtx, apply.GroupId, true)
	if err != nil {
		l.Errorf("HandleGroupApply getAllGroupManager error: %v", err)
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 判断是否是管理员
	isManager := false
	for _, manager := range groupManagers {
		if manager.UserId == in.CommonReq.UserId {
			isManager = true
			break
		}
	}
	if !isManager {
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewAlertErrorResp(
			l.svcCtx.T(in.CommonReq.Language, "操作失败"),
			l.svcCtx.T(in.CommonReq.Language, "您不是群管理员"),
		)}, nil
	}
	updateMap := map[string]interface{}{
		"result":       in.Result,
		"handleTime":   time.Now().UnixMilli(),
		"handleUserId": in.CommonReq.UserId,
	}
	{
		apply.Result = in.Result
		apply.HandleTime = time.Now().UnixMilli()
		apply.HandleUserId = in.CommonReq.UserId
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		// 更新 apply
		err := tx.Model(&groupmodel.GroupApply{}).Where("id = ?", in.ApplyId).Updates(updateMap).Error
		if err != nil {
			l.Errorf("HandleGroupApply Update error: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		for _, manager := range groupManagers {
			data := &pb.NoticeData{
				ConvId:         noticemodel.ConvId_GroupNotice,
				UnreadCount:    0,
				UnreadAbsolute: false,
				NoticeId:       apply.Id,
				CreateTime:     "",
				Title:          "",
				ContentType:    1,
				Content:        []byte(utils.AnyToString(apply)),
				Options: &pb.NoticeData_Options{
					StorageForClient: false,
					UpdateConvMsg:    false,
					OnlinePushOnce:   false,
				},
				Ext: nil,
			}
			m := noticemodel.NoticeFromPB(data, false, manager.UserId)
			err := m.Upsert(tx)
			if err != nil {
				l.Errorf("Upsert failed, err: %v", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		l.Errorf("HandleGroupApply Transaction error: %v", err)
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 通知给群里所有的管理员
	go utils.RetryProxy(l.ctx, 12, time.Second, func() error {
		for _, manager := range groupManagers {
			_, err := l.svcCtx.NoticeService().SendNoticeData(l.ctx, &pb.SendNoticeDataReq{
				CommonReq: in.CommonReq,
				NoticeData: &pb.NoticeData{
					NoticeId: fmt.Sprintf("%s", apply.Id),
					ConvId:   noticemodel.ConvId_FriendNotice,
				},
				UserId:      utils.AnyPtr(manager.UserId),
				IsBroadcast: nil,
				Inserted:    utils.AnyPtr(true),
			})
			if err != nil {
				l.Errorf("ApplyToBeGroupMember SendNoticeData error: %v", err)
				return err
			}
		}
		return nil
	})
	return &pb.HandleGroupApplyResp{}, nil
}
