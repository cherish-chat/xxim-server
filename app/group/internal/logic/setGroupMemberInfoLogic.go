package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupMemberInfoLogic {
	return &SetGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetGroupMemberInfo 设置群成员信息
func (l *SetGroupMemberInfoLogic) SetGroupMemberInfo(in *pb.SetGroupMemberInfoReq) (*pb.SetGroupMemberInfoResp, error) {
	// 获取群信息
	var group *groupmodel.Group
	{
		var resp *pb.MapGroupByIdsResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "mapGroupByIds", func(ctx context.Context) {
			resp, err = NewMapGroupByIdsLogic(ctx, l.svcCtx).MapGroupByIds(&pb.MapGroupByIdsReq{
				CommonReq: in.CommonReq,
				Ids:       []string{in.GroupId},
			})
		})
		if err != nil {
			l.Errorf("getGroupMemberInfoLogic err: %v", err)
			return &pb.SetGroupMemberInfoResp{CommonResp: resp.CommonResp}, err
		}
		value, ok := resp.GroupMap[in.GroupId]
		if !ok {
			return &pb.SetGroupMemberInfoResp{CommonResp: pb.NewAlertErrorResp(
				l.svcCtx.T(in.CommonReq.Language, "操作失败"),
				l.svcCtx.T(in.CommonReq.Language, "群聊不存在"),
			)}, nil
		}
		group = groupmodel.GroupFromBytes(value)
	}
	// 获取自己的群成员信息 看是否有权
	myRole := pb.GroupRole_MEMBER
	{
		var resp *pb.GetGroupMemberInfoResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "getGroupMemberInfoLogic", func(ctx context.Context) {
			resp, err = NewGetGroupMemberInfoLogic(ctx, l.svcCtx).GetGroupMemberInfo(&pb.GetGroupMemberInfoReq{
				CommonReq: in.CommonReq,
				GroupId:   in.GroupId,
				MemberId:  in.CommonReq.UserId,
			})
		})
		if err != nil {
			l.Errorf("getGroupMemberInfoLogic err: %v", err)
			return &pb.SetGroupMemberInfoResp{CommonResp: resp.CommonResp}, err
		}
		myRole = resp.GroupMemberInfo.Role
	}

	// 判断是否有权限
	// 如果是member 只能修改自己的信息
	if myRole == pb.GroupRole_MEMBER && in.MemberId != in.CommonReq.UserId {
		return &pb.SetGroupMemberInfoResp{
			CommonResp: pb.NewAlertErrorResp(
				l.svcCtx.T(in.CommonReq.Language, "操作失败"),
				l.svcCtx.T(in.CommonReq.Language, "没有权限"),
			),
		}, nil
	}
	// 如果是管理员 不能修改群主的信息
	if myRole == pb.GroupRole_MANAGER && group.Owner == in.MemberId {
		return &pb.SetGroupMemberInfoResp{
			CommonResp: pb.NewAlertErrorResp(
				l.svcCtx.T(in.CommonReq.Language, "操作失败"),
				l.svcCtx.T(in.CommonReq.Language, "没有权限"),
			),
		}, nil
	}
	// 可以修改
	{
		// 先清除缓存
		err := groupmodel.FlushGroupMemberCache(l.ctx, l.svcCtx.Redis(), in.GroupId, in.MemberId)
		if err != nil {
			l.Errorf("FlushGroupMemberCache err: %v", err)
			return &pb.SetGroupMemberInfoResp{
				CommonResp: pb.NewRetryErrorResp(),
			}, err
		}
		updateMap := map[string]interface{}{}
		if in.Role != nil {
			updateMap["role"] = *in.Role
		}
		if in.Remark != nil {
			updateMap["remark"] = *in.Remark
		}
		if in.UnbanTime != nil {
			updateMap["unbanTime"] = *in.UnbanTime
		}
		if in.GroupRemark != nil {
			updateMap["groupRemark"] = *in.GroupRemark
		}
		// 再修改数据库
		err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			// 修改数据库
			if len(updateMap) > 0 {
				err := xorm.Update(tx, &groupmodel.GroupMember{}, updateMap, xorm.Where("groupId = ? and userId = ?", in.GroupId, in.MemberId))
				if err != nil {
					l.Errorf("Update err: %v", err)
					return err
				}
			}
			return nil
		}, func(tx *gorm.DB) error {
			// 通知member
			// 发送一条订阅号消息 订阅号的convId = notice:group@groupId  noticeId = memberId
			data := &pb.NoticeData{
				ConvId:         noticemodel.ConvIdGroup(group.Id),
				UnreadCount:    0,
				UnreadAbsolute: false,
				NoticeId:       noticemodel.NoticeIdUpdateMemberInfo(in.MemberId),
				ContentType:    0,
				Content:        []byte(in.Notice),
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
			l.Errorf("Transaction err: %v", err)
			return &pb.SetGroupMemberInfoResp{
				CommonResp: pb.NewRetryErrorResp(),
			}, err
		}
	}
	// 后续操作
	{
		utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
			// 删除缓存
			err := groupmodel.FlushGroupMemberCache(l.ctx, l.svcCtx.Redis(), in.GroupId, in.MemberId)
			if err != nil {
				l.Errorf("FlushGroupMemberCache err: %v", err)
				return err
			}
			// 预热缓存
			go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarmUp", func(ctx context.Context) {
				_, err := groupmodel.ListGroupMemberFromMysql(ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), in.GroupId, []string{in.MemberId})
				if err != nil {
					l.Errorf("ListGroupMemberFromMysql err: %v", err)
				}
			}, propagation.MapCarrier{
				"group_id":  group.Id,
				"member_id": in.MemberId,
			})
			// 发送通知
			_, err = l.svcCtx.NoticeService().SendNoticeData(l.ctx, &pb.SendNoticeDataReq{
				CommonReq: in.CommonReq,
				NoticeData: &pb.NoticeData{
					NoticeId: noticemodel.NoticeIdUpdateMemberInfo(in.MemberId),
					ConvId:   noticemodel.ConvIdGroup(group.Id),
				},
				UserId:      nil,
				IsBroadcast: utils.AnyPtr(true),
				Inserted:    utils.AnyPtr(true),
			})
			if err != nil {
				l.Errorf("SendNoticeData failed, err: %v", err)
			}
			return nil
		})
	}
	return &pb.SetGroupMemberInfoResp{}, nil
}
