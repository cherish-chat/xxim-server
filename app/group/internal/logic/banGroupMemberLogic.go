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

type BanGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBanGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanGroupMemberLogic {
	return &BanGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BanGroupMember 禁言群成员
func (l *BanGroupMemberLogic) BanGroupMember(in *pb.BanGroupMemberReq) (*pb.BanGroupMemberResp, error) {
	if in.UnbanTime < time.Now().UnixMilli()+1000*60 {
		return &pb.BanGroupMemberResp{CommonResp: pb.NewAlertErrorResp(
			l.svcCtx.T(in.CommonReq.Language, "操作失败"),
			l.svcCtx.T(in.CommonReq.Language, "解禁时间不能小于一分钟"),
		)}, nil
	}
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
			return &pb.BanGroupMemberResp{CommonResp: resp.CommonResp}, err
		}
		value, ok := resp.GroupMap[in.GroupId]
		if !ok {
			return &pb.BanGroupMemberResp{CommonResp: pb.NewAlertErrorResp(
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
			return &pb.BanGroupMemberResp{CommonResp: resp.CommonResp}, err
		}
		myRole = resp.GroupMemberInfo.Role
	}

	// 判断是否有权限
	// 如果是member 只能修改自己的信息
	if myRole == pb.GroupRole_MEMBER && in.MemberId != in.CommonReq.UserId {
		return &pb.BanGroupMemberResp{
			CommonResp: pb.NewAlertErrorResp(
				l.svcCtx.T(in.CommonReq.Language, "操作失败"),
				l.svcCtx.T(in.CommonReq.Language, "没有权限"),
			),
		}, nil
	}
	// 如果是管理员 不能修改群主的信息
	if myRole == pb.GroupRole_MANAGER && group.Owner == in.MemberId {
		return &pb.BanGroupMemberResp{
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
			return &pb.BanGroupMemberResp{
				CommonResp: pb.NewRetryErrorResp(),
			}, err
		}
		updateMap := map[string]interface{}{}
		updateMap["unbanTime"] = in.UnbanTime
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
			notice := &noticemodel.Notice{
				ConvId: pb.HiddenConvIdGroup(in.GroupId),
				Options: noticemodel.NoticeOption{
					StorageForClient: false,
					UpdateConvNotice: false,
				},
				ContentType: pb.NoticeContentType_SetGroupMemberInfo,
				Content: utils.AnyToBytes(pb.NoticeContent_SetGroupMemberInfo{
					GroupId:   in.GroupId,
					MemberId:  in.MemberId,
					UpdateMap: updateMap,
				}),
				Title:    "",
				UniqueId: "member",
				Ext:      nil,
			}
			err = notice.Insert(l.ctx, tx, l.svcCtx.Redis())
			if err != nil {
				l.Errorf("insert notice failed, err: %v", err)
				return err
			}
			return nil
		})
		if err != nil {
			l.Errorf("Transaction err: %v", err)
			return &pb.BanGroupMemberResp{
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
			_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
				CommonReq: in.CommonReq,
				ConvId:    pb.HiddenConvIdGroup(in.GroupId),
			})
			if err != nil {
				l.Errorf("SendNoticeData failed, err: %v", err)
			}
			return nil
		})
	}
	return &pb.BanGroupMemberResp{}, nil
}
