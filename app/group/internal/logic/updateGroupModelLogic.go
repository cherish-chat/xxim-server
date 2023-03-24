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

type UpdateGroupModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGroupModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupModelLogic {
	return &UpdateGroupModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateGroupModel 更新群组
func (l *UpdateGroupModelLogic) UpdateGroupModel(in *pb.UpdateGroupModelReq) (*pb.UpdateGroupModelResp, error) {
	err := groupmodel.CleanGroupCache(l.ctx, l.svcCtx.Redis(), in.GroupModel.Id)
	if err != nil {
		l.Errorf("flush group cache failed, err: %v", err)
		return &pb.UpdateGroupModelResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 查询原模型
	model := &groupmodel.Group{}
	err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.GroupModel.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateGroupModelResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	if in.GroupModel.Name != "" {
		updateMap["name"] = in.GroupModel.Name
	}
	if in.GroupModel.Avatar != "" {
		updateMap["avatar"] = in.GroupModel.Avatar
	}
	if in.GroupModel.Description != "" {
		updateMap["description"] = in.GroupModel.Description
	}
	if in.GroupModel.AdminRemark != "" {
		updateMap["adminRemark"] = in.GroupModel.AdminRemark
	}
	updateMap["allMute"] = in.GroupModel.AllMute
	updateMap["allMuterType"] = in.GroupModel.AllMuterType
	updateMap["speakLimit"] = in.GroupModel.SpeakLimit
	updateMap["maxMember"] = in.GroupModel.MaxMember
	updateMap["memberCanStartTempChat"] = in.GroupModel.MemberCanStartTempChat
	updateMap["memberCanInviteFriend"] = in.GroupModel.MemberCanInviteFriend
	updateMap["newMemberHistoryMsgCount"] = in.GroupModel.NewMemberHistoryMsgCount
	updateMap["anonymousChat"] = in.GroupModel.AnonymousChat
	if len(updateMap) > 0 {
		err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			err := tx.Model(model).Where("id = ?", in.GroupModel.Id).Updates(updateMap).Error
			if err != nil {
				l.Errorf("更新失败: %v", err)
				return err
			}
			return nil
		}, func(tx *gorm.DB) error {
			notice := &noticemodel.Notice{
				ConvId: pb.HiddenConvIdGroup(in.GroupModel.Id),
				Options: noticemodel.NoticeOption{
					StorageForClient: false,
					UpdateConvNotice: false,
				},
				ContentType: pb.NoticeContentType_SetGroupInfo,
				Content: utils.AnyToBytes(pb.NoticeContent_DismissGroup{
					GroupId: in.GroupModel.Id,
				}),
				UniqueId: "info",
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
			err = groupmodel.CleanGroupCache(l.ctx, l.svcCtx.Redis(), in.GroupModel.Id)
			if err != nil {
				l.Errorf("CreateGroup CleanGroupCache error: %v", err)
				return err
			}
			return nil
		})
		if err != nil {
			return &pb.UpdateGroupModelResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		{
			// 刷新订阅
			utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
				groupmodel.CleanGroupCache(l.ctx, l.svcCtx.Redis(), in.GroupModel.Id)
				_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
					CommonReq: in.CommonReq,
					ConvId:    pb.HiddenConvIdGroup(in.GroupModel.Id),
				})
				if err != nil {
					l.Errorf("SendNoticeData failed, err: %v", err)
				}
				return err
			})
		}
	}
	return &pb.UpdateGroupModelResp{}, nil
}
