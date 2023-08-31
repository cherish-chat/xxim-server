package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type EditMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEditMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EditMsgLogic {
	return &EditMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EditMsg 编辑消息
func (l *EditMsgLogic) EditMsg(in *pb.EditMsgReq) (*pb.EditMsgResp, error) {
	getMsgByIdResp, err := NewGetMsgByIdLogic(l.ctx, l.svcCtx).GetMsgById(&pb.GetMsgByIdReq{
		ServerMsgId: utils.AnyPtr(in.ServerMsgId),
		ClientMsgId: nil,
		Push:        false,
		CommonReq:   in.GetCommonReq(),
	})
	if err != nil {
		l.Errorf("getMsgByIdLogic err: %v", err)
		return &pb.EditMsgResp{CommonResp: getMsgByIdResp.CommonResp}, err
	}
	convId, _ := pb.ParseConvServerMsgId(in.ServerMsgId)
	if getMsgByIdResp.GetMsgData().GetSenderId() != in.GetCommonReq().GetUserId() && getMsgByIdResp.GetMsgData().GetContentType() != 23 {
		if !pb.IsGroupConv(convId) {
			return &pb.EditMsgResp{CommonResp: pb.NewToastErrorResp("没有权限")}, nil
		}
		// 自己是什么身份
		canEditGroupMemberMsg, err := l.svcCtx.GroupService().CanEditGroupMemberMsg(l.ctx, &pb.CanEditGroupMemberMsgReq{
			CommonReq: in.CommonReq,
			GroupId:   pb.ParseGroupConv(convId),
			MemberId:  getMsgByIdResp.GetMsgData().GetSenderId(),
		})
		if err != nil {
			l.Errorf("CanEditGroupMemberMsg err: %v", err)
			return &pb.EditMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if canEditGroupMemberMsg.GetCommonResp().GetCode() != pb.CommonResp_Success {
			return &pb.EditMsgResp{CommonResp: canEditGroupMemberMsg.GetCommonResp()}, nil
		}
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		return tx.Model(&msgmodel.Msg{}).Table(msgmodel.GetMsgTableNameById(in.ServerMsgId)).Where("id = ?", in.ServerMsgId).Updates(map[string]interface{}{
			"contentType": in.ContentType,
			"content":     in.Content,
			"ext":         in.Ext,
		}).Error
	}, func(tx *gorm.DB) error {
		// 通知
		msgData := &pb.MsgData{
			ClientMsgId: getMsgByIdResp.GetMsgData().GetClientMsgId(),
			ServerMsgId: getMsgByIdResp.GetMsgData().GetServerMsgId(),
			ClientTime:  getMsgByIdResp.GetMsgData().GetClientTime(),
			ServerTime:  getMsgByIdResp.GetMsgData().GetServerTime(),
			SenderId:    getMsgByIdResp.GetMsgData().GetSenderId(),
			SenderInfo:  getMsgByIdResp.GetMsgData().GetSenderInfo(),
			ConvId:      getMsgByIdResp.GetMsgData().GetConvId(),
			AtUsers:     getMsgByIdResp.GetMsgData().GetAtUsers(),
			ContentType: in.ContentType,
			Content:     in.Content,
			Seq:         getMsgByIdResp.GetMsgData().GetSeq(),
			Options:     getMsgByIdResp.GetMsgData().GetOptions(),
			OfflinePush: getMsgByIdResp.GetMsgData().GetOfflinePush(),
			Ext:         in.Ext,
		}
		buf, _ := proto.Marshal(msgData)
		l.Infof("edit msg: %v %v", msgData.ContentType, string(msgData.Content))
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvId(convId),
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvNotice: false,
			},
			ContentType: int32(pb.NoticeType_EDIT),
			Content:     buf,
			UniqueId:    in.ServerMsgId,
			Title:       "",
			Ext:         nil,
		}
		return notice.Insert(l.ctx, tx, l.svcCtx.Redis())
	}, func(tx *gorm.DB) error {
		// 删除消息缓存
		return msgmodel.FlushMsgCache(l.ctx, l.svcCtx.Redis(), []string{in.ServerMsgId})
	})
	if err != nil {
		l.Errorf("edit msg failed, err: %v", err)
		return &pb.EditMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 通知
	xtrace.StartFuncSpan(l.ctx, "SendNotice", func(ctx context.Context) {
		utils.RetryProxy(ctx, 12, time.Second, func() error {
			_, err := l.svcCtx.NoticeService().GetUserNoticeData(ctx, &pb.GetUserNoticeDataReq{
				CommonReq: in.CommonReq,
				UserId:    "",
				ConvId:    pb.HiddenConvId(convId),
			})
			if err != nil {
				l.Errorf("ApplyToBeGroupMember SendNoticeData error: %v", err)
				return err
			}
			return nil
		})
	})
	return &pb.EditMsgResp{}, nil
}
