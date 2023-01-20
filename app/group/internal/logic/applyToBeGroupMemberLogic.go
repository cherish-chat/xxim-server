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

type ApplyToBeGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyToBeGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyToBeGroupMemberLogic {
	return &ApplyToBeGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ApplyToBeGroupMember 申请加入群聊
func (l *ApplyToBeGroupMemberLogic) ApplyToBeGroupMember(in *pb.ApplyToBeGroupMemberReq) (*pb.ApplyToBeGroupMemberResp, error) {
	apply := &groupmodel.GroupApply{
		Id:         utils.GenId(),
		GroupId:    in.GroupId,
		UserId:     in.CommonReq.UserId,
		Result:     pb.GroupApplyHandleResult_UNHANDLED,
		Reason:     in.Reason,
		ApplyTime:  time.Now().UnixMilli(),
		HandleTime: 0,
	}
	// 判断是否申请过了
	var count int64
	err := l.svcCtx.Mysql().Model(&groupmodel.GroupApply{}).Where(
		"groupId = ? and userId = ? and result = ?",
		in.GroupId, in.CommonReq.UserId, pb.GroupApplyHandleResult_UNHANDLED,
	).Count(&count).Error
	if err != nil {
		l.Errorf("ApplyToBeGroupMember Count error: %v", err)
		return &pb.ApplyToBeGroupMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if count > 0 {
		// 直接返回成功
		return &pb.ApplyToBeGroupMemberResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	// 获取群里所有的管理员
	groupManagers, err := getAllGroupManager(l.ctx, l.svcCtx, in.GroupId, true)
	if err != nil {
		l.Errorf("ApplyToBeGroupMember getAllGroupManager error: %v", err)
		return &pb.ApplyToBeGroupMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err := xorm.InsertOne(tx, apply)
		if err != nil {
			l.Errorf("ApplyToBeGroupMember InsertOne error: %v", err)
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
		l.Errorf("ApplyToBeGroupMember Transaction error: %v", err)
		return &pb.ApplyToBeGroupMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
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
	return &pb.ApplyToBeGroupMemberResp{}, nil
}
